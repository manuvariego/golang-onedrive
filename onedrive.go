package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func returnStaticPaths() (string, string) {
	baseUrl := fmt.Sprintf("https://graph.microsoft.com/v1.0/")
	rootUrl := fmt.Sprintf("me/drives/%s", os.Getenv("SHAREPOINT_PATH"))
	return baseUrl, rootUrl
}

func directoryExists(items []Item, cmd string) (Item, bool) {
	for _, item := range items {
		if item.Name == cmd && item.IsFolder != nil {
			return item, true
		}
	}
	return Item{}, false
}

func ChangeDirectory(item Item, currentPath *string) {

	newPath := *currentPath + "/" + item.Name

	*currentPath = newPath

	fmt.Println("Current Path has been updated:", *currentPath)

}

func ListFiles(client *http.Client, currentPath string) ([]Item, error) {
	baseUrl, rootUrl := returnStaticPaths()

	url := baseUrl + rootUrl + currentPath + ":/children"
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		os.Exit(1)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Error fetching files: %s", body)
	}

	var response struct {
		Value []Item `json:"value"`
	}

	json.NewDecoder(resp.Body).Decode(&response)

	return response.Value, nil
}

func Menu2(client *http.Client, currentPath *string) {
	// var cmd string
	// var files []Item

	for {
		items, err := ListFiles(client, *currentPath)

		if err != nil {
			fmt.Println("Error listing files:", err)
			continue
		}
		line := bufio.NewReader(os.Stdin)
		cmd, err := line.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input", err)
			return
		}
		cmd = strings.TrimSpace(cmd)
		if cmd == "ls" {
			for _, item := range items {
				if item.IsFolder != nil {
					fmt.Println("DIR: ", item.Name)
				} else {
					fmt.Println("FILE: ", item.Name)
				}

			}

		} else if strings.HasPrefix(cmd, "cd") {
			fmt.Println("inside change directory statement")
			cmd = strings.TrimPrefix(cmd, "cd")
			cmd = strings.TrimSpace(cmd)
			//TEMP
			// fmt.Println(cmd)
			item, exists := directoryExists(items, cmd)
			if !exists {
				fmt.Println("That directory doesn't exist")
				continue

			}
			ChangeDirectory(item, currentPath)
		}

	}

}
