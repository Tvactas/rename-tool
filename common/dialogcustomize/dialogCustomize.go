package dialogcustomize

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
)

// 普通消息框
func ShowMessageDialog(kind, title, message string, window fyne.Window) {
	bg := getBgColor(kind)
	d := showBaseDialog(title, message, window, bg, 400)
	d.Show()
}


// 多行+复制按钮专用弹窗
func ShowMultiLineCopyDialog(kind, title string, paths []string, window fyne.Window) {
	bg := getBgColor(kind)
	content := strings.Join(paths, "\n")

	textArea := widget.NewMultiLineEntry()
	textArea.SetText(content)
	textArea.Wrapping = fyne.TextWrapWord
	textArea.Disable()
	textArea.SetMinRowsVisible(6)

	contentContainer := createContentContainer(textArea, 400, bg)

	copyBtn := widget.NewButton(dialogTr("copy"), func() {
		window.Clipboard().SetContent(content)
	})
	closeBtn := widget.NewButton(dialogTr("confirm"), nil)
	btns := container.NewHBox(layout.NewSpacer(), copyBtn, closeBtn, layout.NewSpacer())

	finalContent := container.NewVBox(contentContainer, btns)
	dialogErr := dialog.NewCustomWithoutButtons(title, finalContent, window)
	closeBtn.OnTapped = dialogErr.Hide
	dialogErr.Show()
}
