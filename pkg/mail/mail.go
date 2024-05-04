package mail

import (
	"log"

	"github.com/wneessen/go-mail"
)

type Mailer struct {
	SmtpHost     string
	SmtpPort     int
	SmtpUser     string
	SmtpPassword string
}

func (mailer *Mailer) SendMail(to []string, subject string, message string) bool {
	m := mail.NewMsg()

	if err := m.From(mailer.SmtpUser); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}
	if err := m.To(to...); err != nil {
		log.Fatalf("failed to set To address: %s", err)
	}

	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, message)

	c, err := mail.NewClient(mailer.SmtpHost, mail.WithPort(mailer.SmtpPort), mail.WithSMTPAuth(mail.SMTPAuthLogin), mail.WithTLSPolicy(mail.TLSMandatory), mail.WithUsername(mailer.SmtpUser), mail.WithPassword(mailer.SmtpPassword))

	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}

	if err := c.DialAndSend(m); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}

	return true
}
