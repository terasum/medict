package mdict

import (
	"errors"
	gomdict "github.com/terasum/medict/internal/libs/go-mdict"
	"github.com/terasum/medict/pkg/model"
	idxer "github.com/terasum/medict/pkg/service/mdict/mdict-idxer"
	"strconv"
	"sync"
)

type mdictHolder struct {
	lock         *sync.Mutex
	idxFilePath  string
	dictFilePath string
	idxer        idxer.Indexer
	rawdict      *gomdict.Mdict
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

	log.Infof("holder %+v", index.RecordBlock)

	def, err := mh.rawdict.LocateByKeywordIndex(index)
	if err != nil {
		log.Errorf("locate error %s", err.Error())
		return nil, err
	}
	return def, nil
}

func (mh *mdictHolder) Lookup(keyword string) ([]byte, error) {
	entry, err := mh.idxer.Lookup(keyword)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, errors.New("not found")
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

func (mh *mdictHolder) BuildIndex() error {
	err := mh.rawdict.BuildIndex()
	if err != nil {
		return err
	}
	// has already built
	value, err := mh.idxer.GetMeta("entries_num")
	if err == nil && value != "" {
		num, err1 := strconv.ParseInt(value, 10, 64)
		if err1 == nil {
			log.Infof("index has already built, entries number is %d", num)
			return nil
		}
	}

	err = mh.idxer.SetMeta("Title", mh.rawdict.Title())
	if err != nil {
		return err
	}
	err = mh.idxer.SetMeta("Description", mh.rawdict.Description())
	if err != nil {
		return err
	}
	err = mh.idxer.SetMeta("CreationDate", mh.rawdict.CreationDate())
	if err != nil {
		return err
	}
	err = mh.idxer.SetMeta("GenerateEngineVersion", mh.rawdict.GeneratedByEngineVersion())
	if err != nil {
		return err
	}
	err = mh.idxer.SetMeta("filepath", mh.dictFilePath)
	if err != nil {
		return err
	}
	err = mh.idxer.SetMeta("idx_filepath", mh.idxFilePath)
	if err != nil {
		return err
	}
	err = mh.idxer.SetMeta("is_utf16", strconv.FormatBool(mh.rawdict.IsUTF16()))
	if err != nil {
		return err
	}
	err = mh.idxer.SetMeta("is_mdd", strconv.FormatBool(mh.rawdict.IsMDD()))
	if err != nil {
		return err
	}
	err = mh.idxer.SetMeta("is_record_encrypt", strconv.FormatBool(mh.rawdict.IsRecordEncrypted()))
	if err != nil {
		return err
	}

	entries, err := mh.rawdict.GetKeyWordEntries()
	if err != nil {
		return err
	}
	for _, entry := range entries {
		idx, err1 := mh.ConvertKeyWordIndex(entry)
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

func (mh *mdictHolder) Title() string {
	value, err := mh.idxer.GetMeta("Title")
	if err == nil {
		return value
	}
	return ""

}

func (mh *mdictHolder) Description() string {
	value, err := mh.idxer.GetMeta("Description")
	if err == nil {
		return value
	}
	return ""

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
	result, err := mh.idxer.Search(keyword)
	if err != nil {
		return nil, err
	}
	for idx, re := range result {
		re.ID = idx
	}
	return result, nil
}
