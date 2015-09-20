package metric

import (
	"errors"
	"math"
)

func init() {
	MetricMap["Annualized"] = PortfolioAnnualize
	MetricMap["MeanGeometric"] = MeanGeometric
	MetricMap["Variance"] = PortfolioVariance
	MetricMap["StdDev"] = PortfolioStdDev
	MetricMap["StdDev_Annualized"] = PortfolioStdDevAnnualized
	MetricMap["SharpeRatio"] = PortfolioSharpeRatio
	MetricMap["Skewness"] = PortfolioSkewness
	MetricMap["Kurtosis"] = PortfolioKurtosis
	MetricMap["AdjustedSharpeRatio"] = PortfolioAdjustedSharpeRatio
	MetricMap["MaxDrawdown"] = PortfolioMaxDrawDown
	MetricMap["AverageDrawdown"] = PortfolioAverageDrawDown
	MetricMap["AverageLength"] = PortfolioAverageLength
	MetricMap["AverageRecovery"] = PortfolioAverageRecovery
	MetricMap["UpsideFrequency"] = PortfolioUpsideFrequency
	MetricMap["DownsideDeviation2"] = PortfolioDownsideDeviation
	MetricMap["SortinoRatio"] = PortfolioSortinoRatio
	MetricMap["ProspectRatio"] = PortfolioProspectRatio
	MetricMap["UpsidePotentialRatio"] = PortfolioUpsidePotentialRatio
	MetricMap["UpsideRisk"] = PortfolioUpsideRisk
	MetricMap["KellyRatio_Full"] = PortfolioKellyRatioFull
	MetricMap["DRatio"] = PortfolioDRatio
	MetricMap["BernardoLedoitRatio"] = PortfolioBernardoLedoitRatio
	MetricMap["CalmarRatio"] = PortfolioCalmarRatio
	MetricMap["SterlingRatio"] = PortfolioSterlingRatio
	MetricMap["PainIndex"] = PortfolioPainIndex
	MetricMap["PainRatio"] = PortfolioPainRatio
	MetricMap["Kappa"] = PortfolioKappa
	MetricMap["BurkeRatio"] = PortfolioBurkeRatio
}

func PortfolioAnnualize(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioAnnualized", func() (float64, error) {
		return Vector(c.PortfolioRatio()).Annualize(Scale, true), nil
	})
}

func BenchAnnualize(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("BenchAnnualize", func() (float64, error) {
		return Vector(c.BenchRatio()).Annualize(Scale, true), nil
	})
}

func MeanGeometric(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioMeanGeometric", func() (float64, error) {
		pr := c.PortfolioRatio()
		if len(pr) <= 0 {
			return math.NaN(), errors.New("In MeanGeometric, Ra.Count() <= 0")
		}
		sum := 0.0
		for _, dr := range pr {
			sum += math.Log(dr + 1)
		}
		return math.Exp(sum/float64(len(pr))) - 1, nil
	})
}

func PortfolioVariance(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioVariance", func() (float64, error) {
		pr := c.PortfolioRatio()
		return Variance(pr)
	})
}

func BenchVariance(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("BenchVariance", func() (float64, error) {
		pr := c.BenchRatio()
		return Variance(pr)
	})
}

func Variance(array []float64) (float64, error) {
	if len(array) <= 2 {
		return math.NaN(), errors.New("In Variance, data length <= 2")
	}
	var sum, squareSum float64
	var length = float64(len(array))
	for _, dr := range array {
		sum += dr
		squareSum += dr * dr
	}
	return (squareSum - sum*sum/length) / (length - 1.0), nil
}

func PortfolioStdDev(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioStdDev", func() (float64, error) {
		value, err := PortfolioVariance(c)
		if err != nil {
			return math.NaN(), err
		}
		return math.Sqrt(value), nil
	})
}

