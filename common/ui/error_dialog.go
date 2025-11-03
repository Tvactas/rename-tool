package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// FilenameLengthError 表示文件名长度错误
type FilenameLengthError struct {
	Files []string
}

func (e *FilenameLengthError) Error() string {
	return fmt.Sprintf("以下文件名的长度小于指定的插入位置：\n%s", strings.Join(e.Files, "\n"))
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
