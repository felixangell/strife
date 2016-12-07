package strife

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"runtime"
)

// The current render instance, in an ideal
// world this will only be set once. This is exists
// because SDL wants to have the render instance for
// loading images, fonts, etc.
var RenderInstance *Renderer

type RenderWindow struct {
	*sdl.Window
	renderContext *Renderer
}

// CloseRequested will poll for any events. All events
// are unhandled, _except_ for the Quit Event, which will
// destroy the render context, render window, and cause
// this function to return true.
func (w *RenderWindow) CloseRequested() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			{
				w.renderContext.Destroy()
				w.Destroy()
			}
			return true
		}
	}
	return false
}

func (w *RenderWindow) GetRenderContext() *Renderer {
	return w.renderContext
}

// CreateRenderWindow will create a render window of the given
// width and height, with the specified configuration. You can
// specify a default configuration with `strife.DefaultConfig()`.
// Note that this will spawn the window at the centre of the main
// display.
func CreateRenderWindow(w, h int, config *RenderConfig) (*RenderWindow, error) {
	sdl.Init(sdl.INIT_VIDEO)

	windowHandle, err := sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, fmt.Errorf("Failed to create window\n")
	}

	window := &RenderWindow{
		Window:        windowHandle,
		renderContext: nil,
	}

	renderer, err := CreateRenderer(window, config)
	if err != nil {
		return nil, err
	}
	window.renderContext = renderer
	RenderInstance = renderer

	runtime.LockOSThread()
	return window, nil
}
