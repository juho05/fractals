package main

func mapFloat(nr, min, max, mappedMin, mappedMax float64) float64 {
	return (nr-min)/(max-min)*(mappedMax-mappedMin) + mappedMin
}
