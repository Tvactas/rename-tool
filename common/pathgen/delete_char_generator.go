package pathgen

import (
	"fmt"
	"rename-tool/setting/model"
)

// DeleteCharPathGenerator 处理字符删除的路径生成
type DeleteCharPathGenerator struct {
	BasePathGenerator
}

// GeneratePath 生成删除字符后的新路径
func (g *DeleteCharPathGenerator) GeneratePath(file string, config model.RenameConfig) (string, error) {
	dirPath, nameWithoutExt, ext := g.splitPath(file)

	runes := []rune(nameWithoutExt)
	if config.DeleteStartPosition >= len(runes) {
		return "", fmt.Errorf("delete start position %d is out of range for filename length %d",
			config.DeleteStartPosition, len(runes))
	}

	if config.DeleteStartPosition+config.DeleteLength > len(runes) {
		return "", fmt.Errorf("delete length %d from position %d exceeds filename length %d",
			config.DeleteLength, config.DeleteStartPosition, len(runes))
	}

	newName := string(runes[:config.DeleteStartPosition]) + string(runes[config.DeleteStartPosition+config.DeleteLength:])
	return g.joinPath(dirPath, newName, ext), nil
}
