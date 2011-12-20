package cluster

import (
	"sort"
)

type KMedoids struct {
	KMeans
	// Distances between data points [m x m]
	Distances Matrix
}

func NewKMedoids(X Matrix, metric MetricOp) *KMedoids {
	return &KMedoids{
		KMeans: KMeans{X:X, Metric:metric},
		Distances: Distances(X, metric),
	}
}

// Cluster runs the k-medoids algorithm once with random initialization
// Returns the classification information
func (c *KMedoids) Cluster(k int) (classes *Classes) {
	if c.X == nil { return }
	c.K = k
	c.initialize()
	i := 0
	for !c.expectation() && (c.MaxIter == 0 || i < c.MaxIter) {
		c.maximization()
		i++
	}

	// copy classifcation information
	classes = &Classes{
		make([]int, len(c.X)), c.Cost }
	copy(classes.Index, c.Index)

	return
}

type pair struct {
	key int
	value float64
}

type pairs []pair

func (p pairs) Len() int {
	return len(p)
}

func (p pairs) Less(i, j int) bool {
	return p[i].value < p[j].value
}

func (p pairs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Initialize the medoids by choosing the most central k data points
func (c *KMedoids) initialize() {
	// calculate normalized distances
	normalized := make(Matrix, len(c.Distances))
	for i, d := range c.Distances {
		normalized[i] = make(Vector, len(d))
		sum := 0.0
		for j, x := range d {
			normalized[i][j] = x
			sum += x
		}
		for j, _ := range d {
			normalized[i][j] /= sum
		}
	}

	// sum the normalized distances across all rows
	p := make(pairs, len(normalized[0]))
	for i, _ := range(normalized) {
		p[i].key = i
		for _, x := range normalized[i] {
			p[i].value += x
		}
	}

	// sort the summed normalized distances
	sort.Sort(p)

	// initialize centers
	c.Centers, c.Errors = make(Matrix, c.K), make(Vector, c.K)
	c.Index = make([]int, len(c.X))
	for k, _ := range c.Centers {
		// use the first k data points sorted by summed normalized distances
		x := c.X[ p[k].key ]
		c.Centers[k] = make(Vector, len(x))
		copy(c.Centers[k], x)
	}
}

// Maximization step: Swap medoid with another data point in the cluster
// s.t. total distance to the new medoid is minimized.
func (c *KMedoids) maximization() {
	// swap medoid
	swap := func(ii int, chCost chan float64) {
		center := c.Centers[ii]

		// gather members
		n := 0
		memberIdx := make([]int, len(c.Index))
		for i, class := range c.Index {
			if class == ii {
				memberIdx[n] = i
				n++
			}
		}
		memberIdx = memberIdx[:n]

		// calculate total distances for each member
		totalDistances := make(Vector, n)
		n = 0
		for _, i := range memberIdx {
			for _, j := range memberIdx {
				totalDistances[n] += c.Distances[i][j]
			}
			n++
		}

		// find the member with the minimum total distance
		// set this as the new center
		newCenter, min := memberIdx[0], totalDistances[0]
		for i, d := range totalDistances {
			if d < min {
				newCenter, min = memberIdx[i], d
			}
		}
		copy(center, c.X[newCenter])

		// use the minimum total distance as the cost
		c.Errors[ii] = min
		chCost <- min
	}

	// process cluster center concurrently
	ch := make(chan float64)
	for ii, _ := range c.Centers {
		go swap(ii, ch)
	}

	// collect results
	J := 0.0
	for ii := 0; ii < len(c.Centers); ii++ {
		J += <-ch;
	}
	c.Cost = J / float64( len(c.X) )
}

