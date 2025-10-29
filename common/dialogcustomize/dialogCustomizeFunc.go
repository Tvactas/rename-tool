package dialogcustomize

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
)

// =====================================
// 淡色定义（RGBA）
// =====================================
var (
	ColorSuccess = color.RGBA{R: 220, G: 255, B: 220, A: 255} // 淡绿
	ColorWarning = color.RGBA{R: 255, G: 250, B: 210, A: 255} // 淡黄
	ColorError   = color.RGBA{R: 255, G: 225, B: 225, A: 255} // 淡红
	ColorDefault = color.RGBA{R: 255, G: 255, B: 255, A: 255} // 白
)

// =====================================
// 基础工具函数
// =====================================

// 根据类型返回背景色
func getBgColor(kind string) color.Color {
	switch kind {
	case "success":
		return ColorSuccess
	case "warning":
		return ColorWarning
	case "error":
		return ColorError
	default:
		return ColorDefault
	}
}

// 居中文本标签
func createCenteredLabel(msg string) fyne.CanvasObject {
	lbl := canvas.NewText(msg, color.Black)
	lbl.Alignment = fyne.TextAlignCenter
	return lbl
}

// 创建容器（带背景与内边距）
func createContentContainer(content fyne.CanvasObject, width float32, bg color.Color) *fyne.Container {
	// 背景矩形：给定宽度和一定高度，撑开视觉背景
	bgRect := canvas.NewRectangle(bg)
	bgRect.SetMinSize(fyne.NewSize(width, 120)) // ← 高度别太小，否则还是只包文字

	// 内容区域（加内边距）
	padded := container.NewPadded(content)

	// 居中内容
	centered := container.NewVBox(
		layout.NewSpacer(),
		padded,
		layout.NewSpacer(),
	)

	// 用 Stack 将背景放底层
	stack := container.NewStack(bgRect, centered)

	// 给整个对话框再包一层 Padding，避免贴边
	return container.NewPadded(stack)
}

// =====================================
// 核心：统一的底层显示函数
// =====================================
func showBaseDialog(title, message string, window fyne.Window, bg color.Color, width float32) *dialog.CustomDialog {
	content := createCenteredLabel(message)
	contentContainer := createContentContainer(content, width, bg)
	return dialog.NewCustom(title, dialogTr("confirm"), contentContainer, window)
}

// =====================================
// 对外统一接口
// =====================================
