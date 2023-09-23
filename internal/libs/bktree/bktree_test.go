package bktree

import (
	"fmt"
	"math/rand"
	"testing"
)

type entry uint64

func hamming(a, b uint64) int {
	count := 0
	var k uint64 = 1
	for i := 0; i < 64; i++ {
		if a&k != b&k {
			count++
		}
		k <<= 1
	}
	return count
}

func (e entry) Distance(x Entry) int {
	a := uint64(e)
	b := uint64(x.(entry))

	return hamming(a, b)
}

func TestEmptySearch(t *testing.T) {
	var tree BKTree
	results := tree.Search(entry(0), 0, 100)
	if len(results) != 0 {
		t.Fatalf("empty tree should return empty results, bot got %d results", len(results))
	}
}

func TestExactMatch(t *testing.T) {
	var tree BKTree
	for i := 0; i < 100; i++ {
		tree.Add(entry(i))
	}

	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("searching %d", i), func(st *testing.T) {
			results := tree.Search(entry(i), 0, 100)
			if len(results) != 1 {
				st.Fatalf("exact match should return only one result, but got %d results (%#v)", len(results), results)
			}
			if results[0].Distance != 0 {
				st.Fatalf("exact match result should have 0 as Distance field, but got %d", results[0].Distance)
			}
			if int(results[0].Entry.(entry)) != i {
				st.Fatalf("expected result entry value is %d, but got %d", i, int(results[0].Entry.(entry)))
			}
		})
	}
}

func TestFuzzyMatch(t *testing.T) {
	var tree BKTree
	for i := 0; i < 100; i++ {
		tree.Add(entry(i))
	}

	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("searching %d", i), func(st *testing.T) {
			results := tree.Search(entry(i), 2, 100)
			for _, result := range results {
				if result.Distance > 2 {
					st.Fatalf("Distance fields of results should be less than or equal to 2, but got %d", result.Distance)
				}
				if result.Entry.Distance(entry(i)) > 2 {
					st.Fatalf("distances of result entries should be less than or equal to 2, but got %d", result.Distance)
				}
			}
		})
	}
}

const largeSize int = 1000000
const smallSize int = 1000

func BenchmarkConstruct(b *testing.B) {
	randoms := make([]uint64, 100000)
	for i := range randoms {
		randoms[i] = rand.Uint64()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var tree BKTree
		for _, r := range randoms {
			tree.Add(entry(r))
		}
	}
}

func makeRandomTree(size int) *BKTree {
	randoms := make([]int, size)
	for i := range randoms {
		randoms[i] = rand.Int()
	}
	var tree BKTree
	for _, r := range randoms {
		tree.Add(entry(r))
	}
	return &tree
}

func BenchmarkSearch_ExactForLargeTree(b *testing.B) {
	tree := makeRandomTree(largeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		tree.Search(entry(needle), 0, 100)
	}
}

func BenchmarkSearch_Tolerance1ForLargeTree(b *testing.B) {
	tree := makeRandomTree(largeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		tree.Search(entry(needle), 1, 100)
	}
}

func BenchmarkSearch_Tolerance2ForLargeTree(b *testing.B) {
	tree := makeRandomTree(largeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		tree.Search(entry(needle), 2, 100)
	}
}

func BenchmarkSearch_Tolerance4ForLargeTree(b *testing.B) {
	tree := makeRandomTree(largeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		tree.Search(entry(needle), 4, 100)
	}
}

func BenchmarkSearch_Tolerance8ForLargeTree(b *testing.B) {
	tree := makeRandomTree(largeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		tree.Search(entry(needle), 8, 100)
	}
}

func BenchmarkSearch_Tolerance32ForLargeTree(b *testing.B) {
	tree := makeRandomTree(largeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		tree.Search(entry(needle), 32, 100)
	}
}

func BenchmarkLinearSearchForLargeSet(b *testing.B) {
	randoms := make([]uint64, largeSize)
	for i := range randoms {
		randoms[i] = rand.Uint64()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		cnt := 0
		for _, c := range randoms {
			if hamming(c, needle) <= 1 {
				cnt++
			}
		}
	}
}

func BenchmarkSearch_ExactForSmallTree(b *testing.B) {
	tree := makeRandomTree(smallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		tree.Search(entry(needle), 0, 100)
	}
}

func BenchmarkSearch_Tolerance1ForSmallTree(b *testing.B) {
	tree := makeRandomTree(smallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		tree.Search(entry(needle), 1, 100)
	}
}

func BenchmarkSearch_Tolerance2ForSmallTree(b *testing.B) {
	tree := makeRandomTree(smallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		tree.Search(entry(needle), 2, 100)
	}
}

func BenchmarkSearch_Tolerance4ForSmallTree(b *testing.B) {
	tree := makeRandomTree(smallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		tree.Search(entry(needle), 4, 100)
	}
}

func BenchmarkSearch_Tolerance8ForSmallTree(b *testing.B) {
	tree := makeRandomTree(smallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		tree.Search(entry(needle), 8, 100)
	}
}

func BenchmarkSearch_Tolerance32ForSmallTree(b *testing.B) {
	tree := makeRandomTree(smallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		tree.Search(entry(needle), 32, 100)
	}
}

func BenchmarkLinearSearchForSmallSet(b *testing.B) {
	randoms := make([]uint64, smallSize)
	for i := range randoms {
		randoms[i] = rand.Uint64()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.Uint64()
		cnt := 0
		for _, c := range randoms {
			if hamming(c, needle) <= 1 {
				cnt++
			}
		}
	}
}
