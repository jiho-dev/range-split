package main

import (
	"fmt"
	"sort"

	"jiho-dev.com/range-split/iputil"
	"jiho-dev.com/range-split/set"
)

/*
* Range can be represented by StartSpot and EndSpot
* Spots: a pair of spots representing a range
	-. type: start/end/dot
	-. pos: position
	-. data: user defined data
* Stack: remaining spots
* startPos: a start position of a new interval
* curPos: current position
* for all spots:
	curSpots = getNextSpots()
	curPos = curSpot[0].Pos
	if curSpots have StartSopt
		Close the previous interval if Stack is not empty and startPos < curPos
			addInterval(startPos, curPos -1)
		startPos = curPos
	PushStack(curSpots)
	if curSpots have EndSpot
		Close the previous interval if Stack is not empty
			addInterval(startPos, curPos)
		startPos = curPos + 1
	Remove the closed Spots in the Stack, which means if Start/EndSpot are in the Stack
*/

////////////////////////////////////
type Interval struct {
	Low     int64
	High    int64
	NodeIds set.Uint64Set
	Data    set.Uint64Set
}

type IntervalList []*Interval

///////////////////////////

type SpotType uint32

const (
	SpotStart SpotType = 0x00000001
	SpotEnd   SpotType = 0x00000002
)

type Spot struct {
	Type   SpotType
	Pos    int64
	NodeId uint64
	Data   set.Uint64Set // XXX: only for StartSpot
}

type SpotList []*Spot
type SpotMap map[uint64]bool // key: Spot.NodeId, val: bool

///////////////////////////

type RangeSplit struct {
	Intervals IntervalList
	Spots     SpotList // all spots being treated
	Stack     SpotList // the ranges alive
	StartPos  int64
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
	return []string{"None", "StartSpot", "EndSpot", "MixedSpot"}[st]
}

func (s *Spot) String() string {
	return fmt.Sprintf("<Spot: pos:%d, type:%s, id:%d, data:%v>", s.Pos, s.Type, s.NodeId, s.Data)
}

////////////////////////////////

func (sl *SpotList) GetNodeId() set.Uint64Set {
	ids := make(set.Uint64Set, 0)

	for _, sp := range *sl {
		ids = ids.Add(sp.NodeId)
	}

	return ids
}

func (sl *SpotList) GetNodeData() set.Uint64Set {
	data := make(set.Uint64Set, 0)

	for _, sp := range *sl {
		data = data.Add(sp.Data...)
	}

	return data
}

func (sl *SpotList) GetEndSpotId() SpotMap {
	end := SpotMap{}

	for _, sp := range *sl {
		if (sp.Type & SpotEnd) == SpotEnd {
			end[sp.NodeId] = true
		}
	}

	return end
}

func (sl *SpotList) Push(sp []*Spot) {
	*sl = append(*sl, sp...)
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
		Type:   SpotStart,
		Pos:    low,
		NodeId: nodeId,
	}

	end = &Spot{
		Type:   SpotEnd,
		Pos:    high,
		NodeId: nodeId,
	}

	return start, end
}

////////////////////////////////

func (rs *RangeSplit) pushInterval(itv *Interval) {
	rs.Intervals = append(rs.Intervals, itv)
}

// return all spots placed in the same position
// SpotList: list of spot in the same pos
// SpotType: type of SpotList
// nextIdx: next index of slice
func (rs *RangeSplit) getNextSpots(idx int) (SpotList, SpotType, int) {
	var pos int64 = -1
	var i int
	var spType SpotType
	var cnt = len(rs.Spots)

	if idx >= cnt {
		return nil, 0, -1
	}

	curSpots := SpotList{}

	for i = idx; i < cnt; i++ {
		sp := rs.Spots[i]

		if pos != -1 && pos != sp.Pos {
			break
		}

		pos = sp.Pos
		curSpots = append(curSpots, sp)
		spType |= sp.Type
	}

	return curSpots, spType, i
}

func (rs *RangeSplit) addInterval(end int64) {
	ids := rs.Stack.GetNodeId()
	data := rs.Stack.GetNodeData()
	if len(ids) < 1 {
		return
	}

	itv := NewInterval(rs.StartPos, end, ids)
	itv.Data = data

	sort.Sort(itv.Data)
	rs.pushInterval(itv)
}

func (rs *RangeSplit) removeClosedSpot() {
	remain := SpotList{}
	end := rs.Stack.GetEndSpotId()

	for _, sp := range rs.Stack {
		_, ok := end[sp.NodeId]
		if !ok {
			remain = append(remain, sp)
		}
	}

	rs.Stack = remain
}

//////////////////////////////////////
// external function

func (rs *RangeSplit) Init() {
	rs.StartPos = -1
	rs.Stack = make(SpotList, 0)
	rs.Intervals = make(IntervalList, 0)
}

func (rs *RangeSplit) AddRange(low, high int64, nodeId uint64, data []uint64) {
	s, e := NewSpot(low, high, nodeId)
	// XXX: StartSpot only has data
	s.Data = data

	rs.Spots = append(rs.Spots, s, e)
}

func (rs *RangeSplit) Build() {
	var nextIdx int = 0
	var curSpots SpotList
	var spType SpotType
	var curPos int64

	// the first key : position
	// the second key: Type. Start is the first
	sort.Sort(rs.Spots)

	for {
		// Get all Spots placed in the same position
		curSpots, spType, nextIdx = rs.getNextSpots(nextIdx)
		if nextIdx == -1 {
			// end of data
			break
		} else if curSpots == nil || len(curSpots) < 1 {
			continue
		}

		curPos = curSpots[0].Pos

		// exist StartSpot
		if (spType & SpotStart) == SpotStart {
			// exist remaining range ?
			if len(rs.Stack) > 0 && rs.StartPos < curPos {
				rs.addInterval(curPos - 1)
			}

			rs.StartPos = curPos
		}

		rs.Stack.Push(curSpots)

		//exist EndSpot
		if (spType & SpotEnd) == SpotEnd {
			if len(rs.Stack) > 0 {
				rs.addInterval(curPos)
				rs.removeClosedSpot()
			}

			rs.StartPos = curPos + 1
		}
	}
}

func (rs *RangeSplit) Dump(msg string) {
	fmt.Printf(">>> RangeSplit Dump: %s\n", msg)
	fmt.Printf("Pos: StartSpot=%d \n", rs.StartPos)

	for i, sp := range rs.Stack {
		fmt.Printf("%d: %s \n", i, sp)
	}

	for i, itv := range rs.Intervals {
		fmt.Printf("%d: %s \n", i, itv)
	}

	fmt.Printf("<<<\n")
}

func (rs *RangeSplit) DumpAllSpots() {
	fmt.Printf(">>> RangeSplit Spot Dump \n")

	fmt.Printf("All Spots \n")
	for i, sp := range rs.Spots {
		fmt.Printf("%d: %s \n", i, sp)
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
