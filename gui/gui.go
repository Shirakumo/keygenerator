package gui

import (
	"os"
	"keygenerator/keygen"
	"path/filepath"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var a fyne.App = app.New()
var w fyne.Window = a.NewWindow("Keygen Updater")

func ShowKeyPrompt() string {
	input := widget.NewEntry()
	input.SetPlaceHolder("https://keygen.tymoon.eu/access/...")

	w.SetContent(container.NewVBox(
		widget.NewLabel("Please enter your Keygen Key URL"),
		input,
		widget.NewButton("OK", func() {
			w.Close()
		}),
	))
	return input.Text
}

func ShowFailure(err error) {
	w.SetContent(container.NewVBox(
		widget.NewLabel("An error occurred: "+err.Error()),
		widget.NewButton("OK", func() {
			w.Close()
		}),
	))
}

func ShowSuccess(str string) {
	w.SetContent(container.NewVBox(
		widget.NewLabel("Success: "+str),
		widget.NewButton("OK", func() {
			w.Close()
		}),
	))
}

func ShowNewVersionPrompt(file *keygen.File) string {
	path := ""
	input := widget.NewEntry()
	input.Text = filepath.Join(os.TempDir(), file.Filename)
	input.SetPlaceHolder("Package path...")
	
	w.SetContent(container.NewVBox(
		widget.NewLabel("A new version is available: "+file.Version),
		input,
		widget.NewButton("Update", func() {
			path = input.Text
			w.Close()
		}),
	))
	return path
}
