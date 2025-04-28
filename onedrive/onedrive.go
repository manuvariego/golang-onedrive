package onedrive

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
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

func NewRootDir() *Directory {
	return &Directory{
		Name: "root",
		Path: "",
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

func FetchFileTree(client *Client, root *Directory) error {
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

/* Download all the structure and files of a root path */
func DownloadFileTree(client *Client, root *Directory, output string) {
    children, err := root.fetchChildren(client)
    if err != nil {
        slog.Error("error fetching childrens", "err", err.Error())
    }
    if err := createDir(output, ""); err != nil {
        slog.Error("error creating output folder", "err", err.Error())
        return
    }

	for _, item := range children {
		if item.IsFolder != nil {
			newDirectory := &Directory{
				Name:   item.Name,
				Path:   path.Join(root.Path, item.Name),
				Parent: root,
			}
			root.Children = append(root.Children, newDirectory)

            if err := createDir(output, newDirectory.Path); err != nil {
                slog.Error("error creating folder", "path", newDirectory.Path, "err", err.Error())
                return
            }

			DownloadFileTree(client, newDirectory, output)
		} else {
			root.Files = append(root.Files, &File{Name: item.Name, Id: item.ID, DownloadUrl: item.DownloadUrl})
		}
	}

    // Download all the files
    var wg sync.WaitGroup
    for _, file := range root.Files {
        wg.Add(1)
        go func() {
            defer wg.Done()
            if err := downloadFile(client, file, output, root.Path); err != nil {
                slog.Error(err.Error())
            }
        }()
    }
    wg.Wait()
}

func downloadFile(client *Client, f *File, output string, parentPath string) error {
    res, err := client.Get(f.DownloadUrl)
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        return fmt.Errorf("invalid server response: %s", res.Status)
    }

    file, err := os.Create(path.Join(output, parentPath, f.Name))
    if err != nil {
        return fmt.Errorf("error creating the file: %v", err)
    }
    defer file.Close()

    if _, err := io.Copy(file, res.Body); err != nil {
        return fmt.Errorf("error saving the file: %v", err)
    }
    return nil
}

func createDir(out string, p string) error {
    outPath := path.Join(out, p)
    if _, err := os.Stat(outPath); os.IsNotExist(err) {
        err := os.Mkdir(outPath, os.ModeDir)
        if err != nil {
            return err
        }
    }
    return nil
}

func (d *Directory) fetchChildren(client *Client) ([]Item, error) {
    u, err := url.JoinPath(client.baseURL, d.Path)
    url := fmt.Sprintf("%s:/children/", u)
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

	// fmt.Println(response.Items)
    return response.Items, nil
}
