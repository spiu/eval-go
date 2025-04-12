package eval

import (
	"context"
	"fmt"
)

// Evaluation represents a set of metrics to be evaluated
type Evaluation struct {
	Name             string
	Description      string
	pairwiseMetrics  []PairwiseMetric
	pointwiseMetrics []PointwiseMetric
}

// NewEvaluation creates a new evaluation with the given name and description
func NewEvaluation(name, description string, pairwiseMetrics []PairwiseMetric, pointwiseMetrics []PointwiseMetric) *Evaluation {
	return &Evaluation{
		Name:             name,
		Description:      description,
		pairwiseMetrics:  pairwiseMetrics,
		pointwiseMetrics: pointwiseMetrics,
	}
}

// Run executes the evaluation on the given instances and predictions
func (e *Evaluation) Run(ctx context.Context, instances []Instance, predictions []string) ([]Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, fmt.Errorf("no instances provided")
	}

	if len(instances) != len(predictions) {
		return nil, fmt.Errorf("number of instances (%d) does not match number of predictions (%d)", 
			len(instances), len(predictions))
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

		for i, score := range scores {
			results[i].MetricResults[metric.Name] = score
		}
	}

	return results, nil
} 