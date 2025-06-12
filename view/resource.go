package view

import (
	"embed"
	"fmt"
	"rename-tool/common/log"
	"rename-tool/setting/global"

	"io/fs"

	"fyne.io/fyne/v2"
)

// Resource file system
// 资源文件系统
var fontFS embed.FS

// Resource cache
// 资源缓存
var (
	imageCache = make(map[string]fyne.Resource) // Image resource cache // 图片资源缓存
	fontCache  = make(map[string]fyne.Resource) // Font resource cache // 字体资源缓存
)

// Font name constants
// 字体名称常量
const (
	FontJP      = "JP.TTF"       // Japanese font // 日语字体
	FontTimes   = "TIMES.TTF"    // Times New Roman font // Times New Roman字体
	FontXingKai = "STXINGKA.TTF" // Xing Kai font // 行楷字体
)

// SetFontFS sets the embedded file system
// SetFontFS 设置嵌入的文件系统
func SetFontFS(fs embed.FS) {
	fontFS = fs
}

// Init initializes the resource loader and preloads fonts and images
// Init 初始化资源加载器，预加载字体和图片
func Init() {
	// Preload fonts
	// 预加载字体
	fonts := []string{FontTimes, FontXingKai, FontJP}
	for _, font := range fonts {
		if data, err := fontFS.ReadFile("src/font/" + font); err == nil {
			fontCache[font] = fyne.NewStaticResource(font, data)
		} else {
			log.LogError(fmt.Errorf("failed to preload font %s: %v", font, err))
		}
	}

	// Preload images
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

// GetFontNameByLang returns the appropriate font name based on the current language
// GetFontNameByLang 根据当前语言返回对应的字体名称
func GetFontNameByLang() string {
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

// LoadFont loads the appropriate font based on the current language
// LoadFont 根据当前语言加载对应的字体
func LoadFont(style fyne.TextStyle) fyne.Resource {
	fontName := GetFontNameByLang()

	// Check cache
	// 检查缓存
	if font, ok := fontCache[fontName]; ok {
		return font
	}

	// Load from file system
	// 从文件系统加载
	data, err := fontFS.ReadFile("src/font/" + fontName)
	if err != nil {
		log.LogError(fmt.Errorf("failed to load font %s: %v", fontName, err))
		return nil
	}

	// Create resource and cache it
	// 创建资源并缓存
	font := fyne.NewStaticResource(fontName, data)
	fontCache[fontName] = font
	return font
}

// LoadDefaultFont loads the default Times New Roman font
// LoadDefaultFont 加载默认的Times New Roman字体
func LoadDefaultFont() fyne.Resource {
	// Check cache
	// 检查缓存
	if font, ok := fontCache[FontTimes]; ok {
		return font
	}

	// Load from file system
	// 从文件系统加载
	data, err := fontFS.ReadFile("src/font/" + FontTimes)
	if err != nil {
		log.LogError(fmt.Errorf("failed to load font %s: %v", FontTimes, err))
		return nil
	}

	// Create resource and cache it
	// 创建资源并缓存
	font := fyne.NewStaticResource(FontTimes, data)
	fontCache[FontTimes] = font
	return font
}

// LoadImage loads an image resource by name
// LoadImage 根据名称加载图片资源
func LoadImage(name string) fyne.Resource {
	// Check cache
	// 检查缓存
	if img, ok := imageCache[name]; ok {
		return img
	}

	// Load from file system
	// 从文件系统加载
	data, err := fontFS.ReadFile("src/img/" + name)
	if err != nil {
		log.LogError(fmt.Errorf("failed to load image %s: %v", name, err))
		return nil
	}

	// Create resource and cache it
	// 创建资源并缓存
	img := fyne.NewStaticResource(name, data)
	imageCache[name] = img
	return img
}

// ReadDir reads the directory named by dirname and returns a list of directory entries
// ReadDir 读取指定目录并返回目录条目列表
func ReadDir(dirname string) ([]fs.DirEntry, error) {
	return fontFS.ReadDir(dirname)
}
