package keygen

import (
	"path/filepath"
	"strings"
	"runtime"
	"os"
	"fmt"
	"io"
	"regexp"
	"errors"
	"net/http"
	"encoding/json"
	"archive/zip"
)

type File struct {
	Filename string `json:"filename"`
	Types []string `json:"types"`
	Version string `json:"version"`
	LastModified int64 `json:"last-modified"`
	URL string `json:"url"`
}

type FileResponse struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Files []File `json:"data"`
}

type GenericResponse struct {
	Status int `json:"status"`
	Message string `json:"message"`
}

type Key struct {
	Domain string
	Code string
	Authcode string
}

func ExtractPackage(src, dest string) error {
    dest = filepath.Clean(dest) + string(os.PathSeparator)

    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer func() {
        if err := r.Close(); err != nil {
            panic(err)
        }
    }()

    os.MkdirAll(dest, 0755)

    // Closure to address file descriptors issue with all the deferred .Close() methods
    extractAndWriteFile := func(f *zip.File) error {
        path := filepath.Join(dest, f.Name)
        // Check for ZipSlip: https://snyk.io/research/zip-slip-vulnerability
        if !strings.HasPrefix(path, dest) {
            return fmt.Errorf("%s: illegal file path", path)
        }

        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer func() {
            if err := rc.Close(); err != nil {
                panic(err)
            }
        }()

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
        } else {
            os.MkdirAll(filepath.Dir(path), f.Mode())
            f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return err
            }
            defer func() {
                if err := f.Close(); err != nil {
                    panic(err)
                }
            }()

            _, err = io.Copy(f, rc)
            if err != nil {
                return err
            }
        }
        return nil
    }

    for _, f := range r.File {
        err := extractAndWriteFile(f)
        if err != nil {
            return err
        }
    }

    return nil
}

func DownloadPackage(file *File, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(file.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func FetchKeyFiles(key *Key) ([]File, error) {
	res, err := http.Get("https://"+key.Domain+"/api/keygen/key/files?code="+key.Code+"&authcode="+key.Authcode)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 299 {
		return nil, errors.New("Internal server failure")
	}
	
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	var response FileResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		var generic GenericResponse
		err = json.Unmarshal(body, &generic)
		if err != nil {
			return nil, err
		} else {
			return nil, errors.New(generic.Message)
		}
	}
	return response.Files, nil
}

func ParseKeyURL(url string) (*Key, error){
	r, err := regexp.Compile(`^https?:\/\/([^/]+)\/access\/([^?]+?)(?:\?authcode(?:=([^&]*))?)?$`)
	if err != nil {
		panic("Failed to compile regex")
	}
	m := r.FindStringSubmatch(url)
	if m == nil {
		return nil, errors.New("The Key URL is malformed")
	}
	key := &Key{}
	key.Domain = m[1]
	key.Code = m[2]
	key.Authcode = m[3]
	return key, nil
}

func FindUpdatedFile(files []File, match *File) *File {
	for i := 0; i < len(files); i++ {
		if ( match.LastModified < files[i].LastModified && match.Filename == files[i].Filename ) {
			return &files[i]
		}
	}
	return nil
}

func osString() string {
	if runtime.GOOS == "darwin" {
		return "apple"
	} else {
		return runtime.GOOS
	}
}

func FindMatchingOSFile(files []File) *File {
	var local = osString()
	for i := 0; i < len(files); i++ {
		for _, candidate := range files[i].Types {
			if candidate == local {
				return &files[i]
			}
		}
	}
	return nil
}
