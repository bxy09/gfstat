// 是在未来将会废弃的基于天收益序列计算指标的模块，由于大量不恰当的使用了sliding window，导致性能不佳
package performance

import (
	"errors"
	"github.com/bxy09/gfstat/performance/utils"
	"math"
	"time"
)

type Func1Sliding func(Ra *utils.SlidingWindow) (float64, error)
type Func1Sliding1F func(Ra *utils.SlidingWindow, param float64) (float64, error)
type Func1Sliding1F1B func(Ra *utils.SlidingWindow, param float64, flag bool) (float64, error)
type Func1Sliding2F func(Ra *utils.SlidingWindow, param1 float64, param2 float64) (float64, error)
type Func1Sliding1F1S func(Ra *utils.SlidingWindow, param float64, str string) (float64, error)

type Func2Sliding func(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow) (float64, error)
type Func2Sliding1F func(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, param float64) (float64, error)
type Func2Sliding2F func(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, param1 float64, param2 float64) (float64, error)
type Func2Sliding2F1S func(Ra *utils.SlidingWindow, Rb *utils.SlidingWindow, param1 float64, param2 float64, str string) (float64, error)

func getPeriod(recordDate []time.Time) float64 {
	//to estimate the period
	count_OneMinute := 0
	count_OneDay := 0
	count_OneMonth := 0
	count_OneSeason := 0
	for i := 1; i < len(recordDate); i++ {
		diff := recordDate[i].Sub(recordDate[i-1])
		//fmt.Println(diff.Hours(), diff.Minutes(), diff.Seconds())
		if diff.Hours() >= 23.0 && diff.Hours() < 48.0 {
			count_OneDay++
		}
		if diff.Minutes() >= 1.0 && diff.Minutes() <= 5.0 {
			count_OneMinute++
		}
		if diff.Hours() >= 28.0*24.0 && diff.Hours() <= 35.0*24.0 {
			count_OneMonth++
		}
		if diff.Hours() >= 3.8*30.0*24.0 && diff.Hours() <= 4.2*30.0*24.0 {
			count_OneSeason++
		}
	}
	//count the percent of the total length
	len_Date := len(recordDate) - 1
	Minute_Ratio := float64(count_OneMinute) / float64(len_Date)
	Day_Ratio := float64(count_OneDay) / float64(len_Date)
	Month_Ratio := float64(count_OneMonth) / float64(len_Date)
	Season_Ratio := float64(count_OneSeason) / float64(len_Date)
	if Minute_Ratio > Day_Ratio && Minute_Ratio > Month_Ratio && Minute_Ratio > Season_Ratio {
		return 2520.0
	}
	if Day_Ratio > Minute_Ratio && Day_Ratio > Month_Ratio && Day_Ratio > Season_Ratio {
		return 252.0
	}
	if Month_Ratio > Minute_Ratio && Month_Ratio > Day_Ratio && Month_Ratio > Season_Ratio {
		return 12.0
	}
	if Season_Ratio > Minute_Ratio && Season_Ratio > Day_Ratio && Season_Ratio > Month_Ratio {
		return 4.0
	}
	return 252.0
}

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

type Performance interface {
	Process(AssetPriceReturns, AssetPriceBenchMark []float64, date []time.Time) (float64, error)
}

type P1SWrapper struct {
	function Func1Sliding
}

type P1S1FWrapper struct {
	function Func1Sliding1F
	param    float64
}

type P1S1F1BWrapper struct {
	function Func1Sliding1F1B
	param    float64
	flag     bool
}

type P1S2FWrapper struct {
	function Func1Sliding2F
	param1   float64
	param2   float64
}

type P1S1F1SWrapper struct {
	function Func1Sliding1F1S
	param    float64
	str      string
}

func (this *P1SWrapper) Process(AssetPriceReturns, AssetPriceBenchMark []float64, date []time.Time) (float64, error) {
	if AssetPriceReturns == nil {
		return math.NaN(), errors.New("The Input RA is Error !!!")
	}
	var err error
	Period := getPeriod(date)
	if Period == 2520 {
		AssetPriceReturns, err = reorganizeInputPrice(date, AssetPriceReturns)
		if err != nil {
			return math.NaN(), errors.New("Reorganize Minutes Price Error !!!")
		}
	}

	Price, err := utils.NewSlidingWindow(len(AssetPriceReturns))
	if err != nil {
		return math.NaN(), err
	}
	for _, val := range AssetPriceReturns {
		Price.Add(val)
	}
	Ra, err := Calculate(Price, "discrete")
	if err != nil {
		return math.NaN(), err
	}
	return this.function(Ra)
}

