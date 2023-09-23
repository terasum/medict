package main

import (
	"fmt"

	"github.com/agatan/bktree"
	levenshtein "github.com/creasty/go-levenshtein"
)

type word string

// Distance calculates hamming distance.
func (x word) Distance(e bktree.Entry) int {
	a := string(x)
	b := string(e.(word))

	return levenshtein.Distance(a, b)
}

func main() {
	var tree bktree.BKTree
	// add words
	words := []string{"apple", "banana", "orange", "peach", "bean", "tomato", "egg", "pineapple"}
	for _, w := range words {
		tree.Add(word(w))
	}

	// spell check
	results := tree.Search(word("peacn"), 2)
	fmt.Println("Input is peacn. Did you mean:")
	for _, result := range results {
		fmt.Printf("\t%s (distance: %d)\n", result.Entry.(word), result.Distance)
	}
}
