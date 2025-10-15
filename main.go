package main

//power by Tvacats
import (
	"embed"
	"fmt"

	"rename-tool/backend/common/appinit"
	"rename-tool/backend/common/log"
	"rename-tool/backend/common/menu"
	"rename-tool/backend/common/theme"
	"rename-tool/backend/setting/global"
	"rename-tool/backend/setting/i18n"
)

//go:embed src/font/* src/img/*
var resourceFS embed.FS

func main() {
	// Initialize application with default configuration
	if err := appinit.InitializeApp(appinit.DefaultConfig()); err != nil {
		log.LogError(fmt.Errorf("failed to initialize application: %v", err))
		return
	}

	// Show main menu
	menu.ShowMainMenu()

	// Run application
	appinit.RunApp()
}

func init() {
	// Initialize resource loader
	theme.SetFontFS(resourceFS) // Set font file system
	theme.Init()                // Initialize resource loader

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

	// List all embedded files
	files, err := theme.ReadDir(".")
	if err != nil {
		log.LogError(fmt.Errorf("failed to read embedded files: %v", err))
		return
	}
	for _, file := range files {
		log.LogError(fmt.Errorf("embedded file: %s", file.Name()))
	}
}
