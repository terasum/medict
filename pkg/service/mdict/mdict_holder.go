package mdict

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"

	"github.com/creasty/go-levenshtein"
	"github.com/terasum/medict/internal/libs/bktree"
	gomdict "github.com/terasum/medict/internal/libs/go-mdict"
	"github.com/terasum/medict/pkg/model"
	idxer "github.com/terasum/medict/pkg/service/mdict/mdict-idxer"
)

type mdictHolder struct {
	lock         *sync.Mutex
	idxFilePath  string
	dictFilePath string
	idxer        idxer.Indexer
	rawdict      *gomdict.Mdict

	bktree     *bktree.BKTree // 模糊搜索用;BuildIndex 或首次模糊查询时构建
	bktreeDone bool           // BK-tree 是否已构建(lazy 守卫)
}

const (
	fuzzyTolerance = 2 // Levenshtein 容差：覆盖 1-2 字符拼写差异
	fuzzyLimit     = 100
)

// indexSchemaVersion tags the .melev index format. Bumping it forces a one-time
// rebuild on existing installs (and repairs legacy collision-broken indexes) —
// issue #722 #1 (headword-only key -> offset-suffix key).
const indexSchemaVersion = "2"

// fuzzyEntry 包装 *MdictKeyWordIndex 用于 BK-tree 模糊搜索。
// 不复用 MdictKeyWordIndex.Distance：后者基于 utils.StrToUnicode 转义串（每个 rune 展开为 6 字符的 \uXXXX），
// 会把 1 个 rune 的差异放大成 6 个字符的 Levenshtein 距离，使 tolerance 语义失效。
// 这里直接在原始 KeyWord 上算字符级 Levenshtein（与 stardict 的 bkString 一致）。
type fuzzyEntry struct {
	*model.MdictKeyWordIndex
}

func (e *fuzzyEntry) Distance(other bktree.Entry) int {
	return levenshtein.Distance(e.KeyWord, other.(*fuzzyEntry).KeyWord)
}

func newMdictHolder(filePath string) (*mdictHolder, error) {
	indexFile := filePath + ".melev"

	dict, err := gomdict.New(filePath)
	if err != nil {
		return nil, err
	}

	indexer, err := idxer.NewIndexer(indexFile)
	if err != nil {
		return nil, err
	}
	return &mdictHolder{
		lock:         new(sync.Mutex),
		idxFilePath:  indexFile,
		dictFilePath: filePath,
		idxer:        indexer,
		rawdict:      dict,
	}, nil
}

func (mh *mdictHolder) ConvertKeyWordIndex(entry *gomdict.MDictKeywordEntry) (*model.MdictKeyWordIndex, error) {
	index, err1 := mh.rawdict.KeywordEntryToIndex(entry)
	if err1 != nil {
		return nil, err1
	}

	isUtf16 := 0
	isRecordEncrypt := 0
	isMdd := 0

	if mh.rawdict.IsUTF16() {
		isUtf16 = 1
	}
	if mh.rawdict.IsRecordEncrypted() {
		isRecordEncrypt = 1
	}
	if mh.rawdict.IsMDD() {
		isMdd = 1
	}

	return &model.MdictKeyWordIndex{
		KeyWord:                       index.KeywordEntry.KeyWord,
		RecordLocateStartOffset:       index.KeywordEntry.RecordStartOffset,
		RecordLocateEndOffset:         index.KeywordEntry.RecordEndOffset,
		IsUTF16:                       isUtf16,
		IsRecordEncrypt:               isRecordEncrypt,
		IsMDD:                         isMdd,
		RecordBlockDataStartOffset:    index.RecordBlock.DataStartOffset,
		RecordBlockDataCompressSize:   index.RecordBlock.CompressSize,
		RecordBlockDataDeCompressSize: index.RecordBlock.DeCompressSize,
		KeyWordDataStartOffset:        index.RecordBlock.KeyWordPartStartOffset,
		KeyWordDataEndOffset:          index.RecordBlock.KeyWordPartDataEndOffset,
	}, nil
}

