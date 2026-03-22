package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

// SMTPMailer реализует output.Mailer через SMTP.
type SMTPMailer struct {
	host                   string
	port                   int
	skipTLSVerify          bool
	authEnabled            bool
	username               string
	password               string
	from                   string
	baseURL                string
	resetTTL               time.Duration
	fallbackTimeoutSeconds int
}

type SMTPConfig struct {
	Host                   string
	Port                   int
	SkipTLSVerify          bool
	AuthEnabled            bool
	Username               string
	Password               string
	From                   string
	BaseURL                string
	ResetTTL               time.Duration
	FallbackTimeoutSeconds int
}

// NewSMTPMailer создаёт новый SMTP mailer.
func NewSMTPMailer(cfg SMTPConfig) *SMTPMailer {
	return &SMTPMailer{
		host:                   cfg.Host,
		port:                   cfg.Port,
		skipTLSVerify:          cfg.SkipTLSVerify,
		authEnabled:            cfg.AuthEnabled,
		username:               cfg.Username,
		password:               cfg.Password,
		from:                   cfg.From,
		baseURL:                cfg.BaseURL,
		resetTTL:               cfg.ResetTTL,
		fallbackTimeoutSeconds: cfg.FallbackTimeoutSeconds,
	}
}

// SendPasswordResetEmail отправляет письмо с ссылкой для сброса пароля.
// resetToken — сырой (не захешированный) токен, который вставляется в ссылку.
// ctx используется для отмены на этапе dial.
// Для последующих SMTP-шагов на соединение выставляется deadline:
// если у ctx есть deadline — используется он,
// иначе применяется fallback timeout из конфига.
func (m *SMTPMailer) SendPasswordResetEmail(ctx context.Context, toEmail string, resetToken string) error {
	resetLink := fmt.Sprintf("%s?token=%s", m.baseURL, resetToken)
	ttlMinutes := int(m.resetTTL.Minutes())

	subject := "Password Reset Request"
	body := fmt.Sprintf(
		"You have requested a password reset.\r\n\r\nClick the link below to reset your password:\r\n%s\r\n\r\nThis link expires in %d minutes.\r\n\r\nIf you did not request a password reset, please ignore this email.",
		resetLink,
		ttlMinutes,
	)
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		m.from, toEmail, subject, body,
	)

	return m.send(ctx, toEmail, msg)
}

// send выполняет весь SMTP-диалог: dial → STARTTLS → AUTH → DATA → Quit.
func (m *SMTPMailer) send(ctx context.Context, toEmail, msg string) error {
	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}

	fallback := time.Now().Add(time.Duration(m.fallbackTimeoutSeconds) * time.Second)
	deadline := fallback
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	if err := conn.SetDeadline(deadline); err != nil {
		_ = conn.Close()
		return fmt.Errorf("smtp set deadline: %w", err)
	}

	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		// conn может остаться открытым если NewClient упал — закрываем явно
		_ = conn.Close()
		return fmt.Errorf("smtp new client: %w", err)
	}

	// Quit() уже закрывает соединение нормально.
	// defer Close() нужен только как fallback при ошибке до Quit.
	closed := false
	defer func() {
		if !closed {
			_ = client.Close()
		}
	}()

	// STARTTLS — используется если сервер объявляет поддержку расширения.
	// tls.Config с ServerName необходим для корректного TLS-handshake.
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsCfg := &tls.Config{
			ServerName:         m.host,
			InsecureSkipVerify: m.skipTLSVerify,
		}
		if err := client.StartTLS(tlsCfg); err != nil {
			return fmt.Errorf("smtp starttls: %w", err)
		}
	}

	// AUTH — только если явно включён в конфиге
	if m.authEnabled {
		auth := smtp.PlainAuth("", m.username, m.password, m.host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	if err := client.Mail(m.from); err != nil {
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
