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

func TestReplacerEntry_Replace(t *testing.T) {
	entry := &ReplacerEntry{}
	_, html := entry.Replace("X182310003", nil, TESTENTRYHTML)
	t.Logf("%s", html)
}

// Regression for #718: entry:// IDs may contain '=' and other non-word chars
// ('.', '/', ':', ...). The old whitelist [\w#_ -] did not match them, so the
// raw entry:// href survived into the page and the webview handed it to the OS.
func TestReplacerEntry_HandlesNonWordCharsInEntryURL(t *testing.T) {
	entry := &ReplacerEntry{}
	in := `<a href="entry://topic_transport-by-water_level=c1">x</a>`
	_, out := entry.Replace("dict123", nil, in)

	if strings.Contains(out, "entry://") {
		t.Fatalf("raw entry:// href must be rewritten away, got: %s", out)
	}
	want := `__medict_entry_jump('topic_transport-by-water_level=c1', 'dict123')`
	if !strings.Contains(out, want) {
		t.Fatalf("expected rewritten jump call containing %q, got: %s", want, out)
	}
}

// A captured entry word is spliced into a single-quoted JS literal, so quotes
// and backslashes must be escaped or they break/inject into the href.
func TestReplacerEntry_EscapesWordForJSLiteral(t *testing.T) {
	entry := &ReplacerEntry{}
	_, out := entry.Replace("d1", nil, `<a href="entry://it's">x</a>`)

	if strings.Contains(out, `('it's`) {
		t.Fatalf("single quote in entry word must be escaped, got: %s", out)
	}
	if !strings.Contains(out, `it\'s`) {
		t.Fatalf("expected escaped quote in output, got: %s", out)
	}
}
