package strife

import (
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

var fontLoaderIitialized bool = false

type Font struct {
	*ttf.Font
}

func LoadFont(path string) (*Font, bool) {
	if !fontLoaderIitialized {
		ttf.Init()
	}

	font, err := ttf.OpenFont(path, 14)
	if err != nil {
		return nil, true
	}
	return &Font{
		font,
	}, false
}

func (f *Font) Destroy() {
	f.Font.Close()
}
