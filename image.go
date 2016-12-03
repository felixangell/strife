package strife

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type Image struct {
	*sdl.Texture
	*sdl.Surface
}

func LoadImage(path string) (*Image, error) {
	surface, err := img.Load(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to load image '%s'\n", path)
	}

	texture, err := RenderInstance.CreateTextureFromSurface(surface)
	if err != nil {
		surface.Free()
		return nil, fmt.Errorf("Failed to load '%s' into memory\n", path)
	}

	image := &Image{
		texture,
		surface,
	}
	return image, nil
}

func (i *Image) GetSurface() *sdl.Surface {
	return i.Surface
}

func (i *Image) Destroy() {
	i.Texture.Destroy()
	i.Surface.Free()
}
