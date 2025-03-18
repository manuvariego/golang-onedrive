package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

func NewOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("MS_APP_ID"),
		ClientSecret: os.Getenv("MS_APP_SECRET"),
		Scopes:       []string{"User.Read", "Files.Read.All", "offline_access"},
		RedirectURL:  "http://localhost/auth",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		},
	}
}

func CheckTokenFile() bool {
	_, err := os.Open("token.json")
	if err != nil {
		fmt.Println(err)
		return false
	}

	tok, err := LoadToken()
	if tok.RefreshToken == "" {
		return false

	}

	return true
}

func SaveToken(token *oauth2.Token) error {
	data, err := json.Marshal(token)
	os.WriteFile("token.json", data, 0666)

	return err
}

func LoadToken() (*oauth2.Token, error) {
	var token oauth2.Token

	data, err := os.ReadFile("token.json")

	err = json.Unmarshal(data, &token)

	return &token, err

}

func GetInitialTokens(conf *oauth2.Config) *oauth2.Token {

	verifier := oauth2.GenerateVerifier()

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	fmt.Printf("Visitar el URL y copiar el codigo para autorizar: %v\n\n", url)

	var code string

	fmt.Println("Ingrese el codigo: ")
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	tok, err := conf.Exchange(context.Background(), code, oauth2.VerifierOption(verifier))

	if err != nil {
		log.Fatalf("Error al generar el token: %v", err)
	}

	return tok

}

func GetValidToken(conf *oauth2.Config) (*oauth2.Token, error) {
	currentToken, err := LoadToken()
	if err != nil {
		log.Fatal("Error al cargar el token del archivo")
		return nil, err
	}

	//Checks if the currentToken is valid
	if !currentToken.Valid() {
		tokenSource := conf.TokenSource(context.Background(), currentToken)
		fmt.Println("Successfully refreshed the token")
		refreshedToken, err := tokenSource.Token()
		err = SaveToken(refreshedToken)
		if err != nil {
			log.Fatal("Failed to refresh token", err)
			return nil, err
		}

		return refreshedToken, nil
	}

	return currentToken, nil
}

// Add some sort of context throughout everything (no idea)
func GetClient(conf *oauth2.Config) (*http.Client, error) {
	token, err := GetValidToken(conf)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(context.Background(), token)
	return client, err
}
