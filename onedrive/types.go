package onedrive

import (
	"net/http"
)

type OneDriver interface {
	Ls() ([]string, []string, error)
	Cd(directory string) (*Directory, error)
}

type OneDriveClient struct {
	Client  *http.Client
	RootDir Directory
}

// A single DriveItem from the OneDrive API (modifiable)
type Item struct {
	Name        string `json:"name"`
    // If the field 'folder' its present the item its a folder, otherwise its a file. 
    // does not have a boolean isFolder field
	IsFolder    *interface{} `json:"folder,omitempty"`
	ID          string `json:"id"`
}

type Directory struct {
	Path     string       `json:"path"`
	Name     string       `json:"name"`
	Files    []*File      `json:"files,omitempty"`
	Children []*Directory `json:"folders,omitempty"`
	Parent   *Directory   `json:"-"`
}

type File struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
}
