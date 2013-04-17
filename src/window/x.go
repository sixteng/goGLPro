// +build linux

package window

import "../../pkg/xgb"

type XWindow struct {
	width uint
}

func CreateWindow() Window {
	return &XWindow{}
}

func (window *XWindow) Peek() bool {

}

func (window *XWindow) Create(x int, y int, width uint, height uint, title string) bool {

}
