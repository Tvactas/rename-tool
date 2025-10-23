package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"rename-tool/common/applog"
	"rename-tool/common/log"
	"rename-tool/common/ui"
	"rename-tool/setting/global"
	"rename-tool/setting/i18n"

	"fyne.io/fyne/v2/dialog"
)

// SaveLogs handles saving the operation logs
func SaveLogs() {
	if len(global.Logs) == 0 {
		dialog.ShowInformation(tr("info"), tr("no_operations_to_save"), global.MainWindow)
		return
	}

	var sb strings.Builder
	for _, log := range global.Logs {
		fmt.Fprintf(&sb, "%s > %s [%s]\n", log.Original, log.New, log.Time)
	}
	content := sb.String()

	dir := filepath.Dir(log.GetLogPath())
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		applog.Logger.Printf("[SAVE LOG ERROR] %v", err)
		dialog.ShowError(fmt.Errorf(tr("error_creating_directory")+": %v", err), global.MainWindow)
		return
	}

	tempPath := log.GetLogPath() + ".tmp"
	if err := os.WriteFile(tempPath, []byte(content), 0644); err != nil {
		applog.Logger.Printf("[SAVE LOG ERROR] %v", err)
		dialog.ShowError(fmt.Errorf(tr("error_saving_log")+": %v", err), global.MainWindow)
		return
	}

	// Windows 需要先删除目标文件
	_ = os.Remove(log.GetLogPath())
	if err := os.Rename(tempPath, log.GetLogPath()); err != nil {
		applog.Logger.Printf("[SAVE LOG ERROR] %v", err)
		os.Remove(tempPath)
		dialog.ShowError(fmt.Errorf(tr("error_saving_log")+": %v", err), global.MainWindow)
		return
	}

	message := fmt.Sprintf("%s %d %s %s", i18n.DialogTr("SuccessSaved"), len(global.Logs), i18n.DialogTr("files_count_with_path"), log.GetLogPath())
	ui.ShowWidePlainMessage(global.MainWindow, tr("success"), message)
}
