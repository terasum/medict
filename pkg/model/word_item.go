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

type WrappedWordItem struct {
	DictId string `json:"dict_id"`
	*WordItem
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
