package strife

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"runtime"
)

var RenderInstance *Renderer

type RenderWindow struct {
	*sdl.Window
	renderContext *Renderer
}

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
