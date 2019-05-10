package strife

import (
	"fmt"
	"log"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

// RenderInstance is the current render instance, in an ideal
// world this will only be set once. This is exists
// because SDL wants to have the render instance for
// loading images, fonts, etc.
var RenderInstance *Renderer

func init() {
	runtime.LockOSThread()
}

// RenderWindow is a window context
// with a rendererer. It contains the handlers
// for events, as well as the render configuration options.
type RenderWindow struct {
	*sdl.Window
	config         *RenderConfig
	renderContext  *Renderer
	handler        func(StrifeEvent)
	w, h           int
	closeRequested bool
	flags          uint32
}

// SetIconImage will set the window icon from the given Image
func (w *RenderWindow) SetIconImage(img *Image) {
	w.SetIcon(img.Surface)
}

// CloseRequested will return if the window has
// had a CloseRequest even triggered.
func (w *RenderWindow) CloseRequested() bool {
	return w.closeRequested
}

// HandleEvents will set the event handler predicate.
func (w *RenderWindow) HandleEvents(handler func(StrifeEvent)) {
	w.handler = handler
}

func (w *RenderWindow) handleKeyboardEvent(evt *sdl.KeyboardEvent) {
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
}

func (w *RenderWindow) handleMouseButtonEvent(evt *sdl.MouseButtonEvent) {
	if evt.Type == sdl.MOUSEBUTTONUP {
		mouseInstance.ButtonState = NoMouseButtonsDown
		return
	}

	switch evt.Button {
	case sdl.BUTTON_LEFT:
		mouseInstance.ButtonState = LeftMouseButton
	case sdl.BUTTON_MIDDLE:
		mouseInstance.ButtonState = ScrollWheel
	case sdl.BUTTON_RIGHT:
		mouseInstance.ButtonState = RightMouseButton
	}
}

func (w *RenderWindow) handleMouseMotionEvent(evt *sdl.MouseMotionEvent) {
	w.handler(&MouseMoveEvent{BaseEvent{}, int(evt.X), int(evt.Y)})
	mouseInstance.X = int(evt.X)
	mouseInstance.Y = int(evt.Y)

	switch evt.State {
	case sdl.BUTTON_LEFT:
		mouseInstance.ButtonState = LeftMouseButton
	case sdl.BUTTON_MIDDLE:
		mouseInstance.ButtonState = ScrollWheel
	case sdl.BUTTON_RIGHT:
		mouseInstance.ButtonState = RightMouseButton
	default:
		mouseInstance.ButtonState = NoMouseButtonsDown
	}
}

func (w *RenderWindow) handleWindowEvent(event *sdl.WindowEvent) {
	switch event.Event {

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
		w.handler(&WindowResizeEvent{BaseEvent{}, int(event.Data1), int(event.Data2)})
	case sdl.WINDOWEVENT_MOVED:
		w.handler(&WindowMoveEvent{BaseEvent{}, int(event.Data1), int(event.Data2)})

	// TODO: ENTER/LEAVE ... CLOSE?

	// events that are to do with focus!
	case sdl.WINDOWEVENT_FOCUS_GAINED:
		w.handler(&WindowFocusEvent{BaseEvent{}, FocusLost})
	case sdl.WINDOWEVENT_FOCUS_LOST:
		w.handler(&WindowFocusEvent{BaseEvent{}, FocusGained})
	}
}

func (w *RenderWindow) processEvent(event sdl.Event) {
	switch evt := event.(type) {
	case *sdl.QuitEvent:
		w.handler(&CloseEvent{BaseEvent{}})
	case *sdl.KeyboardEvent:
		w.handleKeyboardEvent(evt)
	case *sdl.MouseButtonEvent:
		w.handleMouseButtonEvent(evt)
	case *sdl.MouseMotionEvent:
		w.handleMouseMotionEvent(evt)
	case *sdl.MouseWheelEvent:
		w.handler(&MouseWheelEvent{BaseEvent{}, int(evt.X), int(evt.Y)})
	case *sdl.WindowEvent:
		w.handleWindowEvent(evt)
	default:
		// log.Println("unhandled event!", reflect.TypeOf(evt), evt, " ... please file an issue on GitHub!")
	}
}

// PollEvents will poll for any events. All events
// are unhandled, _except_ for the Quit Event, which will
// destroy the render context, render window, and cause
// this function to return true.
func (w *RenderWindow) PollEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		w.processEvent(event)
	}
}

// GetRenderContext returns the render context
func (w *RenderWindow) GetRenderContext() *Renderer {
	return w.renderContext
}

// Close will close the render window
// and destroy the context.
func (w *RenderWindow) Close() {
	w.closeRequested = true
	w.renderContext.Destroy()
	w.Destroy()
}

// SetResizable will add the resizable flag, must be called
// before Create()
func (w *RenderWindow) SetResizable(resizable bool) {
	w.flags |= sdl.WINDOW_RESIZABLE
}

// AllowHighDPI will allow a high DPI, must be called
// before Create()
func (w *RenderWindow) AllowHighDPI() {
	w.flags |= sdl.WINDOW_ALLOW_HIGHDPI
}

// Create window take all of the settings and create
// a window with a rendering context
func (w *RenderWindow) Create() error {
	windowHandle, err := sdl.CreateWindow("", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(w.w), int32(w.h), w.flags)
	if err != nil {
		return fmt.Errorf("Failed to create window")
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

// GetSize will return the window width and height
func (w *RenderWindow) GetSize() (int, int) {
	ww, hh := w.Window.GetSize()
	return int(ww), int(hh)
}

// SetupRenderWindow will create a render window of the given
// width and height, with the specified configuration. You can
// specify a default configuration with `strife.DefaultConfig()`.
// Note that this will spawn the window at the centre of the main
// display.
func SetupRenderWindow(w, h int, config *RenderConfig) *RenderWindow {
	log.Println("initializing window ", w, "x", h)

	EnableDPI()

	a, b := GetDisplayDPI(0)
	log.Println("dpi, default_dpi = ", a, b)

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
