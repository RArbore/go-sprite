package main

import (
	//"os"
	"fmt"
	//"strings"
	"runtime"
	//"image/font"
	"image/color"
	"github.com/faiface/pixel"
	//"github.com/faiface/pixel/text"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	//"github.com/golang/freetype/truetype"
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

func handleInput(win *pixelgl.Window) {
	if win.Pressed(pixelgl.KeyEscape) {
		mode = Drawing
	}
	if win.Typed() == ":" {
		mode = Command
	}
}

func render(win *pixelgl.Window) {
	win.Clear(color.RGBA{32, 35, 40, 255})

	imd := imdraw.New(nil)
	imd.EndShape = imdraw.SharpEndShape

	imd.Color = color.RGBA{14, 15, 17, 255}
	imd.Push(pixel.V(0, 0), pixel.V(float64(window_w), BOTTOM_BAR_HEIGHT))
	imd.Rectangle(0)

	imd.Color = color.RGBA{7, 7, 8, 255}
	imd.Push(pixel.V(0, BOTTOM_BAR_HEIGHT - 2), pixel.V(float64(window_w), BOTTOM_BAR_HEIGHT))
	imd.Rectangle(0)

	imd.Draw(win)
}

func run() {
	runtime.LockOSThread()

	cfg := pixelgl.WindowConfig {
		Title: "Go-Sprite",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync: true,
		Resizable: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		bounds := win.Bounds()
		window_w = int32(bounds.Max.X - bounds.Min.X)
		window_h = int32(bounds.Max.Y - bounds.Min.Y)

		handleInput(win)
		render(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
