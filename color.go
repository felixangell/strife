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

func (c *Color) Equals(o *Color) bool {
	return c.R == o.R && c.G == o.G && c.B == o.B && c.A == o.A
}

func (c Color) ToSDLColor() sdl.Color {
	return sdl.Color{c.R, c.G, c.B, c.A}
}

func (c Color) AsHex() int {
	res := int(((c.R & 0xff) << 16) | ((c.G & 0xff) << 8) | (c.B & 0xff))
	return res
}

// TODO: alpha channel >> 24.
func HexRGB(col int32) *Color {
	a := uint8(255)
	r := uint8(col & 0xff0000 >> 16)
	g := uint8(col & 0xff00 >> 8)
	b := uint8(col & 0xff)
	return RGBA(r, g, b, a)
}

func RGBA(r, g, b, a uint8) *Color {
	return &Color{r, g, b, a}
}

func RGB(r, g, b uint8) *Color {
	return RGBA(r, g, b, 255)
}
