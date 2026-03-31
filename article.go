package main

import "strings"

// applyAToAn corrects "a" to "an" when the next word begins with a vowel or h.
// Preserves the original capitalisation of the article.
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
