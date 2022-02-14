package generate

import (
	"math/cmplx"
)

func NewMandelbrotGenerator(width, height int) *Generator {
	g := newGenerator(width, height)
	g.symmetric = true

	g.calculatePixel = func(x, y int) int {
		iteration := 0

		c := g.complexNumberFromPixel(x, y)
		z := complex(0, 0)

		if isInKnownShape(c) {
			return g.maxIterations
		}

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

func isInKnownShape(c complex128) bool {
	// cardioid
	q := (real(c)-0.25)*(real(c)-0.25) + imag(c)*imag(c)
	if q*(q+(real(c)-0.25)) < imag(c)*imag(c)*0.25 {
		return true
	}

	// period-2 bulb
	if real(c)*real(c)+(2*real(c))+1+imag(c)*imag(c) < 0.0625 {
		return true
	}

	return false
}
