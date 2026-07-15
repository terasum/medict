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
	"strings"
	"testing"
)

func TestInlineResources_AttributesAndNoMarker(t *testing.T) {
	// note the css uses HTML-escaped &amp;d=0, the img uses raw &d=0
	html := `<head><link rel="stylesheet" href="style.css?dict_id=abc123&amp;d=0"></head>` +
		`<body>` +
		`<img src="photo.png?dict_id=abc123&d=0">` + // referenced twice → exercises cache
		`<img src="photo.png?dict_id=abc123&d=0">` +
		`<img src="missing.png?dict_id=abc123&d=0">` + // unresolvable → left untouched
		`</body>`

	fetch := func(key string) ([]byte, bool) {
		switch key {
		case "style.css":
			return []byte(`.a{color:red}`), true
		case "photo.png":
			return []byte{0x89, 'P', 'N', 'G'}, true
		}
		return nil, false
	}

	out := string(InlineResources([]byte(html), fetch))

	// resolvable resources are fully inlined (no marker left on them)
	if strings.Contains(out, "style.css?dict_id=") {
		t.Errorf("css not inlined:\n%s", out)
	}
	if strings.Contains(out, "photo.png?dict_id=") {
		t.Errorf("image not inlined:\n%s", out)
	}
	if !strings.Contains(out, `src="data:image/png;base64,`) {
		t.Errorf("image not inlined as data url:\n%s", out)
	}
	if !strings.Contains(out, `href="data:text/css;base64,`) {
		t.Errorf("css link not inlined as data url:\n%s", out)
	}
	// unresolvable resource is left in place (marker intact)
	if !strings.Contains(out, "missing.png?dict_id=") {
		t.Errorf("unresolvable resource should be left in place:\n%s", out)
	}
}

func TestInlineCSSURLs_RelativeInlinedAbsoluteSkipped(t *testing.T) {
	css := `.a{background:url(bg.png)} .b{background:url('http://x/y.png')} .c{list-style:url(dot.gif)}`

	out := inlineCSSURLs(css, func(key string) ([]byte, bool) {
		if key == "bg.png" || key == "dot.gif" {
			return []byte{1, 2, 3}, true
		}
		return nil, false
	})

	if !strings.Contains(out, "url(data:image/png;base64,") {
		t.Errorf("relative bg.png should be inlined:\n%s", out)
	}
	if !strings.Contains(out, "http://x/y.png") {
		t.Errorf("absolute http url should be left untouched:\n%s", out)
	}
}

func TestMimeByExt(t *testing.T) {
	cases := map[string]string{
		"style.css":  "text/css",
		"a.JPG":      "image/jpeg",
		"font.woff2": "font/woff2",
		"noise.mp3":  "audio/mpeg",
		"unknown":    "application/octet-stream",
	}
	for key, want := range cases {
		if got := mimeByExt(key); got != want {
			t.Errorf("mimeByExt(%q) = %q, want %q", key, got, want)
		}
	}
}
