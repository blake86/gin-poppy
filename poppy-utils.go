package gin_poppy

import (
	"math"
	"sort"
	"time"
)

func checkAndGetCurTime(curTime int64) int64 {
	if curTime <= 0 {
		return time.Now().Unix()
	}
	return curTime
}

func max(left, right int64) int64 {
	if left > right {
		return left
	}
	return right
}
func min(left, right int64) int64 {
	if left < right {
		return left
	}
	return right
}

func average(orderedList []int64, l int) float64 {
	if l == 0 {
		return 0
	}
	var sum int64 = 0
	for i := 0; i < l; i++ {
		sum += orderedList[i]
	}
	return float64(sum) / float64(l)
}

func mean(orderedList []int64, l int) float64 {
	if l == 0 {
		return 0
	}
	var res int64 = 0
	for i := 0; i < l; i++ {
		res += orderedList[i]
	}

	return float64(res) / float64(l)
}

func p90(orderedList []int64, l int) int64 {
	return percentile(orderedList, l, 0.9)
}

func p95(orderedList []int64, l int) int64 {
	return percentile(orderedList, l, 0.95)
}

func p99(orderedList []int64, l int) int64 {
	return percentile(orderedList, l, 0.99)
}

func percentile(orderedList []int64, l int, p float64) int64 {
	return orderedList[int(p*float64(l))]
}

func stdev(orderedList []int64, mean float64, l int) float64 {
	if l == 1 {
		return 0
	}
	var omega float64
	for i := 0; i < l; i++ {
		omega += math.Pow(float64(orderedList[i])-mean, 2)
	}
	stdev := math.Sqrt(1 / (float64(l) - 1) * omega)
	return stdev
}

func convertStoreToStat(store *PathStore) PathStat {
	result := NewPathStat()

	latencies := make([]int64, 0)
	codeCount := make(map[int]int64)

	var count int64 = 0
	for _, item := range store.Items {
		for _, latencySlice := range item.Latency {
			latencies = append(latencies, latencySlice)
		}
		for code, c := range item.Code {
			codeCount[code] += c
			count += c
		}
	}

	baseStat := BaseStat{}
	baseStat.TotalCount = store.Total
	baseStat.Count = count

	sortedSlice := latencies[:]

	sort.Slice(sortedSlice, func(i, j int) bool { return sortedSlice[i] < sortedSlice[j] })

	length := len(sortedSlice)
	if length >= 0 {
		baseStat.Min = sortedSlice[0]
		baseStat.Max = sortedSlice[baseStat.Count-1]
		baseStat.Mean = mean(sortedSlice, length)
		baseStat.Average = average(sortedSlice, length)
		baseStat.Stdev = stdev(sortedSlice, baseStat.Mean, length)
		baseStat.P90 = p90(sortedSlice, length)
		baseStat.P95 = p95(sortedSlice, length)
		baseStat.P99 = p99(sortedSlice, length)
		baseStat.StartTime = store.StartTime
		baseStat.LastTime = store.LastTime
	}

	result.Result = baseStat
	result.Result.CodeCount = codeCount

	return result
}
