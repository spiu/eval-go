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
        "my_evaluation",
        "Evaluates changes in LLM responses",
        pairwiseMetrics,
        pointwiseMetrics,
    )

    // Create instances from references
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
            Reference: references[i],
        }
    }
    
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Run evaluation on the instances and predictions
    results, err := eval.Run(ctx, instances, predictions)
    if err != nil {
        log.Fatal(err)
    }
    
    // Process results
    for i, result := range results {
        fmt.Printf("Instance %d results:\n", i+1)
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

// Create a single instance and prediction
instance := []eval.Instance{{
    Reference: "The model's performance is critical.",
}}
prediction := []string{"The model's performance is important."}

// Run evaluation on a single instance
results, err := eval.Run(ctx, instance, prediction)
if err != nil {
    log.Fatal(err)
}

result := results[0]
for metricName, score := range result.MetricResults {
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