package utils

import (
	"fmt"
	"rename-tool/common/dirpath"
	"rename-tool/common/preview"
	"rename-tool/common/scan"
	"rename-tool/common/theme"
	"rename-tool/setting/global"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type RenameUIComponents struct {
	Window              fyne.Window
	Title               *widget.Label
	FormatLabel         *widget.Label
	FormatListContainer *fyne.Container
	FormatChecks        map[string]*widget.Check
	SelectAllBtn        *widget.Button
	FormatScroll        *container.Scroll
	DirSelector         fyne.CanvasObject
}

func safeUI(f func()) {
	if fyne.CurrentApp() == nil {
		f()
		return
	}
	fyne.Do(f)
}

func initRenameUI(config RenameUIConfig) (*RenameUIComponents, error) {
	global.MainWindow.Hide()

	window := global.MyApp.NewWindow(config.Title)
	window.Resize(fyne.NewSize(600, 500))
	window.SetFixedSize(false)
	window.SetCloseIntercept(func() {
		global.MyApp.Quit()
	})

	title := widget.NewLabelWithStyle(config.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: false})
	formatLabel := widget.NewLabel(buttonTr("scanFormat") + ": " + buttonTr("scanNotStart"))
	formatListContainer := container.NewGridWithColumns(4)
	selectAllBtn := widget.NewButton(tr("select_all"), nil)
	selectAllBtn.Hide()
	formatChecks := make(map[string]*widget.Check)

	formatScroll := container.NewScroll(formatListContainer)
	formatScroll.SetMinSize(fyne.NewSize(0, 200))
	formatScroll.Resize(fyne.NewSize(0, 200))

	onDirChanged := func() {
		safeUI(func() {
			formatListContainer.Objects = nil
			formatChecks = make(map[string]*widget.Check)
			formatLabel.SetText(buttonTr("scanFormat") + ": " + buttonTr("scanNotStart"))
			selectAllBtn.Hide()
			formatListContainer.Refresh()
			formatScroll.Refresh()
			window.Content().Refresh()
		})
	}

	dirSelector := dirpath.CreateDirSelector(window, onDirChanged)

	return &RenameUIComponents{
		Window:              window,
		Title:               title,
		FormatLabel:         formatLabel,
		FormatListContainer: formatListContainer,
		FormatChecks:        formatChecks,
		SelectAllBtn:        selectAllBtn,
		FormatScroll:        formatScroll,
		DirSelector:         dirSelector,
	}, nil
}

//
// ===== 提取扫描逻辑部分 =====
//

// 纯扫描逻辑（与 UI 无关）
func doScanFormats(dir string) ([]string, error) {
	if dir == "" {
		return nil, fmt.Errorf("no directory selected")
	}
	return scan.ScanFormats(dir)
}

// UI 更新逻辑（和 UI 状态绑定）
func updateFormatListUI(ui *RenameUIComponents, formats []string) {
	safeUI(func() {
		if len(formats) == 0 {
			ui.FormatLabel.SetText(buttonTr("scanFormat") + ": " + tr("scan_no_files"))
			ui.SelectAllBtn.Hide()
			return
		}

		ui.FormatLabel.SetText(fmt.Sprintf(buttonTr("scanFormat")+": "+tr("scan_found_formats"), len(formats)))
		ui.FormatListContainer.Objects = nil
		ui.FormatChecks = make(map[string]*widget.Check)

		for _, format := range formats {
			check := widget.NewCheck(format, nil)
			check.SetChecked(true)
			ui.FormatChecks[format] = check
			ui.FormatListContainer.Add(check)
		}

		ui.SelectAllBtn.OnTapped = func() {
			allChecked := true
			for _, check := range ui.FormatChecks {
				if !check.Checked {
					allChecked = false
					break
				}
			}
			for _, check := range ui.FormatChecks {
				check.SetChecked(!allChecked)
			}
		}
		ui.SelectAllBtn.Show()
		ui.FormatListContainer.Refresh()
		ui.FormatScroll.Refresh()
		ui.Window.Content().Refresh()
	})
}

func setupScanButton(ui *RenameUIComponents, config RenameUIConfig) *widget.Button {
	_ = config
	return widget.NewButton(buttonTr("scanFormat"), func() {
		if global.SelectedDir == "" {
			errorDiaLog(ui.Window, dialogTr("selectDirFirst"))
			return
		}

		formats, err := doScanFormats(global.SelectedDir)
		if err != nil {
			safeUI(func() {
				ui.FormatLabel.SetText(buttonTr("scanFormat") + ": " + tr("scan_failed"))
			})
			return
		}

		updateFormatListUI(ui, formats)
	})
}

//
// ===== 其他按钮逻辑保持不动 =====
//

func setupPreviewButton(ui *RenameUIComponents, config RenameUIConfig) *widget.Button {
	return widget.NewButton(buttonTr("preview"), func() {
		var selectedFormats []string
		for format, check := range ui.FormatChecks {
			if check.Checked {
				selectedFormats = append(selectedFormats, format)
			}
		}
		if len(selectedFormats) == 0 {
			errorDiaLog(ui.Window, dialogTr("selectFormat"))
			return
		}

		renameConfig := config.ConfigBuilder()
		renameConfig.Type = config.RenameType
		renameConfig.SelectedDir = global.SelectedDir
		renameConfig.Formats = selectedFormats

		if err := config.ValidateConfig(renameConfig); err != nil {
			dialog.ShowError(err, ui.Window)
			return
		}

		files, err := dirpath.GetFiles(global.SelectedDir, selectedFormats)
		if err != nil {
			dialog.ShowError(err, ui.Window)
			return
		}

		preview.ShowPreviewWindow(ui.Window, files, renameConfig)
	})
}

func setupRenameButton(ui *RenameUIComponents, config RenameUIConfig) *widget.Button {
	var btn *widget.Button
	btn = widget.NewButton(buttonTr("implement"), func() {
		var selectedFormats []string
		for format, check := range ui.FormatChecks {
			if check.Checked {
				selectedFormats = append(selectedFormats, format)
			}
		}
		if len(selectedFormats) == 0 {
			errorDiaLog(ui.Window, dialogTr("selectFormat"))
			return
		}

		renameConfig := config.ConfigBuilder()
		renameConfig.Type = config.RenameType
		renameConfig.SelectedDir = global.SelectedDir
		renameConfig.Formats = selectedFormats

		if err := config.ValidateConfig(renameConfig); err != nil {
			dialog.ShowError(err, ui.Window)
			return
		}

		btn.Disable()
		performRename(ui.Window, renameConfig)

		time.AfterFunc(500*time.Millisecond, func() {
			safeUI(func() {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   tr("rename_done"),
					Content: tr("rename_success"),
				})
				btn.Enable()
			})
		})
	})
	return btn
}

func setupBackButton(ui *RenameUIComponents) *widget.Button {
	return widget.NewButton(buttonTr("back"), func() {
		ui.Window.Close()
		global.MyApp.Settings().SetTheme(&theme.MainTheme{})
		global.MainWindow.Show()
	})
}

func setupRenameUIEvents(ui *RenameUIComponents, config RenameUIConfig) (scanBtn, previewBtn, renameBtn, backBtn *widget.Button) {
	scanBtn = setupScanButton(ui, config)
	previewBtn = setupPreviewButton(ui, config)
	renameBtn = setupRenameButton(ui, config)
	backBtn = setupBackButton(ui)
	return
}
