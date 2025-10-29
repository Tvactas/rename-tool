package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"rename-tool/common/applog"
	"rename-tool/common/ui"
	"rename-tool/setting/global"
	"rename-tool/setting/i18n"
)

// SaveLogs handles saving the operation logs
func SaveLogs() {
	if len(global.Logs) == 0 {
		warningDiaLog(global.MainWindow, dialogTr("noLogSaved"))
		return
	}

	var sb strings.Builder
	for _, log := range global.Logs {
		fmt.Fprintf(&sb, "%s > %s [%s]\n", log.Original, log.New, log.Time)
	}
	content := sb.String()

	dir := filepath.Dir(applog.GetLogPath())
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		applog.Logger.Printf("[SAVE LOG ERROR] %v", err)
		errorDiaLog(global.MainWindow, fmt.Sprintf(tr("error_creating_directory")+": %v", err))

		return
	}

	tempPath := applog.GetLogPath() + ".tmp"
	if err := os.WriteFile(tempPath, []byte(content), 0644); err != nil {
		applog.Logger.Printf("[SAVE LOG ERROR] %v", err)
		errorDiaLog(global.MainWindow, fmt.Sprintf(tr("error_saving_log")+": %v", err))
		return
	}

	// Windows 需要先删除目标文件
	_ = os.Remove(applog.GetLogPath())
	if err := os.Rename(tempPath, applog.GetLogPath()); err != nil {
		applog.Logger.Printf("[SAVE LOG ERROR] %v", err)
		os.Remove(tempPath)
		errorDiaLog(global.MainWindow, fmt.Sprintf(tr("error_saving_log")+": %v", err))
		return
	}

	message := fmt.Sprintf("%d %s %s", len(global.Logs), i18n.DialogTr("successSavedTo"), applog.GetLogPath())
	ui.ShowWidePlainMessage(global.MainWindow, dialogTr("success"), message)
}
