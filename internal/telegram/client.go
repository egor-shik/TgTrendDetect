package telegram

import (
 "context"
 "fmt"
 "strings"
 "time" 

 "github.com/egor-shik/TgTrendDetect/internal/config"

 "github.com/gotd/td/session"
 "github.com/gotd/td/telegram"
 "github.com/gotd/td/telegram/auth"
 "github.com/gotd/td/telegram/updates"
 "github.com/gotd/td/tg"
)

type Client struct {
 cfg      *config.Config
 tgClient *telegram.Client
 gaps     *updates.Manager
}

func NewClient(cfg *config.Config, gaps *updates.Manager) (*Client, error) {
 sessionStorage := &session.FileStorage{
  Path: cfg.SessionFilePath,
 }
 rawClient := telegram.NewClient(cfg.AppID, cfg.AppHash, telegram.Options{
  SessionStorage: sessionStorage,
  UpdateHandler:  gaps,
 })
 return &Client{
  cfg:      cfg,
  tgClient: rawClient,
  gaps:     gaps,
 }, nil
}

func (c *Client) Start(ctx context.Context) error {
 flow := auth.NewFlow(
  auth.Constant(c.cfg.Phone, "", auth.CodeAuthenticatorFunc(promtCode)),
  auth.SendCodeOptions{},
 )

 return c.tgClient.Run(ctx, func(ctx context.Context) error {
  if err := flow.Run(ctx, c.tgClient.Auth()); err != nil {
   return fmt.Errorf("Authorization error: %w", err)
  }

  self, err := c.tgClient.Self(ctx)
  if err != nil {
   return fmt.Errorf("failed to get self: %w", err)
  }

  fmt.Printf("Logged in as %s (ID: %d)\n", self.FirstName, self.ID)
  fmt.Println("Warming up channel dialogs...")

  api := c.tgClient.API()

  var chats []tg.ChatClass 

  dialogs, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
   OffsetPeer: &tg.InputPeerEmpty{},
   Limit:      100,
  })
  if err != nil {
   fmt.Printf("Warning: failed to fetch dialogs: %v\n", err)
  } else {
   switch d := dialogs.(type) {
   case *tg.MessagesDialogs:
    chats = d.Chats
   case *tg.MessagesDialogsSlice:
    chats = d.Chats
   }
   fmt.Printf("SUCCESS! Loaded %d dialogs/channels into cache!\n", len(chats))
  }

  var targetChannels []*tg.Channel
  for _, chat := range chats {
   if ch, ok := chat.(*tg.Channel); ok {
    targetChannels = append(targetChannels, ch)
   }
  }
  fmt.Printf("SUCCESS! Target channels for polling: %d\n", len(targetChannels))

  poller := NewChannelPoller(c)
  go poller.StartPolling(ctx, targetChannels, 10*time.Second)

  fmt.Println("Starting Gaps synchronization... Listening for posts!")

  return c.gaps.Run(ctx, api, self.ID, updates.AuthOptions{
   IsBot: false,
  })
 })
}

func promtCode(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
 fmt.Print("Enter the code from Telegram: ")
 var code string
 _, err := fmt.Scan(&code)
 return strings.TrimSpace(code), err
}

func promptPassword(ctx context.Context) (string, error) {
 fmt.Print("Enter your 2FA Cloud Password: ")
 var password string
 _, err := fmt.Scan(&password)
 return strings.TrimSpace(password), err
}