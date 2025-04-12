package eval

import (
	"context"
	"fmt"
)

// PairwiseMetric represents a metric that compares reference and prediction
type PairwiseMetric struct {
	Name        string
	Description string
	compute     PairwiseMetricFunc
}

// PointwiseMetric represents a metric that evaluates a prediction
type PointwiseMetric struct {
	Name        string
	Description string
	compute     PointwiseMetricFunc
}

// Compute executes the pairwise metric on the given references and predictions
func (m *PairwiseMetric) Compute(ctx context.Context, references, predictions []string) ([]float64, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if len(references) == 0 {
		return nil, fmt.Errorf("no references provided")
	}

	if len(references) != len(predictions) {
		return nil, fmt.Errorf("number of references (%d) does not match number of predictions (%d)", 
			len(references), len(predictions))
	}

	return m.compute(ctx, references, predictions)
}

// Compute executes the pointwise metric on the given predictions
func (m *PointwiseMetric) Compute(ctx context.Context, predictions []string) ([]float64, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if len(predictions) == 0 {
		return nil, fmt.Errorf("no predictions provided")
	}

	return m.compute(ctx, predictions)
}

// NewPairwiseMetric creates a new pairwise metric
func NewPairwiseMetric(name, description string, compute PairwiseMetricFunc) PairwiseMetric {
	return PairwiseMetric{
		Name:        name,
		Description: description,
		compute:     compute,
	}
}

// NewPointwiseMetric creates a new pointwise metric
func NewPointwiseMetric(name, description string, compute PointwiseMetricFunc) PointwiseMetric {
	return PointwiseMetric{
		Name:        name,
		Description: description,
		compute:     compute,
	}
} 