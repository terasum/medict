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

func (r *ReplacerImage) Replace(dictId string, entry *model.KeyBlockEntry, html string) (*model.KeyBlockEntry, string) {
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
