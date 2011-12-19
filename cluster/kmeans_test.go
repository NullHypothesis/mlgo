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

	/*
	X := Matrix{
			{-10, -20},
			{-10, -18},
			{ -8, -18},
			{ -8, -20},
			{ 10,  20},
			{ 10,  18},
			{  8,  18},
			{  8,  20},
	}
	*/

	/*
	metric := func(a, b Vector) float64 {
		return Minkowski(a, b, 0.5)
	}
	*/

	metric := Euclidean
	
	c := NewKMeans(X, metric)
	classes := c.Cluster(K)

	fmt.Println(classes)
	fmt.Println(c.Centers)
	fmt.Println(c.Errors)

	c2 := NewKMedians(X, metric)
	classes2 := c2.Cluster(K)
	fmt.Println(classes2)
	fmt.Println(c2.Centers)
	fmt.Println(c2.Errors)

	c3 := NewKMedoids(X, metric)
	classes3 := c3.Cluster(K)
	fmt.Println(classes3)
	fmt.Println(c3.Centers)
	fmt.Println(c3.Errors)

}

