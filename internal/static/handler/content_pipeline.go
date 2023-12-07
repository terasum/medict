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
	"encoding/base64"
	"fmt"
	"github.com/sym01/htmlsanitizer"
	"github.com/terasum/medict/internal/static/tmpl"

	"github.com/terasum/medict/pkg/model"
)

type ContentPreHandlePipeline struct {
	replacers []Replacer
	handlers  []Handler
}

var handler = &ContentPreHandlePipeline{
	replacers: []Replacer{
		&ReplacerCss{},
		&ReplacerJs{},
		&ReplacerImage{},
		&ReplacerSound{},
		&ReplacerEntry{},
	},
	handlers: []Handler{
		&HandlerFont{},
	},
}

var sanitizer = htmlsanitizer.NewHTMLSanitizer()
var base64encoder = base64.StdEncoding

func init() {
	sanitizer.RemoveTag("a")
	sanitizer.RemoveTag("img")
}

func WrapDesc(dictid, title, desc string) string {
	rep1 := &ReplacerImage{}
	rep2 := &ReplacerCss{}
	rep3 := &ReplacerEntry{}
	rawHtml := fmt.Sprintf(tmpl.DictDescTempl, title, desc)

	_, rawHtml = rep1.Replace(dictid, nil, rawHtml)
	_, rawHtml = rep2.Replace(dictid, nil, rawHtml)
	_, rawHtml = rep3.Replace(dictid, nil, rawHtml)

	rawHtml = base64encoder.EncodeToString([]byte(rawHtml))
	return rawHtml
	//sanitizedHTML, err := s.SanitizeString(rawHtml)
	//if err != nil {
	//	return "[sanitized failed]"
	//}
	//return sanitizedHTML
}

func WrapContent(dict *model.PlainDictionaryItem, keyEntry *model.MdictKeyWordIndex, definition string) ([]byte, error) {
	content := handleContent(dict, keyEntry, definition)
	return []byte(fmt.Sprintf(tmpl.WordDefinitionTempl, dict.Name, dict.ID, dict.Name, dict.ID, content)), nil
}

func WrapResource(dictId string, keyWord string, resource []byte) ([]byte, error) {
	content := handleResource(dictId, keyWord, resource)
	return content, nil
}

func handleContent(dict *model.PlainDictionaryItem, keyEntry *model.MdictKeyWordIndex, definition string) string {
	for _, replacer := range handler.replacers {
		keyEntry, definition = replacer.Replace(dict.ID, keyEntry, definition)
	}
	return definition
}

func handleResource(dictId string, keyWord string, resource []byte) []byte {
	for _, han := range handler.handlers {
		if han.Match(dictId, keyWord) {
			fmt.Printf("resource handle matched [%s](%s)\n", keyWord, dictId)
			keyWord, resource = han.Replace(dictId, keyWord, resource)
		}
	}
	return resource
}
