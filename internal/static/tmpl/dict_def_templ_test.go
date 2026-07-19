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
	"strings"
	"testing"
)

func TestWordDefinitionTemplateUsesClickToLookup(t *testing.T) {
	required := []string{
		"Intl.Segmenter",
		"border-bottom:1px solid currentColor",
		"__medict_isCurrentWord",
		"__Medict_INNER_FRAME_MSG_EVTP_CLICK_LOOKUP",
	}
	for _, marker := range required {
		if !strings.Contains(WordDefinitionTempl, marker) {
			t.Fatalf("word interaction marker %q is missing", marker)
		}
	}

	removed := []string{
		"__Medict_INNER_FRAME_MSG_EVTP_HOVER_LOOKUP",
		"__Medict_INNER_FRAME_MSG_EVTP_HOVER_LEAVE",
		"__Medict_INNER_FRAME_MSG_EVTP_DBLCLICK_LOOKUP",
	}
	for _, marker := range removed {
		if strings.Contains(WordDefinitionTempl, marker) {
			t.Fatalf("legacy lookup marker %q is still present", marker)
		}
	}
}

func TestWordDefinitionTemplateAppliesExplicitContentZoom(t *testing.T) {
	required := []string{
		"__Medict_TOP_WIN_MSG_EVTP_SET_ZOOM",
		"document.documentElement.style.zoom",
		"Math.min(160, Math.max(80",
	}
	for _, marker := range required {
		if !strings.Contains(WordDefinitionTempl, marker) {
			t.Fatalf("content zoom marker %q is missing", marker)
		}
	}
}
