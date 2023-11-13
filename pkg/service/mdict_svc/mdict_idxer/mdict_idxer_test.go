package mdict_idxer

import (
	"github.com/terasum/medict/pkg/model"
	"testing"
)

func TestIndexEngine_AddRecord(t *testing.T) {
	indexFilePath := "./testdata/testidx.meidx"

	engine, err := NewIdxer(indexFilePath)
	if err != nil {
		t.Fatal(err)
	}

	err = engine.SetMeta("id", "1234")
	if err != nil {
		t.Fatal(err)
	}

	value, err := engine.GetMeta("id")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("value id : %s", value)

	record := &model.MdictKeyWordIndex{
		ID:                            0,
		KeyWord:                       "$100, £50, etc. a throw",
		RecordLocateStartOffset:       0,
		RecordLocateEndOffset:         1077,
		IsUTF16:                       0,
		IsRecordEncrypt:               0,
		RecordBlockDataStartOffset:    833968,
		RecordBlockDataCompressSize:   10119,
		RecordBlockDataDeCompressSize: 65004,
		KeyWordDataStartOffset:        0,
		KeyWordDataEndOffset:          1077,
	}

	if err = engine.AddRecord(record); err != nil {
		t.Errorf("AddRecord() error = %v", err)
	}
}
