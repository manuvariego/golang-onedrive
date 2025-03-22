package main

import (
	// "context"

	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/manuvariego/golang-onedrive/onedrive"

	"github.com/joho/godotenv"
)

// ONLY TESTING PURPOSES
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	scopes := []string{"User.Read", "offline_access", "Sites.Read.All", "Files.ReadWrite.All"}
	tenantID := os.Getenv("MS_TENANT_ID")
	appID := os.Getenv("MS_OPENGRAPH_APP_ID")
	clientSecret := os.Getenv("MS_OPENGRAPH_CLIENT_SECRET")
	sharePointDrive := os.Getenv("SHAREPOINT_DRIVE")
	sharePointPath := os.Getenv("SHAREPOINT_PATH")
	oauthconf := onedrive.NewOauthConfig(tenantID, appID, clientSecret, scopes)

	//Checks if token.json exists, if it doesn't it is created with a new code from user
	if !onedrive.CheckTokenFile() {
		token := onedrive.GetInitialTokens(oauthconf)
		err = onedrive.SaveToken(token)
		if err != nil {
			log.Fatal("Error saving token", err)
		}
	}

	client, err := onedrive.GetClient(oauthconf)

	sharePoint := sharePointDrive + sharePointPath

	rootUrl := onedrive.GetRootUrl(sharePoint)
	driveUrl := onedrive.GetRootUrl(sharePointDrive)

	//Creates onedriveclient with the data
	od := onedrive.OneDriveClient{Client: client, Path: onedrive.Path{CurrentPath: rootUrl}}

	// var oneDriveRoot = onedrive.Directory{Name: "root"}

	// err = od.LoadOneDrive(&oneDriveRoot, rootUrl)
	// data, err := json.Marshal(oneDriveRoot)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// os.WriteFile("output.json", data, 0666)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	var currentDir onedrive.Directory

	data, err := os.ReadFile("output.json")

	err = json.Unmarshal(data, &currentDir)

	//Testing purposes it is iterated

	for {
		// fmt.Println(&currentDir)
		directories, files, err := onedrive.Ls(&currentDir)
		if err != nil {
			fmt.Println("Error listing directories and files:", err)
		} else {
			fmt.Println("Directories:", directories)
			fmt.Println("Files:", files)
		}

		line := bufio.NewReader(os.Stdin)
		cmd, err := line.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input", err)
			return
		}

		cmd = strings.TrimSpace(cmd)
		_, isFile := onedrive.IsFile(&currentDir, cmd)
		isDirectory := onedrive.IsDirectory(directories, cmd)

		if isFile {
			dwnloadUrl, err := od.GetDownloadUrl(cmd, &currentDir, driveUrl)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(dwnloadUrl)
		}

		if isDirectory {
			fmt.Println("here")
			newDir, err := onedrive.Cd(cmd, &currentDir)
			currentDir = *newDir
			if err != nil {
				fmt.Println("Error changing directory: ", err)
			} else {
				// fmt.Println("New Path: ", newPath.CurrentPath)

			}

		}
	}
}
