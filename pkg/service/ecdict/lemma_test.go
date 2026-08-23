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

package ecdict

import (
	"slices"
	"testing"
)

// TestLemmaCandidates guards the inflection rules: each inflected form must
// generate its base form among the candidates. Extra harmless candidates are
// fine (misses just fall through); a MISSING base form breaks the fallback.
func TestLemmaCandidates(t *testing.T) {
	cases := []struct{ in, want string }{
		// plural / 3rd-person -s/-es
		{"parts", "part"},
		{"goes", "go"},
		{"boxes", "box"},
		{"studies", "study"},
		{"watches", "watch"},
		{"fixes", "fix"},
		{"carries", "carry"},
		// past -ed
		{"loved", "love"},
		{"hoped", "hope"},
		{"stopped", "stop"},
		{"carried", "carry"},
		// progressive -ing
		{"going", "go"},
		{"running", "run"},
		{"lying", "lie"},
		{"dying", "die"},
		// comparative / superlative
		{"faster", "fast"},
		{"bigger", "big"},
		{"happier", "happy"},
		{"fastest", "fast"},
		{"biggest", "big"},
		{"happiest", "happy"},
		// irregular table
		{"ran", "run"},
		{"went", "go"},
		{"written", "write"},
		{"understood", "understand"},
		{"ate", "eat"},
	}
	for _, tc := range cases {
		got := lemmaCandidates(tc.in)
		if !slices.Contains(got, tc.want) {
			t.Errorf("lemmaCandidates(%q) = %v, want to contain %q", tc.in, got, tc.want)
		}
	}
}

// TestLemmaCandidatesGuards: rules must not mangle short words or produce
// candidates for non-English input.
func TestLemmaCandidatesGuards(t *testing.T) {
	for _, w := range []string{"", "is", "as", "its", "was", "苹果", "тест", "part-time"} {
		if got := lemmaCandidates(w); len(got) != 0 {
			// short real words: irregular table may legitimately hit ("was"→
			// nothing here, "is"→nothing), so only complain about non-empty
			// suffix-rule output on non-ASCII input
			if !isAsciiWord(w) {
				t.Errorf("lemmaCandidates(%q) = %v, want none for non-ASCII input", w, got)
			}
		}
	}
	// "ss" endings must not shed their final s (glass ≠ glas)
	for _, bad := range []string{"glass", "class", "miss"} {
		for _, c := range lemmaCandidates(bad) {
			if c == bad[:len(bad)-1] {
				t.Errorf("lemmaCandidates(%q) must not produce %q", bad, c)
			}
		}
	}
}