func (this *P1S1FWrapper) Process(AssetPriceReturns, AssetPriceBenchMark []float64, date []time.Time) (float64, error) {
	if AssetPriceReturns == nil {
		return math.NaN(), errors.New("The Input RA is Error !!!")
	}
	var err error
	Period := getPeriod(date)
	if Period == 2520 {
		AssetPriceReturns, err = reorganizeInputPrice(date, AssetPriceReturns)
		if err != nil {
			return math.NaN(), errors.New("Reorganize Minutes Price Error !!!")
		}
		Period = 252
	}
	if this.param == 252 {
		this.param = Period
	}

	Price, err := utils.NewSlidingWindow(len(AssetPriceReturns))
	if err != nil {
		return math.NaN(), err
	}
	for _, val := range AssetPriceReturns {
		Price.Add(val)
	}
	Ra, err := Calculate(Price, "discrete")
	if err != nil {
		return math.NaN(), err
	}
	if this.param == 0.03 {
		this.param = this.param / Period
	}
	return this.function(Ra, this.param)
}

func (this *P1S1F1BWrapper) Process(AssetPriceReturns, AssetPriceBenchMark []float64, date []time.Time) (float64, error) {
	if AssetPriceReturns == nil {
		return math.NaN(), errors.New("The Input RA is Error !!!")
	}
	var err error
	Period := getPeriod(date)
	if Period == 2520 {
		AssetPriceReturns, err = reorganizeInputPrice(date, AssetPriceReturns)
		if err != nil {
			return math.NaN(), errors.New("Reorganize Minutes Price Error !!!")
		}
		Period = 252
	}
	if this.param == 252 {
		this.param = Period
	}

	Price, err := utils.NewSlidingWindow(len(AssetPriceReturns))
	if err != nil {
		return math.NaN(), err
	}
	for _, val := range AssetPriceReturns {
		Price.Add(val)
	}
	Ra, err := Calculate(Price, "discrete")
	if err != nil {
		return math.NaN(), err
	}
	if this.param == 0.03 {
		this.param = this.param / Period
	}
	return this.function(Ra, this.param, this.flag)
}

func (this *P1S2FWrapper) Process(AssetPriceReturns, AssetPriceBenchMark []float64, date []time.Time) (float64, error) {
	if AssetPriceReturns == nil {
		return math.NaN(), errors.New("The Input RA is Error !!!")
	}

	var err error
	Period := getPeriod(date)
	if Period == 2520 {
		AssetPriceReturns, err = reorganizeInputPrice(date, AssetPriceReturns)
		if err != nil {
			return math.NaN(), errors.New("Reorganize Minutes Price Error !!!")
		}
		Period = 252
	}

	if this.param2 == 252 {
		this.param2 = Period
	}

	if this.param1 == 252 {
		this.param1 = Period
	}

	Price, err := utils.NewSlidingWindow(len(AssetPriceReturns))
	if err != nil {
		return math.NaN(), err
	}
	for _, val := range AssetPriceReturns {
		Price.Add(val)
	}
	Ra, err := Calculate(Price, "discrete")
	if err != nil {
		return math.NaN(), err
	}

	if this.param1 == 0.03 {
		this.param1 = this.param1 / Period
	}

	if this.param2 == 0.03 {
		this.param2 = this.param2 / Period
	}
	return this.function(Ra, this.param1, this.param2)
}

func (this *P1S1F1SWrapper) Process(AssetPriceReturns, AssetPriceBenchMark []float64, date []time.Time) (float64, error) {
	if AssetPriceReturns == nil {
		return math.NaN(), errors.New("The Input RA is Error !!!")
	}
	var err error
	Period := getPeriod(date)
	if Period == 2520 {
		AssetPriceReturns, err = reorganizeInputPrice(date, AssetPriceReturns)
		if err != nil {
			return math.NaN(), errors.New("Reorganize Minutes Price Error !!!")
		}
		Period = 252
	}

	if this.param == 252 {
		this.param = Period
	}

	Price, err := utils.NewSlidingWindow(len(AssetPriceReturns))
	if err != nil {
		return math.NaN(), err
	}
	for _, val := range AssetPriceReturns {
		Price.Add(val)
	}
	Ra, err := Calculate(Price, "discrete")
	if err != nil {
		return math.NaN(), err
	}
	if this.param == 0.03 {
		this.param = this.param / Period
	}
	return this.function(Ra, this.param, this.str)
}

