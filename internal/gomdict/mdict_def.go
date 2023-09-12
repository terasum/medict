package gomdict

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
	FilePath string
	FileType MdictType
	Meta     *MDictMeta

	Header       *MDictHeader
	KeyBlockMeta *MDictKeyBlockMeta
	KeyBlockInfo *MDictKeyBlockInfo
	KeyBlockData *MDictKeyBlockData

	RecordBlockMeta *MDictRecordBlockMeta
	RecordBlockInfo *MDictRecordBlockInfo
	RecordBlockData *MDictRecordBlockData
}

type MDictHeader struct {
	HeaderBytesSize          uint32
	HeaderInfoBytes          []byte
	HeaderInfo               string
	Adler32Checksum          uint32
	DictionaryHeaderByteSize int64
}

type MDictMeta struct {
	EncryptType  int
	Version      float32
	NumberWidth  int
	NumberFormat int
	Encoding     int

	// key-block part bytes start offset in the mdx/mdd file
	KeyBlockMetaStartOffset int64
}

type MDictKeyBlockMeta struct {
	// KeyBlockNum key block number size
	KeyBlockNum int64
	// entriesNums entries number size
	EntriesNum int64
	// key-block information size (decompressed)
	KeyBlockInfoDecompressSize int64
	// key-block information size (compressed)
	KeyBlockInfoCompressedSize int64
	// key-block Data Size (decompressed)
	KeyBlockDataTotalSize int64
	// key-block information start position in the mdx/mdd file
	KeyBlockInfoStartOffset int64
}

type MDictKeyBlockInfo struct {
	KeyBlockEntriesStartOffset int64
	KeyBlockInfoList           []*MDictKeyBlockInfoItem
}

type MDictKeyBlockInfoItem struct {
	FirstKey                      string
	FirstKeySize                  int
	LastKey                       string
	LastKeySize                   int
	KeyBlockInfoIndex             int
	KeyBlockCompressSize          int64
	KeyBlockCompAccumulator       int64
	KeyBlockDeCompressSize        int64
	KeyBlockDeCompressAccumulator int64
}

type MDictKeyBlockData struct {
	KeyEntries                 []*MDictKeyBlockEntry
	KeyEntriesSize             int64
	RecordBlockMetaStartOffset int64
}

type MDictKeyBlockEntry struct {
	RecordStartOffset int64
	RecordEndOffset   int64
	KeyWord           string
	KeyBlockIdx       int64
}

type MDictRecordBlockData struct {
	RecordItemList         []*MDictRecordDataItem
	RecordBlockStartOffset int64
	RecordBlockEndOffset   int64
}

type MDictRecordDataItem struct {
	KeyWord                          string
	RecordEntryIndex                 int64
	RecordInfoIndex                  int64
	RecordBlockIndex                 int64
	RecordBlockCompressStart         int64
	RecordBlockCompressEnd           int64
	RecordBlockCompressSize          int64
	RecordBlockDeCompressSize        int64
	RecordBlockCompressType          string
	RecordBlockEncrypted             bool
	RecordBlockFileRelativeOffset    int64
	RecordBlockCompressAccumulator   int64
	RecordBlockDeCompressAccumulator int64

	RecordEntryDecompressStart int64
	RecordEntryDecompressEnd   int64
}
type MDictRecordBlockMeta struct {
	KeyRecordMetaStartOffset int64
	KeyRecordMetaEndOffset   int64

	RecordBlockNum          int64
	EntriesNum              int64
	RecordBlockInfoCompSize int64
	RecordBlockCompSize     int64
}
type MDictRecordBlockInfo struct {
	RecordInfoList             []*MDictRecordBlockInfoListItem
	RecordBlockInfoStartOffset int64
	RecordBlockInfoEndOffset   int64
	RecordBlockDataStartOffset int64
}

type MDictRecordBlockInfoListItem struct {
	CompressSize                int64
	DeCompressSize              int64
	CompressAccumulatorOffset   int64
	DeCompressAccumulatorOffset int64
}
