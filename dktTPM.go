package main

import (
    "flag"
    "fmt"
    "log"
    "github.com/google/go-tpm/legacy/tpm2"
    "time"
    "io"
)

const alphabet = "0123456789ABCDEF"

func main() {
    numSheets := flag.Int("a", 2, "Anzahl der Schlüssel (A und B Paare)")
    dateStr := flag.String("d", "", "Datum für die Schlüssel (Format: DD.MM.YYYY)")
    flag.Parse()

    if *numSheets < 1 {
        log.Fatal("Die Anzahl der Schlüssel muss mindestens 1 betragen")
    }

    var formattedDate string
    if *dateStr != "" {
        date, err := time.Parse("02.01.2006", *dateStr)
        if err != nil {
            log.Fatal("Ungültiges Datumsformat. Bitte DD.MM.YYYY verwenden")
        }
        formattedDate = date.Format("02.01.2006")
    } else {
        now := time.Now().UTC()
        formattedDate = now.Format("02.01.2006")
    }

    rwc, err := tpm2.OpenTPM()
    if err != nil {
        fmt.Printf("TPM öffnen fehlgeschlagen: %v\n", err)
        return
    }
    defer rwc.Close()

    fmt.Printf("Tagesschlüssel für den %s\n", formattedDate)

    for i := 0; i < *numSheets; i++ {
        fmt.Printf("\nTagesschlüssel %d:\n", i+1)
        matrixA := generateMatrix(rwc)
        matrixB := generateMatrix(rwc)
        printKastenPair(matrixA, matrixB)
    }
}

func generateMatrix(rwc io.ReadWriteCloser) [4][4]rune {
    var matrix [4][4]rune
    letters := shuffleTPM(alphabet, rwc)
    index := 0

    for row := 0; row < 4; row++ {
        for col := 0; col < 4; col++ {
            matrix[row][col] = rune(letters[index])
            index++
        }
    }

    return matrix
}

func shuffleTPM(input string, rwc io.ReadWriteCloser) string {
    runes := []rune(input)
    n := len(runes)

    for i := n - 1; i > 0; i-- {
        random, _ := tpm2.GetRandom(rwc, 1)
        j := int(random[0] % uint8(i+1))
        runes[i], runes[j] = runes[j], runes[i]
    }

    return string(runes)
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
