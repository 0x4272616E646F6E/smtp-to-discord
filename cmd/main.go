package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/0x4272616E646F6E/smtp-to-discord/internal/discord"
	"github.com/0x4272616E646F6E/smtp-to-discord/internal/smtp"
)

func main() {
        webhook := os.Getenv("DISCORD_WEBHOOK_URL")
        if webhook == "" {
                log.Fatal("DISCORD_WEBHOOK_URL environment variable not set")
        }

        queue := make(chan smtp.Message, 100)
        var accepting = true
        var wg sync.WaitGroup

        // Handler enqueues messages if still accepting
        handler := func(msg smtp.Message) {
                if accepting {
                        queue <- msg
                }
        }

        // Start a worker per CPU
        numWorkers := runtime.NumCPU()
        log.Printf("Starting %d worker(s) for message processing", numWorkers)

        for i := 0; i < numWorkers; i++ {
                wg.Add(1)
                go func(workerID int) {
                        defer wg.Done()
                        for msg := range queue {
                                err := discord.SendToDiscord(msg.Subject, msg.From, msg.Data)
                                if err != nil {
                                        log.Printf("[Worker %d] Failed to send to Discord: %v", workerID, err)
                                }
                        }
                }(i)
        }

        // Graceful shutdown
        ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
        defer stop()

        addr := ":25"
        serverErr := make(chan error, 1)
        go func() {
                log.Printf("Starting SMTP to Discord bridge on %s", addr)
                serverErr <- smtp.StartServer(addr, handler)
        }()

        select {
        case <-ctx.Done():
                log.Println("Shutdown signal received, finishing queued messages...")
                accepting = false
                close(queue) // stop all workers once queue is drained
                wg.Wait()
                log.Println("All messages processed. Exiting.")
        case err := <-serverErr:
                if err != nil {
                        log.Fatalf("SMTP server error: %v", err)
                }
        }
}