func (mh *mdictHolder) Locate(entry *model.MdictKeyWordIndex) ([]byte, error) {
	index := &gomdict.MDictKeywordIndex{
		KeywordEntry: gomdict.MDictKeywordEntry{
			RecordStartOffset: entry.RecordLocateStartOffset,
			RecordEndOffset:   entry.RecordLocateEndOffset,
			KeyWord:           entry.KeyWord,
			KeyBlockIdx:       0,
		},
		RecordBlock: gomdict.MDictKeywordIndexRecordBlock{
			DataStartOffset:          entry.RecordBlockDataStartOffset,
			CompressSize:             entry.RecordBlockDataCompressSize,
			DeCompressSize:           entry.RecordBlockDataDeCompressSize,
			KeyWordPartStartOffset:   entry.KeyWordDataStartOffset,
			KeyWordPartDataEndOffset: entry.KeyWordDataEndOffset,
		},
	}

	log.Debugf("holder %+v", index.RecordBlock)

	def, err := mh.rawdict.LocateByKeywordIndex(index)
	if err != nil {
		log.Errorf("locate error %s", err.Error())
		return nil, err
	}
	return def, nil
}

// Close releases the indexer's leveldb handle. (The raw .mdx parser handle is
// not released here — it exposes no Close; tracked separately.)
func (mh *mdictHolder) Close() error {
	if mh.idxer != nil {
		return mh.idxer.Close()
	}
	return nil
}

func (mh *mdictHolder) Lookup(keyword string) ([]byte, error) {
	entry, err := mh.idxer.Lookup(keyword)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, model.ErrNotFound
	}
	index := &gomdict.MDictKeywordIndex{
		KeywordEntry: gomdict.MDictKeywordEntry{
			RecordStartOffset: entry.RecordLocateStartOffset,
			RecordEndOffset:   entry.RecordLocateEndOffset,
			KeyWord:           entry.KeyWord,
			KeyBlockIdx:       0,
		},
		RecordBlock: gomdict.MDictKeywordIndexRecordBlock{
			DataStartOffset:          entry.RecordBlockDataStartOffset,
			CompressSize:             entry.RecordBlockDataCompressSize,
			DeCompressSize:           entry.RecordBlockDataDeCompressSize,
			KeyWordPartStartOffset:   entry.KeyWordDataStartOffset,
			KeyWordPartDataEndOffset: entry.KeyWordDataEndOffset,
		},
	}

	return mh.rawdict.LocateByKeywordIndex(index)
}

// indexUpToDate reports whether the persisted .melev index is populated AND at
// the current schema, so BuildIndex can skip the rebuild. Exposed as a helper
// so the migration trigger is unit-testable without a source dict (issue #722 #1).
func indexUpToDate(idx idxer.Indexer) bool {
	v, err := idx.GetMeta("entries_num")
	if err != nil || v == "" {
		return false
	}
	sv, _ := idx.GetMeta("schema_version")
	return sv == indexSchemaVersion
}

