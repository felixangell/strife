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

// returns the dpi, and the default dpi
func GetDisplayDPI(displayIndex int) (dpi float32, def float32) {
	ddpi, hdpi, vdpi, err := sdl.GetDisplayDPI(displayIndex)
	if err != nil {
		panic(err)
		return 0, defaultDpi()
	}

	fmt.Println("dpi stuff", ddpi, hdpi, vdpi)
	return hdpi, defaultDpi()
}
