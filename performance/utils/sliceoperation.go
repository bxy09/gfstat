package utils

import (
	"math"
)

func Prod(values *SlidingWindow) (float64, error) {
	result := 1.0
	for i := 0; i < values.Count(); i++ {
		result *= values.Data()[i]
	}
	return result, nil
}

func Add(x float64, values *SlidingWindow) (*SlidingWindow, error) {
	result, err := NewSlidingWindow(values.Count())
	if err != nil {
		return nil, err
	}
	for i := 0; i < values.Count(); i++ {
		result.Add(values.Data()[i] + x)
	}
	return result, nil
}

func Add2(values1 *SlidingWindow, values2 *SlidingWindow) (*SlidingWindow, error) {
	result, err := NewSlidingWindow(values1.Count())
	if err != nil {
		return nil, err
	}
	for i := 0; i < values1.Count(); i++ {
		result.Add(values1.Data()[i] + values2.Data()[i])
	}
	return result, nil
}

func Sub(values1 *SlidingWindow, values2 *SlidingWindow) (*SlidingWindow, error) {
	result, err := NewSlidingWindow(values1.Count())
	if err != nil {
		return nil, err
	}
	for i := 0; i < values1.Count(); i++ {
		result.Add(values1.Data()[i] - values2.Data()[i])
	}
	return result, nil
}

func Multi(x float64, values *SlidingWindow) (*SlidingWindow, error) {
	result, err := NewSlidingWindow(values.Count())
	if err != nil {
		return nil, err
	}
	for i := 0; i < values.Count(); i++ {
		result.Add(values.Data()[i] * x)
	}
	return result, nil
}

func Power(values *SlidingWindow, x float64) (*SlidingWindow, error) {
	result, err := NewSlidingWindow(values.Count())
	if err != nil {
		return nil, err
	}
	for i := 0; i < values.Count(); i++ {
		result.Add(math.Pow(values.Data()[i], x))
	}
	return result, nil
}

func CreateList(value float64, length int) (*SlidingWindow, error) {
	result, err := NewSlidingWindow(length)
	if err != nil {
		return nil, err
	}
	for i := 0; i < length; i++ {
		result.Add(value)
	}
	return result, nil
}

func Negative(values *SlidingWindow) (*SlidingWindow, error) {
	result, err := NewSlidingWindow(values.Count())
	if err != nil {
		return nil, err
	}
	for i := 0; i < values.Count(); i++ {
		result.Add(-values.Data()[i])
	}
	return result, nil
}

func ElementMulti(values1 *SlidingWindow, values2 *SlidingWindow) (*SlidingWindow, error) {
	result, err := NewSlidingWindow(values1.Count())
	if err != nil {
		return nil, err
	}
	for i := 0; i < values1.Count(); i++ {
		result.Add(values1.Data()[i] * values2.Data()[i])
	}
	return result, nil
}

func Log(values *SlidingWindow) (*SlidingWindow, error) {
	result, err := NewSlidingWindow(values.Count())
	if err != nil {
		return nil, err
	}
	for i := 0; i < values.Count(); i++ {
		result.Add(math.Log(values.Data()[i]))
	}
	return result, nil
}

func Abs(values *SlidingWindow) (*SlidingWindow, error) {
	result, err := NewSlidingWindow(values.Count())
	if err != nil {
		return nil, err
	}
	for i := 0; i < values.Count(); i++ {
		result.Add(math.Abs(values.Data()[i]))
	}
	return result, nil
}

// find out positive and negtive values
func PosNegValues(values *SlidingWindow) (positivevalues, negativevalues *SlidingWindow, err error) {
	positivevalues, err = NewSlidingWindow(values.Count())
	if err != nil {
		return
	}
	negativevalues, err = NewSlidingWindow(values.Count())
	if err != nil {
		return
	}

	for i := 0; i < values.Count(); i++ {
		if values.Data()[i] > 0 {
			positivevalues.Add(values.Data()[i])
		} else {
			negativevalues.Add(values.Data()[i])
		}
	}
	return
}

func AboveValue(Ra *SlidingWindow, v float64) (*SlidingWindow, error) {
	r, err := NewSlidingWindow(Ra.Count())
	if err != nil {
		return nil, err
	}
	for i := 0; i < Ra.Count(); i++ {
		if Ra.Data()[i] > v {
			r.Add(Ra.Data()[i])
		}
	}
	return r, nil
}
