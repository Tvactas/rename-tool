package filestatus

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"rename-tool/setting/global"
	"rename-tool/setting/i18n"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

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
		dialog.ShowInformation(dialogTr("success"), i18n.Tr("copy_success"), window)
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
			dialog.ShowInformation(dialogTr("warning"), i18n.Tr("retry_too_fast"), window)
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
