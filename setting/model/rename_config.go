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
    Type                    RenameType
    SelectedDir             string
    Formats                 []string
    PrefixDigits            int
    PrefixText              string
    SuffixDigits            int
    SuffixText              string
    KeepOriginal            bool
    NewExtension            string
    CaseType                string
    InsertPosition          int
    InsertText              string
    ReplacePattern          string
    ReplaceText             string
    UseRegex                bool
    FormatSpecificNumbering bool
    StartFromZero           bool
    DeleteStartPosition     int
    DeleteLength            int
    Filename                string
}
