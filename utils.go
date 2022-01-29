package main

type Point struct {
	X          int
	Y          int
	Iterations int64
}

type Chunk struct {
	Index  int
	Points []Point
}

func mapFloat(nr, min, max, mappedMin, mappedMax float64) float64 {
	return (nr-min)/(max-min)*(mappedMax-mappedMin) + mappedMin
}
