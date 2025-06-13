package ui

import (
	"fmt"

	"rename-tool/setting/i18n"
	"rename-tool/setting/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// ShowSuccessMessage 显示重命名成功的消息
func ShowSuccessMessage(window fyne.Window, renameType model.RenameType, count int) {
	var message string
	switch renameType {
	case model.RenameTypeBatch:
		message = fmt.Sprintf(tr("success_renamed")+" "+tr("files_count"), count)
	case model.RenameTypeExtension:
		message = fmt.Sprintf(tr("success_modified")+" "+tr("files_count"), count)
	case model.RenameTypeCase:
		message = fmt.Sprintf(tr("success_renamed")+" "+tr("files_count"), count)
	case model.RenameTypeInsertChar:
		message = fmt.Sprintf(tr("success_inserted")+" "+tr("files_count"), count)
	case model.RenameTypeReplace:
		message = fmt.Sprintf(tr("success_replaced")+" "+tr("files_count"), count)
	case model.RenameTypeDeleteChar:
		message = fmt.Sprintf(tr("success_deleted")+" "+tr("files_count"), count)
	}

	dialog.ShowInformation(tr("success"), message, window)
}

// tr 是 i18n.Tr 的包装函数
func tr(key string) string {
	return i18n.Tr(key)
}
