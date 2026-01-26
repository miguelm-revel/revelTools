package randx

type Dist interface {
	Rand() float64
	PDF(x float64) float64
	CDF(x float64) float64
}
