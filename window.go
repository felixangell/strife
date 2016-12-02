package strife

import (
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

func CreateRenderWindow(w, h int) (*RenderWindow, bool) {
	sdl.Init(sdl.INIT_VIDEO)

	windowHandle, err := sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, true
	}

	window := &RenderWindow{
		Window:        windowHandle,
		renderContext: nil,
	}

	renderer, failed := CreateRenderer(window)
	if failed {
		return nil, true
	}
	window.renderContext = renderer
	RenderInstance = renderer

	runtime.LockOSThread()
	return window, false
}
