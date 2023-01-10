package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	accessJwtSecret  = []byte(cfg.Api.AccessTokenSecret)
	refreshJwtSecret = []byte(cfg.Api.AccessTokenSecret)
)

// Constructor
func NewAuthService() JwtAuthService {
	return &authJwtServiceImpl{}
}

// Implementation
type authJwtServiceImpl struct{}

func (s *authJwtServiceImpl) CreateToken(claims JwtClaims, exp time.Duration, isRefreshType bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":       time.Now().Add(exp).Unix(),
		"user_id":   claims.ID.String(),
		"user_role": claims.Role,
	})

	var tokenString string
	var err error

	if isRefreshType {
		tokenString, err = token.SignedString(refreshJwtSecret)
	} else {
		tokenString, err = token.SignedString(accessJwtSecret)
	}

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authJwtServiceImpl) CreateTokens(claims JwtClaims, access_exp, refresh_exp time.Duration) (string, string, error) {
	accessTokenString, err_1 := s.CreateToken(claims, access_exp, false)

	refreshTokenString, err_2 := s.CreateToken(claims, refresh_exp, true)

	if err_1 != nil || err_2 != nil {
		return "", "", fmt.Errorf("error creating tokens")
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *authJwtServiceImpl) DecodeToken(jwtoken string, refreshType bool) (*JwtClaims, error) {
	token, err := jwt.Parse(jwtoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if refreshType {
			return refreshJwtSecret, nil
		} else {
			return accessJwtSecret, nil
		}
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !(ok && token.Valid) {
		return nil, fmt.Errorf("error getting the claims")
	} else {
		userID := claims["user_id"].(string)
		userRole := claims["user_role"].(string)

		return &JwtClaims{
			ID:   uuid.MustParse(userID),
			Role: userRole,
		}, nil
	}
}
