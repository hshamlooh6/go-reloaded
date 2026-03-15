package main

import (
	"strconv"
	"strings"
	"unicode"
)

// -------------------------------------------------------
// MAIN FUNCTION: runs all transformations one by one
// -------------------------------------------------------

// ApplyAllTransformations takes the raw text and applies every rule in order
func ApplyAllTransformations(text string) string {
	// Split the text into a slice of words (tokens)
	words := strings.Fields(text)
	
	// Apply each transformation step by step
	words = applyHexAndBin(words)
	words = applyCaseTransformations(words)
	words = applyAToAn(words)

	// Join the words back into a single string
	result := strings.Join(words, " ")

	// Fix punctuation and quotes (easier to work on the full string)
	result = fixPunctuation(result)
	result = fixQuotes(result)

	return result
}

// -------------------------------------------------------
// STEP 1: Convert hex and bin numbers to decimal
// -------------------------------------------------------

// applyHexAndBin looks for (hex) and (bin) markers and converts the word before them
func applyHexAndBin(words []string) []string {
	// We will build a new list of words as the result
	result := []string{}

	// Go through every word one by one
	i := 0
	for i < len(words) {
		word := words[i]

		// Check if the current word is a (hex) or (bin) marker
		if word == "(hex)" || word == "(bin)" {
			// We need a word BEFORE this marker to convert
			if len(result) > 0 {
				// Take the last word we added to result
				lastWord := result[len(result)-1]

				var converted string
				var err error

				if word == "(hex)" {
					// Convert hex string to decimal number
					converted, err = hexToDecimal(lastWord)
				} else {
					// Convert binary string to decimal number
					converted, err = binToDecimal(lastWord)
				}

				// If conversion worked, replace the last word with the decimal
				if err == nil {
					result[len(result)-1] = converted
				}
				// If conversion failed, just skip the marker (leave text as is)
			}
			// Skip the marker word itself (don't add it to result)
		} else {
			// Normal word, just add it to result
			result = append(result, word)
		}

		i++
	}

	return result
}

// hexToDecimal converts a hex string like "1E" to a decimal string like "30"
func hexToDecimal(hex string) (string, error) {
	// strconv.ParseInt can parse hex when we pass base 16
	// 64 means we want a 64-bit integer
	number, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return "", err
	}
	// Convert the number back to a string (base 10 = decimal)
	return strconv.FormatInt(number, 10), nil
}

