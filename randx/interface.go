package randx

// Dist is implemented by any probability distribution that can generate random samples
// and evaluate its probability density/mass function (PDF) and cumulative distribution function (CDF).
type Dist interface {
	Rand() float64
	PDF(x float64) float64
	CDF(x float64) float64
}
