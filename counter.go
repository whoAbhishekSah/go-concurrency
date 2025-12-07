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

type pair struct {
	word  string
	count int
}

func submitWords(next func() string) chan string {
	pending := make(chan string)
	go func() {
		for {
			word := next()
			if word != "" {
				pending <- word
			} else {
				close(pending)
				break
			}
		}
	}()
	return pending
}

func countWords(pending <-chan string) chan pair {
	counted := make(chan pair)
	go func() {
		for word := range pending {
			counted <- pair{word, countDigits(word)}
		}
		close(counted)
	}()
	return counted
}

func fillStats(counted <-chan pair) counter {
	stats := counter{}
	for countPair := range counted {
		stats[countPair.word] = countPair.count
	}
	return stats
}

// countDigitsInWords counts the number of digits in the words of a phrase.
func countDigitsInWords(next func() string) counter {
	// Returning an output channel from a function and
	// filling it within an internal goroutine
	// is a common pattern in Go.
	pending := submitWords(next)
	counted := countWords(pending)
	return fillStats(counted)
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

func main_bak() {
	phrase := "0ne 1wo thr33 4068"
	next := wordGenerator(phrase)
	counts := countDigitsInWords(next)
	printStats(counts)
}
