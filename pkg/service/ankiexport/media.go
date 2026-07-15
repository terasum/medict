//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package ankiexport

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"regexp"
	"strings"
)

// mediaFile is one extracted resource (image/audio) written into the .apkg.
type mediaFile struct {
	name string
	data []byte
}

// mediaCollector dedups extracted resources by SHA-1 (so the same image reused
// across many cards is stored once) and assigns each a stable filename.
type mediaCollector struct {
	byHash map[string]string // sha1(hex) → filename
	files  []mediaFile
}

func newMediaCollector() *mediaCollector {
	return &mediaCollector{byHash: map[string]string{}}
}

func (m *mediaCollector) add(mime string, data []byte) string {
	sum := sha1.Sum(data)
	h := hex.EncodeToString(sum[:])
	if name, ok := m.byHash[h]; ok {
		return name
	}
	name := "medict-" + h + extForMime(mime)
	m.byHash[h] = name
	m.files = append(m.files, mediaFile{name: name, data: data})
	return name
}

// srcDataReg matches element resource attributes pointing at a data: URL, e.g.
// <img src="data:image/png;base64,...">. It deliberately matches only src= (not
// href=), so inlined CSS on <link href="data:..."> stays put (the Anki reviewer
// renders it via its webview). Group 1 = `src="` prefix (with its spacing),
// group 2 = the data: URL body, group 3 = the closing quote.
var srcDataReg = regexp.MustCompile(`(src\s*=\s*")(data:[^"]*)(")`)

// rewrite replaces every src="data:..." in html with src="<filename>", recording
// the extracted bytes on the collector. Unparseable data: URLs are left as-is.
func (m *mediaCollector) rewrite(html string) string {
	return srcDataReg.ReplaceAllStringFunc(html, func(match string) string {
		sub := srcDataReg.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		mime, data, ok := parseDataURL(sub[2])
		if !ok {
			return match
		}
		return sub[1] + m.add(mime, data) + sub[3]
	})
}

// parseDataURL decodes a `data:[<mime>][;base64],<payload>` URL. Returns ok=false
// for anything that isn't a decodable data: URL.
func parseDataURL(url string) (mime string, data []byte, ok bool) {
	const prefix = "data:"
	if !strings.HasPrefix(url, prefix) {
		return "", nil, false
	}
	rest := url[len(prefix):]
	comma := strings.IndexByte(rest, ',')
	if comma < 0 {
		return "", nil, false
	}
	header, payload := rest[:comma], rest[comma+1:]
	mime = header
	if i := strings.IndexByte(header, ';'); i >= 0 {
		mime = header[:i]
	}
	if mime == "" || mime == "data" {
		mime = "application/octet-stream"
	}
	if strings.Contains(header, ";base64") {
		dec, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return "", nil, false
		}
		return mime, dec, true
	}
	// Rare: URL-encoded payload (no ;base64).
	return mime, []byte(payload), true
}

// extForMime maps a MIME type to a file extension (mirror of handler.mimeByExt).
func extForMime(mime string) string {
	switch strings.ToLower(mime) {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "image/svg+xml":
		return ".svg"
	case "image/webp":
		return ".webp"
	case "image/bmp":
		return ".bmp"
	case "image/x-icon":
		return ".ico"
	case "audio/mpeg", "audio/mp3":
		return ".mp3"
	case "audio/wav":
		return ".wav"
	case "audio/ogg":
		return ".ogg"
	case "audio/mp4", "audio/m4a":
		return ".m4a"
	case "text/css":
		return ".css"
	case "text/javascript", "application/javascript":
		return ".js"
	case "font/woff":
		return ".woff"
	case "font/woff2":
		return ".woff2"
	case "font/ttf":
		return ".ttf"
	case "font/otf":
		return ".otf"
	case "application/vnd.ms-fontobject":
		return ".eot"
	default:
		return ".bin"
	}
}

// styleBlockReg matches a full <style>...</style> block (dotall).
var styleBlockReg = regexp.MustCompile(`(?si)<style[^>]*>.*?</style>`)

// bodyReg captures the inner HTML of <body> (dotall, non-greedy).
var bodyReg = regexp.MustCompile(`(?si)<body[^>]*>(.*?)</body>`)

// extractDefinition turns a full snapshot HTML document into a card field value:
// the <body> inner HTML prefixed with any <style> blocks (so dictionary styling
// survives). If there's no <body>, the whole input is used. The caller rewrites
// src="data:..." images to media filenames afterwards.
func extractDefinition(htmlStr string) string {
	styles := styleBlockReg.FindAllString(htmlStr, -1)
	inner := ""
	if m := bodyReg.FindStringSubmatch(htmlStr); len(m) >= 2 {
		inner = m[1]
	} else {
		inner = htmlStr
	}
	// Drop <style> from the body inner (already hoisted) to avoid duplication.
	inner = styleBlockReg.ReplaceAllString(inner, "")
	var b strings.Builder
	for _, s := range styles {
		b.WriteString(s)
	}
	b.WriteString(inner)
	return b.String()
}
