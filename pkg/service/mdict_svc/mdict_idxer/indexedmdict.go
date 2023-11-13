package mdict_idxer

//type IndexedMedict struct {
//	bktree *bktree.BKTree
//}
//
//func (mdict *IndexedMedict) Search(word string) ([]*MDictKeyBlockEntry, error) {
//	if mdict.bktree == nil {
//		return nil, errors.New("bktree hasn't build yet")
//	}
//	result, err := mdict.SimSearch(word, 1)
//	return result, err
//}
//
//func (mdict *MdictSvcImpl) BuildBKTree() error {
//	mdict.bktree = &bktree.BKTree{}
//	for _, e := range mdict.KeyBlockData.KeyEntries {
//		mdict.bktree.Add(e)
//	}
//	return nil
//}
//
//// SimSearch 近似搜索，搜索容忍区间内的词(<=tolerance)
//// 使用 bktree + levenshtein 距离进行计算
//func (mdict *MdictSvcImpl) SimSearch(word string, tolerance int) ([]*MDictKeyBlockEntry, error) {
//	entry := &MDictKeyBlockEntry{KeyWord: word}
//	results := mdict.bktree.Search(entry, tolerance, 100)
//
//	wrapper := &entryWrapper{
//		list: make([]*entryWrapperItem, 0),
//	}
//	for _, r := range results {
//		wrapper.list = append(wrapper.list, &entryWrapperItem{
//			entry:    r.Entry.(*MDictKeyBlockEntry),
//			distance: r.Distance,
//		})
//	}
//
//	sort.Sort(wrapper)
//	return wrapper.toEntryList(), nil
//}
//
//type entryWrapperItem struct {
//	entry    *MDictKeyBlockEntry
//	distance int
//}

//type entryWrapper struct {
//	list []*entryWrapperItem
//}
//
//func (w *entryWrapper) Less(i, j int) bool {
//	return w.list[i].distance < w.list[j].distance
//}
//
//func (w *entryWrapper) Len() int {
//	return len(w.list)
//}
//
//func (w *entryWrapper) Swap(i, j int) {
//	temp := w.list[i]
//	w.list[i] = w.list[j]
//	w.list[j] = temp
//}
//
//func (w *entryWrapper) toEntryList() []*MDictKeyBlockEntry {
//	entries := make([]*MDictKeyBlockEntry, w.Len())
//	for idx, result := range w.list {
//		entries[idx] = result.entry
//	}
//	return entries
//}
//
//// Distance calculates hamming distance.
//func (x *MDictKeyBlockEntry) Distance(e bktree.Entry) int {
//	a := x.KeyWord
//	b := e.(*MDictKeyBlockEntry).KeyWord
//	a = utils.StrToUnicode(a)
//	b = utils.StrToUnicode(b)
//
//	return levenshtein.Distance(a, b)
//}
