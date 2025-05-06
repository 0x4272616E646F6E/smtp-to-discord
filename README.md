# smtp-to-discord

`smtp-to-discord` is a lightweight, containerized Go application that runs a minimal SMTP server and relays received email messages to a Discord channel via webhook. Ideal for infrastructure alerts, internal system notifications, or bridging legacy SMTP tools into modern workflows.

---

## ✨ Features

- Minimal, custom SMTP server (port 25)
- Parses `From`, `To`, `Subject`, and body
- Converts messages to Discord embeds
- Concurrent processing via worker queue
- Graceful shutdown support
- Docker-ready for deployment

---

## 🚀 Getting Started

### 🔧 Prerequisites

- Go 1.20+
- Docker (optional for container use)
- A Discord Webhook URL ([how to get one](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks))

---

### 🧪 Run Locally

```sh
git clone <your-repo-url>
cd smtp-to-discord
export DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/your_webhook_id/your_webhook_token"
go run ./cmd/main.go

### Running with Docker

1. **Build the Docker image:**
    ```sh
    docker build -t smtp-to-discord .
    ```
2. **Run the container:**
    ```sh
    docker run -e DISCORD_WEBHOOK_URL="<your_webhook_url>" -p 25:25 smtp-to-discord
    ```

> **Note:** Binding to port 25 may require root privileges or port remapping on some systems. For local testing, you can map another port (e.g., `-p 2525:25`) and connect your SMTP client to that port.

### Sending a Test Email
You can use `swaks` or any SMTP client to send a test message:
```sh
swaks --to test@example.com --from you@example.com --server localhost:25 --header "Subject: Hello Discord" --body "This is a test."
```

## Graceful Shutdown
The application supports graceful shutdown on SIGINT/SIGTERM. It will stop accepting new SMTP connections and finish processing any queued messages before exiting.

## Project Structure

- `cmd/main.go` — Application entry point, queue, worker pool, and graceful shutdown
- `internal/smtp/server.go` — Minimal SMTP server implementation
- `internal/shared/queue.go` — Thread-safe message queue
- `internal/shared/parsing.go` — SMTP-to-Discord parsing logic
- `internal/discord/client.go` — Discord webhook integration

## License
[MIT](LICENSE)
