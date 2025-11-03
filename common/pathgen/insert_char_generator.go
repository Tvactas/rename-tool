package pathgen

import (
	"fmt"
	"path/filepath"
	"rename-tool/setting/model"
)

// InsertCharPathGenerator 处理字符插入的路径生成
type InsertCharPathGenerator struct {
	BasePathGenerator
}

// GeneratePath 生成插入字符后的新路径
func (g *InsertCharPathGenerator) GeneratePath(file string, config model.RenameConfig) (string, error) {
	dirPath, nameWithoutExt, ext := g.splitPath(file)

	// 将文件名转换为rune切片以正确处理Unicode字符
	runes := []rune(nameWithoutExt)
	if config.InsertPosition > len(runes) {
		return "", fmt.Errorf("insert position %d exceeds filename length %d for %s", config.InsertPosition, len(runes), filepath.Base(file))
	}

	// 在指定位置插入文本
	newName := string(runes[:config.InsertPosition]) + config.InsertText + string(runes[config.InsertPosition:])
	return g.joinPath(dirPath, newName, ext), nil
}
