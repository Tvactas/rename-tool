package appinit

import (
	"errors"
	"rename-tool/common/dirpath"
	"rename-tool/common/theme"
	"rename-tool/setting/global"

	"fyne.io/fyne/v2/app"
)

func initApp(config AppConfig) error {
	global.MyApp = app.NewWithID(config.AppID)
	if global.MyApp == nil {
		return errors.New(textTr("failInitAppID"))
	}
	return nil
}

func initMainWindow(config AppConfig) error {
	global.MainWindow = global.MyApp.NewWindow(buttonTr("AppName"))
	if global.MainWindow == nil {
		return errors.New(textTr("failCreateMainWindow"))
	}

	global.MainWindow.Resize(config.WindowSize)
	global.MainWindow.CenterOnScreen()
	global.MainWindow.SetMaster()

	if config.FixedSize {
		global.MainWindow.SetFixedSize(true)
	}

	return nil
}

// initializeDirectories 初始化目录
func initializeDirectories() error {
	global.CurrentDir = dirpath.GetCurrentDir()
	if global.CurrentDir == "" {
		return errors.New(textTr("failGetCurrentDir"))

	}
	global.SelectedDir = global.CurrentDir
	return nil
}
func initTheme() {
	global.MyApp.Settings().SetTheme(&theme.MainTheme{})
}
