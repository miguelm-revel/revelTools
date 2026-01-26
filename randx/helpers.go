package randx

import (
	"math"
	"math/rand/v2"
)

/* =============================
   Helpers: combinatoria, gamma, etc.
=============================*/

func logGamma(x float64) float64 {
	lg, _ := math.Lgamma(x)
	return lg
}

func logFactorial(n int) float64 {
	return logGamma(float64(n) + 1.0)
}

func logChoose(n, k int) float64 {
	if k < 0 || k > n {
		return math.Inf(-1)
	}
	return logFactorial(n) - logFactorial(k) - logFactorial(n-k)
}

/* -----------------------------
   Gamma RNG: Marsaglia–Tsang
   Devuelve Gamma(shape, scale=1)
------------------------------*/

func gammaRand(shape float64) float64 {
	if shape <= 0 {
		return math.NaN()
	}

	if shape < 1.0 {
		u := rand.Float64()
		if u == 0 {
			u = math.SmallestNonzeroFloat64
		}
		return gammaRand(shape+1.0) * math.Pow(u, 1.0/shape)
	}

	d := shape - 1.0/3.0
	c := 1.0 / math.Sqrt(9.0*d)

	nd := NormalDist{Mu: 0, Sigma: 1}

	for {
		x := nd.Rand()
		v := 1.0 + c*x
		if v <= 0 {
			continue
		}
		v = v * v * v
		u := rand.Float64()
		if u < 1.0-0.0331*(x*x)*(x*x) {
			return d * v
		}
		if math.Log(u) < 0.5*x*x+d*(1.0-v+math.Log(v)) {
			return d * v
		}
	}
}

/* -----------------------------
   Regularized lower incomplete gamma:
   P(a,x) = γ(a,x)/Γ(a)
   Usada para CDF de Chi² y Poisson.
------------------------------*/

func regLowerGamma(a, x float64) float64 {
	if a <= 0 || x < 0 {
		return math.NaN()
	}
	if x == 0 {
		return 0
	}

	if x < a+1.0 {
		return gammaSeries(a, x)
	}
	return 1.0 - gammaContFrac(a, x)
}

func gammaSeries(a, x float64) float64 {
	const itmax = 200
	const eps = 3e-14

	sum := 1.0 / a
	del := sum
	ap := a

	for n := 1; n <= itmax; n++ {
		ap += 1.0
		del *= x / ap
		sum += del
		if math.Abs(del) < math.Abs(sum)*eps {
			break
		}
	}

	return sum * math.Exp(-x+a*math.Log(x)-logGamma(a))
}

func gammaContFrac(a, x float64) float64 {
	const itmax = 200
	const eps = 3e-14
	const fpmin = 1e-300

	b := x + 1.0 - a
	c := 1.0 / fpmin
	d := 1.0 / b
	h := d

	for i := 1; i <= itmax; i++ {
		an := -float64(i) * (float64(i) - a)
		b += 2.0
		d = an*d + b
		if math.Abs(d) < fpmin {
			d = fpmin
		}
		c = b + an/c
		if math.Abs(c) < fpmin {
			c = fpmin
		}
		d = 1.0 / d
		del := d * c
		h *= del
		if math.Abs(del-1.0) < eps {
			break
		}
	}

	return math.Exp(-x+a*math.Log(x)-logGamma(a)) * h
}

/* -----------------------------
   Poisson PTRS (Hörmann, 1993)
------------------------------*/

func poissonPTRS(lambda float64) int {
	sqrtL := math.Sqrt(lambda)
	logL := math.Log(lambda)

	b := 0.931 + 2.53*sqrtL
	a := -0.059 + 0.02483*b
	invAlpha := 1.1239 + 1.1328/(b-3.4)
	vR := 0.9277 - 3.6224/(b-2.0)

	for {
		u := rand.Float64() - 0.5
		v := rand.Float64()

		us := 0.5 - math.Abs(u)
		k := int(math.Floor((2*a/us+b)*u + lambda + 0.43))
		if k < 0 {
			continue
		}

		if us >= 0.07 && v <= vR {
			return k
		}

		lhs := math.Log(v * invAlpha / (a/(us*us) + b))
		rhs := float64(k)*logL - lambda - logFactorial(k)
		if lhs <= rhs {
			return k
		}
	}
}
