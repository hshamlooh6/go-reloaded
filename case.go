package main

import (
	"strconv"
	"strings"
)

// applyCaseTransformations processes (up), (low), (cap) markers with optional counts.
// It first repairs any split markers, then applies case changes to the preceding words.
func applyCaseTransformations(words []string) []string {
	words = rejoinSplitMarkers(words)
	result := []string{}

	i := 0
	for i < len(words) {
		word := words[i]

		caseType, count := parseCaseMarker(word)

		if caseType != "" {
			if len(result) >= count {
				for j := len(result) - count; j < len(result); j++ {
					result[j] = applyCase(result[j], caseType)
				}
			}
		} else {
			result = append(result, word)
		}

		i++
	}

	return result
}

// rejoinSplitMarkers repairs markers like (up, 3) that were split into two tokens
// by strings.Fields when the input contained a space after the comma.
func rejoinSplitMarkers(words []string) []string {
	result := []string{}
	i := 0
	for i < len(words) {
		word := words[i]
		lower := strings.ToLower(word)
		if (strings.HasPrefix(lower, "(up,") ||
			strings.HasPrefix(lower, "(low,") ||
			strings.HasPrefix(lower, "(cap,")) &&
			strings.HasSuffix(word, ",") &&
			i+1 < len(words) {
			joined := word + words[i+1]
			result = append(result, joined)
			i += 2
		} else {
			result = append(result, word)
			i++
		}
	}
	return result
}

// parseCaseMarker checks whether a token is a case marker and returns
// the case type ("up", "low", "cap") and the number of words to affect.
// Returns ("", 0) if the token is not a marker.
func parseCaseMarker(word string) (string, int) {
	clean := strings.ReplaceAll(word, " ", "")
	clean = strings.ToLower(clean)

	if clean == "(up)" {
		return "up", 1
	}
	if clean == "(low)" {
		return "low", 1
	}
	if clean == "(cap)" {
		return "cap", 1
	}

	for _, caseType := range []string{"up", "low", "cap"} {
		prefix := "(" + caseType + ","

		if strings.HasPrefix(clean, prefix) {
			inner := clean[len(prefix):]
			inner = strings.TrimRight(inner, ")")
			inner = strings.TrimSpace(inner)

			number, err := strconv.Atoi(inner)
			if err == nil && number > 0 {
				return caseType, number
			}
		}
	}

	return "", 0
}

// applyCase transforms a single word to uppercase, lowercase, or capitalised.
func applyCase(word string, caseType string) string {
	if caseType == "up" {
		return strings.ToUpper(word)
	}
	if caseType == "low" {
		return strings.ToLower(word)
	}
	if caseType == "cap" {
		if len(word) == 0 {
			return word
		}
		return strings.ToUpper(string(word[0])) + word[1:]
	}
	return word
}
