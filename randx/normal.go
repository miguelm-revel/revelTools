package randx

import (
	"math"
	"math/rand/v2"
)

// NormalDist represents a Normal(Mu, Sigma) distribution.
// Mu is the mean and Sigma is the standard deviation (must be > 0).
type NormalDist struct {
	Mu, Sigma float64
}

// Rand returns a random sample using the Box-Muller transform.
// Returns NaN if Sigma <= 0.
func (n NormalDist) Rand() float64 {
	if n.Sigma <= 0 {
		return math.NaN()
	}

	u1 := rand.Float64()
	if u1 == 0 {
		u1 = math.SmallestNonzeroFloat64
	}
	u2 := rand.Float64()

	r := math.Sqrt(-2.0 * math.Log(u1))
	theta := 2.0 * math.Pi * u2

	z0 := r * math.Cos(theta)

	return n.Mu + n.Sigma*z0
}

// PDF returns the probability density at x.
// Returns NaN if Sigma <= 0.
func (n NormalDist) PDF(x float64) float64 {
	if n.Sigma <= 0 {
		return math.NaN()
	}
	z := (x - n.Mu) / n.Sigma
	return (1.0 / (n.Sigma * math.Sqrt(2.0*math.Pi))) * math.Exp(-0.5*z*z)
}

// CDF returns the cumulative probability P(X <= x) using the error function.
// Returns NaN if Sigma <= 0.
func (n NormalDist) CDF(x float64) float64 {
	if n.Sigma <= 0 {
		return math.NaN()
	}
	z := (x - n.Mu) / (n.Sigma * math.Sqrt2)
	return 0.5 * (1.0 + math.Erf(z))
}
