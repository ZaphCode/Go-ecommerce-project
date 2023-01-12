package config

import (
	"context"
	"encoding/json"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type api struct {
	Port               string `json:"port"`
	ClientOrigin       string `json:"client_origin"`
	ServerHost         string `json:"server_host"`
	AccessTokenSecret  string `json:"access_token_secret"`
	RefreshTokenSecret string `json:"refresh_token_secret"`
	AccessTokenHeader  string `json:"access_token_header"`
	RefreshTokenHeader string `json:"refresh_token_header"`
	RefreshTokenCookie string `json:"refresh_token_cookie"`
	AccessTokenExp     int    `json:"access_token_exp"`
	RefreshTokenExp    int    `json:"refresh_token_exp"`
}

type oauthServices struct {
	Google  oauthCredentials `json:"google"`
	Discord oauthCredentials `json:"discord"`
	Github  oauthCredentials `json:"github"`
}

type oauthCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type storage struct {
	Bucket       string `json:"bucket"`
	UploadFolder string `json:"upload_folder"`
}

type Config struct {
	Api     api           `json:"api"`
	OAuth   oauthServices `json:"oauth"`
	Storage storage       `json:"storage"`
}

var (
	config      Config
	firebaseApp *firebase.App
)

func LoadConfig() {
	var file []byte
	var err error

	if IsProduction() {
		file, err = os.ReadFile("./config/prod_config.json")
	} else {
		file, err = os.ReadFile("./config/dev_config.json")
	}

	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(file, &config); err != nil {
		log.Fatal(err)
	}

	if config.Api.Port == "" || config.Api.ServerHost == "" {
		log.Fatal("json config not readed")
	}
}

func LoadFirebaseConfig() {
	path := "./firebase_config.json"

	ctx := context.Background()

	opt := option.WithCredentialsFile(path)

	app, err := firebase.NewApp(ctx, nil, opt)

	if err != nil {
		log.Fatal("firebase connection fail:", err)
	}

	firebaseApp = app
}

func GetConfig() *Config {
	return &config
}

func GetFirebaseApp() *firebase.App {
	return firebaseApp
}

func IsProduction() bool {
	return os.Getenv("NODE_ENV") == "production"
}
