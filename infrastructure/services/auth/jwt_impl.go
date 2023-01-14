package auth

import (
	"fmt"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/golang-jwt/jwt/v4"
)

// var (
// 	accessJwtSecret  = []byte(cfg.Api.AccessTokenSecret)
// 	refreshJwtSecret = []byte(cfg.Api.AccessTokenSecret)
// )

// Custom types
type CustomJwtClaims struct {
	Claims
	jwt.RegisteredClaims
}

// Constructor
func NewJwtAuthService() JwtAuthService {
	return &jwtAuthServiceImpl{
		accessJwtSecret:  []byte(config.Get().Api.AccessTokenSecret),
		refreshJwtSecret: []byte(config.Get().Api.RefreshTokenSecret),
	}
}

// Implementation
type jwtAuthServiceImpl struct {
	accessJwtSecret  []byte
	refreshJwtSecret []byte
}

func (s *jwtAuthServiceImpl) CreateToken(claims Claims, exp time.Duration, isRefreshType bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomJwtClaims{
		Claims: claims,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    claims.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	})

	var tokenString string
	var err error

	if isRefreshType {
		tokenString, err = token.SignedString(s.refreshJwtSecret)
	} else {
		tokenString, err = token.SignedString(s.refreshJwtSecret)
	}

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *jwtAuthServiceImpl) CreateTokens(claims Claims, access_exp, refresh_exp time.Duration) (string, string, error) {
	accessTokenString, err_1 := s.CreateToken(claims, access_exp, false)

	refreshTokenString, err_2 := s.CreateToken(claims, refresh_exp, true)

	if err_1 != nil || err_2 != nil {
		return "", "", fmt.Errorf("error creating tokens")
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *jwtAuthServiceImpl) DecodeToken(jwtoken string, refreshType bool) (*Claims, error) {
	var customJwtClais CustomJwtClaims

	token, err := jwt.ParseWithClaims(jwtoken, &customJwtClais, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if refreshType {
			return s.refreshJwtSecret, nil
		} else {
			return s.accessJwtSecret, nil
		}
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return &customJwtClais.Claims, nil
}
