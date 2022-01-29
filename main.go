package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Point struct {
	X          int
	Y          int
	Iterations int64
}

type Chunk struct {
	Index  int
	Points []Point
}

const windowWidth = 800
const windowHeight = 800
const chunkWidth = 80
const chunkHeight = 80
const maxIterations = 500

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

var chunks = make([]Chunk, (windowWidth/chunkWidth)*(windowHeight/chunkHeight))

func colorFromIterations(iterations int64) color.RGBA {
	if iterations == maxIterations {
		return rl.Black
	}
	return colors[iterations%int64(len(colors))]
}

func drawChunk(chunk Chunk) {
	for _, p := range chunk.Points {
		rl.DrawPixel(int32(p.X), int32(p.Y), colorFromIterations(p.Iterations))
	}
}

func receiveChunks(channel <-chan Chunk) {
	for {
		select {
		case chunk := <-channel:
			chunks[chunk.Index] = chunk
		default:
			return
		}
	}
}

func main() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(int32(windowWidth), int32(windowHeight), "Fractals")

	channel := make(chan Chunk, (windowWidth/chunkWidth)*(windowHeight/chunkHeight))

	// generateMandelbrot(channel, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)
	generateJulia(-0.1+0.65i, channel, int(windowWidth), int(windowHeight), chunkWidth, chunkHeight, maxIterations)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		receiveChunks(channel)
		for _, chunk := range chunks {
			drawChunk(chunk)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
