package index

import (
	"github.com/miguelm-revel/revelTools/collections"
	"math"
)

type Vector interface {
	Dist(Vector) float64
}
type vec328 [328]float64

func (v vec328) Dist(vector Vector) float64 {
	var dot, n1, n2 float64
	vec := vector.(vec328)

	for v1, v2 := range collections.ZipSlice(v[:], vec[:]) {
		dot += v1 * v2
		n1 += v1 * v1
		n2 += v2 * v2
	}

	den := math.Sqrt(n1) * math.Sqrt(n2)
	if den == 0 {
		return 0
	}
	return 1 - dot/den
}

type vec4 [4]float64

func (v vec4) Dist(vector Vector) float64 {
	var dot, n1, n2 float64
	vec := vector.(vec4)

	for v1, v2 := range collections.ZipSlice(v[:], vec[:]) {
		dot += v1 * v2
		n1 += v1 * v1
		n2 += v2 * v2
	}

	den := math.Sqrt(n1) * math.Sqrt(n2)
	if den == 0 {
		return 0
	}
	return 1 - dot/den
}
