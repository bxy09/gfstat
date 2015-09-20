package performance

import (
	"errors"
	"github.com/bxy09/gfstat/performance/utils"
	"math"
)

/// <param name="R"></param>
/// <param name="scale"></param>
/// <param name="geometric"></param>
/// <returns></returns>
func Y2Scale(R float64, scale int, geometric bool) (float64, error) {
	if geometric {
		return math.Pow(1.0+R, 1.0/float64(scale)-1.0), nil
	} else {
		return R / float64(scale), nil
	}
}

/// <param name="returns"></param>
/// <param name="scale"></param>
/// <param name="geometric"></param>
/// <returns></returns>
func Annualized(returns *utils.SlidingWindow, scale float64, geometric bool) (float64, error) {
	if returns == nil {
		return math.NaN(), errors.New("Returns Utils Sliding Window is nil")
	}
	if returns.Count() == 0 {
		return math.NaN(), errors.New("Returns Windows content is Zero")
	}
	n := returns.Count()
	if geometric {
		add_Sliding, err := utils.Add(1.0, returns)
		if err != nil {
			return math.NaN(), err
		}
		prod_Data, err := utils.Prod(add_Sliding)
		if err != nil {
			return math.NaN(), err
		}
		return math.Pow(prod_Data, float64(scale)/float64(n)) - 1.0, nil
	} else {
		return returns.Average() * float64(scale), nil
	}
}

/// <param name="prices"></param>
/// <param name="method"></param>
/// <returns></returns>
func Calculate(prices *utils.SlidingWindow, method string) (*utils.SlidingWindow, error) {
	if prices == nil {
		return nil, errors.New("Prices Utils Sliding Window is nil")
	}
	if prices.Count() == 0 {
		return nil, errors.New("Returns Windows content is Zero")
	}

	lastPrice := prices.First()
	returns, err := utils.NewSlidingWindow(prices.Count())
	if err != nil {
		return nil, errors.New("create a Sliding Window is Error !!")
	}

	switch method {
	case "simple":
	case "discrete":
		for i := 0; i < prices.Count(); i++ {
			price := prices.Data()[i]
			if lastPrice != 0.0 {
				returns.Add(price/lastPrice - 1.0)
			} else {
				returns.Add(0.0)
			}
			lastPrice = price
		}
	case "compound":
	case "log":
		for i := 0; i < prices.Count(); i++ {
			price := prices.Data()[i]
			if lastPrice != 0.0 {
				returns.Add(math.Log(price / lastPrice))
			} else {
				returns.Add(0.0)
			}
			lastPrice = price
		}
	default:
		return nil, errors.New("The input Method is nil !!!")
	}
	return returns, nil
}

/// <param name="returns"></param>
/// <param name="Rf"></param>
/// <returns></returns>
func Excess2(returns *utils.SlidingWindow, Rf float64) (*utils.SlidingWindow, error) {
	return utils.Add(-Rf, returns)
}

func Excess(returns, Rf *utils.SlidingWindow) (*utils.SlidingWindow, error) {
	result, err := utils.NewSlidingWindow(returns.Count())
	if err != nil {
		return nil, err
	}
	for i := 0; i < returns.Count(); i++ {
		result.Add(returns.Data()[i] - Rf.Data()[i])
	}
	return result, nil
}

/// <param name="Ra"></param>
/// <param name="Rb"></param>
/// <returns></returns>
func Relative(Ra, Rb *utils.SlidingWindow) (*utils.SlidingWindow, error) {
	res4Ra := 1.0
	res4Rb := 1.0
	result, err := utils.NewSlidingWindow(Ra.Count())
	if err != nil {
		return nil, err
	}
	for i := 0; i < Ra.Count(); i++ {
		res4Ra = res4Ra * (1 + Ra.Data()[i])
		res4Rb = res4Rb * (1 + Rb.Data()[i])
		result.Add(res4Ra / res4Rb)
	}
	return result, nil
}

/// <param name="returns"></param>
/// <returns></returns>
func Centered(returns *utils.SlidingWindow) (*utils.SlidingWindow, error) {
	if returns == nil {
		return nil, errors.New("Centered Sliding window is nil")
	}
	if returns.Count() == 0 {
		return nil, errors.New("Centered Count is Zero !!!")
	}
	return utils.Add(-returns.Average(), returns)
}

/// <param name="returns"></param>
/// <param name="geometric"></param>
/// <returns></returns>
func Cumulative(returns *utils.SlidingWindow, geometric bool) (float64, error) {
	if returns == nil {
		return math.NaN(), errors.New("Cumulative Sliding window is Nil !!!")
	}
	if returns.Count() == 0 {
		return math.NaN(), errors.New("Cumulative Count == 0 !!")
	}
	if !geometric {
		return (returns.Sum()), nil
	} else {
		add_data, err := utils.Add(1.0, returns)
		if err != nil {
			return math.NaN(), err
		}
		prod_data, err := utils.Prod(add_data)
		if err != nil {
			return math.NaN(), err
		}
		return (prod_data - 1.0), nil
	}
}
