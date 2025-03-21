package onedrive

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

func NewOauthConfig(tenantID, appID, clientSecret string, scopes []string) *oauth2.Config {
	if tenantID == "" || appID == "" || clientSecret == "" {
		panic("the enviroment variables cant be empty")
	}

	authUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenantID)
	tokenUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	return &oauth2.Config{
		Scopes:       scopes,
		ClientID:     appID,
		ClientSecret: clientSecret,
		//Temp for development
		RedirectURL: "https://localhost:3000",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authUrl,
			TokenURL: tokenUrl,
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

	fmt.Println("Enterr the code: ")
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	tok, err := conf.Exchange(context.Background(), code, oauth2.VerifierOption(verifier))

	if err != nil {
		log.Fatalf("Error generating token: %v", err)
	}

	return tok
}

func GetValidToken(conf *oauth2.Config) (*oauth2.Token, error) {
	currentToken, err := LoadToken()
	if err != nil {
		log.Fatal("Error while loading token")
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
