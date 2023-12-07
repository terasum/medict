package mdict_idxer

import (
	"github.com/dghubble/trie"
	"github.com/terasum/medict/internal/libs/bktree"
	art "github.com/terasum/medict/internal/libs/go-adaptive-radix-tree"
	"sync"

	"github.com/terasum/medict/pkg/model"
)

type searchTree struct {
	lock     *sync.Mutex
	bktree   *bktree.BKTree
	trie     *trie.RuneTrie
	artTree  art.Tree
	hasBuild bool
}

type entryWrapper struct {
	list []*entryWrapperItem
}

func (w *entryWrapper) Less(i, j int) bool {
	return w.list[i].distance < w.list[j].distance
}

func (w *entryWrapper) Len() int {
	return len(w.list)
}

func (w *entryWrapper) Swap(i, j int) {
	temp := w.list[i]
	w.list[i] = w.list[j]
	w.list[j] = temp
}

func (w *entryWrapper) toEntryList() []*model.MdictKeyWordIndex {
	entries := make([]*model.MdictKeyWordIndex, w.Len())
	for idx, result := range w.list {
		entries[idx] = result.entry
	}
	return entries
}

//func (m *MedictDBIndexer) initSearchTree() error {
//	db, err := m.acquire()
//	if err != nil {
//		return err
//	}
//	defer m.release(db)
//	m.searchTree = &searchTree{
//		lock:     new(sync.Mutex),
//		bktree:   nil,
//		trie:     nil,
//		artTree:  art.New(),
//		hasBuild: false,
//	}
//	ite := db.NewIterator(util.BytesPrefix([]byte(prefixKeyword)), nil)
//	for ite.Next() {
//		key := ite.Key()
//		m.searchTree.artTree.Insert(key, key)
//		log.Infof("build tree key: %s", strip(string(key)))
//	}
//	return nil
//}

func (st *searchTree) addKeyValue(key string, value []byte) error {

	return nil
}

//func (st *searchTree) search(keyword string) ([]string, error) {
//	if st.hasBuild {
//		return nil, errors.New("search hasn't built")
//	}
//	results := make([][]byte, 0)
//	st.artTree.ForEachPrefix(art.Key(keyword), func(node art.Node) (cont bool) {
//		if node == nil || node.Value() == nil {
//			return true
//		}
//		results = append(results, node.Value().([]byte))
//		return true
//	})
//	return results, nil
//}

//func (m *MedictDBIndexer) buildBKTree() error {
//	m.bktreeLock.Lock()
//	defer m.bktreeLock.Unlock()
//	if m.bktree != nil {
//		return nil
//	}
//	m.bktree = &bktree.BKTree{}
//	db, err := m.acquire()
//	if err != nil {
//		return err
//	}
//	defer m.release(db)
//
//	ite := db.NewIterator(util.BytesPrefix([]byte(prefixKeyword)), nil)
//	for ite.Next() {
//		index := new(model.MdictKeyWordIndex)
//		err = json.Unmarshal(ite.Value(), index)
//		if err != nil {
//			log.Errorf("unmarshal key-value from leveldb failed, %s:%s", ite.Key(), ite.Value())
//			continue
//		}
//		//m.bktree.Add(index)
//	}
//	return nil
//}

//func (m *MedictDBIndexer) buildTrie() error {
//	m.bktreeLock.Lock()
//	defer m.bktreeLock.Unlock()
//	if m.bktree != nil {
//		return nil
//	}
//	m.trie = &trie.RuneTrie{}
//	db, err := m.acquire()
//	if err != nil {
//		return err
//	}
//	defer m.release(db)
//
//	ite := db.NewIterator(util.BytesPrefix([]byte(prefixKeyword)), nil)
//	for ite.Next() {
//		index := new(model.MdictKeyWordIndex)
//		err = json.Unmarshal(ite.Value(), index)
//		if err != nil {
//			log.Errorf("unmarshal key-value from leveldb failed, %s:%s", ite.Key(), ite.Value())
//			continue
//		}
//		m.trie.Put(string(ite.Key()), ite.Value())
//	}
//	return nil
//}

func (st *searchTree) buildARTTree() error {
	st.lock.Lock()
	defer st.lock.Unlock()
	st.artTree = art.New()

	return nil

	//m.artTree.Insert(art.Key("Hi, I'm Key"), "Nice to meet you, I'm Value")
	//value, found := tree.Search(art.Key("Hi, I'm Key"))
	//if found {
	//	fmt.Printf("Search value=%v\n", value)
	//}
	//
	//tree.ForEach(func(node art.Node) bool {
	//	fmt.Printf("Callback value=%v\n", node.Value())
	//	return true
	//})
	//
	//for it := tree.Iterator(); it.HasNext(); {
	//	value, _ := it.Next()
	//	fmt.Printf("Iterator value=%v\n", value.Value())
	//}

}

// simSearch 近似搜索，搜索容忍区间内的词(<=tolerance)
// 使用 bktree + levenshtein 距离进行计算
//func (mdict *MedictDBIndexer) simSearch(word string, tolerance int) ([]*model.MdictKeyWordIndex, error) {
//	entry := &model.MdictKeyWordIndex{KeyWord: word}
//	results := mdict.bktree.Search(entry, tolerance, 16)
//
//	wrapper := &entryWrapper{
//		list: make([]*entryWrapperItem, 0),
//	}
//	for _, r := range results {
//		wrapper.list = append(wrapper.list, &entryWrapperItem{
//			entry:    r.Entry.(*model.MdictKeyWordIndex),
//			distance: r.Distance,
//		})
//	}
//
//	sort.Sort(wrapper)
//	return wrapper.toEntryList(), nil
//}

type entryWrapperItem struct {
	entry    *model.MdictKeyWordIndex
	distance int
}
