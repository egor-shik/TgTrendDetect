package main 

import (
    "context"
    "log"
    "os/signal"
    "syscall"
    "fmt"
    "os"
    "TrendDet/internal/config"
    "TrendDet/internal/telegram"
)

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop
    cfg := &config.Config{
    AppID:  
    AppHash:  
    Phone: 
    SessionFilePath:  
    }
    client, err := telegram.NewClient(cfg)
    if err != nil {
        log.Fatalf("Client creation error: %v", err)
    }
    fmt.Println("The app has shut down")
}