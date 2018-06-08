package strife

import "github.com/veandco/go-sdl2/sdl"

// MOUSE INPUT

type MouseButtonState int

const (
	NoMouseButtonsDown MouseButtonState = iota
	LeftMouseButton
	RightMouseButton
	ScrollWheel
)

type MouseHandler struct {
	X, Y        int
	ButtonState MouseButtonState
}

var mouseInstance = &MouseHandler{}

func MouseButtonsState() MouseButtonState {
	return mouseInstance.ButtonState
}

func MouseCoords() (int, int) {
	return mouseInstance.X, mouseInstance.Y
}

// KEY INPUT

type KeyboardHandler struct {
	keys map[int]bool
	buff []int
}

func KeyState() []uint8 {
	return sdl.GetKeyboardState()
}

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
