package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	"gopkg.in/yaml.v3"
)

type api struct {
	Port               string `yaml:"port"`
	ServerHost         string `yaml:"server_host"`
	AccessTokenSecret  string `yaml:"access_token_secret"`
	RefreshTokenSecret string `yaml:"refresh_token_secret"`
	AccessTokenHeader  string `yaml:"access_token_header"`
	RefreshTokenHeader string `yaml:"refresh_token_header"`
	RefreshTokenCookie string `yaml:"refresh_token_cookie"`
}

type oauth struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type storage struct {
	Bucket       string `yaml:"bucket"`
	UploadFolder string `yaml:"upload_folder"`
}

type Config struct {
	Api   api
	OAuth oauth
}

var config Config

var firebaseApp *firebase.App

func LoadYAMLConfig() {
	file, err := os.ReadFile("./config/dev_config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatal(err)
	}

	if config.Api.Port == "" || config.OAuth.ClientID == "" {
		log.Fatal("yaml config not readed")
	}
}

func LoadFireBaseJsonConfig() {
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
	return os.Getenv("NODE_ENV") == "prod"
}
