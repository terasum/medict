package stardict

import (
	"fmt"
	"github.com/creasty/go-levenshtein"
	"github.com/terasum/medict/internal/libs/bktree"
	"github.com/terasum/medict/internal/libs/go-stardict/stardict"
	"github.com/terasum/medict/pkg/model"
)

var _ model.GeneralDictionary = &StarDict{}

/**
 * StarDict support stardict file type
 */
type StarDict struct {
	DzFilePath   string
	DictFilePath string
	IfoFilePath  string
	IdxFilePath  string
	SDict        *stardict.StarDict
	BkTree       *bktree.BKTree
	ready        bool
}

func (s *StarDict) Name() string {
	return s.SDict.GetBookName()
}

func (s *StarDict) Description() *model.PlainDictionaryInfo {
	return &model.PlainDictionaryInfo{
		Title:                 s.SDict.GetBookName(),
		Description:           fmt.Sprintf("StarDict Dictionary"),
		CreateDate:            "",
		GenerateEngineVersion: "",
	}
}

type bkString struct {
	w string
}

func (bs *bkString) Distance(entry bktree.Entry) int {
	return levenshtein.Distance(bs.w, entry.(*bkString).w)
}

// Close is a no-op for StarDict: it holds no leveldb/database handles, only
// in-memory structures (BKTree) and the parsed stardict data.
func (s *StarDict) Close() error {
	return nil
}

func (s *StarDict) BuildIndex() error {
	if s.ready {
		return nil
	}
	err := s.SDict.BuildIndex()
	if err != nil {
		return err
	}

	s.BkTree = &bktree.BKTree{}
	words := s.SDict.KeyList()
	for _, word := range words {
		s.BkTree.Add(&bkString{w: word})
	}
	s.ready = true
	return nil

}

func (s *StarDict) DictType() model.DictType {
	return model.DictTypeStarDict
}

// Lookup returns the definition HTML for keyword. The underlying lib renders a
// "Nothing Found" page on a miss; a not-yet-built dictionary surfaces as
// ErrNotFound rather than an empty body (#739).
func (s *StarDict) Lookup(keyword string) ([]byte, error) {
	if !s.ready {
		return nil, model.ErrNotFound
	}
	return []byte(s.SDict.Lookup(keyword)), nil
}

// LookupResource: stardict resources (css/images/audio under the <name>.res
// folder) are not loaded by the underlying lib today, so there is nothing to
// serve. Return ErrNotFound so the request fails cleanly instead of an empty
// 200 (#739). Full .res/ loading is a future feature.
func (s *StarDict) LookupResource(keyword string) ([]byte, error) {
	return nil, model.ErrNotFound
}

func (s *StarDict) Locate(entry *model.KeyQueryIndex) ([]byte, error) {
	return s.Lookup(entry.KeyWord)
}

func (s *StarDict) Search(keyword string) ([]*model.KeyQueryIndex, error) {
	if !s.ready {
		return nil, fmt.Errorf("dictionary has not built yet")
	}
	res := s.BkTree.Search(&bkString{w: keyword}, 1, 100)
	result := make([]*model.KeyQueryIndex, 0)
	for id, r := range res {
		result = append(result, &model.KeyQueryIndex{
			IndexType: model.IndexTypeStardict,
			MdictKeyWordIndex: &model.MdictKeyWordIndex{
				ID:                      id,
				RecordLocateStartOffset: 0,
				RecordLocateEndOffset:   0,
				KeyWord:                 r.Entry.(*bkString).w,
			},
		})
	}
	return result, nil
}

func NewStardict(dirItem *model.DirItem) (model.GeneralDictionary, error) {
	sdict, err := stardict.NewStarDict(&stardict.Option{
		DictDzFilePath: dirItem.StarDictDzAbsPath,
		DictFilePath:   dirItem.StarDictDzAbsPath,
		IdxFilePath:    dirItem.StarDictIdxAbsPath,
		IfoFilePath:    dirItem.StarDictIfoAbsPath,
	})
	if err != nil {
		return nil, err
	}

	return &StarDict{
		DzFilePath:   dirItem.StarDictDzAbsPath,
		DictFilePath: dirItem.StarDictAbsPath,
		IfoFilePath:  dirItem.StarDictIfoAbsPath,
		IdxFilePath:  dirItem.StarDictIdxAbsPath,
		SDict:        sdict,
	}, nil

}
