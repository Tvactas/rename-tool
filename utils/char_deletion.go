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
		widget.NewLabel(buttonTr("deletePosition")),
		startPositionEntry,
		widget.NewLabel(buttonTr("deleteLength")),
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
		// 1. 检查输入是否为数字（在configBuilder中已转换，但这里做二次校验）
		startPositionStr := startPositionEntry.Text
		deleteLengthStr := deleteLengthEntry.Text
		if _, err := strconv.Atoi(startPositionStr); err != nil {
			return errors.New(textTr("isNotNumber"))
		}
		if _, err := strconv.Atoi(deleteLengthStr); err != nil {
			return errors.New(textTr("delLengthIsNotNumber"))
		}

		// 获取最短文件名长度
		minLen, err := dirpath.GetShortestFilenameLength(global.SelectedDir)
		if err != nil {
			return err
		}

		if config.DeleteStartPosition > minLen {
			return errors.New(textTr("positionExceedsLength"))
		}
		if config.DeleteLength < 0 {
			return errors.New(textTr("delLengthNegative"))
		}
		if config.DeleteStartPosition+config.DeleteLength > minLen {
			return errors.New(textTr("delExceedsFileLenght"))
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
