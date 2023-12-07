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
	ENTRY_REG, err = regexp.Compile(`href=\"entry://([\w#_ -]+)\"`)
	if err != nil {
		panic(err)
	}
}

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
		oldStr := matched[0]
		oldWord := strings.TrimRight(matched[0], "\"")
		oldWord = strings.TrimPrefix(oldWord, "href=\"entry://")

		newStr := fmt.Sprintf("href=\"javascript:__medict_entry_jump('%s', '%s');\"", oldWord, dictId)
		fmt.Printf("old %s => new %s\n", oldStr, newStr)

		newhtml = strings.ReplaceAll(newhtml, oldStr, newStr)
	}

	return entry, newhtml
}
