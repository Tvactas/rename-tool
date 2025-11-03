package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
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
		errorDiaLog(window, dialogTr("failGetFiles"))
		return
	}

	// 统一防重名预检（批量内部重复、命中磁盘已存在路径）
	if stop, err := antisamename.CheckAndShowConflicts(window, files, config); err != nil {
		dialog.ShowError(err, window)
		return
	} else if stop {
		return
	}

	// 检查重名（仅替换类型）
	if config.Type == model.RenameTypeReplace {
		duplicates, err := pathgen.CheckDuplicateNames(files, config)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if len(duplicates) > 0 {
			dialogcustomize.ShowMultiLineCopyDialog("error", dialogTr("duplicateNames"), duplicates, window)
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
				// 生成新路径
				newPath, err := generateRenamePath(file, config, &counter, &counterMutex, counters, &countersMutex)
				if err != nil {
					resultChan <- struct {
						file string
						err  error
					}{file: file, err: err}
					continue
				}

				// 执行重命名
				if err := filestatus.RenameFile(file, newPath); err != nil {
					resultChan <- struct {
						file string
						err  error
					}{file: file, err: err}
				} else {
					appendRenameLog(file, newPath)
					resultChan <- struct {
						file string
						err  error
					}{file: file, err: nil}
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
	errorResults := collectRenameResults(resultChan, pd)
	pd.Hide()

	if pd.IsCancelled() {
		warningDiaLog(window, dialogTr("operationCancelled"))
		return
	}

	// 显示错误或成功消息
	showRenameResults(window, errorResults, len(files))
}

// generateRenamePath 生成重命名路径
func generateRenamePath(file string, config model.RenameConfig, counter *int, counterMutex *sync.Mutex,
	counters map[string]int, countersMutex *sync.Mutex) (string, error) {
	switch config.Type {
	case model.RenameTypeBatch:
		// 批量重命名需要计数器和扩展名计数器
		counterMutex.Lock()
		currentCounter := *counter
		*counter++
		counterMutex.Unlock()

		countersMutex.Lock()
		ext := filepath.Ext(file)
		if _, exists := counters[ext]; !exists {
			counters[ext] = 0
		}
		newPath, err := pathgen.GenerateBatchRenamePath(file, config, currentCounter, counters)
		countersMutex.Unlock()
		return newPath, err

	case model.RenameTypeExtension:
		return pathgen.GenerateExtensionRenamePath(file, config)
	case model.RenameTypeCase:
		return pathgen.GenerateCaseRenamePath(file, config)
	case model.RenameTypeInsertChar:
		return pathgen.GenerateInsertCharRenamePath(file, config)
	case model.RenameTypeReplace:
		return pathgen.GenerateReplaceRenamePath(file, config)
	case model.RenameTypeDeleteChar:
		return pathgen.GenerateDeleteCharRenamePath(file, config)
	default:
		return "", fmt.Errorf("unsupported rename type: %v", config.Type)
	}
}

// errorResults 错误结果集合
type errorResults struct {
	busyFiles   []string
	lengthFiles []string
	otherErrors map[string]error
}

// collectRenameResults 收集重命名结果
func collectRenameResults(resultChan <-chan struct {
	file string
	err  error
}, pd *progress.Dialog) errorResults {
	results := errorResults{
		otherErrors: make(map[string]error),
	}

	for result := range resultChan {
		if result.err != nil {
			isLenErr := false
			if _, ok := result.err.(*ui.FilenameLengthError); ok {
				isLenErr = true
			}
			switch {
			case filestatus.IsFileBusyError(result.err):
				results.busyFiles = append(results.busyFiles, result.file)
			case isLenErr:
				results.lengthFiles = append(results.lengthFiles, result.file)
			default:
				results.otherErrors[result.file] = result.err
			}
		}

		if pd.IsCancelled() {
			break
		}
	}

	return results
}

// appendRenameLog 追加重命名日志
func appendRenameLog(original, newPath string) {
	global.Logs = append(global.Logs, global.RenameLog{
		Original: original,
		New:      newPath,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
	})
}

// showRenameResults 显示重命名结果
func showRenameResults(window fyne.Window, results errorResults, totalFiles int) {
	// 显示文件占用错误
	if len(results.busyFiles) > 0 {
		dialogcustomize.ShowMultiLineCopyDialog("error", "rename_failed_files", results.busyFiles, window)
		return
	}

	// 显示其他类型的错误（权限、磁盘空间等）
	if len(results.otherErrors) > 0 {
		dialogcustomize.ShowMultiLineErrorDialog("error", "rename_failed_files", results.otherErrors, window)
		return
	}

	// 所有文件都成功
	if len(results.busyFiles) == 0 && len(results.lengthFiles) == 0 && len(results.otherErrors) == 0 {
		successDiaLog(window, fmt.Sprintf(dialogTr("successRenameCount"), totalFiles))
	}
}
