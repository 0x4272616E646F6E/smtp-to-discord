package discord

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type DiscordEmbed struct {
        Title       string `json:"title,omitempty"`
        Description string `json:"description,omitempty"`
}

func SendToDiscord(subject, from, body string) error {
        webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
        if webhookURL == "" {
                return nil
        }

        body = strings.TrimSpace(body)

        embed := DiscordEmbed{
                Title:       subject,
                Description: from + "\n\n" + body,
        }

        payload := map[string]interface{}{
                "embeds": []DiscordEmbed{embed},
        }

        b, _ := json.Marshal(payload)
        _, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(b))
        return err
}