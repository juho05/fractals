package main

import (
	"image/color"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera struct {
	Zoom    float64
	offsetX float64
	offsetY float64
}

const windowWidth = 800
const windowHeight = 800
const chunkWidth = 80
const chunkHeight = 80
const maxIterations = 300
const zoomSpeed = 0.05

var camera = Camera{
	Zoom:    1.0,
	offsetX: 0.0,
	offsetY: 0.0,
}

var colors = []color.RGBA{
	rl.NewColor(66, 30, 15, 255),
	rl.NewColor(25, 7, 26, 255),
	rl.NewColor(9, 1, 47, 255),
	rl.NewColor(4, 4, 73, 255),
	rl.NewColor(0, 7, 100, 255),
	rl.NewColor(12, 44, 138, 255),
	rl.NewColor(24, 82, 177, 255),
	rl.NewColor(57, 125, 209, 255),
	rl.NewColor(134, 181, 229, 255),
	rl.NewColor(211, 236, 248, 255),
	rl.NewColor(248, 201, 95, 255),
	rl.NewColor(255, 170, 0, 255),
	rl.NewColor(204, 128, 0, 255),
	rl.NewColor(153, 87, 0, 255),
	rl.NewColor(106, 52, 3, 255),
}

var points = make([]Point, windowWidth*windowHeight)
var pointsLock = sync.RWMutex{}

func colorFromIterations(iterations int64) color.RGBA {
	if iterations == maxIterations {
		return rl.Black
	}
	return colors[iterations%int64(len(colors))]
}

func processInput() {
	changed := false

	wheelMove := float64(rl.GetMouseWheelMove())
	if wheelMove != 0 {
		// TODO: Zoom to cursor
		camera.Zoom -= wheelMove * zoomSpeed * camera.Zoom
		changed = true
	}

	if rl.IsMouseButtonDown(rl.MouseMiddleButton) {
		mouseMove := rl.GetMouseDelta()
		if mouseMove.X != 0 || mouseMove.Y != 0 {
			camera.offsetX -= float64(mouseMove.X) * camera.Zoom / 500
			camera.offsetY -= float64(mouseMove.Y) * camera.Zoom / 500
			changed = true
		}
	}

	if changed {
		go generate()
	}
}

func generate() {
	// TODO: cancel all previous calculations

	temp := generateMandelbrot(camera, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)
	// temp := generateJulia(-0.1+0.65i, camera, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)
	// temp := generateJulia(-0.79+0.125i, camera, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)
	// temp := generateJulia(-1.45+0i, camera, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)
	// temp := generateJulia(-1.37969+0i, camera, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)
	// temp := generateJulia(-0.562292+0.642817i, camera, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)

	pointsLock.Lock()
	points = temp
	pointsLock.Unlock()
}

func main() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(int32(windowWidth), int32(windowHeight), "Fractals")

	go generate()

	for !rl.WindowShouldClose() {
		processInput()

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		pointsLock.RLock()
		for _, p := range points {
			rl.DrawPixel(int32(p.X), int32(p.Y), colorFromIterations(p.Iterations))
		}
		pointsLock.RUnlock()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
