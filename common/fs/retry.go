package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"time"

	"rename-tool/setting/config"
	"rename-tool/setting/global"
	"rename-tool/setting/i18n"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// RenameFile attempts to rename a file with retry and name conflict resolution.
func RenameFile(oldPath, newPath string) error {
	if oldPath == newPath {
		return nil
	}

	baseNewPath := newPath
	counter := 1
	for {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			break
		}
		ext := filepath.Ext(baseNewPath)
		name := baseNewPath[:len(baseNewPath)-len(ext)]
		newPath = fmt.Sprintf("%s_%d%s", name, counter, ext)
		counter++
	}

	var err error
	for i := 0; i < config.MaxRetryAttempts; i++ {
		srcFile, err := os.Open(oldPath)
		if err != nil {
			if IsFileBusyError(err) {
				time.Sleep(config.RetryDelay)
				continue
			}
			return &AppError{
				Code:    "RENAME_OPEN_ERROR",
				Message: i18n.Tr("rename_open_failed") + ": " + oldPath,
				Err:     err,
			}
		}
		srcFile.Close()

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
		Code:    "RENAME_FAILED",
		Message: fmt.Sprintf(i18n.Tr("rename_failed_format"), oldPath, newPath),
		Err:     err,
	}
}

// RetryRenameForFile tries to rename a file in-place up to 3 times.
func RetryRenameForFile(filePath string) bool {
	for i := 0; i < 3; i++ {
		err := RenameFile(filePath, filePath)
		if err == nil {
			return true
		}
		if IsFileBusyError(err) {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}
	return false
}

// RetryRename attempts to retry renaming all given files.
func RetryRename(files []string, window fyne.Window) {
	successCount := 0
	var failedFiles []string

	for _, file := range files {
		if RetryRenameForFile(file) {
			successCount++
		} else {
			failedFiles = append(failedFiles, file)
		}
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf(i18n.Tr("success_retried")+" %d/%d", successCount, len(files)))

	if len(failedFiles) > 0 {
		message.WriteString("\n\n" + i18n.Tr("some_files_may_still_be_in_use") + ":\n")
		for _, file := range failedFiles {
			message.WriteString("  - " + file + "\n")
		}
		ShowBusyFilesDialog(window, failedFiles)
	} else {
		dialog.ShowInformation(i18n.Tr("retry_result"), message.String(), window)
	}
}

// ShowBusyFilesDialog displays a retry dialog for files still in use.
func ShowBusyFilesDialog(window fyne.Window, busyFiles []string) {
	content := strings.Join(busyFiles, "\n")
	textArea := widget.NewMultiLineEntry()
	textArea.SetText(content)
	textArea.Wrapping = fyne.TextWrapWord
	textArea.Disable()

	copyBtn := widget.NewButton(i18n.Tr("copy"), func() {
		global.MyApp.Clipboard().SetContent(content)
		dialog.ShowInformation(i18n.Tr("success"), i18n.Tr("copy_success"), window)
	})

	retryBtn := widget.NewButton(i18n.Tr("retry"), nil)
	cancelBtn := widget.NewButton(i18n.Tr("cancel"), nil)

	var lastRetryTime time.Time
	var retryMutex sync.Mutex

	bottomButtons := container.NewHBox(
		copyBtn,
		layout.NewSpacer(),
		retryBtn,
		cancelBtn,
	)

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
