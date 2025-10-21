package dirbrowser

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// DirBrowser 目录浏览器结构体
type DirBrowser struct {
	parent      fyne.Window
	window      fyne.Window
	currentDir  string
	selectedDir string
	entries     []os.DirEntry
	pathLabel   *widget.Label
	list        *widget.List
	onSelected  func(string)
}

// NewDirBrowser 创建新的目录浏览器实例
func NewDirBrowser(parent fyne.Window, onSelected func(string)) *DirBrowser {
	startDir, err := os.UserHomeDir()
	if err != nil {
		startDir, _ = os.Getwd()
	}

	db := &DirBrowser{
		parent:     parent,
		currentDir: startDir,
		onSelected: onSelected,
		pathLabel:  widget.NewLabel(startDir),
	}

	db.createList()
	return db
}

// createList 创建文件列表组件
func (db *DirBrowser) createList() {
	db.list = widget.NewList(
		func() int { return len(db.entries) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(nil),
				widget.NewLabel(""),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(db.entries) {
				return
			}
			entry := db.entries[i]
			box := o.(*fyne.Container)
			icon := box.Objects[0].(*widget.Icon)
			label := box.Objects[1].(*widget.Label)

			if entry.IsDir() {
				icon.SetResource(theme.FolderIcon())
			} else {
				icon.SetResource(theme.FileIcon())
			}
			label.SetText(entry.Name())
		},
	)

	db.list.OnSelected = func(id widget.ListItemID) {
		if id >= len(db.entries) {
			return
		}
		entry := db.entries[id]
		fullPath := filepath.Join(db.currentDir, entry.Name())

		if entry.IsDir() {
			db.loadDir(fullPath)
		} else {
			// 点击文件时，选中其所在目录
			db.selectedDir = db.currentDir
		}
	}
}

// loadDir 加载指定目录
func (db *DirBrowser) loadDir(path string) {
	// 检查路径是否可访问
	info, err := os.Stat(path)
	if err != nil {
		dialog.ShowError(fmt.Errorf("无法访问目录: %v", err), db.parent)
		return
	}

	if !info.IsDir() {
		dialog.ShowError(fmt.Errorf("路径不是目录: %s", path), db.parent)
		return
	}

	files, err := os.ReadDir(path)
	if err != nil {
		dialog.ShowError(fmt.Errorf("读取目录失败: %v", err), db.parent)
		return
	}

	// 排序：目录在前，然后按名称排序
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir() != files[j].IsDir() {
			return files[i].IsDir()
		}
		return files[i].Name() < files[j].Name()
	})

	db.entries = files
	db.currentDir = path
	db.selectedDir = path // 更新选中目录为当前目录
	db.pathLabel.SetText(path)
	db.list.Refresh()
	db.list.UnselectAll()
}

// goToParentDir 返回上级目录
func (db *DirBrowser) goToParentDir() {
	parentDir := filepath.Dir(db.currentDir)
	if parentDir != db.currentDir {
		db.loadDir(parentDir)
	}
}

// confirm 确认选择
func (db *DirBrowser) confirm() {
	// 如果没有明确选择，使用当前目录
	if db.selectedDir == "" {
		db.selectedDir = db.currentDir
	}

	if db.onSelected != nil {
		db.onSelected(db.selectedDir)
	}

	if db.window != nil {
		db.window.Close()
	}
}

// cancel 取消选择
func (db *DirBrowser) cancel() {
	if db.window != nil {
		db.window.Close()
	}
}

// Show 显示目录浏览器窗口
func (db *DirBrowser) Show() {
	// 创建工具栏
	backBtn := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), db.goToParentDir)
	backBtn.Importance = widget.LowImportance

	homeBtn := widget.NewButtonWithIcon("", theme.HomeIcon(), func() {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			db.loadDir(homeDir)
		}
	})
	homeBtn.Importance = widget.LowImportance

	// 创建按钮
	okBtn := widget.NewButton("确定", db.confirm)
	okBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton("取消", db.cancel)

	// 布局
	toolbar := container.NewBorder(
		nil, nil,
		container.NewHBox(backBtn, homeBtn),
		nil,
		db.pathLabel,
	)

	footer := container.NewBorder(
		nil, nil, nil,
		container.NewHBox(cancelBtn, okBtn),
	)

	content := container.NewBorder(
		toolbar,
		footer,
		nil, nil,
		db.list,
	)

	// 创建窗口
	db.window = fyne.CurrentApp().NewWindow("选择目录")
	db.window.Resize(fyne.NewSize(650, 450))
	db.window.SetContent(content)
	db.window.CenterOnScreen()

	// 加载初始目录
	db.loadDir(db.currentDir)

	db.window.Show()
}

// ShowDirBrowser 快捷函数：显示目录浏览器
func ShowDirBrowser(parent fyne.Window, onSelected func(string)) {
	browser := NewDirBrowser(parent, onSelected)
	browser.Show()
}

// ------------------- Demo 主函数 -------------------

func main() {
	a := app.New()
	win := a.NewWindow("Fyne 目录浏览器 Demo")

	selectedLabel := widget.NewLabel("未选择目录")
	selectedLabel.Wrapping = fyne.TextWrapWord

	btn := widget.NewButton("选择目录", func() {
		ShowDirBrowser(win, func(selected string) {
			selectedLabel.SetText("已选择目录: " + selected)
			fmt.Println("用户选中目录:", selected)
		})
	})
	btn.Importance = widget.HighImportance

	content := container.NewVBox(
		widget.NewLabel("点击按钮弹出目录浏览器"),
		widget.NewSeparator(),
		btn,
		widget.NewSeparator(),
		selectedLabel,
	)

	win.SetContent(container.NewPadded(content))
	win.Resize(fyne.NewSize(500, 300))
	win.CenterOnScreen()
	win.ShowAndRun()
}
