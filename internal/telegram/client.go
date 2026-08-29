package telegram

import (
 "context"
 "fmt"
 "strings"

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
  auth.Constant(c.cfg.Phone, "pipu", auth.CodeAuthenticatorFunc(promtCode)),
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
  if _, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
   Limit: 100,
  }); err != nil {
   fmt.Printf("Warning: failed to fetch dialogs: %v\n", err)
  }

  fmt.Println("Starting Gaps synchronization... Listening for posts from ALL channels!")

  return c.gaps.Run(ctx, api, self.ID, updates.AuthOptions{})
}
}

func promtCode(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
 fmt.Print("Enter the code from Telegram: ")
 var code string
 _, err := fmt.Scan(&code)
 return strings.TrimSpace(code), err
}