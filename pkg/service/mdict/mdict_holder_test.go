package mdict

import (
	"os"
	"testing"

	"github.com/terasum/medict/pkg/model"
)

func TestMdictHolder_BuildIndex(t *testing.T) {
	if _, err := os.Stat("./testdata/mdx/testdict.mdx"); err != nil {
		t.Skip("requires real mdx testdata: ./testdata/mdx/testdict.mdx")
	}
	holder, err := newMdictHolder("./testdata/mdx/testdict.mdx")
	if err != nil {
		t.Fatal(err)
	}

	err = holder.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}

	data, err := holder.Locate(&model.MdictKeyWordIndex{
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
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("data: %s", data)
}
