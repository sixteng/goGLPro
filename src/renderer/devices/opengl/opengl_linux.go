// +build linux

package opengl

/*
#cgo pkg-config: x11
#cgo pkg-config: glew
#cgo pkg-config: gl

#include <GL/glew.h>
#include <GL/gl.h>
#include <GL/glx.h>
#include <X11/X.h>
#include <X11/Xlib.h>

XVisualInfo * Test(Display * dpy) {
	GLint att[] = { GLX_RGBA, GLX_DEPTH_SIZE, 24, GLX_DOUBLEBUFFER, None };
	return glXChooseVisual(dpy, 0, att);
} 

*/
import "C"
import "../../../window"
//import "unsafe"
import "fmt"

type XWinData struct {
	dpy * C.Display
	win C.Window
}

func CreateContext(win window.Window) {
	dpy := win.GetHDC()

	xData := (*XWinData)(dpy)

	if xData == nil {
		//TODO: Handle error somehow ???
		return
	}

	//attr := make([]C.int,)

	/*
	attr := []C.int{
		C.GLX_RGBA, 
		C.GLX_DEPTH_SIZE, 
		24, 
		C.GLX_DOUBLEBUFFER,
		C.None,
	}
*/
	
	//vi := C.glXChooseVisual(xData.dpy, 0, (*C.int)(unsafe.Pointer(&attr)));
	vi := C.Test(xData.dpy)
	tmpContext := C.glXCreateContext(xData.dpy, vi, nil, C.GL_TRUE)

	err := C.glGetError()

	if err != C.GL_NO_ERROR {
		fmt.Println("Error")
	}

	C.glXMakeCurrent(xData.dpy, (C.GLXDrawable)(xData.win), tmpContext)
	C.glewInit()
	/*
	C.glXMakeCurrent(xData.dpy)
	C.glewInit()

	var context C.HGLRC
	if C.wglewIsSupported(C.CString("WGL_ARB_create_context")) == 1 {
		//TODO: Figure out how to send attr to c func so it works!!!
		context = C.openglCreateContextAttribsARB(*(*C.HDC)(hdc), nil)

		C.wglMakeCurrent(nil, nil)
		C.wglDeleteContext(tmpContext)
		C.wglMakeCurrent(*(*C.HDC)(hdc), context)
	}
	*/
	return
}

func Test(win window.Window) {
	C.glClearColor(1.0, 0.0, 0.0, 0.0)
	C.glClear(C.GL_COLOR_BUFFER_BIT | C.GL_DEPTH_BUFFER_BIT | C.GL_STENCIL_BUFFER_BIT)
	//win.SwapBuffers()
}

//http://stackoverflow.com/questions/6005076/building-glew-on-windows-with-mingw
