package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"

	"github.com/AlbinaKonovalova/auth-service/internal/config/modules"
)

// SMTPMailer реализует output.Mailer через SMTP.
type SMTPMailer struct {
	smtp          modules.SMTPConfig
	passwordReset modules.PasswordResetConfig
}

// NewSMTPMailer создаёт новый SMTP mailer.
func NewSMTPMailer(smtpCfg modules.SMTPConfig, passwordResetCfg modules.PasswordResetConfig) *SMTPMailer {
	return &SMTPMailer{
		smtp:          smtpCfg,
		passwordReset: passwordResetCfg,
	}
}

// SendPasswordResetEmail отправляет письмо с ссылкой для сброса пароля.
// resetToken — сырой (не захешированный) токен, который вставляется в ссылку.
func (m *SMTPMailer) SendPasswordResetEmail(ctx context.Context, toEmail string, resetToken string) error {
	resetLink := fmt.Sprintf("%s?token=%s", m.passwordReset.BaseURL, resetToken)
	ttlMinutes := int(m.passwordReset.TTL.Minutes())

	subject := "Password Reset Request"
	body := fmt.Sprintf(
		"You have requested a password reset.\r\n\r\nClick the link below to reset your password:\r\n%s\r\n\r\nThis link expires in %d minutes.\r\n\r\nIf you did not request a password reset, please ignore this email.",
		resetLink,
		ttlMinutes,
	)
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		m.smtp.From, toEmail, subject, body,
	)

	return m.send(ctx, toEmail, msg)
}

// send выполняет весь SMTP-диалог: dial → STARTTLS → AUTH → DATA → Quit.
func (m *SMTPMailer) send(ctx context.Context, toEmail, msg string) error {
	addr := fmt.Sprintf("%s:%d", m.smtp.Host, m.smtp.Port)

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}

	fallback := time.Now().Add(time.Duration(m.smtp.FallbackTimeoutSeconds) * time.Second)
	deadline := fallback
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	if err := conn.SetDeadline(deadline); err != nil {
		_ = conn.Close()
		return fmt.Errorf("smtp set deadline: %w", err)
	}

	client, err := smtp.NewClient(conn, m.smtp.Host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("smtp new client: %w", err)
	}

	closed := false
	defer func() {
		if !closed {
			_ = client.Close()
		}
	}()

	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsCfg := &tls.Config{
			ServerName:         m.smtp.Host,
			InsecureSkipVerify: m.smtp.SkipTLSVerify, //nolint:gosec
		}
		if err := client.StartTLS(tlsCfg); err != nil {
			return fmt.Errorf("smtp starttls: %w", err)
		}
	}

	if m.smtp.AuthEnabled {
		auth := smtp.PlainAuth("", m.smtp.Username, m.smtp.Password, m.smtp.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	if err := client.Mail(m.smtp.From); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err := client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("smtp rcpt to: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := fmt.Fprint(wc, msg); err != nil {
		return fmt.Errorf("smtp write body: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("smtp close body: %w", err)
	}

	if err := client.Quit(); err != nil {
		return fmt.Errorf("smtp quit: %w", err)
	}
	closed = true
	return nil
}
