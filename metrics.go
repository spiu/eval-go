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

// ToPairwise converts a pointwise metric into a pairwise one
// by allowing custom logic to determine the score between reference and prediction
func (m *PointwiseMetric) ToPairwise(scoreFunc PairwiseScoreFunc) PairwiseMetric {
	return NewPairwiseMetric(
		m.Name,
		m.Description,
		func(ctx context.Context, references, predictions []string) ([]float64, error) {
			// Get scores for references
			referenceScores, err := m.Compute(ctx, references)
			if err != nil {
				return nil, err
			}
			
			// Get scores for predictions
			predictionScores, err := m.Compute(ctx, predictions)
			if err != nil {
				return nil, err
			}
			
			// Apply the custom scoring function to each pair
			scores := make([]float64, len(references))
			for i := range references {
				scores[i] = scoreFunc(referenceScores[i], predictionScores[i])
			}
			
			return scores, nil
		},
	)
}

// Default scoring functions for converting pointwise metrics to pairwise

// DifferenceScore calculates the difference between prediction and reference scores
func DifferenceScore(referenceScore, predictionScore float64) float64 {
	return predictionScore - referenceScore
}

// RatioScore calculates the ratio between prediction and reference scores
func RatioScore(referenceScore, predictionScore float64) float64 {
	if referenceScore == 0 {
		return 0 // Avoid division by zero
	}
	return predictionScore / referenceScore
}

// AbsoluteDifferenceScore calculates the absolute difference between prediction and reference scores
func AbsoluteDifferenceScore(referenceScore, predictionScore float64) float64 {
	diff := predictionScore - referenceScore
	if diff < 0 {
		return -diff
	}
	return diff
}

// MaxScore returns the maximum of the reference and prediction scores
func MaxScore(referenceScore, predictionScore float64) float64 {
	if referenceScore > predictionScore {
		return referenceScore
	}
	return predictionScore
}

// MinScore returns the minimum of the reference and prediction scores
func MinScore(referenceScore, predictionScore float64) float64 {
	if referenceScore < predictionScore {
		return referenceScore
	}
	return predictionScore
}

// AverageScore calculates the average of the reference and prediction scores
func AverageScore(referenceScore, predictionScore float64) float64 {
	return (referenceScore + predictionScore) / 2
} 