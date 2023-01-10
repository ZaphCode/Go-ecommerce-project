package config

import (
	"log"
	"os"

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

type firestore struct {
	Bucket       string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type Config struct {
	Api   api
	OAuth oauth
}

var config Config

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

func GetConfig() *Config {
	return &config
}

func IsProduction() bool {
	return os.Getenv("NODE_ENV") == "prod"
}
