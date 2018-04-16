package strife

import (
	"fmt"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

// The current render instance, in an ideal
// world this will only be set once. This is exists
// because SDL wants to have the render instance for
// loading images, fonts, etc.
var RenderInstance *Renderer

func init() {
	runtime.LockOSThread()
}

type KeyboardHandler struct {
	keys map[int]bool
	buff []int
}

var keyboardInstance = &KeyboardHandler{
	keys: map[int]bool{},
	buff: []int{},
}

func PollKeys() bool {
	return len(keyboardInstance.buff) > 0
}

func PopKey() int {
	keyBuffSize := len(keyboardInstance.buff)

	// get the key
	keyPressed := keyboardInstance.buff[keyBuffSize-1]

	// apply pop
	keyboardInstance.buff = keyboardInstance.buff[:keyBuffSize-1]

	return keyPressed
}

func KeyPressed(keyCode int) bool {
	if val, ok := keyboardInstance.keys[keyCode]; ok {
		return val
	}
	return false
}

type RenderWindow struct {
	*sdl.Window
	config         *RenderConfig
	renderContext  *Renderer
	handler        func(StrifeEvent)
	w, h           int
	closeRequested bool
	flags          uint32
}

func (w *RenderWindow) SetIconImage(img *Image) {
	w.SetIcon(img.Surface)
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

		case *sdl.KeyboardEvent:
			keyCode := int(evt.Keysym.Sym)
			if evt.Type == sdl.KEYUP {
				w.handler(&KeyUpEvent{BaseEvent{}, keyCode})
				keyboardInstance.keys[keyCode] = false
			} else if evt.Type == sdl.KEYDOWN {
				w.handler(&KeyDownEvent{BaseEvent{}, keyCode})
				keyboardInstance.keys[keyCode] = true

				// append the key press into a key
				// buffer which can be processed.
				keyboardInstance.buff = append(keyboardInstance.buff, keyCode)
			}

		case *sdl.MouseWheelEvent:
			w.handler(&MouseWheelEvent{BaseEvent{}, int(evt.X), int(evt.Y)})

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
	windowHandle, err := sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(w.w), int32(w.h), w.flags)
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
