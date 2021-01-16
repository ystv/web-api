package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

// Mailer encapsulates the dependency
type Mailer struct {
	*mail.SMTPClient
	Enabled bool
}

// Config represents a configuration to connect to an SMTP server
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
}

// Mail represents an email to be sent
type Mail struct {
	Subject string
	To      string
	Cc      []string
	Bcc     []string
	From    string
	tpl     template.Template
	tplData interface{}
}

// NewMailer creates a new SMTP client
func NewMailer(config Config) (*Mailer, error) {
	smtpServer := mail.SMTPServer{
		Host:           config.Host,
		Port:           config.Port,
		Username:       config.Username,
		Password:       config.Password,
		Encryption:     mail.EncryptionTLS,
		Authentication: mail.AuthPlain,
		ConnectTimeout: 10 * time.Second,
		SendTimeout:    10 * time.Second,
		TLSConfig:      &tls.Config{InsecureSkipVerify: true},
	}

	smtpClient, err := smtpServer.Connect()
	if err != nil {
		return &Mailer{nil, false}, err
	}
	return &Mailer{smtpClient, true}, nil
}

// SendMail sends a template email
func (m *Mailer) SendMail(item Mail) error {
	body := bytes.Buffer{}
	err := item.tpl.Execute(&body, item.tplData)
	if err != nil {
		return fmt.Errorf("failed to exec tpl: %w", err)
	}
	email := mail.NewMSG()
	email.SetFrom(item.From).AddTo(item.To).SetSubject(item.Subject)
	if len(item.Cc) != 0 {
		email.AddCc(item.Cc...)
	}
	if len(item.Bcc) != 0 {
		email.AddBcc(item.Bcc...)
	}
	email.SetBody(mail.TextHTML, body.String())
	if email.Error != nil {
		return fmt.Errorf("failed to set mail data: %w", email.Error)
	}
	return email.Send(m.SMTPClient)
}
