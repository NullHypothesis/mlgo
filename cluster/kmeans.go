package cluster

import (
	"rand"
)

type KMeans struct {
	// Matrix of data points
	X Matrix
	// number of clusters
	K int
	// Matrix of centroids	
	Centers, Errors Matrix
	// cluster center assignment index
	Index []int
	// cost
	Cost float64
	// Maximum number of iterations
	MaxIter int
}

func NewKMeans(X Matrix) *KMeans {
	return &KMeans{X:X}
}

// Cluster runs the k-means algorithm once with random initialization
// Returns the classification information
func (c *KMeans) Cluster(k int) (classes *Classes) {
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

// initialize the cluster centroids by randomly selecting data points
func (c *KMeans) initialize() {
	c.Centers, c.Errors = make(Matrix, c.K), make(Matrix, c.K)
	c.Index = make([]int, len(c.X))
	for k, _ := range c.Centers {
		x := c.X[ rand.Intn(len(c.X)) ]
		c.Centers[k], c.Errors[k] = make(Vector, len(x)), make(Vector, len(x))
		copy(c.Centers[k], x)
	}
}

// expectation step: assign data points to cluster centroids
// Returns whether the algorithm has converged
func (c *KMeans) expectation() (converged bool) {
	// find the centroids that is closest to the current data point
	assign := func(i int, chIndex chan int) {
		index, min :=  0, maxValue
		for ii := 0; ii < len(c.Centers); ii++ {
			// calculate distance
			distance := 0.0
			for j := 0; j < len(c.X[i]); j++ {
				diff := c.X[i][j] - c.Centers[ii][j]
				distance += diff * diff
			}
			if distance < min {
				index = ii
				min = distance
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
		means := c.Centers[ii]
		errors := c.Errors[ii]
		// zero the coordinates
		for j, _ := range means {
			means[j] = 0
			errors[j] = 0
		}
		// compute centroid
		n := 0.0
		for i, class := range c.Index {
			if class == ii {
				for j, _ := range means {
					x := c.X[i][j]
					means[j] += x
					errors[j] += x * x
				}
				n++
			}
		}
		cost := 0.0
		for j, _ := range means {
			mean := means[j] / n
			means[j] = mean
			// complete calculating the variance*N using the sum of squares formula
			errors[j] -= mean*mean * n
			cost += errors[j]
		}
		// mathematically, the cost is equivalent to 1/m * sum_i( ||x_i - center_i||^2 )
		// where m is the number of data points, and center_i is the
		// center to which x_i is assigned
		chCost <- cost;
	}

	// process cluster centers concurrently
	ch := make(chan float64)
	for ii, _ := range c.Centers {
		go move(ii, ch)
	}

	// collect results
	J := 0.0
	for ii := 0; ii < len(c.Centers); ii++ {
		J += <-ch;
	}
	c.Cost = J / float64( len(c.X) )

}

