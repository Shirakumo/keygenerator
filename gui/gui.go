package gui

import (
	"os"
	"keygenerator/keygen"
	"path/filepath"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowKeyPrompt() string {
	a := app.New()
	w := a.NewWindow("Keygen Updater")
	
	input := widget.NewEntry()
	input.SetPlaceHolder("https://keygen.tymoon.eu/access/...")

	w.SetContent(container.NewVBox(
		widget.NewLabel("Please enter your Keygen Key URL"),
		input,
		widget.NewButton("OK", func() {
			w.Close()
		}),
	))
	w.ShowAndRun()
	return input.Text
}

func ShowFailure(err error) {
	a := app.New()
	w := a.NewWindow("Keygen Updater")
	
	w.SetContent(container.NewVBox(
		widget.NewLabel("An error occurred: "+err.Error()),
		widget.NewButton("OK", func() {
			w.Close()
		}),
	))
	w.ShowAndRun()
}

func ShowSuccess(str string) {
	a := app.New()
	w := a.NewWindow("Keygen Updater")
	
	w.SetContent(container.NewVBox(
		widget.NewLabel("Success: "+str),
		widget.NewButton("OK", func() {
			w.Close()
		}),
	))
	w.ShowAndRun()
}

func ShowNewVersionPrompt(file *keygen.File) string {
	path := ""
	a := app.New()
	w := a.NewWindow("Keygen Updater")

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
	w.ShowAndRun()
	return path
}
