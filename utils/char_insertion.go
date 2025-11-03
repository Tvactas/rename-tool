package utils

import (
	"errors"
	"path/filepath"
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
			return errors.New(textTr("isNotNumber"))
		}

		if config.InsertPosition < 0 {
			return errors.New(textTr("insertPositionNegative"))
		}

		if config.InsertText == "" {
			return errors.New(textTr("insertEmptyText"))
		}

		// 检查文件名长度
		if config.SelectedDir == "" {
			return errors.New(dialogTr("selectDirFirst"))
		}

		// 获取所有文件
		files, err := dirpath.GetFiles(config.SelectedDir, config.Formats)
		if err != nil {
			return err
		}

		// 检查哪些文件的文件名长度小于插入位置
		lengthErrorFiles := []string{}
		for _, file := range files {
			// 获取文件名（不包括扩展名）的 rune 长度
			baseName := filepath.Base(file)
			nameWithoutExt := baseName[:len(baseName)-len(filepath.Ext(baseName))]
			runes := []rune(nameWithoutExt)

			if config.InsertPosition > len(runes) {
				lengthErrorFiles = append(lengthErrorFiles, baseName)
			}
		}

		// 若存在长度不足的文件，返回可读错误
		if len(lengthErrorFiles) > 0 {
			return errors.New(textTr("insertPositionExceededLength"))
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
