package performance

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/bxy09/gfstat/performance/utils"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func ReadFromCSV(fileName string) ([]float64, []float64, []time.Time) {
	data := make([]float64, 0)
	date := make([]time.Time, 0)
	bench_mark := make([]float64, 0)
	//file, err := os.Open("testData/test_DownsideFrequency.csv")
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil, nil
	}
	defer file.Close()
	reader := csv.NewReader(file)
	i := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return nil, nil, nil
		}
		i++
		if i > 1 {
			val := record[4]
			if strings.EqualFold(val, "NA") {
				data = append(data, math.NaN())
			} else {
				tmp, _ := strconv.ParseFloat(val, 64)
				data = append(data, tmp)
			}

			val1 := record[1]
			if strings.EqualFold(val1, "NA") {
				bench_mark = append(bench_mark, math.NaN())
			} else {
				tmp, _ := strconv.ParseFloat(val1, 64)
				bench_mark = append(bench_mark, tmp)
			}

			dateTime, err := time.ParseInLocation("2006-01-02 15:04:05", record[5], time.Local)
			if err != nil {
				fmt.Println("Parse Time Error !!")
				return nil, nil, nil
			}
			date = append(date, dateTime)
		}
	}
	return data, bench_mark, date
}

func Save_CSV_Data(data []float64) {
	fileName := "testData/drawdown.csv"
	buf := new(bytes.Buffer)
	r2 := csv.NewWriter(buf)
	for _, val := range data {
		tmp := make([]string, 1)
		str2 := fmt.Sprintf("%f", val)
		tmp[0] = str2
		r2.Write(tmp)
		r2.Flush()
	}

	fout, err := os.Create(fileName)
	defer fout.Close()
	if err != nil {
		fmt.Println(fileName, err)
		return
	}
	fout.WriteString(buf.String())
}

func Test_Performance(t *testing.T) {
	filename := "testData/555665746e955230ea000001_profit_table.csv"
	data, bench_mark, date := ReadFromCSV(filename)

	Price, _ := utils.NewSlidingWindow(len(data))
	for _, val := range data {
		Price.Add(val)
	}
	Ra, _ := Calculate(Price, "discrete")
	Save_CSV_Data(Ra.Data())
	for key, _ := range PerformanceMap {
		result, err := PerformanceMap[key].Process(data, bench_mark, date)
		fmt.Println(key, " Result: ", result, "Error: ", err)
	}
}

func Test_Same_Data(t *testing.T) {
	filename := "testData/test_Same_Data.csv"
	data, bench_mark, date := ReadFromCSV(filename)
	for key, _ := range PerformanceMap {
		result, err := PerformanceMap[key].Process(data, bench_mark, date)
		fmt.Println(key, " Result: ", result, "Error: ", err)
	}
}
