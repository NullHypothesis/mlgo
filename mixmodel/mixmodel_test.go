package mixmodel

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
	
	clusters := Run(X, K, 1)

	fmt.Println(clusters.Mixings)
	fmt.Println(clusters.Means)
	fmt.Println(clusters.Variances)
	fmt.Println(clusters.LogLikelihood)

}
