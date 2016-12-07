package strife

import (
	"fmt"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

var fontLoaderIitialized bool = false

type Font struct {
	*ttf.Font
}

func LoadFont(path string, size int) (*Font, error) {
	if !fontLoaderIitialized {
		ttf.Init()
	}

	font, err := ttf.OpenFont(path, size)
	if err != nil {
		return nil, fmt.Errorf("Failed to load font at '%s'\n", path)
	}

	return &Font{font}, nil
}

func (f *Font) Destroy() {
	f.Font.Close()
}
