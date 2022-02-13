package generate

import (
	"math/cmplx"
)

func NewMandelbrotGenerator(width, height int) *Generator {
	g := newGenerator(width, height)

	g.calculatePixel = func(x, y int) int {
		iteration := 0

		c := g.complexNumberFromPixel(x, y)
		z := complex(0, 0)
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
