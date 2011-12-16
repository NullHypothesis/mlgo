package cluster

import (
	"testing"
	"fmt"
)

func TestMixModel(t *testing.T) {
	K := 2
	X := Matrix{
			{-10, -10},
			{-10,  -8},
			{ -8,  -8},
			{ -8, -10},
			{ 10,  10},
			{ 10,   8},
			{  8,   8},
			{  8,  10} }
	
	c := MixModel{X:X}
	classes := c.Cluster(K)

	fmt.Println(classes)
	fmt.Println(c.Mixings)
	fmt.Println(c.Means)
	fmt.Println(c.Variances)
	fmt.Println(c.posteriors)
	fmt.Println(c.NLogLikelihood)

}

