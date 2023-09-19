package stardict

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Dict implements in-memory dictionary
type Dict struct {
	buffer []byte
}

// GetSequence returns data at the given offset
func (d Dict) GetSequence(offset uint64, size uint64) []byte {
	return d.buffer[offset:(offset + size)]
}

// ReadDict reads dictionary into memory
func ReadDict(filename string, info *Info) (dict *Dict, err error) {
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

	dict = new(Dict)
	dict.buffer = buffer

	return
}
