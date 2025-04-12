package metrics

import (
	"context"
	"regexp"
	"strings"

	eval "github.com/snpu/eval-go"
)

// redditQuoteRegex is a regular expression to match Reddit user quotes in markdown format
// Format: ["some user excerpt"](https://www.reddit.com/...)
var redditQuoteRegex = regexp.MustCompile(`\[(.*?)\]\((https://www\.reddit\.com/.*?)\)`)

// externalLinkRegex is a regular expression to match markdown links that point to external sites
// Format: ["some text"](https://example.com/...)
var externalLinkRegex = regexp.MustCompile(`\[(.*?)\]\((https://[^)]+)\)`)

// postIdRegex is a regular expression to extract post IDs from Reddit URLs
// Format: https://www.reddit.com/r/subreddit/comments/postId/...
var postIdRegex = regexp.MustCompile(`https://www\.reddit\.com/r/.*?/comments/([a-zA-Z0-9]+)/`)

// subredditRegex is a regular expression to extract subreddit names from Reddit URLs
// Format: https://www.reddit.com/r/subreddit/...
var subredditRegex = regexp.MustCompile(`https://www\.reddit\.com/r/([a-zA-Z0-9_]+)/`)

// QuotesCount returns a pointwise metric that counts the number of Reddit user quotes in markdown format
func QuotesCount() eval.PointwiseMetric {
	return eval.NewPointwiseMetric(
		"quotes_count",
		"Counts the number of Reddit user quotes in markdown format",
		func(ctx context.Context, predictions []string) ([]float64, error) {
			scores := make([]float64, len(predictions))
			
			for i, prediction := range predictions {
				matches := redditQuoteRegex.FindAllString(prediction, -1)
				scores[i] = float64(len(matches))
			}
			
			return scores, nil
		},
	)
}

// QuotesRatio returns a pointwise metric that calculates the ratio of Reddit user quotes to the total number of words
func QuotesRatio() eval.PointwiseMetric {
	return eval.NewPointwiseMetric(
		"quotes_ratio",
		"Calculates the ratio of Reddit user quotes to the total number of words",
		func(ctx context.Context, predictions []string) ([]float64, error) {
			scores := make([]float64, len(predictions))
			
			for i, prediction := range predictions {
				matches := redditQuoteRegex.FindAllString(prediction, -1)
				words := strings.Fields(prediction)
				if len(words) == 0 {
					scores[i] = 0.0
				} else {
					scores[i] = float64(len(matches)) / float64(len(words))
				}
			}
			
			return scores, nil
		},
	)
}

// QuotesPresence returns a pointwise metric that checks if there is at least one Reddit user quote in the text
func QuotesPresence() eval.PointwiseMetric {
	return eval.NewPointwiseMetric(
		"quotes_presence",
		"Checks if there is at least one Reddit user quote in the text",
		func(ctx context.Context, predictions []string) ([]float64, error) {
			scores := make([]float64, len(predictions))
			
			for i, prediction := range predictions {
				if redditQuoteRegex.MatchString(prediction) {
					scores[i] = 1.0
				} else {
					scores[i] = 0.0
				}
			}
			
			return scores, nil
		},
	)
}

// QuotesSize returns a pointwise metric that calculates the total size of quoted text in characters
func QuotesSize() eval.PointwiseMetric {
	return eval.NewPointwiseMetric(
		"quotes_size",
		"Calculates the total size of quoted text in characters",
		func(ctx context.Context, predictions []string) ([]float64, error) {
			scores := make([]float64, len(predictions))
			
			for i, prediction := range predictions {
				matches := redditQuoteRegex.FindAllStringSubmatch(prediction, -1)
				totalSize := 0
				for _, match := range matches {
					if len(match) > 1 {
						totalSize += len(match[1])
					}
				}
				scores[i] = float64(totalSize)
			}
			
			return scores, nil
		},
	)
}

