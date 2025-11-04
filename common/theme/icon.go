package theme

import (
	"fyne.io/fyne/v2"
)

// Icon 封装应用图标
type Icon struct {
	resource fyne.Resource
}

// NewIcon 从内嵌资源初始化图标
func NewIcon(name string) *Icon {
	return &Icon{
		resource: LoadImage(name), // 复用 theme 包的 LoadImage
	}
}

// Resource 返回 fyne.Resource 对象
func (i *Icon) Resource() fyne.Resource {
	return i.resource
}

// Apply 设置到 App 和 Window
func (i *Icon) Apply(app fyne.App, win fyne.Window) {
	if i.resource != nil {
		app.SetIcon(i.resource)
		if win != nil {
			win.SetIcon(i.resource)
		}
	}
}
