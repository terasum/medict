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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadMDictFile(t *testing.T) {
	mdict, err := readMDictFileHeader("testdata/dict/testdict.mdx")
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 2030, mdict.HeaderBytesSize)
	t.Log("Header Bytes Size:", mdict.HeaderBytesSize)
	assert.Equal(t, "<Dictionary GeneratedByEngineVersion=\"2.0\" RequiredEngineVersion=\"2.0\" Format=\"Html\" KeyCaseSensitive=\"No\" StripKey=\"Yes\" Encrypted=\"2\" RegisterBy=\"EMail\" Description=\"&lt;font size=5&gt;\r\n&lt;b&gt;The New Oxford Thesaurus of English&lt;/b&gt;\r\n&lt;br&gt;\r\n&lt;br&gt;&lt;b&gt;Author: &lt;/b&gt;OUP\r\n&lt;br&gt;\r\n&lt;br&gt;&lt;b&gt;Description: &lt;/b&gt;This new edition of the Oxford Thesaurus provides comprehensive, reliable help for all writing needs and crossword solving. The synonyms and antonyms listed include unusual and specialist vocabulary as well as everyday words, and there is clear labelling of informal, dialect, literary and technical items. There are also example sentences, and thematic lists provide information on a wide range of subjects. \r\n&lt;br&gt;\r\n&lt;br&gt;&lt;b&gt;No. of Entries: &lt;/b&gt;15,774.\r\n&lt;/font&gt;\" Title=\"The New Oxford Thesaurus of English\" Encoding=\"UTF-8\" CreationDate=\"2009-12-27\" Compact=\"No\" Compat=\"No\" Left2Right=\"Yes\" DataSourceFormat=\"107\" StyleSheet=\"\"/>\r\n\u0000",
		mdict.HeaderInfo)
	t.Log("Header Info:", mdict.HeaderInfo)
	assert.Equal(t, 2992291341, int(mdict.Adler32Checksum))
	//fmt.Println("Adler32 Checksum:", int(mdict.Adler32Checksum))
}

func TestReadMDictFile2(t *testing.T) {
	mdict, err := readMDictFileHeader("testdata/dict/wlghyzd2000.mdx")
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 5552, int(mdict.HeaderBytesSize))
	t.Log("Header Bytes Size:", mdict.HeaderBytesSize)
	t.Log("Header Info Bytes:", string(littleEndianBinUTF16ToUTF8(mdict.HeaderInfoBytes, 0, int(mdict.HeaderBytesSize))))
	assert.Equal(t, 3421787354, int(mdict.Adler32Checksum))
	t.Log("Adler32 Checksum:", int(mdict.Adler32Checksum))
}

func TestMdictBase_ReadDictHeader(t *testing.T) {
	mdictBase := &MdictBase{
		FilePath: "testdata/dict/testdict.mdx",
	}
	err := mdictBase.ReadDictHeader()
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
		FilePath: "testdata/dict/testdict.mdx",
	}
	err := mdictBase.ReadDictHeader()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.ReadKeyBlockMeta()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.ReadKeyBlockInfo()
	if err != nil {
		t.Fatal(err)
	}

	jsondata, err := json.MarshalIndent(mdictBase, "", "  ")
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s\n", jsondata)

	//t.Logf("Dictionary Header keyBlockStartOffset %d / Meta KeyBlockHeaderStartOffset %d\n", mdictBase.Header.KeyBlockOffset, mdictBase.Meta.KeyBlockHeaderStartOffset)
}
func TestMdictBase_ReadDictHeader3(t *testing.T) {
	mdictBase := &MdictBase{
		FilePath: "testdata/dict/oale8.mdx",
	}
	err := mdictBase.ReadDictHeader()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.ReadKeyBlockMeta()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.ReadKeyBlockInfo()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.ReadKeyEntries()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.ReadRecordBlockMeta()
	if err != nil {
		t.Fatal(err)
	}

	err = mdictBase.ReadRecordBlockInfo()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("key entries list len: %d, record block info entry list len %d", len(mdictBase.KeyBlockData.KeyEntries), len(mdictBase.RecordBlockInfo.RecordInfoList))
	t.Logf("entries number size %d\n", mdictBase.KeyBlockData.KeyEntriesSize)
	t.Logf("keylist[0] %+v\n", mdictBase.KeyBlockData.KeyEntries[0])

	item := mdictBase.KeyBlockData.KeyEntries[0]

	data, err := mdictBase.LocateRecordDefinition(item)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("0-0 keyText: %s, data: %s", item.KeyWord, data)

	item = mdictBase.KeyBlockData.KeyEntries[1]

	data, err = mdictBase.LocateRecordDefinition(item)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("13-0 keyText: %s, data: %s", item.KeyWord, data)

	item = mdictBase.KeyBlockData.KeyEntries[3]

	data, err = mdictBase.LocateRecordDefinition(item)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("13-7 keyText: %s, data: %s", item.KeyWord, data)

}

//
///*
//#include "mdict_extern.h"
//#include "mdict.h"
//#include <stdlib.h>
//*/
//import "C"
//
//import (
//	"unsafe"
//)
//
//type Mdict struct {
//	dict unsafe.Pointer
//}
//
//type SimpleKeyItem struct {
//	keyWord string
//}
//
//func MdictInit(dictionaryPath string) *Mdict {
//	dictFilePath := C.CString(dictionaryPath)
//	defer C.free(unsafe.Pointer(dictFilePath))
//	mydict := C.mdict_init(dictFilePath)
//	return &Mdict{dict: mydict}
//}
//
//func (m *Mdict) Lookup(word string) string {
//	queryWord := C.CString(word)
//	defer C.free(unsafe.Pointer(queryWord))
//	var result *C.char
//	C.mdict_lookup(m.dict, queryWord, &result)
//	defer C.free(unsafe.Pointer(result))
//	return C.GoString(result)
//}
//
//func (m *Mdict) ParseDefinition(word string, recordStart uint64) string {
//	queryWord := C.CString(word)
//	defer C.free(unsafe.Pointer(queryWord))
//	var result *C.char
//	C.mdict_parse_definition(m.dict, queryWord, C.ulong(recordStart), &result)
//	defer C.free(unsafe.Pointer(result))
//	return C.GoString(result)
//}
//
//func (m *Mdict) KeyList() []*SimpleKeyItem {
//	var len C.ulong
//	keylist := C.mdict_keylist(m.dict, &len)
//	defer C.free(unsafe.Pointer(keylist))
//	items := make([]*SimpleKeyItem, len)
//	for i := C.ulong(0); i < len; i++ {
//		items[i] = &SimpleKeyItem{keyWord: C.GoString((*C.simple_key_item)(unsafe.Pointer(uintptr(unsafe.Pointer(keylist)) + uintptr(i)*unsafe.Sizeof(*keylist))).key_word)}
//	}
//	return items
//}
//
//func main() {
//	// Example usage
//	dictionaryPath := "path/to/dictionary"
//	mydict := MdictInit(dictionaryPath)
//	result := mydict.Lookup("word")
//	println(result)
//}
