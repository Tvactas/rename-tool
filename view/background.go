package view

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// 背景设置函数
func SetBackground(content fyne.CanvasObject) fyne.CanvasObject {
	// 创建蓝到紫的线性渐变（左上到右下）
	grad1 := canvas.NewLinearGradient(
		color.RGBA{R: 0, G: 128, B: 255, A: 255}, // 蓝色
		color.RGBA{R: 128, G: 0, B: 255, A: 255}, // 紫色
		45,                                       // 角度，左上到右下
	)
	// 叠加紫到绿的半透明渐变
	grad2 := canvas.NewLinearGradient(
		color.RGBA{R: 128, G: 0, B: 255, A: 128}, // 半透明紫色
		color.RGBA{R: 0, G: 255, B: 128, A: 128}, // 半透明绿色
		45,
	)

	return container.NewStack(
		grad1,
		grad2,
		container.NewPadded(content),
	)
}
