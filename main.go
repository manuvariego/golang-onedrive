package main

import (
	// "context"

	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/joho/godotenv"
)

func main() {

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

	client, err := GetClient(oauthconf)

	_, rootUrl := ReturnStaticPaths()

	//Creates onedriveclient with the data
	od := OneDriveClient{Client: client, Path: Path{CurrentPath: rootUrl}}

	//Testing purposes it is iterated
	// fmt.Println("Current Path:", od.Pwd().CurrentPath)

	for {
		folders, files, err := od.Ls()
		if err != nil {
			fmt.Println("Error listing folders and files:", err)
		} else {
			fmt.Println("Folders:", folders)
			fmt.Println("Files:", files)
		}

		line := bufio.NewReader(os.Stdin)
		cmd, err := line.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input", err)
			return
		}

		cmd = strings.TrimSpace(cmd)
		isFile := slices.Contains(files, cmd)
		isFolder := slices.Contains(folders, cmd)

		if isFile {
			dwnloadUrl, err := od.GetDownloadUrl(cmd)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(dwnloadUrl)
		}

		if isFolder {
			newPath, err := od.Cd(cmd)
			if err != nil {
				fmt.Println("Error changing directory: ", err)
			} else {
				fmt.Println("New Path: ", newPath.CurrentPath)

			}

		}
	}
}
