package main

import (
	"fmt"

	"github.com/agatan/bktree"
)

type hashValue uint16

// Distance calculates hamming distance.
func (h hashValue) Distance(e bktree.Entry) int {
	a := uint16(h)
	b := uint16(e.(hashValue))

	d := 0
	var k uint16 = 1
	for i := 0; i < 16; i++ {
		if a&k != b&k {
			d++
		}
		k <<= 1
	}
	return d
}

func main() {
	var tree bktree.BKTree
	// add 0x0000 to 0xffff to the tree.
	for i := 0; i < 0xffff; i++ {
		tree.Add(hashValue(i))
	}

	// search neighbors of 0x00000000 whose distances are less than or equal to 1.
	results := tree.Search(hashValue(0), 1)
	for _, result := range results {
		fmt.Printf("%016b (distance: %d)\n", result.Entry.(hashValue), result.Distance)
	}
}
