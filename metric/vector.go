package metric

import (
	"math"
)

type Vector []float64

func (v Vector) ImplV(op func(float64) float64) Vector {
	for i, value := range v {
		v[i] = op(value)
	}
	return v
}

func (v1 Vector) ImplVectorV(v2 Vector, op func(float64, float64) (float64, float64)) Vector {
	length := len(v1)
	if len(v2) < length {
		length = len(v2)
	}
	for i := 0; i < length; i++ {
		v1[i], v2[i] = op(v1[i], v2[i])
	}
	return v1
}

func (v1 Vector) ReturnRatio(method string) Vector {
	if len(v1) == 0 {
		return nil
	}
	v2 := make([]float64, len(v1))
	switch method {
	case "discrete":
		v2[0] = 0
		for i := 1; i < len(v1); i++ {
			if v1[i-1] != 0.0 {
				v2[i] = (v1[i]/v1[i-1] - 1.0)
			} else {
				v2[i] = 0
			}
		}
	default:
		return nil
	}
	return v2
}

func (v1 Vector) AddScalarV(s float64) Vector {
	return v1.ImplV(func(a float64) float64 { return a + s })
}

func (v1 Vector) AddVectorV(v2 Vector) Vector {
	return v1.ImplVectorV(v2, func(a, b float64) (float64, float64) { return a + b, b })
}

func (v1 Vector) AccumulateSum() float64 {
	sum := 0.0
	for _, p := range v1 {
		sum += p
	}
	return sum
}

func (v1 Vector) AccumulateMul() float64 {
	prod := 1.0
	for _, p := range v1 {
		prod *= p
	}
	return prod
}

func (v1 Vector) Average() float64 {
	if len(v1) == 0 {
		return math.NaN()
	}
	return v1.AccumulateSum() / float64(len(v1))
}

func (v1 Vector) Annualize(scale float64, geometric bool) float64 {
	n := len(v1)
	if n == 0 {
		return math.NaN()
	}
	if geometric {
		prod := 1.0
		for _, p := range v1 {
			prod *= p + 1
		}
		return math.Pow(prod, float64(scale)/float64(n)) - 1.0
	} else {
		return v1.Average() * scale
	}
}

func (v1 Vector) Drawdowns() Vector {
	if len(v1) < 1 {
		return nil
	}
	result := make([]float64, len(v1))
	curMax := v1[0]
	for i, r := range v1 {
		if r > curMax {
			curMax = r
		} else {
			result[i] = r/curMax - 1.0
		}
	}
	return result
}
