package performance

import (
	"github.com/GaryBoone/GoStats/stats"
	"github.com/bxy09/gfstat/performance/utils"
	"math"
)

/// <summary>
/// 资本资产定价模型
/// </summary>

/// <summary>
/// 阿尔法
/// </summary>
/// <param name="Ra"></param>
/// <param name="Rb"></param>
/// <param name="Rf"></param>
/// <returns></returns>????
func Alpha(Ra, Rb, Rf *utils.SlidingWindow) (float64, error) {
	xRb, err := Excess(Rb, Rf)
	if err != nil {
		return math.NaN(), err
	}
	xRa, err := Excess(Ra, Rf)
	if err != nil {
		return math.NaN(), err
	}
	var _, intercept, _, _, _, _ = stats.LinearRegression(xRb.Data(), xRa.Data())
	return intercept, nil
}

func Alpha2(Ra, Rb *utils.SlidingWindow, Rf float64) (float64, error) {
	RfList, err := utils.CreateList(Rf, Ra.Count())
	if err != nil {
		return math.NaN(), err
	}
	return Alpha(Ra, Rb, RfList)
}

func Beta(Ra, Rb, Rf *utils.SlidingWindow) (float64, error) {
	xRb, err := Excess(Rb, Rf)
	if err != nil {
		return math.NaN(), err
	}
	xRa, err := Excess(Ra, Rf)
	if err != nil {
		return math.NaN(), err
	}
	var slope, _, _, _, _, _ = stats.LinearRegression(xRb.Data(), xRa.Data())
	return slope, nil
}

func Beta2(Ra, Rb *utils.SlidingWindow, Rf float64) (float64, error) {
	RfList, err := utils.CreateList(Rf, Ra.Count())
	if err != nil {
		return math.NaN(), err
	}
	return Beta(Ra, Rb, RfList)
}

/// <summary>
/// Epsilon
/// </summary>
func Epsilon(Ra, Rb, Rf *utils.SlidingWindow, scale float64) (float64, error) {
	Rpf, err := Annualized(Rf, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	return Epsilon2(Ra, Rb, Rpf, scale)
}
func Epsilon2(Ra, Rb *utils.SlidingWindow, Rf float64, scale float64) (float64, error) {

	alpha, err := Alpha2(Ra, Rb, Rf)
	if err != nil {
		return math.NaN(), err
	}
	beta, err := Beta2(Ra, Rb, Rf)
	if err != nil {
		return math.NaN(), err
	}
	Rpa, err := Annualized(Ra, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	Rpb, err := Annualized(Rb, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	return Rpa - Rf - alpha - beta*(Rpb-Rf), nil
}

/// <summary>
/// The Jensen's alpha is the intercept of the regression equation in the Capital
/// Asset Pricing Model and is in effect the exess return adjusted for systematic risk.
/// （与alpha相比的区别？年化数值，包含残差的影响）
/// </summary>
/// <param name="Ra"></param>
/// <param name="Rb"></param>
/// <param name="Rf"></param>
/// <param name="scale"></param>
/// <returns></returns>

func JensenAlpha(Ra, Rb, Rf *utils.SlidingWindow, scale float64) (float64, error) {
	Rpf, err := Annualized(Rf, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	return JensenAlpha2(Ra, Rb, Rpf, scale)
}

func JensenAlpha2(Ra, Rb *utils.SlidingWindow, Rf float64, scale float64) (float64, error) {
	beta, err := Beta2(Ra, Rb, Rf)
	if err != nil {
		return math.NaN(), err
	}
	Rpa, err := Annualized(Ra, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	Rpb, err := Annualized(Rb, scale, true)
	if err != nil {
		return math.NaN(), err
	}
	Rf = Rf * scale
	return Rpa - Rf - beta*(Rpb-Rf), nil
}

/// <summary>
/// 斜率
/// </summary>
/// <param name="Rb"></param>
/// <param name="Rf"></param>
/// <returns></returns>
func CMLSlope2(Rb *utils.SlidingWindow, Rf float64) (float64, error) {
	return CMLSlope(Rb, Rf)
}
func CMLSlope(Rb *utils.SlidingWindow, Rf float64) (float64, error) {
	return SharpeRatio(Rb, Rf, 1)
}
