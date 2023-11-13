package gomdict

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func (mdict *Mdict) Digest() string {
	outstr := ""
	outstr += "meta:\n"
	outstr += "-----------------------------\n"
	outstr += mdict.digestMeta()
	outstr += "\n"
	outstr += "-----------------------------\n"
	outstr += "header:\n"
	outstr += "-----------------------------\n"
	outstr += fmt.Sprintf("headerInfo: %+v\n", mdict.header.headerInfo)
	outstr += fmt.Sprintf("headerBytesSize: %+v\n", mdict.header.headerBytesSize)
	outstr += fmt.Sprintf("dictionaryHeaderByteSize: %+v\n", mdict.header.dictionaryHeaderByteSize)
	outstr += fmt.Sprintf("adler32Checksum: %+v\n", mdict.header.adler32Checksum)
	outstr += "-----------------------------\n"
	outstr += "keyBlockMeta:\n"
	outstr += "-----------------------------\n"
	outstr += mdict.digestKeyBlockMeta()
	outstr += "-----------------------------\n"
	outstr += "keyBlockInfo:\n"
	outstr += "-----------------------------\n"
	outstr += fmt.Sprintf("keyBlockEntriesStartOffset: %+v\n", mdict.keyBlockInfo.keyBlockEntriesStartOffset)
	outstr += mdict.digestKeyBlockInfo()

	outstr += "-----------------------------\n"
	outstr += "keyBlockData:\n"
	outstr += "-----------------------------\n"
	outstr += fmt.Sprintf("keyEntriesSize: %+v\n", mdict.keyBlockData.keyEntriesSize)
	outstr += fmt.Sprintf("recordBlockMetaStartOffset: %+v\n", mdict.keyBlockData.recordBlockMetaStartOffset)
	outstr += fmt.Sprintf("keyEntries[0]: %+v\n", mdict.keyBlockData.keyEntries[0])
	outstr += fmt.Sprintf("keyEntries[1]: %+v\n", mdict.keyBlockData.keyEntries[1])
	outstr += fmt.Sprintf("keyEntries[3]: %+v\n", mdict.keyBlockData.keyEntries[2])
	outstr += fmt.Sprintf("keyEntries[4]: %+v\n", mdict.keyBlockData.keyEntries[4])
	outstr += fmt.Sprintf("keyEntries[5]: %+v\n", mdict.keyBlockData.keyEntries[5])
	outstr += fmt.Sprintf("keyEntries.len(): %+v\n", len(mdict.keyBlockData.keyEntries))

	outstr += "-----------------------------\n"
	outstr += "recordBlockMeta:\n"
	outstr += "-----------------------------\n"
	outstr += fmt.Sprintf("keyRecordMetaStartOffset: %+v\n", mdict.recordBlockMeta.keyRecordMetaStartOffset)
	outstr += fmt.Sprintf("keyRecordMetaEndOffset: %+v\n", mdict.recordBlockMeta.keyRecordMetaEndOffset)
	outstr += fmt.Sprintf("entriesNum: %+v\n", mdict.recordBlockMeta.entriesNum)
	outstr += fmt.Sprintf("recordBlockInfoCompSize: %+v\n", mdict.recordBlockMeta.recordBlockInfoCompSize)
	outstr += fmt.Sprintf("recordBlockCompSize: %+v\n", mdict.recordBlockMeta.recordBlockCompSize)

	outstr += "-----------------------------\n"
	outstr += "recordBlockInfo:\n"
	outstr += "-----------------------------\n"
	outstr += fmt.Sprintf("recordBlockDataStartOffset: %+v\n", mdict.recordBlockInfo.recordBlockDataStartOffset)
	outstr += fmt.Sprintf("recordBlockInfoStartOffset: %+v\n", mdict.recordBlockInfo.recordBlockInfoStartOffset)
	outstr += fmt.Sprintf("recordBlockInfoEndOffset: %+v\n", mdict.recordBlockInfo.recordBlockInfoEndOffset)
	outstr += fmt.Sprintf("recordInfoList[0]: %+v\n", mdict.recordBlockInfo.recordInfoList[0])
	outstr += fmt.Sprintf("recordInfoList[1]: %+v\n", mdict.recordBlockInfo.recordInfoList[1])
	outstr += fmt.Sprintf("recordInfoList[2]: %+v\n", mdict.recordBlockInfo.recordInfoList[2])
	outstr += fmt.Sprintf("recordInfoList[3]: %+v\n", mdict.recordBlockInfo.recordInfoList[3])
	outstr += fmt.Sprintf("recordInfoList.len(): %+v\n", len(mdict.recordBlockInfo.recordInfoList))

	//outstr += "-----------------------------\n"
	//outstr += "RecordBlockData:\n"
	//outstr += fmt.Sprintf("RecordBlockData: %+v\n", mdict.RecordBlockData)

	return outstr

}

