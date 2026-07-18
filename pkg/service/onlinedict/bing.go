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

// Package onlinedict provides online dictionary backends (default free sources)
// implementing model.GeneralDictionary, so they plug into the existing dict
// list / search / content pipeline like the offline dicts.
//
// Bing (default, no API key): scrapes cn.bing.com/dict/clientsearch (the Bing
// Dictionary desktop-client endpoint, which returns the definition HTML server-
// side — more stable than the consumer /search page whose content is JS-loaded).
// Reference parsing: oldstone/bing-dictionary-for-goldendict (Python, GoldenDict
// use case). Free, no key; downside: HTML structure can change.
package onlinedict

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/op/go-logging"
	"github.com/terasum/medict/pkg/model"
)

var onlineLog = logging.MustGetLogger("service.onlinedict")

// bingSearchURL is the desktop-client endpoint (ClientVer = Bing Dict TV client).
const bingSearchURL = "https://cn.bing.com/dict/clientsearch?mkt=zh-CN&setLang=zh&form=BDVEHC&ClientVer=BDDTV3.5.1.4320&q="

const bingUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36"

// Bing is an online EN-CN dictionary backed by cn.bing.com/dict (no API key).
type Bing struct {
	client *http.Client
}

// NewBing builds a Bing dict with a short-timeout HTTP client. The client uses
// http.DefaultTransport, so it honors the http(s)_proxy env (clash) when set.
func NewBing() *Bing {
	return &Bing{client: &http.Client{Timeout: 8 * time.Second}}
}

func (b *Bing) DictType() model.DictType { return model.DictTypeOnline }

func (b *Bing) Name() string { return "Bing 在线词典" }

func (b *Bing) Description() *model.PlainDictionaryInfo {
	return &model.PlainDictionaryInfo{
		Title:       "Bing 在线词典",
		Description: "必应在线英汉(cn.bing.com/dict)",
	}
}

// BuildIndex is a no-op for an online dict.
func (b *Bing) BuildIndex() error { return nil }

// Close is a no-op (no persistent resources).
func (b *Bing) Close() error { return nil }

// LookupResource — online dicts carry no packable resources.
func (b *Bing) LookupResource(keyword string) ([]byte, error) {
	return nil, fmt.Errorf("onlinedict: no resources (%q)", keyword)
}

// Search returns the exact word as a single pseudo-match — online dicts don't
// do prefix suggestion; this supports the "type word → Enter → locate" flow.
func (b *Bing) Search(keyword string) ([]*model.KeyQueryIndex, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []*model.KeyQueryIndex{}, nil
	}
	return []*model.KeyQueryIndex{{
		IndexType:         string(model.DictTypeOnline),
		MdictKeyWordIndex: &model.MdictKeyWordIndex{KeyWord: keyword},
	}}, nil
}

// Lookup fetches the word from Bing and returns an HTML fragment ("" if no data).
func (b *Bing) Lookup(keyword string) ([]byte, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []byte(""), nil
	}
	page, err := b.fetch(keyword)
	if err != nil {
		// Definitions are loaded in an iframe. Propagating a transient provider
		// failure makes the embedded Gin handler return HTTP 500 and leaves the
		// whole dictionary section unusable. Keep the failure observable while
		// rendering a safe, provider-local fallback instead.
		onlineLog.Warningf("Bing lookup failed for %q: %v", keyword, err)
		return []byte(renderBingUnavailable(keyword)), nil
	}
	card := parseBing(keyword, page)
	return []byte(card), nil
}

func renderBingUnavailable(word string) string {
	return fmt.Sprintf(`<style>
.bi-unavailable{font-family:-apple-system,"Segoe UI","Microsoft YaHei",sans-serif;color:#6b7280;padding:16px 12px;line-height:1.6;}
.bi-unavailable-word{color:#374151;font-weight:600;}
</style><div class="bi-unavailable" role="status"><span class="bi-unavailable-word">%s</span>：在线词典暂时不可用，请稍后重试。</div>`, html.EscapeString(word))
}

