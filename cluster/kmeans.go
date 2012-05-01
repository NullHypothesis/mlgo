package cluster

import (
	"math/rand"
)

type KMeans struct {
	// Matrix of data points
	X Matrix
	// Distance metric
	Metric MetricOp
	// number of clusters
	K int
	// Matrix of centroids	
	Centers Matrix
	// Total distance of members to each centroid
	Errors Vector
	// cluster center assignment index
	Index []int
	// cost
	Cost float64
	// Maximum number of iterations
	MaxIter int
}

func NewKMeans(X Matrix, metric MetricOp) *KMeans {
	return &KMeans{
		X:      X,
		Metric: metric,
	}
}

// Cluster runs the k-means algorithm once with random initialization
// Returns the classification information
func (c *KMeans) Cluster(k int) (classes *Classes) {
	if c.X == nil {
		return
	}
	c.K = k
	c.initialize()
	i := 0
	for !c.expectation() && (c.MaxIter == 0 || i < c.MaxIter) {
		c.maximization()
		i++
	}

	// copy classifcation information
	classes = &Classes{
		make([]int, len(c.X)), c.Cost}
	copy(classes.Index, c.Index)

	return
}

// initialize the cluster centroids by randomly selecting data points
func (c *KMeans) initialize() {
	c.Centers, c.Errors = make(Matrix, c.K), make(Vector, c.K)
	c.Index = make([]int, len(c.X))
	for k, _ := range c.Centers {
		x := c.X[rand.Intn(len(c.X))]
		c.Centers[k] = make(Vector, len(x))
		copy(c.Centers[k], x)
	}
}

// expectation step: assign data points to cluster centroids
// Returns whether the algorithm has converged
func (c *KMeans) expectation() (converged bool) {
	// find the centroids that is closest to the current data point
	assign := func(i int, chIndex chan int) {
		index, min := 0, maxValue
		// find the center with the minimum distance
		for ii := 0; ii < len(c.Centers); ii++ {
			distance := c.Metric(c.X[i], c.Centers[ii])
			if distance < min {
				index, min = ii, distance
			}
		}
		chIndex <- index
	}

	// process examples concurrently
	ch := make(chan int)
	for i, _ := range c.X {
		go assign(i, ch)
	}

	// collect results
	converged = true
	for i, _ := range c.X {
		if index := <-ch; c.Index[i] != index {
			c.Index[i] = index
			converged = false
		}
	}

	return
}

// maximization step: move cluster centers to centroids of data points
func (c *KMeans) maximization() {
	// move the center of cluster_ii to the mean
	move := func(ii int, chCost chan float64) {
		center := c.Centers[ii]

		// zero the coordinates
		for j, _ := range center {
			center[j] = 0
		}

		// compute centroid and gather members
		n := 0
		memberIdx := make([]int, len(c.Index))
		for i, class := range c.Index {
			if class == ii {
				for j, _ := range center {
					x := c.X[i][j]
					center[j] += x
				}
				memberIdx[n] = i
				n++
			}
		}
		memberIdx = memberIdx[:n]

		fn := float64(n)
		for j, _ := range center {
			center[j] /= fn
		}

		// compute cost
		cost := 0.0
		for _, i := range memberIdx {
			cost += c.Metric(center, c.X[i])
		}

		c.Errors[ii] = cost
		chCost <- cost
	}

	// process cluster centers concurrently
	ch := make(chan float64)
	for ii, _ := range c.Centers {
		go move(ii, ch)
	}

	// collect results
	J := 0.0
	for ii := 0; ii < len(c.Centers); ii++ {
		J += <-ch
	}
	c.Cost = J / float64(len(c.X))

}
