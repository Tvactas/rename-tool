package utils

import (
	"path/filepath"
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
func GetPathGenerator(renameType model.RenameType) PathGenerator {
	switch renameType {
	case model.RenameTypeExtension:
		return &ExtensionPathGenerator{}
	case model.RenameTypeCase:
		return &CasePathGenerator{}
	case model.RenameTypeInsertChar:
		return &InsertCharPathGenerator{}
	case model.RenameTypeReplace:
		return &ReplacePathGenerator{}
	case model.RenameTypeDeleteChar:
		return &DeleteCharPathGenerator{}
	default:
		return nil
	}
}
