package randx

import "math"

// Chi2Dist represents a Chi-squared distribution with K degrees of freedom.
// K must be positive.
type Chi2Dist struct {
	K float64
}

// Rand returns a random sample via the Gamma(K/2, 2) parameterization.
// Returns NaN if K <= 0.
func (c Chi2Dist) Rand() float64 {
	if c.K <= 0 {
		return math.NaN()
	}
	shape := c.K / 2.0
	scale := 2.0
	return scale * gammaRand(shape)
}

// PDF returns the probability density at x.
// Returns NaN if K <= 0; 0 if x < 0.
func (c Chi2Dist) PDF(x float64) float64 {
	if c.K <= 0 {
		return math.NaN()
	}
	if x < 0 {
		return 0
	}
	a := c.K / 2.0
	logf := -(a*math.Log(2.0) + logGamma(a)) + (a-1.0)*math.Log(x) - x/2.0
	return math.Exp(logf)
}

// CDF returns the cumulative probability P(X <= x) via the regularized lower incomplete gamma function.
// Returns NaN if K <= 0; 0 if x < 0.
func (c Chi2Dist) CDF(x float64) float64 {
	if c.K <= 0 {
		return math.NaN()
	}
	if x < 0 {
		return 0
	}
	return regLowerGamma(c.K/2.0, x/2.0)
}