func PortfolioStdDevAnnualized(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PorfolioStdDevAnnualized", func() (float64, error) {
		value, err := PortfolioStdDev(c)
		if err != nil {
			return math.NaN(), err
		}
		return math.Sqrt(float64(Scale)) * value, nil
	})
}

func PortfolioSharpeRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioSharpeRatio", func() (float64, error) {
		denominator, err := PortfolioStdDevAnnualized(c)
		if err != nil {
			return math.NaN(), err
		}
		periodRf := Rf / c.Period()
		pr := c.PortfolioRatio()
		length := len(pr)
		if length == 0 {
			return math.NaN(), errors.New("In SharpeRatio, data lenght == 0")
		}
		prod := 1.0
		for _, p := range pr {
			prod *= (p - periodRf) + 1
		}
		numerator := math.Pow(prod, float64(Scale)/float64(length)) - 1.0
		return numerator / denominator, nil
	})
}

func PortfolioSkewness(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioSkewness", func() (float64, error) {
		return Skewness(c.PortfolioRatio())
	})
}

// Skewness 偏度 using moment
func Skewness(array []float64) (float64, error) {
	if len(array) <= 2 {
		return math.NaN(), errors.New("In Skewness, Ra == nil || Ra.Count() <= 2")
	}
	n := float64(len(array))
	//*"moment"*, "fisher", "sample"
	varData, err := Variance(array)
	if err != nil {
		return math.NaN(), err
	}
	mean := Vector(array).Average()
	sum := 0.0
	factor := 1.0 / math.Pow(varData*(n-1.0)/n, 1.5)
	for _, p := range array {
		diff := (p - mean)
		sum += diff * diff * diff * factor
	}
	return sum / n, nil
}

func PortfolioKurtosis(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioKurtosis", func() (float64, error) {
		return Kurtosis(c.PortfolioRatio())
	})
}

// Kurtosis 峰度 using sample_excess
func Kurtosis(array []float64) (float64, error) {
	if len(array) <= 3 {
		return math.NaN(), errors.New("In Kurtosis, Ra == nil || Ra.Count() <= 3")
	}
	n := float64(len(array))
	varData, err := Variance(array)
	if err != nil {
		return math.NaN(), err
	}
	mean := Vector(array).Average()
	sum := 0.0
	factor := 1.0 / (varData * varData)
	for _, p := range array {
		diff := (p - mean)
		diff = diff * diff
		diff = diff * diff
		sum += diff * factor
	}
	return sum*n*(n+1.0)/((n-1.0)*(n-2.0)*(n-3.0)) - 3*(n-1.0)*(n-1.0)/((n-2.0)*(n-3.0)), nil
}

func PortfolioAdjustedSharpeRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioAdjustedSharpeRatio", func() (float64, error) {
		Rp, err := PortfolioAnnualize(c)
		if err != nil {
			return math.NaN(), err
		}
		Sigp, err := PortfolioStdDevAnnualized(c)
		if err != nil {
			return math.NaN(), err
		}
		SR := (Rp - Rf) / Sigp
		K, err := PortfolioKurtosis(c)
		if err != nil {
			return math.NaN(), err
		}
		S, err := PortfolioSkewness(c)
		if err != nil {
			return math.NaN(), err
		}
		var result = SR * (1.0 + (S/6.0)*SR - ((K-3.0)/24.0)*math.Pow(SR, 2.0))
		return result, nil
	})
}

func PortfolioDrawDown(c MetricCalculator) (Vector, error) {
	return c.GetOrSetVector("PortfolioDrawDown", func() (Vector, error) {
		return Vector(c.portfolio).Drawdowns(), nil
	})
}

func PortfolioMaxDrawDown(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioMaxDrawDown", func() (float64, error) {
		drawDowns, err := PortfolioDrawDown(c)
		if err != nil {
			return math.NaN(), err
		}
		if len(drawDowns) < 1 {
			return math.NaN(), errors.New("Drawdown length < 1")
		}
		max := drawDowns[0]
		for _, dd := range drawDowns {
			if dd < max {
				max = dd
			}
		}
		return -max, nil
	})
}

