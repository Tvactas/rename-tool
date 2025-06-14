package pathgen

import (
	"regexp"
	"rename-tool/setting/model"
	"strings"
)

// ReplacePathGenerator 处理正则替换的路径生成
type ReplacePathGenerator struct {
	BasePathGenerator
}

// GeneratePath 生成正则替换后的新路径
func (g *ReplacePathGenerator) GeneratePath(file string, config model.RenameConfig) (string, error) {
	dirPath, nameWithoutExt, ext := g.splitPath(file)

	var newName string
	if config.UseRegex {
		re, err := regexp.Compile(config.ReplacePattern)
		if err != nil {
			return "", err
		}
		newName = re.ReplaceAllString(nameWithoutExt, config.ReplaceText)
	} else {
		newName = strings.ReplaceAll(nameWithoutExt, config.ReplacePattern, config.ReplaceText)
	}

	return g.joinPath(dirPath, newName, ext), nil
}
