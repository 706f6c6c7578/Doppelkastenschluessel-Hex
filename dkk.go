package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
    "strconv"
)

const matrixSize = 4

type KeyMatrix struct {
    A [matrixSize][matrixSize]string
    B [matrixSize][matrixSize]string
}

func main() {
    decrypt := flag.Bool("e", false, "Eingabe entschlüsseln")
    keyFile := flag.String("k", "keys.txt", "Pfad zur Schlüsseldatei")
    keySequence := flag.String("s", "1,2", "Komma-separierte Folge von Schlüsselzahlen")
    flag.Parse()

    keyNumbers := parseKeySequence(*keySequence)
    allKeys := make(map[int]KeyMatrix)
    
    // Read all required keys into a map
    for _, keyNum := range keyNumbers {
        key, err := readSingleKey(*keyFile, keyNum)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Fehler beim Lesen des Schlüssels %d: %v\n", keyNum, err)
            os.Exit(1)
        }
        allKeys[keyNum] = key
    }

    scanner := bufio.NewScanner(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    for scanner.Scan() {
        line := scanner.Text()
        currentResult := line

        if *decrypt {
            // For decryption, process in exact reverse order of keyNumbers
            for i := len(keyNumbers) - 1; i >= 0; i-- {
                keyNum := keyNumbers[i]
                key := allKeys[keyNum]
                currentResult = encryptSingleKey(currentResult, key, true)
            }
        } else {
            // Encryption remains the same
            for _, keyNum := range keyNumbers {
                key := allKeys[keyNum]
                currentResult = encryptSingleKey(currentResult, key, false)
            }
        }

        fmt.Fprintln(writer, currentResult)
    }
}

func encryptSingleKey(input string, key KeyMatrix, decrypt bool) string {
    fields := strings.Fields(input)
    var output []string

    for _, field := range fields {
        if len(field) != 2 {
            continue
        }

        high := string(field[0])
        low := string(field[1])

        if decrypt {
            output = append(output, decryptHex(high, low, key))
        } else {
            output = append(output, encryptHex(high, low, key))
        }
    }

    return strings.Join(output, " ")
}

func readSingleKey(filename string, targetKeyNum int) (KeyMatrix, error) {
    file, err := os.Open(filename)
    if err != nil {
        return KeyMatrix{}, err
    }
    defer file.Close()

    var currentKeyNum int
    var matrixLines []string
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()

        if strings.Contains(line, "Tagesschlüssel") && !strings.Contains(line, "UTC") {
            currentKeyNum++
            matrixLines = nil
            continue
        }

        if currentKeyNum == targetKeyNum && strings.Contains(line, "      ") {
            matrixLines = append(matrixLines, line)
            if len(matrixLines) == matrixSize {
                return parseKeyMatrix(matrixLines)
            }
        }
    }

    return KeyMatrix{}, fmt.Errorf("Schlüssel Nr. %d nicht gefunden", targetKeyNum)
}

func parseKeySequence(sequence string) []int {
    parts := strings.Split(sequence, ",")
    numbers := make([]int, 0, len(parts))
    
    for _, part := range parts {
        if num, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
            numbers = append(numbers, num)
        }
    }
    
    return numbers
}

func parseKeyMatrix(lines []string) (KeyMatrix, error) {
    var key KeyMatrix
    
    for i := 0; i < matrixSize; i++ {
        parts := strings.SplitN(lines[i], "      ", 2)
        if len(parts) != 2 {
            return key, fmt.Errorf("ungültiges Matrixformat")
        }

        leftFields := strings.Fields(parts[0])
        rightFields := strings.Fields(parts[1])

        if len(leftFields) != matrixSize || len(rightFields) != matrixSize {
            return key, fmt.Errorf("invalid matrix size")
        }

        for j := 0; j < matrixSize; j++ {
            key.A[i][j] = leftFields[j]
            key.B[i][j] = rightFields[j]
        }
    }

    return key, nil
}

func encryptHex(high, low string, key KeyMatrix) string {
    hRow, hCol := findPosition(key.A, high)
    lRow, lCol := findPosition(key.B, low)

    if hRow == -1 || hCol == -1 || lRow == -1 || lCol == -1 {
        return "??"
    }

    encryptedHigh := key.B[hRow][hCol]
    encryptedLow := key.A[lRow][lCol]

    return encryptedHigh + encryptedLow
}

func decryptHex(high, low string, key KeyMatrix) string {
    hRow, hCol := findPosition(key.B, high)
    lRow, lCol := findPosition(key.A, low)

    if hRow == -1 || hCol == -1 || lRow == -1 || lCol == -1 {
        return "??"
    }

    decryptedHigh := key.A[hRow][hCol]
    decryptedLow := key.B[lRow][lCol]

    return decryptedHigh + decryptedLow
}

func findPosition(matrix [matrixSize][matrixSize]string, value string) (int, int) {
    for i := 0; i < matrixSize; i++ {
        for j := 0; j < matrixSize; j++ {
            if matrix[i][j] == value {
                return i, j
            }
        }
    }
    return -1, -1
}
