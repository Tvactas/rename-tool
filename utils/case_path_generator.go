package utils

import (
	"rename-tool/setting/model"
)

// CasePathGenerator 处理大小写转换的路径生成
type CasePathGenerator struct {
	BasePathGenerator
}

// GeneratePath 生成大小写转换后的新路径
func (g *CasePathGenerator) GeneratePath(file string, config model.RenameConfig) (string, error) {
	dirPath, nameWithoutExt, ext := g.splitPath(file)
	return g.joinPath(dirPath, TransformName(nameWithoutExt, config.CaseType), ext), nil
}
