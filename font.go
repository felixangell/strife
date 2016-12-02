package strife

import (
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

type Font struct {
	*ttf.Font
}

func LoadFont(path string) (*Font, bool) {
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
