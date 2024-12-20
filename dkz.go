package main

import (
	"crypto/sha512"
	"encoding/binary"
	"flag"
	"fmt"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/hkdf"
	"math/big"
	"math/rand"
	"time"
)

func generateNonZeroNumber(rng *rand.Rand, maxNum *big.Int) int64 {
	for {
		num := rng.Int63n(maxNum.Int64())
		if num > 0 {
			return num
		}
	}
}

func createDeterministicRNG(password, salt string, runCounter int) *rand.Rand {
	// Get start of current day in UTC
	now := time.Now().UTC()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Argon2id parameters
	timeParam := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLen := uint32(32)

	// Generate key using Argon2id
	argonKey := argon2.IDKey([]byte(password), []byte(salt), timeParam, memory, threads, keyLen)

	// Initialize HKDF with SHA-512
	info := []byte(fmt.Sprintf("%d%d", midnight.Unix(), runCounter))
	hkdf := hkdf.New(sha512.New, argonKey, nil, info)
	
	var seed [8]byte
	if _, err := hkdf.Read(seed[:]); err != nil {
		panic(err)
	}
	
	return rand.New(rand.NewSource(int64(binary.BigEndian.Uint64(seed[:]))))
}

func main() {
	count := flag.Int("a", 4, "Anzahl der zu erzeugenden Zahlen")
	maxDigits := flag.Int("z", 2, "Maximale Anzahl von Dezimalstellen pro Zahl")
	separator := flag.String("t", ",", "Trennzeichen zwischen Zahlen")
	runs := flag.Int("d", 1, "Anzahl der Durchläufe")
	password := flag.String("p", "", "Passwort für deterministische Generierung")
	salt := flag.String("s", "", "Salz für deterministische Generierung")
	flag.Parse()

	maxNum := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(*maxDigits)), nil)

	for run := 0; run < *runs; run++ {
		var rng *rand.Rand
		if *password != "" {
			rng = createDeterministicRNG(*password, *salt, run)
		}

		for i := 0; i < *count; i++ {
			var num int64
			if rng != nil {
				num = generateNonZeroNumber(rng, maxNum)
			} else {
				cryptoNum := generateNonZeroNumber(rand.New(rand.NewSource(time.Now().UnixNano())), maxNum)
				num = cryptoNum
			}

			if i < *count-1 {
				fmt.Printf("%v%s", num, *separator)
			} else {
				fmt.Printf("%v", num)
			}
		}
		fmt.Println()
	}
}
