//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package go_mdict

type MdictType int

const (
	MdictTypeMdd MdictType = 1
	MdictTypeMdx MdictType = 2

	EncryptNoEnc      = 0
	EncryptRecordEnc  = 1
	EncryptKeyInfoEnc = 2
	NumfmtBe8bytesq   = 0
	NumfmtBe4bytesi   = 1
	EncodingUtf8      = 0
	EncodingUtf16     = 1
	EncodingBig5      = 2
	ENCODING_GBK      = 3
	ENCODING_GB2312   = 4
	EncodingGb18030   = 5
)

type MdictBase struct {
	filePath string
	fileType MdictType
	meta     *mdictMeta

	header       *mdictHeader
	keyBlockMeta *mdictKeyBlockMeta
	keyBlockInfo *mdictKeyBlockInfo
	keyBlockData *mdictKeyBlockData

	recordBlockMeta *mdictRecordBlockMeta
	recordBlockInfo *mdictRecordBlockInfo
	//RecordBlockData *MDictRecordBlockData

	rangeTreeRoot *RecordBlockRangeTreeNode
}

/********************************
 *    private data type          *
 ********************************/
type mdictHeader struct {
	headerBytesSize          uint32
	headerInfoBytes          []byte
	headerInfo               string
	adler32Checksum          uint32
	dictionaryHeaderByteSize int64
}

type mdictMeta struct {
	encryptType  int
	version      float32
	numberWidth  int
	numberFormat int
	encoding     int

	// key-block part bytes start offset in the mdx/mdd file
	keyBlockMetaStartOffset int64

	description              string
	title                    string
	creationDate             string
	generatedByEngineVersion string
}

type mdictKeyBlockMeta struct {
	// keyBlockNum key block number size
	keyBlockNum int64
	// entriesNums entries number size
	entriesNum int64
	// key-block information size (decompressed)
	keyBlockInfoDecompressSize int64
	// key-block information size (compressed)
	keyBlockInfoCompressedSize int64
	// key-block Data Size (decompressed)
	keyBlockDataTotalSize int64
	// key-block information start position in the mdx/mdd file
	keyBlockInfoStartOffset int64
}

type mdictKeyBlockInfo struct {
	keyBlockEntriesStartOffset int64
	keyBlockInfoList           []*mdictKeyBlockInfoItem
}

type mdictKeyBlockInfoItem struct {
	firstKey                      string
	firstKeySize                  int
	lastKey                       string
	lastKeySize                   int
	keyBlockInfoIndex             int
	keyBlockCompressSize          int64
	keyBlockCompAccumulator       int64
	keyBlockDeCompressSize        int64
	keyBlockDeCompressAccumulator int64
}

type mdictKeyBlockData struct {
	keyEntries                 []*MDictKeywordEntry
	keyEntriesSize             int64
	recordBlockMetaStartOffset int64
}

type mdictRecordBlockMeta struct {
	keyRecordMetaStartOffset int64
	keyRecordMetaEndOffset   int64

	recordBlockNum          int64
	entriesNum              int64
	recordBlockInfoCompSize int64
	recordBlockCompSize     int64
}
type mdictRecordBlockInfo struct {
	recordInfoList             []*MdictRecordBlockInfoListItem
	recordBlockInfoStartOffset int64
	recordBlockInfoEndOffset   int64
	recordBlockDataStartOffset int64
}

type MdictRecordBlockInfoListItem struct {
	compressSize                int64
	deCompressSize              int64
	compressAccumulatorOffset   int64
	deCompressAccumulatorOffset int64
}

/********************************
 *    public data type          *
 ********************************/

type MDictKeywordEntry struct {
	RecordStartOffset int64
	RecordEndOffset   int64
	KeyWord           string
	KeyBlockIdx       int64
}

type MDictKeywordIndex struct {
	//encoding                            int
	//encryptType                         int
	KeywordEntry MDictKeywordEntry
	RecordBlock  MDictKeywordIndexRecordBlock
}

type MDictKeywordIndexRecordBlock struct {
	DataStartOffset          int64
	CompressSize             int64
	DeCompressSize           int64
	KeyWordPartStartOffset   int64
	KeyWordPartDataEndOffset int64
}
