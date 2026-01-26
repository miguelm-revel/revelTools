package randx

import (
	"math"
	"math/rand/v2"
)

type BinomDist struct {
	N int
	P float64
}

func (b BinomDist) Rand() float64 {
	if b.N < 0 || b.P < 0 || b.P > 1 {
		return math.NaN()
	}
	k := 0
	for i := 0; i < b.N; i++ {
		if rand.Float64() < b.P {
			k++
		}
	}
	return float64(k)
}

func (b BinomDist) PDF(x float64) float64 {
	if b.N < 0 || b.P < 0 || b.P > 1 {
		return math.NaN()
	}
	k := int(math.Round(x))
	if float64(k) != x {
		return 0
	}
	if k < 0 || k > b.N {
		return 0
	}

	if b.P == 0 {
		if k == 0 {
			return 1
		}
		return 0
	}
	if b.P == 1 {
		if k == b.N {
			return 1
		}
		return 0
	}

	logC := logChoose(b.N, k)
	return math.Exp(logC + float64(k)*math.Log(b.P) + float64(b.N-k)*math.Log(1.0-b.P))
}

func (b BinomDist) CDF(x float64) float64 {
	if b.N < 0 || b.P < 0 || b.P > 1 {
		return math.NaN()
	}
	k := int(math.Floor(x))
	if k < 0 {
		return 0
	}
	if k >= b.N {
		return 1
	}
	sum := 0.0
	for i := 0; i <= k; i++ {
		sum += b.PDF(float64(i))
	}
	if sum < 0 {
		return 0
	}
	if sum > 1 {
		return 1
	}
	return sum
}
