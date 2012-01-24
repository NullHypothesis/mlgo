package cluster

// min binary heap
// complete binary tree
// partial order: every node a stores a value that is less than
// or equal to that of its children
type Heap struct {
	values []int
}

func (h *Heap) Init() {
	// heapify existing array
	// assumes data is already stored in array
	// worse case complexity is O(n)
	n := len(h.values)
	for i := n/2 - 1; i >= 0; i-- {
		siftdown(h.values, i)
	}
}

func (h *Heap) Push(value int) {
	n := len(h.values)
	// first place the value at the end of the heap, then sift up
	h.values = append(h.values, value)
	// sift appended value to correct position
	siftup(h.values, n)
}

func siftup(values []int, i int) {
	// sift up until x's parent <= x
	// i, j are the indices of x and x's parent
	for ; i != 0; {
		// parent
		j := (i-1)/2
		if values[j] <= values[i] {
			break
		}
		// sift up
		values[i], values[j] = values[j], values[i]
		i = j
	}
}

func siftdown(values []int, i int) {
	// n is the new size of the values array
	// values array has not been re-resized yet
	n := len(values)

	// loop until i is a leaf
	for ; !(i >= n/2 && i < n); {
		j := 2*i + 1     // left child
		right := j + 1   // right child: 2*i + 2
		if right < n && values[right] < values[j] {
			// set j to lesser child
			j = right
		}
		if values[i] <= values[j] {
			// i is less than both its children: heap order satisified
			break
		}
		// sift down
		values[i], values[j] = values[j], values[i]
		i = j
	}
}

// Pop removes the minimum element
// complexity is O(log(n))
func (h *Heap) Pop() (y int) {
	n := len(h.values)
	if n == 0 { return -1 }
	n--
	// swap min with the last value
	h.values[0], h.values[n] = h.values[n], h.values[0]

	// get the value to remove (min)
	y = h.values[n]
	// re-slice to remove the last element
	h.values = h.values[:n]

	// sift down new root
	if n != 0 { siftdown(h.values, 0) }

	return
}

func (h *Heap) Remove(i int) (y int) {
	n := len(h.values)
	if i < 0 || i >= n {
		// invalid position
		return -1
	}
	n--
	// swap with last value
	h.values[i], h.values[n] = h.values[n],h.values[i]

	// pop value
	y = h.values[n]
	h.values = h.values[:n]
	
	if n != 0 {
		// sift up if small key
		siftup(h.values, i)
		// sift down if large key
		siftdown(h.values, i)
	}

	return
}


// Update updates the value at the i-th position
func (h *Heap) Update(i int, x int) {
	n := len(h.values)
	if i < 0 || i >= n {
		// invalid position
		return
	}
	
	old := h.values[i]
	h.values[i] = x

	if x < old {
		// value became smaller: siftup
		siftup(h.values, i)
	} else {
		// value became bigger: siftdown
		// no real effect if value did not change
		siftdown(h.values, i)
	}
}

func (h *Heap) Len() int {
	return len(h.values)
}

