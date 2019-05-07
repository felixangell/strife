package strife

import "github.com/veandco/go-sdl2/sdl"

// CurrentTimeMillis wraps over SDL_GetTicks. Use
// time.Now() instead.
func CurrentTimeMillis() int64 {
	return int64(sdl.GetTicks())
}
