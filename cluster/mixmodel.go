package cluster

import (
	"rand"
	"math"
	"fmt"
	"mlgo/mlgo"
)

type MixModel struct {
	// Matrix of data points [m x n]
	X Matrix
	// number of clusters
	K int
	// Matrix of posterior probabilities [m x k]
	posteriors Matrix
	// Matrix of Gaussians [k x n]
	Means, Variances Matrix
	// Vector of mixing proportions [k]
	Mixings Vector
	// Negative likelihood to be minimized
	NLogLikelihood float64
	// Maximum number of iterations
	MaxIter int
}

const logProbEpsilon = 0.01

// Cluster runs the algorithm once with random initialization
// Returns the classification information
func (c *MixModel) Cluster(k int) (classes *Classes) {
	if c.X == nil { return }
	c.K = k
	c.initialize()
	i := 0
	for !c.expectation() && (c.MaxIter == 0 || i < c.MaxIter) {
		c.maximization()
		i++
	}

	// copy classification information
	classes = &Classes{
		make([]int, len(c.X)), c.NLogLikelihood }
	for i, pp := range c.posteriors {
		class, maxPosterior := 0, 0.0
		for k, p := range pp {
			if p > maxPosterior {
				maxPosterior = p
				class = k
			}
		}
		classes.Index[i] = class
	}

	return
}

func summary(X Matrix) (means, variances Vector) {
	m := len(X)
	if m < 2 { return }

	n := len(X[0])
	stats := make([]mlgo.Summary, n)

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

// initialize Gaussians randomly
func (c *MixModel) initialize() {
	means, variances := summary(c.X)
	m, n := len(c.X), len(means)

	c.Means, c.Variances, c.Mixings = make(Matrix, c.K), make(Matrix, c.K), make(Vector, c.K)

	c.posteriors = make(Matrix, m)

	c.NLogLikelihood = maxValue

	for i := 0; i < m; i++ {
		c.posteriors[i] = make(Vector, c.K)
	}

	for k := 0; k < c.K; k++ {

		c.Means[k], c.Variances[k] = make(Vector, n), make(Vector, n)

		// use large variance: variance of each feature based on entire data set
		copy(c.Variances[k], variances)

		// use mean of each feature plus some noise
		for j := 0; j < len(means); j++ {
			sd := math.Sqrt(variances[j])
			c.Means[k][j] = means[j] + (rand.Float64() * sd - sd/2)
		}

		// uniform mixing proportions
		c.Mixings[k] = 1/float64(c.K)

	}
	fmt.Println("initial: ", c.Means, c.Variances, c.Mixings)
}

type pdf func(float64) float64

func normPdf(mu, sigma2 float64) func(float64) float64 {
	return func(x float64) float64 {
		d := x-mu
		return 1 / math.Sqrt(2*math.Pi*sigma2) * math.Exp(-d*d / (2*sigma2))
	}
}

// expectation step: assign data points to cluster meanss
// Returns whether the algorithm has converged
func (c *MixModel) expectation() (converged bool) {

	// setup Gaussians
	pnorms := make([][]pdf, c.K)
	for k := 0; k < c.K; k++ {
		pnorms[k] = make([]pdf, len(c.Means[k]))
		for d := 0; d < len(c.Means[k]); d++ {
			pnorms[k][d] = normPdf(c.Means[k][d], c.Variances[k][d])
		}
	}

	// Calculate the posterior probability densities of each data point
	//   being generated by each Gaussian using Bayes theorem
	// Also calculate the negative log likelihood of the model
	//   (i.e. the probability of entire data given the mixture model)
	model := 0.0
	for i, _ := range c.X {
		px := 0.0
		for k := 0; k < c.K; k++ {
			likelihood := 1.0
			for d := 0; d < len(c.X[i]); d++ {
				likelihood *= pnorms[k][d](c.X[i][d])
			}
			p := c.Mixings[k] * likelihood
			c.posteriors[i][k] = p
			px += p
		}
		// normalize posterior
		for k := 0; k < c.K; k++ {
			c.posteriors[i][k] /= px
		}
		model -= math.Log(px)
	}

	// Check that model negative log likelihood is decreasing
	// (negative log likelihood is guaranteed to be non-increasing).
	// If the current value does not differ from the previous,
	// the algorithm has converged (possibly to a local minimum).
	if c.NLogLikelihood - model > logProbEpsilon {
		converged = false
	} else {
		converged = true
	}
	c.NLogLikelihood = model

	return
}

// maximization step: move cluster meanss to centroids of data points
// Returns the cost
func (c *MixModel) maximization() {

	for k := 0; k < c.K; k++ {
		
		// Compute new mixing proportions
		sum := 0.0
		for i, _ := range c.posteriors {
			sum += c.posteriors[i][k]
		}
		c.Mixings[k] = sum / float64( len(c.posteriors) )

		// Compute new means
		for d, _ := range c.Means[k] {
			a, b := 0.0, 0.0
			for i, _ := range c.posteriors {
				p := c.posteriors[i][k]
				a += p * c.X[i][d]
				b += p
			}
			c.Means[k][d] = a / b
		}

		//TODO calculate variance more efficiently
		// Compuete new variances
		for d, _ := range c.Variances[k] {
			a, b := 0.0, 0.0
			for i, _ := range c.posteriors {
				p := c.posteriors[i][k]
				diff := c.X[i][d] - c.Means[k][d]
				a += p * (diff * diff)
				b += p
			}
			c.Variances[k][d] = a / b
		}

	}

}

