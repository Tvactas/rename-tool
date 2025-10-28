package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"rename-tool/setting/i18n"
)

// FilenameLengthError 表示文件名长度错误
type FilenameLengthError struct {
	Files []string
}

func (e *FilenameLengthError) Error() string {
	return fmt.Sprintf("以下文件名的长度小于指定的插入位置：\n%s", strings.Join(e.Files, "\n"))
}

// ShowLengthErrorDialog 显示文件名长度错误对话框
func ShowLengthErrorDialog(window fyne.Window, files []string) {
	// 创建文本内容
	content := strings.Join(files, "\n")
	textArea := widget.NewMultiLineEntry()
	textArea.SetText(content)
	textArea.Wrapping = fyne.TextWrapWord
	textArea.Disable() // 设置为只读

	// 创建按钮
	copyBtn := widget.NewButton(i18n.Tr("copy"), func() {
		window.Clipboard().SetContent(content)
		dialog.ShowInformation(dialogTr("success"), i18n.Tr("copy_success"), window)
	})

	closeBtn := widget.NewButton(i18n.Tr("close"), nil)

	// 创建对话框内容
	dialogContent := container.NewBorder(
		widget.NewLabel(i18n.Tr("filename_length_error")+":"),
		container.NewHBox(copyBtn, layout.NewSpacer(), closeBtn),
		nil,
		nil,
		container.NewStack(textArea),
	)

	dialog := dialog.NewCustom(
		dialogTr("error"),
		"",
		dialogContent,
		window,
	)

	// 设置关闭按钮动作
	closeBtn.OnTapped = dialog.Hide

	dialog.Show()
}

func ShowWidePlainMessage(win fyne.Window, title, message string) {
	label := widget.NewLabel(message)
	label.Wrapping = fyne.TextWrapWord // 自动换行

	// 给内容加一个最小宽度容器，让对话框宽一些
	content := container.NewVBox(
		label,
	)
	scroll := container.NewScroll(content)    // 可滚动，防止内容过多
	scroll.SetMinSize(fyne.NewSize(200, 100)) // 设置对话框最小尺寸

	dialog.ShowCustom(
		title,
		"OK",
		scroll,
		win,
	)
}
