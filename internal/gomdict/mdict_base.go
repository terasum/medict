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
	"github.com/rasky/go-lzo"
	"os"
	"strconv"
	"strings"
)

// ReadDictHeader reads the dictionary header.
func (mdict *MdictBase) ReadDictHeader() error {
	// read dict header info
	dictHeader, err := readMDictFileHeader(mdict.FilePath)
	if err != nil {
		return err
	}

	mdict.Header = dictHeader

	// Parse header XML into a map
	headerInfo, err := parseXMLHeader(dictHeader.HeaderInfo)
	if err != nil {
		return err
	}

	// TODO: Handle Alder32 checksum

	meta := &MDictMeta{}

	// Handle encryption flag
	encrypted := headerInfo.Encrypted
	switch {
	case encrypted == "" || encrypted == "No":
		meta.EncryptType = EncryptNoEnc
	case encrypted == "Yes":
		meta.EncryptType = EncryptRecordEnc
	default:
		if encrypted[0] == '2' {
			meta.EncryptType = EncryptKeyInfoEnc
		} else if encrypted[0] == '1' {
			meta.EncryptType = EncryptRecordEnc
		} else {
			meta.EncryptType = EncryptNoEnc
		}
	}

	// Handle version
	versionStr := headerInfo.GeneratedByEngineVersion
	version, err := strconv.ParseFloat(versionStr, 32)
	if err != nil {
		return err
	}
	meta.Version = float32(version)

	// Handle number format and width based on version
	if meta.Version >= 2.0 {
		meta.NumberWidth = 8
		meta.NumberFormat = NumfmtBe8bytesq
	} else {
		meta.NumberWidth = 4
		meta.NumberFormat = NumfmtBe4bytesi
	}

	// Handle encoding
	encoding := headerInfo.Encoding
	encoding = strings.ToLower(encoding)
	switch encoding {
	case "GBK", "GB2312", "gbk", "gb2312":
		meta.Encoding = EncodingGb18030
	case "Big5", "BIG5", "big5":
		meta.Encoding = EncodingBig5
	case "utf16", "utf-16", "UTF-16":
		meta.Encoding = EncodingUtf16
	default:
		meta.Encoding = EncodingUtf8
	}

	// Fix for MDD type
	if mdict.FileType == MdictTypeMdd {
		meta.Encoding = EncodingUtf16
	}

	// 4 bytes header size + header_bytes_size + 4bytes alder checksum
	meta.KeyBlockMetaStartOffset = int64(4 + dictHeader.HeaderBytesSize + 4)

	meta.Description = headerInfo.Description
	meta.Title = headerInfo.Title
	meta.CreationDate = headerInfo.CreationDate
	meta.GeneratedByEngineVersion = headerInfo.GeneratedByEngineVersion

	mdict.Meta = meta

	return nil
}

func readMDictFileHeader(filename string) (*MDictHeader, error) {
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

	mdict := &MDictHeader{
		HeaderBytesSize:          headerBytesSize,
		HeaderInfoBytes:          headerInfoBytes,
		HeaderInfo:               utfHeaderInfo,
		Adler32Checksum:          adler32Checksum,
		DictionaryHeaderByteSize: dictHeaderPartByteSize,
	}

	return mdict, nil
}