// ShortQuotesCount returns a pointwise metric that counts the number of quotes with fewer words than the specified threshold
func ShortQuotesCount(threshold int) eval.PointwiseMetric {
	return eval.NewPointwiseMetric(
		"short_quotes_count",
		"Counts the number of quotes with fewer words than the specified threshold",
		func(ctx context.Context, predictions []string) ([]float64, error) {
			scores := make([]float64, len(predictions))
			
			for i, prediction := range predictions {
				matches := redditQuoteRegex.FindAllStringSubmatch(prediction, -1)
				shortCount := 0
				for _, match := range matches {
					if len(match) > 1 {
						words := strings.Fields(match[1])
						if len(words) < threshold {
							shortCount++
						}
					}
				}
				scores[i] = float64(shortCount)
			}
			
			return scores, nil
		},
	)
}

// ExternalLinksCount returns a pointwise metric that counts the number of external site references in markdown format
func ExternalLinksCount() eval.PointwiseMetric {
	return eval.NewPointwiseMetric(
		"external_links_count",
		"Counts the number of external site references in markdown format",
		func(ctx context.Context, predictions []string) ([]float64, error) {
			scores := make([]float64, len(predictions))
			
			for i, prediction := range predictions {
				allLinks := regexp.MustCompile(`\[(.*?)\]\((https://[^)]+)\)`).FindAllString(prediction, -1)
				externalCount := 0
				for _, link := range allLinks {
					if !strings.Contains(link, "reddit.com") {
						externalCount++
					}
				}
				scores[i] = float64(externalCount)
			}
			
			return scores, nil
		},
	)
}

// QuoteDiversity returns a pointwise metric that counts the number of unique Reddit links referenced in the text
func QuoteDiversity() eval.PointwiseMetric {
	return eval.NewPointwiseMetric(
		"quote_diversity",
		"Counts the number of unique Reddit links referenced in the text",
		func(ctx context.Context, predictions []string) ([]float64, error) {
			scores := make([]float64, len(predictions))
			
			for i, prediction := range predictions {
				matches := redditQuoteRegex.FindAllString(prediction, -1)
				uniqueLinks := make(map[string]bool)
				for _, match := range matches {
					uniqueLinks[match] = true
				}
				scores[i] = float64(len(uniqueLinks))
			}
			
			return scores, nil
		},
	)
}

// PostDiversity returns a pointwise metric that counts the number of unique Reddit post IDs referenced in the text
func PostDiversity() eval.PointwiseMetric {
	return eval.NewPointwiseMetric(
		"post_diversity",
		"Counts the number of unique Reddit post IDs referenced in the text",
		func(ctx context.Context, predictions []string) ([]float64, error) {
			scores := make([]float64, len(predictions))
			
			for i, prediction := range predictions {
				matches := redditQuoteRegex.FindAllString(prediction, -1)
				uniquePostIds := make(map[string]bool)
				for _, match := range matches {
					postIds := postIdRegex.FindAllStringSubmatch(match, -1)
					for _, postId := range postIds {
						if len(postId) > 1 {
							uniquePostIds[postId[1]] = true
						}
					}
				}
				scores[i] = float64(len(uniquePostIds))
			}
			
			return scores, nil
		},
	)
}

// SubredditDiversity returns a pointwise metric that counts the number of unique subreddits referenced in the text
func SubredditDiversity() eval.PointwiseMetric {
	return eval.NewPointwiseMetric(
		"subreddit_diversity",
		"Counts the number of unique subreddits referenced in the text",
		func(ctx context.Context, predictions []string) ([]float64, error) {
			scores := make([]float64, len(predictions))
			
			for i, prediction := range predictions {
				matches := redditQuoteRegex.FindAllString(prediction, -1)
				uniqueSubreddits := make(map[string]bool)
				for _, match := range matches {
					subreddits := subredditRegex.FindAllStringSubmatch(match, -1)
					for _, subreddit := range subreddits {
						if len(subreddit) > 1 {
							uniqueSubreddits[subreddit[1]] = true
						}
					}
				}
				scores[i] = float64(len(uniqueSubreddits))
			}
			
			return scores, nil
		},
	)
} 