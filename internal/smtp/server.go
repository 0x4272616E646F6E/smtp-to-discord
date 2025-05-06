package smtp

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

type Message struct {
        From    string
        To      string
        Subject string
        Data    string
}

type Handler func(msg Message)

func StartServer(addr string, handler Handler) error {
        ln, err := net.Listen("tcp", addr)
        if err != nil {
                return err
        }
        log.Printf("SMTP server listening on %s", addr)
        for {
                conn, err := ln.Accept()
                if err != nil {
                        log.Printf("Failed to accept connection: %v", err)
                        continue
                }
                go handleConn(conn, handler)
        }
}

func handleConn(conn net.Conn, handler Handler) {
        defer conn.Close()

        reader := bufio.NewReader(conn)
        write := func(s string) {
                conn.Write([]byte(s + "\r\n"))
        }

        write("220 smtp-to-discord ESMTP Service Ready")

        var from, to, subject string
        var dataLines []string

        for {
                line, err := reader.ReadString('\n')
                if err != nil {
                        if err != io.EOF {
                                log.Printf("Read error: %v", err)
                        }
                        return
                }
                line = strings.TrimSpace(line)

                switch {
                case strings.HasPrefix(line, "HELO") || strings.HasPrefix(line, "EHLO"):
                        write("250 Hello")

                case strings.HasPrefix(line, "MAIL FROM:"):
                        from = trimSMTPField(strings.TrimPrefix(line, "MAIL FROM:"))
                        write("250 OK")

                case strings.HasPrefix(line, "RCPT TO:"):
                        to = trimSMTPField(strings.TrimPrefix(line, "RCPT TO:"))
                        write("250 OK")

                case strings.HasPrefix(line, "DATA"):
                        write("354 End data with <CR><LF>.<CR><LF>")
                        dataLines = nil

                        for {
                                dataLine, err := reader.ReadString('\n')
                                if err != nil {
                                        log.Printf("Error reading DATA: %v", err)
                                        return
                                }
                                dataLine = strings.TrimRight(dataLine, "\r\n")
                                if dataLine == "." {
                                        break
                                }
                                if strings.HasPrefix(strings.ToLower(dataLine), "subject:") && subject == "" {
                                        subject = trimSMTPField(dataLine[8:])
                                        continue
                                }
                                dataLines = append(dataLines, dataLine)
                        }

                        write("250 Message accepted")
                        handler(Message{
                                From:    from,
                                To:      to,
                                Subject: subject,
                                Data:    strings.Join(dataLines, "\n"),
                        })

                case strings.HasPrefix(line, "QUIT"):
                        write("221 Bye")
                        return

                default:
                        write("250 OK")
                }
        }
}

func trimSMTPField(s string) string {
        s = strings.TrimSpace(s)
        s = strings.Trim(s, "<>")
        return s
}