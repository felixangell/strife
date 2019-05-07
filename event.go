package strife

import (
	"github.com/veandco/go-sdl2/sdl"
)

// StrifeEvent wraps an SDL event
type StrifeEvent interface {
	sdl.Event
	Trigger()
}

// HandleEvent will trigger the given StrifeEvent
func HandleEvent(event StrifeEvent) {
	event.Trigger()
}

// BASIC EVENT

// BaseEvent is a simple helper structure
// for events to satisfy sdl2
type BaseEvent struct{}

// GetTimestamp returns the timestamp of the event
// FIXME
func (b *BaseEvent) GetTimestamp() uint32 {
	// FIXME
	return 0
}

// GetType returns the type of this event
// FIXME
func (b *BaseEvent) GetType() uint32 {
	// FIXME
	return 0
}

// Trigger is invoked when the event is triggered
func (b *BaseEvent) Trigger() {}

// MOUSE EVENTS

// MouseWheelEvent represents a mouse scroll
// wheel event
type MouseWheelEvent struct {
	BaseEvent
	X, Y int
}

// MouseMoveEvent represents a mouse movement
// event
type MouseMoveEvent struct {
	BaseEvent
	X, Y int
}

// KEYBOARD

// KeyUpEvent is invoked when the key is _released_
type KeyUpEvent struct {
	BaseEvent
	KeyCode int
}

// KeyDownEvent is invoked when a key is pressed and
// held down
type KeyDownEvent struct {
	BaseEvent
	KeyCode int
}

// WINDOW CLOSE

// CloseEvent is invoked when the window
// requests to close.
type CloseEvent struct {
	BaseEvent
}

// WINDOW VISIBILITY

// Visibility of the window, e.g.
// hidden, shown.
type Visibility int

// The types of visibilities available
// for the window.
const (
	Shown Visibility = iota
	Hidden
	Exposed
)

// WindowVisibilityEvent
type WindowVisibilityEvent struct {
	BaseEvent
	Visibility
}

// WINDOW RESIZE

// WindowResizeEvent is invoked when the window
// is resized. contains the new width and height
type WindowResizeEvent struct {
	BaseEvent
	Width, Height int
}

// WINDOW MOVE

// WindowMoveEvent is invoked when the window is moved
// contains the new x, y position of the window.
type WindowMoveEvent struct {
	BaseEvent
	X, Y int
}

// WINDOW FOCUS

// Focus represents the state of focus for the window
type Focus int

// Currently two types: focus gained, and focus lost.
// i.e. the window is clicked onto or the window is clicked off of.
const (
	FocusGained Focus = iota
	FocusLost
)

// WindowFocusEvent
type WindowFocusEvent struct {
	BaseEvent
	Focus
}
