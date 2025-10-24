package appinit

import (
	"errors"
	"fmt"
	"rename-tool/common/dirpath"
	"rename-tool/common/recovery"

	"rename-tool/common/theme"
	"rename-tool/setting/global"
	"rename-tool/setting/i18n"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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
	// 设置错误处理
	defer recovery.RecoverPanic()

	// 初始化应用
	global.MyApp = app.NewWithID(config.AppID)
	if global.MyApp == nil {
		return errors.New(i18n.Tr("init_app_failed"))

	}

	// 创建主窗口
	global.MainWindow = global.MyApp.NewWindow(i18n.Tr("title"))

	if global.MainWindow == nil {
		return errors.New(i18n.Tr("create_main_window_failed"))
	}

	// 配置窗口
	global.MainWindow.Resize(config.WindowSize)
	if config.FixedSize {
		global.MainWindow.SetFixedSize(true)
	}
	global.MainWindow.SetMaster()

	// 初始化目录
	if err := initializeDirectories(); err != nil {
		return fmt.Errorf("%s: %v", i18n.Tr("init_dir_failed"), err)
	}
	// 设置主题
	global.MyApp.Settings().SetTheme(&theme.MainTheme{})

	return nil
}

// initializeDirectories 初始化目录
func initializeDirectories() error {
	// 获取当前目录
	global.CurrentDir = dirpath.GetCurrentDir()
	if global.CurrentDir == "" {
		return fmt.Errorf("%s", i18n.Tr("failed_to_get_current_directory"))
	}
	global.SelectedDir = global.CurrentDir
	return nil
}

// RunApp 运行应用程序
func RunApp() {
	global.MainWindow.ShowAndRun()
}
