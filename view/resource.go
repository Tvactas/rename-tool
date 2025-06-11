package view

import (
	"embed"
	"fmt"
	"rename-tool/common/log"
	"rename-tool/setting/global"

	"io/fs"

	"fyne.io/fyne/v2"
)

// 资源文件系统
var fontFS embed.FS

// 资源缓存
var (
	imageCache = make(map[string]fyne.Resource) // 图片资源缓存
	fontCache  = make(map[string]fyne.Resource) // 字体资源缓存
)

// 字体名称常量
const (
	FontJP      = "JP.TTF"
	FontTimes   = "TIMES.TTF"
	FontXingKai = "STXINGKA.TTF"
)

// SetFontFS 设置嵌入的文件系统
func SetFontFS(fs embed.FS) {
	fontFS = fs
}

// Init 初始化资源加载器
func Init() {
	// 预加载字体
	fonts := []string{FontTimes, FontXingKai, FontJP}
	for _, font := range fonts {
		if data, err := fontFS.ReadFile("src/font/" + font); err == nil {
			fontCache[font] = fyne.NewStaticResource(font, data)
		} else {
			log.LogError(fmt.Errorf("failed to preload font %s: %v", font, err))
		}
	}

	// 预加载图片
	images := []string{"cat.png"}
	for _, img := range images {
		if data, err := fontFS.ReadFile("src/img/" + img); err == nil {
			imageCache[img] = fyne.NewStaticResource(img, data)
		} else {
			log.LogError(fmt.Errorf("failed to preload image %s: %v", img, err))
		}
	}
}

// getFontNameByLang 根据语言获取对应的字体名称
func getFontNameByLang() string {
	switch global.Lang {
	case "zh":
		return FontXingKai
	case "ja":
		return FontJP
	case "en":
		fallthrough
	default:
		return FontTimes
	}
}

// LoadFont 根据语言加载对应的字体
func LoadFont(style fyne.TextStyle) fyne.Resource {
	fontName := getFontNameByLang()

	// 检查缓存
	if font, ok := fontCache[fontName]; ok {
		return font
	}

	// 从文件系统加载
	data, err := fontFS.ReadFile("src/font/" + fontName)
	if err != nil {
		log.LogError(fmt.Errorf("failed to load font %s: %v", fontName, err))
		return nil
	}

	// 创建资源并缓存
	font := fyne.NewStaticResource(fontName, data)
	fontCache[fontName] = font
	return font
}

// LoadDefaultFont 加载默认字体
func LoadDefaultFont() fyne.Resource {
	// 检查缓存
	if font, ok := fontCache[FontTimes]; ok {
		return font
	}

	// 从文件系统加载
	data, err := fontFS.ReadFile("src/font/" + FontTimes)
	if err != nil {
		log.LogError(fmt.Errorf("failed to load font %s: %v", FontTimes, err))
		return nil
	}

	// 创建资源并缓存
	font := fyne.NewStaticResource(FontTimes, data)
	fontCache[FontTimes] = font
	return font
}

// LoadImage 加载图片资源
func LoadImage(name string) fyne.Resource {
	// 检查缓存
	if img, ok := imageCache[name]; ok {
		return img
	}

	// 从文件系统加载
	data, err := fontFS.ReadFile("src/img/" + name)
	if err != nil {
		log.LogError(fmt.Errorf("failed to load image %s: %v", name, err))
		return nil
	}

	// 创建资源并缓存
	img := fyne.NewStaticResource(name, data)
	imageCache[name] = img
	return img
}

// ReadDir reads the directory named by dirname and returns a list of directory entries
func ReadDir(dirname string) ([]fs.DirEntry, error) {
	return fontFS.ReadDir(dirname)
}
