package tmpl

import (
	"fmt"
	"regexp"
	"strings"
)

var _ Handler = &HandlerFont{}

var FONT_REG *regexp.Regexp

func init() {
	var err error
	FONT_REG, err = regexp.Compile(`url\(\"(\S+\.(ttf|otf|woff|woff2))\"\)`)
	if err != nil {
		panic(err)
	}
}

type HandlerFont struct {
}

func (r *HandlerFont) Match(dictId string, key string) bool {
	fmt.Printf("handle font matching %s:%s", dictId, key)
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
