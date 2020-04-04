package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

const (
	screenWidth  = 600 // must be a square
	screenHeight = 600
	frameRate    = 10
)

func main() {
	setup()
}

func setup() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("initializing SDL:", err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"Tic Tac Go",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("initializing window:", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("initializing renderer:", err)
		return
	}
	defer renderer.Destroy()

	img.Init(img.INIT_JPG | img.INIT_PNG)

	textures := loadGophers(renderer)

	b := newBoard(screenWidth, renderer, textures)

	b.drawBoard()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.MouseButtonEvent:
				gameOver, winner := b.handleMouseButtonEvent(t)
				b.drawBoard()

				if gameOver {
					running = false
					if winner == 0 {
						fmt.Printf("It's a draw!!!\n")
					} else {
						fmt.Printf("Player %d wins!!!\n", winner)
					}
					break
				}
			}

			sdl.Delay(1000 / frameRate)
		}
	}
}

func loadGophers(renderer *sdl.Renderer) []*sdl.Texture {
	img1 := make(chan *sdl.Surface, 1)
	img2 := make(chan *sdl.Surface, 1)

	go loadImage("assets/gopher1.png", img1)
	go loadImage("assets/gopher2.png", img2)

	gopher1 := <-img1
	gopher2 := <-img2

	gopher1Texture, err := renderer.CreateTextureFromSurface(gopher1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		os.Exit(5)
	}

	gopher2Texture, err := renderer.CreateTextureFromSurface(gopher2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		os.Exit(5)
	}

	gopher1Texture.Unlock()
	gopher2Texture.Unlock()

	return []*sdl.Texture{gopher1Texture, gopher2Texture}
}

func loadImage(path string, res chan *sdl.Surface) {
	image, err := img.Load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to gopher2 PNG: %s\n", err)
		os.Exit(4)
	}

	res <- image
}
