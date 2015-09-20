package utils

import (
	"math"
)

type MapSorter []Item

type Item struct {
	Key string
	Val float64
}

func NewMapSorter(m map[string]float64) MapSorter {
	ms := make(MapSorter, 0, len(m))
	for k, v := range m {
		ms = append(ms, Item{k, v})
	}
	return ms
}

func (ms MapSorter) Len() int {
	return len(ms)
}

//按值递增排序，加入NaN处理
func (ms MapSorter) Less(i, j int) bool {
	if math.IsNaN(ms[i].Val) {
		return false
	}
	if math.IsNaN(ms[j].Val) {
		return true
	}
	return ms[i].Val > ms[j].Val
}

func (ms MapSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}
