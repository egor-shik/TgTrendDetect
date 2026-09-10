package main 

import (
    "context"
    "log"
    "os/signal"
    "syscall"
    "fmt"
    "os"
    "github.com/egor-shik/TgTrendDetect/internal/config"
    "github.com/egor-shik/TgTrendDetect/internal/telegram"
)

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()
    cfg := config.NewConfig(
        pipo,
        "pipo",
        "pipo",
        "session.json", 
        )
    handler := telegram.NewUpdateHandler()
    client, err := telegram.NewClient(cfg, handler)
    if err != nil {
        log.Fatalf("Client creation error: %v", err)
    }
    fmt.Println("Starting application... Press ctrl+C to exit")

    if err := client.Start(ctx); err != nil && err != context.Canceled {
        log.Fatalf("Client exwecution error: %v", err)
    }
    fmt.Println("The app has shut down")
}