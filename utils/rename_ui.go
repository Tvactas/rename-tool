package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"rename-tool/common/antisamename"
	"rename-tool/common/dialogcustomize"
	"rename-tool/common/dirpath"
	"rename-tool/common/filestatus"
	"rename-tool/common/pathgen"
	"rename-tool/common/progress"
	"rename-tool/common/ui"
	"rename-tool/setting/global"
	"rename-tool/setting/i18n"
	"rename-tool/setting/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// RenameUIConfig 重命名界面配置
type RenameUIConfig struct {
	Title           string
	Window          fyne.Window
	RenameType      model.RenameType
	ConfigBuilder   func() model.RenameConfig
	ValidateConfig  func(config model.RenameConfig) error
	AdditionalItems []fyne.CanvasObject
}

// ✅ 主入口函数，整合UI、事件与布局
func ShowRenameUI(config RenameUIConfig) {
	ui, err := initRenameUI(config)

	if err != nil {
		errorDiaLog(global.MainWindow, err.Error())
		return
	}

	scanBtn, previewBtn, renameBtn, backBtn := setupRenameUIEvents(ui, config)

	dirBox := container.NewHBox(ui.DirSelector, scanBtn)
	formatBox := container.NewHBox(ui.FormatLabel, ui.SelectAllBtn)

	mainContent := container.NewVBox(
		ui.Title,
		widget.NewSeparator(),
		dirBox,
		widget.NewSeparator(),
		formatBox,
		ui.FormatScroll,
		widget.NewSeparator(),
	)

	if len(config.AdditionalItems) > 0 {
		mainContent.Add(container.NewVBox(config.AdditionalItems...))
		mainContent.Add(widget.NewSeparator())
	}

	bottomButtons := container.NewHBox(layout.NewSpacer(), backBtn, previewBtn, renameBtn)
	mainContent.Add(bottomButtons)

	ui.Window.SetContent(mainContent)
	ui.Window.Show()
}

// performRename 执行重命名操作
func performRename(window fyne.Window, config model.RenameConfig) {
	if config.SelectedDir == "" {
		errorDiaLog(window, dialogTr("selectDirFirst"))

		return
	}
	// 获取文件列表
	files, err := dirpath.GetFiles(config.SelectedDir, config.Formats)
	if err != nil {
		errorDiaLog(window, tr("error_getting_files"))

		return
	}

	// 统一防重名预检（批量内部重复、命中磁盘已存在路径）
	if stop, err := antisamename.CheckAndShowConflicts(window, files, config); err != nil {
		dialog.ShowError(err, window)
		return
	} else if stop {
		return
	}

	// 检查重名
	if config.Type == model.RenameTypeReplace {
		duplicates, err := pathgen.CheckDuplicateNames(files, config)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if len(duplicates) > 0 {
			content := strings.Join(duplicates, "\n")
			textArea := widget.NewMultiLineEntry()
			textArea.SetText(content)
			textArea.Wrapping = fyne.TextWrapWord
			textArea.Disable()

			copyBtn := widget.NewButton(tr("copy"), func() {
				window.Clipboard().SetContent(content)
				dialog.ShowInformation(dialogTr("success"), tr("copySuccess"), window)
			})

			closeBtn := widget.NewButton(tr("close"), nil)

			dialogContent := container.NewBorder(
				widget.NewLabel(dialogTr("error")+": "+tr("duplicate_names")),
				container.NewHBox(copyBtn, layout.NewSpacer(), closeBtn),
				nil,
				nil,
				container.NewStack(textArea),
			)

			dialog := dialog.NewCustom(
				dialogTr("error"),
				"",
				dialogContent,
				window,
			)

			closeBtn.OnTapped = dialog.Hide
			dialog.Show()
			return
		}
	}

	// 创建进度对话框
	pd := progress.NewDialog(buttonTr("implement"), window)
	pd.Show()

	// 使用工作池处理文件
	workerCount := runtime.NumCPU()
	fileChan := make(chan string, len(files))
	resultChan := make(chan struct {
		file string
		err  error
	}, len(files))

	// 创建本地计数器map和互斥锁
	counters := make(map[string]int)
	var countersMutex sync.Mutex

	// 启动工作协程
	var wg sync.WaitGroup
	counter := 0
	var counterMutex sync.Mutex
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				var newPath string
				var err error

				switch config.Type {
				case model.RenameTypeBatch:
					counterMutex.Lock()
					currentCounter := counter
					counter++
					counterMutex.Unlock()

					// 使用互斥锁保护counters map的访问
					countersMutex.Lock()
					// 如果是第一次遇到这个扩展名，重置计数器
					ext := filepath.Ext(file)
					if _, exists := counters[ext]; !exists {
						counters[ext] = 0
					}
					newPath, err = pathgen.GenerateBatchRenamePath(file, config, currentCounter, counters)
					countersMutex.Unlock()
				case model.RenameTypeExtension:
					newPath, err = pathgen.GenerateExtensionRenamePath(file, config)
				case model.RenameTypeCase:
					newPath, err = pathgen.GenerateCaseRenamePath(file, config)
				case model.RenameTypeInsertChar:
					newPath, err = pathgen.GenerateInsertCharRenamePath(file, config)
				case model.RenameTypeReplace:
					newPath, err = pathgen.GenerateReplaceRenamePath(file, config)
				case model.RenameTypeDeleteChar:
					newPath, err = pathgen.GenerateDeleteCharRenamePath(file, config)
				}

				if err != nil {
					resultChan <- struct {
						file string
						err  error
					}{
						file: file,
						err:  err,
					}
					continue
				}

				err = filestatus.RenameFile(file, newPath)
				if err == nil {
					global.Logs = append(global.Logs, global.RenameLog{
						Original: file,
						New:      newPath,
						Time:     time.Now().Format("2006-01-02 15:04:05"),
					})
				}
				resultChan <- struct {
					file string
					err  error
				}{
					file: file,
					err:  err,
				}
			}
		}()
	}

	// 发送文件到工作池
	go func() {
		for _, file := range files {
			fileChan <- file
		}
		close(fileChan)
		wg.Wait()
		close(resultChan)
	}()

	// 处理结果
	busyFiles := []string{}
	lengthErrorFiles := []string{}
	for result := range resultChan {
		if result.err != nil {
			if filestatus.IsFileBusyError(result.err) {
				busyFiles = append(busyFiles, result.file)
			} else if _, ok := result.err.(*ui.FilenameLengthError); ok {
				lengthErrorFiles = append(lengthErrorFiles, result.file)
			}
		}

		if pd.IsCancelled() {
			break
		}
	}

	pd.Hide()

	if pd.IsCancelled() {
		dialog.ShowInformation(dialogTr("warning"), tr("operation_cancelled"), window)
		return
	}

	// 显示文件名长度错误
	if len(lengthErrorFiles) > 0 {
		ui.ShowLengthErrorDialog(window, lengthErrorFiles)
		return
	}

	// 直接用弹窗列出未完成重命名的文件
	if len(busyFiles) > 0 {
		dialogcustomize.ShowMultiLineCopyDialog("error", tr("rename_failed_files"), busyFiles, window)
		return
	} else {
		if len(busyFiles) == 0 && len(lengthErrorFiles) == 0 {
			dialog.ShowInformation(dialogTr("success"), fmt.Sprintf(i18n.Tr("rename_success_count"), len(files)), window)
		}
	}
}
