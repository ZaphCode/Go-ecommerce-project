package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Custom types
type CustomJwtClaims struct {
	Claims
	jwt.RegisteredClaims
}

// Constructor
func NewJWTService() JWTService {
	return new(jwtServiceImpl)
}

// Implementation
type jwtServiceImpl struct{}

func (s *jwtServiceImpl) CreateToken(claims Claims, exp time.Duration, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomJwtClaims{
		Claims: claims,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    claims.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	})

	var tokenString string
	var err error

	tokenString, err = token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *jwtServiceImpl) DecodeToken(jwtoken string, secret string) (*Claims, error) {
	var customJwtClais CustomJwtClaims

	token, err := jwt.ParseWithClaims(jwtoken, &customJwtClais, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return &customJwtClais.Claims, nil
}
