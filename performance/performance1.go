package performance

import (
	"errors"
	"github.com/bxy09/gfstat/performance/utils"
	"math"
)

/// <summary>
/// BernardoLedoitRatio：take the sum of the subset of
/// returns that are above 0 and we divide it by the opposite of the sum of
/// the subset of returns that are below 0
/// （正收益率汇总/负收益率汇总；粗略描述胜败比率的总体水平，但注意，负收益率的力量更大，不能满足于1）
/// </summary>
func BernardoLedoitRatio(Ra *utils.SlidingWindow) (float64, error) {
	positivevalues, negativevalues, err := utils.PosNegValues(Ra)
	if err != nil {
		return math.NaN(), err
	}
	return -positivevalues.Sum() / negativevalues.Sum(), nil
}

/// <summary>
/// d ratio of the return distribution
/// The d ratio is similar to the Bernado Ledoit ratio but inverted and
/// taking into account the frequency of positive and negative returns.
/// </summary>
func DRatio(Ra *utils.SlidingWindow) (float64, error) {
	if Ra == nil {
		return math.NaN(), errors.New("In DRatio, Ra == nil")
	}
	if Ra.Count() <= 0 {
		return math.NaN(), errors.New("In DRatio, Ra.Count() <= 0")
	}

	upList, _ := utils.NewSlidingWindow(Ra.Count())
	downList, _ := utils.NewSlidingWindow(Ra.Count())

	for i := 0; i < Ra.Count(); i++ {
		if Ra.Data()[i] < 0 {
			downList.Add(Ra.Data()[i])
		} else if Ra.Data()[i] > 0 {
			upList.Add(Ra.Data()[i])
		}
	}

	return -(downList.Sum() * float64(downList.Count())) / (float64(upList.Sum()) * float64(upList.Count())), nil
}

/// <summary>
/// To calculate Mean absolute deviation we take
/// the sum of the absolute value of the difference between the returns and the mean of the returns
/// and we divide it by the number of returns.
/// （描述收益率偏离均值得一个指标）
/// </summary>
func MeanAbsoluteDeviation(Ra *utils.SlidingWindow) (float64, error) {
	if Ra.Count() <= 0 {
		return math.NaN(), errors.New("In MeanAbsoluteDeviation, Ra.Count() <= 0")
	}
	add_Sliding, _ := utils.Add(-Ra.Average(), Ra)
	ads_Sliding, _ := utils.Abs(add_Sliding)
	return ads_Sliding.Sum() / float64(Ra.Count()), nil
}

/// <summary>
/// 偏度峰度比
/// </summary>
func SkewnessKurtosisRatio(Ra *utils.SlidingWindow) (float64, error) {
	ske, err := Skewness(Ra)
	if err != nil {
		return math.NaN(), err
	}
	kur, err := Kurtosis(Ra)
	if err != nil {
		return math.NaN(), err
	}
	return ske / kur, nil
}

/// <summary>
/// 收益率序列的几何均值，非年化
/// </summary>
func MeanGeometric(Ra *utils.SlidingWindow) (float64, error) {
	if Ra.Count() <= 0 {
		return math.NaN(), errors.New("In MeanGeometric, Ra.Count() <= 0")
	}
	add_Sliding, _ := utils.Add(1, Ra)
	log_Sliding, _ := utils.Log(add_Sliding)
	return math.Exp(log_Sliding.Average()) - 1.0, nil
}

/// <summary>
/// 方差
/// </summary>
func Variance(Ra *utils.SlidingWindow) (float64, error) {
	if Ra == nil || Ra.Count() <= 1 {
		return math.NaN(), errors.New("In Variance, Ra == nil || Ra.Count() <= 1")
	}

	result := 0.0
	mean := Ra.Average()
	for i := 0; i < Ra.Count(); i++ {
		result += (Ra.Data()[i] - mean) * (Ra.Data()[i] - mean)
	}
	return result / (float64)(Ra.Count()-1), nil
}

