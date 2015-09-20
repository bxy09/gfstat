package utils

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestSlidingwindow(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	size := 5
	n := 3
	temp := n - size
	if temp < 0 {
		temp = 0
	}
	total := 0
	test, err := NewSlidingWindow(size)
	if err != nil {
		t.Fatalf("Create NewSliding Window is Failed!!!")
	}

	for i := 0; i < n; i++ {
		test.Add(float64(i))
		if i >= temp {
			total += i
		}
	}

	data := test.Data()

	for i, d := range data {
		if round(d) != temp+i {
			t.Fatalf("bad data in window: %v", data)
		}
	}
	max, maxposition := test.Max()
	min, minposition := test.Min()
	if round(max) != n-1 || maxposition != test.Count()-1 {
		fmt.Printf("%d, %v", test.index, data)
		t.Fatalf("bad sliding window, n = %d, size = %d, max: %f, %d, expected: %d, %d", test.Count(), size, max, maxposition, n-1, n-1)
	}
	if round(min) != temp || minposition != 0 {
		t.Fatalf("bad sliding window for min: %f, %d", min, minposition)
	}
	if round(test.First()) != temp {
		t.Fatalf("bad sliding window for First: %f", test.First())
	}
	fmt.Printf("total = %d\n", total)
	if test.Average() != float64(total)/float64(test.Count()) {
		t.Fatalf("bad sliding window for Average: %f", test.Average())
	}
	if test.Sum() != float64(total) {
		t.Fatalf("bad sliding window for Average: %f", test.Sum())
	}
	fmt.Printf("diffave = %.3f\n", test.DiffAbsAve())
	if test.RefData(1) != float64(n-1) {
		t.Fatalf("bad sliding window for RefData: %f", test.RefData(1))
	}
}

func round(f float64) int {
	return int(math.Floor(f + 0.5))
}
