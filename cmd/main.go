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

	// rootUrl := onedrive.GetRootUrl(sharePoint)
	// driveUrl := onedrive.GetRootUrl(sharePointDrive)

	//Creates onedriveclient with the data
	// od := onedrive.OneDriveClient{Client: client, CurrentDir: &onedrive.Directory{Name: "root"}}

	var root *onedrive.Directory

	fetchTree := false
	data, err := os.ReadFile("output.json")
	if err != nil {
		fmt.Printf("error reading tree file")
		fetchTree = true
	}

	err = json.Unmarshal(data, &root)
	if err != nil {

		fmt.Printf("\nerror unmarshalling tree: %v\n", err)
		fetchTree = true
	}

	// err = od.LoadOneDrive(od.CurrentDir, rootUrl)

	if fetchTree {
		root = onedrive.NewRootDir(sharePoint)
		onedrive.FetchFileTree(client, root)
		fmt.Println(root)

		data, err = json.Marshal(&root)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		err = os.WriteFile("output.json", data, 0666)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	} else {
		onedrive.SetParents(root, nil)
	}

	// var rootDir onedrive.Directory
	// var currentDir onedrive.Directory

	//Testing purposes it is iterated

	for {
		directories, files, err := root.Ls()
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
		// _, isFile := onedrive.IsFile(&currentDir, cmd)
		// isDirectory := onedrive.IsDirectory(directories, cmd)
		//
		// if isFile {
		// 	dwnloadUrl, err := od.GetDownloadUrl(cmd, &currentDir, driveUrl)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// 	fmt.Println(dwnloadUrl)
		// }

		newRoot, err := root.Cd(cmd)
		if err != nil {
			fmt.Println("Error changing directory: ", err)
			continue
		}

		root = newRoot
	}
	// for {
	// 	// fmt.Println(&currentDir)
	// 	directories, files, err := onedrive.Ls(&currentDir)
	// 	if err != nil {
	// 		fmt.Println("Error listing directories and files:", err)
	// 	} else {
	// 		fmt.Println("Directories:", directories)
	// 		fmt.Println("Files:", files)
	// 	}
	//
	// 	line := bufio.NewReader(os.Stdin)
	// 	cmd, err := line.ReadString('\n')
	// 	if err != nil {
	// 		fmt.Println("Error reading input", err)
	// 		return
	// 	}
	//
	// 	cmd = strings.TrimSpace(cmd)
	// 	_, isFile := onedrive.IsFile(&currentDir, cmd)
	// 	isDirectory := onedrive.IsDirectory(directories, cmd)
	//
	// 	if isFile {
	// 		dwnloadUrl, err := od.GetDownloadUrl(cmd, &currentDir, driveUrl)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		fmt.Println(dwnloadUrl)
	// 	}
	//
	// 	if isDirectory {
	// 		fmt.Println("here")
	// 		newDir, err := onedrive.Cd(cmd, &currentDir)
	// 		currentDir = *newDir
	// 		fmt.Println(currentDir)
	// 		if err != nil {
	// 			fmt.Println("Error changing directory: ", err)
	// 		} else {
	// 			// fmt.Println("New Path: ", newPath.CurrentPath)
	//
	// 		}
	//
	// 	}
	// }
}