func (mdict *Mdict) digestMeta() string {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Key", "Value")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	writer := bytes.NewBufferString("")

	tbl.WithWriter(writer)

	tbl.AddRow("encryptType", fmt.Sprintf("%v", mdict.meta.encryptType))
	tbl.AddRow("version", fmt.Sprintf("%v", mdict.meta.version))
	tbl.AddRow("encoding", fmt.Sprintf("%v", mdict.meta.encoding))
	tbl.AddRow("NumberWith", fmt.Sprintf("%v", mdict.meta.numberWidth))
	tbl.AddRow("numberFormat", fmt.Sprintf("%v", mdict.meta.numberFormat))
	tbl.AddRow("keyBlockMetaStartOffset", fmt.Sprintf("%v", mdict.meta.keyBlockMetaStartOffset))
	tbl.AddRow("keyBlockMetaStartOffset", fmt.Sprintf("%v", mdict.meta.description))
	tbl.Print()
	return writer.String()
}

func (mdict *Mdict) digestKeyBlockMeta() string {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Key", "Value")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	writer := bytes.NewBufferString("")
	tbl.WithWriter(writer)
	tbl.AddRow("keyBlockNum", fmt.Sprintf("%v", mdict.keyBlockMeta.keyBlockNum))
	tbl.AddRow("entriesNum", fmt.Sprintf("%v", mdict.keyBlockMeta.entriesNum))
	tbl.AddRow("keyBlockInfoDecompressSize", fmt.Sprintf("%v", mdict.keyBlockMeta.keyBlockInfoDecompressSize))
	tbl.AddRow("keyBlockInfoCompressedSize", fmt.Sprintf("%v", mdict.keyBlockMeta.keyBlockInfoCompressedSize))
	tbl.AddRow("keyBlockDataTotalSize", fmt.Sprintf("%v", mdict.keyBlockMeta.keyBlockDataTotalSize))
	tbl.AddRow("keyBlockInfoStartOffset", fmt.Sprintf("%v", mdict.keyBlockMeta.keyBlockInfoStartOffset))
	tbl.Print()
	return writer.String()
}

func (mdict *Mdict) digestKeyBlockInfo() string {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Key", "Value")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	writer := bytes.NewBufferString("")
	tbl.WithWriter(writer)
	tbl.AddRow("0", fmt.Sprintf("%+v", mdict.keyBlockInfo.keyBlockInfoList[0]))
	tbl.AddRow("1", fmt.Sprintf("%+v", mdict.keyBlockInfo.keyBlockInfoList[1]))
	tbl.AddRow("2", fmt.Sprintf("%+v", mdict.keyBlockInfo.keyBlockInfoList[2]))
	tbl.AddRow("3", fmt.Sprintf("%+v", mdict.keyBlockInfo.keyBlockInfoList[3]))
	tbl.AddRow("4", fmt.Sprintf("%+v", mdict.keyBlockInfo.keyBlockInfoList[4]))
	tbl.Print()
	return writer.String()
}

