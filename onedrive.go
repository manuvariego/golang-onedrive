package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var BaseURL = "https://graph.microsoft.com/v1.0/"
var rootFolder = "/me/drive/root"

// var folderPath = rootFolder

// func ChangeDirectory(folderPath string) {
//
// }

func Menu() {
	var cmd string
	_, err := fmt.Scanln(&cmd)
	if err != nil {
		fmt.Println("Error reading input", err)
		return
	}
	cmd = strings.TrimSpace(cmd)

	if cmd == "ls" {
		fmt.Println("temp")
	}

	if strings.HasPrefix(cmd, "cd") {
		cmd = strings.TrimPrefix(cmd, "cd")
		fmt.Println(cmd)
	}

}

func ListFiles(client *http.Client, folderPath string) ([]DriveItem, error) {

	url := BaseURL + rootFolder + ":/children"
	req, _ := http.NewRequest("GET", url, nil)
	fmt.Println(url)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Error fetching files: %s", body)
	}

	var response struct {
		Value []DriveItem `json:"value"`
	}

	json.NewDecoder(resp.Body).Decode(&response)

	// fmt.Println(response.Value)

	return response.Value, nil
}
