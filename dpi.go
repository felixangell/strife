package strife

import (
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

// returns the dpi, and the default dpi
func GetDisplayDPI(displayIndex int) (float32, float32) {
	_, hdpi, _, err := sdl.GetDisplayDPI(displayIndex)
	if err != nil {
		return 0, defaultDpi()
	}
	return hdpi, defaultDpi()
}
