package main

import (
	// "context"

	"log"

	"github.com/joho/godotenv"
)

func main() {

	//currentPath could be global across packages?
	var currentPath string

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	oauthconf := NewOauthConfig()

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
	Menu2(client, &currentPath)

}
