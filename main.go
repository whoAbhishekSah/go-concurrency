// Counting digits in words.
package main

import (
	"fmt"
	"strings"
	"sync"
	"unicode"
)

// counter stores the number of digits in each word.
// The key is the word, and the value is the number of digits.
type counter map[string]int

// solution start

// countDigitsInWords counts the number of digits in the words of a phrase.
func countDigitsInWords(phrase string) counter {
    var wg sync.WaitGroup
    syncStats := new(sync.Map)
    words := strings.Fields(phrase)

	// wg.Add(len(words))
   for _, word := range words {
	    wg.Add(1)
		go func() {
            defer wg.Done()
            count := countDigits(word)
            // fmt.Println("------", word, count)
            syncStats.Store(word, count)
            // fmt.Println("stored")
        }()
   }
   wg.Wait()

   return asStats(syncStats)
}

// solution end

// countDigits returns the number of digits in a string.
func countDigits(str string) int {
	count := 0
	for _, char := range str {
		if unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

// asStats converts statistics from sync.Map to a regular map.
func asStats(m *sync.Map) counter {
	stats := counter{}
	m.Range(func(word, count any) bool {
		stats[word.(string)] = count.(int)
		return true
	})
	return stats
}

// printStats prints the number of digits in words.
func printStats(stats counter) {
	for word, count := range stats {
		fmt.Printf("%s: %d\n", word, count)
	}
}

func main() {
	phrase := "0ne 1wo thr33 4068"
	counts := countDigitsInWords(phrase)
	printStats(counts)
}
