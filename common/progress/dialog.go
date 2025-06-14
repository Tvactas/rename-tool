package progress

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Dialog 进度对话框
type Dialog struct {
	window    fyne.Window
	dialog    *dialog.CustomDialog
	progress  *widget.ProgressBar
	status    *widget.Label
	cancelBtn *widget.Button
	cancelled bool
	mu        sync.RWMutex
	updateCh  chan struct {
		progress float64
		status   string
	}
}

// NewDialog 创建新的进度对话框
func NewDialog(title string, window fyne.Window) *Dialog {
	progress := widget.NewProgressBar()
	progress.Resize(fyne.NewSize(400, progress.MinSize().Height))

	status := widget.NewLabel("")
	status.Wrapping = fyne.TextWrapWord

	cancelBtn := widget.NewButton("取消", nil)

	content := container.NewVBox(
		container.NewCenter(progress),
		container.NewCenter(status),
		container.NewCenter(cancelBtn),
	)

	dialog := dialog.NewCustom(
		title,
		"",
		content,
		window,
	)

	pd := &Dialog{
		window:    window,
		dialog:    dialog,
		progress:  progress,
		status:    status,
		cancelBtn: cancelBtn,
		cancelled: false,
		updateCh: make(chan struct {
			progress float64
			status   string
		}, 100),
	}

	cancelBtn.OnTapped = func() {
		pd.mu.Lock()
		pd.cancelled = true
		pd.mu.Unlock()
		pd.dialog.Hide()
	}

	// 启动更新处理协程
	go func() {
		for update := range pd.updateCh {
			progress.SetValue(update.progress)
			status.SetText(update.status)
		}
	}()

	return pd
}

// Show 显示对话框
func (pd *Dialog) Show() {
	pd.dialog.Show()
}

// Hide 隐藏对话框
func (pd *Dialog) Hide() {
	pd.dialog.Hide()
	close(pd.updateCh) // 关闭更新通道
}

// Update 更新进度和状态
func (pd *Dialog) Update(progress float64, status string) {
	select {
	case pd.updateCh <- struct {
		progress float64
		status   string
	}{progress, status}:
	default:
		// 如果通道已满，丢弃更新
	}
}

// IsCancelled 检查是否已取消
func (pd *Dialog) IsCancelled() bool {
	pd.mu.RLock()
	defer pd.mu.RUnlock()
	return pd.cancelled
}
