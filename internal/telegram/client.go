package telegram 

import (
	tg "github.com/gotd/td/telegram"
	"TrendDet/internal/config"
	"github.com/gotd/td/session"
)

type Client struct {
	cfg *config.Config 
	tgClient *tg.Client
}

func NewCLient(cfg *config.Config) (*Client, error) {
	sessionStorage := &session.FileStorage{
		Path: cfg.SessionFilePath,
	}
	rawClient := tg.NewCLient(cfg.AppID, cfg.AppHash, tg.Options{
		SessionStorage: sessionStorage,
	})
	return &Client{
		cfg: cfg,
		tgClient: rawClient
	}, nil
}