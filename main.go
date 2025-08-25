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
func countDigitsInWords(next func() string) counter {
	type pair struct {
		word  string
		count int
	}

	pending := make(chan string)
	counted := make(chan pair)

	// sends words to be counted
	go func() {
		for {
			word := next()
			pending <- word
			if word == "" {
				return
			}
		}
	}()

	// count digits in words
	go func() {
		for {
			// should return when
			// there are no more words
			word := <-pending
			if word == "" {
				counted <- pair{"", 0}
				return
			}
			counted <- pair{word, countDigits(word)}
		}
	}()

	// fill stats by words
	stats := counter{}
	for {
		countPair := <-counted
		// how to break from the loop
		// when there are no more words?
		if countPair.word == "" {
			break
		}
		// where should the word come from?
		stats[countPair.word] = countPair.count
	}
	return stats
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

// wordGenerator returns a generator that yields words from a phrase.
func wordGenerator(phrase string) func() string {
	words := strings.Fields(phrase)
	idx := 0
	return func() string {
		if idx == len(words) {
			return ""
		}
		word := words[idx]
		idx++
		return word
	}
}

func main() {
	phrase := "0ne 1wo thr33 4068"
	next := wordGenerator(phrase)
	counts := countDigitsInWords(next)
	printStats(counts)
}
