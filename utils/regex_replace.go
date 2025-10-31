package utils

import (
	"errors"

	"rename-tool/setting/global"
	"rename-tool/setting/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ShowRegexReplace displays the regex replace interface
func ShowRegexReplace() {
	// Create configuration form
	replacePatternEntry := widget.NewEntry()
	replaceTextEntry := widget.NewEntry()
	useRegexCheck := widget.NewCheck(buttonTr("useRegex"), nil)
	configForm := container.NewVBox(
		widget.NewLabel(tr("replace_pattern")),
		replacePatternEntry,
		widget.NewLabel(tr("replace_text")),
		replaceTextEntry,
		useRegexCheck,
	)

	// Create configuration builder
	configBuilder := func() model.RenameConfig {
		return model.RenameConfig{
			Type:           model.RenameTypeReplace,
			ReplacePattern: replacePatternEntry.Text,
			ReplaceText:    replaceTextEntry.Text,
			UseRegex:       useRegexCheck.Checked,
		}
	}

	// Create validation function
	validateConfig := func(config model.RenameConfig) error {
		if config.ReplacePattern == "" {
			return errors.New(tr("replace_pattern_empty"))
		}
		return nil
	}

	// Show rename interface
	ShowRenameUI(RenameUIConfig{
		Title:           buttonTr("regexReplace"),
		Window:          global.MainWindow,
		RenameType:      model.RenameTypeReplace,
		ConfigBuilder:   configBuilder,
		ValidateConfig:  validateConfig,
		AdditionalItems: []fyne.CanvasObject{configForm},
	})
}
