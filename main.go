package main

import (
	// "context"

	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	//currentPath could be global across packages?

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

	var path Path
	client, err := GetClient(oauthconf)
	path.CurrentPath = ""

	//Client used to make reqs to Onedrive API (or sharepoint) depends.. 8D (es una carita)
	for {
		items, err := GetFiles(client, &path)
		if err != nil {
			fmt.Println("t")
		}
		ListFiles(items, &path)
		bool := ChangeDirectory(items, "cd", &path)
		fmt.Println(path.CurrentPath)
		if !bool {
			fmt.Println("cd didn't work")
		}

	}

}
