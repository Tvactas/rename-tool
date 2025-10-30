package utils

import (
	"errors"
	"strconv"

	"rename-tool/common/dirpath"
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
			Filename:            "your_filename_here",
		}
	}

	// Create validation function
	validateConfig := func(config model.RenameConfig) error {
		// 1. 检查输入是否为数字（在configBuilder中已转换，但这里做二次校验）
		startPositionStr := startPositionEntry.Text
		deleteLengthStr := deleteLengthEntry.Text
		if _, err := strconv.Atoi(startPositionStr); err != nil {
			return errors.New(tr("position_must_be_number"))
		}
		if _, err := strconv.Atoi(deleteLengthStr); err != nil {
			return errors.New(tr("delete_length_must_be_number"))
		}

		// 获取最短文件名长度
		minLen, err := dirpath.GetShortestFilenameLength(global.SelectedDir)
		if err != nil {
			return err
		}

		if config.DeleteStartPosition > minLen {
			return errors.New(tr("position_exceeds_length"))
		}
		if config.DeleteLength < 0 {
			return errors.New(tr("delete_length_negative"))
		}
		if config.DeleteStartPosition+config.DeleteLength > minLen {
			return errors.New(tr("delete_range_exceeds_length"))
		}
		return nil
	}

	// Show rename interface
	ShowRenameUI(RenameUIConfig{
		Title:           buttonTr("deleteLetter"),
		Window:          global.MainWindow,
		RenameType:      model.RenameTypeDeleteChar,
		ConfigBuilder:   configBuilder,
		ValidateConfig:  validateConfig,
		AdditionalItems: []fyne.CanvasObject{configForm},
	})
}
