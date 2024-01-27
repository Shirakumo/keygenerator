package main

import (
	"keygenerator/config"
	"keygenerator/keygen"
	"keygenerator/tui"
	"keygenerator/gui"
	"path/filepath"
	"flag"
	"log"
)

func autoUpdate(conf *config.Config) {
	if conf.KeyURL == "" { log.Panic("No KeyURL set.") }
	key, err := keygen.ParseKeyURL(conf.KeyURL)
	if err != nil { log.Panic(err) }
	files, err := keygen.FetchKeyFiles(key)
	if err != nil { log.Panic(err) }
	candidate := keygen.FindMatchingOSFile(files)
	if conf.LocalFile != nil {
		candidate = keygen.FindUpdatedFile(files, conf.LocalFile)
	}
	if candidate == nil {
		log.Print("Already up to date.")
		return
	}
	path := filepath.Join(conf.PackagePath, candidate.Filename)
	err = keygen.DownloadPackage(candidate, path)
	if err != nil { log.Panic(err) }
	err = keygen.ExtractPackage(path, conf.LocalPath)
	if err != nil { log.Panic(err) }
	log.Print("Successfully updated to "+candidate.Version)
}

func main() {
	configFile := flag.String("config", "", "The config file to read (defaults to .key in the binary directory)")
	keyURL := flag.String("key", "", "The Key URL to access")
	localPath := flag.String("path", "", "The directory to extract things to")
	packagePath := flag.String("tmp", "", "The directory to store packages in")
	useGui := flag.Bool("gui", true, "Whether to use the GUI or not")
	useAuto := flag.Bool("auto", false, "Whether to perform a one-shot auto-update. Ipmlies --gui false")
	flag.Parse()

	var conf config.Config
	var err error
	if *configFile == "" {
		conf, err = config.ReadDefault()
		*configFile = config.DefaultPath()
		if err != nil {
			log.Print(err)
		}
	} else {
		conf, err = config.Read(*configFile)
		if err != nil {
			log.Print(err)
			return
		}
	}

	if *keyURL != "" { conf.KeyURL = *keyURL }
	if *localPath != "" { conf.LocalPath = *localPath }
	if *packagePath != "" { conf.PackagePath = *packagePath }

	if *useAuto {
		autoUpdate(&conf)
	} else if *useGui {
		gui.Main(&conf)
	} else {
		tui.Main(&conf)
	}

	err = config.Write(conf, *configFile)
	if err != nil {
		log.Print(err)
	}
}

