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

	// Run the pairwise evaluation
	pairwiseResults, err := pairwiseEval.Run(ctx, instances)
	if err != nil {
		log.Fatalf("Pairwise evaluation failed: %v", err)
	}

	// Run the pointwise evaluation
	pointwiseResults, err := pointwiseEval.Run(ctx, predictions)
	if err != nil {
		log.Fatalf("Pointwise evaluation failed: %v", err)
	}

	// Print pairwise results
	fmt.Println("\nPairwise Evaluation Results:")
	for i, result := range pairwiseResults {
		fmt.Printf("\nInstance %d:\n", i+1)
		fmt.Printf("Reference: %q\n", result.Instance.Reference)
		fmt.Printf("Prediction: %q\n", result.Instance.Prediction)
		
		for metricName, score := range result.MetricResults {
			fmt.Printf("%s: %.2f\n", metricName, score)
		}
	}

	// Print pointwise results
	fmt.Println("\nPointwise Evaluation Results:")
	for i, result := range pointwiseResults {
		fmt.Printf("\nPrediction %d: %q\n", i+1, result.Prediction)
		for metricName, score := range result.MetricResults {
			fmt.Printf("%s: %.2f\n", metricName, score)
		}
	}
	
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
	
	fmt.Println("\nSingle Instance Pairwise Results:")
	pairwiseResult := singlePairwiseResults[0]
	fmt.Printf("Reference: %q\n", pairwiseResult.Instance.Reference)
	fmt.Printf("Prediction: %q\n", pairwiseResult.Instance.Prediction)
	for metricName, score := range pairwiseResult.MetricResults {
		fmt.Printf("%s: %.2f\n", metricName, score)
	}

	fmt.Println("\nSingle Instance Pointwise Results:")
	pointwiseResult := singlePointwiseResults[0]
	fmt.Printf("Prediction: %q\n", pointwiseResult.Prediction)
	for metricName, score := range pointwiseResult.MetricResults {
		fmt.Printf("%s: %.2f\n", metricName, score)
	}
} 