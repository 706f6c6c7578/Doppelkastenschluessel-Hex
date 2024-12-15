package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"flag"
	"fmt"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/hkdf"
	"io"
	"log"
	"math/big"
	"time"
)

const alphabet = "0123456789ABCDEF"

func main() {
	numSheets := flag.Int("a", 1, "Anzahl der Schlüssel (A und B Paare)")
	password := flag.String("p", "", "Passwort")
	salt := flag.String("s", "", "Salz")
	flag.Parse()

	if *numSheets < 1 {
		log.Fatal("Die Anzahl der Schlüssel muss mindestens 1 betragen")
	}

	if *password == "" || *salt == "" {
		log.Fatal("Passwort und Salz dürfen nicht leer sein")
	}

	// Seed basierend auf dem aktuellen Tag (UTC)
	now := time.Now().UTC()
	formattedDate := fmt.Sprintf("%02d.%02d.%d", now.Day(), int(now.Month()), now.Year())
	seed := generateSeed(*password, *salt, formattedDate)
	fmt.Printf("Tagesschlüssel für den %s (UTC)\n", formattedDate)

	for i := 0; i < *numSheets; i++ {
		fmt.Printf("\nTagesschlüssel %d:\n", i+1)
		matrixA := generateMatrix(seed + int64(i))
		matrixB := generateMatrix(seed + int64(i) + 1000) // Offset für Matrix B

		printKastenPair(matrixA, matrixB)
	}
}

func generateSeed(password, salt, date string) int64 {
	saltBytes := []byte(salt)
	passwordBytes := []byte(password + date)
	key := argon2.IDKey(passwordBytes, saltBytes, 1, 64*1024, 4, 64)

	hkdf := hkdf.New(sha512.New, key, saltBytes, nil)
	var seedBytes [8]byte
	if _, err := hkdf.Read(seedBytes[:]); err != nil {
		log.Fatalf("Fehler beim Generieren des Seeds: %v", err)
	}

	seed := int64(binary.BigEndian.Uint64(seedBytes[:]))
	return seed
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
	r := newCryptoRand(seed)
	runes := []rune(input)

	for i := len(runes) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func newCryptoRand(seed int64) *cryptoRandReader {
	h := sha512.New()
	_ = binary.Write(h, binary.BigEndian, seed)
	return &cryptoRandReader{reader: hkdf.New(sha512.New, h.Sum(nil), nil, nil)}
}

type cryptoRandReader struct {
	reader io.Reader
}

func (r *cryptoRandReader) Intn(n int) int {
	bigN := big.NewInt(int64(n))
	val, err := rand.Int(r.reader, bigN)
	if err != nil {
		log.Fatalf("Fehler bei Intn: %v", err)
	}
	return int(val.Int64())
}

func printKastenPair(matrixA, matrixB [4][4]rune) {
	fmt.Printf("%-12s %-12s\n", "Kasten: A", "Kasten: B")

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			fmt.Printf("%c ", matrixA[row][col])
		}
		fmt.Print("     ")
		for col := 0; col < 4; col++ {
			fmt.Printf("%c ", matrixB[row][col])
		}
		fmt.Println()
	}
}
