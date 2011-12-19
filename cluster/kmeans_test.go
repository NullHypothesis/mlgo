package cluster

import (
	"testing"
	"fmt"
)

func TestKMeans(t *testing.T) {
	K := 2
	X := Matrix{
			{-100, -200},  // outlier
			{-10, -20},
			{-10, -18},
			{ -8, -18},
			{ -8, -20},
			{ 10,  20},
			{ 10,  18},
			{  8,  18},
			{  8,  20} }
	
	c := KMeans{X:X}
	classes := c.Cluster(K)

	fmt.Println(classes)
	fmt.Println(c.Centers)
	fmt.Println(c.Errors)

	c2 := KMedians{KMeans{X:X}}
	classes2 := c2.Cluster(K)
	fmt.Println(classes2)
	fmt.Println(c2.Centers)
	fmt.Println(c2.Errors)

}

