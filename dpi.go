package strife

import (
	"fmt"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

func defaultDpi() float32 {
	if runtime.GOOS == "windows" {
		return 96.0
	} else if runtime.GOOS == "darwin" {
		return 72.0
	}
	return 72.0 // hm!?
}

// GetDisplayDPI returns the dpi and default dpi of the
// given display
func GetDisplayDPI(displayIndex int) (dpi float32, def float32) {
	ddpi, hdpi, vdpi, err := sdl.GetDisplayDPI(displayIndex)
	if err != nil {
		return 0, defaultDpi()
	}
	fmt.Println("GetDisplayDPI", ddpi, hdpi, vdpi)
	return hdpi, defaultDpi()
}
