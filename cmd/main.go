package main

import (
	"sync"

	"github.com/Bananenpro/fractals"
	"github.com/Bananenpro/fractals/generate"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	windowWidth  = 800
	windowHeight = 800
)

const zoomSpeed = 0.1

var generator *generate.Generator

var points = [][]fractals.Point{}
var pointsLock = sync.RWMutex{}

var camera generate.Camera
var maxIterations int
var deltaTime int64

func processInput() {
	generator.BeginMovement()
	defer generator.EndMovement()

	wheelMove := float64(rl.GetMouseWheelMove())
	if wheelMove != 0 {
		generator.Zoom(wheelMove*zoomSpeed, int(rl.GetMouseX()), int(rl.GetMouseY()))
	}

	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		mouseMove := rl.GetMouseDelta()
		generator.Move(int(mouseMove.X), int(mouseMove.Y))
	}
}

func generatorCallback(pointsData [][]fractals.Point, usedCamera generate.Camera, usedMaxIterations int, time int64) {
	pointsLock.Lock()
	points = pointsData
	pointsLock.Unlock()

	camera = usedCamera
	maxIterations = usedMaxIterations
	deltaTime = time
}

func newGenerator() *generate.Generator {
	size := windowWidth
	if windowWidth > windowHeight {
		size = windowHeight
	}

	generator = generate.NewMandelbrotGenerator(size, size)

	generator.AddCallback(generatorCallback)

	return generator
}

func main() {
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagWindowResizable)
	rl.InitWindow(int32(windowWidth), int32(windowHeight), "Fractals")

	loadAssets()

	generator := newGenerator()
	generator.Start(true)

	for !rl.WindowShouldClose() {
		if rl.IsWindowResized() {
			generator.Stop()
			camera = generator.GetCamera()
			maxIterations = generator.GetMaxIterations()
			windowWidth = rl.GetScreenWidth()
			windowHeight = rl.GetScreenHeight()
			pointsLock.Lock()
			points = [][]fractals.Point{}
			pointsLock.Unlock()
			generator = newGenerator()
			generator.SetCamera(camera)
			generator.SetMaxIterations(maxIterations)
			generator.Start(true)
		}

		processInput()

		rl.BeginDrawing()

		rl.ClearBackground(rl.NewColor(10, 10, 10, 255))

		pointsLock.RLock()
		for i := range points {
			for _, p := range points[i] {
				x := int32(p.X)
				if windowWidth > windowHeight {
					x += int32(windowWidth-windowHeight) / 2
				}

				y := int32(p.Y)
				if windowHeight > windowWidth {
					y += int32(windowHeight-windowWidth) / 2
				}

				rl.DrawPixel(x, y, fractals.BernsteinPolynomials(p.Iterations, maxIterations))
			}
		}
		pointsLock.RUnlock()

		renderGui()

		rl.EndDrawing()
	}

	generator.Stop()

	rl.CloseWindow()
}
