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
	config         *RenderConfig
	renderContext  *Renderer
	handler        func(StrifeEvent)
	w, h           int
	closeRequested bool
	flags          uint32
}

func (w *RenderWindow) CloseRequested() bool {
	return w.closeRequested
}

func (w *RenderWindow) HandleEvents(handler func(StrifeEvent)) {
	w.handler = handler
}

// CloseRequested will poll for any events. All events
// are unhandled, _except_ for the Quit Event, which will
// destroy the render context, render window, and cause
// this function to return true.
func (w *RenderWindow) PollEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch evt := event.(type) {
		case *sdl.QuitEvent:
			w.handler(&CloseEvent{BaseEvent{}})
		case *sdl.WindowEvent:
			switch evt.Event {

			// events that affect visibility
			case sdl.WINDOWEVENT_HIDDEN:
				w.handler(&WindowVisibilityEvent{BaseEvent{}, Hidden})
			case sdl.WINDOWEVENT_SHOWN:
				w.handler(&WindowVisibilityEvent{BaseEvent{}, Shown})
			case sdl.WINDOWEVENT_EXPOSED:
				w.handler(&WindowVisibilityEvent{BaseEvent{}, Exposed})

			// size/position stuff
			case sdl.WINDOWEVENT_SIZE_CHANGED:
				// should this be handled as its own event
				// or as a resized event?
				fallthrough
			case sdl.WINDOWEVENT_RESIZED:
				w.handler(&WindowResizeEvent{BaseEvent{}, int(evt.Data1), int(evt.Data2)})
			case sdl.WINDOWEVENT_MOVED:
				w.handler(&WindowMoveEvent{BaseEvent{}, int(evt.Data1), int(evt.Data2)})

			// TODO: ENTER/LEAVE ... CLOSE?

			// events that are to do with focus!
			case sdl.WINDOWEVENT_FOCUS_GAINED:
				w.handler(&WindowFocusEvent{BaseEvent{}, FocusLost})
			case sdl.WINDOWEVENT_FOCUS_LOST:
				w.handler(&WindowFocusEvent{BaseEvent{}, FocusGained})
			}
		}
	}
}

func (w *RenderWindow) GetRenderContext() *Renderer {
	return w.renderContext
}

func (w *RenderWindow) Close() {
	w.closeRequested = true
	w.renderContext.Destroy()
	w.Destroy()
}

func (w *RenderWindow) SetResizable(resizable bool) {
	w.flags |= sdl.WINDOW_RESIZABLE
}

func (w *RenderWindow) Create() error {
	windowHandle, err := sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, w.w, w.h, w.flags)
	if err != nil {
		return fmt.Errorf("Failed to create window\n")
	}
	w.Window = windowHandle

	renderer, err := CreateRenderer(w, w.config)
	if err != nil {
		return err
	}
	w.renderContext = renderer
	RenderInstance = renderer

	runtime.LockOSThread()
	return nil
}

func (w *RenderWindow) GetSize() (int, int) {
	ww, hh := w.Window.GetSize()
	return int(ww), int(hh)
}

// CreateRenderWindow will create a render window of the given
// width and height, with the specified configuration. You can
// specify a default configuration with `strife.DefaultConfig()`.
// Note that this will spawn the window at the centre of the main
// display.
func SetupRenderWindow(w, h int, config *RenderConfig) *RenderWindow {
	sdl.Init(sdl.INIT_VIDEO)

	window := &RenderWindow{
		config: config,
		w:      w,
		h:      h,
	}
	window.handler = func(evt StrifeEvent) {
		switch evt.(type) {
		case *CloseEvent:
			window.Close()
		}
	}
	return window
}
