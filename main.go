package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var gmailService *gmail.Service

func initGmailService() {
	config := oauth2.Config{
		ClientID: os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint: google.Endpoint,
		RedirectURL: "http://localhost",
	}

	token := oauth2.Token{
		AccessToken: os.Getenv("ACCESS_TOKEN"),
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		Expiry: time.Now(),
		TokenType: "Bearer",
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	service, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))

	if err != nil {
		log.Fatal("error creating gmail service", err)
	}

	gmailService = service

	if gmailService != nil {
		fmt.Println("gmail service initialized")
	}
}

func main() {
	initGmailService()
}