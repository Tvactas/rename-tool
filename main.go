package main

//power by Tvacats
import (
	"rename-tool/common/appinit"
	"rename-tool/common/applog"
	"rename-tool/common/menu"
	"rename-tool/common/theme"
	"rename-tool/setting/global"
	"rename-tool/setting/i18n"
)

func main() {

	// Initialize application with default configuration
	if err := appinit.InitializeApp(appinit.DefaultConfig()); err != nil {
		applog.Logger.Printf("[INIT ERROR]  %s %v", i18n.LogTr("initAppError"), err)
		return
	}

	// Show main menu
	menu.ShowMainMenu()

	// Run application
	appinit.RunApp()
}

func init() {
	applog.InitLogger("tvacats_rename.log")
	// Initialize resource loader
	theme.Init() // Initialize resource loader

	// Set language change callback
	i18n.GetManager().SetOnLangChange(func() {
		// Refresh main window
		if global.MainWindow != nil {
			// Save current window size
			size := global.MainWindow.Canvas().Size()
			// Recreate main menu
			menu.ShowMainMenu()
			// Restore window size
			global.MainWindow.Resize(size)
		}
	})

}
