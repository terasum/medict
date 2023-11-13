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

package gomdict

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rasky/go-lzo"
)

// readDictHeader reads the dictionary header.
func (mdict *MdictBase) readDictHeader() error {
	// read dict header info
	dictHeader, err := readMDictFileHeader(mdict.filePath)
	if err != nil {
		return err
	}

	mdict.header = dictHeader

	// Parse header XML into a map
	headerInfo, err := parseXMLHeader(dictHeader.headerInfo)
	if err != nil {
		return err
	}

	// TODO: Handle Alder32 checksum

	meta := &mdictMeta{}

	// Handle encryption flag
	encrypted := headerInfo.Encrypted
	switch {
	case encrypted == "" || encrypted == "No":
		meta.encryptType = EncryptNoEnc
	case encrypted == "Yes":
		meta.encryptType = EncryptRecordEnc
	default:
		if encrypted[0] == '2' {
			meta.encryptType = EncryptKeyInfoEnc
		} else if encrypted[0] == '1' {
			meta.encryptType = EncryptRecordEnc
		} else {
			meta.encryptType = EncryptNoEnc
		}
	}

	// Handle version
	versionStr := headerInfo.GeneratedByEngineVersion
	version, err := strconv.ParseFloat(versionStr, 32)
	if err != nil {
		return err
	}
	meta.version = float32(version)

	// Handle number format and width based on version
	if meta.version >= 2.0 {
		meta.numberWidth = 8
		meta.numberFormat = NumfmtBe8bytesq
	} else {
		meta.numberWidth = 4
		meta.numberFormat = NumfmtBe4bytesi
	}

	// Handle encoding
	encoding := headerInfo.Encoding
	encoding = strings.ToLower(encoding)
	switch encoding {
	case "GBK", "GB2312", "gbk", "gb2312":
		meta.encoding = EncodingGb18030
	case "Big5", "BIG5", "big5":
		meta.encoding = EncodingBig5
	case "utf16", "utf-16", "UTF-16":
		meta.encoding = EncodingUtf16
	default:
		meta.encoding = EncodingUtf8
	}

	// Fix for MDD type
	if mdict.fileType == MdictTypeMdd {
		meta.encoding = EncodingUtf16
	}

	// 4 bytes header size + header_bytes_size + 4bytes alder checksum
	meta.keyBlockMetaStartOffset = int64(4 + dictHeader.headerBytesSize + 4)

	meta.description = headerInfo.Description
	meta.title = headerInfo.Title
	meta.creationDate = headerInfo.CreationDate
	meta.generatedByEngineVersion = headerInfo.GeneratedByEngineVersion

	mdict.meta = meta

	return nil
}

func readMDictFileHeader(filename string) (*mdictHeader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dictHeaderPartByteSize := int64(0)

	// Read dictionary header length
	var headerBytesSize uint32
	dictHeaderPartByteSize += 4
	err = binary.Read(file, binary.BigEndian, &headerBytesSize)
	if err != nil {
		return nil, err
	}

	// Read dictionary header info bytes
	headerInfoBytes := make([]byte, headerBytesSize)
	dictHeaderPartByteSize += int64(headerBytesSize)
	_, err = file.Read(headerInfoBytes)
	if err != nil {
		return nil, err
	}

	// Read adler32 checksum
	var adler32Checksum uint32
	dictHeaderPartByteSize += 4
	err = binary.Read(file, binary.BigEndian, &adler32Checksum)
	if err != nil {
		return nil, err
	}

	utfHeaderInfo := littleEndianBinUTF16ToUTF8(headerInfoBytes, 0, int(headerBytesSize))
	utfHeaderInfo = strings.Replace(utfHeaderInfo, "Library_Data", "Dictionary", 1)

	mdict := &mdictHeader{
		headerBytesSize:          headerBytesSize,
		headerInfoBytes:          headerInfoBytes,
		headerInfo:               utfHeaderInfo,
		adler32Checksum:          adler32Checksum,
		dictionaryHeaderByteSize: dictHeaderPartByteSize,
	}

	return mdict, nil
}

