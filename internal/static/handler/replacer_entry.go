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

package handler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/terasum/medict/pkg/model"
)

var _ Replacer = &ReplacerEntry{}

var ENTRY_REG *regexp.Regexp

func init() {
	var err error
	// Capture the entry word as "anything up to the closing quote". The old
	// whitelist [\w#_ -] skipped entry IDs containing '=', '.', '/', ':', ...
	// (e.g. entry://topic_transport-by-water_level=c1, issue #718): the href
	// stayed a raw entry:// link, the webview passed the unknown scheme to the
	// OS, and the user saw "There is no application set to open the URL
	// entry://...".
	ENTRY_REG, err = regexp.Compile(`href=\"entry://([^\"]+)\"`)
	if err != nil {
		panic(err)
	}
}

// entryWordEscaper makes a captured entry word safe to splice into the
// single-quoted JS string literal emitted below. Backslash and single quote
// are the only characters that can break out of that literal.
var entryWordEscaper = strings.NewReplacer(`\`, `\\`, `'`, `\'`)

type ReplacerEntry struct {
}

func (r *ReplacerEntry) Replace(dictId string, entry *model.MdictKeyWordIndex, html string) (*model.MdictKeyWordIndex, string) {

	if html == "" || dictId == "" {
		return entry, html
	}

	newhtml := html
	matchedGroup := ENTRY_REG.FindAllStringSubmatch(html, -1)
	for _, matched := range matchedGroup {
		if len(matched) != 2 {
			continue
		}
		oldStr := matched[0]  // href="entry://<word>"
		oldWord := matched[1] // <word>
		escapedWord := entryWordEscaper.Replace(oldWord)

		newStr := fmt.Sprintf("href=\"javascript:__medict_entry_jump('%s', '%s');\"", escapedWord, dictId)
		newhtml = strings.ReplaceAll(newhtml, oldStr, newStr)
	}

	return entry, newhtml
}
