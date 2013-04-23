// +build !linux

package opengl

/*
#cgo windows CFLAGS: -I../../../../pkg/glew/include
#cgo windows LDFLAGS: -lopengl32 -lglew32 -L../../../../pkg/glew/lib
#
#include <GL/glew.h>
#include <GL/wglew.h>
#include <GL/gl.h>

HGLRC openglCreateContextAttribsARB(HDC hDC, HGLRC hShareContext) {
	int attributes[] = {
		WGL_CONTEXT_MAJOR_VERSION_ARB, 3,
		WGL_CONTEXT_MINOR_VERSION_ARB, 2,
		WGL_CONTEXT_FLAGS_ARB, WGL_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB,
		0
	};

	return wglCreateContextAttribsARB(hDC, hShareContext, attributes);
}

*/
import "C"
import "../../../window"
import "fmt"

func CreateContext(win window.Window) C.HGLRC {
	hdc := win.GetHDC()

	if hdc == nil {
		//TODO: Handle error somehow ???
		return nil
	}

	tmpContext := C.wglCreateContext(*(*C.HDC)(hdc))
	C.wglMakeCurrent(*(*C.HDC)(hdc), tmpContext)
	C.glewInit()

	var context C.HGLRC
	if C.wglewIsSupported(C.CString("WGL_ARB_create_context")) == 1 {
		//TODO: Figure out how to send attr to c func so it works!!!
		context = C.openglCreateContextAttribsARB(*(*C.HDC)(hdc), nil)

		C.wglMakeCurrent(nil, nil)
		C.wglDeleteContext(tmpContext)
		C.wglMakeCurrent(*(*C.HDC)(hdc), context)
	}
	return tmpContext
}

func Test(win window.Window) {
	C.glClearColor(1.0, 0.0, 0.0, 0.0)
	C.glClear(C.GL_COLOR_BUFFER_BIT | C.GL_DEPTH_BUFFER_BIT | C.GL_STENCIL_BUFFER_BIT)
	win.SwapBuffers()
}

//http://stackoverflow.com/questions/6005076/building-glew-on-windows-with-mingw
