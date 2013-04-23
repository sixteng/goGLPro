// +build linux,x11

package window
/*
#cgo pkg-config: x11
#include<X11/X.h>
#include<X11/Xlib.h>

int GetEventType(XEvent * event) {
	return event->type;
}

*/
import "C"
import "unsafe"
import "fmt"

type XWindow struct {
	width uint
	height uint
	Callback  func(Window, uint32, uintptr, uintptr) bool
	dpy * C.Display
}

func CreateWindow() Window {
	return &XWindow{}
}

func (window *XWindow) Peek() bool {
	var xev C.XEvent

	events := C.XCheckWindowEvent(window.dpy)
	fmt.Println(events)
	if(events != 0) {
		fmt.Println("Get Event")
		C.XNextEvent(window.dpy, &xev)
		if (C.GetEventType(&xev) == C.Expose){
			return false
		}

		window.Callback(window, uint32(C.GetEventType(&xev)), 0, 0)

		return true
	}

	return false


	/*
	ev, xerr := window.x.PollForEvent()
	if(ev != nil) {
		//window.Callback(window, ev.)
		return true;
	}

	if(xerr != nil) {
		return true;
	}
	*/
	return false
}

func (window *XWindow) SetCallback(callback func(Window, uint32, uintptr, uintptr) bool) {
	window.Callback = callback
}

func (window *XWindow) GetHDC() unsafe.Pointer {
	return nil
}

func (window *XWindow) SwapBuffers() {

}

func (window *XWindow) Create(width uint, height uint, title string) bool {
	dpy := C.XOpenDisplay(nil)
	scr := C.XDefaultScreen(dpy)
	root := C.XDefaultRootWindow(dpy)

	win := C.XCreateSimpleWindow(dpy, root, 0, 0, C.uint(width), C.uint(height), 0, C.XBlackPixel(dpy, scr), C.XBlackPixel(dpy, scr))

	C.XStoreName(dpy, win, C.CString(title));
	C.XMapWindow(dpy, win);

	window.dpy = dpy

	return true
}
