package main

import (
	// "context"

	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	//TODO: Return http.Client to use to access ONEDRIVE API
	// ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conf := NewOauthConfig()
	fmt.Println(conf)
	//Checks if token.json exists, if it doesn't it is created with new tokens from user
	if !CheckTokenFile() {
		fmt.Println("here")
		token := GetInitialTokens(conf)
		err = SaveToken(token)
		if err != nil {
			log.Fatal("Error saving token")
		}
	}

	client, err := GetClient(conf)

	// endpoint

	x, y := ListFiles(client, "x")
	fmt.Println("Back in main")
	fmt.Println(x)
	fmt.Println(y)

}
