package cluster

import (
	"testing"
	"sort"
)

type KeyValues struct {
	r []KeyValue
}

func (x KeyValues) Len() int {
	return len(x.r)
}

func (x KeyValues) Less(i, j int) bool {
	return x.r[i].Key < x.r[j].Key
}

func (x KeyValues) Swap(i, j int) {
	x.r[i], x.r[j] = x.r[j], x.r[i]
}

func (x *KeyValues) Copy(y KeyValues) {
	x.r = make([]KeyValue, len(y.r))
	copy(x.r, y.r)
}

func TestHeap(t *testing.T) {
	x := KeyValues{
		[]KeyValue {
			KeyValue{1, 4},
			KeyValue{2, 2},
			KeyValue{3, 5},
			KeyValue{4, 8},
			KeyValue{5, 9},
		},
	}
	y := KeyValues{}
	y.Copy(x)

	// create heap by using array
	h := Heap{y.r}
	h.Init()

	// create heap by sequential pushes
	g := Heap{}
	for _, a := range x.r {
		g.Push(a)
	}

	// sort original array
	sort.Sort(x)


	for i := 0; i < x.Len(); i++ {
		a, b, c := x.r[i].Value, h.Pop(), g.Pop()
		if a != b {
			t.Errorf("%d-th min in heapified array = %d, expected %d", i, b, a)
		}
		if a != c {
			t.Errorf("%d-th min in incrementally built heap = %d, expected %d", i, c, a)
		}
	}

	y.Copy(x)
	h = Heap{y.r}
	h.Init()

	// replace min value with large value
	h.Update(0, KeyValue{1,10})

	if d := h.Pop(); d != x.r[1].Value {
		t.Errorf("new min updated heap = %d, expected %d", d, x.r[1].Value)
	}


	// replace the last leaf with the min value
	a := KeyValue{10, 1.0}
	h.Update(h.Len()-1, a)

	if d := h.Pop(); d != a.Value {
		t.Errorf("new min updated (2) heap = %d, expected %d", d, a.Value)
	}
}

