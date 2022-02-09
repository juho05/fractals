package main

import (
	"fmt"
	"strings"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/sqweek/dialog"
)

var font rl.Font

func loadAssets() {
	font = rl.LoadFontEx("assets/fonts/Roboto/Roboto-Regular.ttf", 16, nil, 0)
}

func renderGui() {
	rl.DrawRectangle(0, 0, windowWidth, 30, rl.NewColor(0, 0, 0, 100))
	rl.DrawTextEx(font, fmt.Sprintf("NMAX: %d\tTIME: %dms", maxIterations, deltaTime), rl.Vector2{X: 5, Y: 7}, 16, 0, rl.White)

	rl.DrawRectangle(0, windowHeight-25, windowWidth, 25, rl.NewColor(0, 0, 0, 100))
	rl.DrawTextEx(font, fmt.Sprintf("SCALE: %g\tOFFSET-X: %g\tOFFSET-Y: %g", camera.Scale, camera.OffsetX, camera.OffsetY), rl.Vector2{X: 5, Y: windowHeight - 18}, 16, 0, rl.White)

	if raygui.Button(rl.NewRectangle(windowWidth-55, 5, 50, 20), "SAVE") {
		save()
	}
}

func save() {
	path, err := dialog.File().Filter("PNG image file", "png").Save()
	if err != nil {
		dialog.Message("Failed to open provided path!").Title("Error").Error()
		return
	}

	if !strings.HasSuffix(strings.ToLower(path), ".png") {
		path += ".png"
	}

	err = saveImageToDisk(path)
	if err != nil {
		dialog.Message("Failed to save file: %s", err).Title("Error").Error()
		return
	}

	dialog.Message("Successfully saved image!").Title("Success").Info()
}
