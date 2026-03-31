package main

import "strings"

// fixPunctuation ensures punctuation marks are attached to the preceding word
// and separated from the following word by exactly one space.
// Consecutive punctuation groups like "..." or "!?" are kept together.
func fixPunctuation(text string) string {
	var result strings.Builder

	chars := []rune(text)
	length := len(chars)

	i := 0
	for i < length {
		ch := chars[i]

		if isPunctuation(ch) {
			// Collect all consecutive punctuation into one group.
			punctGroup := []rune{}
			for i < length && isPunctuation(chars[i]) {
				punctGroup = append(punctGroup, chars[i])
				i++
			}

			// Strip the trailing space that was already written to the buffer.
			built := result.String()
			built = strings.TrimRight(built, " ")
			result.Reset()
			result.WriteString(built)

			// Write the punctuation group and add one space after it.
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

// isPunctuation reports whether the rune is one of the six recognised punctuation marks.
func isPunctuation(ch rune) bool {
	return ch == '.' || ch == ',' || ch == '!' || ch == '?' || ch == ':' || ch == ';'
}

// cleanSpaces collapses any sequence of whitespace into a single space.
func cleanSpaces(text string) string {
	words := strings.Fields(text)
	return strings.Join(words, " ")
}