package onedrive

import (
	"net/http"
)

type OneDriver interface {
	Pwd() Path
	Ls() ([]string, error)
	Cd(folder string) (Path, error)
}

type OneDriveClient struct {
	Client     *http.Client
	CurrentDir *Directory
	RootDir    Directory
}

// A single DriveItem from the OneDrive API (modifiable)
type Item struct {
	Name        string `json:"name"`
	IsFolder    *bool  `json:"folder,omitempty"`
	ID          string `json:"id"`
	DownloadUrl string `json:"@microsoft.graph.downloadUrl"`
}

type Directory struct {
	Path     string       `json:"path"`
	Name     string       `json:"name"`
	Files    []File       `json:"files,omitempty"`
	Children []*Directory `json:"folders,omitempty"`
	Parent   *Directory   `json:"-"`
}

type File struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DownloadUrl string `json:"download_url"`
}

type Path struct {
	CurrentPath string
}
