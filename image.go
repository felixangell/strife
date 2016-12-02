package gfx

import (
	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type Image struct {
	*sdl.Texture
}

func LoadImage(path string) (*Image, bool) {
	surface, err := img.Load(path)
	if err != nil {
		return nil, true
	}

	texture, err := RenderInstance.CreateTextureFromSurface(surface)
	if err != nil {
		surface.Free()
		return nil, true
	}

	image := &Image{texture}
	surface.Free()
	return image, false
}

func (i *Image) Destroy() {
	i.Texture.Destroy()
}
