package utils

import (
	"path/filepath"
	"strings"
)

// TransformName 根据指定的大小写类型转换文件名
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
