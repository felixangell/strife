package strife

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type Image struct {
	*sdl.Texture
	*sdl.Surface
	Width, Height int
}

// LoadImage will load the image at the given path. It will
// return the loaded image, and any errors encountered.
func LoadImage(path string) (*Image, error) {
	surface, err := img.Load(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to load image '%s'\n", path)
	}

	if RenderInstance == nil {
		return nil, fmt.Errorf("Render context has not been initialized yet.")
	}

	texture, err := RenderInstance.CreateTextureFromSurface(surface)
	if err != nil {
		surface.Free()
		return nil, fmt.Errorf("Failed to load '%s' into memory\n", path)
	}

	image := &Image{
		texture,
		surface,
		int(surface.W),
		int(surface.H),
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
