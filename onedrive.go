package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
)

func ReturnStaticPaths() (string, string) {
	baseUrl := fmt.Sprintf("https://graph.microsoft.com/v1.0/me")
	rootUrl := fmt.Sprintf("/drives/%s", os.Getenv("SHAREPOINT_PATH"))
	return baseUrl, rootUrl
}

func (od *OneDriveClient) Pwd() Path {
	return od.Path
}

func (od *OneDriveClient) GetDownloadUrl(itemName string) (string, error) {
	items, err := od.GetFiles()
	if err != nil {
		return "", err
	}

	for _, item := range items {
		if item.Name == itemName {
			return item.DownloadUrl, nil
		}

	}
	return "", nil
}

func (od *OneDriveClient) IsDirectory(directories []string, directory string) bool {
	exists := slices.Contains(directories, directory)
	if !exists {
		return false
	}
	return true
}

func (od *OneDriveClient) IsFile(files []string, file string) bool {
	exists := slices.Contains(files, file)
	if !exists {
		return false
	}
	return true
}

func (od *OneDriveClient) Ls() ([]string, []string, error) {
	items, err := od.GetFiles()
	if err != nil {
		return nil, nil, err
	}
	// fmt.Println(items)

	var folders []string
	var files []string
	for _, item := range items {
		if item.IsFolder != nil {
			folders = append(folders, item.Name)
		} else if item.DownloadUrl != "" {
			files = append(files, item.Name)
		}

	}

	return folders, files, nil
}

func (od *OneDriveClient) Cd(directory string) (Path, error) {

	od.Path.CurrentPath += "/" + directory
	return od.Path, nil
}

func (od *OneDriveClient) GetFiles() ([]Item, error) {
	baseUrl, _ := ReturnStaticPaths()
	// path.CurrentPath :
	url := baseUrl + od.Path.CurrentPath + ":/children"

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := od.Client.Do(req)
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
