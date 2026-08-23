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

import "strings"

// irregularLemmas maps common irregular English inflections (mostly past
// tenses/participles that suffix rules cannot undo) to their base form.
var irregularLemmas = map[string]string{
	"ran": "run", "went": "go", "gone": "go",
	"made": "make", "said": "say", "told": "tell",
	"took": "take", "taken": "take", "saw": "see", "seen": "see",
	"came": "come", "got": "get", "gotten": "get",
	"gave": "give", "given": "give", "found": "find",
	"thought": "think", "felt": "feel", "kept": "keep", "held": "hold",
	"brought": "bring", "bought": "buy", "taught": "teach", "caught": "catch",
	"sat": "sit", "stood": "stand", "slept": "sleep", "met": "meet",
	"paid": "pay", "laid": "lay", "lost": "lose", "sent": "send",
	"built": "build", "understood": "understand",
	"spoke": "speak", "spoken": "speak", "broke": "break", "broken": "break",
	"began": "begin", "begun": "begin", "drank": "drink", "drunk": "drink",
	"drove": "drive", "driven": "drive", "ate": "eat", "eaten": "eat",
	"fell": "fall", "fallen": "fall", "flew": "fly", "flown": "fly",
	"knew": "know", "known": "know", "grew": "grow", "grown": "grow",
	"threw": "throw", "thrown": "throw", "wore": "wear", "worn": "wear",
	"wrote": "write", "written": "write", "chose": "choose", "chosen": "choose",
	"forgot": "forget", "forgotten": "forget", "hid": "hide", "hidden": "hide",
	"sang": "sing", "sung": "sing", "swam": "swim", "swum": "swim",
	"rose": "rise", "risen": "rise", "shook": "shake", "shaken": "shake",
}

// lemmaCandidates generates likely base forms of an inflected English word:
// plural / 3rd-person -s/-es (parts→part, boxes→box, studies→study), past
// -ed (loved→love, stopped→stop), progressive -ing (going→go, running→run,
// lying→lie), comparative/superlative -er/-est (faster→fast, biggest→big),
// plus the irregular table above.
//
// ECDICT stores base forms almost exclusively — spot-check on the shipped
// preset DB: parts/boxes/studies/ran/going are all absent — so Search/Lookup
// fall back to these candidates on a miss (roadmap #786). Non-ASCII input
// (e.g. Chinese) yields no candidates: the rules are English-only.
func lemmaCandidates(word string) []string {
	w := strings.ToLower(word)
	if !isAsciiWord(w) {
		return nil
	}
	var out []string
	seen := map[string]bool{}
	add := func(s string) {
		if len(s) >= 2 && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	if base, ok := irregularLemmas[w]; ok {
		add(base)
	}
	if len(w) < 4 {
		// too short for suffix rules ("as"/"its"/"is" must not be mangled);
		// the exact irregular table above already ran
		return out
	}
	hasSuffix := func(s string) bool { return strings.HasSuffix(w, s) }
	stem := func(n int) string { return w[:len(w)-n] }
	doubled := func(base string) bool { // stopped→stop, running→run, biggest→big
		if len(base) >= 2 && base[len(base)-1] == base[len(base)-2] {
			add(base[:len(base)-1])
			return true
		}
		return false
	}
	switch {
	case hasSuffix("ies") && len(w) > 4:
		add(stem(3) + "y") // studies→study, carried→carry
	case hasSuffix("ses"), hasSuffix("xes"), hasSuffix("zes"),
		hasSuffix("ches"), hasSuffix("shes"):
		add(stem(2)) // boxes→box, watches→watch, fixes→fix
	case hasSuffix("ied"):
		add(stem(3) + "y") // carried→carry
	case hasSuffix("es"):
		add(stem(2)) // goes→go, does→do, echoes→echo
	}
	if hasSuffix("s") && !hasSuffix("ss") {
		add(stem(1)) // parts→part, closes→close
	}
	if hasSuffix("ed") {
		base := stem(2)
		if !doubled(base) { // stopped→stop
			add(base) // wanted→want
		}
		add(base + "e") // loved→love, hoped→hope (e-deletion verbs)
	}
	if hasSuffix("ying") && len(w) > 4 {
		add(stem(4) + "ie") // lying→lie, dying→die, tying→tie
	} else if hasSuffix("ing") {
		base := stem(3) // going→go
		add(base)
		if !doubled(base) { // running→run
			add(base + "e") // hoping→hope, basing→base (e-deletion verbs)
		}
	}
	if hasSuffix("iest") {
		add(stem(4) + "y") // happiest→happy
	} else if hasSuffix("est") {
		if !doubled(stem(3)) { // biggest→big
			add(stem(3)) // fastest→fast
		}
	}
	if hasSuffix("ier") {
		add(stem(3) + "y") // happier→happy
	} else if hasSuffix("er") {
		if !doubled(stem(2)) { // bigger→big
			add(stem(2)) // faster→fast
		}
	}
	return out
}

// isAsciiWord reports whether s is non-empty pure ASCII letters.
func isAsciiWord(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
			return false
		}
	}
	return true
}

// lookupCandidates returns the exact-match chain for Lookup: the typed form,
// its lowercase (ECDICT headwords are lowercase; LIKE/exact are
// case-sensitive under case_sensitive_like=ON), then lemma candidates.
func lookupCandidates(keyword string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, 4)
	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	add(keyword)
	add(strings.ToLower(keyword))
	for _, c := range lemmaCandidates(keyword) {
		add(c)
	}
	return out
}
