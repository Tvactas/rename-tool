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

// =======================
// UI 组件结构体
// =======================
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

// =======================
// 线程安全封装
// =======================
func safeUI(f func()) {
	if fyne.CurrentApp() == nil {
		f()
		return
	}
	fyne.Do(f) // Fyne 官方线程安全执行
}

// =======================
// 初始化 UI
// =======================
func initRenameUI(config RenameUIConfig) (*RenameUIComponents, error) {
	global.MainWindow.Hide()

	window := global.MyApp.NewWindow(config.Title)
	window.Resize(fyne.NewSize(600, 500))
	window.SetFixedSize(false)
	window.SetCloseIntercept(func() {
		global.MyApp.Quit()
	})

	title := widget.NewLabelWithStyle(config.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: false})
	formatLabel := widget.NewLabel(tr("scan_format") + ": " + tr("scan_not_started"))
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
			formatLabel.SetText(tr("scan_format") + ": " + tr("scan_not_started"))
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

// =======================
// 按钮逻辑拆分
// =======================
func setupScanButton(ui *RenameUIComponents, config RenameUIConfig) *widget.Button {
	_ = config

	return widget.NewButton(tr("scan_format"), func() {
		if global.SelectedDir == "" {
			dialog.ShowInformation(dialogTr("error"), tr("please_select_dir"), ui.Window)
			return
		}

		formats, err := scan.ScanFormats(global.SelectedDir)
		if err != nil {
			safeUI(func() { ui.FormatLabel.SetText(tr("scan_format") + ": " + tr("scan_failed")) })
			return
		}

		if len(formats) == 0 {
			safeUI(func() { ui.FormatLabel.SetText(tr("scan_format") + ": " + tr("scan_no_files")) })
			return
		}

		safeUI(func() {
			ui.FormatLabel.SetText(fmt.Sprintf(tr("scan_format")+": "+tr("scan_found_formats"), len(formats)))
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
	})
}

func setupPreviewButton(ui *RenameUIComponents, config RenameUIConfig) *widget.Button {
	return widget.NewButton(buttonTr("preview"), func() {
		var selectedFormats []string
		for format, check := range ui.FormatChecks {
			if check.Checked {
				selectedFormats = append(selectedFormats, format)
			}
		}
		if len(selectedFormats) == 0 {
			dialog.ShowInformation(dialogTr("error"), tr("please_select_format"), ui.Window)
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
	var btn *widget.Button // 先声明
	btn = widget.NewButton(tr("rename"), func() {
		var selectedFormats []string
		for format, check := range ui.FormatChecks {
			if check.Checked {
				selectedFormats = append(selectedFormats, format)
			}
		}
		if len(selectedFormats) == 0 {
			dialog.ShowInformation(dialogTr("error"), tr("please_select_format"), ui.Window)
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

		btn.Disable() // ✅ 安全使用
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
	return widget.NewButton(tr("back"), func() {
		ui.Window.Close()
		global.MyApp.Settings().SetTheme(&theme.MainTheme{})
		global.MainWindow.Show()
	})
}

// =======================
// 整合事件绑定
// =======================
func setupRenameUIEvents(ui *RenameUIComponents, config RenameUIConfig) (scanBtn, previewBtn, renameBtn, backBtn *widget.Button) {
	scanBtn = setupScanButton(ui, config)
	previewBtn = setupPreviewButton(ui, config)
	renameBtn = setupRenameButton(ui, config)
	backBtn = setupBackButton(ui)
	return
}

// =======================
// 主入口 ShowRenameUI
// =======================
