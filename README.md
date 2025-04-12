# eval-go

A declarative evaluation library for assessing changes in LLM tasks. This library allows you to define sets of pairwise or pointwise metrics to evaluate changes in LLM outputs.

## Features

- Declarative metric definition
- Support for both pairwise (comparing reference and prediction) and pointwise (evaluating prediction only) metrics
- Batch processing for efficient evaluation of multiple instances
- Extensible metric system
- Built-in common metrics for LLM evaluation
- Simple and intuitive API
- Context-aware operations for cancellation and timeouts
- Ability to convert pointwise metrics to pairwise ones with custom scoring logic
- Built-in scoring functions for common comparison strategies

## Installation

```bash
go get github.com/snpu/eval-go
```

## Usage

### Basic Usage

```go
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

    // Create pairwise evaluation
    pairwiseEval := eval.NewPairwiseEvaluation(
        "my_evaluation",
        "Evaluates changes in LLM responses",
        pairwiseMetrics,
    )

    // Create pointwise evaluation
    pointwiseEval := eval.NewPointwiseEvaluation(
        "my_pointwise_evaluation",
        "Evaluates the quality of LLM responses",
        pointwiseMetrics,
    )

    // Create evaluations with different conversion strategies
    differenceEval := eval.NewPairwiseEvaluation(
        "my_difference_evaluation",
        "Evaluates the difference between reference and prediction scores",
        convertedPairwiseMetrics,
    )

    ratioEval := eval.NewPairwiseEvaluation(
        "my_ratio_evaluation",
        "Evaluates the ratio between prediction and reference scores",
        alternativePairwiseMetrics,
    )

    absoluteDifferenceEval := eval.NewPairwiseEvaluation(
        "my_absolute_difference_evaluation",
        "Evaluates the absolute difference between reference and prediction scores",
        absoluteDifferenceMetrics,
    )

    // Create instances from references and predictions
    references := []string{
        "The model's performance is critical.",
        "Machine learning improves efficiency.",
    }
    
    predictions := []string{
        "The model's performance is important.",
        "AI systems enhance productivity.",
    }
    
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
    
    // Run pairwise evaluation on the instances
    pairwiseResults, err := pairwiseEval.Run(ctx, instances)
    if err != nil {
        log.Fatal(err)
    }
    
    // Run pointwise evaluation on the predictions
    pointwiseResults, err := pointwiseEval.Run(ctx, predictions)
    if err != nil {
        log.Fatal(err)
    }
    
    // Run difference evaluation on the instances
    differenceResults, err := differenceEval.Run(ctx, instances)
    if err != nil {
        log.Fatal(err)
    }
    
    // Run ratio evaluation on the instances
    ratioResults, err := ratioEval.Run(ctx, instances)
    if err != nil {
        log.Fatal(err)
    }
    
    // Run absolute difference evaluation on the instances
    absoluteDifferenceResults, err := absoluteDifferenceEval.Run(ctx, instances)
    if err != nil {
        log.Fatal(err)
    }
    
    // Process results
    for i, result := range pairwiseResults {
        fmt.Printf("Instance %d pairwise results:\n", i+1)
        for metricName, score := range result.MetricResults {
            fmt.Printf("%s: %.2f\n", metricName, score)
        }
    }
    
    for i, result := range pointwiseResults {
        fmt.Printf("Prediction %d pointwise results:\n", i+1)
        for metricName, score := range result.MetricResults {
            fmt.Printf("%s: %.2f\n", metricName, score)
        }
    }
    
    for i, result := range differenceResults {
        fmt.Printf("Instance %d difference results:\n", i+1)
        for metricName, score := range result.MetricResults {
            fmt.Printf("%s: %.2f\n", metricName, score)
        }
    }
    
    for i, result := range ratioResults {
        fmt.Printf("Instance %d ratio results:\n", i+1)
        for metricName, score := range result.MetricResults {
            fmt.Printf("%s: %.2f\n", metricName, score)
        }
    }
    
    for i, result := range absoluteDifferenceResults {
        fmt.Printf("Instance %d absolute difference results:\n", i+1)
        for metricName, score := range result.MetricResults {
            fmt.Printf("%s: %.2f\n", metricName, score)
        }
    }
}
```

### Single Instance Evaluation

