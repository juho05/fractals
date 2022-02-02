package generate

import (
	"github.com/Bananenpro/fractals"
	"github.com/google/uuid"
)

type chunk []fractals.Point

type calculatePixelFunc func(camera Camera, x, y int) int

type callbackFunc func(points []fractals.Point, camera Camera, maxIterations int)

type Generator struct {
	camera         Camera
	previousCamera Camera
	width          int
	height         int
	maxIterations  int
	calculatePixel calculatePixelFunc

	points []fractals.Point

	callbacks map[uuid.UUID]callbackFunc

	regenerateChan chan bool

	running bool
}

// New generator without calculatePixel function
func newGenerator(width, height int) *Generator {
	return &Generator{
		camera: Camera{
			Zoom: 1,
		},
		previousCamera: Camera{
			Zoom: 1,
		},
		width:          width,
		height:         height,
		maxIterations:  100,
		callbacks:      make(map[uuid.UUID]callbackFunc),
		regenerateChan: make(chan bool, 1),
	}
}

func (g *Generator) Start(loop bool) {
	g.running = true
	go func() {
		for loop && g.running {
			go g.generate()
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

func (g *Generator) generate() {
	channel := make(chan chunk, g.height)

	camera := g.camera

	for y := 0; y < g.height; y++ {
		go g.generateChunk(camera, 0, y, g.width, y+1, channel)
	}

	points := make([]fractals.Point, 0, g.width*g.height)

	for i := 0; i < g.height; i++ {
		chunk := <-channel
		points = append(points, chunk...)
	}

	g.points = points
	for _, cb := range g.callbacks {
		cb(points, g.camera, g.maxIterations)
	}
}

func (g *Generator) generateChunk(camera Camera, fromX, fromY, toX, toY int, channel chan<- chunk) {
	points := make(chunk, 0, (toX-fromX)*(toY-fromY))
	for x := fromX; x < toX; x++ {
		for y := fromY; y < toY; y++ {
			points = append(points, fractals.Point{
				X:          x,
				Y:          y,
				Iterations: g.calculatePixel(camera, x, y),
			})
		}
	}

	channel <- points
}

func (g *Generator) complexNumberFromPixel(camera Camera, x, y int) complex128 {
	return complex(camera.OffsetX+(float64(x)/float64(g.width)-0.5)*camera.Zoom*4, camera.OffsetY+(float64(y)/float64(g.height)-0.5)*camera.Zoom*4)
}
