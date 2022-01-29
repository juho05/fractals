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
