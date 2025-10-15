package global

import (
	"fyne.io/fyne/v2"
)

type RenameLog struct {
	Original string
	New      string
	Time     string
}

var (
	MyApp       fyne.App
	MainWindow  fyne.Window
	CurrentDir  string
	Logs        []RenameLog
	SelectedDir string
	Lang        string = "en"
)
