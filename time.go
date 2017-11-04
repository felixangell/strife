package strife

import "github.com/veandco/go-sdl2/sdl"

func CurrentTimeMillis() int64 {
	return int64(sdl.GetTicks())
}