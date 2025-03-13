package infrastructure

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type SMTPClient struct {
	cfg Config
}

func NewSMTPInstance(cfg *Config) *SMTPClient {
	return &SMTPClient{
		cfg: *cfg,
	}
}

func (s *SMTPClient) SendEmail(ctx context.Context, senderEmail, senderName string, recipientEmail []string, subject, body string) error {
	// set up the smtp server
	smtpUsername := s.cfg.SMTP.Username
	smtpPassword := s.cfg.SMTP.Password
	smtpHost := s.cfg.SMTP.Host
	smtpPort := s.cfg.SMTP.Port

	var smtpAuth smtp.Auth
	if s.cfg.Mode == "development" {
		smtpAuth = nil
	} else {
		smtpAuth = smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	}

	// Compose message with MIME headers
	message := "From: " + senderName + " <" + senderEmail + ">\r\n" +
		"To: " + strings.Join(recipientEmail, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + body

	// Send email
	smtpAddr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	if err := smtp.SendMail(smtpAddr, smtpAuth, senderEmail, recipientEmail, []byte(message)); err != nil {
		log.Printf("failed to send email with error: %v\n", err)
		return err
	}

	return nil
}
