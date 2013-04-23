// +build linux

package window

import "unsafe"
//import "fmt"

/*
#cgo pkg-config: gtk+-3.0
#include <gtk/gtk.h>
#include <gdk/gdkx.h>

extern void GtkSetEvent(gpointer * instance, char * event, void * data);
*/
import "C"

//export goHandleEvent
func goHandleEvent(userdata unsafe.Pointer) {
	if userdata != nil {
		(*(*func())(userdata))()
	}
}

type XWinData struct {
	dpy * C.Display
	win C.Window
}

type GTKWindow struct {
	width uint
	height uint
	Callback  func(Window, uint32, uintptr, uintptr) bool
	destroy bool
	render bool
}

type GtkCallbackHandler func()


func CreateWindow() Window {
	return &GTKWindow{destroy : false, render : false}
}

func (window *GTKWindow) Peek() bool {
	if C.gtk_events_pending() == 1 {
		C.gtk_main_iteration()

		return true
	}

	if window.destroy {
		window.Callback(window, 0x0012, 0, 0)
		return true
	}

	if window.render {
		window.render = false
		//window.Callback(window, 0x0012, 0, 0)
		return false
	}

	return true
}

func (window *GTKWindow) SetCallback(callback func(Window, uint32, uintptr, uintptr) bool) {
	window.Callback = callback
}

func (window *GTKWindow) GetHDC() unsafe.Pointer {
	data := XWinData{
		dpy : C.gdk_x11_get_default_xdisplay(),
		win : C.gdk_x11_get_default_root_xwindow(),
	}

	return unsafe.Pointer(&data)
}

func (window *GTKWindow) SwapBuffers() {

}

func (window *GTKWindow) Create(width uint, height uint, title string) bool {

	C.gtk_init (nil, nil);

	gtk_window := C.gtk_window_new (C.GTK_WINDOW_TOPLEVEL);
	C.gtk_window_set_title ( (*C.GtkWindow)(unsafe.Pointer(gtk_window)), (*C.gchar)(C.CString(title)));

	f := func() {
		window.destroy = true
	}

	r := func() {
		window.render = true
	}

	//g_signal_connect (G_OBJECT (canvas), "expose-event", G_CALLBACK (expose_cb), NULL);
	C.GtkSetEvent((*C.gpointer)(unsafe.Pointer(gtk_window)), C.CString("destroy"), unsafe.Pointer(&f))
	C.GtkSetEvent((*C.gpointer)(unsafe.Pointer(gtk_window)), C.CString("expose-event"), unsafe.Pointer(&r))
	C.gtk_widget_show (gtk_window);

	return true
}

