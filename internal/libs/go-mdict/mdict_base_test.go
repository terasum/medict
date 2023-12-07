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

import (
	"bytes"
	"encoding/json"
	"github.com/rasky/go-lzo"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadMDictFile(t *testing.T) {
	mdict, err := readMDictFileHeader("testdata/mdx/testdict.mdx")
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, uint32(1190), mdict.headerBytesSize)
	t.Log("header Bytes Size:", mdict.headerBytesSize)
	assert.Equal(t, "<Dictionary GeneratedByEngineVersion=\"2.0\" RequiredEngineVersion=\"2.0\" Format=\"Html\" KeyCaseSensitive=\"No\" StripKey=\"Yes\" Encrypted=\"2\" RegisterBy=\"EMail\" Description=\"Oxford Advanced Learner’s English-Chinese Dictionary Eighth edition Based on Langheping&apos;s version Modified by EarthWorm&lt;br/&gt;\r\nHeadwords: 41969 &lt;br/&gt;\r\nEntries: 109473 &lt;br/&gt;\r\nVersion: 3.0.0 &lt;br/&gt;\r\nDate: 2018.02.18 &lt;br/&gt;\r\nLast Modified By roamlog&lt;br/&gt;\" Title=\"\" IsUTF16=\"UTF-8\" CreationDate=\"2018-2-18\" Compact=\"Yes\" Compat=\"Yes\" Left2Right=\"Yes\" DataSourceFormat=\"106\" StyleSheet=\"\"/>\r\n\u0000",
		mdict.headerInfo)
	t.Log("header Info:", mdict.headerInfo)
	assert.Equal(t, 3301029905, int(mdict.adler32Checksum))
	//fmt.Println("Adler32 Checksum:", int(mdict.adler32Checksum))
}

func TestReadMDictFile2(t *testing.T) {
	mdict, err := readMDictFileHeader("testdata/dict/wlghyzd2000.mdx")
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 5552, int(mdict.headerBytesSize))
	t.Log("header Bytes Size:", mdict.headerBytesSize)
	t.Log("header Info Bytes:", string(littleEndianBinUTF16ToUTF8(mdict.headerInfoBytes, 0, int(mdict.headerBytesSize))))
	assert.Equal(t, 3421787354, int(mdict.adler32Checksum))
	t.Log("Adler32 Checksum:", int(mdict.adler32Checksum))
}

func TestMdictBase_ReadDictHeader(t *testing.T) {
	mdictBase := &MdictBase{
		filePath: "testdata/dict/testdict.mdx",
	}
	err := mdictBase.readDictHeader()
	if err != nil {
		t.Error(err)
	}
	jsondata, err := json.MarshalIndent(mdictBase, "", "  ")
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s\n", jsondata)
}

func TestMdictBase_ReadDictHeader2(t *testing.T) {
	mdictBase := &MdictBase{
		filePath: "testdata/dict/testdict.mdx",
	}
	err := mdictBase.readDictHeader()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.readKeyBlockMeta()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.readKeyBlockInfo()
	if err != nil {
		t.Fatal(err)
	}

	jsondata, err := json.MarshalIndent(mdictBase, "", "  ")
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s\n", jsondata)

	//t.Logf("Dictionary header keyBlockStartOffset %d / meta KeyBlockHeaderStartOffset %d\n", mdictBase.header.KeyBlockOffset, mdictBase.meta.KeyBlockHeaderStartOffset)
}
func TestMdictBase_ReadDictHeader3(t *testing.T) {
	mdictBase := &MdictBase{
		filePath: "testdata/dict/oale8.mdx",
	}
	err := mdictBase.readDictHeader()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.readKeyBlockMeta()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.readKeyBlockInfo()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.readKeyEntries()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.readRecordBlockMeta()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.readRecordBlockInfo()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("key entries list len: %d, record block info entry list len %d", len(mdictBase.keyBlockData.keyEntries), len(mdictBase.recordBlockInfo.recordInfoList))
	t.Logf("entries number size %d\n", mdictBase.keyBlockData.keyEntriesSize)
	t.Logf("keylist[0] %+v\n", mdictBase.keyBlockData.keyEntries[0])

	item := mdictBase.keyBlockData.keyEntries[0]

	data, err := mdictBase.locateByKeywordEntry(item)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("0-0 keyText: %s, data: %s", item.KeyWord, data)

	item = mdictBase.keyBlockData.keyEntries[1]

	data, err = mdictBase.locateByKeywordEntry(item)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("13-0 keyText: %s, data: %s", item.KeyWord, data)

	item = mdictBase.keyBlockData.keyEntries[3]

	data, err = mdictBase.locateByKeywordEntry(item)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("13-7 keyText: %s, data: %s", item.KeyWord, data)

}

