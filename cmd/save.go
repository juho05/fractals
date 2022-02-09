package main

import (
	"image"
	"image/png"
	"os"

	"github.com/Bananenpro/fractals"
)

func saveImageToDisk(path string) error {
	img := constructImage()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

func constructImage() image.Image {
	pointsLock.RLock()

	img := image.NewRGBA(image.Rect(0, 0, windowWidth, windowHeight))
	for _, p := range points {
		img.Set(p.X, p.Y, fractals.BernsteinPolynomials(p.Iterations, maxIterations))
	}

	pointsLock.RUnlock()

	return img
}
