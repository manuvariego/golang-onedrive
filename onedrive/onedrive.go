package onedrive

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
)

const baseUrl = "https://graph.microsoft.com/v1.0/me"

func (d *Directory) IsFile(fileName string) (*File, bool) {
	for _, file := range d.Files {
		if file.Name == fileName {
			return file, true
		}
	}
	return &File{}, false
}

func SetParents(d *Directory, parent *Directory) {
	d.Parent = parent
	for _, directory := range d.Children {
		SetParents(directory, d)
	}
}

func NewRootDir(driveID, sharePoint string) *Directory {
	return &Directory{
		Name: "root",
		Path: path.Join("drives", driveID, sharePoint),
    }
}

func (d *Directory) Ls() ([]string, []string, error) {
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
    children, err := root.fetchChildren(client)
    if err != nil {
        return err
    }

	for _, item := range children {
		if item.IsFolder != nil {
			newDirectory := &Directory{
				Name:   item.Name,
				Path:   path.Join(root.Path, item.Name),
				Parent: root,
			}
			root.Children = append(root.Children, newDirectory)
			FetchFileTree(client, newDirectory)
		} else {
			root.Files = append(root.Files, &File{Name: item.Name, Id: item.ID})
		}
	}

	return nil
}

func (d *Directory) fetchChildren(client *http.Client) ([]Item, error) {
    url := fmt.Sprintf("%s/%s:/children", baseUrl, d.Path)
	log.Printf("fetching %s\n", d.Path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("fetch children request failed with status %d", res.StatusCode))
	}

    var response struct {
        Items []Item `json:"value"`
	}

    if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
    }

	fmt.Println(response.Items)
    return response.Items, nil
}