// binToDecimal converts a binary string like "10" to a decimal string like "2"
func binToDecimal(bin string) (string, error) {
	// strconv.ParseInt can parse binary when we pass base 2
	number, err := strconv.ParseInt(bin, 2, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(number, 10), nil
}

// -------------------------------------------------------
// STEP 2: Apply (up), (low), (cap) transformations
// -------------------------------------------------------

// rejoinSplitMarkers fixes cases where "(up, 2)" got split into "(up," and "2)"
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

// applyCaseTransformations looks for (up), (low), (cap) markers
// and changes the case of the word(s) before them
func applyCaseTransformations(words []string) []string {
	words = rejoinSplitMarkers(words)
	result := []string{}

	i := 0
	for i < len(words) {
		word := words[i]

		// Check if this word is a case marker (possibly with a number like "(up, 3)")
		caseType, count := parseCaseMarker(word)

		if caseType != "" {
			// We need to change the last `count` words in result
			// First, make sure we have enough words
			if len(result) >= count {
				// Go back `count` positions and change each word
				for j := len(result) - count; j < len(result); j++ {
					result[j] = applyCase(result[j], caseType)
				}
			}
			// Don't add the marker itself to result
		} else {
			// Normal word, add it
			result = append(result, word)
		}

		i++
	}

	return result
}

// parseCaseMarker checks if a word is like "(up)", "(low)", "(cap)",
// "(up, 2)", "(low, 3)", etc.
// Returns the case type ("up", "low", "cap") and the count (default 1)
func parseCaseMarker(word string) (string, int) {
	// Remove spaces inside the word for easier checking
	clean := strings.ReplaceAll(word, " ", "")
	clean = strings.ToLower(clean)

	// Simple cases without a number
	if clean == "(up)" {
		return "up", 1
	}
	if clean == "(low)" {
		return "low", 1
	}
	if clean == "(cap)" {
		return "cap", 1
	}

	// Cases with a number: (up, 2), (low, 3), (cap, 1)
	// We check if the word starts with "(up," or "(low," or "(cap,"
	// Note: clean already has spaces removed, so "(up, 2)" becomes "(up,2)"
	for _, caseType := range []string{"up", "low", "cap"} {
		prefix := "(" + caseType + ","

		if strings.HasPrefix(clean, prefix) {
			// Extract the number part
			// Example: "(up,2)" -> we want "2"
			inner := clean[len(prefix):]          // "2)"
			inner = strings.TrimRight(inner, ")") // "2"
			inner = strings.TrimSpace(inner)      // "2" (remove any spaces)

			number, err := strconv.Atoi(inner)
			if err == nil && number > 0 {
				return caseType, number
			}
		}
	}

	// Also handle the original word (with spaces) in case clean removed too much
	// For example the token might be a combined "(up," followed by "2)"
	// This handles tokens like "(up," that appear as a single word due to splitting

	// Not a case marker
	return "", 0
}

// applyCase converts a word to upper, lower, or capitalized form
func applyCase(word string, caseType string) string {
	if caseType == "up" {
		return strings.ToUpper(word)
	}
	if caseType == "low" {
		return strings.ToLower(word)
	}
	if caseType == "cap" {
		// Capitalize = first letter upper, rest lower
		if len(word) == 0 {
			return word
		}
		// Make the first character uppercase, keep the rest as is
		return strings.ToUpper(string(word[0])) + word[1:]
	}
	return word
}

// -------------------------------------------------------
// STEP 3: Fix "a" -> "an" before vowels and 'h'
// -------------------------------------------------------

// applyAToAn checks each word: if it's "a" or "A" and the next word
// starts with a vowel or 'h', change it to "an" or "An"
func applyAToAn(words []string) []string {
	// Vowels and 'h' trigger the a -> an rule
	vowelsAndH := "aeiouAEIOUhH"

	for i := 0; i < len(words)-1; i++ {
		currentWord := words[i]
		nextWord := words[i+1]

		// Check if current word is "a" or "A"
		if currentWord == "a" || currentWord == "A" {
			// Check if next word starts with a vowel or 'h'
			if len(nextWord) > 0 && strings.ContainsRune(vowelsAndH, rune(nextWord[0])) {
				// Replace "a" with "an", keeping the same case
				if currentWord == "A" {
					words[i] = "An"
				} else {
					words[i] = "an"
				}
			}
		}
	}

	return words
}

// -------------------------------------------------------
// STEP 4: Fix punctuation spacing
// -------------------------------------------------------

// fixPunctuation makes sure punctuation like . , ! ? : ;
// is attached to the word before it, with a space after
func fixPunctuation(text string) string {
	// We will go through the text character by character
	// and build the result in a strings.Builder (efficient way to build strings)
	var result strings.Builder

	// Convert the text to a slice of runes (handles unicode characters safely)
	chars := []rune(text)
	length := len(chars)

	i := 0
	for i < length {
		ch := chars[i]

		// Check if this character is a punctuation we care about
		if isPunctuation(ch) {
			// Collect all consecutive punctuation characters (like "..." or "!?")
			punctGroup := []rune{}
			for i < length && isPunctuation(chars[i]) {
				punctGroup = append(punctGroup, chars[i])
				i++
			}

			// Remove any trailing space from what we've built so far
			// (so the punctuation hugs the previous word)
			built := result.String()
			built = strings.TrimRight(built, " ")
			result.Reset()
			result.WriteString(built)

			// Write the punctuation group
			result.WriteString(string(punctGroup))

			// If there are more characters after, add a space
			if i < length && chars[i] != ' ' {
				result.WriteRune(' ')
			}
		} else {
			// Normal character, just add it
			result.WriteRune(ch)
			i++
		}
	}

	// Clean up: remove extra spaces
	return cleanSpaces(result.String())
}

// isPunctuation returns true if the character is one of our special punctuation marks
func isPunctuation(ch rune) bool {
	return ch == '.' || ch == ',' || ch == '!' || ch == '?' || ch == ':' || ch == ';'
}

// cleanSpaces removes extra spaces (multiple spaces become one, trim edges)
func cleanSpaces(text string) string {
	// Split by whitespace and rejoin with single spaces
	words := strings.Fields(text)
	return strings.Join(words, " ")
}

// -------------------------------------------------------
// STEP 5: Fix single quotes
// -------------------------------------------------------

// fixQuotes finds pairs of ' marks and removes the spaces
// between the quote and the word next to it
// Example: " ' awesome ' " -> " 'awesome' "
func fixQuotes(text string) string {
	// We look for pairs of quotes: ' ... '
	// Strategy: find the first ', then find the matching second '
	// and remove spaces right after the first ' and right before the second '

	result := []rune(text)

	for {
		// Find the first single quote
		firstPos := -1
		for i, ch := range result {
			if ch == '\'' {
				firstPos = i
				break
			}
		}

		// If no quote found, we are done
		if firstPos == -1 {
			break
		}

		// Find the second single quote (after the first one)
		secondPos := -1
		for i := firstPos + 1; i < len(result); i++ {
			if result[i] == '\'' {
				secondPos = i
				break
			}
		}

		// If no matching second quote, we are done
		if secondPos == -1 {
			break
		}

		// Remove spaces right after the first quote
		// Example: ' awesome' -> 'awesome
		for firstPos+1 < len(result) && result[firstPos+1] == ' ' {
			result = removeCharAt(result, firstPos+1)
			// secondPos shifted left by 1 since we removed a character
			secondPos--
		}

		// Remove spaces right before the second quote
		// Example: awesome ' -> awesome'
		for secondPos > 0 && result[secondPos-1] == ' ' {
			result = removeCharAt(result, secondPos-1)
			secondPos--
		}

		// Now we need to move past these quotes to avoid infinite loop
		// We replace them temporarily with a placeholder to not find them again
		// Actually, let's just search after secondPos next iteration
		// We do this by converting already-processed quotes to a rare marker
		// Simpler: mark them as processed by replacing with a different unicode char
		// then replace back at the end
		result[firstPos] = '\''  // keep as is
		result[secondPos] = '\'' // keep as is

		// To avoid re-processing, we break the loop and do it in passes
		// by rebuilding and scanning after the secondPos
		// Simple approach: process the string in one pass using index tracking
		break
	}

	// Use a cleaner approach: process all quote pairs in one pass
	return fixQuotesProperly(text)
}

// fixQuotesProperly is the clean version that handles all quote pairs correctly
func fixQuotesProperly(text string) string {
	runes := []rune(text)
	n := len(runes)

	// Find all positions of single quotes
	quotePositions := []int{}
	for i := 0; i < n; i++ {
		if runes[i] == '\'' {
			quotePositions = append(quotePositions, i)
		}
	}

	// Process pairs: first and second quote, third and fourth, etc.
	// We will mark characters for deletion using a boolean slice
	toDelete := make([]bool, n)

	for p := 0; p+1 < len(quotePositions); p += 2 {
		openPos := quotePositions[p]
		closePos := quotePositions[p+1]

		// Remove spaces after the opening quote
		j := openPos + 1
		for j < closePos && runes[j] == ' ' {
			toDelete[j] = true
			j++
		}

		// Remove spaces before the closing quote
		k := closePos - 1
		for k > openPos && runes[k] == ' ' {
			toDelete[k] = true
			k--
		}
	}

	// Build the result skipping deleted characters
	var result strings.Builder
	for i, ch := range runes {
		if !toDelete[i] {
			result.WriteRune(ch)
		}
	}

	return result.String()
}

// removeCharAt removes the character at position pos from a rune slice
func removeCharAt(runes []rune, pos int) []rune {
	// Build a new slice without the character at pos
	newRunes := []rune{}
	for i, ch := range runes {
		if i != pos {
			newRunes = append(newRunes, ch)
		}
	}
	return newRunes
}

// -------------------------------------------------------
// HELPER: check if a rune is a vowel (used for a -> an)
// -------------------------------------------------------

// isVowel returns true if the character is a vowel
func isVowel(ch rune) bool {
	return unicode.ToLower(ch) == 'a' ||
		unicode.ToLower(ch) == 'e' ||
		unicode.ToLower(ch) == 'i' ||
		unicode.ToLower(ch) == 'o' ||
		unicode.ToLower(ch) == 'u'
}