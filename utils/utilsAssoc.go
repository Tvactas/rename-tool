package utils

import (
	"rename-tool/common/dialogcustomize"
	"rename-tool/setting/global"
	"rename-tool/setting/i18n"
)

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
