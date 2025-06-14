package pathgen

import (
	"path/filepath"
	"rename-tool/setting/model"
	"strings"
)

// CasePathGenerator 处理大小写转换的路径生成
type CasePathGenerator struct {
	BasePathGenerator
}

// TransformName 转换文件名的大小写
func (g *CasePathGenerator) TransformName(name, caseType string) string {
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

// GeneratePath 生成大小写转换后的新路径
func (g *CasePathGenerator) GeneratePath(file string, config model.RenameConfig) (string, error) {
	dirPath, nameWithoutExt, ext := g.splitPath(file)
	newName := g.TransformName(nameWithoutExt+ext, config.CaseType)
	return filepath.Join(dirPath, newName), nil
}
