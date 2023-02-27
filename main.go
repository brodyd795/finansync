package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
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

func fetchEmails() {
	response, err := gmailService.Users.Messages.List("me").Q("label:UNREAD").Do()

	if err != nil {
		log.Fatal("error getting messages \n", err)
	}
	mssgs := response.Messages
	fmt.Println(len(mssgs))
	fullMessage := fetchEmail(mssgs[0].Id)

	re := regexp.MustCompile(`Your transaction of \$.+?\.`)
	match := re.FindString(fullMessage) // why does this line print out the fullMessage?
	fmt.Println("match", match)

	// for _, msg := range mssgs {
	// 	fmt.Println(msg.Id)
	// 	fetchEmail(mssgs[0].Id)
	// }
}

func fetchEmail(id string) string {
	response, err := gmailService.Users.Messages.Get("me", id).Format("RAW").Do()

	modifyRequest := gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}
	gmailService.Users.Messages.Modify("me", id, &modifyRequest)

	if err != nil {
		log.Fatal("error fetching message", err)
	}

	decoded, err := base64.URLEncoding.DecodeString(response.Raw)

	stringified := string(decoded[:])
	fmt.Println(stringified)

	return stringified
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("err getting env", err)
	}
	initGmailService()
	fetchEmails()
}