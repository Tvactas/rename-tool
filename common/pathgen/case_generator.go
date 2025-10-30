package pathgen

import (
	"path/filepath"
	"strings"
)

// TransformName 根据 caseType 转换文件名大小写
func (g *CasePathGenerator) TransformName(name, caseType string) string {
	if name == "" {
		return name
	}

	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	switch strings.ToLower(caseType) {
	case "upper":
		base = strings.ToUpper(base)
	case "lower":
		base = strings.ToLower(base)
	case "title":
		base = transformWords(base, true, false)
	case "camel":
		base = transformWords(base, true, true)
	default:
		// 不识别的类型，原样返回
		return name
	}
	return base + ext
}

// transformWords 将文件名按空格拆分并进行首字母转换
// titleCase: 首字母大写，其余小写
// concat: 是否拼接（camel 模式下拼接为单词）
func transformWords(input string, titleCase, concat bool) string {
	if input == "" {
		return input
	}

	words := strings.FieldsFunc(input, func(r rune) bool {
		return r == '_' || r == '-' || r == ' ' // 兼容多种分隔符
	})

	for i, w := range words {
		if len(w) == 0 {
			continue
		}
		if titleCase {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		} else {
			words[i] = strings.ToLower(w)
		}
	}

	if concat {
		return strings.Join(words, "")
	}
	return strings.Join(words, " ")
}
