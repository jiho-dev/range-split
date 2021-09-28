package main

import (
	"fmt"
	"sort"

	"jiho-dev.com/range-split/iputil"
	"jiho-dev.com/range-split/set"
)

////////////////////////////////////

const (
	AllSpot SpotType = iota
	StartSpot
	EndSpot
	MixedSpot
)

////////////////////////////////////

type CurrentNodeIdMap map[uint64]bool

///////////////////////////
type Interval struct {
	Low     int64
	High    int64
	NodeIds set.Uint64Set
	Data    set.Uint64Set
}

type IntervalList []*Interval

///////////////////////////
type SpotType int
type Spot struct {
	Type   SpotType
	Pos    int64
	NodeId uint64
}

type SpotList []*Spot
type SpotDataMap map[uint64]set.Uint64Set

///////////////////////////

type RangeSplit struct {
	Intervals      IntervalList
	Spots          SpotList
	SpotData       SpotDataMap
	StartPos       int64
	CurrentPos     int64
	CurrentSpots   SpotList
	CurrentNodeIds CurrentNodeIdMap
}

////////////////////////////////
// for Sorting items
func (il IntervalList) Swap(i, j int)      { il[i], il[j] = il[j], il[i] }
func (il IntervalList) Len() int           { return len(il) }
func (il IntervalList) Less(i, j int) bool { return il[i].Low < il[j].Low }

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
	sort.Sort(itv.NodeIds)
	return fmt.Sprintf("Interval: range:%d-%d, IDs:%v, data:%v", itv.Low, itv.High, itv.NodeIds, itv.Data)
}

func (st SpotType) String() string {
	return []string{"AllSpot", "StartSpot", "EndSpot", "MixedSpot"}[st]
}

func (s *Spot) String() string {
	//sort.Sort(s.Data)
	return fmt.Sprintf("Spot: pos:%d, type:%s, id:%d", s.Pos, s.Type, s.NodeId)
}

////////////////////////////////

func (sl *SpotList) GetSpotType() SpotType {
	var t SpotType = AllSpot

	for _, sp := range *sl {
		if t == AllSpot {
			t = sp.Type
		} else if t != sp.Type {
			return MixedSpot
		}
	}

	return t
}

func (sl *SpotList) GetAllNodeId(t SpotType) set.Uint64Set {
	ids := make(set.Uint64Set, 0)

	for _, sp := range *sl {
		if t == AllSpot || sp.Type == t {
			ids = ids.Add(sp.NodeId)
		}
	}

	return ids
}

func (sl *SpotList) Push(sp *Spot) {
	*sl = append(*sl, sp)
}

////////////////////////////////

func NewInterval(low, high int64, ids set.Uint64Set) *Interval {
	return &Interval{
		Low:     low,
		High:    high,
		NodeIds: ids,
	}
}

func NewSpot(low, high int64, nodeId uint64) (start, end *Spot) {
	start = &Spot{
		Type:   StartSpot,
		Pos:    low,
		NodeId: nodeId,
	}

	end = &Spot{
		Type:   EndSpot,
		Pos:    high,
		NodeId: nodeId,
	}

	return start, end
}

////////////////////////////////

func (rs *RangeSplit) pushInterval(itv *Interval) {
	rs.Intervals = append(rs.Intervals, itv)
}

func (rs *RangeSplit) pushCurrentNodeIds(t SpotType) {
	ids := rs.CurrentSpots.GetAllNodeId(t)
	if len(ids) < 1 {
		return
	}

	for _, id := range ids {
		rs.CurrentNodeIds[id] = true
	}
}

func (rs *RangeSplit) getCurrentNodeIds() set.Uint64Set {
	ids := make(set.Uint64Set, 0)

	for id, _ := range rs.CurrentNodeIds {
		ids = ids.Add(id)
	}

	return ids
}

func (rs *RangeSplit) removeCurrentNodeIds() {
	for _, sp := range rs.CurrentSpots {
		if sp.Type != EndSpot {
			continue
		}

		delete(rs.CurrentNodeIds, sp.NodeId)
	}
}

func (rs *RangeSplit) cleanCurrentSpots() {
	rs.CurrentSpots = make(SpotList, 0)
}

func (rs *RangeSplit) addInterval(end int64) {
	ids := rs.getCurrentNodeIds()

	if len(ids) < 1 {
		return
	}

	itv := NewInterval(rs.StartPos, end, ids)

	for _, id := range ids {
		if data, ok := rs.SpotData[id]; ok {
			itv.Data = itv.Data.Add(data...)
		}
	}

	sort.Sort(itv.Data)
	rs.pushInterval(itv)
}

