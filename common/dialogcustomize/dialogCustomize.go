package dialogcustomize

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

// ShowCustomDialog 显示自定义对话框（无图标）
func ShowCustomDialog(title, message string, window fyne.Window) {
	content := createCenteredLabel(message)
	contentContainer := createContentContainer(content, 400)

	customDialog := dialog.NewCustom(title, dialogTr("confirm"), contentContainer, window)
	customDialog.Show()
}

// ShowCustomDialogWithCallback 显示自定义对话框并带回调
func ShowCustomDialogWithCallback(title, message string, window fyne.Window, callback func()) {
	content := createCenteredLabel(message)
	contentContainer := createContentContainer(content, 400)

	customDialog := dialog.NewCustom(title, dialogTr("confirm"), contentContainer, window)
	customDialog.SetOnClosed(func() {
		if callback != nil {
			callback()
		}
	})
	customDialog.Show()
}

// ShowCustomConfirm 显示确认对话框（无图标）
func ShowCustomConfirm(title, message string, window fyne.Window, onConfirm func(), onCancel func()) {
	content := createCenteredLabel(message)
	contentContainer := createContentContainer(content, 400)

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

// ShowCustomDialogMultiLine 显示多行内容的对话框
func ShowCustomDialogMultiLine(title string, messages []string, window fyne.Window) {
	labels := make([]fyne.CanvasObject, len(messages))
	for i, msg := range messages {
		labels[i] = createCenteredLabel(msg)
	}

	content := append([]fyne.CanvasObject{createSpacer(400)}, labels...)
	contentContainer := container.NewVBox(content...)

	customDialog := dialog.NewCustom(title, dialogTr("confirm"), contentContainer, window)
	customDialog.Show()
}

// ShowCustomDialogWithSize 显示自定义尺寸的对话框
func ShowCustomDialogWithSize(title, message string, window fyne.Window, width, height float32) {
	content := createCenteredLabel(message)
	contentContainer := createContentContainer(content, width)

	customDialog := dialog.NewCustom(title, dialogTr("confirm"), contentContainer, window)
	customDialog.Show()
}
