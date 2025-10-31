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

// ShowInsertCharRename displays the character insertion interface
func ShowInsertCharRename() {
	// Create configuration form
	positionEntry := widget.NewEntry()
	insertTextEntry := widget.NewEntry()
	configForm := container.NewVBox(
		widget.NewLabel(buttonTr("insertPosition")),
		positionEntry,
		widget.NewLabel(buttonTr("insertText")),
		insertTextEntry,
	)

	// Create configuration builder
	configBuilder := func() model.RenameConfig {
		position, _ := strconv.Atoi(positionEntry.Text)
		return model.RenameConfig{
			Type:           model.RenameTypeInsertChar,
			InsertPosition: position,
			InsertText:     insertTextEntry.Text,
		}
	}

	// Create validation function
	validateConfig := func(config model.RenameConfig) error {
		// 1. 检查输入是否为数字（在configBuilder中已转换，但这里做二次校验）
		positionStr := positionEntry.Text
		if _, err := strconv.Atoi(positionStr); err != nil {
			return errors.New(tr("position_must_be_number"))
		}

		// 获取最短文件名长度
		minLen, err := dirpath.GetShortestFilenameLength(global.SelectedDir)
		if err != nil {
			return err
		}
		if config.InsertPosition > minLen {
			return errors.New(tr("position_exceeds_length"))
		}
		if config.InsertText == "" {
			return errors.New(tr("insert_text_empty"))
		}
		return nil
	}

	// Show rename interface
	ShowRenameUI(RenameUIConfig{
		Title:           buttonTr("insertLetter"),
		Window:          global.MainWindow,
		RenameType:      model.RenameTypeInsertChar,
		ConfigBuilder:   configBuilder,
		ValidateConfig:  validateConfig,
		AdditionalItems: []fyne.CanvasObject{configForm},
	})
}