// readKeyBlockMeta keyblock header part contains keyblock meta info
func (mdict *MdictBase) readKeyBlockMeta() error {
	file, err := os.Open(mdict.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	keyBlockMeta := &mdictKeyBlockMeta{}

	// Key block meta info part
	// if version > 2.0 key-block meta part bytes length: 40
	// else: length: 16
	keyBlockMetaBytesNum := 0
	if mdict.meta.version >= 2.0 {
		keyBlockMetaBytesNum = 8 * 5
	} else {
		keyBlockMetaBytesNum = 4 * 4
	}

	// Key block meta info buffer
	keyBlockMetaBuffer, err := readFileFromPos(file, mdict.meta.keyBlockMetaStartOffset, int64(keyBlockMetaBytesNum))
	if err != nil {
		return err
	}

	// TODO: Key block info encrypted file not supported yet
	if mdict.meta.encryptType == EncryptRecordEnc {
		return errors.New("user identification is needed to read encrypted file")
	}

	// Key block meta info struct:
	// [0:8]([0:4]) - Number of key blocks
	// [8:16]([4:8]) - Number of entries
	// [16:24] - Key block info decompressed size (if version >= 2.0, otherwise, this section does not exist)
	// [24:32]([8:12]) - Key block info size
	// [32:40]([12:16]) - Key block size
	// Note: If version <2.0, the key info buffer size is 4 * 4
	//       Otherwise, the key info buffer size is 5 * 8

	// 1. [0:8]([0:4]) - Number of key blocks
	keyBlockNumBytes := keyBlockMetaBuffer[0:mdict.meta.numberWidth]

	var keyBlockNumber uint64
	if mdict.meta.numberWidth == 8 {
		keyBlockNumber = beBinToU64(keyBlockNumBytes)
	} else if mdict.meta.numberWidth == 4 {
		keyBlockNumber = uint64(beBinToU32(keyBlockNumBytes))
	}
	keyBlockMeta.keyBlockNum = int64(keyBlockNumber)

	// 2. [8:16]([4:8]) - Number of entries
	entriesNumBytes := keyBlockMetaBuffer[mdict.meta.numberWidth : mdict.meta.numberWidth+mdict.meta.numberWidth]
	if err != nil {
		return err
	}

	var entriesNum uint64
	if mdict.meta.numberWidth == 8 {
		entriesNum = beBinToU64(entriesNumBytes)
	} else if mdict.meta.numberWidth == 4 {
		entriesNum = uint64(beBinToU32(entriesNumBytes))
	}
	keyBlockMeta.entriesNum = int64(entriesNum)

	var keyBlockInfoSizeBytesStartOffset int

	// 3. [16:24] - Key block info decompressed size (if version >= 2.0, this section exists)
	if mdict.meta.version >= 2.0 {
		keyBlockInfoDecompressSizeBytes := keyBlockMetaBuffer[mdict.meta.numberWidth*2 : mdict.meta.numberWidth*2+mdict.meta.numberWidth]

		var keyBlockInfoDecompressSize uint64
		if mdict.meta.numberWidth == 8 {
			keyBlockInfoDecompressSize = beBinToU64(keyBlockInfoDecompressSizeBytes)
		} else if mdict.meta.numberWidth == 4 {
			keyBlockInfoDecompressSize = uint64(beBinToU32(keyBlockInfoDecompressSizeBytes))
		}
		keyBlockMeta.keyBlockInfoDecompressSize = int64(keyBlockInfoDecompressSize)

		keyBlockInfoSizeBytesStartOffset = mdict.meta.numberWidth * 3

	} else {
		keyBlockInfoSizeBytesStartOffset = mdict.meta.numberWidth * 2
	}

	// 4. [24:32]([8:12]) - Key block info size
	keyBlockInfoSizeBytes := keyBlockMetaBuffer[keyBlockInfoSizeBytesStartOffset : keyBlockInfoSizeBytesStartOffset+mdict.meta.numberWidth]

	var keyBlockInfoSize uint64
	if mdict.meta.numberWidth == 8 {
		keyBlockInfoSize = beBinToU64(keyBlockInfoSizeBytes)
	} else if mdict.meta.numberWidth == 4 {
		keyBlockInfoSize = uint64(beBinToU32(keyBlockInfoSizeBytes))
	}

	keyBlockMeta.keyBlockInfoCompressedSize = int64(keyBlockInfoSize)

	// 5. [32:40]([12:16]) - Key block size
	keyBlockDataSizeBytes := keyBlockMetaBuffer[keyBlockInfoSizeBytesStartOffset+mdict.meta.numberWidth : keyBlockInfoSizeBytesStartOffset+mdict.meta.numberWidth+mdict.meta.numberWidth]

	var keyBlockDataSize uint64
	if mdict.meta.numberWidth == 8 {
		keyBlockDataSize = beBinToU64(keyBlockDataSizeBytes)
	} else if mdict.meta.numberWidth == 4 {
		keyBlockDataSize = uint64(beBinToU32(keyBlockDataSizeBytes))
	}
	keyBlockMeta.keyBlockDataTotalSize = int64(keyBlockDataSize)

	// 6. [40:44] - 4 bytes checksum (TODO: Skip if version > 2.0)
	// TODO checksum verification

	// Free key block info buffer
	if mdict.meta.version >= 2.0 {
		keyBlockMeta.keyBlockInfoStartOffset = mdict.meta.keyBlockMetaStartOffset + 40 + 4
	} else {
		keyBlockMeta.keyBlockInfoStartOffset = mdict.meta.keyBlockMetaStartOffset + 16
	}

	mdict.keyBlockMeta = keyBlockMeta

	return nil
}

func (mdict *MdictBase) readKeyBlockInfo() error {
	file, err := os.Open(mdict.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer, err := readFileFromPos(file, mdict.keyBlockMeta.keyBlockInfoStartOffset, mdict.keyBlockMeta.keyBlockInfoCompressedSize)
	if err != nil {
		return err
	}

	err = mdict.decodeKeyBlockInfo(buffer)
	if err != nil {
		return err
	}
	return nil

}

func (mdict *MdictBase) decodeKeyBlockInfo(data []byte) error {
	if data[0] != 2 && data[1] != 0 && data[2] != 0 && data[3] != 0 {
		return errors.New("check key-block info magic number 2000 failed")
	}
	//fmt.Printf("decodeKeyBlockInfo: 4 magic number [%d,%d,%d,%d]\n", data[0], data[1], data[2], data[3])

	// decrypt
	var keyBlockInfoDecryptedBuffer []byte
	if mdict.meta.encryptType == EncryptKeyInfoEnc {
		// TODO decode key info
		keyBlockInfoDecryptedBuffer = mdxDecrypt(data, mdict.keyBlockMeta.keyBlockInfoCompressedSize)
	} else {
		keyBlockInfoDecryptedBuffer = data
	}

	// finally, we need to check adler32 checksum
	// key_block_info_compressed[4:8] => adler32 checksum
	//          uint32_t chksum = be_bin_to_u32((unsigned char*) (kb_info_buff +
	//          4));
	//          uint32_t adlercs = adler32checksum(key_block_info_uncomp,
	//          static_cast<uint32_t>(key_block_info_uncomp_len)) & 0xffffffff;
	//
	//          assert(chksum == adlercs);

	/// here passed, key block info is corrected
	// TODO decode key block info compressed into keys list

	// for version 2.0, will compress by zlib, lzo just just for 1.0
	// key_block_info_buff[0:8] => compress_type
	// TODO zlib decompress
	// TODO:
	// if the size of compressed data original data is unknown,
	// we malloc 8 size of source data len, we cannot estimate the original data
	// size
	// but currently, we know the size of key_block_info decompress size, so we
	// use this

	// note: we should uncompressed key_block_info_buffer[8:] data, so we need
	// (decrypted + 8, and length -8)

	decompressKeyInfoBuffer, err := zlibDecompress(keyBlockInfoDecryptedBuffer, 8, mdict.keyBlockMeta.keyBlockInfoCompressedSize-8)
	if err != nil {
		return err
	}
	if int64(len(decompressKeyInfoBuffer)) != mdict.keyBlockMeta.keyBlockInfoDecompressSize {
		return errors.New("decoded key block info data size not equals to key block meta indicates key block info size")
	}

	// decode key-block entries
	var counter int64 = 0
	var currentEntriesSize int64 = 0
	var numEntriesCounter int64 = 0
	byteWidth := 1
	textTerm := 0

	if mdict.meta.version >= 2.0 {
		byteWidth = 2
		textTerm = 1
	}

	var dataOffset = 0
	var compressSizeAccumulator = 0
	var decompressSizeAccumulator = 0

	keyBlockInfo := &mdictKeyBlockInfo{
		keyBlockEntriesStartOffset: 0,
		keyBlockInfoList:           make([]*mdictKeyBlockInfoItem, 0),
	}

	for counter < mdict.keyBlockMeta.keyBlockNum {
		firstKeySize, lastKeySize := 0, 0
		firstKey := ""
		lastKey := ""

		if mdict.meta.version >= 2.0 {
			currentEntriesSize = int64(beBinToU64(decompressKeyInfoBuffer[dataOffset : dataOffset+mdict.meta.numberWidth]))
			dataOffset += mdict.meta.numberWidth
			firstKeySize = int(beBinToU16(decompressKeyInfoBuffer[dataOffset : dataOffset+byteWidth]))
			dataOffset += byteWidth
		} else {
			currentEntriesSize = int64(beBinToU32(decompressKeyInfoBuffer[dataOffset : dataOffset+mdict.meta.numberWidth]))
			dataOffset += mdict.meta.numberWidth
			firstKeySize = int(int64(beBinToU8(decompressKeyInfoBuffer[dataOffset : dataOffset+byteWidth])))
			dataOffset += byteWidth
		}
		numEntriesCounter += currentEntriesSize

		// step_gap means first key start data_offset to first key end;
		var stepGap = 0
		var termSize = textTerm
		if mdict.meta.encoding == EncodingUtf16 || mdict.fileType == MdictTypeMdd {
			stepGap = (firstKeySize + textTerm) * 2
			termSize = textTerm * 2
		} else {
			stepGap = firstKeySize + textTerm
			termSize = textTerm
		}

		firstKey = bigEndianBinToUTF8(decompressKeyInfoBuffer, dataOffset, stepGap-termSize)

		dataOffset += stepGap

		if mdict.meta.version >= 2.0 {
			lastKeySize = int(beBinToU16(decompressKeyInfoBuffer[dataOffset : dataOffset+byteWidth]))
		} else {
			lastKeySize = int(beBinToU8(decompressKeyInfoBuffer[dataOffset : dataOffset+byteWidth]))
		}
		dataOffset += byteWidth

		if mdict.meta.encoding == EncodingUtf16 || mdict.fileType == MdictTypeMdd {
			stepGap = (lastKeySize + textTerm) * 2
			termSize = textTerm * 2
		} else {
			stepGap = lastKeySize + textTerm
			termSize = textTerm
		}

		lastKey = bigEndianBinToUTF8(decompressKeyInfoBuffer, dataOffset, stepGap-termSize)

		dataOffset += stepGap
		// key block data meta part
		keyBlockCompressSize := 0
		if mdict.meta.version >= 2.0 {
			keyBlockCompressSize = int(beBinToU64(decompressKeyInfoBuffer[dataOffset : dataOffset+mdict.meta.numberWidth]))
		} else {
			keyBlockCompressSize = int(beBinToU32(decompressKeyInfoBuffer[dataOffset : dataOffset+mdict.meta.numberWidth]))
		}
		dataOffset += mdict.meta.numberWidth

		keyBlockDecompressSize := 0
		if mdict.meta.version >= 2.0 {
			keyBlockDecompressSize = int(beBinToU64(decompressKeyInfoBuffer[dataOffset : dataOffset+mdict.meta.numberWidth]))
		} else {
			keyBlockDecompressSize = int(beBinToU32(decompressKeyInfoBuffer[dataOffset : dataOffset+mdict.meta.numberWidth]))
		}

		dataOffset += mdict.meta.numberWidth

		keyBlockInfoItem := &mdictKeyBlockInfoItem{
			firstKey:                      firstKey,
			firstKeySize:                  firstKeySize,
			lastKey:                       lastKey,
			lastKeySize:                   lastKeySize,
			keyBlockInfoIndex:             int(counter),
			keyBlockCompressSize:          int64(keyBlockCompressSize),
			keyBlockCompAccumulator:       int64(compressSizeAccumulator),
			keyBlockDeCompressSize:        int64(keyBlockDecompressSize),
			keyBlockDeCompressAccumulator: int64(decompressSizeAccumulator),
		}

		compressSizeAccumulator += keyBlockCompressSize
		decompressSizeAccumulator += keyBlockDecompressSize

		keyBlockInfo.keyBlockInfoList = append(keyBlockInfo.keyBlockInfoList, keyBlockInfoItem)

		counter++

	}
	//keyBlockInfo.keyBlockEntriesStartOffset = int64(dataOffset) + mdict.keyBlockMeta.keyBlockInfoStartOffset
	keyBlockInfo.keyBlockEntriesStartOffset = mdict.keyBlockMeta.keyBlockInfoCompressedSize + mdict.keyBlockMeta.keyBlockInfoStartOffset

	mdict.keyBlockInfo = keyBlockInfo

	if int64(compressSizeAccumulator) != mdict.keyBlockMeta.keyBlockDataTotalSize {
		return fmt.Errorf("key block data compress size not equals to meta key block data compress size(%d/%d)", compressSizeAccumulator, mdict.keyBlockMeta.keyBlockDataTotalSize)
	}

	return nil

}

func (mdict *MdictBase) readKeyEntries() error {
	file, err := os.Open(mdict.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer, err := readFileFromPos(file,
		mdict.keyBlockInfo.keyBlockEntriesStartOffset,
		mdict.keyBlockMeta.keyBlockDataTotalSize)
	if err != nil {
		return err
	}

	err = mdict.decodeKeyEntries(buffer)
	if err != nil {
		return err
	}
	return nil
}

func (mdict *MdictBase) decodeKeyEntries(keyBlockDataCompressBuffer []byte) error {

	start := int64(0)
	end := int64(0)
	compAccu := int64(0)

	keyBlockData := &mdictKeyBlockData{
		keyEntries:                 make([]*MDictKeywordEntry, 0),
		keyEntriesSize:             0,
		recordBlockMetaStartOffset: 0,
	}

	for idx := 0; idx < len(mdict.keyBlockInfo.keyBlockInfoList); idx++ {

		compressedSize := mdict.keyBlockInfo.keyBlockInfoList[idx].keyBlockCompressSize
		decompressedSize := mdict.keyBlockInfo.keyBlockInfoList[idx].keyBlockDeCompressSize

		compAccu += mdict.keyBlockInfo.keyBlockInfoList[idx].keyBlockCompressSize

		end = start + compressedSize

		if int64(start) != int64(mdict.keyBlockInfo.keyBlockInfoList[idx].keyBlockCompAccumulator) {
			return fmt.Errorf("[%d] the key-block data start offset not equal to key block compress accumulator(%d/%d/%d)\n",
				idx, start, mdict.keyBlockInfo.keyBlockInfoList[idx].keyBlockCompAccumulator, compAccu)
		}

		kbCompType := keyBlockDataCompressBuffer[start : start+4]
		// TODO 4 bytes adler32 checksum
		// # 4 bytes : adler checksum of decompressed key block
		// adler32 = unpack('>I', key_block_compressed[start + 4:start + 8])[0]

		var key_block []byte

		if kbCompType[0] == 0 {
			key_block = keyBlockDataCompressBuffer[start+8 : end]

		} else if kbCompType[0] == 1 {
			// TODO the second part
			header := []byte{0xf0, byte(int(decompressedSize))}
			// # decompress key block
			reader := bytes.NewReader(append(header, keyBlockDataCompressBuffer[start+8:end]...))

			out, err1 := lzo.Decompress1X(reader, 0, 0 /* decompressedSize, 1308672*/)
			if err1 != nil {
				return err1
			}

			key_block = out

			//} else if (kbCompType.toString('hex') === '02000000') {
		} else if kbCompType[0] == 2 {
			// decompress key block, zlib decompress
			out, err2 := zlibDecompress(keyBlockDataCompressBuffer, start+8, end-(start+8))
			if err2 != nil {
				return err2
			}
			key_block = out

			// extract one single key block into a key list
			// notice that adler32 returns signed value
			// TODO compare with previous word
			// assert(adler32 == zlib.adler32(key_block) & 0xffffffff)
		} else {
			return fmt.Errorf("cannot determine the compress type %v", kbCompType)
		}

		splitKeys := mdict.splitKeyBlock(key_block)

		keyBlockData.keyEntries = append(keyBlockData.keyEntries, splitKeys...)
		keyBlockData.keyEntriesSize += int64(len(splitKeys))

		//fmt.Printf("idx(%05d)[start:%05d/end:%05d/comps:%05d->datalen:%05d/compaccu:%d]\n", idx, start, end, compressedSize, len(key_block), compAccu)
		//fmt.Printf("key_list %+v\n", splitKeys)

		start = end
	}

	if keyBlockData.keyEntriesSize != mdict.keyBlockMeta.entriesNum {
		return errors.New("the key list items not equals to entries num")
	}
	keyBlockData.recordBlockMetaStartOffset = mdict.keyBlockInfo.keyBlockEntriesStartOffset + mdict.keyBlockMeta.keyBlockDataTotalSize

	// keep key list in memory
	mdict.keyBlockData = keyBlockData

	return nil
}

func (mdict *MdictBase) splitKeyBlock(keyBlock []byte) []*MDictKeywordEntry {
	// delimiter := ""
	width := 1

	if mdict.meta.encoding == EncodingUtf16 || mdict.fileType == MdictTypeMdd {
		//delimiter = "0000"
		width = 2
	} else {
		//delimiter = "00"
		width = 1
	}

	keyList := make([]*MDictKeywordEntry, 0)

	keyStartIndex := 0
	keyEndIndex := 0

	for keyStartIndex < len(keyBlock) {
		// # the corresponding record's offset in record block
		recordStartOffset := int64(0)

		if mdict.meta.numberWidth == 8 {
			recordStartOffset = int64(beBinToU64(keyBlock[keyStartIndex : keyStartIndex+mdict.meta.numberWidth]))
		} else {
			recordStartOffset = int64(beBinToU32(keyBlock[keyStartIndex : keyStartIndex+mdict.meta.numberWidth]))
		}

		// # key text ends with '\x00'
		i := keyStartIndex + mdict.meta.numberWidth
		for i < len(keyBlock) {
			// delimiter = '0' || // delimiter = '00'
			if (width == 1 && keyBlock[i] == 0) || (width == 2 && keyBlock[i] == 0 && keyBlock[i+1] == 0) {
				keyEndIndex = i
				break
			}
			i += width
		}

		keyTextBytes := keyBlock[keyStartIndex+mdict.meta.numberWidth : keyEndIndex]
		keyText := string(keyTextBytes)
		var err error

		if mdict.meta.encoding == EncodingUtf16 {
			keyText, err = decodeLittleEndianUtf16(keyTextBytes)
			if err != nil {
				keyText = string(keyTextBytes)
			}
		}

		if mdict.fileType == MdictTypeMdd {
			keyText, err = decodeLittleEndianUtf16(keyTextBytes)
			if err != nil {
				panic(err)
			}
		}

		keyStartIndex = keyEndIndex + width
		keyList = append(keyList, &MDictKeywordEntry{
			RecordStartOffset: recordStartOffset,
			KeyWord:           keyText,
			KeyBlockIdx:       int64(keyStartIndex),
		})
		if len(keyList) > 1 {
			keyList[len(keyList)-2].RecordEndOffset = keyList[len(keyList)-1].RecordStartOffset
		}
	}
	//keyList[len(keyList)-1].RecordLocateEndOffset = 0

	return keyList
}

func (mdict *MdictBase) readRecordBlockMeta() error {
	file, err := os.Open(mdict.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	/*
	 * [0:8/4]    - record block number
	 * [8:16/4:8] - num entries the key-value entries number
	 * [16:24/8:12] - record block info size
	 * [24:32/12:16] - record block size
	 *
	 */
	recordBlockMetaBufferLen := int64(16)
	if mdict.meta.version >= 2.0 {
		recordBlockMetaBufferLen = 32
	}

	recordBlockStartOffset := mdict.keyBlockInfo.keyBlockEntriesStartOffset + mdict.keyBlockMeta.keyBlockDataTotalSize

	buffer, err := readFileFromPos(file, recordBlockStartOffset, recordBlockMetaBufferLen)
	if err != nil {
		return err
	}

	err = mdict.decodeRecordBlockMeta(buffer, recordBlockStartOffset, recordBlockStartOffset+recordBlockMetaBufferLen)
	if err != nil {
		return err
	}
	return nil
}

/**
 * STEP 7.
 * decode record header,
 * includes:
 * [0:8/4]    - record block number
 * [8:16/4:8] - num entries the key-value entries number
 * [16:24/8:12] - record block info size
 * [24:32/12:16] - record block size
 */
func (mdict *MdictBase) decodeRecordBlockMeta(data []byte, startOffset, endOffset int64) error {
	recordBlockMeta := &mdictRecordBlockMeta{
		keyRecordMetaStartOffset: startOffset,
		keyRecordMetaEndOffset:   endOffset,
	}

	keyRecordBuffer := data
	offset := 0

	if mdict.meta.version >= 2.0 {
		recordBlockMeta.recordBlockNum = int64(beBinToU64(keyRecordBuffer[offset : offset+mdict.meta.numberWidth]))
	} else {
		recordBlockMeta.recordBlockNum = int64(beBinToU32(keyRecordBuffer[offset : offset+mdict.meta.numberWidth]))
	}

	offset += mdict.meta.numberWidth

	if mdict.meta.version >= 2.0 {
		recordBlockMeta.entriesNum = int64(beBinToU64(keyRecordBuffer[offset : offset+mdict.meta.numberWidth]))
	} else {
		recordBlockMeta.entriesNum = int64(beBinToU32(keyRecordBuffer[offset : offset+mdict.meta.numberWidth]))

	}
	if recordBlockMeta.entriesNum != mdict.keyBlockMeta.entriesNum {
		return fmt.Errorf("keyEntriesNum != meta.entriesNum")
	}

	offset += mdict.meta.numberWidth
	if mdict.meta.version >= 2.0 {
		recordBlockMeta.recordBlockInfoCompSize = int64(beBinToU64(keyRecordBuffer[offset : offset+mdict.meta.numberWidth]))
	} else {
		recordBlockMeta.recordBlockInfoCompSize = int64(beBinToU32(keyRecordBuffer[offset : offset+mdict.meta.numberWidth]))
	}

	offset += mdict.meta.numberWidth

	if mdict.meta.version >= 2.0 {
		recordBlockMeta.recordBlockCompSize = int64(beBinToU64(keyRecordBuffer[offset : offset+mdict.meta.numberWidth]))
	} else {
		recordBlockMeta.recordBlockCompSize = int64(beBinToU32(keyRecordBuffer[offset : offset+mdict.meta.numberWidth]))
	}

	mdict.recordBlockMeta = recordBlockMeta
	return nil
}

func (mdict *MdictBase) readRecordBlockInfo() error {
	file, err := os.Open(mdict.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	/*
	 * [0:8/4]    - record block number
	 * [8:16/4:8] - num entries the key-value entries number
	 * [16:24/8:12] - record block info size
	 * [24:32/12:16] - record block size
	 *
	 */
	recordBlockInfoStartOffset := mdict.recordBlockMeta.keyRecordMetaEndOffset
	recordBlockInfoLen := mdict.recordBlockMeta.recordBlockInfoCompSize

	buffer, err := readFileFromPos(file, recordBlockInfoStartOffset, recordBlockInfoLen)
	if err != nil {
		return err
	}

	err = mdict.decodeRecordBlockInfo(buffer, recordBlockInfoStartOffset, recordBlockInfoStartOffset+recordBlockInfoLen)
	if err != nil {
		return err
	}
	return nil
}

func (mdict *MdictBase) decodeRecordBlockInfo(data []byte, startOffset, endOffset int64) error {

	recordBlockInfoList := make([]*MdictRecordBlockInfoListItem, 0)
	var offset = 0
	var compAccu = int64(0)
	var decompAccu = int64(0)
	var i = int64(0)

	for i = int64(0); i < mdict.recordBlockMeta.recordBlockNum; i++ {
		compSize := int64(0)
		if mdict.meta.version >= 2.0 {
			compSize = int64(beBinToU64(data[offset : offset+mdict.meta.numberWidth]))
		} else {
			compSize = int64(beBinToU32(data[offset : offset+mdict.meta.numberWidth]))
		}
		offset += mdict.meta.numberWidth

		decompSize := int64(0)
		if mdict.meta.version >= 2.0 {
			decompSize = int64(beBinToU64(data[offset : offset+mdict.meta.numberWidth]))
		} else {
			decompSize = int64(beBinToU32(data[offset : offset+mdict.meta.numberWidth]))
		}
		offset += mdict.meta.numberWidth

		// then assign
		recordBlockInfoList = append(recordBlockInfoList, &MdictRecordBlockInfoListItem{
			compressSize:                compSize,
			deCompressSize:              decompSize,
			compressAccumulatorOffset:   compAccu,
			deCompressAccumulatorOffset: decompAccu,
		})

		// accu last
		compAccu += compSize
		decompAccu += decompSize
	}
	if int64(i) != mdict.recordBlockMeta.recordBlockNum {
		return fmt.Errorf("recordBlockInfo (i) not equals to meta.recordBlockNum [%d/%d] compA/decompA(%d/%d)", i, mdict.recordBlockMeta.recordBlockNum, compAccu, decompAccu)
	}
	if int64(offset) != mdict.recordBlockMeta.recordBlockInfoCompSize {
		return errors.New("recordBlockInfo offset not equals to meta.recordBlockInfoCompSize")
	}
	if int64(compAccu) != mdict.recordBlockMeta.recordBlockCompSize {
		return errors.New("recordBlockInfo compAccu not equals to meta.recordBlockCompSize")
	}

	recordBlockInfo := &mdictRecordBlockInfo{
		recordInfoList:             recordBlockInfoList,
		recordBlockInfoStartOffset: startOffset,
		recordBlockInfoEndOffset:   endOffset,
		recordBlockDataStartOffset: endOffset,
	}

	mdict.recordBlockInfo = recordBlockInfo

	return nil
}

func (mdict *MdictBase) buildRecordRangeTree() {
	BuildRangeTree(mdict.recordBlockInfo.recordInfoList, mdict.rangeTreeRoot)
}

func (mdict *MdictBase) keywordEntryToIndex(item *MDictKeywordEntry) (*MDictKeywordIndex, error) {
	recordBlockInfo := QueryRangeData(mdict.rangeTreeRoot, item.RecordStartOffset)

	if recordBlockInfo == nil {
		return nil, errors.New("key-item record info not found")
	}

	recordBlockStartOffset := recordBlockInfo.compressAccumulatorOffset + mdict.recordBlockInfo.recordBlockDataStartOffset
	recordBlockLen := recordBlockInfo.compressSize

	start := item.RecordStartOffset - recordBlockInfo.deCompressAccumulatorOffset
	var end int64
	if item.RecordEndOffset == 0 {
		end = recordBlockLen
	} else {
		end = item.RecordEndOffset - recordBlockInfo.deCompressAccumulatorOffset
	}

	return &MDictKeywordIndex{
		KeywordEntry: *item,
		RecordBlock: MDictKeywordIndexRecordBlock{
			DataStartOffset:          recordBlockStartOffset,
			CompressSize:             recordBlockInfo.compressSize,
			DeCompressSize:           recordBlockInfo.deCompressSize,
			KeyWordPartStartOffset:   start,
			KeyWordPartDataEndOffset: end,
		},
	}, nil

}

func (mdict *MdictBase) keywordEntryToIndex1(item *MDictKeywordEntry) (*MDictKeywordIndex, error) {

	var recordBlockInfo *MdictRecordBlockInfoListItem

	var i = 0
	for ; i < len(mdict.recordBlockInfo.recordInfoList)-1; i++ {
		curr := mdict.recordBlockInfo.recordInfoList[i]
		next := mdict.recordBlockInfo.recordInfoList[i+1]
		if item.RecordStartOffset >= curr.deCompressAccumulatorOffset && item.RecordStartOffset < next.deCompressAccumulatorOffset {
			recordBlockInfo = curr
			break
		}
	}

	// the last one
	if i == len(mdict.recordBlockInfo.recordInfoList)-1 {
		lastOne := mdict.recordBlockInfo.recordInfoList[len(mdict.recordBlockInfo.recordInfoList)-1]
		if item.RecordStartOffset < lastOne.deCompressAccumulatorOffset+lastOne.deCompressSize {
			recordBlockInfo = lastOne
		}
	}

	if recordBlockInfo == nil {
		fmt.Printf("record block info is nil, current keyBlockEntry: %+v, last recordBlockInfo: %+v\n", item, mdict.recordBlockInfo.recordInfoList[len(mdict.recordBlockInfo.recordInfoList)-1])
		return nil, errors.New("key-item record info not found")
	}

	recordBlockStartOffset := recordBlockInfo.compressAccumulatorOffset + mdict.recordBlockInfo.recordBlockDataStartOffset
	recordBlockLen := recordBlockInfo.compressSize

	start := item.RecordStartOffset - recordBlockInfo.deCompressAccumulatorOffset
	var end int64
	if item.RecordEndOffset == 0 {
		end = int64(recordBlockLen)
	} else {
		end = item.RecordEndOffset - recordBlockInfo.deCompressAccumulatorOffset
	}

	return &MDictKeywordIndex{
		KeywordEntry: *item,
		RecordBlock: MDictKeywordIndexRecordBlock{
			DataStartOffset:          recordBlockStartOffset,
			CompressSize:             recordBlockInfo.compressSize,
			DeCompressSize:           recordBlockInfo.deCompressSize,
			KeyWordPartStartOffset:   start,
			KeyWordPartDataEndOffset: end,
		},
	}, nil

}

func (mdict *MdictBase) locateByKeywordIndex(index *MDictKeywordIndex) ([]byte, error) {
	return locateDefByKWIndex(index,
		mdict.filePath,
		mdict.meta.encryptType == EncryptRecordEnc,
		mdict.fileType == MdictTypeMdd,
		mdict.meta.encoding == EncodingUtf16)
}

func locateDefByKWIndex(index *MDictKeywordIndex, filePath string, isRecordEncrypted, isMdd, isUtf16 bool) ([]byte, error) {
	log.Infof("locateDefByKWIndex invoked %+v, filepath %s, isRecordEncrypted %v, isMdd %v, isUTF16 %v", index, filePath, isRecordEncrypted, isMdd, isUtf16)
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("open file err %s", err.Error())
		return nil, err
	}
	defer file.Close()
	recordBlockDataCompBuff, err := readFileFromPos(file, index.RecordBlock.DataStartOffset, index.RecordBlock.CompressSize)
	if err != nil {

		log.Errorf("readFileFromPos %s", err.Error())
		return nil, err
	}

	if recordBlockDataCompBuff == nil {
		log.Errorf("record block data buffer is null, index: %v", index)
		return nil, errors.New("record block data buffer is null")
	}

	// 4 bytes: compression type
	var rbCompType = recordBlockDataCompBuff[0:4]

	// record_block stores the final record data
	var recordBlock []byte

	// TODO: ignore adler32 offset
	// Note: here ignore the checksum part
	// bytes: adler32 checksum of decompressed record block
	// adler32 = unpack('>I', record_block_compressed[4:8])[0]
	if rbCompType[0] == 0 {
		recordBlock = recordBlockDataCompBuff[8:index.RecordBlock.CompressSize]
	} else {
		// decrypt
		var blockBufDecrypted []byte
		// if encrypt type == 1, the record block was encrypted
		if isRecordEncrypted {
			// const passkey = new Uint8Array(8);
			// record_block_compressed.copy(passkey, 0, 4, 8);
			// passkey.set([0x95, 0x36, 0x00, 0x00], 4); // key part 2: fixed data
			blockBufDecrypted = mdxDecrypt(recordBlockDataCompBuff, index.RecordBlock.CompressSize)
		} else {
			blockBufDecrypted = recordBlockDataCompBuff[8:index.RecordBlock.CompressSize]
		}

		// decompress
		if rbCompType[0] == 1 {
			// TODO the second part
			header := []byte{0xf0, byte(int(index.RecordBlock.CompressSize))}
			// # decompress key block
			reader := bytes.NewReader(append(header, blockBufDecrypted...))

			out, err1 := lzo.Decompress1X(reader, 0, 0 /* decompressedSize, 1308672*/)
			if err1 != nil {
				log.Errorf("stopped by Decompress1X %s", err1.Error())
				return nil, err1
			}

			recordBlock = out

		} else if rbCompType[0] == 2 {
			var err2 error
			recordBlock, err2 = zlibDecompress(blockBufDecrypted, 0, int64(len(blockBufDecrypted)))
			if err2 != nil {
				log.Errorf("stopped by zlibDecompress %s", err2.Error())
				return nil, err2
			}
		}
	}

	// TODO: ignore the checksum
	// notice that adler32 return signed value
	// assert(adler32 == zlib.adler32(record_block) & 0xffffffff)

	if int64(len(recordBlock)) != index.RecordBlock.DeCompressSize {
		log.Errorf("stopped by len(recordBlock) != index.RecordBlock.DeCompressSize")
		return nil, errors.New("recordBlock length not equals decompress Size")
	}

	start := index.RecordBlock.KeyWordPartStartOffset
	end := index.RecordBlock.KeyWordPartDataEndOffset

	data := recordBlock[start:end]

	if isMdd {
		log.Errorf("return mdd data")
		return data, nil
	}

	if isUtf16 {
		log.Infof("keyword %s, data len %d", index.KeywordEntry.KeyWord, len(data))
		datastr, err1 := decodeLittleEndianUtf16(data)
		if err1 != nil {
			return nil, err
		}
		return []byte(datastr), nil
	}
	return data, nil
}

func (mdict *MdictBase) locateByKeywordEntry(item *MDictKeywordEntry) ([]byte, error) {
	file, err := os.Open(mdict.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var recordBlockInfo *MdictRecordBlockInfoListItem

	var i = 0
	for ; i < len(mdict.recordBlockInfo.recordInfoList)-1; i++ {
		curr := mdict.recordBlockInfo.recordInfoList[i]
		next := mdict.recordBlockInfo.recordInfoList[i+1]
		if item.RecordStartOffset >= curr.deCompressAccumulatorOffset && item.RecordStartOffset < next.deCompressAccumulatorOffset {
			recordBlockInfo = curr
			break
		}
	}

	// the last one
	if i == len(mdict.recordBlockInfo.recordInfoList)-1 {
		lastOne := mdict.recordBlockInfo.recordInfoList[len(mdict.recordBlockInfo.recordInfoList)-1]
		if item.RecordStartOffset < lastOne.deCompressAccumulatorOffset+lastOne.deCompressSize {
			recordBlockInfo = lastOne
		}
	}

	if recordBlockInfo == nil {
		fmt.Printf("record block info is nil, current keyBlockEntry: %+v, last recordBlockInfo: %+v\n", item, mdict.recordBlockInfo.recordInfoList[len(mdict.recordBlockInfo.recordInfoList)-1])
		return nil, errors.New("key-item record info not found")
	}

	recordBlockStartOffset := recordBlockInfo.compressAccumulatorOffset + mdict.recordBlockInfo.recordBlockDataStartOffset
	recordBlockLen := recordBlockInfo.compressSize

	recordBlockDataCompBuff, err := readFileFromPos(file, recordBlockStartOffset, recordBlockLen)
	if err != nil {
		return nil, err
	}

	// 4 bytes: compression type
	var rbCompType = recordBlockDataCompBuff[0:4]

	// record_block stores the final record data
	var recordBlock []byte

	// TODO: ignore adler32 offset
	// Note: here ignore the checksum part
	// bytes: adler32 checksum of decompressed record block
	// adler32 = unpack('>I', record_block_compressed[4:8])[0]
	if rbCompType[0] == 0 {
		recordBlock = recordBlockDataCompBuff[8:recordBlockInfo.compressSize]
	} else {
		// decrypt
		var blockBufDecrypted []byte
		// if encrypt type == 1, the record block was encrypted
		if mdict.meta.encryptType == EncryptRecordEnc {
			// const passkey = new Uint8Array(8);
			// record_block_compressed.copy(passkey, 0, 4, 8);
			// passkey.set([0x95, 0x36, 0x00, 0x00], 4); // key part 2: fixed data
			blockBufDecrypted = mdxDecrypt(recordBlockDataCompBuff, recordBlockInfo.compressSize)
		} else {
			blockBufDecrypted = recordBlockDataCompBuff[8:recordBlockInfo.compressSize]
		}

		// decompress
		if rbCompType[0] == 1 {
			// TODO the second part
			header := []byte{0xf0, byte(int(recordBlockInfo.compressSize))}
			// # decompress key block
			reader := bytes.NewReader(append(header, blockBufDecrypted...))

			out, err1 := lzo.Decompress1X(reader, 0, 0 /* decompressedSize, 1308672*/)
			if err1 != nil {
				return nil, err1
			}

			recordBlock = out

		} else if rbCompType[0] == 2 {
			var err2 error
			recordBlock, err2 = zlibDecompress(blockBufDecrypted, 0, int64(len(blockBufDecrypted)))
			if err2 != nil {
				return nil, err2
			}
		}
	}

	// TODO: ignore the checksum
	// notice that adler32 return signed value
	// assert(adler32 == zlib.adler32(record_block) & 0xffffffff)

	if int64(len(recordBlock)) != recordBlockInfo.deCompressSize {
		return nil, errors.New("recordBlock length not equals decompress Size")
	}

	start := item.RecordStartOffset - recordBlockInfo.deCompressAccumulatorOffset
	var end int64
	if item.RecordEndOffset == 0 {
		end = int64(len(recordBlock))
	} else {
		end = item.RecordEndOffset - recordBlockInfo.deCompressAccumulatorOffset
	}

	data := recordBlock[start:end]

	if mdict.fileType == MdictTypeMdd {
		return data, nil
	}

	if mdict.meta.encoding == EncodingUtf16 {
		datastr, err1 := decodeLittleEndianUtf16(data)
		if err1 != nil {
			return nil, err
		}
		return []byte(datastr), nil
	}
	return data, nil

}

func (mdict *MdictBase) getKeyWordEntries() ([]*MDictKeywordEntry, error) {
	return mdict.keyBlockData.keyEntries, nil
}
