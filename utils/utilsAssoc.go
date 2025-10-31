package utils

import (
	"rename-tool/common/dialogcustomize"
	"rename-tool/setting/i18n"

	"fyne.io/fyne/v2"
)

// tr 函数用于国际化
func tr(key string) string {
	return i18n.Tr(key)
}

func dialogTr(key string) string {
	return i18n.DialogTr(key)
}

func buttonTr(key string) string {
	return i18n.ButtonTr(key)
}

func warningDiaLog(window fyne.Window, message string) {
	dialogcustomize.ShowMessageDialog(
		"warning",
		dialogTr("warning"),
		message,
		window,
	)
}

// Show an error dialog bound to the specified window, to avoid always falling to global.MainWindow
func errorDiaLog(window fyne.Window, message string) {
	dialogcustomize.ShowMessageDialog(
		"error",
		dialogTr("error"),
		message,
		window,
	)
}

// Show an error dialog bound to the specified window, to avoid always falling to global.MainWindow
func successDiaLog(window fyne.Window, message string) {
	dialogcustomize.ShowMessageDialog(
		"success",
		dialogTr("success"),
		message,
		window,
	)
}

func warningMultiDiaLog(window fyne.Window, paths []string) {
	dialogcustomize.ShowMultiLineCopyDialog(
		"warning",
		dialogTr("warning"),
		paths,
		window,
	)
}