func (mh *mdictHolder) BuildIndex() error {
	err := mh.rawdict.BuildIndex()
	if err != nil {
		return err
	}
	// Skip the rebuild only if the index is populated AND at the current schema.
	// A schema bump forces a one-time rebuild, which also repairs legacy indexes
	// built with the collision-prone headword-only key (issue #722 #1).
	if indexUpToDate(mh.idxer) {
		if v, _ := mh.idxer.GetMeta("entries_num"); v != "" {
			log.Infof("index already built (schema %s), entries: %s", indexSchemaVersion, v)
		}
		// #781: Even when the leveldb index is cached, the BK-tree is in-memory
		// only and must be rebuilt on every app launch. Try loading the persisted
		// BK-tree cache (instant); if not found, build asynchronously.
		mh.lock.Lock()
		loaded := mh.loadBkTreeCache()
		mh.lock.Unlock()
		if !loaded {
			mh.startAsyncBkTree()
		} else {
			log.Infof("BK-tree loaded from cache (bktree.gob)")
		}
		return nil
	}

	metas := []struct{ k, v string }{
		{"Title", mh.rawdict.Title()},
		{"Description", mh.rawdict.Description()},
		{"CreationDate", mh.rawdict.CreationDate()},
		{"GenerateEngineVersion", mh.rawdict.GeneratedByEngineVersion()},
		{"filepath", mh.dictFilePath},
		{"idx_filepath", mh.idxFilePath},
		{"is_utf16", strconv.FormatBool(mh.rawdict.IsUTF16())},
		{"is_mdd", strconv.FormatBool(mh.rawdict.IsMDD())},
		{"is_record_encrypt", strconv.FormatBool(mh.rawdict.IsRecordEncrypted())},
	}
	for _, m := range metas {
		if err = mh.idxer.SetMeta(m.k, m.v); err != nil {
			return err
		}
	}

	entries, err := mh.rawdict.GetKeyWordEntries()
	if err != nil {
		return err
	}
	records := make([]*model.MdictKeyWordIndex, 0, len(entries))
	for _, entry := range entries {
		idx, err1 := mh.ConvertKeyWordIndex(entry)
		if err1 != nil {
			log.Error(err1.Error())
			continue
		}
		records = append(records, idx)
	}

	// Bulk-write all keyword records in a single leveldb batch — far faster
	// than one Put per record during a full index build.
	if err := mh.idxer.AddRecords(records); err != nil {
		return err
	}

	if err := mh.idxer.SetMeta("schema_version", indexSchemaVersion); err != nil {
		return err
	}
	err = mh.idxer.SetMeta("entries_num", strconv.FormatInt(mh.rawdict.GetKeyWordEntriesSize(), 10))
	if err != nil {
		return err
	}

	// #781: Build the BK-tree asynchronously. The eager in-loop build consumed
	// ~125s for a 195K-entry dictionary (92% of BuildIndex wall time), blocking
	// the UI ("卡死"). The BK-tree is only needed for fuzzy fallback on Search
	// misses — defer it to a background goroutine. ensureBkTree() guards lazy
	// access: Search misses skip fuzzy if the tree isn't ready yet.
	recordsSnapshot := records // capture for goroutine
	go func() {
		mh.lock.Lock()
		defer mh.lock.Unlock()
		if mh.bktreeDone {
			return
		}
		tree := &bktree.BKTree{}
		for _, idx := range recordsSnapshot {
			tree.Add(&fuzzyEntry{idx})
		}
		mh.bktree = tree
		mh.bktreeDone = true
		log.Infof("BK-tree built asynchronously: %d entries", len(recordsSnapshot))
		mh.saveBkTreeCache() // #781: persist for instant load on next launch
	}()

	return nil
}

// bktreeCachePath returns the .bktree file path inside the .melev directory.
func (mh *mdictHolder) bktreeCachePath() string {
	return filepath.Join(mh.idxFilePath, "bktree.gob")
}

// saveBkTreeCache serializes the BK-tree to bktree.gob inside the .melev dir.
// Called after async build completes.
func (mh *mdictHolder) saveBkTreeCache() {
	tree := mh.bktree
	if tree == nil {
		return
	}
	path := mh.bktreeCachePath()
	f, err := os.Create(path)
	if err != nil {
		log.Errorf("save BK-tree cache: %v", err)
		return
	}
	defer f.Close()
	marshalEntry := func(e bktree.Entry) ([]byte, error) {
		fe := e.(*fuzzyEntry)
		return json.Marshal(fe.MdictKeyWordIndex)
	}
	if err := tree.Save(f, marshalEntry); err != nil {
		log.Errorf("encode BK-tree cache: %v", err)
	}
}

// loadBkTreeCache tries to load the BK-tree from bktree.gob. Returns true if loaded.
func (mh *mdictHolder) loadBkTreeCache() bool {
	path := mh.bktreeCachePath()
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	unmarshalEntry := func(data []byte) (bktree.Entry, error) {
		var idx model.MdictKeyWordIndex
		if err := json.Unmarshal(data, &idx); err != nil {
			return nil, err
		}
		return &fuzzyEntry{&idx}, nil
	}
	tree, err := bktree.Load(f, unmarshalEntry)
	if err != nil {
		log.Errorf("load BK-tree cache: %v", err)
		return false
	}
	mh.bktree = tree
	mh.bktreeDone = true
	return true
}

// startAsyncBkTree builds the BK-tree in a background goroutine by reading
// records from the leveldb index (AllRecords). Used when BuildIndex is skipped
// because the leveldb index is already cached — the BK-tree is in-memory only
// and must be rebuilt on every app launch. Non-blocking; Search misses skip
// fuzzy if the tree isn't ready yet.
func (mh *mdictHolder) startAsyncBkTree() {
	go func() {
		mh.lock.Lock()
		defer mh.lock.Unlock()
		if mh.bktreeDone {
			return
		}
		records, err := mh.idxer.AllRecords()
		if err != nil {
			log.Errorf("async BK-tree build: AllRecords failed: %v", err)
			return
		}
		tree := &bktree.BKTree{}
		for _, r := range records {
			tree.Add(&fuzzyEntry{r})
		}
		mh.bktree = tree
		mh.bktreeDone = true
		log.Infof("BK-tree built asynchronously from cache: %d entries", len(records))
		mh.saveBkTreeCache() // persist for instant load on next launch
	}()
}

