package eval

import (
	"context"
	"fmt"
)

// PairwiseEvaluation represents a set of metrics that compare references and predictions
type PairwiseEvaluation struct {
	Name    string
	Description string
	metrics []PairwiseMetric
}

// PointwiseEvaluation represents a set of metrics that evaluate predictions
type PointwiseEvaluation struct {
	Name    string
	Description string
	metrics []PointwiseMetric
}

// NewPairwiseEvaluation creates a new pairwise evaluation
func NewPairwiseEvaluation(name, description string, metrics []PairwiseMetric) *PairwiseEvaluation {
	return &PairwiseEvaluation{
		Name:    name,
		Description: description,
		metrics: metrics,
	}
}

// NewPointwiseEvaluation creates a new pointwise evaluation
func NewPointwiseEvaluation(name, description string, metrics []PointwiseMetric) *PointwiseEvaluation {
	return &PointwiseEvaluation{
		Name:    name,
		Description: description,
		metrics: metrics,
	}
}

// Run executes the pairwise evaluation on the given instances
func (e *PairwiseEvaluation) Run(ctx context.Context, instances []Instance) ([]PairwiseResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, fmt.Errorf("no instances provided")
	}

	// Extract references and predictions from instances
	references := make([]string, len(instances))
	predictions := make([]string, len(instances))
	for i, instance := range instances {
		references[i] = instance.Reference
		predictions[i] = instance.Prediction
	}

	// Initialize results
	results := make([]PairwiseResult, len(instances))
	for i := range instances {
		results[i] = PairwiseResult{
			Instance:      instances[i],
			MetricResults: make(map[string]float64),
		}
	}

	// Run metrics
	for _, metric := range e.metrics {
		scores, err := metric.Compute(ctx, references, predictions)
		if err != nil {
			return nil, fmt.Errorf("metric %s failed: %w", metric.Name, err)
		}

		for i, score := range scores {
			results[i].MetricResults[metric.Name] = score
		}
	}

	return results, nil
}

// Run executes the pointwise evaluation on the given predictions
func (e *PointwiseEvaluation) Run(ctx context.Context, predictions []string) ([]PointwiseResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if len(predictions) == 0 {
		return nil, fmt.Errorf("no predictions provided")
	}

	// Initialize results
	results := make([]PointwiseResult, len(predictions))
	for i, prediction := range predictions {
		results[i] = PointwiseResult{
			Prediction:    prediction,
			MetricResults: make(map[string]float64),
		}
	}

	// Run metrics
	for _, metric := range e.metrics {
		scores, err := metric.Compute(ctx, predictions)
		if err != nil {
			return nil, fmt.Errorf("metric %s failed: %w", metric.Name, err)
		}

		for i, score := range scores {
			results[i].MetricResults[metric.Name] = score
		}
	}

	return results, nil
} 