package strife

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

var fontLoaderInitialized bool = false

type Glyph struct {
	texs map[int]*sdl.Texture
	w, h int32
}

func NewGlyph(w, h int32, Col *Color, tex *sdl.Texture) *Glyph {
	textures := map[int]*sdl.Texture{}
	textures[Col.AsHex()] = tex
	return &Glyph{
		textures,
		w, h,
	}
}

type Font struct {
	*ttf.Font
	CharCache map[rune]*Glyph
}

func LoadFont(path string, size int) (*Font, error) {
	if !fontLoaderInitialized {
		ttf.Init()
	}

	font, err := ttf.OpenFont(path, size)
	if err != nil {
		return nil, fmt.Errorf("Failed to load font at '%s'\n", path)
	}

	return &Font{
		font,
		map[rune]*Glyph{},
	}, nil
}

func (f *Font) Destroy() {
	for _, g := range f.CharCache {
		for _, tex := range g.texs {
			tex.Destroy()
		}
	}
	f.Font.Close()
}
