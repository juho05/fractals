package main

import "math/cmplx"

func generateMandelbrot(camera Camera, width, height, chunkWidth, chunkHeight int, maxIterations int64) []Point {
	channel := make(chan Chunk, (windowWidth/chunkWidth)*(windowHeight/chunkHeight))
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

			go generateMandelbrotChunk(index, camera, x, y, toX, toY, width, height, channel, maxIterations)
			index++
		}
	}

	points := make([]Point, 0, windowWidth*windowHeight)

	for i := 0; i < (windowWidth/chunkWidth)*(windowHeight/chunkHeight); i++ {
		chunk := <-channel
		points = append(points, chunk.Points...)
	}

	return points
}

// generate a chunk of the mandelbrot set from (fromX, fromY) (inclusive) to (toX, toY) (exclusive) and send the result through channel
func generateMandelbrotChunk(index int, camera Camera, fromX, fromY, toX, toY, width, height int, channel chan<- Chunk, maxIterations int64) {
	points := make([]Point, 0, (toX-fromX)*(toY-fromY))
	for x := fromX; x < toX; x++ {
		for y := fromY; y < toY; y++ {
			points = append(points, Point{
				X:          x,
				Y:          y,
				Iterations: calculateMandelbrotPoint(camera, x, y, width, height, maxIterations),
			})
		}
	}

	channel <- Chunk{
		Index:  index,
		Points: points,
	}
}

func calculateMandelbrotPoint(camera Camera, winX, winY, winWidth, winHeight int, maxIterations int64) int64 {
	iteration := int64(0)

	c := complex(camera.offsetX+(float64(winX)/float64(winWidth)-0.5)*camera.Zoom*4, camera.offsetY+(float64(winY)/float64(winHeight)-0.5)*camera.Zoom*4)
	z := complex(0, 0)
	zsquared := z * z

	for cmplx.Abs(zsquared) <= 4 && iteration < maxIterations {
		z = zsquared + c
		zsquared = z * z
		iteration++
	}

	return iteration
}
