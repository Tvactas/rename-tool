package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"rename-tool/setting/config"
	"rename-tool/setting/i18n"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// RenameFile 重命名文件
func RenameFile(oldPath, newPath string) error {
	if oldPath == newPath {
		return nil
	}

	// 避免文件名冲突
	counter := 1
	baseNewPath := newPath
	for {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			break
		}
		ext := filepath.Ext(baseNewPath)
		nameWithoutExt := baseNewPath[:len(baseNewPath)-len(ext)]
		newPath = fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext)
		counter++
	}

	// 使用重试机制
	var err error
	for i := 0; i < config.MaxRetryAttempts; i++ {
		// 尝试打开源文件以确保可访问
		srcFile, err := os.Open(oldPath)
		if err != nil {
			if IsFileBusyError(err) {
				time.Sleep(config.RetryDelay)
				continue
			}
			return &AppError{
				Code:    "RENAME_ERROR",
				Message: fmt.Sprintf("Failed to open source file: %s", oldPath),
				Err:     err,
			}
		}
		srcFile.Close()

		// 执行重命名
		err = os.Rename(oldPath, newPath)
		if err == nil {
			return nil
		}

		if !IsFileBusyError(err) {
			break
		}

		time.Sleep(config.RetryDelay)
	}

	return &AppError{
		Code:    "RENAME_ERROR",
		Message: fmt.Sprintf("Failed to rename file: %s -> %s", oldPath, newPath),
		Err:     err,
	}
}

// RetryRenameForFile 重试重命名单个文件
func RetryRenameForFile(filePath string) bool {
	// 尝试3次重命名
	for i := 0; i < 3; i++ {
		err := RenameFile(filePath, filePath)
		if err == nil {
			return true
		}

		// 如果文件被占用，等待一段时间后重试
		if IsFileBusyError(err) {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}
	return false
}

// RetryRename 重试重命名所有被占用的文件
func RetryRename(files []string, window fyne.Window) {
	successCount := 0
	failedFiles := []string{}

	for _, file := range files {
		if RetryRenameForFile(file) {
			successCount++
		} else {
			failedFiles = append(failedFiles, file)
		}
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf(i18n.Tr("success_retried")+" %d/%d 个文件", successCount, len(files)))

	if len(failedFiles) > 0 {
		message.WriteString("\n\n" + i18n.Tr("some_files_may_still_be_in_use") + ":\n")
		for _, file := range failedFiles {
			message.WriteString(fmt.Sprintf("  - %s\n", file))
		}
		// 如果还有失败的文件，显示重试对话框
		ShowBusyFilesDialog(window, failedFiles)
	} else {
		dialog.ShowInformation(i18n.Tr("retry_result"), message.String(), window)
	}
}

// ShowBusyFilesDialog 显示被占用文件的对话框
func ShowBusyFilesDialog(window fyne.Window, busyFiles []string) {
	// 创建文本内容
	content := strings.Join(busyFiles, "\n")
	textArea := widget.NewMultiLineEntry()
	textArea.SetText(content)
	textArea.Wrapping = fyne.TextWrapWord
	textArea.Disable() // 设置为只读

	// 创建按钮
	copyBtn := widget.NewButton(i18n.Tr("copy"), func() {
		window.Clipboard().SetContent(content)
		dialog.ShowInformation(i18n.Tr("success"), i18n.Tr("copy_success"), window)
	})

	retryBtn := widget.NewButton(i18n.Tr("retry"), nil)
	cancelBtn := widget.NewButton(i18n.Tr("cancel"), nil)

	// 添加重试按钮点击限制
	var lastRetryTime time.Time
	var retryMutex sync.Mutex

	// 创建底部按钮容器
	bottomButtons := container.NewHBox(
		copyBtn,
		layout.NewSpacer(),
		retryBtn,
		cancelBtn,
	)

	// 创建对话框内容
	dialogContent := container.NewBorder(
		widget.NewLabel(i18n.Tr("busy_files_message")+":"),
		bottomButtons,
		nil,
		nil,
		container.NewStack(textArea),
	)

	busyFilesDialog := dialog.NewCustom(
		i18n.Tr("busy_files_title"),
		"",
		dialogContent,
		window,
	)

	// 设置按钮动作
	retryBtn.OnTapped = func() {
		retryMutex.Lock()
		if time.Since(lastRetryTime) < 2*time.Second {
			dialog.ShowInformation(i18n.Tr("warning"), i18n.Tr("retry_too_fast"), window)
			retryMutex.Unlock()
			return
		}
		lastRetryTime = time.Now()
		retryMutex.Unlock()

		busyFilesDialog.Hide()
		RetryRename(busyFiles, window)
	}
	cancelBtn.OnTapped = busyFilesDialog.Hide

	busyFilesDialog.Show()
}
