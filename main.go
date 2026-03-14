package main

import (
	"fmt"
	"os"
)

func main() {
	// Step 1: Make sure the user gave us exactly 2 arguments (input file and output file)
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go input.txt output.txt")
		os.Exit(1)
	}

	// Step 2: Save the file names from the arguments
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Step 3: Read the content of the input file
	// os.ReadFile reads the whole file and gives us the bytes
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		os.Exit(1)
	}

	// Step 4: Convert the bytes to a string so we can work with the text
	text := string(data)

	// Step 5: Apply all the transformations to the text
	result := ApplyAllTransformations(text)

	// Step 6: Write the result into the output file
	// 0644 means: the owner can read/write, others can only read
	err = os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		os.Exit(1)
	}

	// Step 7: Tell the user everything went well
	fmt.Println("Done! Modified text written to", outputFile)
}
