package utils

import (
	"rename-tool/common/dialogcustomize"
	"rename-tool/setting/global"
	"rename-tool/setting/i18n"
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

func errorDiaLog(message string) {
	dialogcustomize.ShowMessageDialog(
		"error",
		dialogTr("error"),
		message,
		global.MainWindow, // 直接用全局变量
	)
}

func warningDiaLog(message string) {
	dialogcustomize.ShowMessageDialog(
		"warning",
		dialogTr("warning"),
		message,
		global.MainWindow, // 直接用全局变量
	)
}
