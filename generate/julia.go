package generate

import (
	"math/cmplx"
)

func NewJuliaGenerator(c complex128, width, height int) *Generator {
	g := newGenerator(width, height)

	g.calculatePixel = func(x, y int) int {
		iteration := 0

		z := g.complexNumberFromPixel(x, y)
		zsquared := z * z

		for cmplx.Abs(zsquared) <= 4 && iteration < g.maxIterations {
			z = zsquared + c
			zsquared = z * z
			iteration++
		}

		return iteration
	}

	return g
}