func (rs *RangeSplit) closeInterval(pos int64) {
	t := rs.CurrentSpots.GetSpotType()

	if t == MixedSpot {
		if len(rs.CurrentNodeIds) > 0 && rs.StartPos < rs.CurrentPos {
			// clsoing the previous interval
			rs.addInterval(rs.CurrentPos - 1)
			rs.StartPos = rs.CurrentPos
		}

		// closing a dot
		rs.pushCurrentNodeIds(StartSpot)
		rs.addInterval(rs.CurrentPos)

		rs.removeCurrentNodeIds()
		rs.cleanCurrentSpots()

		if len(rs.CurrentNodeIds) < 1 {
			rs.StartPos = pos
		} else {
			rs.StartPos++
		}
		rs.CurrentPos = pos
	} else if t == StartSpot {
		if len(rs.CurrentNodeIds) > 0 && rs.StartPos < rs.CurrentPos {
			// closing the previous interval
			rs.addInterval(rs.CurrentPos - 1)
		}

		rs.pushCurrentNodeIds(AllSpot)

		rs.cleanCurrentSpots()

		rs.StartPos = rs.CurrentPos
		rs.CurrentPos = pos
	} else if t == EndSpot {
		if len(rs.CurrentNodeIds) > 0 {
			end := rs.CurrentSpots[0].Pos
			rs.addInterval(end)
		}

		rs.removeCurrentNodeIds()
		rs.cleanCurrentSpots()

		rs.StartPos = rs.CurrentPos + 1
		rs.CurrentPos = pos
	}

}

//////////////////////////////////////
// external function

func (rs *RangeSplit) Init() {
	rs.CurrentPos = -1
	rs.StartPos = -1

	rs.CurrentSpots = make(SpotList, 0)
	rs.Intervals = make(IntervalList, 0)
	rs.CurrentNodeIds = make(CurrentNodeIdMap, 0)
	rs.SpotData = make(SpotDataMap, 0)
}

func (rs *RangeSplit) AddRange(low, high int64, nodeId uint64, data []uint64) {
	s, e := NewSpot(low, high, nodeId)

	rs.Spots = append(rs.Spots, s, e)
	rs.SpotData[nodeId] = data
}

func (rs *RangeSplit) Build() {
	sort.Sort(rs.Spots)

	for _, sp := range rs.Spots {
		if rs.StartPos == -1 {
			// the first time
			rs.CurrentPos = sp.Pos
			rs.StartPos = sp.Pos
			rs.CurrentSpots.Push(sp)
		} else if rs.CurrentPos == sp.Pos {
			//rs.CurrentPos = sp.Pos
			rs.CurrentSpots.Push(sp)
		} else {
			rs.closeInterval(sp.Pos)
			rs.CurrentSpots.Push(sp)
		}
	}

	if len(rs.CurrentSpots) > 0 {
		rs.closeInterval(-1)
	}
}

func (rs *RangeSplit) Dump(msg string) {
	fmt.Printf(">>> RangeSplit Dump: %s\n", msg)
	fmt.Printf("Pos: StartSpot=%d, Current=%d\n", rs.StartPos, rs.CurrentPos)

	for i, sp := range rs.CurrentSpots {
		fmt.Printf("%d: %s \n", i, sp)
	}

	for i, itv := range rs.Intervals {
		fmt.Printf("%d: %s \n", i, itv)
	}

	ids := rs.getCurrentNodeIds()
	sort.Sort(ids)
	fmt.Printf("AllNodeIds: %v \n", ids)

	fmt.Printf("<<<\n")
}

func (rs *RangeSplit) DumpAllSpots() {
	fmt.Printf(">>> RangeSplit Spot Dump \n")

	fmt.Printf("All Spots \n")
	for i, sp := range rs.Spots {
		d := rs.SpotData[sp.NodeId]
		fmt.Printf("%d: %s, data:%v \n", i, sp, d)
	}

	fmt.Printf("<<<\n")
}

func (rs *RangeSplit) DumpIntervals(msg string) {
	fmt.Printf(">>> RangeSplit Interval Dump: %s\n", msg)

	for i, itv := range rs.Intervals {
		fmt.Printf("%d: %s \n", i, itv)
	}

	fmt.Printf("<<<\n")
}

func (rs *RangeSplit) DumpIpIntervals(msg string) {
	fmt.Printf(">>> RangeSplit IP Interval Dump: %s\n", msg)

	for i, itv := range rs.Intervals {
		fmt.Printf("%d: %s - %s, data:%v \n", i,
			iputil.Int2ip(uint32(itv.Low)), iputil.Int2ip(uint32(itv.High)), itv.Data)
	}

	fmt.Printf("<<<\n")
}
