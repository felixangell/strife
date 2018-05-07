package strife

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

// The style to render
// primitive shapes (Rectangles, Circles, etc).
// Fill meaning fill the shape with colour, and
// Line meaning to draw the outline of the shape.
type Style uint

const (
	Line Style = iota
	Fill
)

// RenderConfig is the configuration settings
// for the renderer.
// Alias => Controls if the fonts are smoothed;
// Accelerated => The renderer is hardware accelerated if true; and
// VerticalSync => Will synchronize the FPS with the monitors refresh rate
type RenderConfig struct {
	Alias        bool
	Accelerated  bool
	VerticalSync bool
}

// DefaultConfig with all of the rendering options
// enabled.
func DefaultConfig() *RenderConfig {
	return &RenderConfig{
		Alias:        true,
		Accelerated:  true,
		VerticalSync: true,
	}
}

type Renderer struct {
	RenderConfig
	*sdl.Renderer

	color *Color
	font  *Font
}

// Clear will clear the screen to black. By default
// it will immediately set the colour state to render
// things as white.
func (r *Renderer) Clear() {
	r.SetColor(Black)
	w, h, err := r.Renderer.GetOutputSize()
	if err != nil {
		panic(err)
	}
	r.Rect(0, 0, int(w), int(h), Fill)

	r.SetColor(White)
}

func (r *Renderer) Display() {
	r.Renderer.Present()
}

func (r *Renderer) SetColor(col *Color) {
	r.color = col
}

// Rect will draw a rectangle at the given x, y co-ordinates
// of the specified size. It takes the mode to render the
// rectangle as: fill or line.
func (r *Renderer) Rect(x, y, w, h int, mode Style) {
	color := r.color
	r.SetDrawColor(color.R, color.G, color.B, color.A)

	if mode == Line {
		r.DrawRect(&sdl.Rect{int32(x), int32(y), int32(w), int32(h)})
	} else {
		r.FillRect(&sdl.Rect{int32(x), int32(y), int32(w), int32(h)})
	}
}

func (r *Renderer) SetFont(font *Font) {
	r.font = font
}

func maxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

var allocs int

func (r *Renderer) renderRune(char rune) (*sdl.Texture, []int32) {
	message := string(char)

	var surface *sdl.Surface
	var err error
	if r.Alias {
		surface, err = r.font.RenderUTF8Blended(message, r.color.ToSDLColor())
	} else {
		surface, err = r.font.RenderUTF8Solid(message, r.color.ToSDLColor())
	}
	defer surface.Free()
	if err != nil {
		panic(err)
	}

	texture, err := r.Renderer.CreateTextureFromSurface(surface)

	// TODO we could store how many times
	// these get used and then run a thread to free
	// some of the unused textures every now and then?
	allocs++

	log.Println("allocs:", allocs)
	if err != nil {
		panic(err)
	}
	return texture, []int32{surface.W, surface.H}
}

func (r *Renderer) String(message string, x, y int) (int, int) {
	if r.font == nil {
		panic("Attempted to render '" + message + "' but no font is set!")
	}

	var width, height int32
	for _, char := range message {
		glyph, ok := r.font.CharCache[char]

		// no glyph has been cached
		// so cache one
		if !ok {
			texture, dim := r.renderRune(char)
			glyph = NewGlyph(dim[0], dim[1], r.color, texture)

			r.font.CharCache[char] = glyph
		}

		// we don't have the correct color
		glyphTexture, ok := glyph.texs[r.color.AsHex()]
		if !ok {
			glyphTexture, _ = r.renderRune(char)
			glyph.texs[r.color.AsHex()] = glyphTexture

			// store that boy
			r.font.CharCache[char] = glyph
		}

		r.Renderer.Copy(glyphTexture, nil, &sdl.Rect{int32(x) + width, int32(y), glyph.w, glyph.h})
		width += glyph.w
		height = maxInt32(height, glyph.h)
	}
	return int(width), int(height)
}

func (r *Renderer) UncachedString(message string, x, y int) (int, int) {
	if r.font == nil {
		panic("Attempted to render '" + message + "' but no font is set!")
	}

	var surface *sdl.Surface
	var err error
	if r.Alias {
		surface, err = r.font.RenderUTF8Blended(message, r.color.ToSDLColor())
	} else {
		surface, err = r.font.RenderUTF8Solid(message, r.color.ToSDLColor())
	}
	defer surface.Free()
	if err != nil {
		panic(err)
	}

	texture, err := r.Renderer.CreateTextureFromSurface(surface)
	defer texture.Destroy()
	if err != nil {
		panic(err)
	}

	r.Renderer.Copy(texture, nil, &sdl.Rect{int32(x), int32(y), surface.W, surface.H})
	return int(surface.W), int(surface.H)
}

// SubImage will render a sub-section of the given image. tx, ty are
// the x, y pixel coordinates of the sub-image in the image to render. tw, th
// are the size of the sub-image.
func (r *Renderer) SubImage(image *Image, x, y int, tx, ty, tw, th int) {
	r.SubImageScale(image, x, y, tx, ty, tw, th, tw, th)
}

// SubImageScale will render a sub-section of the given image scaled
// to the given width and height. See the documentation for SubImage.
func (r *Renderer) SubImageScale(image *Image, x, y int, tx, ty, tw, th int, sw, sh int) {
	r.Copy(image.Texture, &sdl.Rect{
		int32(tx), int32(ty), int32(tw), int32(th),
	}, &sdl.Rect{int32(x), int32(y), int32(sw), int32(sh)})
}

// Image will render the given image at the given
// x, y co-ordinates at the images full size.
func (r *Renderer) Image(image *Image, x, y int) {
	_, _, w, h, err := image.Texture.Query()
	if err != nil {
		panic(err)
	}
	r.ImageScale(image, x, y, int(w), int(h))
}

// ImageScale will render the image at the given co-ordinate
// scaled to the given size.
func (r *Renderer) ImageScale(image *Image, x, y, w, h int) {
	r.Copy(image.Texture, nil, &sdl.Rect{int32(x), int32(y), int32(w), int32(h)})
}

// CreateRenderer will create a rendering instance for the given
// window. It takes the configuration specifying if the renderer
// is software or hardware accelerated, as well as if the renderer
// should be vertically synchronized.
// It returns the renderer and any error that is encountered
// during the creation.
func CreateRenderer(parent *RenderWindow, config *RenderConfig) (*Renderer, error) {
	var mode uint32
	if config.Accelerated {
		mode |= sdl.RENDERER_ACCELERATED
	} else {
		mode |= sdl.RENDERER_SOFTWARE
	}
	if config.VerticalSync {
		mode |= sdl.RENDERER_PRESENTVSYNC
	}

	renderInst, err := sdl.CreateRenderer(parent.Window, -1, mode)
	if err != nil {
		return nil, fmt.Errorf("Failed to create render context")
	}

	renderer := &Renderer{
		RenderConfig: *config,
		Renderer:     renderInst,
		color:        RGB(255, 255, 255),
	}
	return renderer, nil
}
