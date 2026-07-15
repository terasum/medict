//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package ankiexport

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// ExportRow is one card source, produced by BookmarkStore.ExportRows.
type ExportRow struct {
	Word         string
	DictName     string
	NotebookName string
	HTML         string
}

// ErrNoBookmarks is returned by Build when there is nothing to export.
var ErrNoBookmarks = errors.New("没有可导出的生词")

// Build produces a native Anki .apkg (zip bytes) for the given rows. Each row
// becomes a card (front = word, back = the snapshot's styled definition with
// images extracted as Anki media); each notebook becomes a deck.
func Build(src []ExportRow) ([]byte, error) {
	if len(src) == 0 {
		return nil, ErrNoBookmarks
	}

	media := newMediaCollector()
	prepared := make([]preparedRow, len(src))
	for i, r := range src {
		def := media.rewrite(extractDefinition(r.HTML))
		nb := r.NotebookName
		if nb == "" {
			nb = "Default"
		}
		prepared[i] = preparedRow{
			Word:       r.Word,
			DictName:   r.DictName,
			Notebook:   nb,
			Definition: def,
		}
	}

	now := time.Now()
	colBytes, err := buildCollection(prepared, now.UnixMilli(), now.Unix())
	if err != nil {
		return nil, err
	}
	return zipApkg(colBytes, media.files)
}

// zipApkg assembles the final .apkg: collection.anki2, the `media` JSON map, and
// each media file stored under its numeric key.
func zipApkg(colBytes []byte, files []mediaFile) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	w, err := zw.Create("collection.anki2")
	if err != nil {
		return nil, fmt.Errorf("anki: zip collection: %w", err)
	}
	if _, err := w.Write(colBytes); err != nil {
		return nil, fmt.Errorf("anki: write collection: %w", err)
	}

	mediaMap := make(map[string]string, len(files))
	for i, f := range files {
		mediaMap[strconv.Itoa(i)] = f.name
	}
	mediaJSON, err := json.Marshal(mediaMap)
	if err != nil {
		return nil, fmt.Errorf("anki: marshal media map: %w", err)
	}
	w, err = zw.Create("media")
	if err != nil {
		return nil, fmt.Errorf("anki: zip media map: %w", err)
	}
	if _, err := w.Write(mediaJSON); err != nil {
		return nil, fmt.Errorf("anki: write media map: %w", err)
	}

	for i, f := range files {
		w, err := zw.Create(strconv.Itoa(i))
		if err != nil {
			return nil, fmt.Errorf("anki: zip media file %s: %w", f.name, err)
		}
		if _, err := w.Write(f.data); err != nil {
			return nil, fmt.Errorf("anki: write media file %s: %w", f.name, err)
		}
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("anki: close zip: %w", err)
	}
	return buf.Bytes(), nil
}
