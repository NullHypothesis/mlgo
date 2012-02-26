package cluster

import (
	"testing"
)

var kmediansTests = []struct {
	x Matrix
	metric MetricOp
	k int
	index Partitions
	centers Matrix
}{
	{
		Matrix{
			{-100, -200},
			{-10, -20},
			{-10, -18},
			{ -8, -18},
			{ -8, -20},
			{ 10,  20},
			{ 10,  18},
			{  8,  18},
			{  8,  20},
		},
		Euclidean,
		2,
		Partitions{0, 0, 0, 0, 0, 1, 1, 1, 1},
		Matrix{
			{-10, -20},
			{9, 19},
		},
	},
}

func TestKMedians(t *testing.T) {
	for i, test := range kmediansTests {
		c := NewKMedians(test.x, test.metric)
		classes := c.Cluster(test.k)
		if !classes.Index.Equal(test.index) {
			t.Errorf("#%d KMedians.Cluster(...) got %v, want %v", i, classes.Index, test.index)
		}
		if !CoordinatesSetEqual(c.Centers, test.centers) {
			t.Errorf("#%d KMedians.Cluster(...) got %v, want %w", i, c.Centers, test.centers)
		}
	}
}

