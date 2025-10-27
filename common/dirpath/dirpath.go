package dirpath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	formatsMap := mapExt(formats)

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
		return nil, fmt.Errorf("%s: %w", textTr("FailReadFiles"), err)
	}

	return files, nil
}

func GetCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		logEvent("PATH ERROR", "FailGetCurrentDir", err)
		return ""
	}
	return dir
}

// 创建目录选择器组件
func CreateDirSelector(win fyne.Window, onDirChanged func()) fyne.CanvasObject {
	label := widget.NewLabel(buttonTr("dir") + ": " + truncatePathMiddle(global.SelectedDir, 50))
	button := widget.NewButton(buttonTr("SelectDir"), func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				logEvent("PATH ERROR", "folder_open_error", err)

				return
			}
			if uri != nil {
				global.SelectedDir = uri.Path()
				label.SetText(buttonTr("dir") + ": " + truncatePathMiddle(global.SelectedDir, 50))
				if onDirChanged != nil {
					onDirChanged()
				}
			}
		}, win).Show()
	})

	return container.NewHBox(label, button)
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