//
//type MDictRecordBlockData struct {
//	RecordItemList         []*MDictRecordDataItem
//	RecordBlockStartOffset int64
//	RecordBlockEndOffset   int64
//}
//
//type MDictRecordDataItem struct {
//	KeyWord                          string
//	RecordEntryIndex                 int64
//	RecordInfoIndex                  int64
//	RecordBlockIndex                 int64
//	RecordBlockCompressStart         int64
//	RecordBlockCompressEnd           int64
//	RecordBlockCompressSize          int64
//	RecordBlockDeCompressSize        int64
//	RecordBlockCompressType          string
//	RecordBlockEncrypted             bool
//	RecordBlockFileRelativeOffset    int64
//	RecordBlockCompressAccumulator   int64
//	RecordBlockDeCompressAccumulator int64
//
//	RecordEntryDecompressStart int64
//	RecordEntryDecompressEnd   int64
//}

//func (mdict *MdictBase) ReadRecordBlockData() error {
//	file, err := os.Open(mdict.filePath)
//	if err != nil {
//		return err
//	}
//	defer file.close()
//
//	recordBlockInfoStartOffset := mdict.recordBlockInfo.recordBlockDataStartOffset
//
//	recordBlockInfoLen := mdict.recordBlockInfo.recordInfoList[len(mdict.recordBlockInfo.recordInfoList)-1].compressAccumulatorOffset
//
//	buffer, err := readFileFromPos(file, recordBlockInfoStartOffset, recordBlockInfoLen)
//	if err != nil {
//		return err
//	}
//
//	err = mdict.decodeRecordBlockData(buffer, recordBlockInfoStartOffset, recordBlockInfoStartOffset+recordBlockInfoLen)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// TODO FIX THIS, 这个方法尚未实现
//func (mdict *MdictBase) decodeRecordBlockData(data []byte, startOffset, endOffset int64) error {
//
//	var sizeCounter = int64(0)
//	var itemCounter = int64(0)
//	var recordBlockOffset = int64(0)
//	var recordCompType = "NONE"
//
//	var recordBlockData = &MDictRecordBlockData{
//		RecordBlockStartOffset: startOffset,
//		RecordBlockEndOffset:   endOffset,
//		RecordItemList:         make([]*MDictRecordDataItem, 0),
//	}
//
//	for idx := 0; idx < len(mdict.recordBlockInfo.recordInfoList); idx++ {
//		var recordBlockCompSize = mdict.recordBlockInfo.recordInfoList[idx].RecordBlockCompressSize
//		var recordBlockDecompSize = mdict.recordBlockInfo.recordInfoList[idx].RecordBlockDeCompressSize
//		var recordBlockDataCompBuff = data[recordBlockOffset : recordBlockOffset+recordBlockCompSize]
//		var recordBlockCompressAccumulator = mdict.recordBlockInfo.recordInfoList[idx].compressAccumulatorOffset
//		var recordBlockDeCompressAccumulator = mdict.recordBlockInfo.recordInfoList[idx].deCompressAccumulatorOffset
//
//		// 4 bytes: compression type
//		var rbCompType = recordBlockDataCompBuff[0:4]
//
//		// record_block stores the final record data
//		var recordBlock []byte
//
//		// TODO: igore adler32 offset
//		// Note: here ignore the checksum part
//		// bytes: adler32 checksum of decompressed record block
//		// adler32 = unpack('>I', record_block_compressed[4:8])[0]
//		if rbCompType[0] == 0 {
//			recordBlock = recordBlockDataCompBuff[8 : 8+recordBlockCompSize]
//			recordCompType = "NONE"
//		} else {
//			// decrypt
//			var blockBufDecrypted []byte
//			// if encrypt type == 1, the record block was encrypted
//			if mdict.meta.encryptType == EncryptRecordEnc {
//				// const passkey = new Uint8Array(8);
//				// record_block_compressed.copy(passkey, 0, 4, 8);
//				// passkey.set([0x95, 0x36, 0x00, 0x00], 4); // key part 2: fixed data
//				blockBufDecrypted = mdxDecrypt(recordBlockDataCompBuff, recordBlockCompSize)
//			} else {
//				blockBufDecrypted = recordBlockDataCompBuff[8:recordBlockCompSize]
//			}
//
//			// decompress
//			if rbCompType[0] == 1 {
//				// TODO the second part
//				recordCompType = "LZOX1"
//				header := []byte{0xf0, byte(int(recordBlockDecompSize))}
//				// # decompress key block
//				reader := bytes.NewReader(append(header, blockBufDecrypted...))
//
//				out, err1 := lzo.Decompress1X(reader, 0, 0 /* decompressedSize, 1308672*/)
//				if err1 != nil {
//					return err1
//				}
//
//				recordBlock = out
//
//			} else if rbCompType[0] == 2 {
//				recordCompType = "ZLIB"
//				var err error
//				recordBlock, err = zlibDecompress(blockBufDecrypted, 0, int64(len(blockBufDecrypted)))
//				if err != nil {
//					return err
//				}
//			}
//		}
//
//		// notice that adler32 return signed value
//		// TODO: ignore the checksum
//		// assert(adler32 == zlib.adler32(record_block) & 0xffffffff)
//
//		if int64(len(recordBlock)) != recordBlockDecompSize {
//			return errors.New("recordBlock length not equals decompress Size")
//		}
//
//		/**
//		 * 请注意，block 是会有很多个的，而每个block都可能会被压缩
//		 * 而 key_list中的 record_start, key_text是相对每一个block而言的，end是需要每次解析的时候算出来的
//		 * 所有的record_start/length/end都是针对解压后的block而言的
//		 */
//
//		// split record block according to the offset info from key block
//		//var offset = int64(0)
//		keyblockEntry := mdict.keyBlockData.keyEntries
//		for i := 0; i < len(keyblockEntry); i++ {
//
//			var recordEntryStart = keyblockEntry[i].RecordLocateStartOffset
//			var keyText = keyblockEntry[i].KeyWord
//			// # reach the end of current record block
//			//if recordEntryStart-offset >= int64(len(recordBlock)) {
//			//	break
//			//}
//			// # record end index
//			var recordEntryEnd int64 = 0
//			if i < len(keyblockEntry)-1 {
//				recordEntryEnd = int64(keyblockEntry[i+1].RecordLocateStartOffset)
//			} else {
//				//recordEntryEnd = int64(len(mdict.keyBlockData.KeyList)) + offset
//				recordEntryEnd = endOffset
//			}
//
//			//fmt.Printf("keyText: %s, recordEntryStart: %d, recordEntryEnd:%d\n", keyText, recordEntryStart, recordEntryEnd)
//
//			recordBlockData.RecordItemList = append(recordBlockData.RecordItemList,
//				&MDictRecordDataItem{
//					KeyWord:                          keyText,
//					RecordBlockCompressStart:         recordBlockOffset,
//					RecordBlockCompressEnd:           recordBlockOffset + recordBlockCompSize,
//					RecordBlockCompressSize:          recordBlockCompSize,
//					RecordBlockDeCompressSize:        recordBlockDecompSize,
//					RecordBlockCompressType:          recordCompType,
//					RecordBlockEncrypted:             mdict.meta.encryptType == EncryptRecordEnc,
//					RecordBlockFileRelativeOffset:    recordBlockOffset + mdict.recordBlockInfo.recordBlockDataStartOffset,
//					RecordBlockCompressAccumulator:   recordBlockCompressAccumulator,
//					RecordBlockDeCompressAccumulator: recordBlockDeCompressAccumulator,
//
//					RecordEntryIndex:           int64(itemCounter),
//					RecordEntryDecompressStart: recordEntryStart,
//					RecordEntryDecompressEnd:   recordEntryEnd,
//
//					RecordInfoIndex: int64(idx),
//				})
//
//			itemCounter++
//		}
//
//		recordBlockOffset += recordBlockCompSize
//		//offset += int64(len(recordBlock))
//		sizeCounter += recordBlockCompSize
//	}
//
//	if sizeCounter != mdict.recordBlockMeta.recordBlockCompSize {
//		return fmt.Errorf("record entries compressed size sum not equals record block compress size(%d/%d)", sizeCounter, mdict.recordBlockMeta.recordBlockCompSize)
//	}
//
//	mdict.RecordBlockData = recordBlockData
//
//	return nil
//}
