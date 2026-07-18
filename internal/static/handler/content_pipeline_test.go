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
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestUserCSSOverridesConcurrentAccess(t *testing.T) {
	const dictID = "race-test"
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			SetUserCSS(dictID, fmt.Sprintf("body{%d}", i))
		}(i)
		go func() {
			defer wg.Done()
			_ = GetUserCSS(dictID)
		}()
	}
	wg.Wait()
	SetUserCSS(dictID, "")
}

func TestSafeInlineCSSNeutralisesStyleEndTag(t *testing.T) {
	malicious := `body{color:red}</style><script>alert(1)</script>`
	safe := SafeInlineCSS(malicious)
	if strings.Contains(strings.ToLower(safe), "</style") {
		t.Fatalf("CSS still contains raw style terminator: %q", safe)
	}
	if !strings.Contains(safe, `<\/style><script>`) {
		t.Fatalf("CSS content was not preserved safely: %q", safe)
	}
}
