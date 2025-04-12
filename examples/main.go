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
	// Create metrics
	pairwiseMetrics := []eval.PairwiseMetric{
		metrics.StringSimilarity(),
		metrics.LengthRatio(),
		metrics.WordOverlap(),
	}
	
	pointwiseMetrics := []eval.PointwiseMetric{
		metrics.KeywordPresence(),
	}

	// Create a new evaluation with metrics
	eval := eval.NewEvaluation(
		"llm_response_evaluation",
		"Evaluates changes in LLM responses",
		pairwiseMetrics,
		pointwiseMetrics,
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

	// Create instances from references
	instances := make([]eval.Instance, len(references))
	for i := range references {
		instances[i] = eval.Instance{
			Reference: references[i],
		}
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run the evaluation on the instances and predictions
	results, err := eval.Run(ctx, instances, predictions)
	if err != nil {
		log.Fatalf("Evaluation failed: %v", err)
	}

	// Print results for each instance
	for i, result := range results {
		fmt.Printf("\nResults for Instance %d:\n", i+1)
		fmt.Printf("Reference: %q\n", references[i])
		fmt.Printf("Prediction: %q\n", predictions[i])
		
		for metricName, score := range result.MetricResults {
			fmt.Printf("%s: %.2f\n", metricName, score)
		}
	}
	
	// Example of running on a single instance
	fmt.Println("\nRunning on a single instance:")
	singleInstance := []eval.Instance{{
		Reference: "The model's performance is critical for the system's success.",
	}}
	singlePrediction := []string{"The model's performance is important for achieving good results."}
	
	singleResults, err := eval.Run(ctx, singleInstance, singlePrediction)
	if err != nil {
		log.Fatalf("Single instance evaluation failed: %v", err)
	}
	
	singleResult := singleResults[0]
	fmt.Printf("Reference: %q\n", singleInstance[0].Reference)
	fmt.Printf("Prediction: %q\n", singlePrediction[0])
	for metricName, score := range singleResult.MetricResults {
		fmt.Printf("%s: %.2f\n", metricName, score)
	}
} 