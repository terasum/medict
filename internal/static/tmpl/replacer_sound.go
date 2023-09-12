package tmpl

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

func (r *ReplacerSound) Replace(dictId string, entry *model.KeyBlockEntry, html string) (*model.KeyBlockEntry, string) {

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
