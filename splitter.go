package main

import (
	"fmt"
	"sort"
)

////////////////////////////////////

const (
	None SpotType = iota
	Start
	End
	Mixed
)

////////////////////////////////////

type Uint64Slice []uint64
type Interval struct {
	Low    int64
	High   int64
	NodeId uint64
	Data   Uint64Slice
}
type IntervalList []*Interval
type SpotType int
type Spot struct {
	Type   SpotType
	Pos    int64
	NodeId uint64
	Data   Uint64Slice
}
type SpotList []*Spot

type SpotMap map[int64]*Spot
type CurrentData map[uint64]bool

type RangeSplit struct {
	Intervals    IntervalList
	Spots        SpotList
	StartPos     int64
	CurrentPos   int64
	CurrentSpots SpotList
	CurrentData  CurrentData
}

////////////////////////////////
// for Sorting items
func (il IntervalList) Swap(i, j int)      { il[i], il[j] = il[j], il[i] }
func (il IntervalList) Len() int           { return len(il) }
func (il IntervalList) Less(i, j int) bool { return il[i].Low < il[j].Low }

func (p Uint64Slice) Len() int           { return len(p) }
func (p Uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (sl SpotList) Swap(i, j int) { sl[i], sl[j] = sl[j], sl[i] }
func (sl SpotList) Len() int      { return len(sl) }
func (sl SpotList) Less(i, j int) bool {
	if sl[i].Pos < sl[j].Pos {
		return true
	} else if sl[i].Pos > sl[j].Pos {
		return false
	}

	return sl[i].Type < sl[j].Type
}

//////////////////////////////

func (itv *Interval) String() string {
	sort.Sort(itv.Data)
	return fmt.Sprintf("Interval: range:%d-%d, id:%d, data:%v", itv.Low, itv.High, itv.NodeId, itv.Data)
}

func (st SpotType) String() string {
	return []string{"None", "Start", "End", "Mixed"}[st]
}

func (s *Spot) String() string {
	sort.Sort(s.Data)
	return fmt.Sprintf("Spot: pos:%d, type:%s, id:%d, data:%v", s.Pos, s.Type, s.NodeId, s.Data)
}

/*
func (s *Spot) PutData(d uint64) {
	for _, dd := range s.Data {
		if dd == d {
			return
		}
	}

	s.Data = append(s.Data, d)
}
*/

/*
func (sl *SpotList) GetData(limit int64) Uint64Slice {
	ids := make(Uint64Slice, 0)

	for _, sp := range *sl {
		if sp.Pos < limit && sp.Type == Start {
			ids = append(ids, sp.Data...)
		}
	}

	return ids
}
*/

func (sl *SpotList) GetSpotType() SpotType {
	var t SpotType = None

	for _, sp := range *sl {
		if t == None {
			t = sp.Type
		} else if t != sp.Type {
			return Mixed
		}
	}

	return t
}

func (sl *SpotList) GetAllData(t SpotType) Uint64Slice {
	ids := make(Uint64Slice, 0)

	for _, sp := range *sl {
		if t == None || sp.Type == t {
			ids = append(ids, sp.Data...)
		}
	}

	return ids
}

func (sl *SpotList) Push(sp *Spot) {
	*sl = append(*sl, sp)
}

/*
func (sl *SpotList) IsEmpty() bool {
	if len(*sl) == 0 {
		return true
	}

	return false
}
*/

/*
func (sl *SpotList) Remove(t SpotType) {
	ids := make([]uint64, 0)

	for _, sp := range rs.CurrentSpots {
		if sp.Type == End {
			ids = append(ids, sp.NodeId)
		}
	}

	for _, id := range ids {
		for i := 0; i < len(rs.CurrentSpots); i++ {
			if rs.CurrentSpots[i].NodeId == id {
				rs.CurrentSpots = append(rs.CurrentSpots[:i], rs.CurrentSpots[i+1:]...)
				i--
			}
		}
	}
}
*/

func NewInterval(low, high int64, nodeId uint64, data []uint64) *Interval {
	return &Interval{
		Low:    low,
		High:   high,
		NodeId: nodeId,
		Data:   data,
	}
}

func NewSpot(low, high int64, nodeId uint64, data []uint64) (start, end *Spot) {
	start = &Spot{
		Type:   Start,
		Pos:    low,
		NodeId: nodeId,
		Data:   make(Uint64Slice, 0),
	}
	start.Data = append(start.Data, data...)

	end = &Spot{
		Type:   End,
		Pos:    high,
		NodeId: nodeId,
		Data:   make(Uint64Slice, 0),
	}
	end.Data = append(end.Data, data...)

	return start, end
}

func (rs *RangeSplit) AddRange(low, high int64, nodeId uint64, data []uint64) {
	s, e := NewSpot(low, high, nodeId, data)

	rs.Spots = append(rs.Spots, s, e)
}

func (rs *RangeSplit) PushInterval(itv *Interval) {
	rs.Intervals = append(rs.Intervals, itv)
}

func (rs *RangeSplit) PushCurrentData(t SpotType) {
	data := rs.CurrentSpots.GetAllData(t)
	if len(data) < 1 {
		return
	}

	for _, d := range data {
		rs.CurrentData[d] = true
	}
}

func (rs *RangeSplit) GetCurrentData() Uint64Slice {
	data := make(Uint64Slice, 0)

	for d, _ := range rs.CurrentData {
		data = append(data, d)
	}

	return data
}

func (rs *RangeSplit) RemoveCurrentData() {
	for _, sp := range rs.CurrentSpots {
		if sp.Type != End {
			continue
		}

		for _, d := range sp.Data {
			delete(rs.CurrentData, d)
		}
	}
}

func (rs *RangeSplit) CleanCurrentSpots() {
	rs.CurrentSpots = make(SpotList, 0)
}

func (rs *RangeSplit) Init() {
	rs.CurrentPos = -1
	rs.StartPos = -1

	rs.CurrentSpots = make(SpotList, 0)
	rs.Intervals = make(IntervalList, 0)
	rs.CurrentData = make(CurrentData, 0)
}

func (rs *RangeSplit) AddInterval(end int64) {
	data := rs.GetCurrentData()

	if len(data) > 0 {
		itv := NewInterval(rs.StartPos, end, 0, data)
		rs.PushInterval(itv)
	}
}

func (rs *RangeSplit) Dump(msg string) {
	fmt.Printf("###>> Dump: %s\n", msg)
	fmt.Printf("Pos: Start=%d, Current=%d\n", rs.StartPos, rs.CurrentPos)

	for i, sp := range rs.CurrentSpots {
		fmt.Printf("%d: %s \n", i, sp)
	}

	for i, itv := range rs.Intervals {
		fmt.Printf("%d: %s \n", i, itv)
	}

	data := rs.GetCurrentData()
	sort.Sort(data)
	fmt.Printf("AllData: %v \n", data)

	fmt.Printf("<<###\n")
}

func (rs *RangeSplit) DumpAllSpots() {
	fmt.Printf("=== dump === \n")

	fmt.Printf("All Spots \n")
	for i, sp := range rs.Spots {
		fmt.Printf("%d: %s \n", i, sp)
	}
}

func (rs *RangeSplit) DumpIntervals(msg string) {
	fmt.Printf("###>> Dump: %s\n", msg)

	for i, itv := range rs.Intervals {
		fmt.Printf("%d: %s \n", i, itv)
	}

	fmt.Printf("<<###\n")
}

func (rs *RangeSplit) CreateInterval(pos int64) {
	t := rs.CurrentSpots.GetSpotType()

	if t == Mixed {
		//fmt.Printf("  ==> Begin Mixed: Pos=%d \n", pos)
		//rs.Dump("Begin Mxied")

		if len(rs.CurrentData) > 0 && rs.StartPos < rs.CurrentPos {
			// clsoing the previous interval
			rs.AddInterval(rs.CurrentPos - 1)
			rs.StartPos = rs.CurrentPos
		}

		// closing a dot
		rs.PushCurrentData(Start)
		rs.AddInterval(rs.CurrentPos)

		rs.RemoveCurrentData()
		rs.CleanCurrentSpots()

		rs.StartPos++
		rs.CurrentPos = pos
	} else if t == Start {
		//fmt.Printf("  ==> Begin Start: Pos=%d \n", pos)
		//rs.Dump("Begin Start")

		if len(rs.CurrentData) > 0 && rs.StartPos < rs.CurrentPos {
			// closing the previous interval
			rs.AddInterval(rs.CurrentPos - 1)
		}

		rs.PushCurrentData(None)

		rs.CleanCurrentSpots()

		rs.StartPos = rs.CurrentPos
		rs.CurrentPos = pos
	} else if t == End {
		//fmt.Printf("  ==> Begin End: Pos=%d \n", pos)
		//rs.Dump("Begin End")

		if len(rs.CurrentData) > 0 {
			end := rs.CurrentSpots[0].Pos
			rs.AddInterval(end)
		}

		rs.RemoveCurrentData()
		rs.CleanCurrentSpots()

		rs.StartPos = rs.CurrentPos + 1
		rs.CurrentPos = pos
	} else {
		fmt.Printf("############### Unknown Type: %v \n", t)
	}

}

func (rs *RangeSplit) Build() {
	sort.Sort(rs.Spots)
	rs.DumpAllSpots()

	for _, sp := range rs.Spots {
		if rs.StartPos == -1 {
			// the first time
			rs.CurrentPos = sp.Pos
			rs.StartPos = sp.Pos
			rs.CurrentSpots.Push(sp)
		} else if rs.CurrentPos == sp.Pos {
			rs.CurrentPos = sp.Pos
			rs.CurrentSpots.Push(sp)
		} else {
			rs.CreateInterval(sp.Pos)
			rs.CurrentSpots.Push(sp)
		}
	}

	if len(rs.CurrentSpots) > 0 {
		rs.CreateInterval(-1)
	}

	rs.DumpIntervals("Done")
}
