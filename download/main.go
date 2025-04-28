package main

import (
	"log"
	"os"

	"github.com/manuvariego/golang-onedrive/onedrive"

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
    driveID := os.Getenv("DRIVE_ID")

	oauthconf := onedrive.NewOauthConfig(tenantID, appID, clientSecret, scopes)

	//Checks if token.json exists, if it doesn't it is created with a new code from user
	if !onedrive.CheckTokenFile() {
		token := onedrive.GetInitialTokens(oauthconf)
		err = onedrive.SaveToken(token)
		if err != nil {
			log.Fatal("Error saving token", err)
		}
	}

	client, err := onedrive.GetClient(oauthconf, driveID, sharePoint)

    root := onedrive.NewRootDir()
    onedrive.DownloadFileTree(client, root, "output")
}