// bktreeAdd 向模糊搜索 BK-tree 添加一个词项（BuildIndex 时逐条调用）。
func (mh *mdictHolder) bktreeAdd(e *model.MdictKeyWordIndex) {
	if mh.bktree == nil {
		mh.bktree = &bktree.BKTree{}
	}
	mh.bktree.Add(&fuzzyEntry{e})
}

// ensureBkTree 按需从 leveldb 索引构建 BK-tree（不再重解 mdx），兜底 .melev 缓存
// 命中、跳过 BuildIndex 的场景：用户也能拿到模糊搜索，且不付全量 mdx 解析代价
// (issue #722 #2)。由 mh.lock 保护；eager 路径已在 BuildIndex 构建并置 bktreeDone，
// 本函数对该路径是 no-op。
func (mh *mdictHolder) ensureBkTree() error {
	mh.lock.Lock()
	defer mh.lock.Unlock()
	if mh.bktreeDone {
		return nil
	}
	records, err := mh.idxer.AllRecords()
	if err != nil {
		return err
	}
	tree := &bktree.BKTree{}
	for _, r := range records {
		tree.Add(&fuzzyEntry{r})
	}
	mh.bktree = tree
	mh.bktreeDone = true
	return nil
}

// Title returns the dictionary title. It prefers the leveldb meta (written
// during BuildIndex) but falls back to the mdx header — which is parsed at load
// (readDictHeader) — so the real title is available BEFORE indexing. Without
// this fallback the dict list shows the filename, because GetMeta("Title") is
// empty until BuildIndex runs (#782).
func (mh *mdictHolder) Title() string {
	if v, err := mh.idxer.GetMeta("Title"); err == nil && v != "" {
		return v
	}
	return mh.rawdict.Title()
}

func (mh *mdictHolder) Description() string {
	if v, err := mh.idxer.GetMeta("Description"); err == nil && v != "" {
		return v
	}
	return mh.rawdict.Description()
}

func (mh *mdictHolder) CreationDate() string {
	value, err := mh.idxer.GetMeta("CreationDate")
	if err == nil {
		return value
	}
	return ""

}

func (mh *mdictHolder) GenerateEngineVersion() string {
	value, err := mh.idxer.GetMeta("GenerateEngineVersion")
	if err == nil {
		return value
	}
	return ""

}

func (mh *mdictHolder) Search(keyword string) ([]*model.MdictKeyWordIndex, error) {
	// 精确前缀优先：命中直接返回
	if result, err := mh.idxer.Search(keyword); err == nil && len(result) > 0 {
		for idx, re := range result {
			re.ID = idx
		}
		return result, nil
	}
	// 前缀为空（拼错/记不清词头）→ BK-tree 模糊兜底
	// #781: 如果 BK-tree 还在后台构建中,跳过 fuzzy(返回空),不阻塞用户。
	// 后台构建完成后自动生效,后续 miss 正常返回 fuzzy 结果。
	mh.lock.Lock()
	treeReady := mh.bktreeDone && mh.bktree != nil
	tree := mh.bktree
	mh.lock.Unlock()
	if !treeReady {
		return nil, nil // BK-tree 构建中,暂时跳过 fuzzy
	}
	// tree pointer is stable once bktreeDone=true (async goroutine sets both
	// atomically under mh.lock). Safe to read without lock for Search (read-only).
	needle := &fuzzyEntry{&model.MdictKeyWordIndex{KeyWord: keyword}}
	raw := tree.Search(needle, fuzzyTolerance, fuzzyLimit)
	if len(raw) == 0 {
		return nil, nil
	}
	sort.Slice(raw, func(i, j int) bool { return raw[i].Distance < raw[j].Distance })
	out := make([]*model.MdictKeyWordIndex, 0, len(raw))
	for i, r := range raw {
		e := r.Entry.(*fuzzyEntry).MdictKeyWordIndex
		e.ID = i
		out = append(out, e)
	}
	return out, nil
}