// ReadKeyBlockMeta keyblock header part contains keyblock meta info
func (mdict *MdictBase) ReadKeyBlockMeta() error {
	file, err := os.Open(mdict.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	keyBlockMeta := &MDictKeyBlockMeta{}

	// Key block meta info part
	// if version > 2.0 key-block meta part bytes length: 40
	// else: length: 16
	keyBlockMetaBytesNum := 0
	if mdict.Meta.Version >= 2.0 {
		keyBlockMetaBytesNum = 8 * 5
	} else {
		keyBlockMetaBytesNum = 4 * 4
	}

	// Key block meta info buffer
	keyBlockMetaBuffer, err := readFileFromPos(file, mdict.Meta.KeyBlockMetaStartOffset, int64(keyBlockMetaBytesNum))
	if err != nil {
		return err
	}

	// TODO: Key block info encrypted file not supported yet
	if mdict.Meta.EncryptType == EncryptRecordEnc {
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
	keyBlockNumBytes := keyBlockMetaBuffer[0:mdict.Meta.NumberWidth]

	var keyBlockNumber uint64
	if mdict.Meta.NumberWidth == 8 {
		keyBlockNumber = beBinToU64(keyBlockNumBytes)
	} else if mdict.Meta.NumberWidth == 4 {
		keyBlockNumber = uint64(beBinToU32(keyBlockNumBytes))
	}
	keyBlockMeta.KeyBlockNum = int64(keyBlockNumber)

	// 2. [8:16]([4:8]) - Number of entries
	entriesNumBytes := keyBlockMetaBuffer[mdict.Meta.NumberWidth : mdict.Meta.NumberWidth+mdict.Meta.NumberWidth]
	if err != nil {
		return err
	}

	var entriesNum uint64
	if mdict.Meta.NumberWidth == 8 {
		entriesNum = beBinToU64(entriesNumBytes)
	} else if mdict.Meta.NumberWidth == 4 {
		entriesNum = uint64(beBinToU32(entriesNumBytes))
	}
	keyBlockMeta.EntriesNum = int64(entriesNum)

	var keyBlockInfoSizeBytesStartOffset int

	// 3. [16:24] - Key block info decompressed size (if version >= 2.0, this section exists)
	if mdict.Meta.Version >= 2.0 {
		keyBlockInfoDecompressSizeBytes := keyBlockMetaBuffer[mdict.Meta.NumberWidth*2 : mdict.Meta.NumberWidth*2+mdict.Meta.NumberWidth]

		var keyBlockInfoDecompressSize uint64
		if mdict.Meta.NumberWidth == 8 {
			keyBlockInfoDecompressSize = beBinToU64(keyBlockInfoDecompressSizeBytes)
		} else if mdict.Meta.NumberWidth == 4 {
			keyBlockInfoDecompressSize = uint64(beBinToU32(keyBlockInfoDecompressSizeBytes))
		}
		keyBlockMeta.KeyBlockInfoDecompressSize = int64(keyBlockInfoDecompressSize)

		keyBlockInfoSizeBytesStartOffset = mdict.Meta.NumberWidth * 3

	} else {
		keyBlockInfoSizeBytesStartOffset = mdict.Meta.NumberWidth * 2
	}

	// 4. [24:32]([8:12]) - Key block info size
	keyBlockInfoSizeBytes := keyBlockMetaBuffer[keyBlockInfoSizeBytesStartOffset : keyBlockInfoSizeBytesStartOffset+mdict.Meta.NumberWidth]

	var keyBlockInfoSize uint64
	if mdict.Meta.NumberWidth == 8 {
		keyBlockInfoSize = beBinToU64(keyBlockInfoSizeBytes)
	} else if mdict.Meta.NumberWidth == 4 {
		keyBlockInfoSize = uint64(beBinToU32(keyBlockInfoSizeBytes))
	}

	keyBlockMeta.KeyBlockInfoCompressedSize = int64(keyBlockInfoSize)

	// 5. [32:40]([12:16]) - Key block size
	keyBlockDataSizeBytes := keyBlockMetaBuffer[keyBlockInfoSizeBytesStartOffset+mdict.Meta.NumberWidth : keyBlockInfoSizeBytesStartOffset+mdict.Meta.NumberWidth+mdict.Meta.NumberWidth]

	var keyBlockDataSize uint64
	if mdict.Meta.NumberWidth == 8 {
		keyBlockDataSize = beBinToU64(keyBlockDataSizeBytes)
	} else if mdict.Meta.NumberWidth == 4 {
		keyBlockDataSize = uint64(beBinToU32(keyBlockDataSizeBytes))
	}
	keyBlockMeta.KeyBlockDataTotalSize = int64(keyBlockDataSize)

	// 6. [40:44] - 4 bytes checksum (TODO: Skip if version > 2.0)
	// TODO checksum verification

	// Free key block info buffer
	if mdict.Meta.Version >= 2.0 {
		keyBlockMeta.KeyBlockInfoStartOffset = mdict.Meta.KeyBlockMetaStartOffset + 40 + 4
	} else {
		keyBlockMeta.KeyBlockInfoStartOffset = mdict.Meta.KeyBlockMetaStartOffset + 16
	}

	mdict.KeyBlockMeta = keyBlockMeta

	return nil
}

func (mdict *MdictBase) ReadKeyBlockInfo() error {
	file, err := os.Open(mdict.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer, err := readFileFromPos(file, mdict.KeyBlockMeta.KeyBlockInfoStartOffset, mdict.KeyBlockMeta.KeyBlockInfoCompressedSize)
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
	if mdict.Meta.EncryptType == EncryptKeyInfoEnc {
		// TODO decode key info
		keyBlockInfoDecryptedBuffer = mdxDecrypt(data, mdict.KeyBlockMeta.KeyBlockInfoCompressedSize)
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

	decompressKeyInfoBuffer, err := zlibDecompress(keyBlockInfoDecryptedBuffer, 8, mdict.KeyBlockMeta.KeyBlockInfoCompressedSize-8)
	if err != nil {
		return err
	}
	if int64(len(decompressKeyInfoBuffer)) != mdict.KeyBlockMeta.KeyBlockInfoDecompressSize {
		return errors.New("decoded key block info data size not equals to key block meta indicates key block info size")
	}

	// decode key-block entries
	var counter int64 = 0
	var currentEntriesSize int64 = 0
	var numEntriesCounter int64 = 0
	byteWidth := 1
	textTerm := 0

	if mdict.Meta.Version >= 2.0 {
		byteWidth = 2
		textTerm = 1
	}

	var dataOffset = 0
	var compressSizeAccumulator = 0
	var decompressSizeAccumulator = 0

	keyBlockInfo := &MDictKeyBlockInfo{
		KeyBlockEntriesStartOffset: 0,
		KeyBlockInfoList:           make([]*MDictKeyBlockInfoItem, 0),
	}

	for counter < mdict.KeyBlockMeta.KeyBlockNum {
		firstKeySize, lastKeySize := 0, 0
		firstKey := ""
		lastKey := ""

		if mdict.Meta.Version >= 2.0 {
			currentEntriesSize = int64(beBinToU64(decompressKeyInfoBuffer[dataOffset : dataOffset+mdict.Meta.NumberWidth]))
			dataOffset += mdict.Meta.NumberWidth
			firstKeySize = int(beBinToU16(decompressKeyInfoBuffer[dataOffset : dataOffset+byteWidth]))
			dataOffset += byteWidth
		} else {
			currentEntriesSize = int64(beBinToU32(decompressKeyInfoBuffer[dataOffset : dataOffset+mdict.Meta.NumberWidth]))
			dataOffset += mdict.Meta.NumberWidth
			firstKeySize = int(int64(beBinToU8(decompressKeyInfoBuffer[dataOffset : dataOffset+byteWidth])))
			dataOffset += byteWidth
		}
		numEntriesCounter += currentEntriesSize

		// step_gap means first key start data_offset to first key end;
		var stepGap = 0
		var termSize = textTerm
		if mdict.Meta.Encoding == EncodingUtf16 || mdict.FileType == MdictTypeMdd {
			stepGap = (firstKeySize + textTerm) * 2
			termSize = textTerm * 2
		} else {
			stepGap = firstKeySize + textTerm
			termSize = textTerm
		}

		firstKey = bigEndianBinToUTF8(decompressKeyInfoBuffer, dataOffset, stepGap-termSize)

		dataOffset += stepGap

		if mdict.Meta.Version >= 2.0 {
			lastKeySize = int(beBinToU16(decompressKeyInfoBuffer[dataOffset : dataOffset+byteWidth]))
		} else {
			lastKeySize = int(beBinToU8(decompressKeyInfoBuffer[dataOffset : dataOffset+byteWidth]))
		}
		dataOffset += byteWidth

		if mdict.Meta.Encoding == EncodingUtf16 || mdict.FileType == MdictTypeMdd {
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
		if mdict.Meta.Version >= 2.0 {
			keyBlockCompressSize = int(beBinToU64(decompressKeyInfoBuffer[dataOffset : dataOffset+mdict.Meta.NumberWidth]))
		} else {
			keyBlockCompressSize = int(beBinToU32(decompressKeyInfoBuffer[dataOffset : dataOffset+mdict.Meta.NumberWidth]))
		}
		dataOffset += mdict.Meta.NumberWidth

		keyBlockDecompressSize := 0
		if mdict.Meta.Version >= 2.0 {
			keyBlockDecompressSize = int(beBinToU64(decompressKeyInfoBuffer[dataOffset : dataOffset+mdict.Meta.NumberWidth]))
		} else {
			keyBlockDecompressSize = int(beBinToU32(decompressKeyInfoBuffer[dataOffset : dataOffset+mdict.Meta.NumberWidth]))
		}

		dataOffset += mdict.Meta.NumberWidth

		keyBlockInfoItem := &MDictKeyBlockInfoItem{
			FirstKey:                      firstKey,
			FirstKeySize:                  firstKeySize,
			LastKey:                       lastKey,
			LastKeySize:                   lastKeySize,
			KeyBlockInfoIndex:             int(counter),
			KeyBlockCompressSize:          int64(keyBlockCompressSize),
			KeyBlockCompAccumulator:       int64(compressSizeAccumulator),
			KeyBlockDeCompressSize:        int64(keyBlockDecompressSize),
			KeyBlockDeCompressAccumulator: int64(decompressSizeAccumulator),
		}

		compressSizeAccumulator += keyBlockCompressSize
		decompressSizeAccumulator += keyBlockDecompressSize

		keyBlockInfo.KeyBlockInfoList = append(keyBlockInfo.KeyBlockInfoList, keyBlockInfoItem)

		counter++

	}
	//keyBlockInfo.KeyBlockEntriesStartOffset = int64(dataOffset) + mdict.KeyBlockMeta.KeyBlockInfoStartOffset
	keyBlockInfo.KeyBlockEntriesStartOffset = mdict.KeyBlockMeta.KeyBlockInfoCompressedSize + mdict.KeyBlockMeta.KeyBlockInfoStartOffset

	mdict.KeyBlockInfo = keyBlockInfo

	if int64(compressSizeAccumulator) != mdict.KeyBlockMeta.KeyBlockDataTotalSize {
		return fmt.Errorf("key block data compress size not equals to meta key block data compress size(%d/%d)", compressSizeAccumulator, mdict.KeyBlockMeta.KeyBlockDataTotalSize)
	}

	return nil

}

func (mdict *MdictBase) ReadKeyEntries() error {
	file, err := os.Open(mdict.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer, err := readFileFromPos(file,
		mdict.KeyBlockInfo.KeyBlockEntriesStartOffset,
		mdict.KeyBlockMeta.KeyBlockDataTotalSize)
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

	keyBlockData := &MDictKeyBlockData{
		KeyEntries:                 make([]*MDictKeyBlockEntry, 0),
		KeyEntriesSize:             0,
		RecordBlockMetaStartOffset: 0,
	}

	for idx := 0; idx < len(mdict.KeyBlockInfo.KeyBlockInfoList); idx++ {

		compressedSize := mdict.KeyBlockInfo.KeyBlockInfoList[idx].KeyBlockCompressSize
		decompressedSize := mdict.KeyBlockInfo.KeyBlockInfoList[idx].KeyBlockDeCompressSize

		compAccu += mdict.KeyBlockInfo.KeyBlockInfoList[idx].KeyBlockCompressSize

		end = start + compressedSize

		if int64(start) != int64(mdict.KeyBlockInfo.KeyBlockInfoList[idx].KeyBlockCompAccumulator) {
			return fmt.Errorf("[%d] the key-block data start offset not equal to key block compress accumulator(%d/%d/%d)\n",
				idx, start, mdict.KeyBlockInfo.KeyBlockInfoList[idx].KeyBlockCompAccumulator, compAccu)
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

		keyBlockData.KeyEntries = append(keyBlockData.KeyEntries, splitKeys...)
		keyBlockData.KeyEntriesSize += int64(len(splitKeys))

		//fmt.Printf("idx(%05d)[start:%05d/end:%05d/comps:%05d->datalen:%05d/compaccu:%d]\n", idx, start, end, compressedSize, len(key_block), compAccu)
		//fmt.Printf("key_list %+v\n", splitKeys)

		start = end
	}

	if keyBlockData.KeyEntriesSize != mdict.KeyBlockMeta.EntriesNum {
		return errors.New("the key list items not equals to entries num")
	}
	keyBlockData.RecordBlockMetaStartOffset = mdict.KeyBlockInfo.KeyBlockEntriesStartOffset + mdict.KeyBlockMeta.KeyBlockDataTotalSize

	// keep key list in memory
	mdict.KeyBlockData = keyBlockData

	return nil
}

func (mdict *MdictBase) splitKeyBlock(keyBlock []byte) []*MDictKeyBlockEntry {
	// delimiter := ""
	width := 1

	if mdict.Meta.Encoding == EncodingUtf16 || mdict.FileType == MdictTypeMdd {
		//delimiter = "0000"
		width = 2
	} else {
		//delimiter = "00"
		width = 1
	}

	keyList := make([]*MDictKeyBlockEntry, 0)

	keyStartIndex := 0
	keyEndIndex := 0

	for keyStartIndex < len(keyBlock) {
		// # the corresponding record's offset in record block
		recordStartOffset := int64(0)

		if mdict.Meta.NumberWidth == 8 {
			recordStartOffset = int64(beBinToU64(keyBlock[keyStartIndex : keyStartIndex+mdict.Meta.NumberWidth]))
		} else {
			recordStartOffset = int64(beBinToU32(keyBlock[keyStartIndex : keyStartIndex+mdict.Meta.NumberWidth]))
		}

		// # key text ends with '\x00'
		i := keyStartIndex + mdict.Meta.NumberWidth
		for i < len(keyBlock) {
			// delimiter = '0' || // delimiter = '00'
			if (width == 1 && keyBlock[i] == 0) || (width == 2 && keyBlock[i] == 0 && keyBlock[i+1] == 0) {
				keyEndIndex = i
				break
			}
			i += width
		}

		keyTextBytes := keyBlock[keyStartIndex+mdict.Meta.NumberWidth : keyEndIndex]
		keyText := string(keyTextBytes)
		var err error

		if mdict.Meta.Encoding == EncodingUtf16 {
			keyText, err = decodeLittleEndianUtf16(keyTextBytes)
			if err != nil {
				keyText = string(keyTextBytes)
			}
		}

		if mdict.FileType == MdictTypeMdd {
			keyText, err = decodeLittleEndianUtf16(keyTextBytes)
			if err != nil {
				panic(err)
			}
		}

		keyStartIndex = keyEndIndex + width
		keyList = append(keyList, &MDictKeyBlockEntry{
			RecordStartOffset: recordStartOffset,
			KeyWord:           keyText,
			KeyBlockIdx:       int64(keyStartIndex),
		})
		if len(keyList) > 1 {
			keyList[len(keyList)-2].RecordEndOffset = keyList[len(keyList)-1].RecordStartOffset
		}
	}
	//keyList[len(keyList)-1].RecordEndOffset = 0

	return keyList
}

func (mdict *MdictBase) ReadRecordBlockMeta() error {
	file, err := os.Open(mdict.FilePath)
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
	if mdict.Meta.Version >= 2.0 {
		recordBlockMetaBufferLen = 32
	}

	recordBlockStartOffset := mdict.KeyBlockInfo.KeyBlockEntriesStartOffset + mdict.KeyBlockMeta.KeyBlockDataTotalSize

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
	recordBlockMeta := &MDictRecordBlockMeta{
		KeyRecordMetaStartOffset: startOffset,
		KeyRecordMetaEndOffset:   endOffset,
	}

	keyRecordBuffer := data
	offset := 0

	if mdict.Meta.Version >= 2.0 {
		recordBlockMeta.RecordBlockNum = int64(beBinToU64(keyRecordBuffer[offset : offset+mdict.Meta.NumberWidth]))
	} else {
		recordBlockMeta.RecordBlockNum = int64(beBinToU32(keyRecordBuffer[offset : offset+mdict.Meta.NumberWidth]))
	}

	offset += mdict.Meta.NumberWidth

	if mdict.Meta.Version >= 2.0 {
		recordBlockMeta.EntriesNum = int64(beBinToU64(keyRecordBuffer[offset : offset+mdict.Meta.NumberWidth]))
	} else {
		recordBlockMeta.EntriesNum = int64(beBinToU32(keyRecordBuffer[offset : offset+mdict.Meta.NumberWidth]))

	}
	if recordBlockMeta.EntriesNum != mdict.KeyBlockMeta.EntriesNum {
		return fmt.Errorf("keyEntriesNum != meta.EntriesNum")
	}

	offset += mdict.Meta.NumberWidth
	if mdict.Meta.Version >= 2.0 {
		recordBlockMeta.RecordBlockInfoCompSize = int64(beBinToU64(keyRecordBuffer[offset : offset+mdict.Meta.NumberWidth]))
	} else {
		recordBlockMeta.RecordBlockInfoCompSize = int64(beBinToU32(keyRecordBuffer[offset : offset+mdict.Meta.NumberWidth]))
	}

	offset += mdict.Meta.NumberWidth

	if mdict.Meta.Version >= 2.0 {
		recordBlockMeta.RecordBlockCompSize = int64(beBinToU64(keyRecordBuffer[offset : offset+mdict.Meta.NumberWidth]))
	} else {
		recordBlockMeta.RecordBlockCompSize = int64(beBinToU32(keyRecordBuffer[offset : offset+mdict.Meta.NumberWidth]))
	}

	mdict.RecordBlockMeta = recordBlockMeta
	return nil
}

func (mdict *MdictBase) ReadRecordBlockInfo() error {
	file, err := os.Open(mdict.FilePath)
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
	recordBlockInfoStartOffset := mdict.RecordBlockMeta.KeyRecordMetaEndOffset
	recordBlockInfoLen := mdict.RecordBlockMeta.RecordBlockInfoCompSize

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

	recordBlockInfoList := make([]*MDictRecordBlockInfoListItem, 0)
	var offset = 0
	var compAccu = int64(0)
	var decompAccu = int64(0)
	var i = int64(0)

	for i = int64(0); i < mdict.RecordBlockMeta.RecordBlockNum; i++ {
		compSize := int64(0)
		if mdict.Meta.Version >= 2.0 {
			compSize = int64(beBinToU64(data[offset : offset+mdict.Meta.NumberWidth]))
		} else {
			compSize = int64(beBinToU32(data[offset : offset+mdict.Meta.NumberWidth]))
		}
		offset += mdict.Meta.NumberWidth

		decompSize := int64(0)
		if mdict.Meta.Version >= 2.0 {
			decompSize = int64(beBinToU64(data[offset : offset+mdict.Meta.NumberWidth]))
		} else {
			decompSize = int64(beBinToU32(data[offset : offset+mdict.Meta.NumberWidth]))
		}
		offset += mdict.Meta.NumberWidth

		// then assign
		recordBlockInfoList = append(recordBlockInfoList, &MDictRecordBlockInfoListItem{
			CompressSize:                compSize,
			DeCompressSize:              decompSize,
			CompressAccumulatorOffset:   compAccu,
			DeCompressAccumulatorOffset: decompAccu,
		})

		// accu last
		compAccu += compSize
		decompAccu += decompSize
	}
	if int64(i) != mdict.RecordBlockMeta.RecordBlockNum {
		return fmt.Errorf("RecordBlockInfo (i) not equals to meta.RecordBlockNum [%d/%d] compA/decompA(%d/%d)", i, mdict.RecordBlockMeta.RecordBlockNum, compAccu, decompAccu)
	}
	if int64(offset) != mdict.RecordBlockMeta.RecordBlockInfoCompSize {
		return errors.New("RecordBlockInfo offset not equals to meta.RecordBlockInfoCompSize")
	}
	if int64(compAccu) != mdict.RecordBlockMeta.RecordBlockCompSize {
		return errors.New("RecordBlockInfo compAccu not equals to meta.RecordBlockCompSize")
	}

	recordBlockInfo := &MDictRecordBlockInfo{
		RecordInfoList:             recordBlockInfoList,
		RecordBlockInfoStartOffset: startOffset,
		RecordBlockInfoEndOffset:   endOffset,
		RecordBlockDataStartOffset: endOffset,
	}

	mdict.RecordBlockInfo = recordBlockInfo

	return nil
}

func (mdict *MdictBase) ReadRecordBlockData() error {
	file, err := os.Open(mdict.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	recordBlockInfoStartOffset := mdict.RecordBlockInfo.RecordBlockDataStartOffset

	recordBlockInfoLen := mdict.RecordBlockInfo.RecordInfoList[len(mdict.RecordBlockInfo.RecordInfoList)-1].CompressAccumulatorOffset

	buffer, err := readFileFromPos(file, recordBlockInfoStartOffset, recordBlockInfoLen)
	if err != nil {
		return err
	}

	err = mdict.decodeRecordBlockData(buffer, recordBlockInfoStartOffset, recordBlockInfoStartOffset+recordBlockInfoLen)
	if err != nil {
		return err
	}
	return nil
}

// TODO FIX THIS, 这个方法尚未实现
func (mdict *MdictBase) decodeRecordBlockData(data []byte, startOffset, endOffset int64) error {

	var sizeCounter = int64(0)
	var itemCounter = int64(0)
	var recordBlockOffset = int64(0)
	var recordCompType = "NONE"

	var recordBlockData = &MDictRecordBlockData{
		RecordBlockStartOffset: startOffset,
		RecordBlockEndOffset:   endOffset,
		RecordItemList:         make([]*MDictRecordDataItem, 0),
	}

	for idx := 0; idx < len(mdict.RecordBlockInfo.RecordInfoList); idx++ {
		var recordBlockCompSize = mdict.RecordBlockInfo.RecordInfoList[idx].CompressSize
		var recordBlockDecompSize = mdict.RecordBlockInfo.RecordInfoList[idx].DeCompressSize
		var recordBlockDataCompBuff = data[recordBlockOffset : recordBlockOffset+recordBlockCompSize]
		var recordBlockCompressAccumulator = mdict.RecordBlockInfo.RecordInfoList[idx].CompressAccumulatorOffset
		var recordBlockDeCompressAccumulator = mdict.RecordBlockInfo.RecordInfoList[idx].DeCompressAccumulatorOffset

		// 4 bytes: compression type
		var rbCompType = recordBlockDataCompBuff[0:4]

		// record_block stores the final record data
		var recordBlock []byte

		// TODO: igore adler32 offset
		// Note: here ignore the checksum part
		// bytes: adler32 checksum of decompressed record block
		// adler32 = unpack('>I', record_block_compressed[4:8])[0]
		if rbCompType[0] == 0 {
			recordBlock = recordBlockDataCompBuff[8 : 8+recordBlockCompSize]
			recordCompType = "NONE"
		} else {
			// decrypt
			var blockBufDecrypted []byte
			// if encrypt type == 1, the record block was encrypted
			if mdict.Meta.EncryptType == EncryptRecordEnc {
				// const passkey = new Uint8Array(8);
				// record_block_compressed.copy(passkey, 0, 4, 8);
				// passkey.set([0x95, 0x36, 0x00, 0x00], 4); // key part 2: fixed data
				blockBufDecrypted = mdxDecrypt(recordBlockDataCompBuff, recordBlockCompSize)
			} else {
				blockBufDecrypted = recordBlockDataCompBuff[8:recordBlockCompSize]
			}

			// decompress
			if rbCompType[0] == 1 {
				// TODO the second part
				recordCompType = "LZOX1"
				header := []byte{0xf0, byte(int(recordBlockDecompSize))}
				// # decompress key block
				reader := bytes.NewReader(append(header, blockBufDecrypted...))

				out, err1 := lzo.Decompress1X(reader, 0, 0 /* decompressedSize, 1308672*/)
				if err1 != nil {
					return err1
				}

				recordBlock = out

			} else if rbCompType[0] == 2 {
				recordCompType = "ZLIB"
				var err error
				recordBlock, err = zlibDecompress(blockBufDecrypted, 0, int64(len(blockBufDecrypted)))
				if err != nil {
					return err
				}
			}
		}

		// notice that adler32 return signed value
		// TODO: ignore the checksum
		// assert(adler32 == zlib.adler32(record_block) & 0xffffffff)

		if int64(len(recordBlock)) != recordBlockDecompSize {
			return errors.New("recordBlock length not equals decompress Size")
		}

		/**
		 * 请注意，block 是会有很多个的，而每个block都可能会被压缩
		 * 而 key_list中的 record_start, key_text是相对每一个block而言的，end是需要每次解析的时候算出来的
		 * 所有的record_start/length/end都是针对解压后的block而言的
		 */

		// split record block according to the offset info from key block
		//var offset = int64(0)
		keyblockEntry := mdict.KeyBlockData.KeyEntries
		for i := 0; i < len(keyblockEntry); i++ {

			var recordEntryStart = keyblockEntry[i].RecordStartOffset
			var keyText = keyblockEntry[i].KeyWord
			// # reach the end of current record block
			//if recordEntryStart-offset >= int64(len(recordBlock)) {
			//	break
			//}
			// # record end index
			var recordEntryEnd int64 = 0
			if i < len(keyblockEntry)-1 {
				recordEntryEnd = int64(keyblockEntry[i+1].RecordStartOffset)
			} else {
				//recordEntryEnd = int64(len(mdict.KeyBlockData.KeyList)) + offset
				recordEntryEnd = endOffset
			}

			//fmt.Printf("keyText: %s, recordEntryStart: %d, recordEntryEnd:%d\n", keyText, recordEntryStart, recordEntryEnd)

			recordBlockData.RecordItemList = append(recordBlockData.RecordItemList,
				&MDictRecordDataItem{
					KeyWord:                          keyText,
					RecordBlockCompressStart:         recordBlockOffset,
					RecordBlockCompressEnd:           recordBlockOffset + recordBlockCompSize,
					RecordBlockCompressSize:          recordBlockCompSize,
					RecordBlockDeCompressSize:        recordBlockDecompSize,
					RecordBlockCompressType:          recordCompType,
					RecordBlockEncrypted:             mdict.Meta.EncryptType == EncryptRecordEnc,
					RecordBlockFileRelativeOffset:    recordBlockOffset + mdict.RecordBlockInfo.RecordBlockDataStartOffset,
					RecordBlockCompressAccumulator:   recordBlockCompressAccumulator,
					RecordBlockDeCompressAccumulator: recordBlockDeCompressAccumulator,

					RecordEntryIndex:           int64(itemCounter),
					RecordEntryDecompressStart: recordEntryStart,
					RecordEntryDecompressEnd:   recordEntryEnd,

					RecordInfoIndex: int64(idx),
				})

			itemCounter++
		}

		recordBlockOffset += recordBlockCompSize
		//offset += int64(len(recordBlock))
		sizeCounter += recordBlockCompSize
	}

	if sizeCounter != mdict.RecordBlockMeta.RecordBlockCompSize {
		return fmt.Errorf("record entries compressed size sum not equals record block compress size(%d/%d)", sizeCounter, mdict.RecordBlockMeta.RecordBlockCompSize)
	}

	mdict.RecordBlockData = recordBlockData

	return nil
}

func (mdict *MdictBase) LocateRecordDefinition(item *MDictKeyBlockEntry) ([]byte, error) {
	file, err := os.Open(mdict.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var recordBlockInfo *MDictRecordBlockInfoListItem

	var i = 0
	for ; i < len(mdict.RecordBlockInfo.RecordInfoList)-1; i++ {
		curr := mdict.RecordBlockInfo.RecordInfoList[i]
		next := mdict.RecordBlockInfo.RecordInfoList[i+1]
		if item.RecordStartOffset >= curr.DeCompressAccumulatorOffset && item.RecordStartOffset < next.DeCompressAccumulatorOffset {
			recordBlockInfo = curr
			break
		}
	}

	// the last one
	if i == len(mdict.RecordBlockInfo.RecordInfoList)-1 {
		lastOne := mdict.RecordBlockInfo.RecordInfoList[len(mdict.RecordBlockInfo.RecordInfoList)-1]
		if item.RecordStartOffset < lastOne.DeCompressAccumulatorOffset+lastOne.DeCompressSize {
			recordBlockInfo = lastOne
		}
	}

	if recordBlockInfo == nil {
		fmt.Printf("record block info is nil, current keyBlockEntry: %+v, last RecordBlockInfo: %+v\n", item, mdict.RecordBlockInfo.RecordInfoList[len(mdict.RecordBlockInfo.RecordInfoList)-1])
		return nil, errors.New("key-item record info not found")
	}

	recordBlockStartOffset := recordBlockInfo.CompressAccumulatorOffset + mdict.RecordBlockInfo.RecordBlockDataStartOffset
	recordBlockLen := recordBlockInfo.CompressSize

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
		recordBlock = recordBlockDataCompBuff[8:recordBlockInfo.CompressSize]
	} else {
		// decrypt
		var blockBufDecrypted []byte
		// if encrypt type == 1, the record block was encrypted
		if mdict.Meta.EncryptType == EncryptRecordEnc {
			// const passkey = new Uint8Array(8);
			// record_block_compressed.copy(passkey, 0, 4, 8);
			// passkey.set([0x95, 0x36, 0x00, 0x00], 4); // key part 2: fixed data
			blockBufDecrypted = mdxDecrypt(recordBlockDataCompBuff, recordBlockInfo.CompressSize)
		} else {
			blockBufDecrypted = recordBlockDataCompBuff[8:recordBlockInfo.CompressSize]
		}

		// decompress
		if rbCompType[0] == 1 {
			// TODO the second part
			header := []byte{0xf0, byte(int(recordBlockInfo.CompressSize))}
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

	if int64(len(recordBlock)) != recordBlockInfo.DeCompressSize {
		return nil, errors.New("recordBlock length not equals decompress Size")
	}

	start := item.RecordStartOffset - recordBlockInfo.DeCompressAccumulatorOffset
	var end int64
	if item.RecordEndOffset == 0 {
		end = int64(len(recordBlock))
	} else {
		end = item.RecordEndOffset - recordBlockInfo.DeCompressAccumulatorOffset
	}

	data := recordBlock[start:end]

	if mdict.FileType == MdictTypeMdd {
		return data, nil
	}

	if mdict.Meta.Encoding == EncodingUtf16 {
		datastr, err1 := decodeLittleEndianUtf16(data)
		if err1 != nil {
			return nil, err
		}
		return []byte(datastr), nil
	}
	return data, nil

}
