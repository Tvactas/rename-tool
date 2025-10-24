package dirpath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"rename-tool/common/applog"
	"rename-tool/common/fs"
	"rename-tool/setting/global"
	"rename-tool/setting/i18n"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

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
		applog.Logger.Printf("[PATH ERROR]  %s: %v", i18n.Tr("get_current_dir_failed"), err)

		return "" // 明确表明失败
	}
	return dir
}

// truncatePath 截断路径，超出长度用省略号显示，保持固定长度
func truncatePath(path string, maxLength int) string {
	if len(path) <= maxLength {
		// 如果路径长度小于等于最大长度，用空格填充到固定长度
		return path + strings.Repeat(" ", maxLength-len(path))
	}
	
	if maxLength <= 3 {
		return strings.Repeat(".", maxLength)
	}
	
	// 保留前面的部分，用省略号连接，确保总长度固定
	return path[:maxLength-3] + "..."
}

// 创建目录选择器组件
func CreateDirSelector(win fyne.Window, onDirChanged func()) fyne.CanvasObject {
	label := widget.NewLabel(i18n.Tr("dir") + ": " + truncatePath(global.SelectedDir, 50))
	button := widget.NewButton(i18n.Tr("select_dir"), func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				applog.Logger.Printf("[PATH ERROR]  %s: %v", i18n.Tr("folder_open_error"), err)

				return
			}
			if uri != nil {
				global.SelectedDir = uri.Path()
				label.SetText(i18n.Tr("dir") + ": " + truncatePath(global.SelectedDir, 50))
				if onDirChanged != nil {
					onDirChanged()
				}
			}
		}, win).Show()
	})

	return container.NewHBox(label, button)
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
