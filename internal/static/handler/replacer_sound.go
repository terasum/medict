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
	"github.com/terasum/medict/pkg/model"
	"regexp"
	"strings"
)

var _ Replacer = &ReplacerSound{}

var SOUND_REG *regexp.Regexp

func init() {
	var err error
	// <a class="fayin" href="sound://us/pond__us_1.mp3"><span class="phon-us">pɑːnd</span><img src="us_pron.png?dict_id=f234356c227f82a54afdaa3514de188a&amp;d=0" class="fayin"></a>
	SOUND_REG, err = regexp.Compile(`href=[\"|'](sound://\S+\.mp3)[\"|']`)
	if err != nil {
		panic(err)
	}
}

type ReplacerSound struct {
}

func (r *ReplacerSound) Replace(dictId string, entry *model.MdictKeyWordIndex, html string) (*model.MdictKeyWordIndex, string) {

	if html == "" || dictId == "" {
		return entry, html
	}

	newhtml := html
	matchedGroup := SOUND_REG.FindAllStringSubmatch(html, -1)
	for _, matched := range matchedGroup {
		if len(matched) != 2 {
			continue
		}
		oldStr := matched[1]
		// sound://us/pond__us_1.mp3?dict_id=1236&d=0
		oldStr = strings.TrimPrefix(oldStr, "sound://")
		soundURL := fmt.Sprintf(`%s?dict_id=%s&d=0`, oldStr, dictId)
		newStr := fmt.Sprintf("javascript:__medict_play_sound('%s')", soundURL)
		newhtml = strings.ReplaceAll(newhtml, "sound://"+oldStr, newStr)
	}

	return entry, newhtml
}
