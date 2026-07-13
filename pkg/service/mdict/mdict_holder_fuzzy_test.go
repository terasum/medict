package mdict

import (
	"errors"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/terasum/medict/internal/libs/bktree"
	"github.com/terasum/medict/pkg/model"
	idxer "github.com/terasum/medict/pkg/service/mdict/mdict-idxer"
)

// fakeIdxer implements idxer.Indexer; its Search always misses, forcing the
// holder's Search onto the fuzzy fallback path. (interface defined in
// pkg/service/mdict/mdict-idxer/indexer.go)
type fakeIdxer struct{}

func (f *fakeIdxer) Lookup(keyword string) (*model.MdictKeyWordIndex, error) {
	return nil, errors.New("not found")
}
func (f *fakeIdxer) SetMeta(key, value string) error { return nil }
func (f *fakeIdxer) GetMeta(key string) (string, error) {
	return "", errors.New("not found")
}
func (f *fakeIdxer) AddRecord(record *model.MdictKeyWordIndex) error { return nil }
func (f *fakeIdxer) AddRecords(records []*model.MdictKeyWordIndex) error { return nil }
func (f *fakeIdxer) AllRecords() ([]*model.MdictKeyWordIndex, error)      { return nil, nil }
func (f *fakeIdxer) Search(keyword string) ([]*model.MdictKeyWordIndex, error) {
	return nil, errors.New("result not found")
}
func (f *fakeIdxer) Close() error { return nil }

// TestFuzzyFallback verifies that a prefix miss routes through the BK-tree
// fallback and returns the closest matches (sorted by distance).
func TestFuzzyFallback(t *testing.T) {
	mh := &mdictHolder{
		lock:       &sync.Mutex{},
		idxer:      &fakeIdxer{},
		bktree:     &bktree.BKTree{},
		bktreeDone: true,
	}
	mh.bktreeAdd(&model.MdictKeyWordIndex{KeyWord: "hello", RecordLocateStartOffset: 100})
	mh.bktreeAdd(&model.MdictKeyWordIndex{KeyWord: "help", RecordLocateStartOffset: 200})
	mh.bktreeAdd(&model.MdictKeyWordIndex{KeyWord: "hallo", RecordLocateStartOffset: 300})
	mh.bktreeAdd(&model.MdictKeyWordIndex{KeyWord: "world", RecordLocateStartOffset: 400})

	// "helo" is a typo of "hello" (distance 1) → fuzzy must surface it.
	res, err := mh.Search("helo")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(res), 1)

	// "world" is far (> tolerance 2) → must NOT appear.
	for _, r := range res {
		assert.NotEqual(t, "world", r.KeyWord, "world is beyond tolerance, should not match")
	}

	// The nearest hit ("hello", distance 1) must rank first.
	assert.Equal(t, "hello", res[0].KeyWord)

	// Results carry their offsets through (frontend Locate depends on these).
	assert.Equal(t, int64(100), res[0].RecordLocateStartOffset)

	// Distances non-decreasing (ascending).
	for i := 1; i < len(res); i++ {
		assert.GreaterOrEqual(t, res[i].ID, res[i-1].ID)
	}
}

// TestLookup_MissReturnsErrNotFound: a keyword miss returns the sentinel
// model.ErrNotFound (not a plain string error), so callers can errors.Is it
// (issue #732).
func TestLookup_MissReturnsErrNotFound(t *testing.T) {
	idx, err := idxer.NewIndexer(filepath.Join(t.TempDir(), "miss"))
	if err != nil {
		t.Fatal(err)
	}
	mh := &mdictHolder{lock: &sync.Mutex{}, idxer: idx}
	_, err = mh.Lookup("does-not-exist")
	if !assert.True(t, errors.Is(err, model.ErrNotFound), "want ErrNotFound, got %v", err) {
		t.FailNow()
	}
}

