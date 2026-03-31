package main

import "strings"

// fixQuotesProperly removes spaces immediately inside pairs of single-quote characters.
// It works in three passes: find all quote positions, mark bad spaces, rebuild without them.
func fixQuotesProperly(text string) string {
	runes := []rune(text)
	n := len(runes)

	// Pass 1 — record the index of every single-quote character.
	quotePositions := []int{}
	for i := 0; i < n; i++ {
		if runes[i] == '\'' {
			quotePositions = append(quotePositions, i)
		}
	}

	// Pass 2 — pair quotes and mark leading/trailing spaces inside each pair.
	toDelete := make([]bool, n)

	for p := 0; p+1 < len(quotePositions); p += 2 {
		openPos := quotePositions[p]
		closePos := quotePositions[p+1]

		// Mark spaces immediately after the opening quote.
		j := openPos + 1
		for j < closePos && runes[j] == ' ' {
			toDelete[j] = true
			j++
		}

		// Mark spaces immediately before the closing quote.
		k := closePos - 1
		for k > openPos && runes[k] == ' ' {
			toDelete[k] = true
			k--
		}
	}

	// Pass 3 — rebuild the string, skipping every marked position.
	var result strings.Builder
	for i, ch := range runes {
		if !toDelete[i] {
			result.WriteRune(ch)
		}
	}

	return result.String()
}