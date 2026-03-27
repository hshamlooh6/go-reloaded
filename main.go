package main

import (
	"fmt"
	"os"
)

// The entry point. If this breaks, everything breaks. No pressure.
func main() {
	// Make sure the user remembered to pass both files, or we bail.
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go input.txt output.txt")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Try to read the input file. If it doesn't exist, that's on you.
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		os.Exit(1)
	}

	text := string(data)

	// The real work happens here. This one line does a lot.
	result := ApplyAllTransformations(text)

	// Write the result out. If this fails, the disk probably hates us.
	err = os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		os.Exit(1)
	}

	fmt.Println("Done! Modified text written to", outputFile)
}