func PortfolioAnalysisDrawDown(c MetricCalculator) (draw, length, recovery Vector, err error) {
	const (
		PDraws    = "PortfolioDraw"
		PLength   = "PortfolioDrawLength"
		PRecovery = "PortfolioDrawVector"
	)
	var exist bool
	draw, exist = c.vectorCache[PDraws]
	if !exist {
		var drawdowns Vector
		drawdowns, err = PortfolioDrawDown(c)
		if err != nil {
			return
		}
		draw, length, recovery, err = AnalysisDrawDown(drawdowns)
		if err != nil {
			c.vectorCache[PDraws] = draw
			c.vectorCache[PLength] = length
			c.vectorCache[PRecovery] = recovery
		}
		return
	}
	length, exist = c.vectorCache[PDraws]
	if !exist {
		err = errors.New(PDraws + " absent, internal error!")
		return
	}
	recovery = c.vectorCache[PRecovery]
	if !exist {
		err = errors.New(PRecovery + " absent, internal error!")
		return
	}
	return
}

func AnalysisDrawDown(drawdowns Vector) (draw, length, recovery Vector, err error) {
	var begin []float64
	var end []float64
	var trough []float64

	if len(drawdowns) < 1 {
		err = errors.New("Cannot analysis drawDown.len < 1")
	}

	priorSign := 0
	if len(drawdowns) < 1 {
		return nil, nil, nil, errors.New("In AnalysisDrawDown, drawdown.len < 1")
	}
	if drawdowns[0] >= 0 {
		priorSign = 1
	} else {
		priorSign = 0
	}

	from := 0.0
	sofar := drawdowns[0]
	to := 0.0
	dmin := 0.0

	for i, _ := range drawdowns {
		thisSign := 0
		if drawdowns[i] < 0 {
			thisSign = 0
		} else {
			thisSign = 1
		}

		if thisSign == priorSign {
			if drawdowns[i] < sofar {
				sofar = drawdowns[i]
				dmin = float64(i)
			}
			to = float64(i) + 1.0
		} else {
			draw = append(draw, sofar)
			begin = append(begin, from)
			trough = append(trough, dmin)
			end = append(end, to)

			from = float64(i)
			sofar = drawdowns[i]
			to = float64(i) + 1
			dmin = float64(i)
			priorSign = thisSign
		}
	}

	draw = append(draw, sofar)
	begin = append(begin, from)
	trough = append(trough, dmin)
	end = append(end, to)

	length = make([]float64, len(end))
	recovery = make([]float64, len(end))
	for i, _ := range end {
		length[i] = end[i] - begin[i] + 1.0
		recovery[i] = end[i] - trough[i]
	}
	return
}

func PortfolioUpsideFrequency(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioUpsideFrequency", func() (float64, error) {
		vector := c.PortfolioRatio()
		if len(vector) < 1 {
			return math.NaN(), nil
		}
		above := 0
		periodMAR := MAR / c.Period()
		for _, value := range vector {
			if value > periodMAR {
				above++
			}
		}
		return float64(above) / float64(len(vector)), nil
	})
}
func PortfolioAverageDrawDown(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioAverageDrawDown", func() (float64, error) {
		draw, _, _, err := PortfolioAnalysisDrawDown(c)
		if err != nil {
			return math.NaN(), err
		}
		sum := 0.0
		nzrCount := 0
		for _, val := range draw {
			if val < 0 {
				nzrCount++
				sum += val
			}
		}
		return -sum / float64(nzrCount), nil
	})
}
func PortfolioAverageLength(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioAverageLength", func() (float64, error) {
		draw, length, _, err := PortfolioAnalysisDrawDown(c)
		if err != nil {
			return math.NaN(), err
		}
		nzrCount := 0
		sum := 0.0
		for i, val := range draw {
			if val < 0 {
				nzrCount++
				sum += length[i]
			}
		}
		return sum / float64(nzrCount), nil
	})
}
func PortfolioAverageRecovery(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioAverageRecovery", func() (float64, error) {
		draw, _, recovery, err := PortfolioAnalysisDrawDown(c)
		if err != nil {
			return math.NaN(), err
		}
		nzrCount := 0
		sum := 0.0
		for i, val := range draw {
			if val < 0 {
				nzrCount++
				sum += recovery[i]
			}
		}
		return sum / float64(nzrCount), nil
	})
}

