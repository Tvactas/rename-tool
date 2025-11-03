package dirpath

import (
	"errors"
	"os"

	"rename-tool/setting/global"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// GetFiles 获取指定目录下符合格式的所有文件
func GetFiles(root string, formats []string) ([]string, error) {
	var files []string

	err := walkDirFiltered(root, formats, func(path string, _ os.FileInfo) {
		files = append(files, path)
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}

// GetShortestFilenameLength 返回目录中文件名的最短长度（忽略子目录）
func GetShortestFilenameLength(dir string) (int, error) {
	minLen := -1

	err := walkDirFiltered(dir, nil, func(path string, info os.FileInfo) {
		nameLen := len(info.Name())
		if minLen == -1 || nameLen < minLen {
			minLen = nameLen
		}
	})

	if err != nil {
		return 0, err
	}
	if minLen == -1 {
		return 0, errors.New(textTr("noFiles"))
	}
	return minLen, nil
}

// GetCurrentDir 返回当前工作目录
func GetCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		logEvent("PATH ERROR", "failGetCurrentDir", err)
		return ""
	}
	return dir
}

// CreateDirSelector 创建目录选择器组件（Fyne UI）
func CreateDirSelector(win fyne.Window, onDirChanged func()) fyne.CanvasObject {
	label := widget.NewLabel(buttonTr("dir") + ": " + truncatePathMiddle(global.SelectedDir, 50))
	button := widget.NewButton(buttonTr("selectDir"), func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				logEvent("PATH ERROR", "folderOpenError", err)
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
