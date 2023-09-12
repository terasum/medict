package tmpl

import (
	"github.com/terasum/medict/pkg/model"
	"regexp"
	"strings"
)

var _ Replacer = &ReplacerJs{}

var JS_REG *regexp.Regexp

func init() {
	var err error
	JS_REG, err = regexp.Compile(`src=\"(\S+\.js)\"`)
	if err != nil {
		panic(err)
	}
}

type ReplacerJs struct {
	mdict *model.DictionaryItem
}

func (r *ReplacerJs) SetDictContext(item *model.DictionaryItem) {
	r.mdict = item
}

func (r *ReplacerJs) Replace(dictId string, entry *model.KeyBlockEntry, html string) (*model.KeyBlockEntry, string) {
	if html == "" || dictId == "" {
		return entry, html
	}

	newhtml := html
	matchedGroup := JS_REG.FindAllStringSubmatch(html, -1)
	for _, matched := range matchedGroup {
		if len(matched) != 2 {
			continue
		}
		oldStr := matched[1]
		newStr := oldStr + "?dict_id=" + dictId + "&d=0"
		newhtml = strings.ReplaceAll(newhtml, oldStr, newStr)
	}

	return entry, newhtml
}
