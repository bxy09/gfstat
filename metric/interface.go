// 提供了一系列基于天收益序列的指标
package metric

import (
	"errors"
	"math"
	"time"
)

type MetricCalculator struct {
	portfolio, bench []float64
	dates            []time.Time
	vectorCache      map[string][]float64
	scalarCache      map[string]float64
	period           float64
}

func NewMetricCalculator(portfolio, bench []float64, dates []time.Time) *MetricCalculator {
	return &MetricCalculator{
		portfolio:   portfolio,
		bench:       bench,
		dates:       dates,
		vectorCache: map[string][]float64{},
		scalarCache: map[string]float64{},
	}
}

func NewDailyMetricCalculatorNoBench(portfolio []float64) *MetricCalculator {
	calculator := NewMetricCalculator(portfolio, nil, nil)
	calculator.period = 252.0
	return calculator
}

func (m *MetricCalculator) Period() float64 {
	if m.period > 0.01 {
		return m.period
	}
	recordDate := m.dates
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
		m.period = 2520.0
		return m.period
	}
	if Day_Ratio > Minute_Ratio && Day_Ratio > Month_Ratio && Day_Ratio > Season_Ratio {
		m.period = 252.0
		return m.period
	}
	if Month_Ratio > Minute_Ratio && Month_Ratio > Day_Ratio && Month_Ratio > Season_Ratio {
		m.period = 12.0
		return m.period
	}
	if Season_Ratio > Minute_Ratio && Season_Ratio > Day_Ratio && Season_Ratio > Month_Ratio {
		m.period = 4.0
		return m.period
	}
	m.period = 252.0
	return m.period
}

func (c MetricCalculator) PortfolioRatio() []float64 {
	key := "PortfolioRatio"
	if c.vectorCache[key] != nil {
		return c.vectorCache[key]
	}
	c.vectorCache[key] = Vector(c.portfolio).ReturnRatio("discrete")
	return c.vectorCache[key]
}

func (c MetricCalculator) BenchRatio() []float64 {
	key := "BenchRatio"
	if c.vectorCache[key] != nil {
		return c.vectorCache[key]
	}
	c.vectorCache[key] = Vector(c.bench).ReturnRatio("discrete")
	return c.vectorCache[key]
}

func (c MetricCalculator) GetOrSetScalar(name string, process func() (float64, error)) (float64, error) {
	if value, exist := c.scalarCache[name]; exist {
		return value, nil
	}
	value, err := process()
	c.scalarCache[name] = value
	return value, err
}

func (c MetricCalculator) GetOrSetVector(name string, process func() (Vector, error)) (Vector, error) {
	if vector, exist := c.vectorCache[name]; exist {
		return vector, nil
	}
	vector, err := process()
	c.vectorCache[name] = vector
	return vector, err
}

func (c MetricCalculator) Process(name string) (float64, error) {
	metric, exist := MetricMap[name]
	if !exist {
		return math.NaN(), errors.New("No such metric")
	}
	return metric(c)
}

type Metric func(c MetricCalculator) (float64, error)

var MetricMap = map[string]Metric{}

//	Params:scale
//		number of periods in a year (daily scale = Scale, monthly scale = 12, quarterly scale = 4)
//	Params:MAR
//		MAR=0.03
//	Params:Rf
//		Rf=0.03

var MAR = 0.03 //---->scale 252 MAR=MAR/Scale
var Rf = 0.03
var Scale = 252.0
