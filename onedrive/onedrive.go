package onedrive

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	// "strings"
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

func (d *Directory) IsFile(fileName string) (File, bool) {
	for _, file := range d.Files {
		if file.Name == fileName {
			return file, true
		}
	}
	return File{}, false
}

func SetParents(d *Directory, parent *Directory) {
	d.Parent = parent
	for _, directory := range d.Children {
		SetParents(directory, d)
	}
}

func NewRootDir(sharePoint string) *Directory {
	return &Directory{
		Name: "root",
		Path: fmt.Sprintf("%s/drives/%s", GetBaseUrl(), sharePoint),
	}
}

func (d *Directory) Ls() ([]string, []string, error) {
	// fmt.Println(items)
	var directories []string
	var files []string

	for _, directory := range d.Children {
		directories = append(directories, directory.Name)
	}

	for _, file := range d.Files {
		files = append(files, file.Name)
	}

	return directories, files, nil
}

func (d *Directory) Cd(cmd string) (*Directory, error) {
	if cmd == ".." && d.Parent != nil {
		return d.Parent, nil
	}

	for _, directory := range d.Children {
		if directory.Name == cmd {
			return directory, nil
		}
	}
	return nil, fmt.Errorf("Error not a valid directory")
}

func FetchFileTree(client *http.Client, root *Directory) error {
	url := root.Path + ":/children"
	fmt.Printf("Current Path: %s", root.Path)

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		return err
	}

	var response struct {
		Value []Item
	}

	json.NewDecoder(resp.Body).Decode(&response)
	fmt.Println(response)

	for _, item := range response.Value {
		if item.IsFolder != nil {
			newDirectory := &Directory{
				Name:   item.Name,
				Path:   root.Path + "/" + item.Name,
				Parent: root,
			}
			root.Children = append(root.Children, newDirectory)
			FetchFileTree(client, newDirectory)
		} else {
			root.Files = append(root.Files, File{Name: item.Name, Id: item.ID})
		}
	}

	return nil
}
