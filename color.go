package fractals

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func BernsteinPolynomials(iterations, maxIterations int) color.RGBA {
	normalized := float64(iterations) / float64(maxIterations)
	if normalized == 0 || normalized == 1 {
		return rl.Black
	}
	r := uint8(9 * (1 - normalized) * normalized * normalized * normalized * 255)
	g := uint8(15 * (1 - normalized) * (1 - normalized) * normalized * normalized * 255)
	b := uint8(8.5 * (1 - normalized) * (1 - normalized) * (1 - normalized) * normalized * 255)
	return rl.NewColor(r, g, b, 255)
}

var palette = []color.RGBA{
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

func ColorPalette(iterations, maxIterations int) color.RGBA {
	if iterations == maxIterations {
		return rl.Black
	}
	return palette[iterations%len(palette)]
}
