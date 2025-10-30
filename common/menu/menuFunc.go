package menu

import (
	"rename-tool/common/theme"
	"rename-tool/setting/global"
	"rename-tool/setting/i18n"
	"rename-tool/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// buildHeader constructs the top bar with language selector and admin status.
func buildHeader() fyne.CanvasObject {
	adminStatus := widget.NewLabel("")
	UpdateAdminStatusLabel(adminStatus)

	return container.NewHBox(
		i18n.LangSelect(),
		layout.NewSpacer(),
		adminStatus,
		layout.NewSpacer(),
		widget.NewLabel(buttonTr("AppName")),
	)
}

// buildMainImage loads and returns the themed image (with fallback).
func buildMainImage() *canvas.Image {
	imgResource := theme.LoadImage("cat.png")
	var image *canvas.Image
	if imgResource != nil {
		image = canvas.NewImageFromResource(imgResource)
	} else {
		image = canvas.NewImageFromFile("")
		logEvent("THEME ERROR", "loadThemeError")
	}
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(250, 380))
	return image
}

func buildBody() fyne.CanvasObject {
	image := buildMainImage()

	// 按钮列表
	buttons := []struct {
		text   string
		action func()
	}{
		{buttonTr("sequenceRename"), utils.ShowBatchRenameNormal},
		{buttonTr("extensionModify"), utils.ShowChangeExtension},
		{buttonTr("toUpper"), func() { utils.ShowRenameToCase("upper") }},
		{buttonTr("toLower"), func() { utils.ShowRenameToCase("lower") }},
		{buttonTr("titlecase"), func() { utils.ShowRenameToCase("title") }},
		{buttonTr("camel"), func() { utils.ShowRenameToCase("camel") }},
		{buttonTr("insertLetter"), utils.ShowInsertCharRename},
		{buttonTr("deleteLetter"), utils.ShowDeleteCharRename},
		{buttonTr("regexReplace"), utils.ShowRegexReplace},
		{buttonTr("undoRename"), utils.UndoRename},
		{buttonTr("logSaved"), utils.SaveLogs},
		{buttonTr("exit"), func() { global.MyApp.Quit() }},
	}

	// 构建按钮网格（两列，每个按钮高度固定）
	buttonGrid := buildFixedButtonGrid(buttons)

	// 将按钮网格放入滚动容器，控制显示高度
	buttonScroll := container.NewVScroll(buttonGrid)
	buttonScroll.SetMinSize(fyne.NewSize(0, 400)) // 可调，显示区域高度

	// 左图右按钮布局
	mainContent := container.NewBorder(nil, nil, image, nil, buttonScroll)

	// 整体居中并设置背景
	return theme.SetBackground(container.NewCenter(mainContent))
}

// buildFixedButtonGrid creates a two-column button layout with fixed height
func buildFixedButtonGrid(buttons []struct {
	text   string
	action func()
}) fyne.CanvasObject {
	var rows []fyne.CanvasObject

	makeTextBtn := func(text string, onTap func()) fyne.CanvasObject {
		btn := widget.NewButton(text, onTap)
		btn.Importance = widget.LowImportance
		btn.Resize(fyne.NewSize(150, 30)) // 固定大小，可调整
		return btn
	}

	for i := 0; i < len(buttons); i += 2 {
		row := []fyne.CanvasObject{makeTextBtn(buttons[i].text, buttons[i].action)}

		if i+1 < len(buttons) {
			row = append(row, makeTextBtn(buttons[i+1].text, buttons[i+1].action))
		} else {
			row = append(row, layout.NewSpacer()) // 占位
		}

		rows = append(rows, container.NewHBox(row...))
	}

	return container.NewVBox(rows...)
}
