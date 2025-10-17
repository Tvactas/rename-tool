package i18n

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// LanguageOption 表示语言选项
type LanguageOption struct {
	DisplayName string // 显示名称
	Code        string // 语言代码
}

// I18nManager 管理多语言支持
type I18nManager struct {
	currentLang  string
	onLangChange func()
	languages    []LanguageOption
}

var manager = &I18nManager{
	currentLang: "en", // 默认使用英文
	languages: []LanguageOption{
		{DisplayName: "中文", Code: "zh"},
		{DisplayName: "English", Code: "en"},
		{DisplayName: "日本語", Code: "ja"},
	},
}

// GetManager 获取I18nManager的单例实例
func GetManager() *I18nManager {
	return manager
}

// SetLanguage 设置当前语言
func (i *I18nManager) SetLanguage(lang string) {
	if i.currentLang != lang {
		i.currentLang = lang
		if i.onLangChange != nil {
			i.onLangChange()
		}
	}
}

// GetLanguage 获取当前语言
func (i *I18nManager) GetLanguage() string {
	return i.currentLang
}

// SetOnLangChange 设置语言变更回调
func (i *I18nManager) SetOnLangChange(callback func()) {
	i.onLangChange = callback
}

// Tr 获取当前语言的翻译文本
func Tr(key string) string {
	return Translations[manager.currentLang][key]
}

func LogTr(key string) string {
	return log_translations[manager.currentLang][key]
}

func ButtonTr(key string) string {
	return button_translations[manager.currentLang][key]
}

// LangSelect 创建语言选择器组件
func LangSelect() fyne.CanvasObject {
	// 创建显示名称列表
	displayNames := make([]string, len(manager.languages))
	for i, lang := range manager.languages {
		displayNames[i] = lang.DisplayName
	}

	langSelect := widget.NewSelect(displayNames, func(selected string) {
		// 根据显示名称查找对应的语言代码
		for _, lang := range manager.languages {
			if lang.DisplayName == selected {
				manager.SetLanguage(lang.Code)
				break
			}
		}
	})

	// 设置当前语言
	currentLang := manager.GetLanguage()
	for _, lang := range manager.languages {
		if lang.Code == currentLang {
			langSelect.SetSelected(lang.DisplayName)
			break
		}
	}

	return langSelect
}
