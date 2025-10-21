package theme

import (
	"embed"
	"fmt"
	"image/color"
	"io/fs"
	"sync"

	"rename-tool/common/log"
	"rename-tool/setting/global"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// 路径常量
const (
	fontPath  = "src/font/"
	imagePath = "src/img/"
)

// 字体名称常量
const (
	FontJP      = "JP.TTF"
	FontTimes   = "TIMES.TTF"
	FontXingKai = "STXINGKA.TTF"
)

// Resource file system
var fontFS embed.FS

// Resource cache with mutex protection
var (
	imageCache = make(map[string]fyne.Resource)
	fontCache  = make(map[string]fyne.Resource)
	cacheMu    sync.RWMutex
)

// SetFontFS sets the embedded file system
func SetFontFS(fs embed.FS) {
	fontFS = fs
}

// Init initializes the resource loader and preloads fonts and images
func Init() {
	// 预加载字体
	fonts := []string{FontTimes, FontXingKai, FontJP}
	for _, font := range fonts {
		if data, err := fontFS.ReadFile(fontPath + font); err == nil {
			cacheMu.Lock()
			fontCache[font] = fyne.NewStaticResource(font, data)
			cacheMu.Unlock()
		} else {
			log.LogError(fmt.Errorf("failed to preload font %s: %v", font, err))
		}
	}

	// 预加载图片
	images := []string{"cat.png"}
	for _, img := range images {
		if data, err := fontFS.ReadFile(imagePath + img); err == nil {
			cacheMu.Lock()
			imageCache[img] = fyne.NewStaticResource(img, data)
			cacheMu.Unlock()
		} else {
			log.LogError(fmt.Errorf("failed to preload image %s: %v", img, err))
		}
	}
}

// GetFontNameByLang returns the appropriate font name based on the current language
func GetFontNameByLang() string {
	switch global.Lang {
	case "zh":
		return FontXingKai
	case "ja":
		return FontJP
	default:
		return FontTimes
	}
}

// loadFontByName 统一的字体加载函数，失败时返回系统默认字体
func loadFontByName(fontName string) fyne.Resource {
	// 检查缓存
	cacheMu.RLock()
	if font, ok := fontCache[fontName]; ok {
		cacheMu.RUnlock()
		return font
	}
	cacheMu.RUnlock()

	// 从文件系统加载
	data, err := fontFS.ReadFile(fontPath + fontName)
	if err != nil {
		log.LogError(fmt.Errorf("failed to load font %s: %v", fontName, err))
		// 返回 Fyne 默认字体，而不是 nil
		return theme.DefaultTheme().Font(fyne.TextStyle{})
	}

	// 创建资源并缓存
	font := fyne.NewStaticResource(fontName, data)
	cacheMu.Lock()
	fontCache[fontName] = font
	cacheMu.Unlock()

	return font
}

// LoadFont loads the appropriate font based on the current language
func LoadFont(style fyne.TextStyle) fyne.Resource {
	font := loadFontByName(GetFontNameByLang())
	if font == nil {
		return theme.DefaultTheme().Font(style)
	}
	return font
}

// LoadDefaultFont loads the default Times New Roman font
func LoadDefaultFont() fyne.Resource {
	font := loadFontByName(FontTimes)
	if font == nil {
		return theme.DefaultTheme().Font(fyne.TextStyle{})
	}
	return font
}

// LoadImage loads an image resource by name
func LoadImage(name string) fyne.Resource {
	// 检查缓存
	cacheMu.RLock()
	if img, ok := imageCache[name]; ok {
		cacheMu.RUnlock()
		return img
	}
	cacheMu.RUnlock()

	// 从文件系统加载
	data, err := fontFS.ReadFile(imagePath + name)
	if err != nil {
		log.LogError(fmt.Errorf("failed to load image %s: %v", name, err))
		return nil
	}

	// 创建资源并缓存
	img := fyne.NewStaticResource(name, data)
	cacheMu.Lock()
	imageCache[name] = img
	cacheMu.Unlock()

	return img
}

// ReadDir reads the directory named by dirname and returns a list of directory entries
func ReadDir(dirname string) ([]fs.DirEntry, error) {
	return fontFS.ReadDir(dirname)
}

// SetBackground sets the background with gradient colors
func SetBackground(content fyne.CanvasObject) fyne.CanvasObject {
	grad1 := canvas.NewLinearGradient(
		color.RGBA{R: 0, G: 128, B: 255, A: 255},
		color.RGBA{R: 128, G: 0, B: 255, A: 255},
		45,
	)
	grad2 := canvas.NewLinearGradient(
		color.RGBA{R: 128, G: 0, B: 255, A: 128},
		color.RGBA{R: 0, G: 255, B: 128, A: 128},
		45,
	)

	return container.NewStack(
		grad1,
		grad2,
		container.NewPadded(content),
	)
}

// MainTheme implements the main theme
type MainTheme struct{}

func (m *MainTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameForeground {
		return color.Black
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m *MainTheme) Font(style fyne.TextStyle) fyne.Resource {
	return LoadFont(style)
}

func (m *MainTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *MainTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// OtherTheme implements the other theme
type OtherTheme struct{}

func (m *OtherTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameForeground {
		return color.Black
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m *OtherTheme) Font(style fyne.TextStyle) fyne.Resource {
	return LoadDefaultFont()
}

func (m *OtherTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *OtherTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
