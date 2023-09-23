package bktree

type BKTree struct {
	root *node
}

type node struct {
	entry    Entry
	children []struct {
		distance int
		node     *node
	}
}

func (n *node) addChild(e Entry) {
	newnode := &node{entry: e}
loop:
	d := n.entry.Distance(e)
	for _, c := range n.children {
		if c.distance == d {
			n = c.node
			goto loop
		}
	}
	n.children = append(n.children, struct {
		distance int
		node     *node
	}{d, newnode})
}

type Entry interface {
	Distance(Entry) int
}

func (bk *BKTree) Add(entry Entry) {
	if bk.root == nil {
		bk.root = &node{
			entry: entry,
		}
		return
	}
	bk.root.addChild(entry)
}

type Result struct {
	Distance int
	Entry    Entry
}

func (bk *BKTree) Search(needle Entry, tolerance int, limit int) []*Result {
	results := make([]*Result, 0)
	if bk.root == nil {
		return results
	}
	candidates := []*node{bk.root}
	count := 0
	for len(candidates) != 0 {
		c := candidates[len(candidates)-1]
		candidates = candidates[:len(candidates)-1]
		d := c.entry.Distance(needle)
		if d <= tolerance {
			results = append(results, &Result{
				Distance: d,
				Entry:    c.entry,
			})
			count += 1
			if count >= limit {
				return results
			}

		}

		low, high := d-tolerance, d+tolerance
		for _, c := range c.children {
			if low <= c.distance && c.distance <= high {
				candidates = append(candidates, c.node)
			}
		}
	}
	return results
}
