package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Encoding: Create pairs of characters from the same position in both lines
func seriation(line1, line2 string) string {
	var result []string

	// Create the character pairs
	for i := 0; i < len(line1); i++ {
		pair := string(line1[i]) + string(line2[i])
		result = append(result, pair)
	}

	// Join the pairs with spaces
	return strings.Join(result, " ")
}

// Decoding: Split the pairs and reconstruct the original lines
func deseriation(encoded string) (string, string) {
	pairs := strings.Fields(encoded)
	var line1, line2 []rune

	// Split the pairs and rebuild the lines
	for _, pair := range pairs {
		line1 = append(line1, rune(pair[0]))
		line2 = append(line2, rune(pair[1]))
	}

	// Rebuild the two lines
	return string(line1), string(line2)
}

// Function to check if the number of lines is even (only important for encoding)
func checkEvenLines(lines []string) bool {
	return len(lines)%2 == 0
}

func main() {
	// Decode flag
	decode := flag.Bool("d", false, "Dekodiermodus")
	flag.Parse()

	// Read input data
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fehler beim Lesen der Eingabe: %v\n", err)
		return
	}

	// Input data as string
	text := string(bytes)

	// Split the text into lines
	lines := strings.Split(text, "\n")

	// Remove empty lines (e.g., due to trailing newline at the end)
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) == "" {
			lines = lines[:i] // Remove the empty line
		} else {
			break
		}
	}

	// Check if the number of lines is even (only for encoding)
	if !*decode && !checkEvenLines(lines) {
		fmt.Println("Error: The number of lines is odd. Please provide an even number of lines.")
		return
	}

	// Encoding or Decoding based on the flag
	if !*decode {
		// Encoding (Seriation)
		var encodedResult []string
		for i := 0; i < len(lines); i += 2 {
			encoded := seriation(lines[i], lines[i+1])
			encodedResult = append(encodedResult, encoded)
		}

		// Output the encoded pairs
		//fmt.Println("Seriation:")
		fmt.Println(strings.Join(encodedResult, "\n"))
	} else {
		// Decoding (Deseriation)
		var decodedResult []string
		for i := 0; i < len(lines); i++ {
			// Decoding: We need to split each line (encoded pairs) into two separate lines.
			decodedLine1, decodedLine2 := deseriation(lines[i])
			decodedResult = append(decodedResult, decodedLine1)
			decodedResult = append(decodedResult, decodedLine2)
		}

		// Output the decoded lines
		//fmt.Println("Decodiert:")
		fmt.Println(strings.Join(decodedResult, "\n"))
	}
}
