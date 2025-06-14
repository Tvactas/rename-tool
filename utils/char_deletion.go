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

// ShowDeleteCharRename displays the character deletion interface
func ShowDeleteCharRename() {
	// Create configuration form
	startPositionEntry := widget.NewEntry()
	deleteLengthEntry := widget.NewEntry()
	configForm := container.NewVBox(
		widget.NewLabel(tr("delete_start_position")),
		startPositionEntry,
		widget.NewLabel(tr("delete_length")),
		deleteLengthEntry,
	)

	// Create configuration builder
	configBuilder := func() model.RenameConfig {
		startPosition, _ := strconv.Atoi(startPositionEntry.Text)
		deleteLength, _ := strconv.Atoi(deleteLengthEntry.Text)
		return model.RenameConfig{
			Type:                model.RenameTypeDeleteChar,
			DeleteStartPosition: startPosition,
			DeleteLength:        deleteLength,
		}
	}

	// Create validation function
	validateConfig := func(config model.RenameConfig) error {
		if config.DeleteStartPosition < 0 {
			return errors.New(tr("position_negative"))
		}
		if config.DeleteLength <= 0 {
			return errors.New(tr("delete_length_invalid"))
		}
		return nil
	}

	// Show rename interface
	ShowRenameUI(RenameUIConfig{
		Title:           tr("delete_char_title"),
		Window:          global.MainWindow,
		RenameType:      model.RenameTypeDeleteChar,
		ConfigBuilder:   configBuilder,
		ValidateConfig:  validateConfig,
		AdditionalItems: []fyne.CanvasObject{configForm},
	})
}
