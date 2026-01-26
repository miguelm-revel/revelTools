package randx

import "math"

type Chi2Dist struct {
	K float64
}

func (c Chi2Dist) Rand() float64 {
	if c.K <= 0 {
		return math.NaN()
	}
	shape := c.K / 2.0
	scale := 2.0
	return scale * gammaRand(shape)
}

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

func (c Chi2Dist) CDF(x float64) float64 {
	if c.K <= 0 {
		return math.NaN()
	}
	if x < 0 {
		return 0
	}
	return regLowerGamma(c.K/2.0, x/2.0)
}
