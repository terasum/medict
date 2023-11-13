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

package tmpl

import (
	"github.com/terasum/medict/pkg/model"
	"regexp"
	"strings"
)

var _ Replacer = &ReplacerImage{}

var imageReg *regexp.Regexp

func init() {
	var err error
	imageReg, err = regexp.Compile(`src=["|'](\S+\.(png|jpg|gif|jpeg|svg))["|']`)
	if err != nil {
		panic(err)
	}
}

type ReplacerImage struct {
}

func (r *ReplacerImage) Replace(dictId string, entry *model.MdictKeyWordIndex, html string) (*model.MdictKeyWordIndex, string) {
	if html == "" || dictId == "" {
		return entry, html
	}
	cache := make(map[string]bool)

	newHtml := html
	matchedGroup := imageReg.FindAllStringSubmatch(html, -1)
	for _, matched := range matchedGroup {
		if len(matched) != 3 {
			continue
		}
		imgstr := matched[1]
		if _, ok := cache[imgstr]; ok {
			continue
		} else {
			cache[imgstr] = true
			newStr := imgstr + "?dict_id=" + dictId + "&d=0"
			newHtml = strings.ReplaceAll(newHtml, imgstr, newStr)
		}
	}

	return entry, newHtml
}
