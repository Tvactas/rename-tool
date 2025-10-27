package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"rename-tool/common/dirpath"
	"rename-tool/common/filestatus"
	"rename-tool/common/pathgen"
	"rename-tool/common/preview"
	"rename-tool/common/progress"
	"rename-tool/common/scan"
	"rename-tool/common/theme"
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

// ShowRenameUI 显示重命名界面
func ShowRenameUI(config RenameUIConfig) {
	global.MyApp.Settings().SetTheme(&theme.OtherTheme{})
	global.MainWindow.Hide()
	window := global.MyApp.NewWindow(config.Title)
	window.Resize(fyne.NewSize(600, 500))
	window.SetFixedSize(false)
	window.SetCloseIntercept(func() {
		global.MyApp.Quit()
	})

	// 标题
	title := widget.NewLabelWithStyle(config.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// 目录选择
	formatLabel := widget.NewLabel(tr("scan_format") + ": " + tr("scan_not_started"))
	formatListContainer := container.NewGridWithColumns(4)
	selectAllBtn := widget.NewButton(tr("select_all"), nil)
	selectAllBtn.Hide()

	// 存储格式复选框的映射
	formatChecks := make(map[string]*widget.Check)

	// 创建格式列表的滚动容器（初始化时就创建）
	formatScroll := container.NewScroll(formatListContainer)
	formatScroll.SetMinSize(fyne.NewSize(0, 200))
	formatScroll.Resize(fyne.NewSize(0, 200))

	var onDirChanged func()
	onDirChanged = func() {
		// 路径变更时清空格式相关内容
		formatListContainer.Objects = nil
		formatChecks = make(map[string]*widget.Check)
		formatLabel.SetText(tr("scan_format") + ": " + tr("scan_not_started"))
		selectAllBtn.Hide()
		formatListContainer.Refresh()
		formatScroll.Refresh()
		window.Content().Refresh()
	}
	dirSelector := dirpath.CreateDirSelector(window, onDirChanged)

	// 扫描按钮
	scanBtn := widget.NewButton(tr("scan_format"), func() {
		if global.SelectedDir == "" {
			dialog.ShowInformation(tr("error"), tr("please_select_dir"), window)
			return
		}

		formats, err := scan.ScanFormats(global.SelectedDir)
		if err != nil {
			formatLabel.SetText(tr("scan_format") + ": " + tr("scan_failed"))
			return
		}

		if len(formats) == 0 {
			formatLabel.SetText(tr("scan_format") + ": " + tr("scan_no_files"))
			return
		}

		formatLabel.SetText(fmt.Sprintf(tr("scan_format")+": "+tr("scan_found_formats"), len(formats)))

		// 清空现有格式列表
		formatListContainer.Objects = nil
		formatChecks = make(map[string]*widget.Check)

		// 为每个格式创建复选框
		for _, format := range formats {
			check := widget.NewCheck(format, nil)
			check.SetChecked(true)
			formatChecks[format] = check
			formatListContainer.Add(check)
		}

		// 设置全选按钮功能
		selectAllBtn.OnTapped = func() {
			allChecked := true
			for _, check := range formatChecks {
				if !check.Checked {
					allChecked = false
					break
				}
			}

			for _, check := range formatChecks {
				check.SetChecked(!allChecked)
			}
		}
		selectAllBtn.Show()
		formatListContainer.Refresh()
		formatScroll.Refresh()
		window.Content().Refresh()

	})

	// 底部按钮
	backBtn := widget.NewButton(tr("back"), func() {
		window.Close()
		global.MyApp.Settings().SetTheme(&theme.MainTheme{})
		global.MainWindow.Show()
	})

	var renameBtn *widget.Button
	renameBtn = widget.NewButton(tr("rename"), func() {
		// 获取选中的格式
		var selectedFormats []string
		for format, check := range formatChecks {
			if check.Checked {
				selectedFormats = append(selectedFormats, format)
			}
		}

		if len(selectedFormats) == 0 {
			dialog.ShowInformation(tr("error"), tr("please_select_format"), window)
			return
		}

		// 创建重命名配置
		renameConfig := config.ConfigBuilder()
		renameConfig.Type = config.RenameType
		renameConfig.SelectedDir = global.SelectedDir
		renameConfig.Formats = selectedFormats

		// 验证配置
		if err := config.ValidateConfig(renameConfig); err != nil {
			dialog.ShowError(err, window)
			return
		}

		// 禁用重命名按钮
		renameBtn.Disable()

		// 执行重命名
		performRename(window, renameConfig)

		// 0.5秒后重新启用重命名按钮
		go func() {
			time.Sleep(500 * time.Millisecond)
			fyne.Do(func() {
				renameBtn.Enable()
			})
		}()
	})

	// 预览按钮
	previewBtn := widget.NewButton(tr("preview"), func() {
		// 获取选中的格式
		var selectedFormats []string
		for format, check := range formatChecks {
			if check.Checked {
				selectedFormats = append(selectedFormats, format)
			}
		}

		if len(selectedFormats) == 0 {
			dialog.ShowInformation(tr("error"), tr("please_select_format"), window)
			return
		}

		// 创建重命名配置
		renameConfig := config.ConfigBuilder()
		renameConfig.Type = config.RenameType
		renameConfig.SelectedDir = global.SelectedDir
		renameConfig.Formats = selectedFormats

		// 验证配置
		if err := config.ValidateConfig(renameConfig); err != nil {
			dialog.ShowError(err, window)
			return
		}

		// 获取文件列表
		files, err := dirpath.GetFiles(global.SelectedDir, selectedFormats)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		// 打开预览界面
		preview.ShowPreviewWindow(window, files, renameConfig)
	})

	// 布局
	dirBox := container.NewHBox(dirSelector, scanBtn)
	formatBox := container.NewHBox(formatLabel, selectAllBtn)

	// 创建主内容
	mainContent := container.NewVBox(
		title,
		widget.NewSeparator(),
		dirBox,
		widget.NewSeparator(),
		formatBox,
		formatScroll,
		widget.NewSeparator(),
	)

	// 添加额外组件
	if len(config.AdditionalItems) > 0 {
		mainContent.Add(container.NewVBox(config.AdditionalItems...))
		mainContent.Add(widget.NewSeparator())
	}

	// 底部按钮
	bottomButtons := container.NewHBox(layout.NewSpacer(), previewBtn, backBtn, renameBtn)
	mainContent.Add(bottomButtons)

	window.SetContent(mainContent)
	window.Show()
}

// tr 函数用于国际化
func tr(key string) string {
	return i18n.Tr(key)
}

func buttonTr(key string) string {
	return i18n.ButtonTr(key)
}

// performRename 执行重命名操作
func performRename(window fyne.Window, config model.RenameConfig) {
	if config.SelectedDir == "" {
		dialog.ShowInformation(tr("error"), tr("please_select_dir"), window)
		return
	}

	// 获取文件列表
	files, err := dirpath.GetFiles(config.SelectedDir, config.Formats)
	if err != nil {
		dialog.ShowError(&filestatus.AppError{
			Code:    "FILE_LIST_ERROR",
			Message: tr("error_getting_files"),
			Err:     err,
		}, window)
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
				dialog.ShowInformation(tr("success"), tr("copy_success"), window)
			})

			closeBtn := widget.NewButton(tr("close"), nil)

			dialogContent := container.NewBorder(
				widget.NewLabel(tr("error")+": "+tr("duplicate_names")),
				container.NewHBox(copyBtn, layout.NewSpacer(), closeBtn),
				nil,
				nil,
				container.NewStack(textArea),
			)

			dialog := dialog.NewCustom(
				tr("error"),
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
	pd := progress.NewDialog(tr("rename"), window)
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
		dialog.ShowInformation(tr("info"), tr("operation_cancelled"), window)
		return
	}

	// 显示文件名长度错误
	if len(lengthErrorFiles) > 0 {
		ui.ShowLengthErrorDialog(window, lengthErrorFiles)
		return
	}

	// 显示文件占用错误
	if len(busyFiles) > 0 {
		filestatus.ShowBusyFilesDialog(window, busyFiles)
	} else {
		if len(busyFiles) == 0 && len(lengthErrorFiles) == 0 {
			dialog.ShowInformation(i18n.Tr("success"), fmt.Sprintf(i18n.Tr("rename_success_count"), len(files)), window)
		}
	}
}
