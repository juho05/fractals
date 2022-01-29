package main

import "math/cmplx"

func generateMandelbrot(channel chan<- Chunk, width, height, chunkWidth, chunkHeight int, maxIterations int64) {
	index := 0
	for x := 0; x < width; x += chunkWidth {
		for y := 0; y < height; y += chunkHeight {
			toX := x + chunkWidth
			if toX >= width {
				toX = width
			}
			toY := y + chunkHeight
			if toY >= height {
				toY = height
			}

			go generateMandelbrotChunk(index, x, y, toX, toY, width, height, channel, maxIterations)
			index++
		}
	}
}

// generate a chunk of the mandelbrot set from (fromX, fromY) (inclusive) to (toX, toY) (exclusive) and send the result through channel
func generateMandelbrotChunk(index, fromX, fromY, toX, toY, width, height int, channel chan<- Chunk, maxIterations int64) {
	points := make([]Point, 0, (toX-fromX)*(toY-fromY))
	for x := fromX; x < toX; x++ {
		for y := fromY; y < toY; y++ {
			points = append(points, Point{
				X:          x,
				Y:          y,
				Iterations: calculateMandelbrotPoint(x, y, width, height, maxIterations),
			})
		}
	}

	channel <- Chunk{
		Index:  index,
		Points: points,
	}
}

func calculateMandelbrotPoint(winX, winY, winWidth, winHeight int, maxIterations int64) int64 {
	iteration := int64(0)

	c := complex(mapFloat(float64(winX), 0, float64(winWidth), -2, 2), mapFloat(float64(winY), 0, float64(winHeight), -2, 2))
	z := complex(0, 0)
	zsquared := z * z

	for cmplx.Abs(zsquared) <= 4 && iteration < maxIterations {
		z = zsquared + c
		zsquared = z * z
		iteration++
	}

	return iteration
}