/// <summary>
/// 标准差
/// </summary>
func StdDev(Ra *utils.SlidingWindow) (float64, error) {
	data, err := Variance(Ra)
	if err != nil {
		return math.NaN(), err
	}
	return math.Sqrt(data), nil
}

/// <summary>
/// 年化标准差
/// </summary>
func StdDev_Annualized(Ra *utils.SlidingWindow, scale float64) (float64, error) {
	data, err := StdDev(Ra)
	if err != nil {
		return math.NaN(), err
	}
	return math.Sqrt(float64(scale)) * data, nil
}

/// <summary>
/// Sortino proposed an improvement on the Sharpe Ratio to better account for
/// skill and excess performance by using only downside semivariance as the
/// measure of risk.Sortino contends that risk should be measured in terms of not meeting the
/// investment goal.
/// （引入MAR，并开始在调整收益率的分子分母进行Excess return与DownsideDeviation的MAR调整）
/// </summary>
func SortinoRatio(Ra *utils.SlidingWindow, MAR float64) (float64, error) {
	exce_Sliding, err := Excess2(Ra, MAR)
	if err != nil {
		return math.NaN(), err
	}

	ddata, err := DownsideDeviation2(Ra, MAR)
	if err != nil {
		return math.NaN(), err
	}
	var SR = exce_Sliding.Average() / ddata
	return SR, nil
}

/// <summary>
/// Prospect ratio is a ratio used to penalise loss since most people feel loss
/// greater than gain
/// （经验类型调整收益率，给损失赋予更大的权重）
/// </summary>
func ProspectRatio(Ra *utils.SlidingWindow, MAR float64) (float64, error) {
	var n = Ra.Count()
	SigD, err := DownsideDeviation2(Ra, MAR)
	if err != nil {
		return math.NaN(), err
	}

	positivevalues, negativevalues, err := utils.PosNegValues(Ra)
	if err != nil {
		return math.NaN(), err
	}

	var result = ((positivevalues.Sum()+2.25*negativevalues.Sum())/float64(n) - MAR) / SigD
	return result, nil
}

func DownsideFrequency2(Ra *utils.SlidingWindow, MAR float64) (float64, error) {
	if Ra == nil || Ra.Count() <= 0 {
		return math.NaN(), errors.New("In DownsideFrequency2, Ra == nil || Ra.Count() <= 0 !!")
	}
	newMAR, _ := utils.CreateList(MAR, Ra.Count())

	return DownsideFrequency(Ra, newMAR)
}

/// <summary>
/// 最大回撤，默认为返回其相反数
/// </summary>
func MaxDrawdown(Ra *utils.SlidingWindow) (float64, error) {
	drawdowns, err := Drawdowns(Ra)
	if err != nil {
		return math.NaN(), err
	}
	result := drawdowns[0]
	for _, d := range drawdowns {
		if d < result {
			result = d
		}
	}
	return -result, nil
}

/// <summary>
/// Calculate the drawdown levels in a timeseries
/// </summary>
//= true
func Drawdowns(Rb *utils.SlidingWindow) ([]float64, error) {
	Ra := Rb.Data()
	if Ra == nil || len(Ra) <= 0 {
		return nil, errors.New("In Drawdowns, Ra == nil")
	}

	geometric := 1

	curReturn := 1.0
	curMaxReturn := 1.0 + Ra[0]
	result := []float64{}
	if geometric == 1 {
		for _, r := range Ra {
			curReturn = curReturn * (1.0 + r)
			if curReturn > curMaxReturn {
				curMaxReturn = curReturn
			}
			result = append(result, curReturn/curMaxReturn-1.0)
		}
	} else {
		for _, r := range Ra {
			curReturn = curReturn + r
			if curReturn > curMaxReturn {
				curMaxReturn = curReturn
			}
			result = append(result, curReturn/curMaxReturn-1.0)
		}
	}

	return result, nil
}

