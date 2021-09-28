package set

// https://github.com/goware/set/blob/master/set.go

import (
	"sort"

	"github.com/xtgo/set"
)

type BoolOp func(sort.Interface, int) bool

//////////////////////////////////////////////
// Int64 set
type Int64Set []int64

func NewInt64Set(v ...int64) Int64Set {
	s := Int64Set{}
	if len(v) > 0 {
		s = s.Add(v...)
	}
	return s
}

func (s Int64Set) Len() int           { return len(s) }
func (s Int64Set) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Int64Set) Less(i, j int) bool { return s[i] < s[j] }

func (s Int64Set) Do(op set.Op, t Int64Set) Int64Set {
	sort.Sort(t)
	data := append(s, t...)
	n := op(data, len(s))
	return data[:n]
}

func (s Int64Set) DoBool(op BoolOp, t Int64Set) bool {
	data := append(s, t...)
	return op(data, len(s))
}

func (s Int64Set) Union(t Int64Set) Int64Set  { return s.Do(set.Union, t) }
func (s Int64Set) Inter(t Int64Set) Int64Set  { return s.Do(set.Inter, t) }
func (s Int64Set) Diff(t Int64Set) Int64Set   { return s.Do(set.Diff, t) }
func (s Int64Set) Add(v ...int64) Int64Set    { return s.Union(Int64Set(v)) }
func (s Int64Set) Remove(v ...int64) Int64Set { return s.Diff(NewInt64Set(v...)) }
func (s Int64Set) Exists(v int64) bool        { return s.DoBool(set.IsInter, Int64Set{v}) }

//////////////////////////////////////////////
// uint64 set
type Uint64Set []uint64

func NewUint64Set(v ...uint64) Uint64Set {
	obj := Uint64Set{}
	if len(v) > 0 {
		obj = obj.Add(v...)
	}
	return obj
}

func (obj Uint64Set) Len() int           { return len(obj) }
func (obj Uint64Set) Swap(i, j int)      { obj[i], obj[j] = obj[j], obj[i] }
func (obj Uint64Set) Less(i, j int) bool { return obj[i] < obj[j] }

func (obj Uint64Set) Do(op set.Op, t Uint64Set) Uint64Set {
	sort.Sort(t)
	data := append(obj, t...)
	n := op(data, len(obj))
	return data[:n]
}

func (obj Uint64Set) DoBool(op BoolOp, t Uint64Set) bool {
	data := append(obj, t...)
	return op(data, len(obj))
}

func (obj Uint64Set) Union(t Uint64Set) Uint64Set  { return obj.Do(set.Union, t) }
func (obj Uint64Set) Inter(t Uint64Set) Uint64Set  { return obj.Do(set.Inter, t) }
func (obj Uint64Set) Diff(t Uint64Set) Uint64Set   { return obj.Do(set.Diff, t) }
func (obj Uint64Set) Add(v ...uint64) Uint64Set    { return obj.Union(Uint64Set(v)) }
func (obj Uint64Set) Remove(v ...uint64) Uint64Set { return obj.Diff(NewUint64Set(v...)) }
func (obj Uint64Set) Exists(v uint64) bool         { return obj.DoBool(set.IsInter, Uint64Set{v}) }

//////////////////////////////////////////////
// Uint32 set
type Uint32Set []uint32

func NewUint32Set(v ...uint32) Uint32Set {
	s := Uint32Set{}
	if len(v) > 0 {
		s = s.Add(v...)
	}
	return s
}

func (s Uint32Set) Len() int           { return len(s) }
func (s Uint32Set) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Uint32Set) Less(i, j int) bool { return s[i] < s[j] }

func (s Uint32Set) Do(op set.Op, t Uint32Set) Uint32Set {
	sort.Sort(t)
	data := append(s, t...)
	n := op(data, len(s))
	return data[:n]
}

func (s Uint32Set) DoBool(op BoolOp, t Uint32Set) bool {
	data := append(s, t...)
	return op(data, len(s))
}

func (s Uint32Set) Union(t Uint32Set) Uint32Set  { return s.Do(set.Union, t) }
func (s Uint32Set) Inter(t Uint32Set) Uint32Set  { return s.Do(set.Inter, t) }
func (s Uint32Set) Diff(t Uint32Set) Uint32Set   { return s.Do(set.Diff, t) }
func (s Uint32Set) Add(v ...uint32) Uint32Set    { return s.Union(Uint32Set(v)) }
func (s Uint32Set) Remove(v ...uint32) Uint32Set { return s.Diff(NewUint32Set(v...)) }
func (s Uint32Set) Exists(v uint32) bool         { return s.DoBool(set.IsInter, Uint32Set{v}) }

////////////////////////////////////////////////////
// String set
type StringSet []string

func NewStringSet(v ...string) StringSet {
	s := StringSet{}
	if len(v) > 0 {
		s = s.Add(sort.StringSlice(v)...)
	}
	return s
}

func (s StringSet) Do(op set.Op, t StringSet) StringSet {
	data := sort.StringSlice(append(s, t...))
	n := op(data, len(s))
	return StringSet(data[:n])
}

func (s StringSet) DoBool(op BoolOp, t StringSet) bool {
	data := sort.StringSlice(append(s, t...))
	return op(data, len(s))
}

func (s StringSet) Union(t StringSet) StringSet  { return s.Do(set.Union, t) }
func (s StringSet) Inter(t StringSet) StringSet  { return s.Do(set.Inter, t) }
func (s StringSet) Diff(t StringSet) StringSet   { return s.Do(set.Diff, t) }
func (s StringSet) Add(v ...string) StringSet    { return s.Union(StringSet(v)) }
func (s StringSet) Remove(v ...string) StringSet { return s.Diff(NewStringSet(v...)) }
func (s StringSet) Exists(v string) bool         { return s.DoBool(set.IsInter, StringSet{v}) }
