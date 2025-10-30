package pathgen

import (
	"path/filepath"
	"rename-tool/setting/model"
)

// CasePathGenerator 用于生成文件名大小写转换后的路径
type CasePathGenerator struct {
	BasePathGenerator
}

// GeneratePath 生成转换后路径
func (g *CasePathGenerator) GeneratePath(file string, config model.RenameConfig) (string, error) {
	dir, base, ext := g.splitPath(file)
	newName := g.TransformName(base+ext, config.CaseType)
	return filepath.Join(dir, newName), nil
}
