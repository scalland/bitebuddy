package utils

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"net/mail"
)

func (u *Utils) IsValidMailAddress(address string) (string, bool) {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return "", false
	}
	return addr.Address, true
}

type SMTPEmail struct {
	SMTPServer     string
	SMTPServerPort int
	SMTPUsername   string
	SMTPPassword   string
	u              *Utils
}

func (u *Utils) NewSMTPEmail() *SMTPEmail {
	return &SMTPEmail{
		SMTPServer:     "",
		SMTPServerPort: 0,
		SMTPUsername:   "",
		SMTPPassword:   "",
		u:              u,
	}
}

func (u *Utils) NewSMTPEmailWithConfig(smtpPort int, smtpServer, smtpUser, smtpPass string) *SMTPEmail {
	return &SMTPEmail{
		SMTPServer:     smtpServer,
		SMTPServerPort: smtpPort,
		SMTPUsername:   smtpUser,
		SMTPPassword:   smtpPass,
		u:              u,
	}
}

func (s *SMTPEmail) Send(fromName, fromEmail, subject, htmlBody string, to, cc, bcc, attachmentPaths []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", fromName, fromEmail))
	m.SetHeader("To", to...)
	m.SetHeader("Cc", cc...)
	m.SetHeader("Bcc", bcc...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	for _, path := range attachmentPaths {
		m.Attach(path)
	}

	d := gomail.NewDialer(s.SMTPServer, s.SMTPServerPort, s.SMTPUsername, s.SMTPPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
