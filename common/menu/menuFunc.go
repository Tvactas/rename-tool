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
	image.SetMinSize(fyne.NewSize(200, 380))
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

	// 左图右按钮布局 - 图片距离左边20px，按钮靠右
	mainContent := container.NewHBox(
		image,              // 左侧图片
		layout.NewSpacer(), // 中间弹性空间，将按钮推向右侧
		buttonScroll,       // 右侧按钮
	)

	// 添加20px左边距和上下边距
	paddedContent := container.NewPadded(mainContent)

	return theme.SetBackground(paddedContent)
}

// buildFixedButtonGrid creates a two-column button layout with fixed height
func buildFixedButtonGrid(buttons []struct {
	text   string
	action func()
}) fyne.CanvasObject {
	// 定义统一的按钮尺寸 - 降低高度
	const buttonWidth = 180
	const buttonHeight = 35 // 从30改为35，更合适的高度

	// 创建按钮列表
	var buttonWidgets []fyne.CanvasObject

	for _, btn := range buttons {
		button := widget.NewButton(btn.text, btn.action)
		button.Importance = widget.LowImportance

		// 创建固定尺寸的容器
		fixedButton := container.NewStack(button)
		fixedButton.Resize(fyne.NewSize(buttonWidth, buttonHeight))

		buttonWidgets = append(buttonWidgets, fixedButton)
	}

	// 如果按钮数量是奇数，添加一个空白占位
	if len(buttonWidgets)%2 != 0 {
		spacer := layout.NewSpacer()
		buttonWidgets = append(buttonWidgets, spacer)
	}

	// 使用 GridWithColumns 创建两列布局
	grid := container.NewGridWithColumns(2, buttonWidgets...)

	return grid
}
