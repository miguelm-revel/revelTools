package randx

import (
	"math"
	"math/rand/v2"
)

// ExpDist represents an Exponential(Lambda) distribution.
// Lambda is the rate parameter (inverse of the mean); it must be positive.
type ExpDist struct {
	Lambda float64
}

// Rand returns a random sample using the inverse-CDF (inversion) method.
func (e ExpDist) Rand() float64 {
	u := rand.Float64()
	if u == 0 {
		u = math.SmallestNonzeroFloat64
	}

	return (-1.0 / e.Lambda) * math.Log(u)
}

// PDF returns the probability density at x.
// Returns NaN if Lambda <= 0; 0 if x < 0.
func (e ExpDist) PDF(x float64) float64 {
	if e.Lambda <= 0 {
		return math.NaN()
	}
	if x < 0 {
		return 0
	}
	return e.Lambda * math.Exp(-e.Lambda*x)
}

// CDF returns the cumulative probability P(X <= x).
// Returns NaN if Lambda <= 0; 0 if x < 0.
func (e ExpDist) CDF(x float64) float64 {
	if e.Lambda <= 0 {
		return math.NaN()
	}
	if x < 0 {
		return 0
	}
	return 1 - math.Exp(-e.Lambda*x)
}
