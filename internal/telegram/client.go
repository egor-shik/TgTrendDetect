package telegram

import (
 "context"
 "fmt"
 "strings"

 "github.com/egor-shik/TgTrendDetect/internal/config"

 "github.com/gotd/td/session"
 "github.com/gotd/td/telegram"
 "github.com/gotd/td/telegram/auth"
 "github.com/gotd/td/tg"
)

type Client struct {
 cfg      *config.Config
 tgClient *telegram.Client
}

func NewClient(cfg *config.Config) (*Client, error) {
 sessionStorage := &session.FileStorage{
  Path: cfg.SessionFilePath,
 }
 rawClient := telegram.NewClient(cfg.AppID, cfg.AppHash, telegram.Options{
  SessionStorage: sessionStorage,
 })
 return &Client{
  cfg:      cfg,
  tgClient: rawClient,
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
  
	  fmt.Println("successfully connected to Telegram!")
  
	  <-ctx.Done()
	  return ctx.Err()
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