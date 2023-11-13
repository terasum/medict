package mdict_svc

import (
	"github.com/terasum/medict/pkg/model"
	"testing"
)

func TestMdictHolder_BuildIndex(t *testing.T) {
	holder, err := newMdictHolder("./testdata/mdx/testdict.mdx")
	if err != nil {
		t.Fatal(err)
	}

	err = holder.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}

	data, err := holder.Locate(&model.KeyQueryIndex{
		IndexType: "mdx",
		MdictKeyWordIndex: &model.MdictKeyWordIndex{
			ID:                            0,
			KeyWord:                       "accessorized",
			RecordLocateStartOffset:       735416,
			RecordLocateEndOffset:         735438,
			IsUTF16:                       0,
			IsRecordEncrypt:               0,
			IsMDD:                         0,
			RecordBlockDataStartOffset:    920434,
			RecordBlockDataCompressSize:   8573,
			RecordBlockDataDeCompressSize: 64463,
			KeyWordDataStartOffset:        35367,
			KeyWordDataEndOffset:          35389,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("data: %s", data)
}
