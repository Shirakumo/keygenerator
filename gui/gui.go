package gui

import (
	"keygenerator/keygen"
	"keygenerator/config"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"log"
	"time"
)

func saveConfig(conf config.Config){
	err := config.WriteDefault(conf)
	if err != nil {
		log.Fatal(err)
	}
}

func showEntries(conf config.Config, w fyne.Window){
	key, _ := keygen.ParseKeyURL(conf.KeyURL)
	log.Print("Fetching entries for "+key.Code)
	files, err := keygen.FetchKeyFiles(key)
	if err != nil {
		dialog.ShowError(err, w)
		log.Fatal(err)
	} else {
		conf.LastChecked = time.Now().Unix()
		saveConfig(conf)

		list := widget.NewList(
			func() int {
				return len(files)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("template")
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(files[i].Version)
			})
		w.SetContent(list)
	}
}

func Main(conf config.Config){
	a := app.New()
	w := a.NewWindow("Keygen Updater")

	input := widget.NewEntry()
	input.SetPlaceHolder("https://keygen.tymoon.eu/access/...")
	input.Validator = func(str string) error {
		_, err := keygen.ParseKeyURL(str)
		return err
	}
	input.OnSubmitted = func(str string) {
		if input.Validator(str) == nil {
			input.Disable()
			defer input.Enable()
			conf.KeyURL = input.Text
			saveConfig(conf)
			showEntries(conf, w)
		}
	}
	
	w.SetContent(container.NewVBox(
		widget.NewLabel("Enter your Key URL:"),
		input,
	))
	w.ShowAndRun()
}
