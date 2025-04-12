package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/snpu/eval-go"
	"github.com/snpu/eval-go/metrics"
)

func main() {
	// Create pairwise metrics
	pairwiseMetrics := []eval.PairwiseMetric{
		metrics.StringSimilarity(),
		metrics.LengthRatio(),
		metrics.WordOverlap(),
	}
	
	// Create pointwise metrics
	pointwiseMetrics := []eval.PointwiseMetric{
		metrics.KeywordPresence(),
	}

	// Convert pointwise metrics to pairwise with different scoring strategies
	convertedPairwiseMetrics := make([]eval.PairwiseMetric, len(pointwiseMetrics))
	for i, metric := range pointwiseMetrics {
		// Example 1: Use the difference between reference and prediction scores
		convertedPairwiseMetrics[i] = metric.ToPairwise(metrics.DifferenceScore)
	}

	// Create another set of converted metrics with a different scoring strategy
	alternativePairwiseMetrics := make([]eval.PairwiseMetric, len(pointwiseMetrics))
	for i, metric := range pointwiseMetrics {
		// Example 2: Use the ratio between prediction and reference scores
		alternativePairwiseMetrics[i] = metric.ToPairwise(metrics.RatioScore)
	}

	// Create a third set with absolute difference scoring
	absoluteDifferenceMetrics := make([]eval.PairwiseMetric, len(pointwiseMetrics))
	for i, metric := range pointwiseMetrics {
		absoluteDifferenceMetrics[i] = metric.ToPairwise(metrics.AbsoluteDifferenceScore)
	}

	// Create a fourth set with average scoring
	averageMetrics := make([]eval.PairwiseMetric, len(pointwiseMetrics))
	for i, metric := range pointwiseMetrics {
		averageMetrics[i] = metric.ToPairwise(metrics.AverageScore)
	}

	// Create pairwise evaluation
	pairwiseEval := eval.NewPairwiseEvaluation(
		"llm_response_comparison",
		"Evaluates changes in LLM responses by comparing with references",
		pairwiseMetrics,
	)

	// Create pointwise evaluation
	pointwiseEval := eval.NewPointwiseEvaluation(
		"llm_response_quality",
		"Evaluates the quality of LLM responses",
		pointwiseMetrics,
	)

	// Create evaluations with different conversion strategies
	differenceEval := eval.NewPairwiseEvaluation(
		"llm_response_difference",
		"Evaluates the difference between reference and prediction scores",
		convertedPairwiseMetrics,
	)

	ratioEval := eval.NewPairwiseEvaluation(
		"llm_response_ratio",
		"Evaluates the ratio between prediction and reference scores",
		alternativePairwiseMetrics,
	)

	absoluteDifferenceEval := eval.NewPairwiseEvaluation(
		"llm_response_absolute_difference",
		"Evaluates the absolute difference between reference and prediction scores",
		absoluteDifferenceMetrics,
	)

	averageEval := eval.NewPairwiseEvaluation(
		"llm_response_average",
		"Evaluates the average of reference and prediction scores",
		averageMetrics,
	)

	// Example: Evaluate changes in LLM responses
	references := []string{
		"The model's performance is critical for the system's success.",
		"Machine learning algorithms can improve efficiency.",
		"The new approach significantly reduces processing time.",
	}
	
	predictions := []string{
		"The model's performance is important for achieving good results.",
		"AI systems enhance productivity through automation.",
		"The novel method substantially decreases computation duration.",
	}

	// Create instances from references and predictions
	instances := make([]eval.Instance, len(references))
	for i := range references {
		instances[i] = eval.Instance{
			Reference:  references[i],
			Prediction: predictions[i],
		}
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run all evaluations
	pairwiseResults, err := pairwiseEval.Run(ctx, instances)
	if err != nil {
		log.Fatalf("Pairwise evaluation failed: %v", err)
	}

	pointwiseResults, err := pointwiseEval.Run(ctx, predictions)
	if err != nil {
		log.Fatalf("Pointwise evaluation failed: %v", err)
	}

	differenceResults, err := differenceEval.Run(ctx, instances)
	if err != nil {
		log.Fatalf("Difference evaluation failed: %v", err)
	}

	ratioResults, err := ratioEval.Run(ctx, instances)
	if err != nil {
		log.Fatalf("Ratio evaluation failed: %v", err)
	}

	absoluteDifferenceResults, err := absoluteDifferenceEval.Run(ctx, instances)
	if err != nil {
		log.Fatalf("Absolute difference evaluation failed: %v", err)
	}

	averageResults, err := averageEval.Run(ctx, instances)
	if err != nil {
		log.Fatalf("Average evaluation failed: %v", err)
	}

	// Print results for all evaluations
	fmt.Println("\nPairwise Evaluation Results:")
	printResults(pairwiseResults)

	fmt.Println("\nPointwise Evaluation Results:")
	printPointwiseResults(pointwiseResults)

	fmt.Println("\nDifference Evaluation Results (prediction - reference):")
	printResults(differenceResults)

	fmt.Println("\nRatio Evaluation Results (prediction / reference):")
	printResults(ratioResults)

	fmt.Println("\nAbsolute Difference Evaluation Results (|prediction - reference|):")
	printResults(absoluteDifferenceResults)

	fmt.Println("\nAverage Evaluation Results ((prediction + reference) / 2):")
	printResults(averageResults)
	
	// Example of running on a single instance
	fmt.Println("\nRunning on a single instance:")
	singleInstance := []eval.Instance{{
		Reference:  "The model's performance is critical for the system's success.",
		Prediction: "The model's performance is important for achieving good results.",
	}}
	
	singlePairwiseResults, err := pairwiseEval.Run(ctx, singleInstance)
	if err != nil {
		log.Fatalf("Single instance pairwise evaluation failed: %v", err)
	}
	
	singlePointwiseResults, err := pointwiseEval.Run(ctx, []string{singleInstance[0].Prediction})
	if err != nil {
		log.Fatalf("Single instance pointwise evaluation failed: %v", err)
	}
	
	singleDifferenceResults, err := differenceEval.Run(ctx, singleInstance)
	if err != nil {
		log.Fatalf("Single instance difference evaluation failed: %v", err)
	}
	
	fmt.Println("\nSingle Instance Pairwise Results:")
	printSingleResult(singlePairwiseResults[0])

	fmt.Println("\nSingle Instance Pointwise Results:")
	printSinglePointwiseResult(singlePointwiseResults[0])

	fmt.Println("\nSingle Instance Difference Results:")
	printSingleResult(singleDifferenceResults[0])

	// Example of using a custom scoring function
	fmt.Println("\nUsing a custom scoring function:")
	customMetric := metrics.KeywordPresence().ToPairwise(func(refScore, predScore float64) float64 {
		// Custom scoring logic: weighted average favoring the prediction
		return 0.3*refScore + 0.7*predScore
	})
	
	customEval := eval.NewPairwiseEvaluation(
		"custom_scoring",
		"Uses a custom weighted scoring function",
		[]eval.PairwiseMetric{customMetric},
	)
	
	customResults, err := customEval.Run(ctx, singleInstance)
	if err != nil {
		log.Fatalf("Custom scoring evaluation failed: %v", err)
	}
	
	fmt.Println("\nCustom Scoring Results (0.3*reference + 0.7*prediction):")
	printSingleResult(customResults[0])
}

// Helper function to print results
func printResults(results []eval.PairwiseResult) {
	for i, result := range results {
		fmt.Printf("\nInstance %d:\n", i+1)
		fmt.Printf("Reference: %q\n", result.Instance.Reference)
		fmt.Printf("Prediction: %q\n", result.Instance.Prediction)
		
		for metricName, score := range result.MetricResults {
			fmt.Printf("%s: %.2f\n", metricName, score)
		}
	}
}

// Helper function to print pointwise results
func printPointwiseResults(results []eval.PointwiseResult) {
	for i, result := range results {
		fmt.Printf("\nPrediction %d: %q\n", i+1, result.Prediction)
		for metricName, score := range result.MetricResults {
			fmt.Printf("%s: %.2f\n", metricName, score)
		}
	}
}

// Helper function to print a single result
func printSingleResult(result eval.PairwiseResult) {
	fmt.Printf("Reference: %q\n", result.Instance.Reference)
	fmt.Printf("Prediction: %q\n", result.Instance.Prediction)
	for metricName, score := range result.MetricResults {
		fmt.Printf("%s: %.2f\n", metricName, score)
	}
}

// Helper function to print a single pointwise result
func printSinglePointwiseResult(result eval.PointwiseResult) {
	fmt.Printf("Prediction: %q\n", result.Prediction)
	for metricName, score := range result.MetricResults {
		fmt.Printf("%s: %.2f\n", metricName, score)
	}
} 