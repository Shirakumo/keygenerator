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

	var key *keygen.Key
	for true {
		if conf.KeyURL == "" {
			conf.KeyURL = gui.ShowKeyPrompt()
			if conf.KeyURL == "" {
				return
			}
			saveConfig(conf)
		}

		key, err = keygen.ParseKeyURL(conf.KeyURL)
		if err != nil {
			log.Fatal(err)
			gui.ShowFailure(err)
		} else {
			break
		}
	}

	log.Print("Fetching entries for "+key.Code)
	files, err := keygen.FetchKeyFiles(key)
	if err != nil {
		log.Fatal(err)
	}
	conf.LastChecked = time.Now().Unix()
	saveConfig(conf)

	log.Print("Checking for new updates...")
	file := keygen.FindMatchingOSFile(files)
	if conf.LocalFile != nil {
		file = keygen.FindUpdatedFile(files, conf.LocalFile)
	}
	
	if file == nil {
		log.Print("No file matched")
		gui.ShowFailure(errors.New("No new updates found."))
	} else {
		log.Print("Found: "+file.Version)
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
