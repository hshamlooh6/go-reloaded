package main

import "strconv"

// applyHexAndBin scans the word list for (hex) and (bin) markers.
// When found, it replaces the preceding word with its decimal equivalent
// and discards the marker itself.
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

// hexToDecimal parses a hexadecimal string and returns its decimal representation.
func hexToDecimal(hex string) (string, error) {
	number, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(number, 10), nil
}

// binToDecimal parses a binary string and returns its decimal representation.
func binToDecimal(bin string) (string, error) {
	number, err := strconv.ParseInt(bin, 2, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(number, 10), nil
}
