package utils

import (
	"fmt"
	"os"
	"time"

	"rename-tool/common/filestatus"
	"rename-tool/setting/global"

	"fyne.io/fyne/v2/dialog"
)

// UndoRename handles the undo operation for file renaming
func UndoRename() {
	if len(global.Logs) == 0 {
		dialog.ShowInformation(dialogTr("warning"), tr("no_undo_operations"), global.MainWindow)
		return
	}

	busyFiles := []string{} // Record busy files
	successCount := 0

	for i := len(global.Logs) - 1; i >= 0; i-- {
		log := global.Logs[i]
		if _, err := os.Stat(log.New); err == nil {
			if err := os.Rename(log.New, log.Original); err == nil {
				successCount++
				// Remove the undone record from logs
				global.Logs = append(global.Logs[:i], global.Logs[i+1:]...)
				global.Logs = append(global.Logs, global.RenameLog{
					Original: log.Original,
					New:      log.New,
					Time:     time.Now().Format("2006-01-02 15:04:05"),
				})
			} else if filestatus.IsFileBusyError(err) {
				busyFiles = append(busyFiles, log.New)
			}
		}
	}

	if len(busyFiles) > 0 {
		filestatus.ShowBusyFilesDialog(global.MainWindow, busyFiles)
	} else {
		dialog.ShowInformation(dialogTr("success"), fmt.Sprintf(tr("undo_success"), successCount), global.MainWindow)
	}
}
