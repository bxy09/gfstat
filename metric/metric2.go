package metric

import (
	"errors"
	"github.com/GaryBoone/GoStats/stats"
	"math"
)

func init() {
	MetricMap["ActivePremium"] = ActivePremium
	MetricMap["TrackingError"] = TrackingError
	MetricMap["InformationRatio"] = InformationRatio
	MetricMap["MSquared"] = MSquared
	MetricMap["JensenAlpha2"] = JensenAlpha
	MetricMap["TreynorRatio"] = TreynorRatio
	MetricMap["AppraisalRatio"] = AppraisalRatio
	MetricMap["SpecificRisk"] = SpecificRisk
	MetricMap["SystematicRisk"] = SystematicRisk
	MetricMap["TotalRisk"] = TotalRisk
}

func ActivePremium(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("ActivePremium", func() (float64, error) {
		pa, err := PortfolioAnnualize(c)
		if err != nil {
			return math.NaN(), err
		}
		ba, err := BenchAnnualize(c)
		if err != nil {
			return math.NaN(), err
		}
		return pa - ba, nil
	})
}

func TrackingError(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("TrackingError", func() (float64, error) {
		pr := c.PortfolioRatio()
		br := c.BenchRatio()
		if len(pr) != len(br) {
			return math.NaN(), errors.New("In TrackingError, len(portfolio) != len(bench)")
		}
		var sum, squareSum float64
		var length = float64(len(pr))
		for i := 0; i < len(pr); i++ {
			diff := pr[i] - br[i]
			sum += diff
			squareSum += diff * diff
		}
		variance := (squareSum - sum*sum/length) / (length - 1.0)
		return math.Sqrt(variance * Scale), nil
	})
}

func InformationRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("InformationRatio", func() (float64, error) {
		ap, err := ActivePremium(c)
		if err != nil {
			return math.NaN(), err
		}
		te, err := TrackingError(c)
		if err != nil {
			return math.NaN(), err
		}
		return ap / te, nil
	})

}

func MSquared(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("MSquared", func() (float64, error) {
		pa, err := PortfolioAnnualize(c)
		if err != nil {
			return math.NaN(), err
		}
		pVar, err := PortfolioVariance(c)
		if err != nil {
			return math.NaN(), err
		}
		bVar, err := BenchVariance(c)
		if err != nil {
			return math.NaN(), err
		}
		var n = float64(len(c.PortfolioRatio()))
		sigp := math.Sqrt(pVar*(n-1)/n) * math.Sqrt(Scale)
		sigm := math.Sqrt(bVar*(n-1)/n) * math.Sqrt(Scale)
		return (pa-Rf)*sigm/sigp + Rf, nil
	})
}

func AlphaBeta(c MetricCalculator) (alpha, beta float64, err error) {
	const (
		AlphaKey = "Alpha"
		BetaKey  = "Beta"
	)
	var exist bool
	alpha, exist = c.scalarCache[AlphaKey]
	if !exist {
		rb := c.BenchRatio()
		rp := c.PortfolioRatio()
		xRBench := make([]float64, len(rb))
		xRPortfolio := make([]float64, len(rp))
		periodRf := Rf / c.Period()
		for i, p := range rb {
			xRBench[i] = p - periodRf
		}
		for i, p := range rp {
			xRPortfolio[i] = p - periodRf
		}
		beta, alpha, _, _, _, _ = stats.LinearRegression(rb, rp)
		c.scalarCache[AlphaKey] = alpha
		c.scalarCache[BetaKey] = beta
		return
	}
	beta, exist = c.scalarCache[BetaKey]
	if !exist {
		err = errors.New(BetaKey + " absent, internal error!")
		return
	}
	return
}

func JensenAlpha(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("JensenAlpha", func() (float64, error) {
		_, beta, err := AlphaBeta(c)
		if err != nil {
			return math.NaN(), err
		}
		Rpa, err := PortfolioAnnualize(c)
		if err != nil {
			return math.NaN(), err
		}
		Rpb, err := BenchAnnualize(c)
		if err != nil {
			return math.NaN(), err
		}
		return Rpa - Rf - beta*(Rpb-Rf), nil
	})
}

func TreynorRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("TreynorRatio", func() (float64, error) {
		periodRf := Rf / c.Period()
		rp := c.PortfolioRatio()
		localRP := make([]float64, len(rp))
		copy(localRP, rp)
		tr := Vector(localRP).AddScalarV(-periodRf).Annualize(Scale, true)
		_, beta, err := AlphaBeta(c)
		if err != nil {
			return math.NaN(), err
		}
		return tr / beta, nil
	})
}

func AppraisalRatio(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("AppraisalRatio", func() (float64, error) {
		js, err := JensenAlpha(c)
		if err != nil {
			return math.NaN(), err
		}
		var result = 0.0
		method := "modified"
		switch method {
		case "appraisal":
			specificRisk, err := SpecificRisk(c)
			if err != nil {
				return math.NaN(), err
			}
			result = js / specificRisk
			break
		case "modified":
			_, beta, err := AlphaBeta(c)
			if err != nil {
				return math.NaN(), err
			}
			result = js / beta
			break
		case "alternative":
			sr_data, err := SystematicRisk(c)
			if err != nil {
				return math.NaN(), err
			}
			result = js / sr_data
			break
		default:
			return math.NaN(), errors.New("In AppraisalRatio, method is default !!!")
		}
		return result, nil
	})
}

func SystematicRisk(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("SystematicRisk", func() (float64, error) {
		periodRf := Rf / c.Period()
		rp := c.BenchRatio()
		localRP := make([]float64, len(rp))
		copy(localRP, rp)
		varData, err := Variance(Vector(localRP).AddScalarV(-periodRf))
		if err != nil {
			return math.NaN(), err
		}
		stdDev := math.Sqrt(varData * Scale)
		_, beta, err := AlphaBeta(c)
		if err != nil {
			return math.NaN(), err
		}
		return beta * stdDev, nil
	})
}

func SpecificRisk(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("SpecificRisk", func() (float64, error) {
		alpha, beta, err := AlphaBeta(c)
		if err != nil {
			return math.NaN(), err
		}
		Rb := c.BenchRatio()
		Ra := c.PortfolioRatio()

		vector := make([]float64, len(Ra))
		copy(vector, Ra)
		epsilon := Vector(vector).ImplVectorV(Rb, func(a, b float64) (float64, float64) {
			return a - beta*b - alpha, b
		})
		epAverage := epsilon.Average()
		sum := Vector(vector).ImplV(func(a float64) float64 {
			diff := (a - epAverage)
			return diff * diff
		}).AccumulateSum()
		specificRisk := math.Sqrt(sum / float64(len(Rb)) * float64(Scale))
		return specificRisk, nil
	})
}

func TotalRisk(c MetricCalculator) (float64, error) {
	return c.GetOrSetScalar("TotalRisk", func() (float64, error) {
		sysR, err := SystematicRisk(c)
		if err != nil {
			return math.NaN(), err
		}
		speR, err := SpecificRisk(c)
		if err != nil {
			return math.NaN(), err
		}
		return math.Sqrt(math.Pow(sysR, 2) + math.Pow(speR, 2)), nil
	})
}
