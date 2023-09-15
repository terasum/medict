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
	"fmt"
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
	},
	handlers: []Handler{
		&HandlerFont{},
	},
}

func WrapContent(dictId string, keyEntry *model.KeyBlockEntry, definition string) ([]byte, error) {
	content := handleContent(dictId, keyEntry, definition)
	return []byte(fmt.Sprintf(WordDefinitionTempl, content)), nil
}

func WrapResource(dictId string, keyWord string, resource []byte) ([]byte, error) {
	content := handleResource(dictId, keyWord, resource)
	return content, nil
}

func handleContent(dictId string, keyEntry *model.KeyBlockEntry, definition string) string {
	for _, replacer := range handler.replacers {
		keyEntry, definition = replacer.Replace(dictId, keyEntry, definition)
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
