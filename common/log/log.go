package log

import (
	"fmt"
	"os"
	"path/filepath"
	"rename-tool/common/filepathTvacats"
	"rename-tool/setting/config"
	"time"
)

// 修改日志记录函数
func LogError(err error) {
	if err == nil {
		return
	}

	logPath := getErrorLogPath()
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprint("[", timestamp, "] ", err, "\n")
	if _, err := f.WriteString(logEntry); err != nil {
		// 记录写入错误，但不返回错误以避免循环
		fmt.Fprintf(os.Stderr, "Failed to write to error log: %v\n", err)
	}
}

// 获取错误日志文件路径
func getErrorLogPath() string {
	appDir := filepathTvacats.GetAppDataDir()
	logDir := filepath.Join(appDir, config.LogDir)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return filepath.Join(appDir, "error.log")
	}
	return filepath.Join(logDir, "error.log")
}
