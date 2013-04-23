package window

import "unsafe"

type Window interface {
	Create(width uint, height uint, title string) bool
	Peek() bool
	SetCallback(func(Window, uint32, uintptr, uintptr) bool)
	GetHDC() unsafe.Pointer
	SwapBuffers()
}