func TestMdictBase_ReadDictFixBug1(t *testing.T) {
	mdictBase := &MdictBase{
		filePath: "testdata/bugdict/教育部重編國語辭典(第五版)/教育部重編國語辭典(第五版).mdx",
	}
	err := mdictBase.readDictHeader()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.readKeyBlockMeta()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.readKeyBlockInfo()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.readKeyEntries()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.readRecordBlockMeta()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.readRecordBlockInfo()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("key entries list len: %d, record block info entry list len %d", len(mdictBase.keyBlockData.keyEntries), len(mdictBase.recordBlockInfo.recordInfoList))
	t.Logf("entries number size %d\n", mdictBase.keyBlockData.keyEntriesSize)
	t.Logf("keylist[0] %+v\n", mdictBase.keyBlockData.keyEntries[0])

	item := mdictBase.keyBlockData.keyEntries[0]

	data, err := mdictBase.locateByKeywordEntry(item)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("0-0 keyText: %s, data: %s", item.KeyWord, data)

	item = mdictBase.keyBlockData.keyEntries[1]

	data, err = mdictBase.locateByKeywordEntry(item)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("13-0 keyText: %s, data: %s", item.KeyWord, data)

	item = mdictBase.keyBlockData.keyEntries[3]

	data, err = mdictBase.locateByKeywordEntry(item)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("13-7 keyText: %s, data: %s", item.KeyWord, data)

}

func TestExtractContentByRecordIndex(t *testing.T) {
	keyWord := "a"
	//recordStartOffset := 30239433
	//recordEndOffset := 30255629
	recordBlockDataStartOffset := int64(4676923)
	recordBlockDataLen := int64(10013)
	myFile, err := os.Open("testdata/mdx/testdict.mdx")
	if err != nil {
		t.Fatal(err)
	}
	defer myFile.Close()

	data, err := testExtractData(myFile, recordBlockDataStartOffset, recordBlockDataLen, false)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s, data %s", keyWord, string(data))
	//t.Logf("%s, data %s", keyWord, string(data[recordStartOffset:recordEndOffset]))

}

func testExtractData(file *os.File, recordBlockStartOffset, CompressSize int64, EncryptRecordEnc bool) ([]byte, error) {
	recordBlockDataCompBuff, err := readFileFromPos(file, recordBlockStartOffset, CompressSize)
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
		recordBlock = recordBlockDataCompBuff[8:CompressSize]
		return recordBlock, nil
	}

	// decrypt
	var blockBufDecrypted []byte
	// if encrypt type == 1, the record block was encrypted
	if EncryptRecordEnc {
		// const passkey = new Uint8Array(8);
		// record_block_compressed.copy(passkey, 0, 4, 8);
		// passkey.set([0x95, 0x36, 0x00, 0x00], 4); // key part 2: fixed data
		blockBufDecrypted = mdxDecrypt(recordBlockDataCompBuff, CompressSize)
	} else {
		blockBufDecrypted = recordBlockDataCompBuff[8:CompressSize]
	}

	// decompress
	if rbCompType[0] == 1 {
		// TODO the second part
		header := []byte{0xf0, byte(int(CompressSize))}
		// # decompress key block
		reader := bytes.NewReader(append(header, blockBufDecrypted...))

		out, err1 := lzo.Decompress1X(reader, 0, 0 /* decompressedSize, 1308672*/)
		if err1 != nil {
			return nil, err1
		}
		return out, nil

	} else if rbCompType[0] == 2 {
		var err2 error
		recordBlock, err2 = zlibDecompress(blockBufDecrypted, 0, int64(len(blockBufDecrypted)))
		if err2 != nil {
			return nil, err2
		}

		return recordBlock, nil
	}

	return recordBlock, nil
}

//func testExtraceDefData(recordBlock []byte, DeCompress) {
//	// TODO: ignore the checksum
//	// notice that adler32 return signed value
//	// assert(adler32 == zlib.adler32(record_block) & 0xffffffff)
//
//	if int64(len(recordBlock)) != recordBlockInfo.RecordBlockDeCompressSize {
//		return nil, errors.New("recordBlock length not equals decompress Size")
//	}
//
//	start := item.RecordLocateStartOffset - recordBlockInfo.deCompressAccumulatorOffset
//	var end int64
//	if item.RecordLocateEndOffset == 0 {
//		end = int64(len(recordBlock))
//	} else {
//		end = item.RecordLocateEndOffset - recordBlockInfo.deCompressAccumulatorOffset
//	}
//
//	data := recordBlock[start:end]
//
//	if mdict.fileType == MdictTypeMdd {
//		return data, nil
//	}
//
//	if mdict.meta.encoding == EncodingUtf16 {
//		datastr, err1 := decodeLittleEndianUtf16(data)
//		if err1 != nil {
//			return nil, err
//		}
//		return []byte(datastr), nil
//	}
//	return data, nil
//}
