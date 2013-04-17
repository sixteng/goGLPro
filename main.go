package main

import "fmt"
import "./src/window"

import "./src/renderer/devices/opengl"

func main() {
	shutdown := false
	mainwindow := window.CreateWindow()

	mainwindow.Create(200, 200, "Hello")
	mainwindow.SetCallback(func(window window.Window, message uint32, wParam uintptr, lparam uintptr) bool {
		if message == 0x0012 {
			fmt.Println("Shutdown!!!")
			shutdown = true
		}

		return false
	})

	opengl.CreateContext(mainwindow)

	mainwindow.GetHDC()
	for !shutdown {
		if mainwindow.Peek() {
			continue
		}

		//Free Loop
		//Render and do engine logic

		opengl.Test(mainwindow)
	}
}
