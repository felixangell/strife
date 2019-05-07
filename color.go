package strife

import (
	"github.com/veandco/go-sdl2/sdl"
)

// List of preset colours
var (
	White = RGB(255, 255, 255)
	Red   = RGB(255, 0, 0)
	Green = RGB(0, 255, 0)
	Blue  = RGB(0, 0, 255)
	Black = RGB(0, 0, 0)
)

// Color is an RGBA colour
type Color struct {
	R, G, B, A uint8
}

// Equals will compare this colour with another colour o
// including alpha channels.
func (c *Color) Equals(o *Color) bool {
	return c.R == o.R && c.G == o.G && c.B == o.B && c.A == o.A
}

// ToSDLColor Converts to an SDL colour object
func (c Color) ToSDLColor() sdl.Color {
	return sdl.Color{c.R, c.G, c.B, c.A}
}

// AsHex will return the Color as a uint32/hex
func (c Color) AsHex() uint32 {
	hex := uint32(c.R)
	hex = (hex << 8) + uint32(c.G)
	hex = (hex << 8) + uint32(c.B)
	return hex
}

var colorCache = map[uint32]*Color{}

// HexRGB will convert the given hex uint32
// value to a Color
// TODO: alpha channel >> 24.
func HexRGB(col uint32) *Color {
	r := uint8((col >> 16) & 0xff)
	g := uint8((col >> 8) & 0xff)
	b := uint8((col) & 0xff)

	if col, ok := colorCache[col]; ok {
		return col
	}
	colour := &Color{r, g, b, 255}
	colorCache[col] = colour
	return colour
}

// RGBA will create a colour from the given r, g, b, a
func RGBA(r, g, b, a int) *Color {
	result := uint32(((r & 0xff) << 16) + ((g & 0xff) << 8) + (b & 0xff))
	return HexRGB(result)
}

// RGB will create a colour from the given RGB, alpha
// is set to full.
func RGB(r, g, b int) *Color {
	return RGBA(r, g, b, 255)
}
