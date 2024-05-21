package mailer

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/mail"
	"net/textproto"

	gomail "github.com/wneessen/go-mail"
)

var _ Mailer = (*SmtpClient)(nil)

type SmtpAuthType string

const (
	SmtpAuthPlain   SmtpAuthType = "PLAIN"
	SmtpAuthLogin   SmtpAuthType = "LOGIN"
	SmtpAuthCramMD5 SmtpAuthType = "CRAM-MD5"
	SmtpXOauth2     SmtpAuthType = "XOAUTH2"
)

type SmtpClient struct {
	Log        *slog.Logger
	Host       string
	Port       int
	Username   string
	Password   string
	TLS        bool
	AuthMethod SmtpAuthType
}

func getAuthType(authType SmtpAuthType) gomail.SMTPAuthType {
	switch authType {
	case SmtpAuthPlain:
		return gomail.SMTPAuthPlain
	case SmtpAuthLogin:
		return gomail.SMTPAuthLogin
	case SmtpAuthCramMD5:
		return gomail.SMTPAuthCramMD5
	case SmtpXOauth2:
		return gomail.SMTPAuthXOAUTH2
	default:
		return gomail.SMTPAuthPlain
	}
}

func formatAddresses(addresses []mail.Address, withName bool) []string {
	result := make([]string, len(addresses))

	for i, addr := range addresses {
		if withName && addr.Name != "" {
			result[i] = addr.String()
		} else {
			result[i] = addr.Address
		}
	}

	return result
}

func (client *SmtpClient) Send(ctx context.Context, msg *Message) error {
	if client.Log == nil {
		client.Log = slog.Default()
	}

	opts := []gomail.Option{
		gomail.WithPort(client.Port),
		gomail.WithUsername(client.Username),
		gomail.WithPassword(client.Password),
		gomail.WithSMTPAuth(getAuthType(client.AuthMethod)),
		gomail.WithTLSPortPolicy(gomail.TLSMandatory),
		gomail.WithSSLPort(true),
	}

	c, err := gomail.NewClient(client.Host, opts...)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	m := gomail.NewMsg()

	if err := m.From(msg.From.Address); err != nil {
		client.Log.Error("invalid `from` header: %w", err)
		return err
	}
	if err := m.To(formatAddresses(msg.To, true)...); err != nil {
		client.Log.Error("invalid `to` header: %w", err)
		return err
	}
	if err := m.Cc(formatAddresses(msg.Cc, true)...); err != nil {
		client.Log.Error("invalid `Cc` header: %w", err)
		return err
	}
	if err := m.Bcc(formatAddresses(msg.Bcc, true)...); err != nil {
		client.Log.Error("invalid `Bcc` header: %w", err)
		return err
	}

	m.Subject(msg.Subject)
	m.SetBodyString(gomail.TypeTextHTML, msg.HTML)
	m.AddAlternativeString(gomail.TypeTextPlain, msg.Text)
	m.SetGenHeader(gomail.HeaderReplyTo, msg.ReplyTo)

	for k, v := range msg.Headers {
		m.SetGenHeader(gomail.Header(k), v)
	}

	attachments := make([]*gomail.File, len(msg.Attachments))
	for i, att := range msg.Attachments {
		attachments[i] = &gomail.File{
			Name:        att.FileName,
			Desc:        att.Description,
			ContentType: gomail.ContentType(att.ContentType),
			Header:      textproto.MIMEHeader{},
			Enc:         gomail.EncodingQP,
			Writer: func(w io.Writer) (int64, error) {
				return io.Copy(w, att.Content)
			},
		}
	}
	if len(attachments) > 0 {
		m.SetAttachements(attachments)
	}

	if err := c.DialAndSendWithContext(ctx, m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
