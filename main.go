package main

import (
	"os"
	"fmt"
	"math"
	"image"
	"strings"
	"runtime"
	"io/ioutil"
	"image/color"
	_ "image/png"
	"golang.org/x/image/font"
	"github.com/faiface/pixel"
	"github.com/flopp/go-findfont"
	"github.com/faiface/pixel/text"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/golang/freetype/truetype"
)

const (
	BOTTOM_BAR_HEIGHT = 24.0
	CHECKER_SIZE = 16
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

	x    float64 = 0.0
	y    float64 = 0.0
	zoom float64 = 1.0
	ipos pixel.Vec
	mat pixel.Matrix

	command string
	notif string

	loaded pixel.Picture
	sprite *pixel.Sprite

	checker_background *imdraw.IMDraw
)

func loadTTF(command string, size float64) (font.Face, error) {
	fontPath, err := findfont.Find(command)
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func redrawChecker(win *pixelgl.Window) {
	checker_background = imdraw.New(nil)
	checker_background.EndShape = imdraw.SharpEndShape
	if (sprite != nil) {
		checker_background.SetMatrix(mat)
		PB := loaded.Bounds()
		i := 0
		for cx := -(PB.Max.X - PB.Min.X) / 2; cx < PB.W() - (PB.Max.X - PB.Min.X) / 2; cx += CHECKER_SIZE {
			for cy := -(PB.Max.Y - PB.Min.Y) / 2; cy < PB.H() - (PB.Max.Y - PB.Min.Y) / 2; cy += CHECKER_SIZE {
				if (i % 2 == 0) {
					checker_background.Color = color.RGBA{200, 200, 200, 255}
				} else {
					checker_background.Color = color.RGBA{150, 150, 150, 255}
				}
				checker_background.Push(pixel.V(cx, PB.Min.Y - cy), pixel.V(math.Min(cx + CHECKER_SIZE, PB.W() - (PB.Max.X - PB.Min.X) / 2), PB.Min.Y - math.Min(cy + CHECKER_SIZE, PB.H() - (PB.Max.Y - PB.Min.Y) / 2)))
				checker_background.Rectangle(0)
				i++
			}
		}
	}
}

func executeCommand(win *pixelgl.Window) {
	tokens := strings.Split(command, " ")
	switch tokens[0] {
	case "e":
		if len(tokens) > 2 {
			notif = "Provide one file path for opening only."
			return
		}
		pic, err := loadPicture(tokens[1])
		if err != nil {
			panic(err)
		}
		loaded = pic
		sprite = pixel.NewSprite(pic, pic.Bounds())
		zoom = 1.0
		redrawChecker(win)
	}
}

func handleInput(win *pixelgl.Window) {
	if mode == Drawing {
		if win.Typed() == ":" {
			mode = Command
			command = ""
			notif = ""
		}
		if win.Pressed(pixelgl.KeyEscape) {
			notif = ""
		}
	} else if mode == Command {
		if win.Pressed(pixelgl.KeyEscape) {
			mode = Drawing
			command = ""
			notif = ""
		} else if win.Pressed(pixelgl.KeyEnter){
			executeCommand(win)
			mode = Drawing
			command = ""
		} else if (win.JustPressed(pixelgl.KeyBackspace) || win.Repeated(pixelgl.KeyBackspace)) && len(command) > 0 {
			command = command[:len(command)-1]
		} else {
			command += win.Typed()
		}
	}
	if sprite != nil {
		if win.Pressed(pixelgl.MouseButton2) {
			mcpos := win.MousePosition()
			mppos := win.MousePreviousPosition()
			x += (mppos.X - mcpos.X) / zoom
			y += (mppos.Y - mcpos.Y) / zoom
		}
		if win.MouseScroll().Y != 0 {
			zoom *= math.Pow(2.0, win.MouseScroll().Y)
		}
		mat = pixel.IM.Moved(pixel.Vec{-x, -y}).Scaled(loaded.Bounds().Center(), zoom).Moved(win.Bounds().Center())
		mpos := win.MousePosition()
		ipos = mat.Unproject(mpos).Add(loaded.Bounds().Center())
		ipos.X = math.Round(ipos.X-0.5)
		ipos.Y = math.Round(ipos.Y-0.5)
		redrawChecker(win)
		fmt.Println(x, y, ipos)
	}
}

func render(win *pixelgl.Window, txt *text.Text) {
	win.Clear(color.RGBA{32, 35, 40, 255})

	imd := imdraw.New(nil)
	imd.EndShape = imdraw.SharpEndShape

	imd.Color = color.RGBA{14, 15, 17, 255}
	imd.Push(pixel.V(0, 0), pixel.V(float64(window_w), BOTTOM_BAR_HEIGHT))
	imd.Rectangle(0)

	imd.Color = color.RGBA{7, 7, 8, 255}
	imd.Push(pixel.V(0, BOTTOM_BAR_HEIGHT - 2), pixel.V(float64(window_w), BOTTOM_BAR_HEIGHT))
	imd.Rectangle(0)

	txt.Clear()

	if mode == Command {
		txt.Color = color.RGBA{44, 151, 244, 255}
		txt.WriteString(":")

		txt.Color = color.RGBA{156, 160, 164, 255}
		txt.WriteString(command)

		imd.Color = color.RGBA{44, 151, 244, 255}
		imd.Push(pixel.V(txt.Bounds().Max.X + 8, txt.Bounds().Min.Y + 8), pixel.V(txt.Bounds().Max.X + 10, txt.Bounds().Max.Y + 4))
		imd.Rectangle(0)
	} else if mode == Drawing {
		txt.Color = color.RGBA{156, 160, 164, 255}
		txt.WriteString(notif)
	}

	if (sprite != nil) {
		checker_background.Draw(win)
		sprite.Draw(win, mat)
		mbox := imdraw.New(nil)
		mbox.Color = color.RGBA{255, 255, 255, 255}
		mbox.Push(mat.Project(ipos.Sub(loaded.Bounds().Center())), mat.Project(ipos.Sub(loaded.Bounds().Center()).Add(pixel.Vec{1.0, 1.0})))
		mbox.Rectangle(1)
		mbox.Draw(win)
	}

	imd.Draw(win)
	txt.Draw(win, pixel.IM.Moved(pixel.V(6, 4)))

	win.Update()
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

	face, err := loadTTF("DejaVuSansMono.ttf", 16)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(pixel.V(0, 0), atlas)

	for !win.Closed() {
		bounds := win.Bounds()
		n_window_w := int32(bounds.Max.X - bounds.Min.X)
		n_window_h := int32(bounds.Max.Y - bounds.Min.Y)
		if window_w != n_window_w || window_h != n_window_h {
			fmt.Println(n_window_w, n_window_h)
			window_w = n_window_w
			window_h = n_window_h
			redrawChecker(win)
		}

		handleInput(win)
		render(win, txt)
	}
}

func main() {
	pixelgl.Run(run)
}
