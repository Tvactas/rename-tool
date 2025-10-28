package preview

import (
	"rename-tool/setting/model"

	"fyne.io/fyne/v2"
)

// ShowPreviewWindow 显示预览窗口
func ShowPreviewWindow(parentWindow fyne.Window, files []string, config model.RenameConfig) {
	previewWindow := createPreviewWindow()
	previewList := createPreviewList(files, config)
	content := buildWindowContent(previewList, len(files), previewWindow)

	previewWindow.SetContent(content)
	previewWindow.Show()
}
