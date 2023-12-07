package stardict

import (
	"encoding/binary"
	"io/ioutil"
)

// Sense has information belonging to single item position in dictionary
type Sense struct {
	Offset uint64
	Size   uint64
}

// RawIdx implements an in-memory index for a dictionary
type RawIdx struct {
	items map[string][]*Sense
}

// NewIdx initializes idx struct
func NewIdx() *RawIdx {
	idx := new(RawIdx)
	idx.items = make(map[string][]*Sense)
	return idx
}

// Add adds an item to in-memory index
func (idx *RawIdx) Add(item string, offset uint64, size uint64) {
	idx.items[item] = append(idx.items[item], &Sense{Offset: offset, Size: size})
}

// Get gets all translations for an item
func (idx RawIdx) Get(item string) []*Sense {
	return idx.items[item]
}

// readIndex reads dictionary index into a memory and returns in-memory index structure
func readIndex(filename string, info *RawInfo) (idx *RawIdx, err error) {
	data, err := ioutil.ReadFile(filename)

	// unable to read index
	if err != nil {
		return
	}

	idx = NewIdx()

	var a [255]byte // temporary buffer
	var aIdx int
	var expect int

	var dataStr string
	var dataOffset uint64
	var dataSize uint64

	var maxIntBytes int

	if info.Is64 == true {
		maxIntBytes = 8
	} else {
		maxIntBytes = 4
	}

	for _, b := range data {
		if expect == 0 {
			a[aIdx] = b
			if b == 0 {
				dataStr = string(a[:aIdx])

				aIdx = 0
				expect++
				continue
			}
			aIdx++
		} else {
			if expect == 1 {
				a[aIdx] = b
				if aIdx == maxIntBytes-1 {
					if info.Is64 {
						dataOffset = binary.BigEndian.Uint64(a[:maxIntBytes])
					} else {
						dataOffset = uint64(binary.BigEndian.Uint32(a[:maxIntBytes]))
					}

					aIdx = 0
					expect++
					continue
				}
				aIdx++
			} else {
				a[aIdx] = b
				if aIdx == maxIntBytes-1 {
					if info.Is64 {
						dataSize = binary.BigEndian.Uint64(a[:maxIntBytes])
					} else {
						dataSize = uint64(binary.BigEndian.Uint32(a[:maxIntBytes]))
					}

					aIdx = 0
					expect = 0

					// finished with one record
					idx.Add(dataStr, dataOffset, dataSize)

					continue
				}
				aIdx++
			}
		}
	}

	return idx, err
}
