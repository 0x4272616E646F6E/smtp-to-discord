package shared

import (
	"fmt"

	"github.com/0x4272616E646F6E/smtp-to-discord/internal/smtp"
)

type DiscordCard struct {
	Subject string             `json:"subject"`
	From    string             `json:"from"`
	Fields  []DiscordCardField `json:"fields,omitempty"`
}

type DiscordCardField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func SMTPToDiscordCard(msg smtp.Message) DiscordCard {
	subject := msg.Subject
	body := msg.Data
	if subject == "" && msg.Data != "" {
		lines := splitLines(msg.Data)
		for i, line := range lines {
			if len(line) >= 8 && line[:8] == "Subject:" {
				subject = line[8:]
				subject = trimSpace(subject)
				lines = append(lines[:i], lines[i+1:]...)
				break
			}
		}
		body = joinLines(lines)
	}
	return DiscordCard{
		Subject: subject,
		From:    msg.From,
		Fields: []DiscordCardField{
			{Name: "Body", Value: body, Inline: false},
		},
	}
}

func splitLines(s string) []string {
	var lines []string
	l := ""
	for i := 0; i < len(s); i++ {
		if s[i] == '\r' {
			if i+1 < len(s) && s[i+1] == '\n' {
				lines = append(lines, l)
				l = ""
				i++
			} else {
				lines = append(lines, l)
				l = ""
			}
		} else if s[i] == '\n' {
			lines = append(lines, l)
			l = ""
		} else {
			l += string(s[i])
		}
	}
	if l != "" {
		lines = append(lines, l)
	}
	return lines
}

func joinLines(lines []string) string {
	result := ""
	for i, l := range lines {
		if i > 0 {
			result += "\n"
		}
		result += l
	}
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

func FormatDiscordCardAsContent(card DiscordCard) string {
	content := fmt.Sprintf("**%s**\n%s\n", card.Subject, card.From)
	for _, f := range card.Fields {
		content += fmt.Sprintf("**%s:** %s\n", f.Name, f.Value)
	}
	return content
}