type P2SWrapper struct {
	function Func2Sliding
}

type P2S1FWrapper struct {
	function Func2Sliding1F
	param    float64
}

type P2S2FWrapper struct {
	function Func2Sliding2F
	param1   float64
	param2   float64
}

type P2S2F1SWrapper struct {
	function Func2Sliding2F1S
	param1   float64
	param2   float64
	str      string
}

func (this *P2SWrapper) Process(AssetPriceReturns, AssetPriceBenchMark []float64, date []time.Time) (float64, error) {
	if AssetPriceReturns == nil || AssetPriceBenchMark == nil {
		return math.NaN(), errors.New("The Input RA or RB are Error !!!")
	}

	var err error
	Period := getPeriod(date)
	if Period == 2520 {
		AssetPriceReturns, err = reorganizeInputPrice(date, AssetPriceReturns)
		AssetPriceBenchMark, err = reorganizeInputPrice(date, AssetPriceBenchMark)
		if err != nil {
			return math.NaN(), errors.New("Reorganize Minutes Price Error !!!")
		}
		Period = 252
	}

	Price, err := utils.NewSlidingWindow(len(AssetPriceReturns))
	if err != nil {
		return math.NaN(), err
	}
	for _, val := range AssetPriceReturns {
		Price.Add(val)
	}
	Ra, err := Calculate(Price, "discrete")
	if err != nil {
		return math.NaN(), err
	}

	Price_Bench, err := utils.NewSlidingWindow(len(AssetPriceBenchMark))
	if err != nil {
		return math.NaN(), err
	}
	for _, val := range AssetPriceBenchMark {
		Price_Bench.Add(val)
	}
	Rb, err := Calculate(Price_Bench, "discrete")
	if err != nil {
		return math.NaN(), err
	}
	return this.function(Ra, Rb)
}

func (this *P2S1FWrapper) Process(AssetPriceReturns, AssetPriceBenchMark []float64, date []time.Time) (float64, error) {
	if AssetPriceReturns == nil || AssetPriceBenchMark == nil {
		return math.NaN(), errors.New("The Input RA or RB are Error !!!")
	}

	var err error
	Period := getPeriod(date)
	if Period == 2520 {
		AssetPriceReturns, err = reorganizeInputPrice(date, AssetPriceReturns)
		AssetPriceBenchMark, err = reorganizeInputPrice(date, AssetPriceBenchMark)
		if err != nil {
			return math.NaN(), errors.New("Reorganize Minutes Price Error !!!")
		}
		Period = 252
	}

	if this.param == 252 {
		this.param = Period
	}
	if this.param == 0.03 {
		this.param = this.param / Period
	}

	Price, err := utils.NewSlidingWindow(len(AssetPriceReturns))
	if err != nil {
		return math.NaN(), err
	}
	for _, val := range AssetPriceReturns {
		Price.Add(val)
	}
	Ra, err := Calculate(Price, "discrete")
	if err != nil {
		return math.NaN(), err
	}

	Price_Bench, err := utils.NewSlidingWindow(len(AssetPriceBenchMark))
	if err != nil {
		return math.NaN(), err
	}
	for _, val := range AssetPriceBenchMark {
		Price_Bench.Add(val)
	}
	Rb, err := Calculate(Price_Bench, "discrete")
	if err != nil {
		return math.NaN(), err
	}

	return this.function(Ra, Rb, this.param)
}

