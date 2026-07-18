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
	"path/filepath"
	"regexp"
	"strings"
)

// InlineResources rewrites resource references in html into self-contained
// data: URLs so the snapshot renders without the embedded Gin server (i.e.
// after the source dictionary is unloaded). It handles:
//   - attribute refs rewritten by the mdict replacers — img/link/script
//     `src|href="KEY?dict_id=<token>&d=0"` (the marker the replacers append;
//     `&` may be HTML-escaped as `&amp;`),
//   - `url(...)` references inside <style> blocks and inside linked CSS,
//     (relative keys only — http(s):/data:/#/entry:// are skipped).
//
// Linked CSS is fetched, its own url() refs inlined, then embedded as a data
// URL on the <link>. Sound refs (a `javascript:` handler) and anything fetch
// can't resolve are left untouched. fetch resolves a resource key (as written
// by the replacers, before the `?dict_id` marker) to its raw bytes; ok=false
// leaves the reference in place. Repeated keys are cached within one call.
func InlineResources(html []byte, fetch func(key string) ([]byte, bool)) []byte {
	cache := map[string]struct {
		data []byte
		ok   bool
	}{}
	fetchCached := func(key string) ([]byte, bool) {
		if r, hit := cache[key]; hit {
			return r.data, r.ok
		}
		d, ok := fetch(key)
		cache[key] = struct {
			data []byte
			ok   bool
		}{d, ok}
		return d, ok
	}

	h := string(html)
	// 1) url() refs already inline in <style> blocks
	h = styleBlockReg.ReplaceAllStringFunc(h, func(m string) string {
		sub := styleBlockReg.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		return strings.Replace(m, sub[1], inlineCSSURLs(sub[1], fetchCached), 1)
	})
	// 2) attribute resource refs (img/link/script)
	h = attrRefReg.ReplaceAllStringFunc(h, func(m string) string {
		sub := attrRefReg.FindStringSubmatch(m)
		if len(sub) < 3 {
			return m
		}
		attr, key := sub[1], sub[2]
		data, ok := fetchCached(key)
		if !ok {
			return m // leave as-is (resource missing)
		}
		if mimeByExt(key) == "text/css" {
			data = []byte(inlineCSSURLs(string(data), fetchCached))
		}
		return attr + `="data:` + mimeByExt(key) + ";base64," + base64.StdEncoding.EncodeToString(data) + `"`
	})
	return []byte(h)
}

// attrRefReg: (src|href)="KEY?dict_id=<token>&d=0" (optionally &amp;).
var attrRefReg = regexp.MustCompile(`(src|href)\s*=\s*"([^"]+?)\?dict_id=[^"&]+&(?:amp;)?d=0"`)

// styleBlockReg: the inner CSS of a <style>...</style> (dotall).
var styleBlockReg = regexp.MustCompile(`(?s)<style[^>]*>(.*?)</style>`)

// urlRefReg: url( KEY ) with optional quotes; KEY = anything but )'".
var urlRefReg = regexp.MustCompile(`url\(\s*['"]?([^'")]+)['"]?\s*\)`)

func inlineCSSURLs(css string, fetch func(string) ([]byte, bool)) string {
	return urlRefReg.ReplaceAllStringFunc(css, func(m string) string {
		sub := urlRefReg.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		key := strings.TrimSpace(sub[1])
		if skipURLKey(key) {
			return m
		}
		data, ok := fetch(key)
		if !ok {
			return m
		}
		return `url(data:` + mimeByExt(key) + ";base64," + base64.StdEncoding.EncodeToString(data) + `)`
	})
}

// InlineCSSResources makes a standalone stylesheet self-contained by replacing
// relative url(...) references with data URLs. It is used when dictionary CSS
// becomes the editable user-override seed.
func InlineCSSResources(css string, fetch func(string) ([]byte, bool)) string {
	return inlineCSSURLs(css, fetch)
}

// skipURLKey leaves absolute/special references alone (can't or shouldn't inline).
func skipURLKey(key string) bool {
	low := strings.ToLower(key)
	switch {
	case strings.HasPrefix(low, "http://"), strings.HasPrefix(low, "https://"):
		return true
	case strings.HasPrefix(low, "//"):
		return true
	case strings.HasPrefix(low, "data:"):
		return true
	case strings.HasPrefix(low, "#"):
		return true
	case strings.HasPrefix(low, "entry://"), strings.HasPrefix(low, "sound://"):
		return true
	}
	return false
}

func mimeByExt(key string) string {
	switch strings.ToLower(filepath.Ext(key)) {
	case ".css":
		return "text/css"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".js", ".mjs":
		return "text/javascript"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	case ".otf":
		return "font/otf"
	case ".eot":
		return "application/vnd.ms-fontobject"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".m4a":
		return "audio/mp4"
	default:
		return "application/octet-stream"
	}
}
