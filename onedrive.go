package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var BaseURL = "https://graph.microsoft.com/v1.0/"
var rootFolder = "/me/drive/root"

// var folderPath = rootFolder

type DriveItem struct {
	Name     string `json:"name"`
	IsFolder bool   `json:"folder,omitempty"`
	ID       string `json:"id"`
}

// func ChangeDirectory(folderPath string) {
//
// }

func menu() {
	var op string
	_, err := fmt.Scanln(&op)
	if err != nil {
		fmt.Println("Error reading input", err)
		return
	}
	switch op {

	case "ls":

	case "cd":
	}

}

func ListFiles(client *http.Client, folderPath string) ([]DriveItem, error) {

	url := BaseURL + rootFolder + "/children"
	req, _ := http.NewRequest("GET", url, nil)
	fmt.Println(url)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error fetching files: %s", body)
	}

	var response struct {
		Value []DriveItem `json:"value"`
	}

	json.NewDecoder(resp.Body).Decode(&response)

	// fmt.Println(response.Value)

	return response.Value, nil
}
