package tmpl

import (
	"github.com/terasum/medict/pkg/model"
	"regexp"
	"strings"
)

var _ Replacer = &ReplacerCss{}

var CSS_REG *regexp.Regexp

func init() {
	var err error
	CSS_REG, err = regexp.Compile(`href=\"(\S+\.css)\"`)
	if err != nil {
		panic(err)
	}
}

type ReplacerCss struct {
	mdict *model.DictionaryItem
}

func (r ReplacerCss) SetDictContext(item *model.DictionaryItem) {
	r.mdict = item
}

func (r *ReplacerCss) Replace(dictId string, entry *model.KeyBlockEntry, html string) (*model.KeyBlockEntry, string) {

	if html == "" || dictId == "" {
		return entry, html
	}

	newhtml := html
	matchedGroup := CSS_REG.FindAllStringSubmatch(html, -1)
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
