package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	GoogleLoginConfig oauth2.Config
	GitHubLoginConfig oauth2.Config
}

var AppConfig Config

func GoogleConfig() oauth2.Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	AppConfig.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  "http://localhost:3000/google-callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return AppConfig.GoogleLoginConfig
}
