package mdict_svc

import (
	"github.com/terasum/medict/internal/gomdict"
	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service/mdict_svc/mdict_idxer"
	"strconv"
	"sync"
)

type mdictHolder struct {
	lock          *sync.Mutex
	idxFilePath   string
	dictFilePath  string
	idxer         *mdict_idxer.MdictIdxer
	rawdict       *gomdict.Mdict
	hasBuildIndex bool
}

func newMdictHolder(filePath string) (*mdictHolder, error) {
	indexFile := filePath + ".meidx"

	dict, err := gomdict.New(filePath)
	if err != nil {
		return nil, err
	}

	idxer, err := mdict_idxer.NewIdxer(indexFile)
	if err != nil {
		return nil, err
	}
	return &mdictHolder{
		lock:         new(sync.Mutex),
		idxFilePath:  indexFile,
		dictFilePath: filePath,
		idxer:        idxer,
		rawdict:      dict,
	}, nil
}

func (mh *mdictHolder) BuildIndex() error {
	err := mh.rawdict.BuildIndex()
	if err != nil {
		return err
	}
	err = mh.idxer.SetMeta("title", mh.rawdict.Title())
	err = mh.idxer.SetMeta("description", mh.rawdict.Description())
	err = mh.idxer.SetMeta("filepath", mh.dictFilePath)
	err = mh.idxer.SetMeta("idx_filepath", mh.idxFilePath)
	err = mh.idxer.SetMeta("is_utf16", strconv.FormatBool(mh.rawdict.IsUTF16()))
	err = mh.idxer.SetMeta("is_mdd", strconv.FormatBool(mh.rawdict.IsMDD()))
	err = mh.idxer.SetMeta("is_record_encrypt", strconv.FormatBool(mh.rawdict.IsRecordEncrypted()))

	entries, err := mh.rawdict.GetKeyWordEntries()
	if err != nil {
		return err
	}
	for _, entry := range entries {
		idx, err1 := mh.convertKeyWordIndex(entry)
		if err1 != nil {
			log.Error(err1.Error())
			continue
		}
		err1 = mh.idxer.AddRecord(idx)
		if err1 != nil {
			log.Error(err1.Error())
			continue
		}
	}

	err = mh.idxer.SetMeta("entries_num", strconv.FormatInt(mh.rawdict.GetKeyWordEntriesSize(), 10))
	return nil
}

func (mh *mdictHolder) convertKeyWordIndex(entry *gomdict.MDictKeywordEntry) (*model.MdictKeyWordIndex, error) {
	index, err1 := mh.rawdict.KeywordEntryToIndex(entry)
	if err1 != nil {
		return nil, err1
	}

	is_utf16 := 0
	is_record_encrypt := 0
	is_mdd := 0

	if mh.rawdict.IsUTF16() {
		is_utf16 = 1
	}
	if mh.rawdict.IsRecordEncrypted() {
		is_record_encrypt = 1
	}
	if mh.rawdict.IsMDD() {
		is_mdd = 1
	}

	return &model.MdictKeyWordIndex{
		KeyWord:                       index.KeywordEntry.KeyWord,
		RecordLocateStartOffset:       index.KeywordEntry.RecordStartOffset,
		RecordLocateEndOffset:         index.KeywordEntry.RecordEndOffset,
		IsUTF16:                       is_utf16,
		IsRecordEncrypt:               is_record_encrypt,
		IsMDD:                         is_mdd,
		RecordBlockDataStartOffset:    index.RecordBlock.DataStartOffset,
		RecordBlockDataCompressSize:   index.RecordBlock.CompressSize,
		RecordBlockDataDeCompressSize: index.RecordBlock.DeCompressSize,
		KeyWordDataStartOffset:        index.RecordBlock.KeyWordPartStartOffset,
		KeyWordDataEndOffset:          index.RecordBlock.KeyWordPartDataEndOffset,
	}, nil
}

func (mh *mdictHolder) Locate(entry *model.KeyQueryIndex) ([]byte, error) {
	//if !mh.hasBuildIndex {
	//	return nil, errors.New("dictionary not ready, building index first")
	//}

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