// TestIndexUpToDate: the schema migration trigger — a populated but old-schema
// (or schema-less) index is NOT up to date, forcing a rebuild; current schema +
// entries is up to date (issue #722 #1).
func TestIndexUpToDate(t *testing.T) {
	idx, err := idxer.NewIndexer(filepath.Join(t.TempDir(), "schema"))
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, indexUpToDate(idx), "empty index should not be up to date")

	assert.NoError(t, idx.SetMeta("entries_num", "100"))
	assert.False(t, indexUpToDate(idx), "old/missing schema should force rebuild")

	assert.NoError(t, idx.SetMeta("schema_version", indexSchemaVersion))
	assert.True(t, indexUpToDate(idx), "current schema + entries should skip rebuild")
}

// TestEnsureBkTree_FromIndexer: with bktreeDone=false (the .melev cache-hit
// case), ensureBkTree builds the BK-tree from the leveldb index — not by
// re-parsing the mdx — so fuzzy search still works (issue #722 #2).
func TestEnsureBkTree_FromIndexer(t *testing.T) {
	idx, err := idxer.NewIndexer(filepath.Join(t.TempDir(), "bktree"))
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, idx.AddRecords([]*model.MdictKeyWordIndex{
		{KeyWord: "hello", RecordLocateStartOffset: 100},
		{KeyWord: "help", RecordLocateStartOffset: 200},
	}))

	mh := &mdictHolder{lock: &sync.Mutex{}, idxer: idx} // bktreeDone defaults to false
	assert.NoError(t, mh.ensureBkTree())
	assert.True(t, mh.bktreeDone)

	// "helo" (typo of hello) resolves via the indexer-built BK-tree.
	res, err := mh.Search("helo")
	assert.NoError(t, err)
	if assert.NotEmpty(t, res) {
		assert.Equal(t, "hello", res[0].KeyWord)
	}
}

// TestFuzzyRanking builds a small BK-tree directly and asserts distance ordering
// and that an exact hit has distance 0.
func TestFuzzyRanking(t *testing.T) {
	tree := &bktree.BKTree{}
	words := []struct {
		word string
		off  int64
	}{
		{"hello", 100},
		{"hell", 200},
		{"help", 300},
		{"heloo", 400},
		{"hallo", 600},
	}
	for _, w := range words {
		tree.Add(&fuzzyEntry{&model.MdictKeyWordIndex{KeyWord: w.word, RecordLocateStartOffset: w.off}})
	}

	raw := tree.Search(&fuzzyEntry{&model.MdictKeyWordIndex{KeyWord: "hello"}}, fuzzyTolerance, fuzzyLimit)
	assert.GreaterOrEqual(t, len(raw), 4)

	// 必含精确匹配 hello (distance 0)，且 offset 保留（BK-tree 不保证顺序，遍历查找）
	var exact *fuzzyEntry
	for _, r := range raw {
		if r.Distance == 0 {
			exact = r.Entry.(*fuzzyEntry)
		}
	}
	assert.NotNil(t, exact, "应含 distance 0 的精确匹配")
	assert.Equal(t, "hello", exact.KeyWord)
	assert.Equal(t, int64(100), exact.RecordLocateStartOffset, "offset 应保留")

	// 所有结果 distance <= tolerance
	for _, r := range raw {
		assert.LessOrEqual(t, r.Distance, fuzzyTolerance)
	}
}

// TestFuzzyTypoCorrection: a single-edit typo resolves to the target word.
func TestFuzzyTypoCorrection(t *testing.T) {
	mh := &mdictHolder{
		lock:       &sync.Mutex{},
		idxer:      &fakeIdxer{},
		bktree:     &bktree.BKTree{},
		bktreeDone: true,
	}
	mh.bktreeAdd(&model.MdictKeyWordIndex{KeyWord: "hello", RecordLocateStartOffset: 100})

	res, err := mh.Search("helo") // missing one 'l'
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "hello", res[0].KeyWord)
}

// TestFuzzyNoMatch: when neither prefix nor fuzzy matches, Search returns an
// empty result with no error (empty != error, issue #722 P3).
func TestFuzzyNoMatch(t *testing.T) {
	mh := &mdictHolder{
		lock:       &sync.Mutex{},
		idxer:      &fakeIdxer{},
		bktree:     &bktree.BKTree{},
		bktreeDone: true,
	}
	mh.bktreeAdd(&model.MdictKeyWordIndex{KeyWord: "hello"})

	res, err := mh.Search("zzzzzzzz") // far from everything
	assert.NoError(t, err)
	assert.Empty(t, res)
}
