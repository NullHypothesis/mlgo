package cluster

import (
	"testing"
	"fmt"
)

func TestKMeans(t *testing.T) {
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
	
	c := KMeans{X:X}
	classes := c.Cluster(K)

	fmt.Println(classes)
	fmt.Println(c.Means)

}

