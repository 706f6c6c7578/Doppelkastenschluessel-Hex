package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
)

func generateNonZeroNumber(maxNum *big.Int) *big.Int {
	for {
		num, err := rand.Int(rand.Reader, maxNum)
		if err != nil {
			continue
		}
		if num.Cmp(big.NewInt(0)) > 0 {
			return num
		}
	}
}

func main() {
	count := flag.Int("a", 32, "Anzahl der zu erzeugenden Zahlen")
	maxDigits := flag.Int("z", 2, "Maximale Anzahl von Ziffern pro Nummer")
	separator := flag.String("t", ",", "Trennzeichen zwischen Zahlen")
	flag.Parse()

	maxNum := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(*maxDigits)), nil)

	for i := 0; i < *count; i++ {
		num := generateNonZeroNumber(maxNum)

		if i < *count-1 {
			fmt.Printf("%v%s", num, *separator)
		} else {
			fmt.Printf("%v\n", num)
		}
	}
}
