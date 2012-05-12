package cluster

import (
	"mlgo/base"
)

// Validation measures

// Silhouette considerations
// use different linkage methods to pre-calculate distance between an element and a cluster
// pre-calculate these linkage distances for efficiency
// use update formula?

// Segregations return a matrix of distances between data points and clusters
func Segregations(distances Matrix, classes *Classes) (S Matrix) {
	// each row of x is considered one data point
	m := len(distances)

	// allocate space
	S = make(Matrix, m)
	for i := 0; i < m; i++ {
		S[i] = make(Vector, classes.K)
	}

	index := classes.Index

	// determine cluster sizes
	sizes := classes.Sizes()

	// calculate the average distances from data point i to data points, for each cluster
	// TODO option to aggregate by median/max/min instead of mean?
	for i := 0; i < m; i++ {
		// accumulate sum
		for j := 0; j < m; j++ {
			S[i][ index[j] ] += distances[i][j]
		}
		// derive mean via division by cluster sizes
		for jj := 0; jj < classes.K; jj++ {
			S[i][jj] /= float64(sizes[jj])
		}
		// correct mean for own cluster (divide by size-1 instead of size)
		c := index[i]
		size := float64(sizes[c])
		S[i][c] *= size / (size - 1)
	}
	return
}

// Separations return a matrix of distances between data points and cluster centers
func Separations(X, centers Matrix, metric MetricOp) (S Matrix) {
	// each row of x is considered one data point
	m := len(X)
	k := len(centers)

	// allocate space
	S = make(Matrix, m)
	for i := 0; i < m; i++ {
		S[i] = make(Vector, k)
	}

	// calculate distance from data point i to center of cluster j
	for i := 0; i < m; i++ {
		for j := 0; j < k; j++ {
			S[i][j] = metric(X[i], centers[j])
		}
	}
	return
}

// Silhouettes returns a vector of silhouettes for data points.
// If S is a segregation matrix, then the returned values are conventionally considered as silhouettes.
// If S is a separation matrix, then the returned values are can be considered as shadows.
// TODO special case: i is in a singleton cluster; silhouette should be 0 intead of 1
// TODO special case: silhouette is not defined for two singleton clusters
// TODO faithful calculation of "shadow" as defined by Friedrich Leisch (average two nearest centroids for 'b')
func Silhouettes(S Matrix, index Partitions) (s Vector)  {
	m := len(S)
	k := len(S[0])

	s = make(Vector, m)

	// calculate silouettes
	for i := 0; i < m; i++ {
		// distance to own cluster
		c := index[i]
		a := S[i][c]
		// distance to nearest cluster
		b := maxValue
		for j := 0; j < k; j++ {
			if j != c && S[i][j] < b {
				b = S[i][j]
			}
		}
		max := a
		if a < b {
			max = b
		}
		s[i] = (b - a) / max
	}
	return
}


type Split struct {
	K int
	Cost float64
}

func Mean(x Vector) (m float64) {
	for i := 0; i < len(x); i++ {
		m += x[i]
	}
	m /= float64(len(x))
	return
}

func SplitByAvgSil(distances Matrix, clust Clusterer, K int) (s Split) {
	m := len(distances)

	// silhouette can only be calculated for 2 <= k <= m - 1

	if K <= 0 || K > m - 1 {
		K = m - 1
	}

	// maximize average silhouette
	avgSil := -1.0
	opt_k := 0
	for k := 2; k <= K; k++ {
		classes := clust.Cluster(k)
		S := Segregations(distances, classes)
		sil := Silhouettes(S, classes.Index)
		t := Mean(sil)
		if t > avgSil {
			avgSil = t
			opt_k = k
		}
	}

	s.K = opt_k
	s.Cost = 1 - avgSil
	return
}

//TODO allow functions to take an index vector s.t. Xsub and distancesSub do not need to be copied?

// K is the maximum number of clusters.
// L is the maximum number of children clusters for any cluster.
func SplitByAvgSplitSil(distances Matrix, clust Clusterer, K, L int) (s Split) {
	m := len(distances)

	// average split silhouette can be only be calculated for 1 <= k <= m/3
	// if k > m/3, at least one cluster would have < 3 elements
	// each cluster needs >= 3 elements to be further split into at least 2 clusters 
	//  for silhouette calculation

	if K <= 0 || K > m / 3 {
		K = m / 3
	}

	avgSplitSil := maxValue
	opt_k := 0
	splitSil := make(Vector, K)
	for k := 1; k <= K; k++ {
		classes := clust.Cluster(k)
		partitions := classes.Partitions()
		for kk := 0; kk < classes.K; kk++ {
			distancesSub := Matrix(mlgo.Matrix(distances).Slice(partitions[kk]))
			clustSplit := SplitByAvgSil(distancesSub, clust, L)
			splitSil[kk] = 1 - clustSplit.Cost
		}
		t := Mean(splitSil)
		if t < avgSplitSil {
			avgSplitSil = t
			opt_k = k
		}
	}

	s.K = opt_k
	s.Cost = avgSplitSil
	return
}

