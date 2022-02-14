package generate

import (
	"sync"
	"time"

	"github.com/Bananenpro/fractals"
	"github.com/google/uuid"
)

type chunk struct {
	points []fractals.Point
	y      int
}

type calculatePixelFunc func(x, y int) int

type callbackFunc func(point [][]fractals.Point, camera Camera, maxIterations int, time int64)

const increaseIterationsThreshold = 0.93
const maxIterationsStep = 50

type Generator struct {
	camera         Camera
	previousCamera Camera
	cameraLock     sync.RWMutex
	deltaX         int
	deltaY         int

	width                 int
	height                int
	maxIterations         int
	previousMaxIterations int

	calculatePixel calculatePixelFunc
	symmetric      bool

	points [][]fractals.Point // y x

	callbacks map[uuid.UUID]callbackFunc

	regenerateChan chan bool

	running bool
}

// New generator without calculatePixel function
func newGenerator(width, height int) *Generator {
	return &Generator{
		camera: Camera{
			Scale: 1,
		},
		previousCamera: Camera{
			Scale: 1,
		},
		width:                 width,
		height:                height,
		maxIterations:         100,
		previousMaxIterations: 100,
		callbacks:             make(map[uuid.UUID]callbackFunc),
		regenerateChan:        make(chan bool, 3),
	}
}

func (g *Generator) Start(loop bool) {
	g.running = true
	go func() {
		for loop && g.running {
			g.cameraLock.RLock()
			g.generate()
			g.previousMaxIterations = g.maxIterations
			if g.camera.Scale != g.previousCamera.Scale {
				g.updateMaxIterations()
			}
			g.cameraLock.RUnlock()
			g.deltaX = 0
			g.deltaY = 0

			if !<-g.regenerateChan {
				return
			}
		}
	}()
}

func (g *Generator) Stop() {
	g.running = false
	select {
	case g.regenerateChan <- false:
		return
	default:
		return
	}
}

func (g *Generator) AddCallback(callback callbackFunc) uuid.UUID {
	id := uuid.New()
	g.callbacks[id] = callback
	return id
}

func (g *Generator) RemoveCallback(id uuid.UUID) {
	delete(g.callbacks, id)
}

func (g *Generator) GetMaxIterations() int {
	return g.maxIterations
}

func (g *Generator) SetMaxIterations(maxIterations int) {
	g.previousMaxIterations = g.maxIterations
	g.maxIterations = maxIterations
}

func (g *Generator) updateMaxIterations() {
	pixelsAboveIncreaseIterationsThreshold := 0
	for i := range g.points {
		for _, p := range g.points[i] {
			if float64(p.Iterations) > float64(g.maxIterations)*increaseIterationsThreshold && p.Iterations < g.maxIterations {
				pixelsAboveIncreaseIterationsThreshold++
			}
		}
	}

	if pixelsAboveIncreaseIterationsThreshold > 1500 {
		g.maxIterations += int(maxIterationsStep * (float64(pixelsAboveIncreaseIterationsThreshold) / 700))
	} else if pixelsAboveIncreaseIterationsThreshold < 1000 {
		g.maxIterations -= int(maxIterationsStep * (1000 / float64(pixelsAboveIncreaseIterationsThreshold)))
	}

	if g.maxIterations < 100 {
		g.maxIterations = 100
	}

	if g.maxIterations != g.previousMaxIterations {
		select {
		case g.regenerateChan <- true:
			return
		default:
			return
		}
	}
}

func (g *Generator) generate() {
	startTime := time.Now()
	channel := make(chan chunk, g.height)

	axisY := g.height
	if g.symmetric && g.camera.Scale < 0.5 {
		for y := 0; y < g.height; y++ {
			if imag(g.complexNumberFromPixel(0, y)) > 0 {
				axisY = y
				break
			}
		}
	}

	goroutineCount := 0
	if axisY >= g.height/2 {
		for y := 0; y < axisY; y++ {
			goroutineCount++
			go g.generateChunk(y, channel)
		}
	} else {
		for y := axisY; y < g.height; y++ {
			goroutineCount++
			go g.generateChunk(y, channel)
		}
	}

	points := make([][]fractals.Point, g.height)

	for i := 0; i < g.height; i++ {
		points[i] = make([]fractals.Point, g.width)
	}

	for i := 0; i < goroutineCount; i++ {
		chunk := <-channel
		points[chunk.y] = chunk.points
	}

	if axisY >= g.height/2 {
		for y := axisY; y < g.height; y++ {
			copy(points[y], points[axisY-(y-axisY)-1])
			for x := 0; x < g.width; x++ {
				points[y][x].Y = y
			}
		}
	} else {
		for y := 0; y < axisY; y++ {
			copy(points[y], points[2*axisY-y+1])
			for x := 0; x < g.width; x++ {
				points[y][x].Y = y
			}
		}
	}

	if g.running {
		for _, cb := range g.callbacks {
			cb(points, g.camera, g.maxIterations, time.Since(startTime).Milliseconds())
		}
	}

	g.points = points
}

func (g *Generator) generateChunk(y int, channel chan<- chunk) {
	points := make([]fractals.Point, 0, g.width)
	for x := 0; x < g.width; x++ {
		points = append(points, fractals.Point{
			X:          x,
			Y:          y,
			Iterations: g.generatePixel(x, y),
		})
	}

	channel <- chunk{
		points: points,
		y:      y,
	}
}

func (g *Generator) generatePixel(x, y int) int {
	if g.camera.Scale == g.previousCamera.Scale && g.maxIterations == g.previousMaxIterations {
		if len(g.points) > 0 {
			if x-g.deltaX >= 0 && y-g.deltaY >= 0 && x-g.deltaX < g.width && y-g.deltaY < g.height {
				return g.points[y-g.deltaY][x-g.deltaX].Iterations
			}
		}
	}

	return g.calculatePixel(x, y)
}

func (g *Generator) complexNumberFromPixel(x, y int) complex128 {
	return complex(g.camera.OffsetX+(float64(x)/float64(g.width)-0.5)*g.camera.Scale*4, g.camera.OffsetY+(float64(y)/float64(g.height)-0.5)*g.camera.Scale*4)
}