func (this *P2S2FWrapper) Process(AssetPriceReturns, AssetPriceBenchMark []float64, date []time.Time) (float64, error) {
	if AssetPriceReturns == nil || AssetPriceBenchMark == nil {
		return math.NaN(), errors.New("The Input RA or RB are Error !!!")
	}

	var err error
	Period := getPeriod(date)
	if Period == 2520 {
		AssetPriceReturns, err = reorganizeInputPrice(date, AssetPriceReturns)
		AssetPriceBenchMark, err = reorganizeInputPrice(date, AssetPriceBenchMark)
		if err != nil {
			return math.NaN(), errors.New("Reorganize Minutes Price Error !!!")
		}
		Period = 252
	}

	if this.param1 == 252 {
		this.param1 = Period
	}
	if this.param2 == 252 {
		this.param2 = Period
	}

	if this.param1 == 0.03 {
		this.param1 = this.param1 / Period
	}
	if this.param2 == 0.03 {
		this.param2 = this.param2 / Period
	}

	Price, err := utils.NewSlidingWindow(len(AssetPriceReturns))
	if err != nil {
		return math.NaN(), err
	}
	for _, val := range AssetPriceReturns {
		Price.Add(val)
	}
	Ra, err := Calculate(Price, "discrete")
	if err != nil {
		return math.NaN(), err
	}

	Price_Bench, err := utils.NewSlidingWindow(len(AssetPriceBenchMark))
	if err != nil {
		return math.NaN(), err
	}
	for _, val := range AssetPriceBenchMark {
		Price_Bench.Add(val)
	}
	Rb, err := Calculate(Price_Bench, "discrete")
	if err != nil {
		return math.NaN(), err
	}

	return this.function(Ra, Rb, this.param1, this.param2)
}

func (this *P2S2F1SWrapper) Process(AssetPriceReturns, AssetPriceBenchMark []float64, date []time.Time) (float64, error) {
	if AssetPriceReturns == nil || AssetPriceBenchMark == nil {
		return math.NaN(), errors.New("The Input RA or RB are Error !!!")
	}

	var err error
	Period := getPeriod(date)
	if Period == 2520 {
		AssetPriceReturns, err = reorganizeInputPrice(date, AssetPriceReturns)
		AssetPriceBenchMark, err = reorganizeInputPrice(date, AssetPriceBenchMark)
		if err != nil {
			return math.NaN(), errors.New("Reorganize Minutes Price Error !!!")
		}
		Period = 252
	}
	if this.param1 == 252 {
		this.param1 = Period
	}
	if this.param2 == 252 {
		this.param2 = Period
	}

	Price, err := utils.NewSlidingWindow(len(AssetPriceReturns))
	if err != nil {
		return math.NaN(), err
	}
	for _, val := range AssetPriceReturns {
		Price.Add(val)
	}
	Ra, err := Calculate(Price, "discrete")
	if err != nil {
		return math.NaN(), err
	}

	Price_Bench, err := utils.NewSlidingWindow(len(AssetPriceBenchMark))
	if err != nil {
		return math.NaN(), err
	}
	for _, val := range AssetPriceBenchMark {
		Price_Bench.Add(val)
	}
	Rb, err := Calculate(Price_Bench, "discrete")
	if err != nil {
		return math.NaN(), err
	}

	if this.param1 == 0.03 {
		this.param1 = this.param1 / Period
	}
	if this.param2 == 0.03 {
		this.param2 = this.param2 / Period
	}

	return this.function(Ra, Rb, this.param1, this.param2, this.str)
}

var PerformanceMap map[string]Performance

