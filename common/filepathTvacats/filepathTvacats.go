package filepathTvacats

import (
	"os"
	"path/filepath"
	"rename-tool/setting/config"
)

// 获取应用数据目录
func GetAppDataDir() string {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	// 在用户主目录下创建应用目录
	appDir := filepath.Join(homeDir, "."+config.AppName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "."
	}
	return appDir
}
