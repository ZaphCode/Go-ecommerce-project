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
	t, err := template.ParseFiles("./templates/" + data.Template)

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

	host := cfg.Smtp.Host

	auth := smtp.PlainAuth("", sender, cfg.Smtp.Password, host)

	recibers := []string{data.Email}

	err = smtp.SendMail(
		host+":"+cfg.Api.Port,
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
