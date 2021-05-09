package main

import (
	//"os"
	//"fmt"
	//"strings"
	"github.com/veandco/go-sdl2/sdl"
	"runtime"
)

const (
	BOTTOM_BAR_HEIGHT = 32.0
)

type Mode uint

const (
	Drawing = 0
	Command = 1
)

var (
	mode Mode

	window_w int32
	window_h int32

	x    float32 = 0.0
	y    float32 = 0.0
	zoom float32 = 1.0

	path string
)

func render(r *sdl.Renderer) {
	r.SetDrawColor(32, 35, 40, 255)
	r.Clear()
	r.SetDrawColor(14, 15, 17, 255)
	r.FillRect(&sdl.Rect{0, window_h - BOTTOM_BAR_HEIGHT, window_w, BOTTOM_BAR_HEIGHT})
	r.SetDrawColor(7, 7, 8, 255)
	r.FillRect(&sdl.Rect{0, window_h - BOTTOM_BAR_HEIGHT, window_w, 2})
	r.Present()
}

func main() {
	runtime.LockOSThread()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Go-Sprite", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	window.SetResizable(true)

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	running := true
	for running {
		window_w, window_h = window.GetSize()

		render(renderer)

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			case *sdl.KeyboardEvent:
				keyCode := t.Keysym.Sym
				switch keyCode {
				case 27:
					mode = Drawing
				}
			case *sdl.TextInputEvent:
				if t.GetText() == ":" {
					mode = Command
				}
			}
		}
	}
}
