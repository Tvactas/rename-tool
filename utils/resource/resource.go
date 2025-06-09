package resource

import (
	"embed"
	"fmt"
	"rename-tool/common/log"

	"fyne.io/fyne/v2"
)

var (
	// 资源缓存
	fontCache  = make(map[string]fyne.Resource)
	imageCache = make(map[string]fyne.Resource)
)

// Init 初始化资源加载器
func Init(fs embed.FS) {
	// 列出所有嵌入的文件
	files, err := fs.ReadDir(".")
	if err != nil {
		log.LogError(fmt.Errorf("failed to read embedded files: %v", err))
		return
	}
	for _, file := range files {
		log.LogError(fmt.Errorf("embedded file: %s", file.Name()))
	}
}

// LoadFont 加载字体资源
func LoadFont(fs embed.FS, name string) fyne.Resource {
	// 检查缓存
	if font, ok := fontCache[name]; ok {
		return font
	}

	// 尝试从嵌入的文件系统加载
	data, err := fs.ReadFile("src/font/" + name)
	if err != nil {
		log.LogError(fmt.Errorf("failed to load font %s: %v", name, err))
		return nil
	}

	// 创建资源并缓存
	font := fyne.NewStaticResource(name, data)
	fontCache[name] = font
	return font
}

// LoadImage 加载图片资源
func LoadImage(fs embed.FS, name string) fyne.Resource {
	// 检查缓存
	if img, ok := imageCache[name]; ok {
		return img
	}

	// 尝试从嵌入的文件系统加载
	data, err := fs.ReadFile("src/img/" + name)
	if err != nil {
		log.LogError(fmt.Errorf("failed to load image %s: %v", name, err))
		return nil
	}

	// 创建资源并缓存
	img := fyne.NewStaticResource(name, data)
	imageCache[name] = img
	return img
}
