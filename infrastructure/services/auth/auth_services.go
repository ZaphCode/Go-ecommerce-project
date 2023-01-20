package auth

import (
	"time"

	"github.com/ZaphCode/clean-arch/domain"
	"github.com/google/uuid"
)

const (
	GoogleProvider  = "google"
	DiscordProvider = "discord"
	GithubProvider  = "github"
)

func GetOAuthProviders() []string {
	return []string{GoogleProvider, DiscordProvider, GithubProvider}
}

type Claims struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
}

//* Services

type JwtAuthService interface {
	CreateToken(claims Claims, exp time.Duration, secret string) (string, error)
	DecodeToken(jwtoken string, secret string) (*Claims, error)
}

type OAuthService interface {
	GetOAuthUser(code string) (*domain.User, error)
	GetOAuthUrl() string
}
