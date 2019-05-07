package strife

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

var fontLoaderInitialized = false

type fontStyle int

// types of font style
const (
	plain fontStyle = iota
	underline
	bold
	strike
)

// glyph is a wrapper over a texture with dimension stored
type glyph struct {
	tex *sdl.Texture
	dim []int32
}

// glpyInfo contains the rune value, colour as a hex
// and the font style.
type glyphInfo struct {
	val   rune
	col   uint32
	style fontStyle
}

// asKey pretty prints the glpyhInfo
func (g *glyphInfo) asKey() string {
	return string(g.val) + fmt.Sprintf("%d", g.col) + fmt.Sprintf("%d", int(g.style))
}

// encode will build a glpyhInfo object from the given
// values
func encode(col uint32, style fontStyle, val rune) glyphInfo {
	return glyphInfo{
		val: val, col: col, style: style,
	}
}

// Font is a TrueTypeFont, stores the path
// as well as the texture cache for the glyphs
type Font struct {
	*ttf.Font
	path     string
	texCache map[string]*glyph
}

// DeriveFont will create a new font object from
// this font of a different size.
func (f *Font) DeriveFont(size int) (*Font, error) {
	return LoadFont(f.path, size)
}

func (f *Font) hasGlyph(g glyphInfo) (*glyph, bool) {
	if val, ok := f.texCache[g.asKey()]; ok {
		return val, true
	}
	return nil, false
}

func (f *Font) cache(g glyphInfo, texture *sdl.Texture, dim []int32) *glyph {
	// todo cache collision?
	glyph := &glyph{texture, dim}
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	f.texCache[g.asKey()] = glyph
	return glyph
}

// LoadFont will try and load the font from the given
// path of the given size.
func LoadFont(path string, size int) (*Font, error) {
	if !fontLoaderInitialized {
		ttf.Init()
	}

	font, err := ttf.OpenFont(path, size)
	if err != nil {
		return nil, fmt.Errorf("Failed to load font at '%s'", path)
	}

	return &Font{
		font,
		path,
		map[string]*glyph{},
	}, nil
}

// Destroy will destroy the given font
// as well as clear the texture cache.
func (f *Font) Destroy() {
	for _, glyph := range f.texCache {
		glyph.tex.Destroy()
	}
	f.Font.Close()
}
