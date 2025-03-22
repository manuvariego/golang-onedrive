package onedrive

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	// "strings"
	"time"
)

func GetBaseUrl() string {
	return fmt.Sprintf("https://graph.microsoft.com/v1.0/me")
}

func GetRootUrl(sharePoint string) string {
	return fmt.Sprintf("/drives/%s", sharePoint)

}

func IsDirectory(directories []string, directory string) bool {
	return slices.Contains(directories, directory)
}

func IsFile(currentDir *Directory, fileName string) (File, bool) {
	for _, file := range currentDir.Files {
		if file.Name == fileName {
			return file, true
		}
	}
	return File{}, false
}

func (od *OneDriveClient) LoadOneDrive(rootDir *Directory, rootUrl string) error {
	return od.FetchFileTree(rootDir, rootUrl)
}

func Ls(currentDir *Directory) ([]string, []string, error) {
	// fmt.Println(items)
	var directories []string
	var files []string

	for _, directory := range currentDir.Children {
		directories = append(directories, directory.Name)
	}

	for _, file := range currentDir.Files {
		files = append(files, file.Name)
	}

	return directories, files, nil
}

func Cd(cmd string, currentDir *Directory) (*Directory, error) {
	for _, directory := range currentDir.Children {
		if directory.Name == cmd {
			return &directory, nil
		}
	}
	return nil, fmt.Errorf("Error not a valid directory")

}

func (od *OneDriveClient) GetDownloadUrl(fileName string, currentDir *Directory, drivePath string) (string, error) {
	file, exists := IsFile(currentDir, fileName)
	if !exists {
		return "", fmt.Errorf("Not a valid file")

	} else {
		fmt.Println(file.Id)
		return od.FetchDownloadUrl(file.Id, drivePath)
	}
}

func (od *OneDriveClient) FetchFileTree(parentDirectory *Directory, path string) error {
	duration := time.Duration(100) * time.Millisecond
	time.Sleep(duration)
	baseUrl := GetBaseUrl()
	// path.CurrentPath :
	url := baseUrl + path + ":/children"
	fmt.Printf("Current Path: %s", path)

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := od.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		return err
	}

	var response struct {
		Value []struct {
			Name   string    `json:"name"`
			Id     string    `json:"id"`
			Folder *struct{} `json:"folder,omitempty"`
			File   *struct{} `json:"file,omitempty"`
			// DownloadUrl string    `json:"@microsoft.graph.downloadUrl,omitempty"`
		}
	}

	json.NewDecoder(resp.Body).Decode(&response)
	fmt.Println(response)

	for _, item := range response.Value {
		if item.Folder != nil {
			newDirectory := Directory{Name: item.Name}
			parentDirectory.Children = append(parentDirectory.Children, newDirectory)
			//Recursivamente recorre el arbol, [len(parentDirectory.Children)-1] esto hace referencia al directorio recien agregado
			od.FetchFileTree(&parentDirectory.Children[len(parentDirectory.Children)-1], path+"/"+item.Name)
		} else if item.File != nil {
			parentDirectory.Files = append(parentDirectory.Files, File{Name: item.Name, Id: item.Id})
		}
	}

	return nil
}

func (od *OneDriveClient) FetchDownloadUrl(itemId string, path string) (string, error) {
	baseUrl := GetBaseUrl()

	url := baseUrl + path + "/items/" + itemId

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := od.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		return "", err
	}

	var response struct {
		DownloadUrl string `json:"@microsoft.graph.downloadUrl"`
	}

	json.NewDecoder(resp.Body).Decode(&response)

	return response.DownloadUrl, nil
}
