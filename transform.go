package main

import (
	"strconv"
	"strings"
	"unicode"
)

func ApplyAllTransformations(text string) string {
	words := strings.Fields(text)

	words = applyHexAndBin(words)
	words = applyCaseTransformations(words)
	words = applyAToAn(words)

	result := strings.Join(words, " ")

	result = fixPunctuation(result)
	result = fixQuotesProperly(result)

	return result
}

func applyHexAndBin(words []string) []string {
	result := []string{}

	i := 0
	for i < len(words) {
		word := words[i]

		if word == "(hex)" || word == "(bin)" {
			if len(result) > 0 {
				lastWord := result[len(result)-1]

				var converted string
				var err error

				if word == "(hex)" {
					converted, err = hexToDecimal(lastWord)
				} else {
					converted, err = binToDecimal(lastWord)
				}

				if err == nil {
					result[len(result)-1] = converted
				}
			}
		} else {
			result = append(result, word)
		}

		i++
	}

	return result
}

func hexToDecimal(hex string) (string, error) {
	number, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(number, 10), nil
}

func binToDecimal(bin string) (string, error) {
	number, err := strconv.ParseInt(bin, 2, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(number, 10), nil
}

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

func applyAToAn(words []string) []string {
	vowelsAndH := "aeiouAEIOUhH"

	for i := 0; i < len(words)-1; i++ {
		currentWord := words[i]
		nextWord := words[i+1]

		if currentWord == "a" || currentWord == "A" {
			if len(nextWord) > 0 && strings.ContainsRune(vowelsAndH, rune(nextWord[0])) {
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

func fixPunctuation(text string) string {
	var result strings.Builder

	chars := []rune(text)
	length := len(chars)

	i := 0
	for i < length {
		ch := chars[i]

		if isPunctuation(ch) {
			punctGroup := []rune{}
			for i < length && isPunctuation(chars[i]) {
				punctGroup = append(punctGroup, chars[i])
				i++
			}

			built := result.String()
			built = strings.TrimRight(built, " ")
			result.Reset()
			result.WriteString(built)

			result.WriteString(string(punctGroup))

			if i < length && chars[i] != ' ' {
				result.WriteRune(' ')
			}
		} else {
			result.WriteRune(ch)
			i++
		}
	}

	return cleanSpaces(result.String())
}

func isPunctuation(ch rune) bool {
	return ch == '.' || ch == ',' || ch == '!' || ch == '?' || ch == ':' || ch == ';'
}

func cleanSpaces(text string) string {
	words := strings.Fields(text)
	return strings.Join(words, " ")
}

func fixQuotesProperly(text string) string {
	runes := []rune(text)
	n := len(runes)

	quotePositions := []int{}
	for i := 0; i < n; i++ {
		if runes[i] == '\'' {
			quotePositions = append(quotePositions, i)
		}
	}

	toDelete := make([]bool, n)

	for p := 0; p+1 < len(quotePositions); p += 2 {
		openPos := quotePositions[p]
		closePos := quotePositions[p+1]

		j := openPos + 1
		for j < closePos && runes[j] == ' ' {
			toDelete[j] = true
			j++
		}

		k := closePos - 1
		for k > openPos && runes[k] == ' ' {
			toDelete[k] = true
			k--
		}
	}

	var result strings.Builder
	for i, ch := range runes {
		if !toDelete[i] {
			result.WriteRune(ch)
		}
	}

	return result.String()
}

func isVowel(ch rune) bool {
	return unicode.ToLower(ch) == 'a' ||
		unicode.ToLower(ch) == 'e' ||
		unicode.ToLower(ch) == 'i' ||
		unicode.ToLower(ch) == 'o' ||
		unicode.ToLower(ch) == 'u'
}