```go
// Create a context
ctx := context.Background()

// Create a single instance
instance := []eval.Instance{{
    Reference:  "The model's performance is critical.",
    Prediction: "The model's performance is important.",
}}

// Run pairwise evaluation on a single instance
pairwiseResults, err := pairwiseEval.Run(ctx, instance)
if err != nil {
    log.Fatal(err)
}

pairwiseResult := pairwiseResults[0]
for metricName, score := range pairwiseResult.MetricResults {
    fmt.Printf("%s: %.2f\n", metricName, score)
}

// Run pointwise evaluation on a single prediction
pointwiseResults, err := pointwiseEval.Run(ctx, []string{instance[0].Prediction})
if err != nil {
    log.Fatal(err)
}

pointwiseResult := pointwiseResults[0]
for metricName, score := range pointwiseResult.MetricResults {
    fmt.Printf("%s: %.2f\n", metricName, score)
}
```

### Creating Custom Metrics

You can create custom metrics using the `NewPairwiseMetric` and `NewPointwiseMetric` functions:

```go
// Create a pairwise metric
pairwiseMetric := eval.NewPairwiseMetric(
    "my_pairwise_metric",
    "Description of my pairwise metric",
    func(ctx context.Context, references, predictions []string) ([]float64, error) {
        // Your metric computation logic here
        scores := make([]float64, len(references))
        for i := range references {
            // Compute score for each reference-prediction pair
            scores[i] = computeScore(references[i], predictions[i])
        }
        return scores, nil
    },
)

// Create a pointwise metric
pointwiseMetric := eval.NewPointwiseMetric(
    "my_pointwise_metric",
    "Description of my pointwise metric",
    func(ctx context.Context, predictions []string) ([]float64, error) {
        // Your metric computation logic here
        scores := make([]float64, len(predictions))
        for i, prediction := range predictions {
            // Compute score for each prediction
            scores[i] = computeScore(prediction)
        }
        return scores, nil
    },
)
```

### Converting Pointwise to Pairwise

You can convert a pointwise metric to a pairwise one using the `ToPairwise` method with a custom scoring function:

```go
// Create a pointwise metric
pointwiseMetric := eval.NewPointwiseMetric(
    "my_pointwise_metric",
    "Evaluates the quality of text",
    func(ctx context.Context, predictions []string) ([]float64, error) {
        // Your metric computation logic here
        scores := make([]float64, len(predictions))
        for i, prediction := range predictions {
            // Compute score for each prediction
            scores[i] = computeScore(prediction)
        }
        return scores, nil
    },
)

// Convert to pairwise metric with a built-in scoring function
pairwiseMetric := pointwiseMetric.ToPairwise(metrics.DifferenceScore)

// Or use a custom scoring function
customPairwiseMetric := pointwiseMetric.ToPairwise(func(refScore, predScore float64) float64 {
    // Custom scoring logic: weighted average favoring the prediction
    return 0.3*refScore + 0.7*predScore
})

// Now you can use the pairwise metric in a pairwise evaluation
pairwiseEval := eval.NewPairwiseEvaluation(
    "my_evaluation",
    "Evaluates changes in LLM responses",
    []eval.PairwiseMetric{pairwiseMetric},
)
```

The `ToPairwise` method allows you to define how to calculate the final score between the reference and prediction scores. This is useful when you want to evaluate predictions both in isolation and in comparison to references, or when you want to use the same metrics in both contexts with different scoring strategies.

### Built-in Scoring Functions

The library provides several built-in scoring functions for converting pointwise metrics to pairwise ones:

- `DifferenceScore`: Calculates the difference between prediction and reference scores (`prediction - reference`)
- `RatioScore`: Calculates the ratio between prediction and reference scores (`prediction / reference`)
- `AbsoluteDifferenceScore`: Calculates the absolute difference between prediction and reference scores (`|prediction - reference|`)
- `MaxScore`: Returns the maximum of the reference and prediction scores
- `MinScore`: Returns the minimum of the reference and prediction scores
- `AverageScore`: Calculates the average of the reference and prediction scores

You can also define your own custom scoring functions by implementing the `PairwiseScoreFunc` type:

```go
type PairwiseScoreFunc func(referenceScore, predictionScore float64) float64
```

## Built-in Metrics

The library includes several common metrics:

### Pairwise Metrics
- `StringSimilarity()`: Computes similarity between two strings
- `LengthRatio()`: Computes the ratio of lengths between two strings
- `WordOverlap()`: Computes Jaccard similarity between words in two strings

### Pointwise Metrics
- `KeywordPresence()`: Checks if text contains specific keywords

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License 