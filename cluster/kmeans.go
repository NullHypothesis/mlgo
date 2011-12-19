package cluster

import (
	"rand"
	"sort"
	"math"
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

// maximization step: move cluster meanss to centroids of data points
// Returns the cost
func (c *KMeans) maximization() {
	// move cluster meanss
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
	J /= float64( len(c.X) )

	c.Cost = J
}


type KMedians struct {
	KMeans
}

// Cluster runs the k-medians algorithm once with random initialization
// Returns the classification information
// N.B. Must explicitly override KMeans.Cluster s.t. KMedians.maximization is called
// instead of KMeans.maximization.
func (c *KMedians) Cluster(k int) (classes *Classes) {
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

// Override KMeans.maximization
// Calculate the median instead of mean;
// total absolute deviation instead of total sum of squares
func (c *KMedians) maximization() {
	// move cluster centroid_ii
	move := func(ii int, chCost chan float64) {
		centers := c.Centers[ii]
		errors := c.Errors[ii]
		// hold coordinate of each dimension for each member
		members := make(Matrix, len(centers))
		// initialize
		for j, _ := range centers {
			members[j] = make(Vector, len(c.X))
		}

		// gather all member data points
		n := 0
		for i, class := range c.Index {
			if class == ii {
				for j, _ := range centers {
					members[j][n] = c.X[i][j]
				}
				n++
			}
		}

		// compute centers and errors
		cost := 0.0
		for j, _ := range centers {
			// find median
			centers[j], errors[j] = median(members[j][:n])
			cost += errors[j]
		}
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
	J /= float64( len(c.X) )

	c.Cost = J
}

// find median and tad
// side-effect: x becomes sorted
func median(x Vector) (med, tad float64) {
	sort.Float64s(x)
	n := len(x)

	// calculate median
	if n % 2 == 0 {
		i := n/2
		med = (x[i] + x[i-1]) / 2
	} else {
		med = x[n/2]
	}

	// calculate total absolute deviation
	for _, z := range x {
		tad += math.Fabs( z - med )
	}

	return
}

