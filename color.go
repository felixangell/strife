package strife

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

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
	return int(((c.R & 0xff) << 16) + ((c.G & 0xff) << 8) + (c.B & 0xff))
}

var colorCache = map[int32]*Color{}

// TODO: alpha channel >> 24.
func HexRGB(col int32) *Color {
	r := uint8(col & 0xff0000 >> 16)
	g := uint8(col & 0xff00 >> 8)
	b := uint8(col & 0xff)

	if col, ok := colorCache[col]; ok {
		return col
	}
	colour := &Color{r, g, b, 255}
	colorCache[col] = colour
	log.Println("cached color ", fmt.Sprintf("0x%x", colour.AsHex()), " caches: ", len(colorCache))
	return colour
}

func RGBA(r, g, b, a uint8) *Color {
	res := int32(((r & 0xff) << 16) | ((g & 0xff) << 8) | (b & 0xff))
	return HexRGB(res)
}

func RGB(r, g, b uint8) *Color {
	return RGBA(r, g, b, 255)
}
