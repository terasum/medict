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
	"regexp"
	"strings"
)

var _ Handler = &HandlerFont{}

var FONT_REG *regexp.Regexp

func init() {
	var err error
	FONT_REG, err = regexp.Compile(`url\([\"|\'](\S+\.(ttf|otf|woff|woff2))[\"|\']\)`)
	if err != nil {
		panic(err)
	}
}

type HandlerFont struct {
}

func (r *HandlerFont) Match(dictId string, key string) bool {
	log.Infof("handle font matching %s:%s", dictId, key)
	if key == "" || dictId == "" {
		return false
	}
	for _, suff := range []string{"css", "scss"} {
		if strings.HasSuffix(key, suff) {
			return true
		}
	}
	return false
}

func (r *HandlerFont) ReplaceRaw(dictId string, templ string, content string) string {
	return content
}

func (r *HandlerFont) Replace(dictId string, keyWord string, resource []byte) (string, []byte) {

	if resource == nil {
		return keyWord, resource
	}

	newHtml := string(resource)
	matchedGroup := FONT_REG.FindAllStringSubmatch(newHtml, -1)
	for _, matched := range matchedGroup {
		if len(matched) != 3 {
			continue
		}
		oldStr := matched[1]
		newStr := oldStr + "?dict_id=" + dictId + "&d=0"
		newHtml = strings.ReplaceAll(newHtml, oldStr, newStr)
	}

	return keyWord, []byte(newHtml)
}
