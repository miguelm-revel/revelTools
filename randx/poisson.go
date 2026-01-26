package randx

import (
	"math"
	"math/rand/v2"
)

type PoissonDist struct {
	Lambda float64
}

func (p PoissonDist) Rand() float64 {
	if p.Lambda < 0 {
		return math.NaN()
	}
	if p.Lambda == 0 {
		return 0
	}
	if p.Lambda < 30 {
		L := math.Exp(-p.Lambda)
		k := 0
		prod := 1.0
		for prod > L {
			k++
			prod *= rand.Float64()
		}
		return float64(k - 1)
	}
	return float64(poissonPTRS(p.Lambda))
}

func (p PoissonDist) PDF(x float64) float64 {
	if p.Lambda < 0 {
		return math.NaN()
	}
	k := int(math.Round(x))
	if float64(k) != x {
		return 0
	}
	if k < 0 {
		return 0
	}
	return math.Exp(float64(k)*math.Log(p.Lambda) - p.Lambda - logFactorial(k))
}

func (p PoissonDist) CDF(x float64) float64 {
	if p.Lambda < 0 {
		return math.NaN()
	}
	k := int(math.Floor(x))
	if k < 0 {
		return 0
	}
	return regLowerGamma(float64(k+1), p.Lambda)
}
