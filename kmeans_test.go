package kmeans

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
	
	clusters := Run(X, K, 10)

	fmt.Println(clusters.Classes)
	fmt.Println(clusters.Means)

}
