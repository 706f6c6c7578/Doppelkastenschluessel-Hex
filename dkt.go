package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"
)

const alphabet = "0123456789ABCDEF"

func main() {
	numSheets := flag.Int("a", 2, "Anzahl der Schlüssel (A und B Paare)")
	flag.Parse()

	if *numSheets < 1 {
		log.Fatal("Die Anzahl der Schlüssel muss mindestens 1 betragen")
	}

	// Seed basierend auf dem aktuellen Tag (UTC)
	now := time.Now().UTC()
	seed := now.Year()*10000 + int(now.Month())*100 + now.Day() // YYYYMMDD als Seed
	formattedDate := fmt.Sprintf("%02d.%02d.%d", now.Day(), int(now.Month()), now.Year())

	fmt.Printf("Tagesschlüssel für den %s (UTC)\n", formattedDate)

	for i := 0; i < *numSheets; i++ {
		fmt.Printf("\nTagesschlüssel %d:\n", i+1)
		matrixA := generateMatrix(int64(seed + i))
		matrixB := generateMatrix(int64(seed + i + 1000)) // Offset für Matrix B

		printKastenPair(matrixA, matrixB)
	}
}

func generateMatrix(seed int64) [4][4]rune {
	var matrix [4][4]rune
	letters := shuffle(alphabet, seed)
	index := 0

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			matrix[row][col] = rune(letters[index])
			index++
		}
	}

	return matrix
}

func shuffle(input string, seed int64) string {
	r := rand.New(rand.NewSource(seed)) // Deterministischer Zufallsgenerator basierend auf Seed
	runes := []rune(input)

	for i := len(runes) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func printKastenPair(matrixA, matrixB [4][4]rune) {
	// Überschrift mit festen Spaltenabständen
	fmt.Printf("%-12s %-12s\n", "Kasten: A", "Kasten: B")

	// Matrixinhalt
	for row := 0; row < 4; row++ {
		// Kasten A
		for col := 0; col < 4; col++ {
			fmt.Printf("%c ", matrixA[row][col])
		}
		// Abstand zwischen A und B
		fmt.Print("     ")
		// Kasten B
		for col := 0; col < 4; col++ {
			fmt.Printf("%c ", matrixB[row][col])
		}
		fmt.Println()
	}
}
