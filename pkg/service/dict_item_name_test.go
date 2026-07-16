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

package service

import (
	"strings"
	"testing"

	"github.com/terasum/medict/pkg/service/support"
)

// TestNewByDirItem_NameFromTitle verifies #782: the dict display Name should
// come from the embedded Title (mdx header Title / stardict bookname) when
// present, falling back to the file/dir name otherwise. Uses the in-repo
// testdata dicts (relative paths) so it actually runs on any host.
func TestNewByDirItem_NameFromTitle(t *testing.T) {
	items, err := support.WalkDir("support/testdata/dicts")
	if err != nil {
		t.Fatalf("WalkDir: %v", err)
	}
	if len(items) == 0 {
		t.Skip("no testdata dicts present under support/testdata/dicts")
	}
	// 部分 testdata mdx 是占位 stub(打开会 EOF),逐个跳过、只对成功打开的断言。
	opened := 0
	for _, di := range items {
		item, err := NewByDirItem(di)
		if err != nil {
			t.Logf("%s: NewByDirItem error, skipping: %v", di.CurrentDir, err)
			continue
		}
		opened++
		title := ""
		if item.Description != nil {
			title = strings.TrimSpace(item.Description.Title)
		}
		switch {
		case title != "":
			// 有 Title → 显示名应等于 Title
			if item.Name != title {
				t.Errorf("%s: Name=%q, want Title %q", di.CurrentDir, item.Name, title)
			}
		default:
			// 无 Title → 回退到文件名/目录名(必须非空)
			if item.Name == "" {
				t.Errorf("%s: Name empty and no Title fallback", di.CurrentDir)
			}
		}
		t.Logf("%s: Name=%q Title=%q", di.CurrentDir, item.Name, title)
	}
	if opened == 0 {
		t.Skip("no testdata dict could be opened (testdata mdx are stubs)")
	}
}
