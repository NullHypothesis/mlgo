package kmeans

import (
	"rand"
	"math"
	"fmt"
)

type Vector []float64
type Matrix [][]float64

const maxValue = math.MaxFloat64

type Clusters struct {
	// Matrix of data points
	X Matrix
	// number of clusters
	K int
	// Matrix of centroids	
	Means, Errors Matrix
	// cluster means assignment index
	Classes []int
	// cost
	Cost float64
}

// Run runs the k-means algorightm for iter random rounds
func Run(X Matrix, K int, iter int) (clusters *Clusters) {
	// run k-means concurrently
	ch := make(chan *Clusters)
	for i := 0; i < iter; i++ {
		c := Clusters{X:X, K:K}
		go c.Solve(ch)
	}

	// determine best clustering
	minCost := maxValue
	for i := 0; i < iter; i++ {
		if c := <-ch; c.Cost < minCost {
			clusters = c
			minCost = c.Cost
		}
	}
	return
}

// Solve runs the k-means algorithm once with random initialization
// Returns the cost
func (c *Clusters) Solve(chClusters chan<- *Clusters) {
	var cost float64
	c.initialize()
	for !c.expectation() {
		cost = c.maximization()
	}
	c.Cost = cost
	fmt.Println("c.Classes", c.Classes)
	fmt.Println("c.Means", c.Means)
	chClusters <- c
}

// initialize the cluster meanss randomly
func (c *Clusters) initialize() {
	c.Means, c.Errors = make(Matrix, c.K), make(Matrix, c.K)
	c.Classes = make([]int, len(c.X))
	for k, _ := range c.Means {
		x := c.X[ rand.Intn(len(c.X)) ]
		c.Means[k], c.Errors[k] = make(Vector, len(x)), make(Vector, len(x))
		copy(c.Means[k], x)
	}
}

// expectation step: assign data points to cluster meanss
// Returns whether the algorithm has converged
func (c *Clusters) expectation() (converged bool) {
	// find the means that is closest to the current data point
	assign := func(i int, chIndex chan int) {
		index, min :=  0, maxValue
		for ii := 0; ii < len(c.Means); ii++ {
			// calculate distance
			distance := 0.0
			for j := 0; j < len(c.X[i]); j++ {
				diff := c.X[i][j] - c.Means[ii][j]
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
		if index := <- ch; c.Classes[i] != index {
			c.Classes[i] = index
			converged = false
		}
	}

	return
}

// maximization step: move cluster meanss to centroids of data points
// Returns the cost
func (c *Clusters) maximization() float64 {
	// move cluster meanss
	move := func(ii int, chCost chan float64) {
		means := c.Means[ii]
		errors := c.Errors[ii]
		// zero the coordinates
		for j, _ := range means {
			means[j] = 0
			errors[j] = 0
		}
		// compute centroid
		n := 0.0
		for i, class := range c.Classes {
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
			// calculate variance * N using sum of squares formula
			errors[j] -= mean*mean * n
			cost = errors[j]
		}
		chCost <- cost;
	}

	// process cluster centers concurrently
	ch := make(chan float64)
	for ii, _ := range c.Means {
		go move(ii, ch)
	}

	// collect results
	J := 0.0
	for ii := 0; ii < len(c.Means); ii++ {
		J += <-ch;
	}
	J /= float64( len(c.X) )

	return J
}

