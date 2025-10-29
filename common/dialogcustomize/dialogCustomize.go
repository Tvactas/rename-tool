package dialogcustomize

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

// 普通消息框
func ShowMessageDialog(kind, title, message string, window fyne.Window) {
	bg := getBgColor(kind)
	d := showBaseDialog(title, message, window, bg, 400)
	d.Show()
}

// 带回调消息框
func ShowMessageDialogWithCallback(kind, title, message string, window fyne.Window, callback func()) {
	bg := getBgColor(kind)
	d := showBaseDialog(title, message, window, bg, 400)
	d.SetOnClosed(callback)
	d.Show()
}

// 确认框（带确认/取消）
func ShowMessageConfirm(kind, title, message string, window fyne.Window, onConfirm, onCancel func()) {
	bg := getBgColor(kind)
	content := createCenteredLabel(message)
	contentContainer := createContentContainer(content, 400, bg)

	customDialog := dialog.NewCustomConfirm(
		title,
		dialogTr("confirm"),
		"Cancel",
		contentContainer,
		func(confirmed bool) {
			if confirmed && onConfirm != nil {
				onConfirm()
			} else if !confirmed && onCancel != nil {
				onCancel()
			}
		},
		window,
	)
	customDialog.Show()
}

// 多行内容框
func ShowMessageDialogMultiLine(kind, title string, messages []string, window fyne.Window) {
	bg := getBgColor(kind)

	labels := make([]fyne.CanvasObject, len(messages))
	for i, msg := range messages {
		labels[i] = createCenteredLabel(msg)
	}

	content := container.NewVBox(labels...)
	contentContainer := createContentContainer(content, 400, bg)

	dialog.NewCustom(title, dialogTr("confirm"), contentContainer, window).Show()
}

// 可自定义尺寸
func ShowMessageDialogWithSize(kind, title, message string, window fyne.Window, width, height float32) {
	bg := getBgColor(kind)
	content := createCenteredLabel(message)
	contentContainer := createContentContainer(content, width, bg)

	customDialog := dialog.NewCustom(title, dialogTr("confirm"), contentContainer, window)
	customDialog.Resize(fyne.NewSize(width, height))
	customDialog.Show()
}
