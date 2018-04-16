package strife

import (
	"github.com/veandco/go-sdl2/sdl"
)

type StrifeEvent interface {
	sdl.Event
	Trigger()
}

func HandleEvent(event StrifeEvent) {
	event.Trigger()
}

// BASIC EVENT

type BaseEvent struct{}

func (b *BaseEvent) Trigger() {}

// MOUSE EVENTS

type MouseWheelEvent struct {
	BaseEvent
	X, Y int
}

// KEYBOARD

type KeyUpEvent struct {
	BaseEvent
	KeyCode int
}

type KeyDownEvent struct {
	BaseEvent
	KeyCode int
}

// WINDOW CLOSE

type CloseEvent struct {
	BaseEvent
}

// WINDOW VISIBILITY

type Visibility int

const (
	Shown Visibility = iota
	Hidden
	Exposed
)

type WindowVisibilityEvent struct {
	BaseEvent
	Visibility
}

// WINDOW RESIZE

type WindowResizeEvent struct {
	BaseEvent
	Width, Height int
}

// WINDOW MOVE

type WindowMoveEvent struct {
	BaseEvent
	X, Y int
}

// WINDOW FOCUS

type Focus int

const (
	FocusGained Focus = iota
	FocusLost
)

type WindowFocusEvent struct {
	BaseEvent
	Focus
}
