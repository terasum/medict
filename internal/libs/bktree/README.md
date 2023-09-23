## bktree

This is an implementation of [BK-tree](https://en.wikipedia.org/wiki/BK-tree) for golang.
BK-tree is a tree data structure for similarity search in a metric space.
Using BK-tree, you can search neighbors of a data from the metric space efficiently.

### Performance

Search similar values from 1,000,000 data.
Data is 64 bits integer, and distance function is hamming distance.
(see `bktree_test.go` for detail)

```
BenchmarkSearch_ExactForLargeTree-2         	 1000000	      3540 ns/op
BenchmarkSearch_Tolerance1ForLargeTree-2    	   30000	     85227 ns/op
BenchmarkSearch_Tolerance2ForLargeTree-2    	    3000	    907897 ns/op
BenchmarkSearch_Tolerance4ForLargeTree-2    	     200	  17273811 ns/op
BenchmarkSearch_Tolerance8ForLargeTree-2    	      10	 234244862 ns/op
BenchmarkSearch_Tolerance32ForLargeTree-2   	       3	 722881937 ns/op

BenchmarkLinearSearchForLargeSet-2          	      10	 281977270 ns/op
```

If the tolerance is small enough, BK-tree is much faster than naive linear search.


### Example

see `_example/` direcotry.

### Install

```
$ go get github.com/agatan/bktree
```

### License

MIT

