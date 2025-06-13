package theme

import (
	"embed"
	"fmt"
	"image/color"
	"io/fs"

	"rename-tool/common/log"
	"rename-tool/setting/global"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// Resource file system
var fontFS embed.FS

// Resource cache
var (
	imageCache = make(map[string]fyne.Resource) // Image resource cache
	fontCache  = make(map[string]fyne.Resource) // Font resource cache
)

// Font name constants
const (
	FontJP      = "JP.TTF"       // Japanese font
	FontTimes   = "TIMES.TTF"    // Times New Roman font
	FontXingKai = "STXINGKA.TTF" // Xing Kai font
)

// SetFontFS sets the embedded file system
func SetFontFS(fs embed.FS) {
	fontFS = fs
}

// Init initializes the resource loader and preloads fonts and images
func Init() {
	// Preload fonts
	fonts := []string{FontTimes, FontXingKai, FontJP}
	for _, font := range fonts {
		if data, err := fontFS.ReadFile("src/font/" + font); err == nil {
			fontCache[font] = fyne.NewStaticResource(font, data)
		} else {
			log.LogError(fmt.Errorf("failed to preload font %s: %v", font, err))
		}
	}

	// Preload images
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
func LoadFont(style fyne.TextStyle) fyne.Resource {
	fontName := GetFontNameByLang()

	// Check cache
	if font, ok := fontCache[fontName]; ok {
		return font
	}

	// Load from file system
	data, err := fontFS.ReadFile("src/font/" + fontName)
	if err != nil {
		log.LogError(fmt.Errorf("failed to load font %s: %v", fontName, err))
		return nil
	}

	// Create resource and cache it
	font := fyne.NewStaticResource(fontName, data)
	fontCache[fontName] = font
	return font
}

// LoadDefaultFont loads the default Times New Roman font
func LoadDefaultFont() fyne.Resource {
	// Check cache
	if font, ok := fontCache[FontTimes]; ok {
		return font
	}

	// Load from file system
	data, err := fontFS.ReadFile("src/font/" + FontTimes)
	if err != nil {
		log.LogError(fmt.Errorf("failed to load font %s: %v", FontTimes, err))
		return nil
	}

	// Create resource and cache it
	font := fyne.NewStaticResource(FontTimes, data)
	fontCache[FontTimes] = font
	return font
}

// LoadImage loads an image resource by name
func LoadImage(name string) fyne.Resource {
	// Check cache
	if img, ok := imageCache[name]; ok {
		return img
	}

	// Load from file system
	data, err := fontFS.ReadFile("src/img/" + name)
	if err != nil {
		log.LogError(fmt.Errorf("failed to load image %s: %v", name, err))
		return nil
	}

	// Create resource and cache it
	img := fyne.NewStaticResource(name, data)
	imageCache[name] = img
	return img
}

// ReadDir reads the directory named by dirname and returns a list of directory entries
func ReadDir(dirname string) ([]fs.DirEntry, error) {
	return fontFS.ReadDir(dirname)
}

// SetBackground sets the background with gradient colors
func SetBackground(content fyne.CanvasObject) fyne.CanvasObject {
	// Create blue to purple linear gradient (top-left to bottom-right)
	grad1 := canvas.NewLinearGradient(
		color.RGBA{R: 0, G: 128, B: 255, A: 255}, // Blue
		color.RGBA{R: 128, G: 0, B: 255, A: 255}, // Purple
		45,                                       // Angle, top-left to bottom-right
	)
	// Overlay purple to green semi-transparent gradient
	grad2 := canvas.NewLinearGradient(
		color.RGBA{R: 128, G: 0, B: 255, A: 128}, // Semi-transparent purple
		color.RGBA{R: 0, G: 255, B: 128, A: 128}, // Semi-transparent green
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
