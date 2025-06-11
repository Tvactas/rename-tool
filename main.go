package main

//power by Tvacats
import (
	"embed"
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"errors"

	fs "rename-tool/common/FileStatus"
	sc "rename-tool/common/scan"
	"rename-tool/setting/config"
	gb "rename-tool/setting/global"
	"rename-tool/setting/i18n"
	"rename-tool/setting/model"
	"rename-tool/utils"

	"rename-tool/view"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed src/font/* src/img/*
var resourceFS embed.FS

type mainTheme struct{}
type otherTheme struct{}

func (m *mainTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameForeground {
		return color.Black
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m *mainTheme) Font(style fyne.TextStyle) fyne.Resource {
	return view.LoadFont(style)
}

func (m *mainTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *mainTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (m *otherTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameForeground {
		return color.Black
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m *otherTheme) Font(style fyne.TextStyle) fyne.Resource {
	return view.LoadDefaultFont()
}

func (m *otherTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *otherTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func main() {
	// 设置错误处理
	defer func() {
		if r := recover(); r != nil {
			logError(fmt.Errorf("panic: %v", r))
		}
	}()

	// 初始化应用
	gb.MyApp = app.NewWithID("com.yourdomain.renametool")
	gb.MainWindow = gb.MyApp.NewWindow(tr("title"))
	gb.MainWindow.Resize(fyne.NewSize(600, 400))
	gb.MainWindow.SetFixedSize(true) // 禁止调整窗口大小
	gb.MainWindow.SetMaster()        // 设置为主窗口

	// 获取当前目录
	dir, err := os.Getwd()
	if err == nil {
		gb.CurrentDir = dir
	} else {
		gb.CurrentDir = "."
		logError(fmt.Errorf("failed to get current directory: %v", err))
	}
	gb.SelectedDir = gb.CurrentDir

	// 设置自定义主题
	gb.MyApp.Settings().SetTheme(&mainTheme{})
	showMainMenu()
	gb.MainWindow.ShowAndRun()
}

func showMainMenu() {
	gb.MyApp.Settings().SetTheme(&mainTheme{})

	// 使用嵌入的图片资源
	imgResource := view.LoadImage("cat.png")
	var image *canvas.Image
	if imgResource == nil {
		image = canvas.NewImageFromFile("")
		logError(fmt.Errorf("failed to load cat.png"))
	} else {
		image = canvas.NewImageFromResource(imgResource)
	}
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(250, 380))

	// 优化按钮创建
	makeTextBtn := func(text string, onTap func()) fyne.CanvasObject {
		btn := widget.NewButton(text, onTap)
		btn.Importance = widget.LowImportance
		return container.NewHBox(btn, layout.NewSpacer())
	}

	// 使用预定义的按钮列表
	buttons := []struct {
		text   string
		action func()
	}{
		{tr("batch"), func() { showBatchRenameNormal() }},
		{tr("ext"), func() { showChangeExtension() }},
		{tr("upper"), func() { showRenameToCase("upper") }},
		{tr("lower"), func() { showRenameToCase("lower") }},
		{tr("titlecase"), func() { showRenameToCase("title") }},
		{tr("camel"), func() { showRenameToCase("camel") }},
		{tr("insert_char"), func() { showInsertCharRename() }},
		{tr("delete_char"), func() { showDeleteCharRename() }},
		{tr("regex_replace"), func() { showRegexReplace() }},
		{tr("undo"), undoRename},
		{tr("log"), saveLogs},
		{tr("exit"), func() { gb.MyApp.Quit() }},
	}

	// 创建按钮网格
	var buttonGridItems []fyne.CanvasObject
	for _, btn := range buttons {
		buttonGridItems = append(buttonGridItems, makeTextBtn(btn.text, btn.action))
	}
	buttonGrid := container.NewGridWithColumns(2, buttonGridItems...)

	// 优化布局
	rightBox := container.NewVBox(buttonGrid)
	mainContent := container.NewBorder(nil, nil, image, rightBox)
	centered := container.NewVBox(
		layout.NewSpacer(),
		mainContent,
		layout.NewSpacer(),
	)

	bgContent := setBackground(centered)
	langSelector := langSelect()

	header := container.NewHBox(
		langSelector,
		layout.NewSpacer(),
		widget.NewLabel(tr("title")),
	)

	content := container.NewVBox(
		header,
		bgContent,
	)

	gb.MainWindow.SetContent(content)
	gb.MainWindow.Show()
}

// 背景设置函数
func setBackground(content fyne.CanvasObject) fyne.CanvasObject {
	// 创建蓝到紫的线性渐变（左上到右下）
	grad1 := canvas.NewLinearGradient(
		color.RGBA{R: 0, G: 128, B: 255, A: 255}, // 蓝色
		color.RGBA{R: 128, G: 0, B: 255, A: 255}, // 紫色
		45,                                       // 角度，左上到右下
	)
	// 叠加紫到绿的半透明渐变
	grad2 := canvas.NewLinearGradient(
		color.RGBA{R: 128, G: 0, B: 255, A: 128}, // 半透明紫色
		color.RGBA{R: 0, G: 255, B: 128, A: 128}, // 半透明绿色
		45,
	)

	return container.NewStack(
		grad1,
		grad2,
		container.NewPadded(content),
	)
}

// 执行重命名操作
func performRename(window fyne.Window, config model.RenameConfig) {
	if config.SelectedDir == "" {
		dialog.ShowInformation(tr("error"), tr("please_select_dir"), window)
		return
	}

	// 获取文件列表
	files, err := getFiles(config.SelectedDir, config.Formats)
	if err != nil {
		dialog.ShowError(&AppError{
			Code:    "FILE_LIST_ERROR",
			Message: tr("error_getting_files"),
			Err:     err,
		}, window)
		return
	}

	// 检查重名
	if config.Type == model.RenameTypeReplace {
		duplicates, err := checkDuplicateNames(files, config)
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
				container.NewMax(textArea),
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
	pd := newProgressDialog(tr("rename"), window)
	pd.Show()

	// 使用工作池处理文件
	workerCount := runtime.NumCPU()
	fileChan := make(chan string, len(files))
	resultChan := make(chan struct {
		file string
		err  error
	}, len(files))

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
					newPath = generateBatchRenamePath(file, config, currentCounter)
				case model.RenameTypeExtension:
					newPath = generateExtensionRenamePath(file, config)
				case model.RenameTypeCase:
					newPath = generateCaseRenamePath(file, config)
				case model.RenameTypeInsertChar:
					newPath, err = generateInsertCharRenamePath(file, config)
					if err != nil {
						resultChan <- struct {
							file string
							err  error
						}{file, err}
						continue
					}
				case model.RenameTypeReplace:
					newPath, err = generateReplaceRenamePath(file, config)
					if err != nil {
						resultChan <- struct {
							file string
							err  error
						}{file, err}
						continue
					}
				case model.RenameTypeDeleteChar:
					newPath, err = generateDeleteCharRenamePath(file, config)
					if err != nil {
						resultChan <- struct {
							file string
							err  error
						}{file, err}
						continue
					}
				}

				err = renameFile(file, newPath)
				resultChan <- struct {
					file string
					err  error
				}{file, err}
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
			if fs.IsFileBusyError(result.err) {
				busyFiles = append(busyFiles, result.file)
			} else if _, ok := result.err.(*FilenameLengthError); ok {
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
		showLengthErrorDialog(window, lengthErrorFiles)
		return
	}

	// 显示文件占用错误
	if len(busyFiles) > 0 {
		showBusyFilesDialog(window, busyFiles)
		return
	}

	showSuccessMessage(window, config.Type, len(files))
}

// 生成批量重命名的新路径
func generateBatchRenamePath(file string, config model.RenameConfig, counter int) string {
	dirPath, oldName := filepath.Split(file)
	ext := filepath.Ext(oldName)
	nameWithoutExt := oldName[:len(oldName)-len(ext)]

	newName := ""

	// 如果从1开始编号，则counter加1
	if !config.StartFromZero {
		counter++
	}

	// 构建前缀序号
	if config.PrefixDigits > 0 {
		newName += fmt.Sprintf("%0*d", config.PrefixDigits, counter)
	}

	// 添加前缀文本
	newName += config.PrefixText

	// 保留原文件名
	if config.KeepOriginal {
		newName += nameWithoutExt
	}

	// 添加后缀文本
	newName += config.SuffixText

	// 构建后缀序号
	if config.SuffixDigits > 0 {
		newName += fmt.Sprintf("%0*d", config.SuffixDigits, counter)
	}

	// 添加扩展名
	newName += ext

	return filepath.Join(dirPath, newName)
}

// 生成扩展名修改的新路径
func generateExtensionRenamePath(file string, config model.RenameConfig) string {
	dirPath, oldName := filepath.Split(file)
	ext := filepath.Ext(oldName)
	nameWithoutExt := oldName[:len(oldName)-len(ext)]
	return filepath.Join(dirPath, nameWithoutExt+config.NewExtension)
}

// 生成大小写重命名的新路径
func generateCaseRenamePath(file string, config model.RenameConfig) string {
	dirPath, oldName := filepath.Split(file)
	newName := transformName(oldName, config.CaseType)
	return filepath.Join(dirPath, newName)
}

// 生成字符插入重命名的新路径
func generateInsertCharRenamePath(file string, config model.RenameConfig) (string, error) {
	dirPath, oldName := filepath.Split(file)
	ext := filepath.Ext(oldName)
	nameWithoutExt := oldName[:len(oldName)-len(ext)]

	// 将文件名转换为rune切片以正确处理Unicode字符
	runes := []rune(nameWithoutExt)
	if config.InsertPosition > len(runes) {
		return "", &FilenameLengthError{Files: []string{oldName}}
	}

	// 在指定位置插入文本
	newName := string(runes[:config.InsertPosition]) + config.InsertText + string(runes[config.InsertPosition:])
	return filepath.Join(dirPath, newName+ext), nil
}

// 生成正则替换重命名的新路径
func generateReplaceRenamePath(file string, config model.RenameConfig) (string, error) {
	generator := utils.GetPathGenerator(model.RenameTypeReplace)
	if generator == nil {
		return "", fmt.Errorf("unsupported rename type: %v", config.Type)
	}
	return generator.GeneratePath(file, config)
}

// 生成删除字符重命名的新路径
func generateDeleteCharRenamePath(file string, config model.RenameConfig) (string, error) {
	generator := utils.GetPathGenerator(model.RenameTypeDeleteChar)
	if generator == nil {
		return "", fmt.Errorf("unsupported rename type: %v", config.Type)
	}
	return generator.GeneratePath(file, config)
}

// 显示成功消息
func showSuccessMessage(window fyne.Window, renameType model.RenameType, count int) {
	var message string
	switch renameType {
	case model.RenameTypeBatch:
		message = fmt.Sprintf(tr("success_renamed")+" "+tr("files_count"), count)
	case model.RenameTypeExtension:
		message = fmt.Sprintf(tr("success_modified")+" "+tr("files_count"), count)
	case model.RenameTypeCase:
		message = fmt.Sprintf(tr("success_renamed")+" "+tr("files_count"), count)
	case model.RenameTypeInsertChar:
		message = fmt.Sprintf(tr("success_inserted")+" "+tr("files_count"), count)
	case model.RenameTypeReplace:
		message = fmt.Sprintf(tr("success_replaced")+" "+tr("files_count"), count)
	case model.RenameTypeDeleteChar:
		message = fmt.Sprintf(tr("success_deleted")+" "+tr("files_count"), count)
	}

	dialog.ShowInformation(tr("success"), message, window)
}

// ================= 批量重命名(普通)界面 =================
func showBatchRenameNormal() {
	gb.MyApp.Settings().SetTheme(&otherTheme{})
	gb.MainWindow.Hide()
	window := gb.MyApp.NewWindow(tr("batch_rename_title"))
	window.Resize(fyne.NewSize(600, 500))
	window.SetFixedSize(true)
	window.SetCloseIntercept(func() {
		gb.MyApp.Quit()
	})

	// 顶部标题
	title := widget.NewLabelWithStyle(tr("title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// 目录选择
	dirSelector := createDirSelector(window)

	// 文件格式扫描
	formatLabel := widget.NewLabel(tr("scan_format") + ": " + tr("scan_not_started"))

	// 创建格式列表容器
	formatListContainer := container.NewGridWithColumns(4)

	// 全选按钮
	selectAllBtn := widget.NewButton(tr("select_all"), nil)
	selectAllBtn.Hide() // 初始隐藏，扫描后显示

	// 格式单独计数选项
	formatSpecificNumbering := widget.NewCheck(tr("format_specific_numbering"), nil)
	formatSpecificNumbering.SetChecked(false)

	// 从0开始编号
	startFromZero := widget.NewCheck(tr("start_from_zero"), nil)
	startFromZero.SetChecked(true)

	// 两个选项同行
	optionRow := container.NewHBox(formatSpecificNumbering, startFromZero)

	// 存储格式复选框的映射
	formatChecks := make(map[string]*widget.Check)

	// 创建格式列表的滚动容器
	formatScroll := container.NewScroll(formatListContainer)
	formatScroll.SetMinSize(fyne.NewSize(0, 0))
	formatContainer := container.NewStack(formatScroll)
	formatContainer.Resize(fyne.NewSize(0, config.FormatListHeight))

	scanBtn := widget.NewButton(tr("scan_format"), func() {
		if gb.SelectedDir == "" {
			dialog.ShowInformation(tr("error"), tr("please_select_dir"), window)
			return
		}

		formats, err := sc.ScanFormats(gb.SelectedDir)
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
			check.SetChecked(true) // 默认选中
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
		// 强制刷新布局
		formatListContainer.Refresh()
		formatScroll.Refresh()
	})

	// 配置选项
	prefixDigits := widget.NewSelect([]string{"0", "1", "2", "3", "4", "5"}, nil)
	prefixDigits.SetSelected("0")
	prefixText := widget.NewEntry()
	prefixText.SetPlaceHolder(tr("prefix_placeholder"))

	keepOriginal := widget.NewCheck(tr("keep_original"), nil)
	keepOriginal.SetChecked(true)

	suffixText := widget.NewEntry()
	suffixText.SetPlaceHolder(tr("suffix_placeholder"))

	suffixDigits := widget.NewSelect([]string{"0", "1", "2", "3", "4", "5"}, nil)
	suffixDigits.SetSelected("0")

	// 预览区域
	previewLabel := widget.NewLabel(tr("preview") + ":")
	previewList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)

	// 底部按钮
	backBtn := widget.NewButton(tr("back"), func() {
		window.Close()
		gb.MyApp.Settings().SetTheme(&mainTheme{})
		gb.MainWindow.Show()
	})
	renameBtn := widget.NewButton(tr("rename"), func() {
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

		// 获取配置
		preDig, _ := strconv.Atoi(prefixDigits.Selected)
		sufDig, _ := strconv.Atoi(suffixDigits.Selected)
		prefix := prefixText.Text
		suffix := suffixText.Text
		keep := keepOriginal.Checked

		// 验证配置
		if !keep && preDig == 0 && sufDig == 0 {
			dialog.ShowError(errors.New(tr("error_no_prefix_suffix")), window)
			return
		}

		// 创建重命名配置
		config := model.RenameConfig{
			Type:                    model.RenameTypeBatch,
			SelectedDir:             gb.SelectedDir,
			Formats:                 selectedFormats,
			PrefixDigits:            preDig,
			PrefixText:              prefix,
			SuffixDigits:            sufDig,
			SuffixText:              suffix,
			KeepOriginal:            keep,
			FormatSpecificNumbering: formatSpecificNumbering.Checked,
			StartFromZero:           startFromZero.Checked,
		}

		// 执行重命名
		performRename(window, config)
	})

	// 添加预览按钮
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

		// 获取配置
		preDig, _ := strconv.Atoi(prefixDigits.Selected)
		sufDig, _ := strconv.Atoi(suffixDigits.Selected)
		prefix := prefixText.Text
		suffix := suffixText.Text
		keep := keepOriginal.Checked

		// 验证配置
		if !keep && preDig == 0 && sufDig == 0 {
			dialog.ShowInformation(tr("error"), tr("error_no_prefix_suffix"), window)
			return
		}

		// 创建重命名配置
		config := model.RenameConfig{
			Type:                    model.RenameTypeBatch,
			SelectedDir:             gb.SelectedDir,
			Formats:                 selectedFormats,
			PrefixDigits:            preDig,
			PrefixText:              prefix,
			SuffixDigits:            sufDig,
			SuffixText:              suffix,
			KeepOriginal:            keep,
			FormatSpecificNumbering: formatSpecificNumbering.Checked,
			StartFromZero:           startFromZero.Checked,
		}

		// 获取文件列表
		files, err := getFiles(gb.SelectedDir, selectedFormats)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		// 更新预览
		updatePreview(previewList, files, config)
	})

	// 布局
	configForm := widget.NewForm(
		widget.NewFormItem(tr("prefix_digits"), prefixDigits),
		widget.NewFormItem(tr("prefix_text"), prefixText),
		widget.NewFormItem("", keepOriginal),
		widget.NewFormItem(tr("suffix_text"), suffixText),
		widget.NewFormItem(tr("suffix_digits"), suffixDigits),
		widget.NewFormItem("", optionRow),
	)

	dirBox := container.NewHBox(dirSelector, scanBtn)
	formatBox := container.NewHBox(formatLabel, selectAllBtn)

	previewBox := container.NewBorder(previewLabel, nil, nil, nil, previewList)

	// 使用Border布局来更好地控制各个部分的大小
	content := container.NewBorder(
		container.NewVBox(
			title,
			widget.NewSeparator(),
			dirBox,
			widget.NewSeparator(),

			formatBox,
			formatContainer,
			widget.NewSeparator(),
		),
		container.NewHBox(layout.NewSpacer(), previewBtn, backBtn, renameBtn),
		nil,
		nil,
		container.NewVBox(
			configForm,
			widget.NewSeparator(),

			previewBox,
		),
	)

	window.SetContent(content)
	window.Show()
}

// ================= 扩展名修改界面 =================
func showChangeExtension() {
	gb.MyApp.Settings().SetTheme(&otherTheme{})
	gb.MainWindow.Hide()
	window := gb.MyApp.NewWindow(tr("change_ext_title"))
	window.Resize(fyne.NewSize(600, 500))
	window.SetFixedSize(true)
	window.SetCloseIntercept(func() {
		gb.MyApp.Quit()
	})

	title := widget.NewLabelWithStyle(tr("title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// 目录选择
	dirSelector := createDirSelector(window)

	// 文件格式选择
	formatLabel := widget.NewLabel(tr("scan_format"))
	formatSelect := widget.NewSelect([]string{}, nil)

	// 预览区域
	previewLabel := widget.NewLabel(tr("preview") + ":")
	previewList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)
	previewBox := container.NewBorder(previewLabel, nil, nil, nil, previewList)

	scanBtn := widget.NewButton(tr("scan_format"), func() {
		if gb.SelectedDir == "" {
			dialog.ShowInformation(tr("error"), tr("please_select_dir"), window)
			return
		}

		formats, err := sc.ScanFormats(gb.SelectedDir)
		if err != nil {
			formatSelect.Options = []string{tr("scan_failed")}
			return
		}

		if len(formats) == 0 {
			formatSelect.Options = []string{tr("scan_no_files")}
			return
		}

		formatSelect.Options = formats
		formatSelect.SetSelected(formats[0])
		formatSelect.Refresh()
	})

	// 新扩展名输入
	newExtLabel := widget.NewLabel(tr("new_extension"))
	newExtEntry := widget.NewEntry()
	newExtEntry.SetPlaceHolder(tr("example_ext"))
	newExtEntry.Resize(fyne.NewSize(600, newExtEntry.MinSize().Height))

	// 底部按钮
	backBtn := widget.NewButton(tr("back"), func() {
		window.Close()
		gb.MyApp.Settings().SetTheme(&mainTheme{})
		gb.MainWindow.Show()
	})
	renameBtn := widget.NewButton(tr("rename"), func() {
		oldExt := formatSelect.Selected
		if oldExt == "" {
			dialog.ShowInformation(tr("error"), tr("please_select_format"), window)
			return
		}

		newExt := newExtEntry.Text
		if newExt == "" {
			dialog.ShowInformation(tr("error"), tr("please_enter_new_extension"), window)
			return
		}

		// 确保扩展名以点开头
		if !strings.HasPrefix(newExt, ".") {
			newExt = "." + newExt
		}

		// 创建重命名配置
		config := model.RenameConfig{
			Type:         model.RenameTypeExtension,
			SelectedDir:  gb.SelectedDir,
			Formats:      []string{oldExt},
			NewExtension: newExt,
		}

		// 执行重命名
		performRename(window, config)
	})

	// 添加预览按钮
	previewBtn := widget.NewButton(tr("preview"), func() {
		oldExt := formatSelect.Selected
		if oldExt == "" {
			dialog.ShowInformation(tr("error"), tr("please_select_format"), window)
			return
		}

		newExt := newExtEntry.Text
		if newExt == "" {
			dialog.ShowInformation(tr("error"), tr("please_enter_new_extension"), window)
			return
		}

		// 确保扩展名以点开头
		if !strings.HasPrefix(newExt, ".") {
			newExt = "." + newExt
		}

		// 创建重命名配置
		config := model.RenameConfig{
			Type:         model.RenameTypeExtension,
			SelectedDir:  gb.SelectedDir,
			Formats:      []string{oldExt},
			NewExtension: newExt,
		}

		// 获取文件列表
		files, err := getFiles(gb.SelectedDir, []string{oldExt})
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		// 更新预览
		updatePreview(previewList, files, config)
	})

	// 布局
	dirBox := container.NewHBox(dirSelector, scanBtn)
	formatBox := container.NewHBox(formatLabel, formatSelect)

	// 修改新扩展名输入框的布局
	newExtBox := container.NewVBox(
		newExtLabel,
		newExtEntry,
	)

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		dirBox,
		widget.NewSeparator(),
		formatBox,
		newExtBox,
		widget.NewSeparator(),
		previewBox,
		widget.NewSeparator(),
		container.NewHBox(layout.NewSpacer(), previewBtn, backBtn, renameBtn),
	)

	window.SetContent(content)
	window.Show()
}

// ================= 大小写重命名界面 =================
func showRenameToCase(caseType string) {
	gb.MyApp.Settings().SetTheme(&otherTheme{})
	gb.MainWindow.Hide()
	window := gb.MyApp.NewWindow(tr(caseType + "_case_title"))
	window.Resize(fyne.NewSize(600, 500))
	window.SetFixedSize(true)
	window.SetCloseIntercept(func() {
		gb.MyApp.Quit()
	})

	title := widget.NewLabelWithStyle(tr(caseType+"_case_title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// 目录选择
	dirSelector := createDirSelector(window)

	// 预览区域
	previewLabel := widget.NewLabel(tr("preview") + ":")
	previewList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)

	// 扫描文件
	scanBtn := widget.NewButton(tr("scan_format"), func() {
		if gb.SelectedDir == "" {
			dialog.ShowInformation(tr("error"), tr("please_select_dir"), window)
			return
		}

		files, err := getFiles(gb.SelectedDir, nil)
		if err != nil {
			previewList.Length = func() int { return 1 }
			previewList.CreateItem = func() fyne.CanvasObject { return widget.NewLabel(tr("scan_failed")) }
			previewList.Refresh()
			return
		}

		if len(files) == 0 {
			previewList.Length = func() int { return 1 }
			previewList.CreateItem = func() fyne.CanvasObject { return widget.NewLabel(tr("scan_no_files")) }
			previewList.Refresh()
			return
		}

		previewList.Length = func() int { return len(files) }
		previewList.CreateItem = func() fyne.CanvasObject { return widget.NewLabel("") }
		previewList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
			_, oldName := filepath.Split(files[id])
			obj.(*widget.Label).SetText(oldName + " → " + transformName(oldName, caseType))
		}
		previewList.Refresh()
	})

	// 底部按钮
	backBtn := widget.NewButton(tr("back"), func() {
		window.Close()
		gb.MyApp.Settings().SetTheme(&mainTheme{})
		gb.MainWindow.Show()
	})
	renameBtn := widget.NewButton(tr("rename"), func() {
		// 创建重命名配置
		config := model.RenameConfig{
			Type:        model.RenameTypeCase,
			SelectedDir: gb.SelectedDir,
			CaseType:    caseType,
		}

		// 执行重命名
		performRename(window, config)
	})

	// 添加预览按钮
	previewBtn := widget.NewButton(tr("preview"), func() {
		// 创建重命名配置
		config := model.RenameConfig{
			Type:        model.RenameTypeCase,
			SelectedDir: gb.SelectedDir,
			CaseType:    caseType,
		}

		// 获取文件列表
		files, err := getFiles(gb.SelectedDir, nil)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		// 更新预览
		updatePreview(previewList, files, config)
	})

	// 布局
	dirBox := container.NewHBox(dirSelector, scanBtn)
	previewBox := container.NewBorder(previewLabel, nil, nil, nil, previewList)

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		dirBox,
		widget.NewSeparator(),
		previewBox,
		widget.NewSeparator(),
		container.NewHBox(layout.NewSpacer(), previewBtn, backBtn, renameBtn),
	)

	window.SetContent(content)
	window.Show()
}

// ================= 字符插入重命名界面 =================
func showInsertCharRename() {
	gb.MyApp.Settings().SetTheme(&otherTheme{})
	gb.MainWindow.Hide()
	window := gb.MyApp.NewWindow(tr("insert_char_title"))
	window.Resize(fyne.NewSize(600, 500))
	window.SetFixedSize(true)
	window.SetCloseIntercept(func() {
		gb.MyApp.Quit()
	})

	title := widget.NewLabelWithStyle(tr("title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// 目录选择
	dirSelector := createDirSelector(window)

	// 文件格式扫描
	formatLabel := widget.NewLabel(tr("scan_format") + ": " + tr("scan_not_started"))
	formatListContainer := container.NewGridWithColumns(4)
	selectAllBtn := widget.NewButton(tr("select_all"), nil)
	selectAllBtn.Hide()

	// 存储格式复选框的映射
	formatChecks := make(map[string]*widget.Check)

	// 创建格式列表的滚动容器
	formatScroll := container.NewScroll(formatListContainer)
	formatScroll.SetMinSize(fyne.NewSize(0, 0))
	formatContainer := container.NewStack(formatScroll)
	formatContainer.Resize(fyne.NewSize(0, config.FormatListHeight))

	// 插入位置输入
	positionEntry := widget.NewEntry()
	positionEntry.SetPlaceHolder(tr("insert_position_placeholder"))

	// 插入文本输入
	insertTextEntry := widget.NewEntry()
	insertTextEntry.SetPlaceHolder(tr("insert_text_placeholder"))

	scanBtn := widget.NewButton(tr("scan_format"), func() {
		if gb.SelectedDir == "" {
			dialog.ShowInformation(tr("error"), tr("please_select_dir"), window)
			return
		}

		formats, err := sc.ScanFormats(gb.SelectedDir)
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
	})

	// 预览区域
	previewLabel := widget.NewLabel(tr("preview") + ":")
	previewList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)

	// 底部按钮
	backBtn := widget.NewButton(tr("back"), func() {
		window.Close()
		gb.MyApp.Settings().SetTheme(&mainTheme{})
		gb.MainWindow.Show()
	})
	renameBtn := widget.NewButton(tr("rename"), func() {
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

		// 获取插入位置
		position, err := strconv.Atoi(positionEntry.Text)
		if err != nil || position < 0 {
			dialog.ShowInformation(tr("error"), tr("invalid_position"), window)
			return
		}

		// 获取插入文本
		insertText := insertTextEntry.Text
		if insertText == "" {
			dialog.ShowInformation(tr("error"), tr("please_enter_insert_text"), window)
			return
		}

		// 创建重命名配置
		config := model.RenameConfig{
			Type:           model.RenameTypeInsertChar,
			SelectedDir:    gb.SelectedDir,
			Formats:        selectedFormats,
			InsertPosition: position,
			InsertText:     insertText,
		}

		// 执行重命名
		performRename(window, config)
	})

	// 添加预览按钮
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

		// 获取插入位置
		position, err := strconv.Atoi(positionEntry.Text)
		if err != nil || position < 0 {
			dialog.ShowInformation(tr("error"), tr("invalid_position"), window)
			return
		}

		// 获取插入文本
		insertText := insertTextEntry.Text
		if insertText == "" {
			dialog.ShowInformation(tr("error"), tr("please_enter_insert_text"), window)
			return
		}

		// 创建重命名配置
		config := model.RenameConfig{
			Type:           model.RenameTypeInsertChar,
			SelectedDir:    gb.SelectedDir,
			Formats:        selectedFormats,
			InsertPosition: position,
			InsertText:     insertText,
		}

		// 获取文件列表
		files, err := getFiles(gb.SelectedDir, selectedFormats)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		// 更新预览
		updatePreview(previewList, files, config)
	})

	// 布局
	configForm := widget.NewForm(
		widget.NewFormItem(tr("insert_position"), positionEntry),
		widget.NewFormItem(tr("insert_text"), insertTextEntry),
	)

	dirBox := container.NewHBox(dirSelector, scanBtn)
	formatBox := container.NewHBox(formatLabel, selectAllBtn)
	previewBox := container.NewBorder(previewLabel, nil, nil, nil, previewList)

	content := container.NewBorder(
		container.NewVBox(
			title,
			widget.NewSeparator(),
			dirBox,
			widget.NewSeparator(),
			formatBox,
			formatContainer,
			widget.NewSeparator(),
		),
		container.NewHBox(layout.NewSpacer(), previewBtn, backBtn, renameBtn),
		nil,
		nil,
		container.NewVBox(
			configForm,
			widget.NewSeparator(),
			previewBox,
		),
	)

	window.SetContent(content)
	window.Show()
}

// ================= 撤销重命名 =================
func undoRename() {
	if len(gb.Logs) == 0 {
		dialog.ShowInformation(tr("info"), tr("no_undo_operations"), gb.MainWindow)
		return
	}

	busyFiles := []string{} // 记录被占用的文件
	successCount := 0

	for i := len(gb.Logs) - 1; i >= 0; i-- {
		log := gb.Logs[i]
		if _, err := os.Stat(log.New); err == nil {
			if err := os.Rename(log.New, log.Original); err == nil {
				successCount++
				// 从日志中移除已撤销的记录
				gb.Logs = append(gb.Logs[:i], gb.Logs[i+1:]...)
			} else if fs.IsFileBusyError(err) {
				busyFiles = append(busyFiles, log.New)
			}
		}
	}

	if len(busyFiles) > 0 {
		showBusyFilesDialog(gb.MainWindow, busyFiles)
	} else {
		dialog.ShowInformation(tr("success"), fmt.Sprintf(tr("undo_success"), successCount), gb.MainWindow)
	}
}

// ================= 保存日志 =================
func saveLogs() {
	if len(gb.Logs) == 0 {
		dialog.ShowInformation(tr("info"), tr("no_operations_to_save"), gb.MainWindow)
		return
	}

	content := ""
	for _, log := range gb.Logs {
		content += fmt.Sprintf("%s > %s [%s]\n", log.Original, log.New, log.Time)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(getLogPath()), os.ModePerm); err != nil {
		dialog.ShowError(fmt.Errorf(tr("error_creating_directory")+": %v", err), gb.MainWindow)
		return
	}

	// 使用临时文件进行写入
	tempPath := getLogPath() + ".tmp"
	if err := os.WriteFile(tempPath, []byte(content), 0644); err != nil {
		dialog.ShowError(fmt.Errorf(tr("error_saving_log")+": %v", err), gb.MainWindow)
		return
	}

	// 原子性地重命名临时文件
	if err := os.Rename(tempPath, getLogPath()); err != nil {
		// 清理临时文件
		os.Remove(tempPath)
		dialog.ShowError(fmt.Errorf(tr("error_saving_log")+": %v", err), gb.MainWindow)
		return
	}

	dialog.ShowInformation(tr("success"), fmt.Sprintf(tr("success_saved")+" "+tr("logs_count")+" "+tr("files_count_with_path"), len(gb.Logs), getLogPath()), gb.MainWindow)
}

// ============== 文件占用处理函数 ==============

// 显示被占用文件的对话框
func showBusyFilesDialog(window fyne.Window, busyFiles []string) {
	// 创建文本内容
	content := strings.Join(busyFiles, "\n")
	textArea := widget.NewMultiLineEntry()
	textArea.SetText(content)
	textArea.Wrapping = fyne.TextWrapWord
	textArea.Disable() // 设置为只读

	// 创建按钮
	copyBtn := widget.NewButton(tr("copy"), func() {
		window.Clipboard().SetContent(content)
		dialog.ShowInformation(tr("success"), tr("copy_success"), window)
	})

	killBtn := widget.NewButton(tr("kill_and_retry"), nil)
	retryBtn := widget.NewButton(tr("retry_no_kill"), nil)
	cancelBtn := widget.NewButton(tr("cancel"), nil)

	// 创建底部按钮容器
	bottomButtons := container.NewHBox(
		copyBtn,
		layout.NewSpacer(),
		killBtn,
		retryBtn,
		cancelBtn,
	)

	// 创建对话框内容
	dialogContent := container.NewBorder(
		widget.NewLabel(tr("busy_files_message")+":"),
		bottomButtons,
		nil,
		nil,
		container.NewMax(textArea),
	)

	dialog := dialog.NewCustom(
		tr("busy_files_title"),
		"",
		dialogContent,
		window,
	)

	// 设置按钮动作
	killBtn.OnTapped = func() {
		dialog.Hide()
		killProcessesAndRetry(window, busyFiles)
	}
	retryBtn.OnTapped = func() {
		dialog.Hide()
		retryRename(busyFiles, window)
	}
	cancelBtn.OnTapped = dialog.Hide

	dialog.Show()
}

// 尝试结束占用文件的进程并重试重命名
func killProcessesAndRetry(window fyne.Window, files []string) {
	// 尝试以管理员权限运行
	if !isAdmin() {
		runAsAdmin()
		return
	}

	// 作为管理员运行，尝试结束进程
	successCount := 0
	failedFiles := []string{}
	processInfo := make(map[string][]string) // 记录每个文件对应的进程信息

	for _, file := range files {
		processes := getProcessesUsingFile(file)
		if len(processes) == 0 {
			// 如果没有找到进程，可能是其他原因导致的占用
			failedFiles = append(failedFiles, file)
			continue
		}

		processInfo[file] = processes
		allKilled := true

		for _, proc := range processes {
			if !killProcess(proc) {
				allKilled = false
				break
			}
		}

		if allKilled {
			// 等待一小段时间确保进程完全结束
			time.Sleep(500 * time.Millisecond)
			// 成功结束进程后尝试重命名
			if retryRenameForFile(file) {
				successCount++
			} else {
				failedFiles = append(failedFiles, file)
			}
		} else {
			failedFiles = append(failedFiles, file)
		}
	}

	// 显示详细结果
	showKillResult(window, successCount, len(files), failedFiles, processInfo)
}

// 获取占用文件的进程信息
func getProcessesUsingFile(filePath string) []string {
	var processes []string

	// 根据操作系统选择不同的命令
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// 使用handle.exe来查找使用文件的进程
		cmd = exec.Command("handle", filePath)
	} else {
		// 在Linux/Unix上使用lsof
		cmd = exec.Command("lsof", filePath)
	}

	output, err := cmd.Output()
	if err != nil {
		logError(&AppError{
			Code:    "PROCESS_LIST_ERROR",
			Message: "Failed to get process list",
			Err:     err,
		})
		return processes
	}

	// 解析输出
	if runtime.GOOS == "windows" {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "pid:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					processes = append(processes, fmt.Sprintf("%s (PID: %s)", fields[0], fields[1]))
				}
			}
		}
	} else {
		// Linux/Unix 系统处理
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				processes = append(processes, fmt.Sprintf("%s (PID: %s)", fields[0], fields[1]))
			}
		}
	}

	return processes
}

// 结束指定进程
func killProcess(processInfo string) bool {
	// 从进程信息中提取PID
	pid := ""
	if strings.Contains(processInfo, "PID:") {
		pid = strings.TrimSpace(strings.Split(processInfo, "PID:")[1])
		pid = strings.Trim(pid, ")")
	}

	if pid == "" {
		return false
	}

	cmd := exec.Command("taskkill", "/F", "/PID", pid)
	return cmd.Run() == nil
}

// 显示进程结束结果
func showKillResult(window fyne.Window, successCount, totalCount int, failedFiles []string, processInfo map[string][]string) {
	var content strings.Builder
	content.WriteString(fmt.Sprintf(tr("success_killed")+" %d/%d 个文件\n\n", successCount, totalCount))

	if len(failedFiles) > 0 {
		content.WriteString(tr("some_files_may_still_be_in_use") + ":\n")
		for _, file := range failedFiles {
			content.WriteString(fmt.Sprintf("\n%s:\n", file))
			if processes, ok := processInfo[file]; ok {
				for _, proc := range processes {
					content.WriteString(fmt.Sprintf("  - %s\n", proc))
				}
			}
		}
	}

	dialog.ShowInformation(tr("operation_completed"), content.String(), window)
}

// 重试重命名单个文件
func retryRenameForFile(filePath string) bool {
	// 第一次尝试
	err := renameFile(filePath, filePath)
	if err == nil {
		return true
	}

	// 如果文件被占用，等待一段时间后重试
	if fs.IsFileBusyError(err) {
		time.Sleep(500 * time.Millisecond)
		err = renameFile(filePath, filePath)
		if err == nil {
			return true
		}
	}

	return false
}

// 重试重命名所有被占用的文件
func retryRename(files []string, window fyne.Window) {
	successCount := 0
	failedFiles := []string{}

	for _, file := range files {
		if retryRenameForFile(file) {
			successCount++
		} else {
			failedFiles = append(failedFiles, file)
		}
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf(tr("success_retried")+" %d/%d 个文件", successCount, len(files)))

	if len(failedFiles) > 0 {
		message.WriteString("\n\n" + tr("some_files_may_still_be_in_use") + ":\n")
		for _, file := range failedFiles {
			message.WriteString(fmt.Sprintf("  - %s\n", file))
		}
	}

	dialog.ShowInformation(tr("retry_result"), message.String(), window)
}

// 检查当前是否以管理员权限运行
func isAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

// 以管理员权限重新运行程序
func runAsAdmin() {
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()

	args := strings.Join(os.Args[1:], " ")

	// 在Windows上请求管理员权限
	cmd := exec.Command("cmd", "/C", "start", "runas", exe, args)
	cmd.Dir = cwd
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = cmd.Start() // 忽略错误，用户可能取消UAC提示
}

func getFiles(dir string, formats []string) ([]string, error) {
	var result []string
	formatSet := make(map[string]bool)
	for _, f := range formats {
		formatSet[f] = true
	}

	// 使用缓冲通道优化文件遍历
	fileChan := make(chan string, 100)
	errorChan := make(chan error, 1)

	go func() {
		defer close(fileChan)
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				// 尝试打开文件以确保可访问
				file, err := os.Open(path)
				if err != nil {
					// 如果文件被占用，记录错误但继续处理其他文件
					if fs.IsFileBusyError(err) {
						logError(fmt.Errorf("file busy: %s", path))
						return nil
					}
					return err
				}
				file.Close()

				fileChan <- path
			}
			return nil
		})
		if err != nil {
			errorChan <- err
		}
	}()

	// 处理文件
	for file := range fileChan {
		ext := strings.ToLower(filepath.Ext(file))
		if len(formats) == 0 || formatSet[ext] {
			result = append(result, file)
		}
	}

	// 检查错误
	select {
	case err := <-errorChan:
		return nil, err
	default:
	}

	sort.Strings(result)
	return result, nil
}

// 添加错误类型定义
type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprint(e.Message, ": ", e.Err)
	}
	return e.Message
}

// 优化文件操作函数
func renameFile(oldPath, newPath string) error {
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
			if fs.IsFileBusyError(err) {
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
			gb.Logs = append(gb.Logs, gb.RenameLog{
				Original: oldPath,
				New:      newPath,
				Time:     time.Now().Format("2006-01-02 15:04:05"),
			})
			return nil
		}

		if !fs.IsFileBusyError(err) {
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

// 优化进度对话框
type ProgressDialog struct {
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

func newProgressDialog(title string, window fyne.Window) *ProgressDialog {
	progress := widget.NewProgressBar()
	progress.Resize(fyne.NewSize(config.ProgressBarWidth, progress.MinSize().Height))

	status := widget.NewLabel("")
	status.Wrapping = fyne.TextWrapWord

	cancelBtn := widget.NewButton(tr("cancel"), nil)

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

	pd := &ProgressDialog{
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

func (pd *ProgressDialog) Show() {
	pd.dialog.Show()
}

func (pd *ProgressDialog) Hide() {
	pd.dialog.Hide()
	close(pd.updateCh) // 关闭更新通道
}

func (pd *ProgressDialog) Update(progress float64, status string) {
	select {
	case pd.updateCh <- struct {
		progress float64
		status   string
	}{progress, status}:
	default:
		// 如果通道已满，丢弃更新
	}
}

func (pd *ProgressDialog) IsCancelled() bool {
	pd.mu.RLock()
	defer pd.mu.RUnlock()
	return pd.cancelled
}

func init() {
	// 初始化资源加载器
	view.SetFontFS(resourceFS) // 设置字体文件系统
	view.Init()                // 初始化资源加载器

	// 设置语言变更回调
	i18n.GetManager().SetOnLangChange(func() {
		// 刷新主窗口
		if gb.MainWindow != nil {
			// 保存当前窗口大小
			size := gb.MainWindow.Canvas().Size()
			// 重新创建主菜单
			showMainMenu()
			// 恢复窗口大小
			gb.MainWindow.Resize(size)
		}
	})

	// 列出所有嵌入的文件
	files, err := view.ReadDir(".")
	if err != nil {
		logError(fmt.Errorf("failed to read embedded files: %v", err))
		return
	}
	for _, file := range files {
		logError(fmt.Errorf("embedded file: %s", file.Name()))
	}
}

// 获取应用数据目录
func getAppDataDir() string {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	// 在用户主目录下创建应用目录
	appDir := filepath.Join(homeDir, "."+config.AppName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "."
	}
	return appDir
}

// 获取日志文件路径
func getLogPath() string {
	appDir := getAppDataDir()
	logDir := filepath.Join(appDir, config.LogDir)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return filepath.Join(appDir, "rename.log")
	}
	return filepath.Join(logDir, "rename.log")
}

// 获取错误日志文件路径
func getErrorLogPath() string {
	return filepath.Join(getAppDataDir(), "error.log")
}

// 修改日志记录函数
func logError(err error) {
	if err == nil {
		return
	}

	logPath := getErrorLogPath()
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprint("[", timestamp, "] ", err, "\n")
	if _, err := f.WriteString(logEntry); err != nil {
		// 记录写入错误，但不返回错误以避免循环
		fmt.Fprintf(os.Stderr, "Failed to write to error log: %v\n", err)
	}
}

// 修改tr函数，使用i18n包的Tr函数
func tr(key string) string {
	return i18n.Tr(key)
}

func langSelect() fyne.CanvasObject {
	return i18n.LangSelect()
}

// 创建目录选择组件
func createDirSelector(window fyne.Window) fyne.CanvasObject {
	dirLabel := widget.NewLabel(tr("dir") + ": " + gb.SelectedDir)
	dirBtn := widget.NewButton(tr("select_dir"), func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				gb.SelectedDir = uri.Path()
				// 替换"父母"为".."
				gb.SelectedDir = strings.Replace(gb.SelectedDir, "父母", "..", -1)
				dirLabel.SetText(tr("dir") + ": " + gb.SelectedDir)
			}
		}, window).Show()
	})
	return container.NewHBox(dirLabel, dirBtn)
}

// 文件名转换函数
func transformName(name, caseType string) string {
	ext := filepath.Ext(name)
	nameWithoutExt := name[:len(name)-len(ext)]

	switch caseType {
	case "upper":
		return strings.ToUpper(nameWithoutExt) + ext
	case "lower":
		return strings.ToLower(nameWithoutExt) + ext
	case "title":
		words := strings.Fields(nameWithoutExt)
		for i, word := range words {
			if len(word) > 0 {
				words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
			}
		}
		return strings.Join(words, " ") + ext
	case "camel":
		words := strings.Fields(nameWithoutExt)
		for i, word := range words {
			if len(word) > 0 {
				words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
			}
		}
		return strings.Join(words, "") + ext
	default:
		return name
	}
}

// 添加新的错误类型
type FilenameLengthError struct {
	Files []string
}

func (e *FilenameLengthError) Error() string {
	return fmt.Sprintf("以下文件名的长度小于指定的插入位置：\n%s", strings.Join(e.Files, "\n"))
}

// 添加显示文件名长度错误的对话框
func showLengthErrorDialog(window fyne.Window, files []string) {
	// 创建文本内容
	content := strings.Join(files, "\n")
	textArea := widget.NewMultiLineEntry()
	textArea.SetText(content)
	textArea.Wrapping = fyne.TextWrapWord
	textArea.Disable() // 设置为只读

	// 创建按钮
	copyBtn := widget.NewButton(tr("copy"), func() {
		window.Clipboard().SetContent(content)
		dialog.ShowInformation(tr("success"), tr("copy_success"), window)
	})

	closeBtn := widget.NewButton(tr("close"), nil)

	// 创建对话框内容
	dialogContent := container.NewBorder(
		widget.NewLabel(tr("filename_length_error")+":"),
		container.NewHBox(copyBtn, layout.NewSpacer(), closeBtn),
		nil,
		nil,
		container.NewMax(textArea),
	)

	dialog := dialog.NewCustom(
		tr("error"),
		"",
		dialogContent,
		window,
	)

	// 设置关闭按钮动作
	closeBtn.OnTapped = dialog.Hide

	dialog.Show()
}

// 添加预览函数
func updatePreview(previewList *widget.List, files []string, config model.RenameConfig) {
	if len(files) == 0 {
		previewList.Length = func() int { return 1 }
		previewList.CreateItem = func() fyne.CanvasObject { return widget.NewLabel(tr("no_files_found")) }
		previewList.Refresh()
		return
	}

	previewList.Length = func() int { return len(files) }
	previewList.CreateItem = func() fyne.CanvasObject { return widget.NewLabel("") }
	previewList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
		oldPath := files[id]
		_, oldName := filepath.Split(oldPath)
		var newPath string
		var err error

		switch config.Type {
		case model.RenameTypeBatch:
			newPath = generateBatchRenamePath(oldPath, config, id)
		case model.RenameTypeExtension:
			newPath = generateExtensionRenamePath(oldPath, config)
		case model.RenameTypeCase:
			newPath = generateCaseRenamePath(oldPath, config)
		case model.RenameTypeInsertChar:
			newPath, err = generateInsertCharRenamePath(oldPath, config)
			if err != nil {
				obj.(*widget.Label).SetText(fmt.Sprintf("%s → %s", oldName, tr("error")))
				return
			}
		case model.RenameTypeReplace:
			newPath, err = generateReplaceRenamePath(oldPath, config)
			if err != nil {
				obj.(*widget.Label).SetText(fmt.Sprintf("%s → %s", oldName, tr("error")))
				return
			}
		case model.RenameTypeDeleteChar:
			newPath, err = generateDeleteCharRenamePath(oldPath, config)
			if err != nil {
				obj.(*widget.Label).SetText(fmt.Sprintf("%s → %s", oldName, tr("error")))
				return
			}
		}

		_, newName := filepath.Split(newPath)
		obj.(*widget.Label).SetText(fmt.Sprintf("%s → %s", oldName, newName))
	}
	previewList.Refresh()
}

// ================= 正则替换界面 =================
func showRegexReplace() {
	gb.MyApp.Settings().SetTheme(&otherTheme{})
	gb.MainWindow.Hide()
	window := gb.MyApp.NewWindow(tr("regex_replace"))
	window.Resize(fyne.NewSize(600, 500))
	window.SetFixedSize(true)
	window.SetCloseIntercept(func() {
		gb.MyApp.Quit()
	})

	// 标题
	title := widget.NewLabel(tr("regex_replace"))
	title.TextStyle = fyne.TextStyle{Bold: true}

	// 目录选择器
	dirSelector := widget.NewEntry()
	dirSelector.SetPlaceHolder(tr("select_dir"))
	dirSelector.Disable()

	// 添加目录选择按钮
	selectDirBtn := widget.NewButton(tr("select_dir"), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if uri != nil {
				gb.SelectedDir = uri.Path()
				dirSelector.SetText(gb.SelectedDir)
			}
		}, window)
	})

	// 预览列表
	previewLabel := widget.NewLabel(tr("preview"))
	previewList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)

	// 替换模式输入框
	replacePatternEntry := widget.NewEntry()
	replacePatternEntry.SetPlaceHolder(tr("replace_pattern_placeholder"))

	// 替换文本输入框
	replaceTextEntry := widget.NewEntry()
	replaceTextEntry.SetPlaceHolder(tr("replace_text"))

	// 使用正则表达式复选框
	useRegexCheck := widget.NewCheck(tr("use_regex"), nil)

	// 添加预览按钮
	previewBtn := widget.NewButton(tr("preview"), func() {
		if gb.SelectedDir == "" {
			dialog.ShowInformation(tr("error"), tr("please_select_dir"), window)
			return
		}

		// 获取文件列表
		files, err := getFiles(gb.SelectedDir, nil)
		if err != nil {
			previewList.Length = func() int { return 1 }
			previewList.CreateItem = func() fyne.CanvasObject { return widget.NewLabel(tr("scan_failed")) }
			previewList.Refresh()
			return
		}

		if len(files) == 0 {
			previewList.Length = func() int { return 1 }
			previewList.CreateItem = func() fyne.CanvasObject { return widget.NewLabel(tr("scan_no_files")) }
			previewList.Refresh()
			return
		}

		if replacePatternEntry.Text == "" {
			dialog.ShowInformation(tr("error"), tr("please_enter_replace_pattern"), window)
			return
		}

		config := model.RenameConfig{
			Type:           model.RenameTypeReplace,
			SelectedDir:    gb.SelectedDir,
			ReplacePattern: replacePatternEntry.Text,
			ReplaceText:    replaceTextEntry.Text,
			UseRegex:       useRegexCheck.Checked,
		}

		// 更新预览
		updatePreview(previewList, files, config)
	})

	// 重命名按钮
	renameBtn := widget.NewButton(tr("rename"), func() {
		if gb.SelectedDir == "" {
			dialog.ShowInformation(tr("error"), tr("please_select_dir"), window)
			return
		}

		if replacePatternEntry.Text == "" {
			dialog.ShowInformation(tr("error"), tr("please_enter_replace_pattern"), window)
			return
		}

		config := model.RenameConfig{
			Type:           model.RenameTypeReplace,
			SelectedDir:    gb.SelectedDir,
			ReplacePattern: replacePatternEntry.Text,
			ReplaceText:    replaceTextEntry.Text,
			UseRegex:       useRegexCheck.Checked,
		}

		performRename(window, config)
	})

	// 返回按钮
	backBtn := widget.NewButton(tr("back"), func() {
		window.Close()
		gb.MyApp.Settings().SetTheme(&mainTheme{})
		gb.MainWindow.Show()
	})

	// 布局
	configForm := widget.NewForm(
		widget.NewFormItem(tr("replace_pattern"), replacePatternEntry),
		widget.NewFormItem(tr("replace_text"), replaceTextEntry),
		widget.NewFormItem(tr("use_regex"), useRegexCheck),
	)

	dirBox := container.NewVBox(
		dirSelector,
		selectDirBtn,
	)
	previewBox := container.NewBorder(previewLabel, nil, nil, nil, previewList)

	content := container.NewBorder(
		container.NewVBox(
			title,
			widget.NewSeparator(),
			dirBox,
			widget.NewSeparator(),
		),
		container.NewHBox(layout.NewSpacer(), previewBtn, backBtn, renameBtn),
		nil,
		nil,
		container.NewVBox(
			configForm,
			widget.NewSeparator(),
			previewBox,
		),
	)

	window.SetContent(content)
	window.Show()
}

// 检查重名文件
func checkDuplicateNames(files []string, config model.RenameConfig) ([]string, error) {
	nameMap := make(map[string]string)
	var duplicates []string

	for _, file := range files {
		newPath, err := generateReplaceRenamePath(file, config)
		if err != nil {
			return nil, err
		}
		_, newName := filepath.Split(newPath)

		if oldFile, exists := nameMap[newName]; exists {
			duplicates = append(duplicates, fmt.Sprintf("%s 和 %s 将重命名为相同的名称: %s",
				filepath.Base(oldFile), filepath.Base(file), newName))
		} else {
			nameMap[newName] = file
		}
	}

	return duplicates, nil
}

// ================= 删除字符重命名界面 =================
func showDeleteCharRename() {
	gb.MyApp.Settings().SetTheme(&otherTheme{})
	gb.MainWindow.Hide()
	window := gb.MyApp.NewWindow(tr("delete_char_title"))
	window.Resize(fyne.NewSize(600, 500))
	window.SetFixedSize(true)
	window.SetCloseIntercept(func() {
		gb.MyApp.Quit()
	})

	title := widget.NewLabelWithStyle(tr("title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// 目录选择
	dirSelector := createDirSelector(window)

	// 文件格式扫描
	formatLabel := widget.NewLabel(tr("scan_format") + ": " + tr("scan_not_started"))
	formatListContainer := container.NewGridWithColumns(4)
	selectAllBtn := widget.NewButton(tr("select_all"), nil)
	selectAllBtn.Hide()

	// 存储格式复选框的映射
	formatChecks := make(map[string]*widget.Check)

	// 创建格式列表的滚动容器
	formatScroll := container.NewScroll(formatListContainer)
	formatScroll.SetMinSize(fyne.NewSize(0, 0))
	formatContainer := container.NewStack(formatScroll)
	formatContainer.Resize(fyne.NewSize(0, config.FormatListHeight))

	// 删除位置输入
	startPositionEntry := widget.NewEntry()
	startPositionEntry.SetPlaceHolder(tr("delete_start_position_placeholder"))

	// 删除长度输入
	deleteLengthEntry := widget.NewEntry()
	deleteLengthEntry.SetPlaceHolder(tr("delete_length_placeholder"))

	scanBtn := widget.NewButton(tr("scan_format"), func() {
		if gb.SelectedDir == "" {
			dialog.ShowInformation(tr("error"), tr("please_select_dir"), window)
			return
		}

		formats, err := sc.ScanFormats(gb.SelectedDir)
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
	})

	// 预览区域
	previewLabel := widget.NewLabel(tr("preview") + ":")
	previewList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)

	// 底部按钮
	backBtn := widget.NewButton(tr("back"), func() {
		window.Close()
		gb.MyApp.Settings().SetTheme(&mainTheme{})
		gb.MainWindow.Show()
	})
	renameBtn := widget.NewButton(tr("rename"), func() {
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

		// 获取删除起始位置
		startPos, err := strconv.Atoi(startPositionEntry.Text)
		if err != nil || startPos < 0 {
			dialog.ShowInformation(tr("error"), tr("invalid_start_position"), window)
			return
		}

		// 获取删除长度
		deleteLen, err := strconv.Atoi(deleteLengthEntry.Text)
		if err != nil || deleteLen <= 0 {
			dialog.ShowInformation(tr("error"), tr("invalid_delete_length"), window)
			return
		}

		// 创建重命名配置
		config := model.RenameConfig{
			Type:                model.RenameTypeDeleteChar,
			SelectedDir:         gb.SelectedDir,
			Formats:             selectedFormats,
			DeleteStartPosition: startPos,
			DeleteLength:        deleteLen,
		}

		// 执行重命名
		performRename(window, config)
	})

	// 添加预览按钮
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

		// 获取删除起始位置
		startPos, err := strconv.Atoi(startPositionEntry.Text)
		if err != nil || startPos < 0 {
			dialog.ShowInformation(tr("error"), tr("invalid_start_position"), window)
			return
		}

		// 获取删除长度
		deleteLen, err := strconv.Atoi(deleteLengthEntry.Text)
		if err != nil || deleteLen <= 0 {
			dialog.ShowInformation(tr("error"), tr("invalid_delete_length"), window)
			return
		}

		// 创建重命名配置
		config := model.RenameConfig{
			Type:                model.RenameTypeDeleteChar,
			SelectedDir:         gb.SelectedDir,
			Formats:             selectedFormats,
			DeleteStartPosition: startPos,
			DeleteLength:        deleteLen,
		}

		// 获取文件列表
		files, err := getFiles(gb.SelectedDir, selectedFormats)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		// 更新预览
		updatePreview(previewList, files, config)
	})

	// 布局
	configForm := widget.NewForm(
		widget.NewFormItem(tr("delete_start_position"), startPositionEntry),
		widget.NewFormItem(tr("delete_length"), deleteLengthEntry),
	)

	dirBox := container.NewHBox(dirSelector, scanBtn)
	formatBox := container.NewHBox(formatLabel, selectAllBtn)
	previewBox := container.NewBorder(previewLabel, nil, nil, nil, previewList)

	content := container.NewBorder(
		container.NewVBox(
			title,
			widget.NewSeparator(),
			dirBox,
			widget.NewSeparator(),
			formatBox,
			formatContainer,
			widget.NewSeparator(),
		),
		container.NewHBox(layout.NewSpacer(), previewBtn, backBtn, renameBtn),
		nil,
		nil,
		container.NewVBox(
			configForm,
			widget.NewSeparator(),
			previewBox,
		),
	)

	window.SetContent(content)
	window.Show()
}
