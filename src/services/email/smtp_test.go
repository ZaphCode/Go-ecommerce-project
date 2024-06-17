package email

import (
	"testing"

	"github.com/ZaphCode/clean-arch/config"
)

func TestMain(m *testing.M) {
	config.MustLoadConfig("./../../../config")
	m.Run()
}

func Test_smtpEmailServiceImpl_SendEmail(t *testing.T) {
	type args struct {
		data EmailData
	}
	tests := []struct {
		name    string
		s       *smtpEmailServiceImpl
		args    args
		wantErr bool
	}{
		{
			name:    "Test SendEmail",
			s:       NewSmtpEmailService().(*smtpEmailServiceImpl),
			wantErr: false,
			args: args{
				data: EmailData{
					Email:    "om4r4rm4@gmail.com",
					Subject:  "Change Password Request",
					Template: "change_password.html",
					Data: map[string]interface{}{
						"Name": "Omar",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &smtpEmailServiceImpl{}
			if err := s.sendEmail(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("smtpEmailServiceImpl.SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// fqpt wsko rgke qmni new

//  /* "password": "tgru uyki xcvt nzvj", */ old

func Test_smtpEmailServiceImpl_SendChangePasswordEmail(t *testing.T) {
	type args struct {
		email      string
		name       string
		secretCode string
	}
	tests := []struct {
		name    string
		s       *smtpEmailServiceImpl
		args    args
		wantErr bool
	}{
		{
			name:    "Test SendChangePasswordEmail",
			s:       NewSmtpEmailService().(*smtpEmailServiceImpl),
			wantErr: false,
			args: args{
				email:      "rebornmcz2@gmail.com",
				name:       "Reborn",
				secretCode: "dfadf",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &smtpEmailServiceImpl{}
			if err := s.SendChangePasswordEmail(tt.args.email, tt.args.name, tt.args.secretCode); (err != nil) != tt.wantErr {
				t.Errorf("smtpEmailServiceImpl.SendChangePasswordEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
