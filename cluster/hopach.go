package cluster

type Hopacher interface {
	Clusterer
	Subset(index []int) Hopacher
	// Heterogeneity of a given partitioning
	Heterogeneity(classes *Classes) float64
	// Sort elements in some order
	Sort()
}

type Hopach struct {
	Base Splitter

	maxK, maxL int
	// implemented parameters
	// clusters = best
	// coll = seq
	// newmed = nn
	// mss = med
	// initord = co
	// ord = neighbour
}

func NewHopach() {

}

