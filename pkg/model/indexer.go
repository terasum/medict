package model

import (
	"container/heap"
	"github.com/agatan/bktree"
	"github.com/terasum/medict/internal/mdictparser"
)

type Indexer struct {
	dtype   DictType
	rawList []*WordItem
	handler *mdictparser.MdictParser
	indexed bool
	counter uint64
	tree    bktree.BKTree
}

func NewIndexer(dtype DictType, handler *mdictparser.MdictParser) *Indexer {
	return &Indexer{
		dtype:   dtype,
		rawList: make([]*WordItem, 0, 100),
		handler: handler,
	}
}

func (idx *Indexer) Indexing() (error, uint64) {
	if idx.indexed {
		return nil, idx.counter
	}
	items, _, err := idx.handler.AllWords()
	if err != nil {
		return err, 0
	}
	for _, item := range items {
		newWordItem, err2 := NewWordItem(idx.dtype, item.KeyWord, item.RecordStart)
		if err2 != nil {
			return err2, 0
		}
		if idx.dtype == MDD {
			idx.rawList = append(idx.rawList, newWordItem)
		} else {
			//idx.tree.Add(newWordItem)
			idx.rawList = append(idx.rawList, newWordItem)
		}
		idx.counter++
	}
	idx.indexed = true
	return nil, idx.counter
}

func (idx *Indexer) listSearch(needle *WordItem, tolerance int) []*WordItem {
	// reduce search number
	counter := 0
	//upperLimit := 20000
	// use priority queue
	results := NewPQ(needle)
	if len(idx.rawList) == 0 {
		return results.list
	}
	for _, wi := range idx.rawList {
		d := wi.Distance(needle)
		if d <= tolerance {
			//results = append(results, wi)
			results.Push(wi)
			counter++
			if counter > 50 {
				break
			}
		}
	}
	heap.Init(results)
	return results.list
}

func (idx *Indexer) SimWord(word string, distance int) ([]*WordItem, error) {
	// spell check
	target, err := NewWordItem(MDX, word, 0)
	if err != nil {
		return nil, err
	}
	//results := idx.tree.Search(target, distance)
	//resArr := make([]*WordItem, len(results))
	//for i, result := range results {
	//	resArr[i] = (result.Entry).(*WordItem)
	//}
	resArr := idx.listSearch(target, distance)
	return resArr, nil
}

func (idx *Indexer) Size() uint64 {
	return idx.counter
}
