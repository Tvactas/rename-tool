package menu

import (
	"rename-tool/common/admin"
	"rename-tool/common/theme"
	"rename-tool/setting/global"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ShowMainMenu renders the main UI menu.
func ShowMainMenu() {
	global.MyApp.Settings().SetTheme(&theme.MainTheme{})

	content := container.NewVBox(
		buildHeader(),
		buildBody(),
	)

	global.MainWindow.SetContent(content)
	global.MainWindow.Show()
}

// UpdateAdminStatusLabel updates the user permission label.
func UpdateAdminStatusLabel(label *widget.Label) {
	if admin.IsAdmin() {
		label.SetText(buttonTr("userPermissionsAD"))
	} else {
		label.SetText(buttonTr("userPermissionsUser"))
	}
}