/// <summary>
/// subset of returns that are
/// more than the target (or Minimum Acceptable Returns (MAR)) returns and
/// divide the length of this subset by the total number of returns.
/// （超过MAR的频率）
/// </summary>
func UpsideFrequency(Ra *utils.SlidingWindow, MAR float64) (float64, error) {
	aboveMAR, err := utils.AboveValue(Ra, MAR)
	if err != nil {
		return math.NaN(), err
	}
	return float64(aboveMAR.Count()) / float64(Ra.Count()), nil
}

/// <summary>
/// 偏度
/// </summary>
// default = "moment"
func Skewness(Ra *utils.SlidingWindow) (float64, error) {
	if Ra == nil || Ra.Count() <= 2 {
		return math.NaN(), errors.New("In Skewness, Ra == nil || Ra.Count() <= 2")
	}

	n := float64(Ra.Count())
	method := "moment"
	switch method {
	//"moment", "fisher", "sample"
	case "moment": //skewness = sum((x-mean(x))^3/sqrt(var(x)*(n-1)/n)^3)/length(x)
		var_data, err := Variance(Ra)
		if err != nil {
			return math.NaN(), err
		}
		add_Sliding, err := utils.Add(-Ra.Average(), Ra)
		if err != nil {
			return math.NaN(), err
		}
		pow_Sliding, err := utils.Power(add_Sliding, 3.0)
		if err != nil {
			return math.NaN(), err
		}
		multi_Sliding, err := utils.Multi(1.0/math.Pow(var_data*(n-1.0)/n, 1.5), pow_Sliding)
		if err != nil {
			return math.NaN(), err
		}
		return multi_Sliding.Sum() / n, nil
	default:
		return math.NaN(), errors.New("In Skewness, method is default")
	}
	return math.NaN(), nil
}

/// <summary>
/// 峰度
/// </summary>
// = "sample"
func Kurtosis(Ra *utils.SlidingWindow) (float64, error) {
	if Ra == nil || Ra.Count() <= 3 {
		return math.NaN(), errors.New("In Kurtosis, Ra == nil || Ra.Count() <= 3")
	}

	n := float64(Ra.Count())
	method := "sample_excess"
	switch method {
	case "sample_excess": //kurtosis = sum((x-mean(x))^4/var(x)^2)*n*(n+1)/((n-1)*(n-2)*(n-3)) - 3*(n-1)^2/((n-2)*(n-3))
		var_data, err := Variance(Ra)
		if err != nil {
			return math.NaN(), err
		}
		add_Sliding, err := utils.Add(-Ra.Average(), Ra)
		if err != nil {
			return math.NaN(), err
		}
		pow_Sliding, err := utils.Power(add_Sliding, 4.0)
		if err != nil {
			return math.NaN(), err
		}
		multi_Sliding, err := utils.Multi(1.0/math.Pow(var_data, 2.0), pow_Sliding)
		if err != nil {
			return math.NaN(), err
		}
		return multi_Sliding.Sum()*n*(n+1.0)/((n-1.0)*(n-2.0)*(n-3.0)) - 3*(n-1.0)*(n-1.0)/((n-2.0)*(n-3.0)), nil
	default:
		return math.NaN(), errors.New("In Kurtosis, method is default")
	}
	return math.NaN(), nil
}

//= "full"
// = false
func DownsideDeviation2(Ra *utils.SlidingWindow, MAR float64) (float64, error) {
	if Ra == nil || Ra.Count() <= 0 {
		return math.NaN(), errors.New("In DownsideDeviation2, Ra == nil || Ra.Count() <= 0")
	}
	newMAR, _ := utils.CreateList(MAR, Ra.Count())
	return DownsideDeviation(Ra, newMAR)
}

