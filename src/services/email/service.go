package email

type EmailService interface {
	SendChangePasswordEmail(email, name, secretCode string) error
	SendVerifyEmail(email, name, secretCode string) error
}

type EmailData struct {
	Email    string
	Template string
	Subject  string
	Data     map[string]any
}
