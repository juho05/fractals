package main

import (
	"fmt"
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
var deltaTime int64

func processInput() {
	generator.BeginMovement()
	defer generator.EndMovement()

	wheelMove := float64(rl.GetMouseWheelMove())
	if wheelMove != 0 {
		generator.Zoom(wheelMove*zoomSpeed, int(rl.GetMouseX()), int(rl.GetMouseY()))
	}

	if rl.IsMouseButtonDown(rl.MouseMiddleButton) {
		mouseMove := rl.GetMouseDelta()
		generator.Move(int(mouseMove.X), int(mouseMove.Y))
	}
}

func generatorCallback(pointsData []fractals.Point, usedCamera generate.Camera, usedMaxIterations int, time int64) {
	pointsLock.Lock()
	points = pointsData
	pointsLock.Unlock()

	camera = usedCamera
	maxIterations = usedMaxIterations
	deltaTime = time
}

func main() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(int32(windowWidth), int32(windowHeight), "Fractals")

	font := rl.LoadFontEx("assets/fonts/Roboto/Roboto-Regular.ttf", 16, nil, 0)

	generator = generate.NewMandelbrotGenerator(windowWidth, windowHeight)
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

		rl.DrawRectangle(0, 0, windowWidth, 25, rl.NewColor(0, 0, 0, 100))
		rl.DrawTextEx(font, fmt.Sprintf("NMAX: %d\tTIME: %dms", maxIterations, deltaTime), rl.Vector2{X: 5, Y: 5}, 16, 0, rl.White)

		rl.DrawRectangle(0, windowHeight - 25, windowWidth, 25, rl.NewColor(0, 0, 0, 100))
		rl.DrawTextEx(font, fmt.Sprintf("SCALE: %g\tOFFSET-X: %g\tOFFSET-Y: %g", camera.Scale, camera.OffsetX, camera.OffsetY), rl.Vector2{X: 5, Y: windowHeight - 18}, 16, 0, rl.White)

		rl.EndDrawing()
	}

	generator.Stop()

	rl.CloseWindow()
}
