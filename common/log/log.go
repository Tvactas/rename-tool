package log

import (
	"os"
	"path/filepath"
	"rename-tool/common/dirpath"
	"rename-tool/setting/config"
)

// 获取日志文件路径
func GetLogPath() string {
	appDir := dirpath.GetAppDataDir()
	logDir := filepath.Join(appDir, config.LogDir)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return filepath.Join(appDir, "rename.log")
	}
	return filepath.Join(logDir, "rename.log")
}
