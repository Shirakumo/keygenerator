package main

import (
	"keygenerator/config"
	"keygenerator/tui"
	"keygenerator/gui"
	"fmt"
	"os"
	flag "github.com/spf13/pflag"
)

func main() {
	configFile := flag.String("config", "", "The config file to read")
	keyURL := flag.String("key", "", "The key URL to access")
	localPath := flag.String("path", "", "The directory to extract things to")
	packagePath := flag.String("tmp", "", "The directory to store packages in")
	flag.Parse()

	var conf config.Config
	var err error
	if *configFile == "" {
		conf, _ = config.ReadDefault()
		*configFile = config.DefaultPath()
	} else {
		conf, err = config.Read(*configFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if os.Getenv("KEY") != "" { conf.KeyURL = os.Getenv("KEY") }
	if *keyURL != "" { conf.KeyURL = *keyURL }
	if *localPath != "" { conf.LocalPath = *localPath }
	if *packagePath != "" { conf.PackagePath = *packagePath }

	action := "gui"
	if len(flag.Args()) >= 1 {
		action = flag.Args()[0]
	}

	switch {
	case action == "update": tui.Update(&conf)
	case action == "list": tui.List(&conf)
	case action == "tui": tui.Main(&conf)
	case action == "gui": gui.Main(&conf)
	case action == "help":
		fmt.Println(`Manage an installation of files from a Keygen key.

Usage:
      keygenerator [action] flag...

Actions:
      update            Update the application from the key
      list              List available files from the key
      tui               Show an interactive terminal UI
      gui               Show an interactive graphical UI (default)
      help              Show this help

Flags:`)
		flag.PrintDefaults()
		fmt.Println(`
Environment Variables:
      KEY               Alternative way to pass the key URL

By default the updater will download packages to a temporary directory
and extract them into the same directory it is residing in. It will
also try to read a config file in the same directory named ".key"
where it stores the key URL, the corresponding remote file, and so
on. Once the key URL has been set once, it will be remembered for
future sessions that way and does not need to be set again.

You can access the source code of this updater application at the
following URL:

     https://shirakumo.org/projects/keygenerator

You can access precompiled binaries of this updater application at the
following URL:

     https://shirakumo.org/projects/keygenerator/releases/latest

- (c) 2024 Yukari Hafner, Shirakumo`)
	}

	err = config.Write(conf, *configFile)
	if err != nil {
		fmt.Print(err)
	}
}

