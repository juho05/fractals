package generate

type Camera struct {
	Scale   float64
	OffsetX float64
	OffsetY float64
}

func (g *Generator) BeginMovement() {
	g.cameraLock.Lock()
	g.previousCamera = g.camera
}

func (g *Generator) EndMovement() {
	g.cameraLock.Unlock()
	if g.camera != g.previousCamera {
		g.regenerateChan <- true
	}
}

func (g *Generator) Zoom(delta float64, pixelCenterX, pixelCenterY int) {
	g.camera.OffsetX -= float64(g.width/2-pixelCenterX) * g.camera.Scale / float64(g.width/4)
	g.camera.OffsetY -= float64(g.height/2-pixelCenterY) * g.camera.Scale / float64(g.height/4)
	g.camera.Scale -= delta * g.camera.Scale
	g.camera.OffsetX += float64(g.width/2-pixelCenterX) * g.camera.Scale / float64(g.width/4)
	g.camera.OffsetY += float64(g.height/2-pixelCenterY) * g.camera.Scale / float64(g.height/4)
}

func (g *Generator) Move(dPixelX, dPixelY int) {
	g.deltaX += dPixelX
	g.deltaY += dPixelY

	g.camera.OffsetX -= float64(dPixelX) * g.camera.Scale / float64(g.width/4)
	g.camera.OffsetY -= float64(dPixelY) * g.camera.Scale / float64(g.height/4)
}
