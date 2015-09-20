package metric

import (
	"errors"
	"time"
)

func reorganizeInputPrice(date []time.Time, MinutesPrice []float64) ([]float64, error) {
	if len(date) != len(MinutesPrice) {
		return nil, errors.New("The length of date and Minutes Price is not Equal !!!")
	}
	Price := make([]float64, 0)
	for i := 1; i < len(date); i++ {
		diff := date[i].Sub(date[i-1])
		if diff.Hours() > 10.0 {
			Price = append(Price, MinutesPrice[i-1])
		}
	}
	Price = append(Price, MinutesPrice[len(date)-1])
	return Price, nil
}
