package gin_poppy

type TotalStat struct {
	StartTime      int64
	LastTime       int64
	TotalCount     int64 //总数
	TotalCodeCount map[int]int64
}

type BaseStat struct {
	TotalStat
	Count     int64
	CodeCount map[int]int64 //http状态码统计
	Min       int64
	Max       int64
	Average   float64
	Mean      float64
	Stdev     float64
	P90       int64
	P95       int64
	P99       int64
	StartTime int64
	LastTime  int64
}

type PathStat struct {
	Result BaseStat
}

func NewPathStat() PathStat {
	return PathStat{}
}

type PoppyStat struct {
	UriResult    map[string]PathStat
	GlobalResult TotalStat
}

func NewPoppyStat() PoppyStat {
	return PoppyStat{
		UriResult: make(map[string]PathStat),
	}
}

type RawBaseData struct {
	Timestamp int64
	BaseStore
}

type PoppyRawPathData struct {
	StartTime int64
	LastTime  int64
	Total     int64 //启动以来的请求总数，不受轮转数组限制
	TotalCode map[int]int64
	Items     []RawBaseData
}

func NewPoppyRawPathData() PoppyRawPathData {
	return PoppyRawPathData{
		Items: make([]RawBaseData, 0),
	}
}

type PoppyRawData struct {
	UriRawData map[string]PoppyRawPathData
}

func NewPoppyRawData() PoppyRawData {
	return PoppyRawData{
		UriRawData: make(map[string]PoppyRawPathData),
	}
}
