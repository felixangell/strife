package strife

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

var fontLoaderIitialized bool = false

type Glyph struct {
	*sdl.Texture
	w, h int32
	Col  *Color
}

type Font struct {
	*ttf.Font
	CharCache map[rune]*Glyph
}

func LoadFont(path string, size int) (*Font, error) {
	if !fontLoaderIitialized {
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
		g.Destroy()
	}
	f.Font.Close()
}
