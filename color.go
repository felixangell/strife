package strife

import "github.com/veandco/go-sdl2/sdl"

var (
	Black *Color = RGB(0, 0, 0)
	White        = RGB(255, 255, 255)
	Red          = RGB(255, 0, 0)
	Green        = RGB(0, 255, 0)
	Blue         = RGB(0, 0, 255)
)

type Color struct {
	R, G, B, A uint8
}

func (c Color) ToSDLColor() sdl.Color {
	return sdl.Color{c.R, c.G, c.B, c.A}
}

func RGBA(r, g, b, a uint8) *Color {
	return &Color{r, g, b, a}
}

func RGB(r, g, b uint8) *Color {
	return RGBA(r, g, b, 255)
}
