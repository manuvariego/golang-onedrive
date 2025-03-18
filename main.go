package main

import (
	// "context"

	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	conf := NewOauthConfig()
	//TEMP
	fmt.Println(conf)

	//Checks if token.json exists, if it doesn't it is created with new tokens from user
	if !CheckTokenFile() {
		token := GetInitialTokens(conf)
		err = SaveToken(token)
		if err != nil {
			log.Fatal("Error saving token", err)
		}
	}

	//Client used to make reqs to Onedrive API (or sharepoint) depends.. 8D (es una carita)
	client, err := GetClient(conf)

	x, y := ListFiles(client, "x")
	// fmt.Println("Back in main")
	//TEMP
	fmt.Println(x)
	fmt.Println(y)

}
