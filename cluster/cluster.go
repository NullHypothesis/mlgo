package cluster

import (
	"mlgo/mlgo"
)

type Vector mlgo.Vector
type Matrix mlgo.Matrix

const maxValue = mlgo.MaxValue

type Classes struct {
	// classification index
	Index Partitions
	Cost float64
}

type Clusterer interface {
	Cluster(k int) (classes *Classes)
}

// FindClusters runs the clustering algorithm for the specified number of repeats.
func FindClusters(c Clusterer, k int, repeats int) (classes *Classes) {
	// repeat clustering concurrently
	ch := make(chan *Classes)
	for i := 0; i < repeats; i++ {
		go func() {
			ch <- c.Cluster(k)
		}()
	}

	// determine best clustering
	minCost := maxValue
	for i := 0; i < repeats; i++ {
		if cl := <-ch; cl.Cost < minCost {
			classes = cl
			minCost = cl.Cost
		}
	}
	return
}

