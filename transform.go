package main
 
import "strings"
 
// ApplyAllTransformations is the central orchestrator.
// It runs all transformation steps in the correct order and returns the final result.
func ApplyAllTransformations(text string) string {
	words := strings.Fields(text)
 
	words = applyHexAndBin(words)
	words = applyAToAn(words)
	words = applyCaseTransformations(words)
 
	result := strings.Join(words, " ")
 
	result = fixPunctuation(result)
	result = fixQuotesProperly(result)
 
	return result
}