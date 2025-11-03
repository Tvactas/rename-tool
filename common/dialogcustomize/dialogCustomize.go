package dialogcustomize

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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
	textArea.SetMinRowsVisible(6)

    // 保持可见文本为常规深色，同时禁止用户编辑
    original := content
    isUpdating := false
    textArea.OnChanged = func(s string) {
        if isUpdating {
            return
        }
        if s != original {
            isUpdating = true
            textArea.SetText(original)
            isUpdating = false
        }
    }

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

// ShowMultiLineErrorDialog 显示带错误信息的多行弹窗（文件路径 + 错误信息）
func ShowMultiLineErrorDialog(kind, title string, errors map[string]error, window fyne.Window) {
	bg := getBgColor(kind)
	
	var lines []string
	for file, err := range errors {
		if err != nil {
			lines = append(lines, fmt.Sprintf("%s\n  └─ %s", file, err.Error()))
		} else {
			lines = append(lines, file)
		}
	}
	content := strings.Join(lines, "\n\n")

	textArea := widget.NewMultiLineEntry()
	textArea.SetText(content)
	textArea.Wrapping = fyne.TextWrapWord
	textArea.SetMinRowsVisible(6)

    // 保持可见文本为常规深色，同时禁止用户编辑
    original := content
    isUpdating := false
    textArea.OnChanged = func(s string) {
        if isUpdating {
            return
        }
        if s != original {
            isUpdating = true
            textArea.SetText(original)
            isUpdating = false
        }
    }

	contentContainer := createContentContainer(textArea, 500, bg)

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