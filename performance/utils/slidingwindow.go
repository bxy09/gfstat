package utils

import (
	"errors"
	"math"
)

type SlidingWindow struct {
	window    []float64
	index     int
	count     int
	container int
}

func NewSlidingWindow(size int) (*SlidingWindow, error) {
	if size <= 0 {
		return nil, errors.New("In NewSlidingWindow, the length invalid size <=0")
	}
	return &SlidingWindow{window: make([]float64, size), index: 0, count: 0, container: size}, nil
}

func (w *SlidingWindow) Data() []float64 {
	ret := make([]float64, w.count)
	index := w.index % w.count
	for i := 0; i < w.count; i++ {
		ret[i] = w.window[index]
		index++
		if index >= w.count {
			index = 0
		}
	}
	return ret
}

func (w *SlidingWindow) Add(value float64) {
	if math.IsNaN(value) {
		value = 0.0
	}
	w.window[w.index] = value
	w.index++
	if w.index >= len(w.window) {
		w.index = 0
	}
	w.count++
	w.count = int(math.Min(float64(len(w.window)), float64(w.count)))
}

func (w *SlidingWindow) Average() (average float64) {
	average = 0.0
	for i := 0; i < w.count; i++ {
		average += w.window[i]
	}
	average /= float64(w.count)
	return
}

func (w *SlidingWindow) StdDev() (stddev float64) {
	stddev = 0.0
	ave := w.Average()
	for i := 0; i < w.count; i++ {
		stddev += w.window[i] * w.window[i]
	}
	stddev /= float64(w.count)
	stddev -= ave * ave
	stddev = math.Sqrt(stddev)
	return
}

func (w *SlidingWindow) Sum() (total float64) {
	total = 0.0
	for i := 0; i < w.count; i++ {
		total += w.window[i]
	}
	return
}

func (w *SlidingWindow) Max() (value float64, position int) {
	value = math.Inf(-1)
	position = -1
	for i := 0; i < w.count; i++ {
		if w.window[i] > value {
			value = w.window[i]
			position = i
		}
	}
	position -= w.index
	if position < 0 {
		position += w.count
	}
	return
}

func (w *SlidingWindow) Min() (value float64, position int) {
	value = math.Inf(1)
	position = -1
	for i := 0; i < w.count; i++ {
		if w.window[i] < value {
			value = w.window[i]
			position = i
		}
	}
	position -= w.index
	if position < 0 {
		position += w.count
	}
	return
}

func (w *SlidingWindow) Count() int {
	return w.count
}

func (w *SlidingWindow) DiffAbsAve() (diffAbsAverage float64) {
	diffAbsAverage = 0.0
	ave := w.Average()
	for i := 0; i < w.count; i++ {
		diffAbsAverage += math.Abs(w.window[i] - ave)
	}
	diffAbsAverage /= float64(w.count)
	return
}

func (w *SlidingWindow) RefData(ref int) (value float64) {
	index := w.index - ref
	if index < 0 {
		index += w.count
	}
	value = w.window[index]
	return
}

func (w *SlidingWindow) First() (value float64) {
	index := w.index % w.count
	value = w.window[index]
	return
}
