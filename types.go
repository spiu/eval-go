package eval

import "context"

// Instance represents a single evaluation instance with reference and prediction texts
type Instance struct {
	Reference  string
	Prediction string
}

// Result represents the output of an evaluation
type Result struct {
	Instance      Instance
	MetricResults map[string]float64
}

// PairwiseMetricFunc is a function that computes scores by comparing references and predictions
type PairwiseMetricFunc func(ctx context.Context, references, predictions []string) ([]float64, error)

// PointwiseMetricFunc is a function that computes scores for predictions
type PointwiseMetricFunc func(ctx context.Context, predictions []string) ([]float64, error) 