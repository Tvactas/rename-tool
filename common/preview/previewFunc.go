package preview

import (
	"fmt"
	"path/filepath"
	"rename-tool/common/pathgen"
	"rename-tool/setting/global"
	"rename-tool/setting/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// createPreviewWindow 创建预览窗口
func createPreviewWindow() fyne.Window {
	window := global.MyApp.NewWindow(buttonTr("preview"))
	window.Resize(fyne.NewSize(800, 600))
	window.SetFixedSize(false)
	window.SetCloseIntercept(func() {
		window.Close()
	})
	return window
}

// createPreviewList 创建预览列表
func createPreviewList(files []string, config model.RenameConfig) *widget.List {
	return widget.NewList(
		func() int { return len(files) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			displayPreviewItem(obj.(*widget.Label), files[id], config, id)
		},
	)
}

// displayPreviewItem 显示单个预览项
func displayPreviewItem(label *widget.Label, file string, config model.RenameConfig, id int) {
	_, oldName := filepath.Split(file)
	newPath, err := generateNewPath(file, config, id)

	if err != nil {
		label.SetText(fmt.Sprintf("%s → %s", oldName, err.Error()))
		return
	}

	_, newName := filepath.Split(newPath)
	label.SetText(fmt.Sprintf("%s → %s", oldName, newName))
}

// generateNewPath 根据配置生成新路径
func generateNewPath(file string, config model.RenameConfig, id int) (string, error) {
	counters := make(map[string]int)

	switch config.Type {
	case model.RenameTypeBatch:
		return pathgen.GenerateBatchRenamePath(file, config, id, counters)
	case model.RenameTypeExtension:
		return pathgen.GenerateExtensionRenamePath(file, config)
	case model.RenameTypeCase:
		return pathgen.GenerateCaseRenamePath(file, config)
	case model.RenameTypeInsertChar:
		return pathgen.GenerateInsertCharRenamePath(file, config)
	case model.RenameTypeReplace:
		return pathgen.GenerateReplaceRenamePath(file, config)
	case model.RenameTypeDeleteChar:
		return pathgen.GenerateDeleteCharRenamePath(file, config)
	default:
		return "", fmt.Errorf("unknown rename type")
	}
}

// buildWindowContent 构建窗口内容
func buildWindowContent(previewList *widget.List, fileCount int, window fyne.Window) *fyne.Container {
	topBar := createTopBar(fileCount)
	bottomBar := createBottomBar(window)

	return container.NewBorder(topBar, bottomBar, nil, nil, previewList)
}

// createTopBar 创建顶部栏
func createTopBar(fileCount int) *fyne.Container {
	title := widget.NewLabelWithStyle(buttonTr("preview"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	countLabel := widget.NewLabel(fmt.Sprintf(dialogTr("totalFiles")+": %d", fileCount))
	return container.NewHBox(title, layout.NewSpacer(), countLabel)
}

// createBottomBar 创建底部栏
func createBottomBar(window fyne.Window) *fyne.Container {
	closeBtn := widget.NewButton(dialogTr("confirm"), func() {
		window.Close()
	})
	return container.NewHBox(layout.NewSpacer(), closeBtn)
}
