package dtos

import "github.com/ZaphCode/clean-arch/src/domain"

type None *struct{}

type SignInResp struct {
	User         *domain.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}
