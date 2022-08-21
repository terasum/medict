package model

import (
	"github.com/agatan/bktree"
	"github.com/creasty/go-levenshtein"
)

type WordItem struct {
	RawKeyWord  string `json:"raw_key_word"`
	KeyWord     string `json:"key_word"`
	RecordStart uint64 `json:"record_start"`
}

func NewWordItem(wtype DictType, rawkw string, rs uint64) (*WordItem, error) {
	realKw := rawkw
	if wtype == MDD {
		// 1. convert hex to bytes
		// 2. convert bytes to utf8 string
		var bytesKw []byte
		var err error
		bytesKw, err = DecodeHex(rawkw)
		if err != nil {
			return nil, err
		}
		realKw, err = DecodeUTF16(bytesKw)
		if err != nil {
			return nil, err
		}
	}

	return &WordItem{
		RawKeyWord:  rawkw,
		KeyWord:     realKw,
		RecordStart: rs,
	}, nil
}

// Distance calculates levenshtein distance.
func (wi *WordItem) Distance(e bktree.Entry) int {
	a := wi.KeyWord
	b := e.(*WordItem).KeyWord
	return levenshtein.Distance(a, b)
}

// HMDistance calculates hamming distance.
func (wi *WordItem) HMDistance(e bktree.Entry) int {
	distance := 0
	s1 := wi.KeyWord
	s2 := e.(*WordItem).KeyWord

	// index by code point, not byte
	r1 := []rune(s1)
	r2 := []rune(s2)

	if len(r1) != len(r2) {
		dis := len(r1) - len(r2)
		if dis < 0 {
			return -dis
		}
		return dis
	}

	for i, v := range r1 {
		if r2[i] != v {
			distance += 1
		}
	}

	return distance
}
