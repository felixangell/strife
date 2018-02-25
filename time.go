package strife

import "github.com/veandco/go-sdl2/sdl"

// dont use this! go has some nice
// utilities instead.
func CurrentTimeMillis() int64 {
	return int64(sdl.GetTicks())
}
