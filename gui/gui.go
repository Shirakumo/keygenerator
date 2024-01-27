package gui

import (
	"keygenerator/keygen"
	"keygenerator/config"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"log"
	"time"
	"path/filepath"
)

type FileItem struct {
	widget.BaseWidget
	File *keygen.File
	Version *widget.Label
	Date *widget.Label
	Types *widget.Label
	Filename *widget.Label
}

func NewFileItem() *FileItem {
	item := &FileItem{
		File: nil,
		Version: widget.NewLabel(""),
		Date: widget.NewLabel(""),
		Types: widget.NewLabel(""),
		Filename: widget.NewLabel(""),
	}
	item.ExtendBaseWidget(item)
	return item
}

func (item *FileItem) SetFile(file *keygen.File) {
	item.File = file
	item.Version.SetText(file.Version)
	item.Date.SetText(time.Unix(file.LastModified, 0).Format("2006-01-02 15:04:05"))
	item.Types.SetText(file.Types[0])
	item.Filename.SetText(file.Filename)
}

func (item *FileItem) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewHBox(item.Version, item.Date, item.Types, item.Filename)
	return widget.NewSimpleRenderer(c)
}

func performDownload(file *keygen.File, conf *config.Config, w fyne.Window){
	progress := widget.NewProgressBar()
	progress.Min = 0.0
	progress.Max = 100.0
	w.SetContent(container.NewVBox(
		widget.NewLabel("Downloading ..."),
		progress,
	))

	path := filepath.Join(conf.PackagePath, file.Filename)
	log.Printf("Downloading %v to %v\n", file.URL, path)
	err := keygen.DownloadPackageProgress(file, path, func(prog float64){
		progress.SetValue(prog)
	})
	if err != nil {
		log.Print(err)
		dialog.ShowError(err, w)
		showEntries(conf, w)
		return
	}
	
	w.SetContent(container.NewVBox(
		widget.NewLabel("Extracting ..."),
		widget.NewProgressBarInfinite(),
	))

	log.Printf("Extracting %v to %v\n", path, conf.LocalPath)
	err = keygen.ExtractPackage(path, conf.LocalPath)
	if err != nil {
		log.Print(err)
		dialog.ShowError(err, w)
		showEntries(conf, w)
		return
	}

	conf.LocalFile = file;
	showEntries(conf, w)
	log.Printf("Successfully installed %v\n", file.Version)
	dialog.ShowInformation("Success", "Installed "+file.Version, w)
}

func showEntries(conf *config.Config, w fyne.Window){
	w.SetContent(container.NewVBox(
		widget.NewLabel("Loading ..."),
		widget.NewProgressBarInfinite(),
	))
	
	key, _ := keygen.ParseKeyURL(conf.KeyURL)
	files, err := keygen.FetchKeyFiles(key)
	if err != nil {
		log.Print(err)
		dialog.ShowError(err, w)
		showKeyEntry(conf, w)
	} else {
		conf.LastChecked = time.Now().Unix()

		candidate := keygen.FindMatchingOSFile(files)
		if conf.LocalFile != nil {
			candidate = keygen.FindUpdatedFile(files, conf.LocalFile)
		}

		list := widget.NewList(
			func() int {
				return len(files)
			},
			func() fyne.CanvasObject {
				return NewFileItem()
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*FileItem).SetFile(&files[i])
			})
		list.OnSelected = func(i widget.ListItemID) {
			candidate = &files[i]
		}
		
		w.SetContent(container.NewBorder(
			widget.NewLabel("Available Files:"),
			widget.NewButtonWithIcon("Install", theme.DownloadIcon(), func() {
				performDownload(candidate, conf, w)
			}),
			nil, nil, list,
		))
		for i := 0; i < len(files); i++ {
			if &files[i] == candidate {
				list.Select(i)
				break
			}
		}
	}
}

func showKeyEntry(conf *config.Config, w fyne.Window){
	input := widget.NewEntry()
	input.SetPlaceHolder("https://keygen.tymoon.eu/access/...")
	input.Text = conf.KeyURL
	input.Validator = func(str string) error {
		_, err := keygen.ParseKeyURL(str)
		return err
	}
	input.OnSubmitted = func(str string) {
		if input.Validator(str) == nil {
			input.Disable()
			conf.KeyURL = input.Text
			showEntries(conf, w)
		}
	}
	
	w.SetContent(container.NewVBox(
		widget.NewLabel("Enter your Key URL:"),
		input,
	))
}

func Main(conf *config.Config){
	a := app.New()
	w := a.NewWindow("Keygen Updater")

	if conf.KeyURL == "" {
		showKeyEntry(conf, w)
	} else {
		showEntries(conf, w)
	}
	w.ShowAndRun()
}
