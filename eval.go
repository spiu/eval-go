package eval

import (
	"context"
	"fmt"
)

// Instance represents a single evaluation instance with reference and prediction texts
type Instance struct {
	Reference  string
	Prediction string
}

// PairwiseMetricFunc is a function that computes scores by comparing references and predictions
type PairwiseMetricFunc func(ctx context.Context, references, predictions []string) ([]float64, error)

// PointwiseMetricFunc is a function that computes scores for predictions
type PointwiseMetricFunc func(ctx context.Context, predictions []string) ([]float64, error)

// PairwiseMetric represents a metric that compares reference and prediction
type PairwiseMetric struct {
	Name        string
	Description string
	Compute     PairwiseMetricFunc
}

// PointwiseMetric represents a metric that evaluates a prediction
type PointwiseMetric struct {
	Name        string
	Description string
	Compute     PointwiseMetricFunc
}

// Evaluation represents a set of metrics to be evaluated
type Evaluation struct {
	Name              string
	Description       string
	pairwiseMetrics   []PairwiseMetric
	pointwiseMetrics  []PointwiseMetric
}

// Result represents the output of an evaluation
type Result struct {
	EvaluationName string
	MetricResults  map[string]float64
	Error         error
}

// NewEvaluation creates a new evaluation with the given name and description
func NewEvaluation(name, description string, pairwiseMetrics []PairwiseMetric, pointwiseMetrics []PointwiseMetric) *Evaluation {
	return &Evaluation{
		Name:              name,
		Description:       description,
		pairwiseMetrics:   pairwiseMetrics,
		pointwiseMetrics:  pointwiseMetrics,
	}
}

// Run executes the evaluation on the given instances and predictions
func (e *Evaluation) Run(ctx context.Context, instances []Instance, predictions []string) ([]Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Extract references and predictions from instances
	references := make([]string, len(instances))
	for i, instance := range instances {
		references[i] = instance.Reference
	}

	// Run pairwise metrics
	results := make([]Result, len(instances))
	for i := range instances {
		results[i] = Result{
			Instance:      instances[i],
			MetricResults: make(map[string]float64),
		}
	}

	// Run pairwise metrics
	for _, metric := range e.pairwiseMetrics {
		scores, err := metric.Compute(ctx, references, predictions)
		if err != nil {
			return nil, fmt.Errorf("pairwise metric %s failed: %w", metric.Name, err)
		}

		if len(scores) != len(instances) {
			return nil, fmt.Errorf("pairwise metric %s returned %d scores, expected %d", metric.Name, len(scores), len(instances))
		}

		for i, score := range scores {
			results[i].MetricResults[metric.Name] = score
		}
	}

	// Run pointwise metrics
	for _, metric := range e.pointwiseMetrics {
		scores, err := metric.Compute(ctx, predictions)
		if err != nil {
			return nil, fmt.Errorf("pointwise metric %s failed: %w", metric.Name, err)
		}

		if len(scores) != len(predictions) {
			return nil, fmt.Errorf("pointwise metric %s returned %d scores, expected %d", metric.Name, len(scores), len(predictions))
		}

		for i, score := range scores {
			results[i].MetricResults[metric.Name] = score
		}
	}

	return results, nil
}

// NewPairwiseMetric creates a new pairwise metric
func NewPairwiseMetric(name, description string, compute PairwiseMetricFunc) PairwiseMetric {
	return PairwiseMetric{
		Name:        name,
		Description: description,
		Compute:     compute,
	}
}

// NewPointwiseMetric creates a new pointwise metric
func NewPointwiseMetric(name, description string, compute PointwiseMetricFunc) PointwiseMetric {
	return PointwiseMetric{
		Name:        name,
		Description: description,
		Compute:     compute,
	}
} 
} 