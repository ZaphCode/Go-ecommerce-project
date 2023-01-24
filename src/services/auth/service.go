package auth

import (
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

type Claims struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
}

//* Services

type JWTService interface {
	CreateToken(claims Claims, exp time.Duration, secret string) (string, error)
	DecodeToken(jwtoken string, secret string) (*Claims, error)
}

type OAuthService interface {
	GetOAuthUser(code string) (*domain.User, error)
	GetOAuthUrl() string
}
