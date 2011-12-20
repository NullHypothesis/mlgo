package cluster

import (
	"math"
)

type MetricOp func(a, b Vector) float64

// EuclideanSq returns the Euclidean squared distance metric between points a and b
func EuclideanSq(a, b Vector) (d float64) {
	if len(a) != len(b) { return }
	for i := 0; i < len(a); i++ {
		t := b[i] - a[i]
		d += t*t
	}
	return
}

// Euclidean returns the Euclidean distance metric between points a and b
func Euclidean(a, b Vector) (d float64) {
	if len(a) != len(b) { return }
	for i := 0; i < len(a); i++ {
		t := b[i] - a[i]
		d += t*t
	}
	d = math.Sqrt(d)
	return
}

// Manhattan returns the Manhattan distance metric beetween points a and b
func Manhattan(a, b Vector) (d float64) {
	if len(a) != len(b) { return }
	for i := 0; i < len(a); i++ {
		d += math.Fabs(b[i] - a[i])
	}
	return
}

// Chebyshev returns the Chebyshev distance metric between points a and b
func Chebyshev(a, b Vector) (d float64) {
	if len(a) != len(b) { return }
	for i := 0; i < len(a); i++ {
		t := math.Fabs(b[i] - a[i])
		if t > d { d = t }
	}
	return
}

// Minkowski returns the Minkowski distance metric between points a and b
func Minkowski(a, b Vector, p float64) (d float64) {
	if len(a) != len(b) { return }
	for i := 0; i < len(a); i++ {
		d += math.Pow(math.Fabs(b[i] - a[i]), p)
	}
	d = math.Pow(d, 1/p)
	return
}

func Distances(X Matrix, metric MetricOp) (D Matrix) {
	// each row of X is considered one data point
	m := len(X)

	D = make(Matrix, m)
	for i := 0; i < m; i++ {
		D[i] = make(Vector, m)
		for j := 0; j < m; j++ {
			if i != j {
				D[i][j] = metric(X[i], X[j])
			}
		}
	}
	return
}

