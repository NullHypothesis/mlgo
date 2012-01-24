package cluster

import (
	"testing"
	"sort"
)

func TestHeap(t *testing.T) {
	x := []int{4, 2, 5, 8, 9}
	y := make([]int, len(x))
	copy(y, x)

	// create heap by using array
	h := Heap{ y }
	h.Init()

	// create heap by sequential pushes
	g := Heap{}
	for _, a := range x {
		g.Push(a)
	}

	// sort original array
	sort.Ints(x)

	for i := 0; i < len(x); i++ {
		a, b, c := x[i], h.Pop(), g.Pop()
		if a != b {
			t.Errorf("%d-th min in heapified array = %d, expected %d", i, b, a)
		}
		if a != c {
			t.Errorf("%d-th min in incrementally built heap = %d, expected %d", i, c, a)
		}
	}
}
