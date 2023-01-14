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
	CreateToken(Claims, time.Duration, bool) (string, error)
	DecodeToken(string, bool) (*Claims, error)
	CreateTokens(Claims, time.Duration, time.Duration) (string, string, error)
}

type OAuthService interface {
	GetOAuthUser(code string) (*domain.User, error)
	GetOAuthUrl() string
}
