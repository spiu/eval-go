package eval

import "context"

// Instance represents a single evaluation instance with reference and prediction texts
type Instance struct {
	Reference  string
	Prediction string
}

// PairwiseResult represents the output of a pairwise evaluation
type PairwiseResult struct {
	Instance      Instance
	MetricResults map[string]float64
}

// PointwiseResult represents the output of a pointwise evaluation
type PointwiseResult struct {
	Prediction    string
	MetricResults map[string]float64
}

// PairwiseMetricFunc is a function that computes scores by comparing references and predictions
type PairwiseMetricFunc func(ctx context.Context, references, predictions []string) ([]float64, error)

// PointwiseMetricFunc is a function that computes scores for predictions
type PointwiseMetricFunc func(ctx context.Context, predictions []string) ([]float64, error)

// PairwiseScoreFunc is a function that determines how to calculate the score between reference and prediction scores
type PairwiseScoreFunc func(referenceScore, predictionScore float64) float64 