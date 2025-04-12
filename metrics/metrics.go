package metrics

import (
	"context"
	"strings"
	"unicode"

	"github.com/snpu/eval-go"
)

// StringSimilarity returns a pairwise metric that computes similarity between two strings
func StringSimilarity() eval.PairwiseMetric {
	return eval.NewPairwiseMetric(
		"string_similarity",
		"Computes similarity between two strings",
		func(ctx context.Context, references, predictions []string) ([]float64, error) {
			scores := make([]float64, len(references))
			for i := range references {
				if references[i] == predictions[i] {
					scores[i] = 1.0
				} else if strings.Contains(references[i], predictions[i]) || strings.Contains(predictions[i], references[i]) {
					scores[i] = 0.5
				} else {
					scores[i] = 0.0
				}
			}
			return scores, nil
		},
	)
}

// LengthRatio returns a pairwise metric that computes the ratio of lengths between two strings
func LengthRatio() eval.PairwiseMetric {
	return eval.NewPairwiseMetric(
		"length_ratio",
		"Computes the ratio of lengths between two strings",
		func(ctx context.Context, references, predictions []string) ([]float64, error) {
			scores := make([]float64, len(references))
			for i := range references {
				if len(references[i]) == 0 {
					scores[i] = 0.0
				} else {
					scores[i] = float64(len(predictions[i])) / float64(len(references[i]))
				}
			}
			return scores, nil
		},
	)
}

// WordOverlap returns a pairwise metric that computes Jaccard similarity between words in two strings
func WordOverlap() eval.PairwiseMetric {
	return eval.NewPairwiseMetric(
		"word_overlap",
		"Computes Jaccard similarity between words in two strings",
		func(ctx context.Context, references, predictions []string) ([]float64, error) {
			scores := make([]float64, len(references))
			for i := range references {
				// Split strings into words
				refWords := splitIntoWords(references[i])
				predWords := splitIntoWords(predictions[i])
				
				if len(refWords) == 0 && len(predWords) == 0 {
					scores[i] = 1.0
				} else if len(refWords) == 0 || len(predWords) == 0 {
					scores[i] = 0.0
				} else {
					// Create maps for faster lookup
					refMap := make(map[string]bool)
					predMap := make(map[string]bool)
					
					for _, word := range refWords {
						refMap[word] = true
					}
					
					for _, word := range predWords {
						predMap[word] = true
					}
					
					// Count intersection
					intersection := 0
					for word := range refMap {
						if predMap[word] {
							intersection++
						}
					}
					
					// Count union
					union := len(refMap) + len(predMap) - intersection
					
					scores[i] = float64(intersection) / float64(union)
				}
			}
			return scores, nil
		},
	)
}

// KeywordPresence returns a pointwise metric that checks if text contains specific keywords
func KeywordPresence() eval.PointwiseMetric {
	return eval.NewPointwiseMetric(
		"keyword_presence",
		"Checks if text contains specific keywords",
		func(ctx context.Context, predictions []string) ([]float64, error) {
			// Example keywords - these should be configurable in practice
			keywords := []string{"important", "critical", "significant"}
			
			if len(keywords) == 0 {
				scores := make([]float64, len(predictions))
				return scores, nil
			}
			
			scores := make([]float64, len(predictions))
			for i, prediction := range predictions {
				count := 0
				for _, keyword := range keywords {
					if strings.Contains(strings.ToLower(prediction), strings.ToLower(keyword)) {
						count++
					}
				}
				scores[i] = float64(count) / float64(len(keywords))
			}
			return scores, nil
		},
	)
}

// splitIntoWords splits a string into words, removing punctuation and converting to lowercase
func splitIntoWords(text string) []string {
	// Convert to lowercase
	text = strings.ToLower(text)
	
	// Replace punctuation with spaces
	for _, char := range text {
		if unicode.IsPunct(char) {
			text = strings.ReplaceAll(text, string(char), " ")
		}
	}
	
	// Split by whitespace
	words := strings.Fields(text)
	
	// Remove empty strings
	result := make([]string, 0, len(words))
	for _, word := range words {
		if word != "" {
			result = append(result, word)
		}
	}
	
	return result
}
