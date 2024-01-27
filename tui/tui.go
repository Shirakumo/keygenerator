package tui

import (
	"keygenerator/keygen"
	"keygenerator/config"
	"path/filepath"
	"strconv"
	"bufio"
	"time"
	"fmt"
	"os"
)

func prompt(message string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(message)
	fmt.Print("> ")
	scanner.Scan()
	return scanner.Text()
}

func Main(conf *config.Config) {
	var key *keygen.Key = nil
	var err error = nil
	var files []keygen.File = nil

	for true {
		if conf.KeyURL == "" {
			conf.KeyURL = prompt("Enter the Key URL:")
			if conf.KeyURL == "" {
				return
			}
		}

		key, err = keygen.ParseKeyURL(conf.KeyURL)
		if err != nil {
			fmt.Println(err)
			conf.KeyURL = ""
		} else {
			fmt.Println("Fetching entries for "+key.Code)
			files, err = keygen.FetchKeyFiles(key)
			if err != nil {
				fmt.Println(err)
				conf.KeyURL = ""
			} else {
				break
			}
		}
	}

	conf.LastChecked = time.Now().Unix()

	candidate := keygen.FindMatchingOSFile(files)
	if conf.LocalFile != nil {
		candidate = keygen.FindUpdatedFile(files, conf.LocalFile)
	}

	selected := 0
	fmt.Println("")
	for i := 0; i < len(files); i++ {
		f := &files[i]
		t := time.Unix(f.LastModified, 0)
		if f == candidate {
			selected = i
			fmt.Printf("=> ")
		} else {
			fmt.Printf("   ")
		}
		fmt.Printf("%2v: %v %v %10v %v\n", i, f.Version, t.Format("2006-01-02 15:04:05"), f.Types, f.Filename)
	}
	fmt.Println("")
	for true {
		res := prompt(fmt.Sprintf("Select the version to download [%v]", selected))
		if res != "" {
			selected, err := strconv.Atoi(res);
			if err != nil {
				fmt.Println(err)
			} else {
				candidate = &files[selected]
				break
			}
		} else {
			break
		}
	}
	
	if candidate == nil {
		return
	} else {
		path := filepath.Join(conf.PackagePath, candidate.Filename)
		res := prompt(fmt.Sprintf("Where should the package be saved? [%v]", path))
		if res != "" {
			path = res
		}
		
		fmt.Println("Downloading to "+path)
		err = keygen.DownloadPackage(candidate, path)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Extracting "+path+" to "+conf.LocalPath)
		err = keygen.ExtractPackage(path, conf.LocalPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		conf.LocalFile = candidate;
		fmt.Println("Successfully updated to "+candidate.Version)
		return
	}
}
