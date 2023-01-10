package auth

import (
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/domain"
	"github.com/google/uuid"
)

var (
	cfg = config.GetConfig()
)

type JwtClaims struct {
	ID   uuid.UUID
	Role string
}

//* Services

type JwtAuthService interface {
	CreateToken(JwtClaims, time.Duration, bool) (string, error)
	DecodeToken(string, bool) (*JwtClaims, error)
	CreateTokens(JwtClaims, time.Duration, time.Duration) (string, string, error)
}

type OAuthService interface {
	GetOAuthUser(code string) (*domain.User, error)
	GetOAuthUrl() string
}
