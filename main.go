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
const chunkWidth = 100
const chunkHeight = 100
const maxIterations = 500
const zoomSpeed = 0.05

var camera = Camera{
	Zoom:    1.0,
	offsetX: 0.0,
	offsetY: 0.0,
}
var cameraMutex = sync.RWMutex{}
var regenerateChan = make(chan bool, 1)

var running = true

var points = make([]Point, windowWidth*windowHeight)
var pointsLock = sync.RWMutex{}

func colorFromIterations(iterations int64) color.RGBA {
	normalized := float64(iterations) / float64(maxIterations)
	if normalized == 0 || normalized == 1 {
		return rl.Black
	}
	// return colors[iterations%int64(len(colors))]
	r := uint8(9 * (1 - normalized) * normalized * normalized * normalized * 255)
	g := uint8(15 * (1 - normalized) * (1 - normalized) * normalized * normalized * 255)
	b := uint8(8.5 * (1 - normalized) * (1 - normalized) * (1 - normalized) * normalized * 255)
	return rl.NewColor(r, g, b, 255)
}

func processInput() {
	changed := false

	cameraMutex.Lock()

	wheelMove := float64(rl.GetMouseWheelMove())
	if wheelMove != 0 {
		camera.offsetX -= float64(windowWidth / 2 - rl.GetMouseX()) * camera.Zoom / (windowWidth / 4)
		camera.offsetY -= float64(windowHeight / 2 - rl.GetMouseY()) * camera.Zoom / (windowHeight / 4)
		camera.Zoom -= wheelMove * zoomSpeed * camera.Zoom
		camera.offsetX += float64(windowWidth / 2 - rl.GetMouseX()) * camera.Zoom / (windowWidth / 4)
		camera.offsetY += float64(windowHeight / 2 - rl.GetMouseY()) * camera.Zoom / (windowHeight / 4)
		changed = true
	}

	if rl.IsMouseButtonDown(rl.MouseMiddleButton) {
		mouseMove := rl.GetMouseDelta()
		if mouseMove.X != 0 || mouseMove.Y != 0 {
			camera.offsetX -= float64(mouseMove.X) * camera.Zoom / (windowWidth / 4)
			camera.offsetY -= float64(mouseMove.Y) * camera.Zoom / (windowHeight / 4)
			changed = true
		}
	}

	cameraMutex.Unlock()

	if changed {
		select {
		case regenerateChan <- true:
			return
		default:
			return
		}
	}
}

func generateLoop() {
	for running {
		cameraMutex.RLock()
		cam := camera
		cameraMutex.RUnlock()

		temp := generateMandelbrot(cam, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)
		// temp := generateJulia(-0.1+0.65i, cam, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)
		// temp := generateJulia(-0.79+0.125i, cam, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)
		// temp := generateJulia(-1.45+0i, cam, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)
		// temp := generateJulia(-1.37969+0i, cam, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)
		// temp := generateJulia(-0.562292+0.642817i, cam, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)

		pointsLock.Lock()
		points = temp
		pointsLock.Unlock()

		<-regenerateChan
	}
}

func main() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(int32(windowWidth), int32(windowHeight), "Fractals")

	go generateLoop()

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

	running = false
	select {
	case regenerateChan <- true:
		break
	default:
		break
	}
}
