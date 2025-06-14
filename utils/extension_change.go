package utils

import (
	"errors"
	"strings"

	"rename-tool/setting/global"
	"rename-tool/setting/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ShowChangeExtension displays the extension change interface
func ShowChangeExtension() {
	// Create configuration form
	newExtLabel := widget.NewLabel(tr("new_extension"))
	newExtEntry := widget.NewEntry()
	configForm := container.NewVBox(
		newExtLabel,
		newExtEntry,
	)

	// Create configuration builder
	configBuilder := func() model.RenameConfig {
		newExt := newExtEntry.Text
		if !strings.HasPrefix(newExt, ".") {
			newExt = "." + newExt
		}
		return model.RenameConfig{
			Type:         model.RenameTypeExtension,
			NewExtension: newExt,
		}
	}

	// Create validation function
	validateConfig := func(config model.RenameConfig) error {
		if config.NewExtension == "" {
			return errors.New(tr("new_extension_empty"))
		}
		return nil
	}

	// Show rename interface
	ShowRenameUI(RenameUIConfig{
		Title:           tr("change_extension_title"),
		Window:          global.MainWindow,
		RenameType:      model.RenameTypeExtension,
		ConfigBuilder:   configBuilder,
		ValidateConfig:  validateConfig,
		AdditionalItems: []fyne.CanvasObject{configForm},
	})
}
