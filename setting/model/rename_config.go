package model

// RenameType 重命名类型
type RenameType string

const (
	RenameTypeBatch      RenameType = "batch"
	RenameTypeExtension  RenameType = "extension"
	RenameTypeCase       RenameType = "case"
	RenameTypeInsertChar RenameType = "insert_char"
	RenameTypeReplace    RenameType = "replace"
	RenameTypeDeleteChar RenameType = "delete_char"
)

// RenameConfig 重命名配置
type RenameConfig struct {
	Type                    RenameType `json:"type"`
	SelectedDir             string     `json:"selected_dir"`
	Formats                 []string   `json:"formats"`
	PrefixDigits            int        `json:"prefix_digits"`
	PrefixText              string     `json:"prefix_text"`
	SuffixDigits            int        `json:"suffix_digits"`
	SuffixText              string     `json:"suffix_text"`
	KeepOriginal            bool       `json:"keep_original"`
	NewExtension            string     `json:"new_extension"`
	CaseType                string     `json:"case_type"`
	InsertPosition          int        `json:"insert_position"`
	InsertText              string     `json:"insert_text"`
	ReplacePattern          string     `json:"replace_pattern"`
	ReplaceText             string     `json:"replace_text"`
	UseRegex                bool       `json:"use_regex"`
	FormatSpecificNumbering bool       `json:"format_specific_numbering"`
	StartFromZero           bool       `json:"start_from_zero"`
	DeleteStartPosition     int        `json:"delete_start_position"`
	DeleteLength            int        `json:"delete_length"`
}
