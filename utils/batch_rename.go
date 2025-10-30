package utils

import (
	"errors"
	"strconv"

	"rename-tool/setting/global"
	"rename-tool/setting/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ShowBatchRenameNormal displays the batch rename interface
func ShowBatchRenameNormal() {
	// Create configuration options
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

	// Format specific numbering option
	formatSpecificNumbering := widget.NewCheck(tr("format_specific_numbering"), nil)
	formatSpecificNumbering.SetChecked(false)

	// Start from zero option
	startFromZero := widget.NewCheck(tr("start_from_zero"), nil)
	startFromZero.SetChecked(true)

	// Place options in the same row
	optionRow := container.NewHBox(formatSpecificNumbering, startFromZero)

	// Create configuration form
	configForm := widget.NewForm(
		widget.NewFormItem(tr("prefix_digits"), prefixDigits),
		widget.NewFormItem(tr("prefix_text"), prefixText),
		widget.NewFormItem("", keepOriginal),
		widget.NewFormItem(tr("suffix_text"), suffixText),
		widget.NewFormItem(tr("suffix_digits"), suffixDigits),
		widget.NewFormItem("", optionRow),
	)

	// Create configuration builder
	configBuilder := func() model.RenameConfig {
		preDig, _ := strconv.Atoi(prefixDigits.Selected)
		sufDig, _ := strconv.Atoi(suffixDigits.Selected)
		return model.RenameConfig{
			Type:                    model.RenameTypeBatch,
			PrefixDigits:            preDig,
			PrefixText:              prefixText.Text,
			SuffixDigits:            sufDig,
			SuffixText:              suffixText.Text,
			KeepOriginal:            keepOriginal.Checked,
			FormatSpecificNumbering: formatSpecificNumbering.Checked,
			StartFromZero:           startFromZero.Checked,
		}
	}

	// Create configuration validator
	validateConfig := func(config model.RenameConfig) error {
		if !config.KeepOriginal && config.PrefixDigits == 0 && config.SuffixDigits == 0 {
			return errors.New(tr("error_no_prefix_suffix"))
		}
		return nil
	}

	// Use common UI display
	ShowRenameUI(RenameUIConfig{
		Title:           buttonTr("sequenceRename"),
		Window:          global.MainWindow,
		RenameType:      model.RenameTypeBatch,
		ConfigBuilder:   configBuilder,
		ValidateConfig:  validateConfig,
		AdditionalItems: []fyne.CanvasObject{configForm},
	})
}
