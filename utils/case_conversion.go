package utils

import (
	"rename-tool/setting/global"
	"rename-tool/setting/model"
)

// ShowRenameToCase displays the case conversion interface
func ShowRenameToCase(caseType string) {
	// Create configuration builder
	configBuilder := func() model.RenameConfig {
		return model.RenameConfig{
			Type:     model.RenameTypeCase,
			CaseType: caseType,
		}
	}

	// Use common UI display
	ShowRenameUI(RenameUIConfig{
		Title:          tr(caseType + "_case_title"),
		Window:         global.MainWindow,
		RenameType:     model.RenameTypeCase,
		ConfigBuilder:  configBuilder,
		ValidateConfig: func(config model.RenameConfig) error { return nil },
	})
}
