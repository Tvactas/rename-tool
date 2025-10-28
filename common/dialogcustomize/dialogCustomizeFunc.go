package dialogcustomize

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// createCenteredLabel 创建居中对齐的 Label
func createCenteredLabel(message string) *widget.Label {
	content := widget.NewLabel(message)
	content.Wrapping = fyne.TextWrapWord
	content.Alignment = fyne.TextAlignCenter
	return content
}

// createSpacer 创建透明占位符
func createSpacer(width float32) *canvas.Rectangle {
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(width, 1))
	return spacer
}

// createContentContainer 创建带占位符的内容容器
func createContentContainer(content fyne.CanvasObject, width float32) *fyne.Container {
	return container.NewVBox(
		createSpacer(width),
		content,
	)
}
