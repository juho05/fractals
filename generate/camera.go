package generate

type Camera struct {
	Zoom    float64
	OffsetX float64
	OffsetY float64
}

func (g *Generator) BeginMovement() {
	g.previousCamera = g.camera
}

func (g *Generator) EndMovement() {
	g.regenerateChan <- true
}

func (g *Generator) Zoom(delta float64, pixelCenterX, pixelCenterY int) {
	g.camera.OffsetX -= float64(g.width/2-pixelCenterX) * g.camera.Zoom / float64(g.width/4)
	g.camera.OffsetY -= float64(g.height/2-pixelCenterY) * g.camera.Zoom / float64(g.height/4)
	g.camera.Zoom -= delta * g.camera.Zoom
	g.camera.OffsetX += float64(g.width/2-pixelCenterX) * g.camera.Zoom / float64(g.width/4)
	g.camera.OffsetY += float64(g.height/2-pixelCenterY) * g.camera.Zoom / float64(g.height/4)
}

func (g *Generator) Move(dPixelX, dPixelY int) {
	g.camera.OffsetX -= float64(dPixelX) * g.camera.Zoom / float64(g.width/4)
	g.camera.OffsetY -= float64(dPixelY) * g.camera.Zoom / float64(g.height/4)
}
