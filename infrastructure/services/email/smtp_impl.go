package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/ZaphCode/clean-arch/config"
)

type smtpEmailServiceImpl struct{}

func NewSmtpEmailService() EmailService {
	return &smtpEmailServiceImpl{}
}

func (s *smtpEmailServiceImpl) SendEmail(data EmailData) error {
	cfg := config.Get()

	subject := fmt.Sprintf("Subject: %s!\n", data.Subject)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	t, err := template.ParseFiles(
		"./infrastructure/services/email/templates/" + data.Template,
	)

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	if err = t.Execute(buf, data.Data); err != nil {
		return err
	}

	body := buf.String()

	msg := []byte(subject + mime + body)

	sender := cfg.Smtp.Email

	auth := smtp.PlainAuth("", sender, "tgruuykixcvtnzvj", "smtp.gmail.com")

	recibers := []string{data.Email}

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		sender,
		recibers,
		[]byte(msg),
	)

	if err != nil {
		return err
	}

	return nil
}
