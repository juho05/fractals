package main

import (
	"sync"

	"github.com/Bananenpro/fractals"
	"github.com/Bananenpro/fractals/generate"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const windowWidth = 800
const windowHeight = 800
const zoomSpeed = 0.05

var generator *generate.Generator

var points = make([]fractals.Point, windowWidth*windowHeight)
var pointsLock = sync.RWMutex{}

var camera generate.Camera
var maxIterations int

func processInput() {
	generator.BeginMovement()
	defer generator.EndMovement()

	wheelMove := float64(rl.GetMouseWheelMove())
	generator.Zoom(wheelMove*zoomSpeed, int(rl.GetMouseX()), int(rl.GetMouseY()))

	if rl.IsMouseButtonDown(rl.MouseMiddleButton) {
		mouseMove := rl.GetMouseDelta()
		generator.Move(int(mouseMove.X), int(mouseMove.Y))
	}
}

func generatorCallback(pointsData []fractals.Point, usedCamera generate.Camera, usedMaxIterations int) {
	pointsLock.Lock()
	points = pointsData
	pointsLock.Unlock()

	camera = usedCamera
	maxIterations = usedMaxIterations
}

func main() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(int32(windowWidth), int32(windowHeight), "Fractals")

	generator = generate.NewJuliaGenerator(-0.1+0.65i, windowWidth, windowHeight)
	generator.AddCallback(generatorCallback)

	generator.Start(true)

	for !rl.WindowShouldClose() {
		processInput()

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		pointsLock.RLock()
		for _, p := range points {
			rl.DrawPixel(int32(p.X), int32(p.Y), fractals.BernsteinPolynomials(p.Iterations, maxIterations))
		}
		pointsLock.RUnlock()

		rl.EndDrawing()
	}

	generator.Stop()

	rl.CloseWindow()
}
