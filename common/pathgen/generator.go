package pathgen

import (
	"fmt"
	"path/filepath"
	"strings"

	"rename-tool/setting/model"
)

// PathGenerator 定义路径生成器接口
type PathGenerator interface {
	GeneratePath(file string, config model.RenameConfig) (string, error)
}

// BasePathGenerator 提供基础路径处理功能
type BasePathGenerator struct{}

// splitPath 将文件路径分割为目录路径、文件名（不含扩展名）和扩展名
func (g *BasePathGenerator) splitPath(file string) (dirPath, nameWithoutExt, ext string) {
	dirPath, oldName := filepath.Split(file)
	ext = filepath.Ext(oldName)
	nameWithoutExt = oldName[:len(oldName)-len(ext)]
	return dirPath, nameWithoutExt, ext
}

// joinPath 将目录路径、文件名和扩展名组合成完整路径
func (g *BasePathGenerator) joinPath(dirPath, nameWithoutExt, ext string) string {
	return filepath.Join(dirPath, nameWithoutExt+ext)
}

// GetPathGenerator 工厂函数，根据重命名类型返回对应的生成器
func GetPathGenerator(renameType model.RenameType) (PathGenerator, error) {
	switch renameType {
	case model.RenameTypeExtension:
		return &ExtensionPathGenerator{}, nil
	case model.RenameTypeCase:
		return &CasePathGenerator{}, nil
	case model.RenameTypeInsertChar:
		return &InsertCharPathGenerator{}, nil
	case model.RenameTypeReplace:
		return &ReplacePathGenerator{}, nil
	case model.RenameTypeDeleteChar:
		return &DeleteCharPathGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported rename type: %v", renameType)
	}
}

// getFormatCounter 获取并更新指定格式的计数器
func getFormatCounter(ext string, counters map[string]int) int {
	// 如果是第一次遇到这个扩展名，重置计数器
	if _, exists := counters[ext]; !exists {
		counters[ext] = 0
	}
	// 返回当前计数，不递增
	return counters[ext]
}

// incrementFormatCounter 递增指定格式的计数器
func incrementFormatCounter(ext string, counters map[string]int) {
	counters[ext]++
}

// GenerateBatchRenamePath 生成批量重命名的新路径
func GenerateBatchRenamePath(file string, config model.RenameConfig, counter int, counters map[string]int) (string, error) {
	dirPath, oldName := filepath.Split(file)
	ext := filepath.Ext(oldName)
	nameWithoutExt := oldName[:len(oldName)-len(ext)]

	var parts []string

	// 构建前缀序号
	if config.PrefixDigits > 0 {
		var number int
		if config.FormatSpecificNumbering {
			number = getFormatCounter(ext, counters)
			if !config.StartFromZero {
				number++
			}
		} else {
			number = counter
			if !config.StartFromZero {
				number++
			}
		}
		parts = append(parts, fmt.Sprintf("%0*d", config.PrefixDigits, number))
	}

	// 添加前缀文本
	if config.PrefixText != "" {
		parts = append(parts, config.PrefixText)
	}

	// 保留原文件名
	if config.KeepOriginal {
		parts = append(parts, nameWithoutExt)
	}

	// 添加后缀文本
	if config.SuffixText != "" {
		parts = append(parts, config.SuffixText)
	}

	// 构建后缀序号
	if config.SuffixDigits > 0 {
		var number int
		if config.FormatSpecificNumbering {
			// 使用不同的计数器键来避免与前缀序号冲突
			suffixKey := ext + "_suffix"
			number = getFormatCounter(suffixKey, counters)
			if !config.StartFromZero {
				number++
			}
		} else {
			number = counter
			if !config.StartFromZero {
				number++
			}
		}
		parts = append(parts, fmt.Sprintf("%0*d", config.SuffixDigits, number))
	}

	// 组合新文件名
	newName := strings.Join(parts, "")
	if newName == "" {
		newName = nameWithoutExt
	}

	// 如果是格式特定编号，递增计数器
	if config.FormatSpecificNumbering {
		incrementFormatCounter(ext, counters)
		if config.SuffixDigits > 0 {
			incrementFormatCounter(ext+"_suffix", counters)
		}
	}

	return filepath.Join(dirPath, newName+ext), nil
}

// GenerateExtensionRenamePath 生成扩展名修改的新路径
func GenerateExtensionRenamePath(file string, config model.RenameConfig) (string, error) {
	generator, err := GetPathGenerator(model.RenameTypeExtension)
	if err != nil {
		return file, err
	}
	return generator.GeneratePath(file, config)
}

// GenerateCaseRenamePath 生成大小写重命名的新路径
func GenerateCaseRenamePath(file string, config model.RenameConfig) (string, error) {
	generator, err := GetPathGenerator(model.RenameTypeCase)
	if err != nil {
		return file, err
	}
	return generator.GeneratePath(file, config)
}

// GenerateInsertCharRenamePath 生成字符插入重命名的新路径
func GenerateInsertCharRenamePath(file string, config model.RenameConfig) (string, error) {
	generator, err := GetPathGenerator(model.RenameTypeInsertChar)
	if err != nil {
		return "", err
	}
	return generator.GeneratePath(file, config)
}

// GenerateReplaceRenamePath 生成正则替换重命名的新路径
func GenerateReplaceRenamePath(file string, config model.RenameConfig) (string, error) {
	generator, err := GetPathGenerator(model.RenameTypeReplace)
	if err != nil {
		return "", err
	}
	return generator.GeneratePath(file, config)
}

// GenerateDeleteCharRenamePath 生成删除字符重命名的新路径
func GenerateDeleteCharRenamePath(file string, config model.RenameConfig) (string, error) {
	generator, err := GetPathGenerator(model.RenameTypeDeleteChar)
	if err != nil {
		return "", err
	}
	return generator.GeneratePath(file, config)
}

// CheckDuplicateNames 检查重名文件
func CheckDuplicateNames(files []string, config model.RenameConfig) ([]string, error) {
	nameMap := make(map[string]string)
	var duplicates []string

	for _, file := range files {
		newPath, err := GenerateReplaceRenamePath(file, config)
		if err != nil {
			return nil, fmt.Errorf("failed to generate path for %s: %w", file, err)
		}
		_, newName := filepath.Split(newPath)

		if oldFile, exists := nameMap[newName]; exists {
			duplicates = append(duplicates, fmt.Sprintf("%s 和 %s 将重命名为相同的名称: %s",
				filepath.Base(oldFile), filepath.Base(file), newName))
		} else {
			nameMap[newName] = file
		}
	}

	return duplicates, nil
}
