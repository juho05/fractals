package main

import (
	"sync"

	"github.com/Bananenpro/fractals"
	"github.com/Bananenpro/fractals/generate"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const windowWidth = 800
const windowHeight = 800
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
	if time > 15 {
		deltaTime = time
	}
}

func main() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(int32(windowWidth), int32(windowHeight), "Fractals")

	loadAssets()

	generator = generate.NewMandelbrotGenerator(windowWidth, windowHeight)
	generator.AddCallback(generatorCallback)

	generator.Start(true)

	for !rl.WindowShouldClose() {
		processInput()

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		pointsLock.RLock()
		for i := range points {
			for _, p := range points[i] {
				rl.DrawPixel(int32(p.X), int32(p.Y), fractals.BernsteinPolynomials(p.Iterations, maxIterations))
			}
		}
		pointsLock.RUnlock()

		renderGui()

		rl.EndDrawing()
	}

	generator.Stop()

	rl.CloseWindow()
}
