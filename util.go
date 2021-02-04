package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/imroc/req"
)

type Releases struct {
	Assets []struct {
		Name               string `json:"name"`
		ContentType        string `json:"content_type"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getLatestRelease(repo string, name string, mime string) {
	url := "https://api.github.com/repos/" + repo + "/releases/latest"
	r, err := req.Get(url)
	if err != nil {
		log.Fatal("Failed to get releases\n")
	} else {
		resp := r.Response()
		b, _ := ioutil.ReadAll(resp.Body)
		var latest Releases
		json.Unmarshal(b, &latest)
		for _, asset := range latest.Assets {
			if asset.ContentType == strings.ToLower(mime) {
				if strings.Contains(strings.ToLower(asset.BrowserDownloadURL), strings.ToLower(name)) {
					// Found asset
					filename := asset.Name
					Download(asset.BrowserDownloadURL, filename)
					break
				} else {
					log.Fatalf("Failed to find assets contains %s", name)
				}
			}
		}
	}
}

func Download(url string, filename string) {
	progress := func(current, total int64) {
		fmt.Printf("Download %f %% ...", float32(current)/float32(total)*100)
	}
	r, err := req.Get(url, req.DownloadProgress(progress))
	r.ToFile(filename)

	if err != nil {
		log.Fatal("Failed to download %s\n", filename)
	}
}

func LogToFile(path string, prefix string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	log.SetOutput(f)
	log.SetPrefix(prefix)
	return f, nil
}