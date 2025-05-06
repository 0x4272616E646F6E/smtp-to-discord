FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o smtp-to-discord ./cmd/main.go

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/smtp-to-discord .
EXPOSE 25
ENV DISCORD_WEBHOOK_URL=""
CMD ["/app/smtp-to-discord"]
