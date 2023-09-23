package stardict

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// RawDict implements in-memory dictionary
type RawDict struct {
	buffer []byte
}

// GetSequence returns data at the given offset
func (d *RawDict) GetSequence(offset uint64, size uint64) []byte {
	return d.buffer[offset:(offset + size)]
}

// readDict reads dictionary into memory
func readDict(filename string, info *RawInfo) (dict *RawDict, err error) {
	reader, err := os.Open(filename)
	if err != nil {
		return
	}
	defer reader.Close()

	var r io.Reader

	if strings.HasSuffix(filename, ".dz") { // if file is compressed then read it from archive
		r, err = gzip.NewReader(reader)
	} else {
		r = reader
	}

	if err != nil {
		return
	}

	buffer, err := ioutil.ReadAll(r)

	if err != nil {
		return
	}

	dict = new(RawDict)
	dict.buffer = buffer

	return
}
