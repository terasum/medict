package bktree

import (
	"encoding/gob"
	"io"
)

// sNode is a serializable representation of node for gob persistence.
// Entry is stored as []byte via caller-provided marshal/unmarshal functions.
type sNode struct {
	EntryData []byte
	Children  []sChild
}

type sChild struct {
	Distance int
	Node     *sNode
}

// Save serializes the BK-tree structure to w. The caller provides marshalEntry
// to serialize each Entry to bytes (bktree doesn't know the concrete type).
// Format: gob-encoded *sNode (pre-order tree with Entry as []byte).
func (bk *BKTree) Save(w io.Writer, marshalEntry func(Entry) ([]byte, error)) error {
	if bk.root == nil {
		return nil
	}
	root := toSerializable(bk.root, marshalEntry)
	return gob.NewEncoder(w).Encode(root)
}

// Load deserializes a BK-tree from r. The caller provides unmarshalEntry to
// reconstruct each Entry from bytes. The reconstructed tree has the exact same
// structure as the original — no Distance computation needed (instant load).
func Load(r io.Reader, unmarshalEntry func([]byte) (Entry, error)) (*BKTree, error) {
	var sn sNode
	if err := gob.NewDecoder(r).Decode(&sn); err != nil {
		return nil, err
	}
	return &BKTree{root: fromSerializable(&sn, unmarshalEntry)}, nil
}

func toSerializable(n *node, marshal func(Entry) ([]byte, error)) *sNode {
	data, _ := marshal(n.entry)
	sn := &sNode{EntryData: data}
	for _, c := range n.children {
		sn.Children = append(sn.Children, sChild{
			Distance: c.distance,
			Node:     toSerializable(c.node, marshal),
		})
	}
	return sn
}

func fromSerializable(sn *sNode, unmarshal func([]byte) (Entry, error)) *node {
	entry, _ := unmarshal(sn.EntryData)
	n := &node{entry: entry}
	for _, c := range sn.Children {
		n.children = append(n.children, struct {
			distance int
			node     *node
		}{c.Distance, fromSerializable(c.Node, unmarshal)})
	}
	return n
}
