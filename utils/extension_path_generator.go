package utils

import (
	"rename-tool/setting/model"
)

// ExtensionPathGenerator 处理扩展名修改的路径生成
type ExtensionPathGenerator struct {
	BasePathGenerator
}

// GeneratePath 生成修改扩展名后的新路径
func (g *ExtensionPathGenerator) GeneratePath(file string, config model.RenameConfig) (string, error) {
	dirPath, nameWithoutExt, _ := g.splitPath(file)
	return g.joinPath(dirPath, nameWithoutExt, config.NewExtension), nil
}
