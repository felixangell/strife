package strife

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// Style to render
// primitive shapes (Rectangles, Circles, etc).
// Fill meaning fill the shape with colour, and
// Line meaning to draw the outline of the shape.
type Style uint

// Types of render styles, Line for stroke
// e.g. outline of a shape, vs fill which will
// fill a rectangle.
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

// GoGoStrifeFast is the default configuration
// with all pretty settings turned off.
// vertical synchronisation is enabled.
func GoGoStrifeFast() *RenderConfig {
	sdl.Init(sdl.INIT_VIDEO)

	return &RenderConfig{
		Alias:        false,
		Accelerated:  true,
		VerticalSync: true,
	}
}

// DefaultConfig with all of the rendering options
// enabled.
func DefaultConfig() *RenderConfig {
	sdl.Init(sdl.INIT_VIDEO)
	return &RenderConfig{
		Alias:        true,
		Accelerated:  true,
		VerticalSync: true,
	}
}

// Renderer contains the
// renderers current configuration, as well
// as the colour state and font state and a wrapper over
// the SDL renderer.
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

// GetSize returns the size of the renderer. on error
// it will return -1, -1
func (r *Renderer) GetSize() (int, int) {
	w, h, err := r.Renderer.GetOutputSize()
	if err != nil {
		return -1, -1
	}
	return int(w), int(h)
}

// Display the renderer to the window
func (r *Renderer) Display() {
	r.Renderer.Present()
}

// SetColor sets the current colour state
func (r *Renderer) SetColor(color *Color) {
	r.color = color
	r.SetDrawColor(color.R, color.G, color.B, color.A)
}

// Rect will draw a rectangle at the given x, y co-ordinates
// of the specified size. It takes the mode to render the
// rectangle as: fill or line.
func (r *Renderer) Rect(x, y, w, h int, mode Style) {
	if mode == Line {
		r.DrawRect(&sdl.Rect{int32(x), int32(y), int32(w), int32(h)})
	} else {
		r.FillRect(&sdl.Rect{int32(x), int32(y), int32(w), int32(h)})
	}
}

// GetFont returns the current font that was last set
func (r *Renderer) GetFont() *Font {
	return r.font
}

// SetFont will set the font state
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

func (r *Renderer) renderRune(color *Color, char rune) (*sdl.Texture, []int32) {
	message := string(char)

	var surface *sdl.Surface
	var err error

	if r.Alias {
		surface, err = r.font.RenderUTF8Blended(message, color.ToSDLColor())
	} else {
		surface, err = r.font.RenderUTF8Solid(message, color.ToSDLColor())
	}

	defer surface.Free()

	if err != nil {
		panic(err)
	}

	texture, err := r.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}

	// TODO we could store how many times
	// these get used and then run a thread to free
	// some of the unused textures every now and then?
	// or an LRU cache or something?
	// or memory pool allocation?
	allocs++

	// log.Println("allocs ", allocs)

	return texture, []int32{surface.W, surface.H}
}

func (r *Renderer) GetStringDimension(message string) (int, int) {
	w, h := r.GetSize()
	return r.Text(message, -w*2, -h*2)
}

// Text renders the given text to the given x, y coordinates.
// Note that the text is cached, i.e. each glyph rendered will be cached
// and re-used. UncachedText is the alternative, though it's slower.
func (r *Renderer) Text(message string, x, y int) (int, int) {
	if r.font == nil {
		panic("Attempted to render '" + message + "' but no font is set!")
	}

	var width, height int32

	col := r.color.AsHex()
	for _, char := range message {
		encoding := encode(col, plain, char)

		glyph, ok := r.font.hasGlyph(encoding)
		if !ok {
			texture, dim := r.renderRune(r.color, char)
			glyph = r.font.cache(encoding, texture, dim)
		}

		dim := glyph.dim

		r.Renderer.Copy(glyph.tex, nil, &sdl.Rect{int32(x) + width, int32(y), dim[0], dim[1]})
		width += dim[0]
		height = maxInt32(height, dim[1])
	}

	return int(width), int(height)
}

// UncachedText is the same as Text, it draws the given
// string to the given x,y
// note that it doesn't cache the glyphs, however, so
// should not be used for lots of text rendering!
func (r *Renderer) UncachedText(message string, x, y int) (int, int) {
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

	// find a default font
	fontFolder := func() string {
		switch runtime.GOOS {
		case "windows":
			return filepath.Join(os.Getenv("WINDIR"), "fonts")
		case "darwin":
			return "/Library/Fonts/"
		case "linux":
			// FIXME
			return "/usr/share/fonts/"
		default:
			log.Fatal("no font folder set for OS ", runtime.GOOS)
			panic("oh boy!")
		}
	}()

	fontChoices := map[string]bool{}
	filepath.Walk(fontFolder, func(path string, r os.FileInfo, err error) error {
		fontName := filepath.Base(path)
		fontChoices[fontName] = true
		return nil
	})

	chosenFont := func() string {
		defaultFontChoices := []string{
			// todo add more fonts.
			// maybe depending on the platform
			// we can sort them manually
			// e.g. georgia first on mac, calibiri first on windows
			// for faster lookup?
			"Verdana.ttf", "calibri.ttf", "verdana.ttf",
			"Ubuntu.ttf", "DejaVuSans.ttf",
		}
		for _, font := range defaultFontChoices {
			if _, exists := fontChoices[font]; exists {
				return font
			}
		}
		return ""
	}()

	fontPath := filepath.Join(fontFolder, chosenFont)
	log.Println("Loading font ", fontPath)

	renderInst.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	renderer := &Renderer{
		RenderConfig: *config,
		Renderer:     renderInst,
		color:        RGB(255, 255, 255),
	}

	// load a default font to render with.
	defaultFont, err := LoadFont(fontPath, 24)
	if err == nil {
		renderer.SetFont(defaultFont)
	} else {
		log.Println(err.Error(), "' try setting a font yourself with strife.LoadFont")
	}

	return renderer, nil
}
