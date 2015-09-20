package metric_test

import (
	"math"
	"math/rand"
	"testing"
	"time"
	"github.com/bxy09/gfstat/metric"
	"github.com/bxy09/gfstat/performance"
)

func TestIdentity(t *testing.T) {
	for _, length := range []int{0, 1, 2, 3, 4, 5, 10, 100, 500, 1000, 2000} {
		t.Log("length=", length)
		dates := make([]time.Time, length)
		now := time.Now().Add(-time.Duration(length) * 24 * time.Hour)
		assets := make([]float64, length)
		bench := make([]float64, length)
		var v1, v2 float64
		v1 = 10000.0
		v2 = 100.0
		for i := range dates {
			dates[i] = now
			now = now.Add(time.Hour * 24)
			assets[i] = v1
			bench[i] = v2
			v1 += 0.2 * v1 * (rand.Float64() - 0.5)
			v2 += 0.2 * v2 * (rand.Float64() - 0.5)
		}

		caculator := metric.NewMetricCalculator(assets, bench, dates)
		for key, _ := range metric.MetricMap {
			var d1, d2 time.Duration
			var now = time.Now()
			pr, err := performance.PerformanceMap[key].Process(assets, bench, dates)
			if err != nil && length > 4 {
				t.Fatal("key", key, err)
			}
			d1 = time.Now().Sub(now)
			mr, err := caculator.Process(key)
			if err != nil && length > 4 {
				t.Fatal("key", key, err)
			}
			if math.Abs(pr-mr) > 0.000001 && length > 4 {
				t.Fatal("not identity", key, pr, mr)
			}
			d2 = time.Now().Sub(now) - d1
			t.Logf("Check OK, key:%s, performance:%f(%s), metric:%f(%s) ", key, pr, d1, mr, d2)
		}
	}
}
