package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"rename-tool/common/log"
	"rename-tool/setting/global"

	"fyne.io/fyne/v2/dialog"
)

// SaveLogs handles saving the operation logs
func SaveLogs() {
	if len(global.Logs) == 0 {
		dialog.ShowInformation(tr("info"), tr("no_operations_to_save"), global.MainWindow)
		return
	}

	content := ""
	for _, log := range global.Logs {
		content += fmt.Sprintf("%s > %s [%s]\n", log.Original, log.New, log.Time)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(log.GetLogPath()), os.ModePerm); err != nil {
		dialog.ShowError(fmt.Errorf(tr("error_creating_directory")+": %v", err), global.MainWindow)
		return
	}

	// Use temporary file for writing
	tempPath := log.GetLogPath() + ".tmp"
	if err := os.WriteFile(tempPath, []byte(content), 0644); err != nil {
		dialog.ShowError(fmt.Errorf(tr("error_saving_log")+": %v", err), global.MainWindow)
		return
	}

	// Atomically rename temporary file
	if err := os.Rename(tempPath, log.GetLogPath()); err != nil {
		// Clean up temporary file
		os.Remove(tempPath)
		dialog.ShowError(fmt.Errorf(tr("error_saving_log")+": %v", err), global.MainWindow)
		return
	}

	dialog.ShowInformation(tr("success"), fmt.Sprintf(tr("success_saved")+" "+tr("logs_count")+" "+tr("files_count_with_path"), len(global.Logs), log.GetLogPath()), global.MainWindow)
}
