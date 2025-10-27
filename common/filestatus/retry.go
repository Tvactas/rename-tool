package filestatus

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
	ext := filepath.Ext(baseNewPath)
	name := baseNewPath[:len(baseNewPath)-len(ext)]
	for {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			break
		}
		newPath = fmt.Sprintf("%s_%d%s", name, counter, ext)
		counter++
	}

	var err error
	delay := config.RetryDelay
	for i := 0; i < config.MaxRetryAttempts; i++ {
		err = os.Rename(oldPath, newPath)
		if err == nil {
			return nil
		}
		if !IsFileBusyError(err) {
			break
		}
		time.Sleep(delay)
		delay *= 2
	}

	return fmt.Errorf("%s: %s → %s", i18n.Tr("rename_failed_format"), oldPath, newPath)

}

// RetryRenameForFile attempts to rename the file to a temp path and revert to check file busy.
func RetryRenameForFile(filePath string) bool {
	tempPath := filePath + ".tmp_retry"
	err := RenameFile(filePath, tempPath)
	if err != nil {
		return false
	}
	// Try revert back
	err = RenameFile(tempPath, filePath)
	return err == nil
}

// RetryRename attempts to retry renaming all given files concurrently.
func RetryRename(files []string, window fyne.Window) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	var failedFiles []string
	sem := make(chan struct{}, 5) // max concurrency = 5

	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if RetryRenameForFile(f) {
				mu.Lock()
				successCount++
				mu.Unlock()
			} else {
				mu.Lock()
				failedFiles = append(failedFiles, f)
				mu.Unlock()
			}
		}(file)
	}

	wg.Wait()

	message := fmt.Sprintf(i18n.Tr("success_retried")+" %d/%d", successCount, len(files))
	if len(failedFiles) > 0 {
		message += "\n\n" + i18n.Tr("some_files_may_still_be_in_use") + ":\n  - " + strings.Join(failedFiles, "\n  - ")
		fyne.CurrentApp().SendNotification(&fyne.Notification{Title: i18n.Tr("retry_result"), Content: message})
		ShowBusyFilesDialog(window, failedFiles)
	} else {
		dialog.ShowInformation(i18n.Tr("retry_result"), message, window)
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
		go RetryRename(busyFiles, window)
	}

	cancelBtn.OnTapped = busyFilesDialog.Hide
	busyFilesDialog.Show()
}