func DownsideDeviation(array []float64, MAR float64) (float64, error) {
	if len(array) < 1 {
		return math.NaN(), errors.New("In DownsideDeviation, Ra.Count() < 1")
	}
	sum := 0.0
	potential := false
	isSubset := true
	downSideLength := 0
	length := 0.0
	for _, value := range array {
		if value < MAR {
			downSideLength++
			diff := MAR - value
			if potential {
				sum += diff
			} else {
				sum += diff * diff
			}
		}
	}
	if isSubset {
		length = float64(downSideLength)
	} else {
		length = float64(len(array))
	}
	sum = sum / length
	if !potential {
		sum = math.Sqrt(sum)
	}
	return sum, nil
}

func PortfolioDownsideDeviation(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioDownsideDeviation", func() (float64, error) {
		return DownsideDeviation(c.PortfolioRatio(), MAR/c.Period())
	})
}

func PortfolioSortinoRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioSortinoRatio", func() (float64, error) {
		dd, err := PortfolioDownsideDeviation(c)
		if err != nil {
			return math.NaN(), err
		}
		return (Vector(c.PortfolioRatio()).Average() - MAR/c.Period()) / dd, nil
	})
}

func PortfolioProspectRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioProspectRatio", func() (float64, error) {
		dd, err := PortfolioDownsideDeviation(c)
		if err != nil {
			return math.NaN(), err
		}
		var posSum, negSum float64
		pr := c.PortfolioRatio()
		for _, value := range pr {
			if value > 0 {
				posSum += value
			} else {
				negSum += value
			}
		}

		return ((posSum+2.25*negSum)/float64(len(pr)) - MAR/c.Period()) / dd, nil
	})
}

func PortfolioUpsidePotentialRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioUpsidePotentialRatio", func() (float64, error) {
		pr := c.PortfolioRatio()
		sum := 0.0
		upsideCount := 0
		periodMAR := MAR / c.Period()
		for _, value := range pr {
			if value > periodMAR {
				upsideCount++
				diff := value - periodMAR
				sum += diff
			}
		}
		isSubset := true
		if isSubset {
			sum = sum / float64(upsideCount)
		} else {
			sum = sum / float64(len(pr))
		}
		dd2Data, err := PortfolioDownsideDeviation(c)
		if err != nil {
			return math.NaN(), err
		}
		return sum / dd2Data, nil
	})
}

func PortfolioUpsideRisk(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioUpsideRisk", func() (float64, error) {
		pr := c.PortfolioRatio()
		sum := 0.0
		upsideCount := 0
		periodMAR := MAR / c.Period()
		for _, value := range pr {
			if value > periodMAR {
				upsideCount++
				diff := value - periodMAR
				sum += diff * diff
			}
		}
		isSubset := true
		if isSubset {
			sum = sum / float64(upsideCount)
		} else {
			sum = sum / float64(len(pr))
		}
		return math.Sqrt(sum), nil
	})
}
func PortfolioKellyRatioFull(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioKellyRatioFull", func() (float64, error) {
		varData, err := PortfolioVariance(c)
		if err != nil {
			return math.NaN(), err
		}
		pr := c.PortfolioRatio()
		return (Vector(pr).Average() - Rf/c.Period()) / varData, nil
	})
}

func PortfolioDRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioDRatio", func() (float64, error) {
		var negSum, posSum float64
		var negCount, posCount int
		pr := c.PortfolioRatio()
		for _, value := range pr {
			if value > 0 {
				posSum += value
				posCount++
			} else if value < 0 {
				negSum += value
				negCount++
			}
		}
		return -(negSum * float64(negCount)) / (posSum * float64(posCount)), nil
	})

}

func PortfolioBernardoLedoitRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioBernardoLedoitRatio", func() (float64, error) {
		var negSum, posSum float64
		pr := c.PortfolioRatio()
		for _, value := range pr {
			if value > 0 {
				posSum += value
			} else if value < 0 {
				negSum += value
			}
		}
		return -posSum / negSum, nil
	})
}

func PortfolioCalmarRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioCalmarRatio", func() (float64, error) {
		ar, err := PortfolioAnnualize(c)
		if err != nil {
			return math.NaN(), err
		}
		md, err := PortfolioMaxDrawDown(c)
		if err != nil {
			return math.NaN(), err
		}
		return ar / md, nil
	})
}

func PortfolioSterlingRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioSterlingRatio", func() (float64, error) {
		excess := 0.1
		ar, err := PortfolioAnnualize(c)
		if err != nil {
			return math.NaN(), err
		}
		md, err := PortfolioMaxDrawDown(c)
		if err != nil {
			return math.NaN(), err
		}
		dd := math.Abs(md + excess)
		if dd < 0.0000001 {
			return math.NaN(), errors.New("In SterlingRatio, draw_down == 0.0")
		}
		return ar / dd, nil
	})
}

func PortfolioPainIndex(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioPainIndex", func() (float64, error) {
		vector, err := PortfolioDrawDown(c)
		if err != nil {
			return math.NaN(), err
		}
		total := 0.0
		for _, val := range vector {
			total += math.Abs(val)
		}
		if len(vector) == 0 {
			return math.NaN(), errors.New("In PainIndex, len(data) == 0")
		}
		return total / float64(len(vector)), nil
	})
}

func PortfolioPainRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioPainRatio", func() (float64, error) {
		a, err := PortfolioAnnualize(c)
		if err != nil {
			return math.NaN(), err
		}
		pi, err := PortfolioPainIndex(c)
		if err != nil {
			return math.NaN(), err
		}
		return (a - Rf) / pi, nil
	})
}

func Kappa(rRatio []float64, MAR float64, l float64) float64 {
	sum := 0.0
	for _, r := range rRatio {
		if MAR > r {
			sum += math.Pow(MAR-r, l)
		}
	}
	n := float64(len(rRatio))
	m := Vector(rRatio).Average()
	temp := sum / n
	return (m - MAR) / math.Pow(temp, 1.0/float64(l))
}

func PortfolioKappa(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioKappa", func() (float64, error) {
		return Kappa(c.PortfolioRatio(), MAR/c.Period(), 1.0), nil
	})
}
func PortfolioBurkeRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("PortfolioBurkeRatio", func() (float64, error) {
		ra, err := PortfolioAnnualize(c)
		if err != nil {
			return math.NaN(), err
		}
		denominator := ra - Rf
		peak := 0
		portfolio := c.portfolio
		inDrawDown := false
		sum := 0.0
		for i := 1; i < len(portfolio); i++ {
			if portfolio[i] < portfolio[i-1] {
				if !inDrawDown {
					peak = i - 1
					inDrawDown = true
				}
			} else {
				if inDrawDown {
					draw := portfolio[i-1]/portfolio[peak] - 1.0
					sum += draw * draw
					inDrawDown = false
				}
			}
		}
		if inDrawDown {
			draw := portfolio[len(portfolio)-1]/portfolio[peak] - 1.0
			sum += draw * draw
			inDrawDown = false
		}
		modified := true
		if modified {
			sum /= float64(len(portfolio))
		}
		return denominator / math.Sqrt(sum), nil

	})
}
