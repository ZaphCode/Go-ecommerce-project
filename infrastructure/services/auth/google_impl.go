package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/domain"
	"github.com/ZaphCode/clean-arch/infrastructure/utils"
)

var (
	googleOAuthUrl = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenUrl = "https://oauth2.googleapis.com/token"
	googleUserUrl  = "https://www.googleapis.com/oauth2/v1/userinfo"
)

// Custom types
type GoogleTokens struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (u GoogleUser) AdaptToUser() (user domain.User) {
	user.Username = u.Name
	user.Email = u.Email
	user.Password = utils.RandomString(20)
	user.ImageUrl = u.Picture
	user.VerifiedEmail = u.VerifiedEmail
	return
}

// Construtor
func NewGoogleOAuthService() OAuthService {
	return &googleOAuthServiceImpl{}
}

// Implementation
type googleOAuthServiceImpl struct{}

func (s *googleOAuthServiceImpl) GetOAuthUrl() string {
	cfg := config.Get()
	params := url.Values{}

	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	}

	params.Add("redirect_uri", cfg.Api.ServerHost+"/api/auth/google/callback")
	params.Add("client_id", cfg.OAuth.Google.ClientID)
	params.Add("access_type", "offline")
	params.Add("response_type", "code")
	params.Add("prompt", "consent")
	params.Add("scope", strings.Join(scopes, " "))

	url := fmt.Sprintf("%s?%s", googleOAuthUrl, params.Encode())

	return url
}

func (s *googleOAuthServiceImpl) GetOAuthUser(code string) (*domain.User, error) {
	tokens, err := s.getGoogleTokens(code)

	if err != nil {
		return nil, fmt.Errorf("getGoogleTokens() error: %v", err)
	}

	googleUser, err := s.getGoogleUser(tokens)

	if err != nil {
		return nil, fmt.Errorf("getGoogleUser() error: %v", err)
	}

	user := googleUser.AdaptToUser()

	return &user, nil
}

func (s *googleOAuthServiceImpl) getGoogleTokens(code string) (*GoogleTokens, error) {
	cfg := config.Get()
	form := url.Values{}

	form.Add("code", code)
	form.Add("client_id", cfg.OAuth.Google.ClientID)
	form.Add("client_secret", cfg.OAuth.Google.ClientSecret)
	form.Add("redirect_uri", cfg.Api.ServerHost+"/api/auth/google/callback")
	form.Add("grant_type", "authorization_code")

	req, err := http.NewRequest(http.MethodPost, googleTokenUrl, strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error fetching to google api | %w", err)
	}

	resBody, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		fmt.Println(string(resBody))
		return nil, fmt.Errorf("error fetching to google api")
	}

	tokens := GoogleTokens{}

	if err := json.Unmarshal(resBody, &tokens); err != nil || tokens.AccessToken == "" {
		return nil, fmt.Errorf("error getting tokens: %w", err)
	}

	return &tokens, nil
}

func (s *googleOAuthServiceImpl) getGoogleUser(tokens *GoogleTokens) (*GoogleUser, error) {
	url := googleUserUrl + "?alt=json&access_token=" + tokens.AccessToken

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+tokens.IDToken)

	res, err := http.DefaultClient.Do(req)

	if err != nil || res.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching to google api | %w", err)
	}

	resBody, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	user := GoogleUser{}

	if err := json.Unmarshal(resBody, &user); err != nil || user.Email == "" {
		return nil, fmt.Errorf("error parsing user. %w", err)
	}

	return &user, nil
}
