package appinit

import (
	"rename-tool/common/recovery"

	"rename-tool/setting/global"

	"fyne.io/fyne/v2"
)

// AppConfig 应用程序配置
type AppConfig struct {
	AppID      string
	WindowSize fyne.Size
	FixedSize  bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() AppConfig {
	return AppConfig{
		AppID:      "com.tencats.renametool",
		WindowSize: fyne.NewSize(600, 400),
		FixedSize:  false,
	}
}

// InitializeApp 初始化应用程序
func InitializeApp(config AppConfig) error {
	defer recovery.RecoverPanic()

	if err := initAppID(config); err != nil {
		return err
	}

	if err := initMainWindow(config); err != nil {
		return err
	}

	if err := initializeDirectories(); err != nil {
		return err
	}

	initTheme()

	return nil
}

// RunApp 运行应用程序
func RunApp() {
	global.MainWindow.ShowAndRun()
}