/// <summary>
/// Adjusted Sharpe ratio of the return distribution
/// Adjusted Sharpe ratio was introduced by Pezier and White (2006) to adjusts
/// for skewness and kurtosis by incorporating a penalty factor for negative skewness
/// and excess kurtosis.
/// </summary>
func AdjustedSharpeRatio(Ra *utils.SlidingWindow, Rf float64, scale float64) (float64, error) {
	Rp, err := Annualized(Ra, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	Sigp, err := StdDev_Annualized(Ra, scale)
	if err != nil {
		return math.NaN(), err
	}
	Rf = Rf * scale
	SR := (Rp - Rf) / Sigp
	K, err := Kurtosis(Ra)
	if err != nil {
		return math.NaN(), err
	}
	S, err := Skewness(Ra)
	if err != nil {
		return math.NaN(), err
	}
	var result = SR * (1.0 + (S/6.0)*SR - ((K-3.0)/24.0)*math.Pow(SR, 2.0))
	return result, nil
}

/// <summary>
/// Kappa is a generalized downside risk-adjusted performance measure.
/// To calculate it, we take the difference of the mean of the distribution
/// to the target and we divide it by the l-root of the lth lower partial
/// moment. To calculate the lth lower partial moment we take the subset of
/// returns below the target and we sum the differences of the target to
/// these returns. We then return return this sum divided by the length of
/// the whole distribution.
/// （非年化的超MAR平均收益率通过l阶根的低于MAR的收益率序列的l阶矩）
/// </summary>
func Kappa(Ra *utils.SlidingWindow, MAR float64, l float64) (float64, error) {
	undervalues, err := utils.NewSlidingWindow(Ra.Count())
	if err != nil {
		return math.NaN(), err
	}
	for i := 0; i < Ra.Count(); i++ {
		if Ra.Data()[i] < MAR {
			undervalues.Add(Ra.Data()[i])
		}
	}

	var n = float64(Ra.Count())
	var m = float64(Ra.Average())
	neg_Sliding, err := utils.Negative(undervalues)
	if err != nil {
		return math.NaN(), err
	}
	add_Sliding, err := utils.Add(MAR, neg_Sliding)
	if err != nil {
		return math.NaN(), err
	}
	pow_Sliding, err := utils.Power(add_Sliding, float64(l))
	if err != nil {
		return math.NaN(), err
	}
	var temp = pow_Sliding.Sum() / n
	return (m - MAR) / math.Pow(temp, (1.0/float64(l))), nil
}

/// <summary>
/// To calculate Burke ratio we take the difference between the portfolio
/// return and the risk free rate and we divide it by the square root of the
/// sum of the square of the drawdowns. To calculate the modified Burke ratio
/// we just multiply the Burke ratio by the square root of the number of datas.
/// （一种调整收益率的计算方式，调整是通过drawdown的平方和进行的）
/// </summary>
func BurkeRatio(Ra *utils.SlidingWindow, Rf float64, scale float64) (float64, error) {
	var len = Ra.Count()
	var in_drawdown = false
	var peak = 1
	var temp = 0.0
	drawdown, err := utils.NewSlidingWindow(len)
	if err != nil {
		return math.NaN(), err
	}
	for i := 1; i < len; i++ {
		if Ra.Data()[i] < 0 {
			if !in_drawdown {
				peak = i - 1
				in_drawdown = true
			}
		} else {
			if in_drawdown {
				temp = 1.0
				for j := peak + 1; j < i; j++ {
					temp = temp * (1.0 + Ra.Data()[j])
				}
				drawdown.Add(temp - 1.0) //Source
				in_drawdown = false
			}
		}
	}

	if in_drawdown {
		temp = 1.0
		for j := peak + 1; j < len; j++ {
			temp = temp * (1.0 + Ra.Data()[j])
		}
		drawdown.Add(temp - 1.0) //Source
		//drawdown.Add((temp - 1.0) * 100.0)
		in_drawdown = false
	}
	//var Rp = Annualized(Ra, scale, true) - 1.0--->Source
	Rp, err := Annualized(Ra, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	var result float64

	if drawdown.Count() != 0 {
		pow_Sliding, err := utils.Power(drawdown, 2)
		if err != nil {
			return math.NaN(), err
		}
		Rf = Rf * scale
		result = (Rp - Rf) / math.Sqrt(pow_Sliding.Sum())
	} else {
		result = 0
	}

	modified := true
	if modified {
		result = result * math.Sqrt(float64(len))
	}
	return result, nil
}

/// <summary>
///  Upside Risk is the similar of semideviation taking the return above the
///  Minimum Acceptable Return instead of using the mean return or zero.
///  （一般来说，非对称类的比较，单求此统计量意义有限）
/// </summary>
func UpsideRisk(Ra *utils.SlidingWindow, MAR float64, stat string) (float64, error) {
	r, err := utils.AboveValue(Ra, MAR)
	if err != nil {
		return math.NaN(), err
	}
	var length float64
	method := "subset"
	switch method {
	case "full":
		length = float64(Ra.Count())
		break
	case "subset":
		length = float64(r.Count())
		break
	default:
		return math.NaN(), errors.New("In Upside Risk, method is default !!!")
	}
	if length <= 0 {
		return 0, nil
	}
	var result float64
	switch stat {
	case "risk":
		add_Sliding, err := utils.Add(-MAR, r)
		if err != nil {
			return math.NaN(), err
		}
		pow_Sliding, err := utils.Power(add_Sliding, 2.0)
		if err != nil {
			return math.NaN(), err
		}
		multi_Sliding, err := utils.Multi(1.0/length, pow_Sliding)
		if err != nil {
			return math.NaN(), err
		}
		result = math.Sqrt(multi_Sliding.Sum())
		break
	case "variance":
		add_Sliding, err := utils.Add(-MAR, r)
		if err != nil {
			return math.NaN(), err
		}
		pow_Sliding, err := utils.Power(add_Sliding, 2.0)
		if err != nil {
			return math.NaN(), err
		}
		multi_Sliding, err := utils.Multi(1.0/length, pow_Sliding)
		if err != nil {
			return math.NaN(), err
		}
		result = multi_Sliding.Sum()
		break
	case "potential":
		add_Sliding, err := utils.Add(-MAR, r)
		if err != nil {
			return math.NaN(), err
		}
		multi_Slding, err := utils.Multi(1.0/length, add_Sliding)
		if err != nil {
			return math.NaN(), err
		}
		result = multi_Slding.Sum()
		break
	default:
		return math.NaN(), errors.New("In UpSide Risk, method is default !!!")
	}

	return result, nil
}

/// <summary>
/// the Kelly criterion is equal to the expected excess return of the strategy
/// divided by the expected variance of the excess return
/// （非年化的平均超额收益除以非年化的方差）
/// </summary>
func KellyRatio_Full(Ra *utils.SlidingWindow, Rf float64) (float64, error) {
	xR, err := Excess2(Ra, Rf)
	if err != nil {
		return math.NaN(), err
	}
	var_data, err := Variance(Ra)
	if err != nil {
		return math.NaN(), err
	}
	KR := xR.Average() / var_data
	return KR, nil
}

func KellyRatio_Half(Ra *utils.SlidingWindow, Rf float64) (float64, error) {
	xR, err := Excess2(Ra, Rf)
	if err != nil {
		return math.NaN(), err
	}
	var_data, err := Variance(Ra)
	if err != nil {
		return math.NaN(), err
	}
	var KR = xR.Average() / var_data
	KR = KR / 2
	return KR, nil
}

/// <summary>
/// Upside Potential Ratio,compared to Sortino, was a further improvement, extending the
/// measurement of only upside on the numerator, and only downside of the
/// denominator of the ratio equation.
/// （分子只考虑超过MAR部分，分母只考虑DownsideDeviation的下跌风险）
/// </summary>
func UpsidePotentialRatio(Ra *utils.SlidingWindow, MAR float64) (float64, error) {
	//var r = Ra.Where<float64>(singleData => singleData > MAR).ToList<float64>();
	r, err := utils.AboveValue(Ra, MAR)
	if err != nil {
		return math.NaN(), err
	}
	var length int
	method := "subset"
	switch method {
	case "full":
		length = Ra.Count()
		break
	case "subset":
		length = r.Count()
		break
	default:
		return math.NaN(), errors.New("In UpsidePotentialRatio, method is default !!!")
	}
	add_Sliding, err := utils.Add(-MAR, r)
	if err != nil {
		return math.NaN(), err
	}
	dd2Data, err := DownsideDeviation2(Ra, MAR)
	if err != nil {
		return math.NaN(), err
	}
	var result = (add_Sliding.Sum() / float64(length)) / dd2Data
	return result, nil
}

/// <summary>
/// Volatility skewness is a similar measure to omega but using the second
/// partial moment. It's the ratio of the upside variance compared to the
/// downside variance. Variability skewness is the ratio of the upside risk
/// compared to the downside risk.
/// （评价收益率分布的偏度，应该是越大越好，与1的关系要看UpsideRisk与DownsideDeviation定义是否一致）
/// </summary>
func VolatilitySkewness_Variance(Ra *utils.SlidingWindow, MAR float64) (float64, error) {
	usr_data, err := UpsideRisk(Ra, MAR, "variance")
	if err != nil {
		return math.NaN(), err
	}
	dd2, err := DownsideDeviation2(Ra, MAR)
	if err != nil {
		return math.NaN(), err
	}
	return usr_data / math.Pow(dd2, 2.0), nil
}

func VolatilitySkewness_Risk(Ra *utils.SlidingWindow, MAR float64) (float64, error) {
	usr, err := UpsideRisk(Ra, MAR, "risk")
	if err != nil {
		return math.NaN(), err
	}
	dd2, err := DownsideDeviation2(Ra, MAR)
	if err != nil {
		return math.NaN(), err
	}
	return usr / dd2, nil
}

func CalmarRatio(Ra *utils.SlidingWindow, scale float64) (float64, error) {
	annualized_return, err := Annualized(Ra, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	md_data, err := MaxDrawdown(Ra)
	if err != nil {
		return math.NaN(), err
	}
	draw_down := math.Abs(md_data)
	return annualized_return / draw_down, nil
}

func SterlingRatio(Ra *utils.SlidingWindow, scale, excess float64) (float64, error) {
	annualized_return, err := Annualized(Ra, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	md_data, err := MaxDrawdown(Ra)
	if err != nil {
		return math.NaN(), err
	}
	draw_down := math.Abs(md_data + excess)
	if draw_down == 0.0 {
		return math.NaN(), errors.New("In SterlingRatio, draw_down == 0.0")
	}
	return annualized_return / draw_down, nil
}

func PainIndex(Ra *utils.SlidingWindow) (float64, error) {
	data, err := Drawdowns(Ra)
	if err != nil {
		return math.NaN(), err
	}
	total := 0.0
	for _, val := range data {
		total += math.Abs(val)
	}
	if len(data) == 0 {
		return math.NaN(), errors.New("In PainIndex, len(data) == 0")
	}
	return total / float64(len(data)), nil
}

func PainRatio(Ra *utils.SlidingWindow, Rf float64, scale float64) (float64, error) {
	PI, err := PainIndex(Ra)
	if err != nil {
		return math.NaN(), err
	}
	n := Ra.Count()
	add_Sliding, err := utils.Add(1.0, Ra)
	if err != nil {
		return math.NaN(), err
	}
	prod_Sliding, err := utils.Prod(add_Sliding)
	if err != nil {
		return math.NaN(), err
	}
	Rp := math.Pow(prod_Sliding, float64(scale)/float64(n)) - 1.0
	Rf = Rf * scale
	return (Rp - Rf) / PI, nil
}

func FindDrawdowns(Ra *utils.SlidingWindow) map[string][]float64 {
	drawdowns, err := Drawdowns(Ra)
	if err != nil {
		return nil
	}

	var draw []float64
	var begin []float64
	var end []float64
	var trough []float64
	var length []float64
	var recovery []float64

	priorSign := 0
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

	results := map[string][]float64{"draw": draw, "begin": begin, "trough": trough, "end": end, "length": length, "recovery": recovery}

	return results
}

func AverageDrawdown(Ra *utils.SlidingWindow) (float64, error) {
	Dj := FindDrawdowns(Ra)["draw"]
	if len(Dj) <= 0 {
		return math.NaN(), errors.New("In AverageDrawdown, len(Dj) <= 0")
	}
	len_NoneZero := 0
	total_NoneZero := 0.0
	for _, val := range Dj {
		if val < 0 {
			len_NoneZero++
			total_NoneZero += val
		}
	}
	result := math.Abs(total_NoneZero / float64(len_NoneZero))
	return result, nil
}

func AverageLength(Ra *utils.SlidingWindow) (float64, error) {
	Dj := FindDrawdowns(Ra)["draw"]
	Dr := FindDrawdowns(Ra)["length"]
	if len(Dj) <= 0 || len(Dr) <= 0 {
		return math.NaN(), errors.New("In AverageLength, len(Dj) <= 0 || len(Dr) <= 0")
	}
	length_NoneZero := 0.0
	total_Dr := 0.0
	for i, val := range Dj {
		if val < 0 {
			length_NoneZero = length_NoneZero + 1.0
			total_Dr += Dr[i]
		}
	}
	result := math.Abs(total_Dr / length_NoneZero)
	return result, nil
}

func AverageRecovery(Ra *utils.SlidingWindow) (float64, error) {
	Dj := FindDrawdowns(Ra)["draw"]
	Dr := FindDrawdowns(Ra)["recovery"]
	if len(Dj) <= 0 || len(Dr) <= 0 {
		return math.NaN(), errors.New("In AverageLength, len(Dj) <= 0 || len(Dr) <= 0")
	}
	length_NoneZero := 0.0
	total_Dr := 0.0
	for i, val := range Dj {
		if val < 0 {
			length_NoneZero += 1.0
			total_Dr += Dr[i]
		}
	}
	result := math.Abs(total_Dr / length_NoneZero)
	return result, nil
}

/// <summary>
/// calculate a traditional or modified Sharpe Ratio of Return over StdDev or
/// VaR or ES
///
/// The Sharpe ratio is simply the return per unit of risk (represented by
/// variability).  In the classic case, the unit of risk is the standard
/// deviation of the returns.
/// </summary>
func SharpeRatio(Ra *utils.SlidingWindow, Rf_val float64, scale float64) (float64, error) {
	Rf, err := utils.CreateList(Rf_val, Ra.Count())
	if err != nil {
		return math.NaN(), err
	}
	xR, err := Excess(Ra, Rf)
	if err != nil {
		return math.NaN(), err
	}
	numerator := 0.0
	denominator := 0.0
	annualize := 1
	if annualize == 1 {
		denominator, err = StdDev_Annualized(Ra, scale)
		if err != nil {
			return math.NaN(), err
		}
		numerator, err = Annualized(xR, scale, true)
		if err != nil {
			return math.NaN(), err
		}
	} else {
		denominator, err = StdDev(Ra)
		if err != nil {
			return math.NaN(), err
		}
		numerator = xR.Average()
	}

	return numerator / denominator, nil
}

/// <summary>
/// calculate annualized Sharpe Ratio
/// The Sharpe Ratio is a risk-adjusted measure of return that uses standard
/// deviation to represent risk.
/// </summary>
func SharpeRatio_Annualized(Ra *utils.SlidingWindow, Rf float64, scale float64) (float64, error) {
	xR, err := Excess2(Ra, Rf)
	if err != nil {
		return math.NaN(), err
	}
	xR_Ann, err := Annualized(xR, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	std_Ann, err := StdDev_Annualized(Ra, scale)
	if err != nil {
		return math.NaN(), err
	}
	return xR_Ann / std_Ann, nil
}