func init() {
	/*
		Params:scale
			number of periods in a year (daily scale = Scale, monthly scale = 12, quarterly scale = 4)
		Params:MAR
			MAR=0.03
		Params:Rf
			Rf=0.03
	*/
	MAR := 0.03 //---->scale 252 MAR=MAR/Scale
	Rf := 0.03
	Scale := 252.0

	PerformanceMap = map[string]Performance{
		//One utils.SlidingWindow
		"BernardoLedoitRatio":   &P1SWrapper{BernardoLedoitRatio},   //
		"DRatio":                &P1SWrapper{DRatio},                //
		"MeanAbsoluteDeviation": &P1SWrapper{MeanAbsoluteDeviation}, //
		"SkewnessKurtosisRatio": &P1SWrapper{SkewnessKurtosisRatio}, //
		"MeanGeometric":         &P1SWrapper{MeanGeometric},         //
		"Variance":              &P1SWrapper{Variance},              //
		"StdDev":                &P1SWrapper{StdDev},                //
		"MaxDrawdown":           &P1SWrapper{MaxDrawdown},           //
		"Skewness":              &P1SWrapper{Skewness},              //
		"Kurtosis":              &P1SWrapper{Kurtosis},              //
		"PainIndex":             &P1SWrapper{PainIndex},             //
		"AverageDrawdown":       &P1SWrapper{AverageDrawdown},       //
		"AverageLength":         &P1SWrapper{AverageLength},         //
		"AverageRecovery":       &P1SWrapper{AverageRecovery},       //

		"StdDev_Annualized":           &P1S1FWrapper{StdDev_Annualized, Scale},         //
		"SortinoRatio":                &P1S1FWrapper{SortinoRatio, MAR},                //
		"ProspectRatio":               &P1S1FWrapper{ProspectRatio, MAR},               //
		"DownsideFrequency2":          &P1S1FWrapper{DownsideFrequency2, MAR},          //
		"UpsideFrequency":             &P1S1FWrapper{UpsideFrequency, MAR},             //
		"DownsideDeviation2":          &P1S1FWrapper{DownsideDeviation2, MAR},          //
		"KellyRatio_Full":             &P1S1FWrapper{KellyRatio_Full, MAR},             //
		"KellyRatio_Half":             &P1S1FWrapper{KellyRatio_Half, MAR},             //
		"UpsidePotentialRatio":        &P1S1FWrapper{UpsidePotentialRatio, MAR},        //
		"VolatilitySkewness_Variance": &P1S1FWrapper{VolatilitySkewness_Variance, MAR}, //
		"VolatilitySkewness_Risk":     &P1S1FWrapper{VolatilitySkewness_Risk, MAR},     //
		"CalmarRatio":                 &P1S1FWrapper{CalmarRatio, Scale},               //

		"Annualized": &P1S1F1BWrapper{Annualized, Scale, true}, //

		"AdjustedSharpeRatio":    &P1S2FWrapper{AdjustedSharpeRatio, MAR, Scale},   //
		"BurkeRatio":             &P1S2FWrapper{BurkeRatio, MAR, Scale},            //
		"Kappa":                  &P1S2FWrapper{Kappa, MAR, 1.0},                   //
		"SterlingRatio":          &P1S2FWrapper{SterlingRatio, Scale, 0.1},         //
		"PainRatio":              &P1S2FWrapper{PainRatio, Rf, Scale},              //
		"SharpeRatio":            &P1S2FWrapper{SharpeRatio, Rf, Scale},            //
		"SharpeRatio_Annualized": &P1S2FWrapper{SharpeRatio_Annualized, Rf, Scale}, //0.63747766

		"UpsideRisk": &P1S1F1SWrapper{UpsideRisk, MAR, "risk"}, //

		//Two utils.SlidingWindow
		"UpDownRatios": &P2SWrapper{UpDownRatios}, //

		"ActivePremium":    &P2S1FWrapper{ActivePremium, Scale},    //
		"TrackingError":    &P2S1FWrapper{TrackingError, Scale},    //
		"InformationRatio": &P2S1FWrapper{InformationRatio, Scale}, //

		"M2Sortino":         &P2S2FWrapper{M2Sortino, Scale, MAR},         //
		"FamaBeta":          &P2S2FWrapper{FamaBeta, Scale, Scale},        //
		"SpecificRisk":      &P2S2FWrapper{SpecificRisk, Scale, Rf},       //
		"SystematicRisk":    &P2S2FWrapper{SystematicRisk, Scale, Rf},     //
		"TotalRisk":         &P2S2FWrapper{TotalRisk, Scale, Rf},          //
		"TreynorRatio":      &P2S2FWrapper{TreynorRatio, Scale, Rf},       //
		"OmegaExcessReturn": &P2S2FWrapper{OmegaExcessReturn, Scale, MAR}, //
		"MSquared":          &P2S2FWrapper{MSquared, Scale, Rf},           //
		"JensenAlpha2":      &P2S2FWrapper{JensenAlpha2, Rf, Scale},       //

		"AppraisalRatio": &P2S2F1SWrapper{AppraisalRatio, Scale, Rf, "modified"},   //
		"MSquaredExcess": &P2S2F1SWrapper{MSquaredExcess, Scale, Rf, "arithmetic"}, //
	}
}
