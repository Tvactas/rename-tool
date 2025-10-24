package menu

import (
	"rename-tool/common/admin"
	"rename-tool/common/applog"
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

// ShowMainMenu displays the main menu interface
func ShowMainMenu() {
	global.MyApp.Settings().SetTheme(&theme.MainTheme{})

	// Use embedded image resource
	imgResource := theme.LoadImage("cat.png")
	var image *canvas.Image
	if imgResource == nil {
		image = canvas.NewImageFromFile("")
		applog.Logger.Printf("[THEME ERROR]  %s", i18n.LogTr("LoadThemeError"))
	} else {
		image = canvas.NewImageFromResource(imgResource)
	}
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(250, 380))

	// Show permission status
	adminStatus := widget.NewLabel("")
	UpdateAdminStatusLabel(adminStatus)

	// Optimize button creation
	makeTextBtn := func(text string, onTap func()) fyne.CanvasObject {
		btn := widget.NewButton(text, onTap)
		btn.Importance = widget.LowImportance
		return container.NewHBox(btn, layout.NewSpacer())
	}

	// Use predefined button list
	buttons := []struct {
		text   string
		action func()
	}{
		{i18n.ButtonTr("SequenceRename"), func() { utils.ShowBatchRenameNormal() }},
		{tr("ext"), func() { utils.ShowChangeExtension() }},
		{tr("upper"), func() { utils.ShowRenameToCase("upper") }},
		{tr("lower"), func() { utils.ShowRenameToCase("lower") }},
		{tr("titlecase"), func() { utils.ShowRenameToCase("title") }},
		{tr("camel"), func() { utils.ShowRenameToCase("camel") }},
		{tr("insert_char"), func() { utils.ShowInsertCharRename() }},
		{tr("delete_char"), func() { utils.ShowDeleteCharRename() }},
		{tr("regex_replace"), func() { utils.ShowRegexReplace() }},
		{tr("undo"), utils.UndoRename},
		{tr("log"), utils.SaveLogs},
		{tr("exit"), func() { global.MyApp.Quit() }},
	}

	// Create button grid
	var buttonGridItems []fyne.CanvasObject
	for _, btn := range buttons {
		buttonGridItems = append(buttonGridItems, makeTextBtn(btn.text, btn.action))
	}
	buttonGrid := container.NewGridWithColumns(2, buttonGridItems...)

	// Optimize layout
	rightBox := container.NewVBox(buttonGrid)
	mainContent := container.NewBorder(nil, nil, image, rightBox)
	centered := container.NewVBox(
		layout.NewSpacer(),
		mainContent,
		layout.NewSpacer(),
	)

	bgContent := theme.SetBackground(centered)
	langSelector := i18n.LangSelect()

	header := container.NewHBox(
		langSelector,
		layout.NewSpacer(),
		adminStatus,
		layout.NewSpacer(),
		widget.NewLabel(tr("title")),
	)

	content := container.NewVBox(
		header,
		bgContent,
	)

	global.MainWindow.SetContent(content)
	global.MainWindow.Show()
}

// Helper function for translation
func tr(key string) string {
	return i18n.Tr(key)
}

func UpdateAdminStatusLabel(label *widget.Label) {
	if admin.IsAdmin() {
		label.SetText(i18n.ButtonTr("UserPermissionsAD"))
	} else {
		label.SetText(i18n.ButtonTr("UserPermissionsUser"))
	}
}
