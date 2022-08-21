package mdictparser

/*
#cgo CFLAGS: -I../libmdict -I.
#cgo CXXFLAGS: -stdlib=libc++
#cgo LDFLAGS: -L../libmdict/build/lib -L. -lmdict -lminilzo -lminiz -lstdc++
#cgo arm64 darwin CFLAGS: -DX86=1

#include <stdlib.h>
#include "mdict_extern.h"
#include "mdict_simple_key.h"

char* fetch_keyword(simple_key_item** item_list, unsigned long idx) {
	return item_list[idx]->key_word;
}

unsigned long fetch_record_start(simple_key_item** item_list, unsigned long idx) {
	return item_list[idx]->record_start;
}

*/
import "C"
import (
	"errors"
	"github.com/terasum/medict/internal/mdictparser/def"
	"unsafe"
)

type RawDict struct {
	pointer unsafe.Pointer
}

func LoadDict(dictName string) (*RawDict, error) {
	cname := C.CString(dictName)
	dict := C.mdict_init(cname)
	if dict == nil {
		return nil, errors.New("load dictionary failed")
	}

	return &RawDict{
		pointer: unsafe.Pointer(dict),
	}, nil
}

func LookUpDict(dict *RawDict, word string) string {
	cword := C.CString(word)

	result := C.CString("")
	C.mdict_lookup(dict.pointer, cword, &result)
	resultStr := C.GoString(result)

	goResult := make([]rune, len(resultStr))
	copy(goResult, []rune(resultStr[:]))
	C.free(unsafe.Pointer(result))

	return string(goResult)
}

func DictType(dict *RawDict) string {
	if C.mdict_filetype(dict.pointer) == C.int(0) {
		return "MDX"
	} else {
		return "MDD"
	}
}

func AllKeyList(dict *RawDict) ([]*def.KeyItem, uint64, error) {
	cKeyListLen := (C.ulong)(uint64(0))
	keyListArray := C.mdict_keylist(dict.pointer, &cKeyListLen)
	goKeyListLen := (uint64)(cKeyListLen)

	finalList := make([]*def.KeyItem, goKeyListLen)

	for i := uint64(0); i < goKeyListLen; i++ {
		rawKeyWord := C.GoString(C.fetch_keyword(keyListArray, C.ulong(i)))
		finalList[i] = &def.KeyItem{
			KeyWord:     rawKeyWord,
			RecordStart: uint64(C.fetch_record_start(keyListArray, C.ulong(i))),
		}
	}

	return finalList, goKeyListLen, nil
}

func ParseDefinition(dict *RawDict, word string, recordStart uint64) string {
	defResult := C.CString("")
	C.mdict_parse_definition(dict.pointer, C.CString(word), C.ulong(recordStart), &defResult)

	resultStr := C.GoString(defResult)

	goResult := make([]rune, len(resultStr))
	copy(goResult, []rune(resultStr[:]))
	C.free(unsafe.Pointer(defResult))

	return string(goResult)
}

func DestroyDict(dict *RawDict) error {
	ret := C.mdict_destory(dict.pointer)
	if ret != 0 {
		return errors.New("destroy failed")
	}
	return nil
}
