package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func encode(input []byte) string {
	// Calculate padding needed
	inputLen := len(input)
	paddingLen := 16 - (inputLen % 16)
	if paddingLen == 16 {
		paddingLen = 0
	}

	// Generate random padding
	padding := make([]byte, paddingLen)
	rand.Read(padding)
	for i := range padding {
		padding[i] = byte(65 + (padding[i] % 26)) // Uppercase A-Z
	}

	// Combine input and padding
	fullData := append(input, padding...)

	// Convert to hex strings, 8 bytes (16 chars) per line
	var result strings.Builder
	for i := 0; i < len(fullData); i += 8 {
		end := i + 8
		if end > len(fullData) {
			end = len(fullData)
		}
		chunk := fullData[i:end]
		for _, b := range chunk {
			fmt.Fprintf(&result, "%02X", b)
		}
		result.WriteString("\n")
	}
	return strings.TrimRight(result.String(), "\n")
}

func decode(input string) ([]byte, error) {
	// Remove newlines
	input = strings.ReplaceAll(input, "\n", "")

	if len(input)%2 != 0 {
		return nil, fmt.Errorf("ungültige Länge der Hex-Zeichenkette")
	}

	decoded := make([]byte, len(input)/2)
	for i := 0; i < len(input); i += 2 {
		var value byte
		_, err := fmt.Sscanf(input[i:i+2], "%02X", &value)
		if err != nil {
			return nil, err
		}
		decoded[i/2] = value
	}
	return decoded, nil
}

func main() {
	decodeFlag := flag.Bool("d", false, "Dekodiermodus")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	input, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fehler beim Lesen der Eingabe: %v\n", err)
		os.Exit(1)
	}

	if *decodeFlag {
		decoded, err := decode(string(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fehler bei der Dekodierung: %v\n", err)
			os.Exit(1)
		}
		os.Stdout.Write(decoded) // Original format
	} else {
		output := encode(input)
		os.Stdout.Write([]byte(output))
		fmt.Println()
	}
}
