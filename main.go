package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	windowWidth := int32(800)
	windowHeight := int32(800)

	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(windowWidth, windowHeight, "fractals")

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
