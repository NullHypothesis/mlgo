package cluster

import (
	"mlgo/base"
	"testing"
)

var silhouetteTests = []struct {
	x Matrix
	metric MetricOp
	classes *Classes
	silhouettes Vector
}{
	{
		Matrix{
			{101, 102, 103}, {102, 103, 104}, {103, 104, 105},
			{111, 112, 113}, {112, 113, 114}, {113, 114, 115},
			{ 21,  22,  23}, { 22,  23,  24}, { 23,  24,  25},
			{ 29,  30,  31}, { 32,  33,  34}, { 33,  34,  35},
		},
		Manhattan,
		&Classes{ Index:Partitions{0, 0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3}, K:4 },
		Vector{
			0.8636364, 0.9000000, 0.8333333, 
			0.8333333, 0.9000000, 0.8636364,
			0.8548387, 0.8928571, 0.8200000,
			0.5000000, 0.8000000, 0.7727273,
		},
	},
}

func TestSilhouettes(t *testing.T) {
	for i, test := range silhouetteTests {
		d := Distances(test.x, test.metric)
		sil := Silhouettes( Segregations(d, test.classes), test.classes.Index )
		if !mlgo.Vector(test.silhouettes).Equal(mlgo.Vector(sil)) {
			t.Errorf("#%d Silhouettes(Segregations(...), ...) got %v, want %v", i, sil, test.silhouettes)
		}
	}
}

var shadowTests = []struct {
	x, centers Matrix
	metric MetricOp
	index Partitions
	shadows Vector
}{
	{
		Matrix{
			{101, 102, 103}, {102, 103, 104}, {103, 104, 105},
			{111, 112, 113}, {112, 113, 114}, {113, 114, 115},
			{ 21,  22,  23}, { 22,  23,  24}, { 23,  24,  25},
			{ 29,  30,  31}, { 32,  33,  34}, { 33,  34,  35},
		},
		Matrix{
			{102, 103, 104},
			{112, 113, 114},
			{ 22,  23,  24},
			{ 32,  33,  34},
		},
		Manhattan,
		Partitions{0, 0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3},
		Vector{
			0.9090909, 1.0000000, 0.8888888,
			0.8888888, 1.0000000, 0.9090909,
			0.9090909, 1.0000000, 0.8888888,
			0.5714286, 1.0000000, 0.9090909,
		},
	},
}

func TestShadows(t *testing.T) {
	for i, test := range shadowTests {
		S := Separations(test.x, test.centers, test.metric)
		shadows := Silhouettes(S, test.index)
		if !mlgo.Vector(test.shadows).Equal(mlgo.Vector(shadows)) {
			t.Errorf("#%d Silhouettes(Separations(...), ...) got %v, want %v", i, shadows, test.shadows)
		}
	}
}