// Locate serves the word carried by the entry.
func (b *Bing) Locate(entry *model.KeyQueryIndex) ([]byte, error) {
	kw := ""
	if entry != nil && entry.MdictKeyWordIndex != nil {
		kw = entry.KeyWord
	}
	return b.Lookup(kw)
}

func (b *Bing) fetch(word string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, bingSearchURL+url.QueryEscape(word), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", bingUserAgent)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	resp, err := b.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bing: HTTP %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// bingDef is one sense: a part-of-speech (or "网络" web label) + translation.
type bingDef struct {
	web       bool
	pos, trns string
}

var (
	rxBingHead = regexp.MustCompile(`client_def_hd_hd[^>]*>([^<]+)`)
	rxBingPhon = regexp.MustCompile(`client_def_hd_pn[^>]*>([^<]+)`)
	// client_def_title(_web)?">POS</span> … client_def_list_word_bar">CN</span>
	rxBingDef = regexp.MustCompile(`client_def_title(_web)?">([^<]+)</span>[\s\S]*?client_def_list_word_bar">([^<]+)</span>`)
)

// parseBing extracts the headword/phonetic + senses from a clientsearch page and
// returns an HTML card fragment ("" if nothing recognizable was found).
func parseBing(word, page string) string {
	phon := ""
	if m := rxBingPhon.FindStringSubmatch(page); len(m) > 1 {
		phon = cleanText(m[1])
	}
	var defs []bingDef
	for _, m := range rxBingDef.FindAllStringSubmatch(page, -1) {
		defs = append(defs, bingDef{web: m[1] == "_web", pos: cleanText(m[2]), trns: cleanText(m[3])})
	}
	if phon == "" && len(defs) == 0 {
		return ""
	}
	_ = rxBingHead // (reserved for an explicit headword check / future use)
	return renderBingCard(word, phon, defs)
}

// cleanText unescapes HTML entities (e.g. &#160; → U+00A0) then normalizes the
// resulting nbsp to a normal space. Unescape must run BEFORE the nbsp replace.
func cleanText(s string) string {
	s = html.UnescapeString(s)
	s = strings.ReplaceAll(s, "\u00a0", " ")
	return strings.TrimSpace(s)
}

func renderBingCard(word, phon string, defs []bingDef) string {
	var b strings.Builder
	b.WriteString(`<style>
.bi-card{font-family:-apple-system,"Segoe UI","Microsoft YaHei",sans-serif;color:#222;max-width:680px;padding:4px 2px;}
.bi-head{display:flex;align-items:baseline;gap:10px;margin-bottom:6px;flex-wrap:wrap;}
.bi-word{font-size:26px;font-weight:700;}
.bi-phon{color:#6b7785;font-size:14px;}
.bi-row{display:flex;align-items:baseline;gap:8px;padding:3px 0;font-size:15px;line-height:1.6;}
.bi-pos{flex:0 0 auto;min-width:40px;color:#fff;background:#5a6;font-size:12px;padding:1px 6px;border-radius:3px;text-align:center;}
.bi-web{flex:0 0 auto;min-width:40px;color:#fff;background:#888;font-size:12px;padding:1px 6px;border-radius:3px;text-align:center;}
.bi-trans{flex:1 1 auto;}
</style>`)
	b.WriteString(`<div class="bi-card"><div class="bi-head">`)
	fmt.Fprintf(&b, `<span class="bi-word">%s</span>`, html.EscapeString(word))
	if phon != "" {
		fmt.Fprintf(&b, `<span class="bi-phon">%s</span>`, html.EscapeString(phon))
	}
	b.WriteString(`</div>`)
	for _, d := range defs {
		cls := "bi-pos"
		if d.web {
			cls = "bi-web"
		}
		fmt.Fprintf(&b, `<div class="bi-row"><span class="%s">%s</span><span class="bi-trans">%s</span></div>`,
			cls, html.EscapeString(d.pos), html.EscapeString(d.trns))
	}
	b.WriteString(`</div>`)
	return b.String()
}
