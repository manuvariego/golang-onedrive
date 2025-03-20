package main

import (
	// "context"

	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	scopes := []string{"User.Read", "offline_access", "Sites.Read.All", "Files.ReadWrite.All"}
	tenantID := os.Getenv("MS_TENANT_ID")
	appID := os.Getenv("MS_OPENGRAPH_APP_ID")
	clientSecret := os.Getenv("MS_OPENGRAPH_CLIENT_SECRET")
	sharePoint := os.Getenv("SHAREPOINT_PATH")
	oauthconf := NewOauthConfig(tenantID, appID, clientSecret, scopes)

	//Checks if token.json exists, if it doesn't it is created with a new code from user
	if !CheckTokenFile() {
		token := GetInitialTokens(oauthconf)
		err = SaveToken(token)
		if err != nil {
			log.Fatal("Error saving token", err)
		}
	}

	//Client used to make reqs to Onedrive API (or sharepoint) depends.. 8D (es una carita)
	client, err := GetClient(oauthconf)
	Menu2(client, sharePoint)
}
