package main

import (
	"keygenerator/keygen"
	"keygenerator/config"
	"keygenerator/gui"
	"errors"
	"time"
	"log"
)

func saveConfig(conf config.Config){
	err := config.WriteDefault(conf)
	if err != nil {
		log.Fatal(err)
		gui.ShowFailure(err)
	}
}

func main() {
	conf, err := config.ReadDefault()
	if err != nil {
		log.Print(err)
	}
	
	if conf.KeyURL == "" {
		conf.KeyURL = gui.ShowKeyPrompt()
		if conf.KeyURL == "" {
			return
		}
		saveConfig(conf)
	}

	key, err := keygen.ParseKeyURL(conf.KeyURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	files, err := keygen.FetchKeyFiles(key)
	if err != nil {
		log.Fatal(err)
	}
	conf.LastChecked = time.Now().Unix()
	saveConfig(conf)

	file := keygen.FindMatchingOSFile(files)
	if conf.LocalFile != nil {
		file = keygen.FindUpdatedFile(files, conf.LocalFile)
	}
	
	if file == nil {
		log.Fatal("No File Matched")
		gui.ShowFailure(errors.New("No new updates found."))
	} else {
		filepath := gui.ShowNewVersionPrompt(file)
		if filepath != "" {
			log.Print("Downloading "+file.URL+" to "+filepath)
			err = keygen.DownloadPackage(file, filepath)
			if err != nil {
				log.Fatal(err)
				gui.ShowFailure(err)
				return
			}

			log.Print("Extracting "+filepath+" to "+conf.LocalPath)
			err = keygen.ExtractPackage(filepath, conf.LocalPath)
			if err != nil {
				log.Fatal(err)
				gui.ShowFailure(err)
				return
			}

			conf.LocalFile = file;
			saveConfig(conf)
			gui.ShowSuccess("Successfully updated to "+file.Version)
		}
	}
}
