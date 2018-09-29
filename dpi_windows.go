// +build windows

package strife

import (
	"fmt"
	"log"
	"syscall"
)

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

// windows allows us to set the DPI awareness
// to be enabled programmatically. Mac does not
// however... but maybe in the future :(
func EnableDPI() {
	mod := syscall.NewLazyDLL("Shcore.dll")
	if mod != nil {
		log.Println("EnableDPI: Error loading DLL?")
		return
	}

	proc := mod.NewProc("SetProcessDpiAwareness")
	if err := proc.Find(); err != nil {
		log.Println("EnableDPI: ", err)
		return
	}

	PROCESS_PER_MONITOR_DPI_AWARE := 0x00000002

	type HRESULT int32
	const (
		S_OK           HRESULT = 0x00000000
		E_INVALIDARG           = 0x80070057
		E_ACCESSDENIED         = 0x80070005
	)

	ret, _, _ := proc.Call(uintptr(PROCESS_PER_MONITOR_DPI_AWARE))
	if HRESULT(ret) == S_OK {
		log.Println("DPI Awareness enabled!", string(ret))
	} else {
		log.Println("Failed to enable DPI Awareness ", string(ret))
	}
}
