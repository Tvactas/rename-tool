package dirpath

import (
	"fmt"
	"os"
	"path/filepath"

	"rename-tool/common/fs"

	"rename-tool/setting/config"
	"rename-tool/setting/i18n"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"rename-tool/setting/global"
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

func GetFiles(dir string, formats []string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if fs.IsFileBusyError(err) {
				return nil
			}
			return err
		}
		if !info.IsDir() {
			if len(formats) == 0 {
				files = append(files, path)
			} else {
				ext := strings.ToLower(filepath.Ext(path))
				for _, format := range formats {
					if ext == format {
						files = append(files, path)
						break
					}
				}
			}
		}
		return nil
	})
	return files, err
}

// GetCurrentDir 获取当前工作目录
func GetCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get current directory: %v\n", err)
		return "."
	}
	return dir
}

// CreateDirSelector 创建目录选择组件
func CreateDirSelector(window fyne.Window) fyne.CanvasObject {
	dirLabel := widget.NewLabel(tr("dir") + ": " + global.SelectedDir)
	dirBtn := widget.NewButton(tr("select_dir"), func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				global.SelectedDir = uri.Path()
				// 替换"父母"为".."
				global.SelectedDir = strings.Replace(global.SelectedDir, "父母", "..", -1)
				dirLabel.SetText(tr("dir") + ": " + global.SelectedDir)
			}
		}, window).Show()
	})
	return container.NewHBox(dirLabel, dirBtn)
}

// 修改tr函数，使用i18n包的Tr函数
func tr(key string) string {
	return i18n.Tr(key)
}
