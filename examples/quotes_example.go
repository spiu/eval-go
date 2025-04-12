package main

import (
	"context"
	"fmt"
	"log"
	"time"

	eval "github.com/snpu/eval-go"
	metrics "github.com/snpu/eval-go/metrics"
)

func main() {
	// Create pointwise metrics for Reddit quotes
	pointwiseMetrics := []eval.PointwiseMetric{
		metrics.QuotesCount(),
		metrics.QuotesRatio(),
		metrics.QuotesPresence(),
		metrics.QuotesSize(),
		metrics.ShortQuotesCount(4), // Count quotes with fewer than 4 words
		metrics.ExternalLinksCount(),
		metrics.QuoteDiversity(),
		metrics.PostDiversity(),
		metrics.SubredditDiversity(),
	}

	// Create pointwise evaluation
	pointwiseEval := eval.NewPointwiseEvaluation(
		"quotes_evaluation",
		"Evaluates the presence of Reddit quotes in summaries",
		pointwiseMetrics,
	)

	// Example summaries with Reddit quotes
	summaries := []string{
		`According to the discussion, many users are concerned about privacy. ["I don't trust these companies with my data"](https://www.reddit.com/r/privacy/comments/abc123/user_comment) was a common sentiment. Another user noted ["We need better regulations"](https://www.reddit.com/r/privacy/comments/def456/another_comment). For more information, see [this article](https://example.com/privacy-article).`,
		
		`The community is divided on this issue. Some believe it's a step forward, while others see it as problematic. ["This will only benefit the wealthy"](https://www.reddit.com/r/politics/comments/ghi789/political_comment) was one perspective shared. ["I disagree completely"](https://www.reddit.com/r/politics/comments/ghi789/another_comment) was another view.`,
		
		`This summary has a short quote: ["I agree"](https://www.reddit.com/r/example/comments/123456/comment) and a longer one: ["This is a much longer quote that exceeds the threshold for being considered short"](https://www.reddit.com/r/example/comments/789012/comment). It also references [an external site](https://example.org/some-article).`,
		
		`This summary has quotes from multiple subreddits: ["I love this feature"](https://www.reddit.com/r/technology/comments/abc123/comment), ["This is amazing"](https://www.reddit.com/r/programming/comments/def456/comment), and ["I can't believe this works"](https://www.reddit.com/r/coding/comments/ghi789/comment).`,
		
		`This summary doesn't contain any Reddit quotes, but it does link to [an external article](https://example.com/article) and [another external article](https://example.org/another-article).`,
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run pointwise evaluation on the summaries
	pointwiseResults, err := pointwiseEval.Run(ctx, summaries)
	if err != nil {
		log.Fatalf("Pointwise evaluation failed: %v", err)
	}

	// Print results
	fmt.Println("Reddit Quotes Metrics Results:")
	for i, result := range pointwiseResults {
		fmt.Printf("\nSummary %d: %q\n", i+1, result.Prediction)
		for metricName, score := range result.MetricResults {
			fmt.Printf("%s: %.2f\n", metricName, score)
		}
	}

	// Example of converting to pairwise metrics
	fmt.Println("\nConverting to pairwise metrics:")
	
	// Convert pointwise metrics to pairwise with different scoring strategies
	convertedPairwiseMetrics := make([]eval.PairwiseMetric, len(pointwiseMetrics))
	for i, metric := range pointwiseMetrics {
		// Use the difference between reference and prediction scores
		convertedPairwiseMetrics[i] = metric.ToPairwise(eval.DifferenceScore)
	}

	// Create pairwise evaluation
	pairwiseEval := eval.NewPairwiseEvaluation(
		"quotes_comparison",
		"Compares the number of Reddit quotes between reference and prediction",
		convertedPairwiseMetrics,
	)

	// Example instances with references and predictions
	instances := []eval.Instance{
		{
			Reference:  `The discussion highlighted privacy concerns. ["Privacy is a fundamental right"](https://www.reddit.com/r/privacy/comments/xyz789/reference_comment).`,
			Prediction: `According to the discussion, many users are concerned about privacy. ["I don't trust these companies with my data"](https://www.reddit.com/r/privacy/comments/abc123/user_comment) was a common sentiment. Another user noted ["We need better regulations"](https://www.reddit.com/r/privacy/comments/def456/another_comment).`,
		},
		{
			Reference:  `The community is divided on this issue. ["This will only benefit the wealthy"](https://www.reddit.com/r/politics/comments/ghi789/political_comment) was one perspective shared.`,
			Prediction: `The community is divided on this issue. Some believe it's a step forward, while others see it as problematic. ["This will only benefit the wealthy"](https://www.reddit.com/r/politics/comments/ghi789/political_comment) was one perspective shared. ["I disagree completely"](https://www.reddit.com/r/politics/comments/ghi789/another_comment) was another view.`,
		},
		{
			Reference:  `This reference has one user comment. ["I agree with this approach"](https://www.reddit.com/r/example/comments/123456/comment).`,
			Prediction: `This summary has a short quote: ["I agree"](https://www.reddit.com/r/example/comments/123456/comment) and a longer one: ["This is a much longer quote that exceeds the threshold for being considered short"](https://www.reddit.com/r/example/comments/789012/comment). It also references [an external site](https://example.org/some-article).`,
		},
		{
			Reference:  `This reference has one user comment. ["I agree with this approach"](https://www.reddit.com/r/example/comments/123456/comment).`,
			Prediction: `This summary has quotes from multiple subreddits: ["I love this feature"](https://www.reddit.com/r/technology/comments/abc123/comment), ["This is amazing"](https://www.reddit.com/r/programming/comments/def456/comment), and ["I can't believe this works"](https://www.reddit.com/r/coding/comments/ghi789/comment).`,
		},
		{
			Reference:  `This reference has one user comment. ["I agree with this approach"](https://www.reddit.com/r/example/comments/123456/comment).`,
			Prediction: `This summary doesn't contain any Reddit quotes, but it does link to [an external article](https://example.com/article) and [another external article](https://example.org/another-article).`,
		},
	}

	// Run pairwise evaluation on the instances
	pairwiseResults, err := pairwiseEval.Run(ctx, instances)
	if err != nil {
		log.Fatalf("Pairwise evaluation failed: %v", err)
	}

	// Print results
	fmt.Println("\nReddit Quotes Pairwise Metrics Results:")
	for i, result := range pairwiseResults {
		fmt.Printf("\nInstance %d:\n", i+1)
		fmt.Printf("Reference: %q\n", result.Instance.Reference)
		fmt.Printf("Prediction: %q\n", result.Instance.Prediction)
		for metricName, score := range result.MetricResults {
			fmt.Printf("%s: %.2f\n", metricName, score)
		}
	}
} 