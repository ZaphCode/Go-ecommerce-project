package email

type EmailService interface {
	SendEmail(data EmailData) error
}

type EmailData struct {
	Email    string
	Template string
	Subject  string
	Data     map[string]any
}
