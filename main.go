package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
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

func fetchEmail(id string) string {
	// seems like the default format is FULL – https://developers.google.com/gmail/api/reference/rest/v1/Format
	// seems like I have to call it with FULL to get the headers,
	// and call it again with RAW to get the payload 🤷‍♂️
	rawResponse, err := gmailService.Users.Messages.Get("me", id).Format("RAW").Do()
	fullResponse, err := gmailService.Users.Messages.Get("me", id).Format("FULL").Do()

	// modifyRequest := gmail.ModifyMessageRequest{
	// 	RemoveLabelIds: []string{"UNREAD"},
	// }
	// gmailService.Users.Messages.Modify("me", id, &modifyRequest)

	if err != nil {
		log.Fatal("error fetching message", err)
	}

	for _, header := range fullResponse.Payload.Headers {
		if header.Name == "From" {
			fmt.Println(header.Value)
		}
	}
	decoded, err := base64.URLEncoding.DecodeString(rawResponse.Raw)

	stringified := string(decoded[:])
	// fmt.Println(stringified)

	return stringified
}

func fetchEmails() {
	emailsResponse, err := gmailService.Users.Messages.List("me").Q("label:UNREAD").Do()

	if err != nil {
		log.Fatal("error getting messages \n", err)
	}
	emails := emailsResponse.Messages

	// fullMessage := fetchEmail(emails[0].Id)
	// fmt.Println("fullMessage", fullMessage)

	// re := regexp.MustCompile(`Your transaction of \$.+?\.`)
	// match := re.FindString(fullMessage) // why does this line print out the fullMessage?
	// fmt.Println("match", match)

	// for _, msg := range emails {
		// fmt.Println(msg.Id)
		fetchEmail(emails[0].Id)
	// }
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("err getting env", err)
	}
	initGmailService()
	fetchEmails()
}