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
	h := Heap{y}
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

	y = make([]int, len(x))
	copy(y, x)
	h = Heap{y}
	h.Init()


	// replace min value with large value
	h.Update(0, 10)


	if d := h.Pop(); d != x[1] {
		t.Errorf("new min updated heap = %d, expected %d", d, x[1])
	}

	// replace the last leaf with the min value
	a := 1
	h.Update(h.Len()-1, a)

	if d := h.Pop(); d != a {
		t.Errorf("new min updated (2) heap = %d, expected %d", d, a)
	}
}
