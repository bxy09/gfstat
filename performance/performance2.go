package performance

import (
	"errors"

	"github.com/bxy09/gfstat/performance/utils"
	"math"
)

/// <summary>
///Active Premium
/// The return on an investment's annualized return minus the benchmark's
/// annualized return.
/// Active Premium = Investment's annualized return - Benchmark's annualized
/// </summary>
func ActivePremium(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, scale float64) (float64, error) {
	ra_ann, err := Annualized(Ra, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	rb_ana, err := Annualized(Rb, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	return ra_ann - rb_ana, nil
}

/// <summary>
/// downside risk (deviation, variance) of the return distribution
/// Downside deviation, semideviation, and semivariance are measures of downside
/// risk.
/// </summary>
// = "full"
// = false
//func DownsideDeviation(Ra *utils.SlidingWindow, MAR *utils.SlidingWindow, method string, potential bool) float64 {
func DownsideDeviation(Ra *utils.SlidingWindow, MAR *utils.SlidingWindow) (float64, error) {
	if Ra == nil {
		return math.NaN(), errors.New("In DownsideDeviation, Ra == nil")
	}
	if Ra.Count() <= 0 {
		return math.NaN(), errors.New("In DownsideDeviation, Ra.Count() <= 0")
	}

	r, err := utils.NewSlidingWindow(Ra.Count())
	if err != nil {
		return math.NaN(), err
	}

	newMAR, err := utils.NewSlidingWindow(Ra.Count())
	if err != nil {
		return math.NaN(), err
	}
	len := 0.0
	result := 0.0
	for i := 0; i < Ra.Count(); i++ {
		if Ra.Data()[i] < MAR.Data()[i] {
			r.Add(Ra.Data()[i])
			newMAR.Add(MAR.Data()[i])
		}
	}

	potential := false
	method := "subset"

	if method == "full" {
		len = float64(Ra.Count())
	} else if method == "subset" {
		len = float64(r.Count())
	} else {
		return math.NaN(), errors.New("In DownsideDeviation, method default !!!")
	}
	if newMAR.Count() <= 0 || r.Count() <= 0 || len <= 0 {
		return math.NaN(), errors.New("In DownsideDeviation, newMAR.Count() <= 0 || r.Count() <= 0 || len <= 0")
	}
	if potential {
		sub_Sliding, err := utils.Sub(newMAR, r)
		if err != nil {
			return math.NaN(), err
		}
		result = sub_Sliding.Sum() / len
	} else {
		sub_Sliding, err := utils.Sub(newMAR, r)
		if err != nil {
			return math.NaN(), err
		}
		pow_Sliding, err := utils.Power(sub_Sliding, 2.0)
		if err != nil {
			return math.NaN(), err
		}
		result = math.Sqrt(pow_Sliding.Sum() / len)
	}
	return result, nil
}

/// <summary>
/// downside frequency of the return distribution
/// To calculate Downside Frequency, we take the subset of returns that are
/// less than the target (or Minimum Acceptable Returns (MAR)) returns and
/// divide the length of this subset by the total number of returns.
/// </summary>
func DownsideFrequency(Ra *utils.SlidingWindow, MAR *utils.SlidingWindow) (float64, error) {
	if Ra == nil {
		return math.NaN(), errors.New("In DownsideFrequency, Ra == nil")
	}
	if Ra.Count() <= 0 {
		return math.NaN(), errors.New("In DownsideFrequency, Ra.Count() <= 0")
	}
	len := 0.0
	for i := 0; i < Ra.Count(); i++ {
		if Ra.Data()[i] < MAR.Data()[i] {
			len++
		}
	}

	return len / float64(Ra.Count()), nil
}

/// <summary>
/// A measure of the unexplained portion of performance relative to a benchmark.
/// （年化的超指数收益率序列标准差）
/// </summary>
func TrackingError(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, scale float64) (float64, error) {
	temp, err := Excess(Ra, Rb)
	if err != nil {
		return math.NaN(), err
	}
	return StdDev_Annualized(temp, scale)
}

/*
data(managers)
TrackingError(managers[,1,drop=FALSE], managers[,8,drop=FALSE])
TrackingError(managers[,1:6], managers[,8,drop=FALSE])
TrackingError(managers[,1:6], managers[,8:7,drop=FALSE])
*/

/// <summary>
/// InformationRatio:ActivePremium/TrackingError
/// （经TrackingError调整的超额收益率）
/// </summary>
func InformationRatio(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, scale float64) (float64, error) {
	ap, err := ActivePremium(Ra, Rb, scale)
	if err != nil {
		return math.NaN(), err
	}
	te, err := TrackingError(Ra, Rb, scale)
	if err != nil {
		return math.NaN(), err
	}
	var IR = ap / te
	return IR, nil
}

/// <summary>
/// M squared for Sortino is a M^2 calculated for Downside risk instead of Total Risk
/// （基于SortinoRatio进行的收益率调整）
/// </summary>
func M2Sortino(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, scale float64, MAR float64) (float64, error) {
	Rp, err := Annualized(Ra, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	ra_dd2, err := DownsideDeviation2(Ra, MAR)
	if err != nil {
		return math.NaN(), err
	}
	SigmaD := ra_dd2 * math.Sqrt(float64(scale))
	rb_dd2, err := DownsideDeviation2(Rb, MAR)
	if err != nil {
		return math.NaN(), err
	}
	SigmaDM := rb_dd2 * math.Sqrt(float64(scale))
	SR, err := SortinoRatio(Ra, MAR)
	if err != nil {
		return math.NaN(), err
	}
	var result = Rp + SR*(SigmaDM-SigmaD)
	return result, nil
}

/// <summary>
/// Appraisal ratio is the Jensen's alpha adjusted for specific risk. The numerator
/// is divided by specific risk instead of total risk.
/// </summary>
func AppraisalRatio(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, scale float64, Rf float64, method string) (float64, error) {
	var result = 0.0
	switch method {
	case "appraisal":
		be_data, err := Beta2(Ra, Rb, Rf)
		if err != nil {
			return math.NaN(), err
		}
		multi_Sliding, err := utils.Multi(be_data, Rb)
		if err != nil {
			return math.NaN(), err
		}
		sub_Sliding, err := utils.Sub(Ra, multi_Sliding)
		if err != nil {
			return math.NaN(), err
		}
		al_data, err := Alpha2(Ra, Rb, Rf)
		if err != nil {
			return math.NaN(), err
		}
		epsilon, err := utils.Add(-al_data, sub_Sliding)
		if err != nil {
			return math.NaN(), err
		}
		add_Sliding, err := utils.Add(-epsilon.Average(), epsilon)
		if err != nil {
			return math.NaN(), err
		}
		pow_Sliding, err := utils.Power(add_Sliding, 2)
		if err != nil {
			return math.NaN(), err
		}
		specifikRisk := math.Sqrt(pow_Sliding.Sum()/float64(epsilon.Count())) * math.Sqrt(float64(scale))
		jsa_data, err := JensenAlpha2(Ra, Rb, Rf, scale)
		if err != nil {
			return math.NaN(), err
		}
		result = jsa_data / specifikRisk
		break
	case "modified":
		jsa2_data, err := JensenAlpha2(Ra, Rb, Rf, scale)
		if err != nil {
			return math.NaN(), err
		}
		be2_data, err := Beta2(Ra, Rb, Rf)
		if err != nil {
			return math.NaN(), err
		}
		result = jsa2_data / be2_data
		break
	case "alternative":
		jsa2_data, err := JensenAlpha2(Ra, Rb, Rf, scale)
		if err != nil {
			return math.NaN(), err
		}
		sr_data, err := SystematicRisk(Ra, Rb, scale, Rf)
		if err != nil {
			return math.NaN(), err
		}
		result = jsa2_data / sr_data
		break
	default:
		return math.NaN(), errors.New("In AppraisalRatio, method is default !!!")
	}
	return result, nil
}

/// <summary>
/// Fama beta is a beta used to calculate the loss of diversification. It is made
/// so that the systematic risk is equivalent to the total portfolio risk.
/// </summary>
func FamaBeta(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, Ra_sclae float64, Rb_scale float64) (float64, error) {
	var n1 = Ra.Count()
	var n2 = Rb.Count()
	var_Ra, err := Variance(Ra)
	if err != nil {
		return math.NaN(), err
	}
	var_Rb, err := Variance(Rb)
	if err != nil {
		return math.NaN(), err
	}
	var result = math.Sqrt(var_Ra*float64(n1-1)/float64(n1)) * math.Sqrt(float64(Ra_sclae)) / (math.Sqrt(var_Rb*float64(n2-1)/float64(n2)) * math.Sqrt(float64(Rb_scale)))
	return result, nil
}

/// <summary>
/// Selectivity is the same as Jensen's alpha
/// </summary>
func Selectivity(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, scale float64, Rf float64) (float64, error) {
	return JensenAlpha2(Ra, Rb, Rf, scale)
}

/// <summary>
/// epsilon与R中不同,但似乎没有影响
/// Specific risk is the standard deviation of the error term in the
/// regression equation.
/// </summary>
func SpecificRisk(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, scale float64, Rf float64) (float64, error) {
	//Period = Frequency(Ra)
	alpha, err := Alpha2(Ra, Rb, Rf)
	if err != nil {
		return math.NaN(), err
	}
	beta, err := Beta2(Ra, Rb, Rf)
	if err != nil {
		return math.NaN(), err
	}
	add_Ra_Sliding, err := utils.Add(-Rf, Ra)
	if err != nil {
		return math.NaN(), err
	}
	add_Rb_Sliding, err := utils.Add(-Rf, Rb)
	if err != nil {
		return math.NaN(), err
	}
	multi_beta_Slidinig, err := utils.Multi(beta, add_Rb_Sliding)
	if err != nil {
		return math.NaN(), err
	}
	sub_Ra_Beta, err := utils.Sub(add_Ra_Sliding, multi_beta_Slidinig)
	if err != nil {
		return math.NaN(), err
	}
	epsilon, err := utils.Add(-alpha, sub_Ra_Beta)
	if err != nil {
		return math.NaN(), err
	}
	var_eps, err := Variance(epsilon)
	if err != nil {
		return math.NaN(), err
	}
	var result = math.Sqrt(var_eps*float64(epsilon.Count()-1)/float64(epsilon.Count())) * math.Sqrt(float64(scale))
	return result, nil
}

/// <summary>
/// Systematic risk as defined by Bacon(2008) is the product of beta by market
/// risk. Be careful ! It's not the same definition as the one given by Michael
/// Jensen. Market risk is the standard deviation of the benchmark. The systematic
/// risk is annualized
/// </summary>
func SystematicRisk(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, scale float64, Rf float64) (float64, error) {
	beta_Sliding, err := Beta2(Ra, Rb, Rf)
	if err != nil {
		return math.NaN(), err
	}
	exce_Data, err := Excess2(Rb, Rf)
	if err != nil {
		return math.NaN(), err
	}
	stdDev, err := StdDev_Annualized(exce_Data, scale)
	if err != nil {
		return math.NaN(), err
	}
	var result = beta_Sliding * stdDev
	return result, nil
}

/// <summary>
/// The square of total risk is the sum of the square of systematic risk and the square
/// of specific risk. Specific risk is the standard deviation of the error term in the
/// regression equation. Both terms are annualized to calculate total risk.
/// （总风险,注意这是一个经过开方之后的数值：SystematicRisk+SpecificRisk）
/// </summary>
func TotalRisk(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, scale float64, Rf float64) (float64, error) {
	SR_data, err := SystematicRisk(Ra, Rb, scale, Rf)
	if err != nil {
		return math.NaN(), err
	}
	Specific, err := SpecificRisk(Ra, Rb, scale, Rf)
	if err != nil {
		return math.NaN(), err
	}
	var result = math.Sqrt(math.Pow(SR_data, 2) + math.Pow(Specific, 2))
	return result, nil
}

/// <summary>
/// The Treynor ratio is similar to the Sharpe Ratio, except it uses beta as the
/// volatility measure (to divide the investment's excess return over the beta).
/// （组合诸多的风险调整收益率之一）
/// </summary>
func TreynorRatio(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, scale float64, Rf float64) (float64, error) {
	beta, err := Beta2(Ra, Rb, Rf)
	if err != nil {
		return math.NaN(), err
	}
	e2, err := Excess2(Ra, Rf)
	if err != nil {
		return math.NaN(), err
	}
	TR, err := Annualized(e2, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	return TR / beta, nil
}

/// <summary>
/// 只测试了默认参数
/// Calculate metrics on how the asset in R performed in up and down markets,
/// measured by periods when the benchmark asset was up or down.
/// Up (Down) Capture Ratio: this is a measure of an investment's compound
/// return when the benchmark was up (down) divided by the benchmark's compound
/// return when the benchmark was up (down). The greater (lower) the value, the
/// better.(Up越大越好，Down越小越好)
///
/// Up (Down) Number Ratio: similarly, this is a measure of the number of
/// periods that the investment was up (down) when the benchmark was up (down),
/// divided by the number of periods that the Benchmark was up (down).(Up越大越好，Down越小越好)
///
/// Up (Down) Percentage Ratio: this is a measure of the number of periods that
/// the investment outperformed the benchmark when the benchmark was up (down),
/// divided by the number of periods that the benchmark was up (down). Unlike
/// the prior two metrics, in both cases a higher value is better.(Up、Down均为越大越好)
/// （当市场涨跌时，组合收益率涨跌所占比率，）
/// </summary>
func UpDownRatios(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow) (float64, error) {
	var cumRa = 0.0
	var cumRb = 0.0
	var result = 0.0

	method := "Capture"
	side := "Up"

	switch method {
	case "Capture":

		switch side {
		case "Up":

			UpRa, err := utils.NewSlidingWindow(Ra.Count())
			if err != nil {
				return math.NaN(), err
			}
			UpRb, err := utils.NewSlidingWindow(Ra.Count())
			if err != nil {
				return math.NaN(), err
			}
			for i := 0; i < Ra.Count(); i++ {
				if Rb.Data()[i] > 0 {
					UpRa.Add(Ra.Data()[i])
					UpRb.Add(Rb.Data()[i])
				}
			}
			cumRa = UpRa.Sum()
			cumRb = UpRb.Sum()
			result = cumRa / cumRb
			return result, nil

		case "Down":

			DnRa, err := utils.NewSlidingWindow(Ra.Count())
			if err != nil {
				return math.NaN(), err
			}
			DnRb, err := utils.NewSlidingWindow(Ra.Count())
			if err != nil {
				return math.NaN(), err
			}
			for i := 0; i < Ra.Count(); i++ {
				if Rb.Data()[i] <= 0 {
					DnRa.Add(Ra.Data()[i])
					DnRb.Add(Rb.Data()[i])
				}
			}
			cumRa = DnRa.Sum()
			cumRb = DnRb.Sum()
			result = cumRa / cumRb
			return result, nil

		default:
			return math.NaN(), errors.New("In UpDownRatios, method Default!!")
		}

	case "Number":

		switch side {
		case "Up":

			UpRa, err := utils.NewSlidingWindow(Ra.Count())
			if err != nil {
				return math.NaN(), err
			}
			UpRb, err := utils.NewSlidingWindow(Ra.Count())
			if err != nil {
				return math.NaN(), err
			}
			for i := 0; i < Ra.Count(); i++ {
				if Ra.Data()[i] > 0 && Rb.Data()[i] > 0 {
					UpRa.Add(Ra.Data()[i])
				}
			}
			for i := 0; i < Ra.Count(); i++ {
				if Rb.Data()[i] > 0 {
					UpRb.Add(Rb.Data()[i])
				}
			}

			cumRa = float64(UpRa.Count())
			cumRb = float64(UpRb.Count())
			result = cumRa / cumRb
			return result, nil

		case "Down":

			DnRa, err := utils.NewSlidingWindow(Ra.Count())
			if err != nil {
				return math.NaN(), err
			}
			DnRb, err := utils.NewSlidingWindow(Ra.Count())
			if err != nil {
				return math.NaN(), err
			}
			for i := 0; i < Ra.Count(); i++ {
				if Ra.Data()[i] < 0 && Rb.Data()[i] < 0 {
					DnRa.Add(Ra.Data()[i])
				}
			}
			for i := 0; i < Ra.Count(); i++ {
				if Rb.Data()[i] < 0 {
					DnRb.Add(Rb.Data()[i])
				}
			}

			cumRa = float64(DnRa.Count())
			cumRb = float64(DnRb.Count())
			result = cumRa / cumRb
			return result, nil

		default:
			return math.NaN(), errors.New("In UpDownRatios, method default 2 error !!!")
		}

	case "Percent":

		switch side {
		case "Up":

			UpRa, err := utils.NewSlidingWindow(Ra.Count())
			if err != nil {
				return math.NaN(), err
			}
			UpRb, err := utils.NewSlidingWindow(Ra.Count())
			if err != nil {
				return math.NaN(), err
			}
			for i := 0; i < Ra.Count(); i++ {
				if Ra.Data()[i] > Rb.Data()[i] && Rb.Data()[i] > 0 {
					UpRa.Add(Ra.Data()[i])
				}
			}
			for i := 0; i < Ra.Count(); i++ {
				if Rb.Data()[i] > 0 {
					UpRb.Add(Rb.Data()[i])
				}
			}

			cumRa = float64(UpRa.Count())
			cumRb = float64(UpRb.Count())
			result = cumRa / cumRb
			return result, nil

		case "Down":

			DnRa, err := utils.NewSlidingWindow(Ra.Count())
			if err != nil {
				return math.NaN(), err
			}
			DnRb, err := utils.NewSlidingWindow(Ra.Count())
			if err != nil {
				return math.NaN(), err
			}
			for i := 0; i < Ra.Count(); i++ {
				if Ra.Data()[i] > Rb.Data()[i] && Rb.Data()[i] < 0 {
					DnRa.Add(Ra.Data()[i])
				}
			}
			for i := 0; i < Ra.Count(); i++ {
				if Rb.Data()[i] < 0 {
					DnRb.Add(Rb.Data()[i])
				}
			}

			cumRa = float64(DnRa.Count())
			cumRb = float64(DnRb.Count())
			result = cumRa / cumRb
			return result, nil

		default:
			return math.NaN(), errors.New("In UpDownRatios, method default 3 is Error !!!")
		}

	default:
		return math.NaN(), errors.New("In UpDownRatios, method default 4 is Error !!!")
	}
	return math.NaN(), nil
}

/// <summary>
/// Omega excess return is another form of downside risk-adjusted return. It is
/// calculated by multiplying the downside variance of the style benchmark by 3
/// times the style beta.
/// （经过Ra,Rb的DownsideDeviation调整后的年化收益率，严格意义并不十分明确）
/// </summary>
func OmegaExcessReturn(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, scale float64, MAR float64) (float64, error) {
	Rp, err := Annualized(Ra, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	SigmaD, err := DownsideDeviation2(Ra, MAR)
	if err != nil {
		return math.NaN(), err
	}
	SigmaD = SigmaD * math.Sqrt(float64(scale))
	SigmaDM, err := DownsideDeviation2(Rb, MAR)
	if err != nil {
		return math.NaN(), err
	}
	SigmaDM = SigmaDM * math.Sqrt(float64(scale))
	var result = Rp - 3.0*SigmaD*SigmaDM
	return result, nil
}

/// <summary>
/// M squared is a risk adjusted return useful to judge the size of relative
/// performance between differents portfolios. With it you can compare portfolios
/// with different levels of risk.
/// （使得不同组合的收益率可比的调整措施）
/// </summary>
func MSquared(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, scale float64, Rf float64) (float64, error) {
	var n = Ra.Count()
	Rp, err := Annualized(Ra, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	var_Ra_data, err := Variance(Ra)
	if err != nil {
		return math.NaN(), err
	}
	sigp := math.Sqrt(var_Ra_data*float64(n-1)/float64(n)) * math.Sqrt(float64(scale))
	if err != nil {
		return math.NaN(), err
	}
	var_Rb_data, err := Variance(Rb)
	if err != nil {
		return math.NaN(), err
	}
	var sigm = math.Sqrt(var_Rb_data*float64(n-1)/float64(n)) * math.Sqrt(float64(scale))
	//var result = (Rp-Rf)*sigp/sigm + Rf//Source
	Rf = Rf * scale
	var result = (Rp-Rf)*sigm/sigp + Rf
	return result, nil
}

/// <summary>
/// M squared excess is the quantity above the standard M.
/// There is a geometric excess return which is better for Bacon and an arithmetic excess return
/// （是与Rb的年化收益率进行的excess比较）
/// </summary>
func MSquaredExcess(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, scale float64, Rf float64, method string) (float64, error) {
	//var n = Rb.Count() //Ra&Rb等长
	Rbp, err := Annualized(Rb, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	var result float64
	switch method {
	case "geometric":
		msq_data, err := MSquared(Ra, Rb, scale, Rf)
		if err != nil {
			return math.NaN(), err
		}
		result = (1.0+msq_data)/(1.0+Rbp) - 1.0
		break
	case "arithmetic":
		msq_data, err := MSquared(Ra, Rb, scale, Rf)
		if err != nil {
			return math.NaN(), err
		}
		result = msq_data - Rbp
		break
	default:
		return math.NaN(), errors.New("In MSquaredExcess, method default !!!")
	}
	return result, nil
}
