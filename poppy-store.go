package gin_poppy

import (
	"fmt"
	"sync"
	"time"
)

const (
	GWindowSize = 3600
)

type BaseStore struct {
	Code    map[int]int64 //记录每一种httpcode的数量
	Latency []int64       //分别记录每一次的耗时
}

func NewBaseStore() BaseStore {
	return BaseStore{
		Code:    make(map[int]int64),
		Latency: make([]int64, 0),
	}
}

type PathStore struct {
	rw        sync.RWMutex
	StartTime int64
	LastTime  int64
	Total     int64 //启动以来的请求总数，不受轮转数组限制
	TotalCode map[int]int64

	Items [GWindowSize]BaseStore //轮准数组记录最近一小时的详细数据
}

func NewPathStore() *PathStore {
	store := &PathStore{}
	for i, _ := range store.Items {
		store.Items[i].Code = make(map[int]int64)
		store.Items[i].Latency = make([]int64, 0)
	}
	store.TotalCode = make(map[int]int64)
	return store
}

func (p *PathStore) Add(code int, latency int64) {
	p.rw.Lock()
	defer p.rw.Unlock()
	timestamp := time.Now().Unix()
	pos := timestamp % GWindowSize
	if p.StartTime <= 0 {
		p.StartTime = timestamp
		p.LastTime = timestamp

		p.Items[pos].Latency = append(p.Items[pos].Latency, latency)
		p.Items[pos].Code[code] = 1
		p.Total++
		p.TotalCode[code]++
		return
	}
	if timestamp == p.LastTime {
		p.Items[pos].Latency = append(p.Items[pos].Latency, latency)
		p.Items[pos].Code[code]++
		p.Total++
		p.TotalCode[code]++
		return
	}

	i := p.LastTime + 1
	for i < timestamp {
		p.Items[i%GWindowSize] = NewBaseStore()
		i++
	}

	p.StartTime = max(p.StartTime, timestamp-GWindowSize+1)
	p.LastTime = timestamp
	p.Items[pos].Latency = append(p.Items[pos].Latency, latency)
	p.Items[pos].Code[code]++
	p.Total++
	p.TotalCode[code]++

	return
}

//poppyStore是更上一层的数据抽象，负责存储每个path对应的数据
type PoppyStore struct {
	rw          sync.RWMutex
	Items       map[string]*PathStore
	GlobalStore TotalStat
}

func NewPoppyStore() PoppyStore {
	store := PoppyStore{}
	store.rw = sync.RWMutex{}
	store.Items = make(map[string]*PathStore)
	store.GlobalStore.TotalCodeCount = make(map[int]int64)
	return store
}

func (s *PoppyStore) Add(path string, code int, latency int64) {
	s.GlobalStore.TotalCount++
	s.GlobalStore.TotalCodeCount[code]++
	s.GlobalStore.LastTime = time.Now().Unix()
	if s.GlobalStore.StartTime == 0 {
		s.GlobalStore.StartTime = time.Now().Unix()
	}

	if pathStore, ok := s.Items[path]; ok {
		pathStore.Add(code, latency)
	} else {
		s.rw.Lock()
		s.Items[path] = NewPathStore()
		s.rw.Unlock()
		s.Items[path].Add(code, latency)
	}
}

func (s *PoppyStore) GenerateStat() PoppyStat {
	s.rw.RLock()
	defer s.rw.RUnlock()

	now := time.Now()
	poppyResult := NewPoppyStat()
	for path, store := range s.Items {
		if _, ok := poppyResult.UriResult[path]; !ok {
			poppyResult.UriResult[path] = NewPathStat()
		}
		poppyResult.UriResult[path] = convertStoreToStat(store)
	}

	poppyResult.GlobalResult = s.GlobalStore
	fmt.Println(time.Now().Sub(now))
	return poppyResult
}

func (s *PoppyStore) GenerateRawData() PoppyRawData {
	s.rw.RLock()
	defer s.rw.RUnlock()

	raw := NewPoppyRawData()
	for path, item := range s.Items {
		rawPathData := NewPoppyRawPathData()
		rawPathData.StartTime = item.StartTime
		rawPathData.LastTime = item.LastTime
		rawPathData.TotalCode = item.TotalCode
		rawPathData.Total = item.Total

		var i int64
		for i = 0; i < item.LastTime-item.StartTime+1; i++ {
			index := item.StartTime%GWindowSize + i
			if len(item.Items[index].Latency) == 0 {
				continue
			}

			rawBaseData := RawBaseData{
				Timestamp: item.StartTime + i,
			}
			rawBaseData.Latency = append(rawBaseData.BaseStore.Latency, item.Items[index].Latency...)
			rawBaseData.Code = item.Items[index].Code

			rawPathData.Items = append(rawPathData.Items, rawBaseData)
		}
		raw.UriRawData[path] = rawPathData

	}
	fmt.Println(raw)

	return raw
}
