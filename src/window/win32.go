// +build !linux

package window

/*
#cgo windows LDFLAGS: -lgdi32
#include <Windows.h>
#include <WinUser.h>
#include <WinGDI.h>

BOOL RegisterWindow(char * class, HINSTANCE hInstance, WNDPROC proc) {
	WNDCLASSEX wcex;
    wcex.cbSize = sizeof( WNDCLASSEX );
    wcex.style = CS_HREDRAW | CS_VREDRAW;
    wcex.lpfnWndProc = proc;
    wcex.cbClsExtra = 0;
    wcex.cbWndExtra = 0;
    wcex.hInstance = hInstance;
    wcex.hIcon = NULL;
    wcex.hCursor = NULL;
    wcex.hbrBackground = NULL;
    wcex.lpszMenuName = NULL;
    wcex.lpszClassName = class;
    wcex.hIconSm = NULL;

    if( !RegisterClassEx( &wcex ) ) {
		return FALSE;
	}

	return TRUE;
}

HWND CreateWindowTest(char * class, char * title, UINT width, UINT height, HINSTANCE hInstance) {
	return CreateWindow( class, title, WS_OVERLAPPEDWINDOW, 0, 0, width, height, NULL, NULL, hInstance, NULL);
}

int ChoosePX(HDC hdc, PIXELFORMATDESCRIPTOR *ppfd) {
	return ChoosePixelFormat(hdc, ppfd);
}

*/
import "C"
import "syscall"
import "fmt"
import "unsafe"

func CreateWindow() Window {
	return &Win32Window{}
}

type Win32Window struct {
	width     uint
	height    uint
	hInstance C.HMODULE
	hwnd      C.HWND
	hdc       C.HDC
	Callback  func(Window, uint32, uintptr, uintptr) bool
}

func (window *Win32Window) SetCallback(callback func(Window, uint32, uintptr, uintptr) bool) {
	window.Callback = callback
}

func (window *Win32Window) WndProc(hwnd C.HWND, message C.UINT, wParam C.WPARAM, lParam C.LPARAM) C.LRESULT {
	if message == 0x0002 {
		C.PostQuitMessage(0)
		return 0
	}

	if window.Callback != nil {
		if window.Callback(window, uint32(message), uintptr(unsafe.Pointer(&wParam)), uintptr(unsafe.Pointer(&lParam))) == true {
			return 0
		}
	}

	return C.DefWindowProc(hwnd, message, wParam, lParam)
}

func (window *Win32Window) GetHDC() unsafe.Pointer {
	if window.hdc != nil {
		return unsafe.Pointer(&window.hdc)
	}

	window.hdc = C.GetDC(window.hwnd)

	px := C.PIXELFORMATDESCRIPTOR{
		dwFlags:    0x00000001 | 0x00000020 | 0x00000004,
		iPixelType: 0,
		cColorBits: 32,
		cDepthBits: 32,
		iLayerType: 0,
	}

	px.nSize = (C.WORD)(unsafe.Sizeof(px))

	pxFormat := C.ChoosePX(window.hdc, &px)
	if pxFormat == 0 {
		return nil
	}

	res := C.SetPixelFormat(window.hdc, pxFormat, &px)
	if res == 0 {
		return nil
	}

	return unsafe.Pointer(&window.hdc)
}

func (window *Win32Window) SetPixelFormat(depth byte) {

}

func (window *Win32Window) Peek() bool {
	var msg C.MSG

	res := C.PeekMessageA((*C.struct_tagMSG)(unsafe.Pointer(&msg)), nil, 0, 0, C.PM_REMOVE)
	if res == 1 {
		C.TranslateMessage(&msg)
		C.DispatchMessage(&msg)

		if msg.message == C.WM_QUIT {
			window.WndProc(window.hwnd, msg.message, 0, 0)
		}

		return true
	}

	return false
}

func (window *Win32Window) Create(width uint, height uint, title string) bool {

	hInstance := C.GetModuleHandle(nil)

	if hInstance == nil {
		fmt.Println("Error")
		return false
	}

	res := C.RegisterWindow(C.CString("WindowClass"), hInstance, (C.WNDPROC)(unsafe.Pointer(syscall.NewCallback(func(hwnd C.HWND, message C.UINT, wParam C.WPARAM, lparam C.LPARAM) C.LRESULT {
		return window.WndProc(hwnd, message, wParam, lparam)
	}))))

	if res == 0 {
		fmt.Println("ERRRROR")
	}

	hwnd := C.CreateWindowTest(C.CString("WindowClass"), C.CString(title), (C.UINT)(width), (C.UINT)(height), hInstance)

	window.width = width
	window.height = height
	window.hInstance = hInstance
	window.hwnd = hwnd

	C.ShowWindow(hwnd, C.SW_SHOW)

	return true
}

func (window *Win32Window) SwapBuffers() {
	C.SwapBuffers(*(*C.HDC)(window.GetHDC()))
}
