package pathgen

import (
	"fmt"
	"path/filepath"
	"strings"

	"rename-tool/setting/model"
	"rename-tool/utils"
)

// 生成批量重命名的新路径
func GenerateBatchRenamePath(file string, config model.RenameConfig, counter int) string {
	dirPath, oldName := filepath.Split(file)
	ext := filepath.Ext(oldName)
	nameWithoutExt := oldName[:len(oldName)-len(ext)]

	newName := ""

	// 如果从1开始编号，则counter加1
	if !config.StartFromZero {
		counter++
	}

	// 构建前缀序号
	if config.PrefixDigits > 0 {
		newName += fmt.Sprintf("%0*d", config.PrefixDigits, counter)
	}

	// 添加前缀文本
	newName += config.PrefixText

	// 保留原文件名
	if config.KeepOriginal {
		newName += nameWithoutExt
	}

	// 添加后缀文本
	newName += config.SuffixText

	// 构建后缀序号
	if config.SuffixDigits > 0 {
		newName += fmt.Sprintf("%0*d", config.SuffixDigits, counter)
	}

	// 添加扩展名
	newName += ext

	return filepath.Join(dirPath, newName)
}

// 生成扩展名修改的新路径
func GenerateExtensionRenamePath(file string, config model.RenameConfig) string {
	dirPath, oldName := filepath.Split(file)
	ext := filepath.Ext(oldName)
	nameWithoutExt := oldName[:len(oldName)-len(ext)]
	return filepath.Join(dirPath, nameWithoutExt+config.NewExtension)
}

// 生成大小写重命名的新路径
func GenerateCaseRenamePath(file string, config model.RenameConfig) string {
	dirPath, oldName := filepath.Split(file)
	newName := TransformName(oldName, config.CaseType)
	return filepath.Join(dirPath, newName)
}

// 生成字符插入重命名的新路径
func GenerateInsertCharRenamePath(file string, config model.RenameConfig) (string, error) {
	dirPath, oldName := filepath.Split(file)
	ext := filepath.Ext(oldName)
	nameWithoutExt := oldName[:len(oldName)-len(ext)]

	// 将文件名转换为rune切片以正确处理Unicode字符
	runes := []rune(nameWithoutExt)
	if config.InsertPosition > len(runes) {
		return "", &FilenameLengthError{Files: []string{oldName}}
	}

	// 在指定位置插入文本
	newName := string(runes[:config.InsertPosition]) + config.InsertText + string(runes[config.InsertPosition:])
	return filepath.Join(dirPath, newName+ext), nil
}

// 生成正则替换重命名的新路径
func GenerateReplaceRenamePath(file string, config model.RenameConfig) (string, error) {
	generator := utils.GetPathGenerator(model.RenameTypeReplace)
	if generator == nil {
		return "", fmt.Errorf("unsupported rename type: %v", config.Type)
	}
	return generator.GeneratePath(file, config)
}

// 生成删除字符重命名的新路径
func GenerateDeleteCharRenamePath(file string, config model.RenameConfig) (string, error) {
	generator := utils.GetPathGenerator(model.RenameTypeDeleteChar)
	if generator == nil {
		return "", fmt.Errorf("unsupported rename type: %v", config.Type)
	}
	return generator.GeneratePath(file, config)
}

// 文件名转换函数
func TransformName(name, caseType string) string {
	ext := filepath.Ext(name)
	nameWithoutExt := name[:len(name)-len(ext)]

	switch caseType {
	case "upper":
		return strings.ToUpper(nameWithoutExt) + ext
	case "lower":
		return strings.ToLower(nameWithoutExt) + ext
	case "title":
		words := strings.Fields(nameWithoutExt)
		for i, word := range words {
			if len(word) > 0 {
				words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
			}
		}
		return strings.Join(words, " ") + ext
	case "camel":
		words := strings.Fields(nameWithoutExt)
		for i, word := range words {
			if len(word) > 0 {
				words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
			}
		}
		return strings.Join(words, "") + ext
	default:
		return name
	}
}

// 添加新的错误类型
type FilenameLengthError struct {
	Files []string
}

func (e *FilenameLengthError) Error() string {
	return fmt.Sprintf("以下文件名的长度小于指定的插入位置：\n%s", strings.Join(e.Files, "\n"))
}
