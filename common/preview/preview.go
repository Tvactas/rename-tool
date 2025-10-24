package preview

import (
	"fmt"
	"path/filepath"

	"rename-tool/common/pathgen"
	"rename-tool/setting/global"
	"rename-tool/setting/i18n"
	"rename-tool/setting/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// ShowPreviewWindow 显示预览窗口
func ShowPreviewWindow(parentWindow fyne.Window, files []string, config model.RenameConfig) {
	// 创建预览窗口
	previewWindow := global.MyApp.NewWindow(tr("preview"))
	previewWindow.Resize(fyne.NewSize(800, 600))
	previewWindow.SetFixedSize(false)

	// 设置窗口关闭事件
	previewWindow.SetCloseIntercept(func() {
		previewWindow.Close()
	})

	// 创建预览列表
	previewList := widget.NewList(
		func() int { return len(files) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			file := files[id]
			_, oldName := filepath.Split(file)

			var newPath string
			var err error

			// 根据重命名类型生成新路径
			switch config.Type {
			case model.RenameTypeBatch:
				// 为批量重命名创建计数器
				counters := make(map[string]int)
				newPath, err = pathgen.GenerateBatchRenamePath(file, config, id, counters)
			case model.RenameTypeExtension:
				newPath, err = pathgen.GenerateExtensionRenamePath(file, config)
			case model.RenameTypeCase:
				newPath, err = pathgen.GenerateCaseRenamePath(file, config)
			case model.RenameTypeInsertChar:
				newPath, err = pathgen.GenerateInsertCharRenamePath(file, config)
			case model.RenameTypeReplace:
				newPath, err = pathgen.GenerateReplaceRenamePath(file, config)
			case model.RenameTypeDeleteChar:
				newPath, err = pathgen.GenerateDeleteCharRenamePath(file, config)
			}

			if err != nil {
				obj.(*widget.Label).SetText(fmt.Sprintf("%s → %s", oldName, err.Error()))
				return
			}

			_, newName := filepath.Split(newPath)
			obj.(*widget.Label).SetText(fmt.Sprintf("%s → %s", oldName, newName))
		},
	)

	// 创建标题
	title := widget.NewLabelWithStyle(tr("preview"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// 创建文件计数标签
	countLabel := widget.NewLabel(fmt.Sprintf(tr("total_files")+": %d", len(files)))

	// 创建关闭按钮
	closeBtn := widget.NewButton(tr("close"), func() {
		previewWindow.Close()
	})

	// 创建布局
	topBar := container.NewHBox(title, layout.NewSpacer(), countLabel)
	bottomBar := container.NewHBox(layout.NewSpacer(), closeBtn)

	content := container.NewBorder(
		topBar,
		bottomBar,
		nil,
		nil,
		previewList,
	)

	previewWindow.SetContent(content)
	previewWindow.Show()
}

// tr 函数用于国际化
func tr(key string) string {
	return i18n.Tr(key)
}
