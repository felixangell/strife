package strife

import "github.com/veandco/go-sdl2/sdl"

// MOUSE INPUT

// MouseButtonState contains the state of
// the mouse button, i.e. if its left/right is down
// no buttons are down, or the scroll wheel is clicked
type MouseButtonState int

// The types of mouse button state.
const (
	NoMouseButtonsDown MouseButtonState = iota
	LeftMouseButton
	RightMouseButton
	ScrollWheel
)

// MouseHandler is a wrapper to store the X, Y location
// of the mouse + the current state
type MouseHandler struct {
	X, Y        int
	ButtonState MouseButtonState
}

// FIXME
var mouseInstance = &MouseHandler{}

// MouseButtonsState returns the current state of the mouse
func MouseButtonsState() MouseButtonState {
	return mouseInstance.ButtonState
}

// MouseCoords returns the coords of the mouse
func MouseCoords() (int, int) {
	return mouseInstance.X, mouseInstance.Y
}

// KEY INPUT

// KeyboardHandler is a wrapper to handle key press
type KeyboardHandler struct {
	keys map[int]bool
	buff []int
}

// KeyState wraps over SDL GetKeyboardState
func KeyState() []uint8 {
	return sdl.GetKeyboardState()
}

// KeyPressed will query if a key is pressed in the KeyboardHandler
func KeyPressed(keyCode int) bool {
	if val, ok := keyboardInstance.keys[keyCode]; ok {
		return val
	}
	return false
}

var keyboardInstance = &KeyboardHandler{
	keys: map[int]bool{},
	buff: []int{},
}

// PollKeys returns if there are in key presses
// in the buffer
func PollKeys() bool {
	return len(keyboardInstance.buff) > 0
}

// PopKey pops a key from the key press queue.
func PopKey() int {
	keyBuffSize := len(keyboardInstance.buff)

	// get the key
	keyPressed := keyboardInstance.buff[keyBuffSize-1]

	// apply pop
	keyboardInstance.buff = keyboardInstance.buff[:keyBuffSize-1]

	return keyPressed
}
