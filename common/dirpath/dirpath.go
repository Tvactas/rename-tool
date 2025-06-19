package dirpath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"rename-tool/common/fs"
	"rename-tool/setting/config"
	"rename-tool/setting/global"
	"rename-tool/setting/i18n"
)

// 获取应用数据目录 ~/.AppName/
func GetAppDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		printError(i18n.Tr("get_home_dir_failed"), err)
		return "."
	}

	appDir := filepath.Join(homeDir, "."+config.AppName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		printError(i18n.Tr("create_app_dir_failed"), err)
		return "."
	}
	return appDir
}

// 获取指定目录下的符合格式的所有文件
func GetFiles(root string, formats []string) ([]string, error) {
	var files []string
	formatsMap := toExtMap(formats)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if fs.IsFileBusyError(err) {
				return nil // 忽略文件占用错误
			}
			return errors.New(i18n.Tr("walk_error") + ": " + err.Error())
		}
		if !info.IsDir() && (len(formatsMap) == 0 || formatsMap[strings.ToLower(filepath.Ext(path))]) {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Tr("read_files_failed"), err)
	}

	return files, nil
}

// 获取当前工作目录
func GetCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		printError(i18n.Tr("get_current_dir_failed"), err)
		return "" // 明确表明失败
	}
	return dir
}

// 创建目录选择器组件
func CreateDirSelector(win fyne.Window, onDirChanged func()) fyne.CanvasObject {
	label := widget.NewLabel(i18n.Tr("dir") + ": " + global.SelectedDir)
	button := widget.NewButton(i18n.Tr("select_dir"), func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				printError(i18n.Tr("folder_open_error"), err)
				return
			}
			if uri != nil {
				global.SelectedDir = uri.Path()
				label.SetText(i18n.Tr("dir") + ": " + global.SelectedDir)
				if onDirChanged != nil {
					onDirChanged()
				}
			}
		}, win).Show()
	})

	return container.NewHBox(label, button)
}

// 打印错误信息（用于调试日志）
func printError(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
}

// 将扩展名列表转换为map便于快速匹配
func toExtMap(formats []string) map[string]bool {
	m := make(map[string]bool)
	for _, ext := range formats {
		ext = strings.ToLower(ext)
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		m[ext] = true
	}
	return m
}

// GetShortestFilenameLength returns the length of the shortest filename in the given directory (ignores subdirectories)
func GetShortestFilenameLength(dir string) (int, error) {
	minLength := -1
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		nameLen := len(entry.Name())
		if minLength == -1 || nameLen < minLength {
			minLength = nameLen
		}
	}
	if minLength == -1 {
		return 0, errors.New("no files in directory")
	}
	return minLength, nil
}
