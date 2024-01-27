package main

import (
	"keygenerator/config"
	"keygenerator/tui"
	"keygenerator/gui"
	"flag"
	"log"
)

func main() {
	configFile := flag.String("config", "", "The config file to read (defaults to .key in the binary directory)")
	keyURL := flag.String("key", "", "The Key URL to access")
	localPath := flag.String("path", "", "The directory to extract things to")
	packagePath := flag.String("tmp", "", "The directory to store packages in")
	useGui := flag.Bool("gui", true, "Whether to use the GUI or not")
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

	if *useGui {
		gui.Main(&conf)
	} else {
		tui.Main(&conf)
	}

	err = config.Write(conf, *configFile)
	if err != nil {
		log.Print(err)
	}
}

