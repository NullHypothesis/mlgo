package cluster

import (
	"math/rand"
	"mlgo/base"
	"fmt"
)

// TODO Make repeat runs internal to KMeans, KMedians, and KMedoids
// FIXME KMeans may initialize to 'duplicate' data points when there are ties,
//       which will result in fewer than k clusters
// FIXME KMeans may also return fewer than k clusters when a centroid loses all members during expectation

type KMeans struct {
	// Matrix of data points
	X Matrix
	// Distance metric
	Metric MetricOp
	// number of clusters
	K int
	// Distances between data points [m x m]
	D *Distances
	// Matrix of centroids	
	Centers Matrix
	// Total distance of members to each centroid
	Errors Vector
	// cluster center assignment index
	Clusters []int
	// cost
	Cost float64
	// Maximum number of iterations
	MaxIter int
	// ordered index of elements subset
	index []int
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
	if i == 0 {
		// convergence is achieved right after initialization...
		// run maximization at least once to calculate cost
		c.maximization()
	}

	// copy classifcation information
	classes = &Classes{
		make([]int, len(c.X)), k, c.Cost}
	copy(classes.Index, c.Clusters)

	return
}

func (c *KMeans) Len() int {
	return len(c.index)
}

func (c *KMeans) Segregations(classes *Classes) (S Matrix) {
	if c.D == nil {
		c.D = NewDistances(c.X, c.Metric)
	}
	S = Segregations(c.D.rep, classes)
	return
}

func (c *KMeans) Subset(index []int) Splitter {
	X := Matrix(mlgo.Matrix(c.X).Slice(index))
	var D *Distances
	if c.D != nil {
		D = c.D.Subset(index)
	}
	d := NewKMeans(X, c.Metric);
	d.D = D
	return d
}

// initialize the cluster centroids by randomly selecting data points
func (c *KMeans) initialize() {
	c.Centers, c.Errors = make(Matrix, c.K), make(Vector, c.K)

	// initialize index to complete index iterating over all elements of X, unless defined otherwise
	if c.index == nil {
		c.index = mlgo.Range(0, len(c.X))
	}
	m := len(c.index)

	c.Clusters = make([]int, m)

	activeSet := NewActiveSet(m)
	for k, _ := range c.Centers {
		i := activeSet.Get( rand.Intn(activeSet.Len()) )
		x := c.X[ c.index[i] ]
		activeSet.Remove(i)
		// copy data vector
		c.Centers[k] = make(Vector, len(x))
		copy(c.Centers[k], x)
	}
}

// expectation step: assign data points to cluster centroids
// Returns whether the algorithm has converged
func (c *KMeans) expectation() (converged bool) {
	// find the centroids that is closest to the current data point
	assign := func(i int, chClusters chan int) {
		clusters, min := 0, maxValue
		// find the center with the minimum distance
		for ii := 0; ii < len(c.Centers); ii++ {
			distance := c.Metric(c.X[i], c.Centers[ii])
			if distance < min {
				clusters, min = ii, distance
			}
		}
		chClusters <- clusters
	}

	// process examples concurrently
	ch := make(chan int)
	for i, _ := range c.X {
		go assign(i, ch)
	}

	// collect results
	converged = true
	for i, _ := range c.X {
		if clusters := <-ch; c.Clusters[i] != clusters {
			c.Clusters[i] = clusters
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
		memberIdx := make([]int, len(c.Clusters))
		for i, class := range c.Clusters {
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
