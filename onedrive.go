package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func returnStaticPaths() (string, string) {
	baseUrl := fmt.Sprintf("https://graph.microsoft.com/v1.0/me")
	rootUrl := fmt.Sprintf("/drives/%s", os.Getenv("SHAREPOINT_PATH"))
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

func ChangeDirectory(items []Item, directory string, path *Path) bool {
	item, exists := directoryExists(items, directory)
	if !exists {
		fmt.Println("That directory doesn't exist")
		return false
	} else {
		path.CurrentPath += "/" + item.Name
		return true
	}

	// fmt.Println("Current Path has been updated:", *currentPath)
}

func ListFiles(items []Item, path *Path) {

	for _, item := range items {
		if item.IsFolder != nil {
			fmt.Println("DIR: ", item.Name)
		} else {
			fmt.Println("FILE: ", item.Name)
		}

	}
}

// Implement a way for the api only to be called once when getting parent information. (remove getParentPath)
func GetFiles(client *http.Client, path *Path) ([]Item, error) {
	baseUrl, _ := returnStaticPaths()
	var url string
	// path.CurrentPath :
	url = baseUrl + path.CurrentPath + ":/children"

	fmt.Println(url)

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		return nil, err
	}

	var response struct {
		Value []Item `json:"value"`
	}

	json.NewDecoder(resp.Body).Decode(&response)

	return response.Value, nil
}
