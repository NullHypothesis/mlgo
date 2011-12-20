package mlgo

import (
	"math"
)

type Vector []float64
type Matrix [][]float64

const (
	MaxValue = math.MaxFloat64
)

func (X Matrix) Summarize() (means, variances Vector) {
	m := len(X)
	if m < 2 { return }

	n := len(X[0])
	stats := make([]Summary, n)

	means, variances = make(Vector, n), make(Vector, n)

	for i := 0; i < m; i++ {
		// accumulate statistics for each feature
		for j, x := range X[i] {
			stats[j].Add(x)
		}
	}

	for j, _ := range stats {
		means[j] = stats[j].Mean
		variances[j] = stats[j].VarP()
	}

	return
}

func (x Vector) Summarize() (mean, variance float64) {
	var stats Summary
	for _, v := range x {
		// accumulate statistics
		stats.Add(v)
	}

	mean, variance = stats.Mean, stats.VarP()
	return
}

