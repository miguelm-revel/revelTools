package randx

import (
	"math"
	"math/rand/v2"
)

type ExpDist struct {
	Lambda float64
}

func (e ExpDist) Rand() float64 {
	u := rand.Float64()
	if u == 0 {
		u = math.SmallestNonzeroFloat64
	}

	return (-1.0 / e.Lambda) * math.Log(u)
}

func (e ExpDist) PDF(x float64) float64 {
	if e.Lambda <= 0 {
		return math.NaN()
	}
	if x < 0 {
		return 0
	}
	return e.Lambda * math.Exp(-e.Lambda*x)
}

func (e ExpDist) CDF(x float64) float64 {
	if e.Lambda <= 0 {
		return math.NaN()
	}
	if x < 0 {
		return 0
	}
	return 1 - math.Exp(-e.Lambda*x)
}
