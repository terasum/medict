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
	"regexp"
	"sync"

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

var base64encoder = base64.StdEncoding

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
}

// userCSSOverrides holds per-dictionary CSS overrides in memory (populated from
// the app config dir at startup, updated by SaveDictUserCSS IPC). WrapContent
// reads from here — no file I/O per render.
var userCSSOverrides = map[string]string{}
var userCSSOverridesMu sync.RWMutex
var inlineStyleEndPattern = regexp.MustCompile(`(?i)</style`)

// SetUserCSS stores a per-dictionary CSS override in the in-memory map.
func SetUserCSS(dictId, css string) {
	userCSSOverridesMu.Lock()
	defer userCSSOverridesMu.Unlock()
	if css == "" {
		delete(userCSSOverrides, dictId)
	} else {
		userCSSOverrides[dictId] = css
	}
}

// GetUserCSS returns the in-memory override for a dictionary ("" if none).
func GetUserCSS(dictId string) string {
	userCSSOverridesMu.RLock()
	defer userCSSOverridesMu.RUnlock()
	return userCSSOverrides[dictId]
}

// SafeInlineCSS prevents CSS raw-text from terminating the surrounding style
// element. Dictionary CSS is third-party input; `</style` is inert in a .css
// response but would otherwise become executable HTML when used as an inline
// user override.
func SafeInlineCSS(css string) string {
	return inlineStyleEndPattern.ReplaceAllString(css, `<\/style`)
}

func WrapContent(dict *model.PlainDictionaryItem, keyEntry *model.MdictKeyWordIndex, definition string) ([]byte, error) {
	content := handleContent(dict, keyEntry, definition)
	// #783: inject per-dictionary user CSS override from the in-memory map.
	if css := GetUserCSS(dict.ID); css != "" {
		content += `<style id="medict-user-css">` + SafeInlineCSS(css) + `</style>`
	}
	return []byte(fmt.Sprintf(tmpl.WordDefinitionTempl, dict.Name, dict.ID, dict.Name, dict.ID, dict.ID, content)), nil
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
			log.Debugf("resource handle matched [%s](%s)", keyWord, dictId)
			keyWord, resource = han.Replace(dictId, keyWord, resource)
		}
	}
	return resource
}
