package main

import (
	"keygenerator/config"
	"keygenerator/tui"
	"log"
)

func main() {
	conf, err := config.ReadDefault()
	if err != nil {
		log.Print(err)
	}

	tui.Main(conf)
}
