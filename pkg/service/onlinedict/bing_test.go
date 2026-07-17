//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// GPL-3.0.

package onlinedict

import (
	"strings"
	"testing"
)

// fixtureBing is a trimmed, realistic cn.bing.com/dict/clientsearch page for
// "hello" (structure verified against the live endpoint 2026-07).
const fixtureBing = `<!DOCTYPE html><html><head><style>.client_def_hd_hd{font-size:20px}</style></head>
<body><script>var _G={};</script>
<div class="client_def_hd_area"><div class="client_def_hd_pn_bar"><div class="client_def_hd_pn_list">
<span class="client_def_hd_pn">美国:&#160;[he&#712;l&#601;&#650;]</span>
<span class="client_def_hd_pn">英国:&#160;[h&#601;&#712;l&#601;&#650;]</span>
</div></div><div class="client_def_hd_hd">hello</div></div>
<div class="client_def_bar"><div class="client_def_title_bar"><span class="client_def_title">int.</span></div>
<div class="client_def_list"><div class="client_def_list_item"><div class="client_def_list_word_item">
<div class="client_def_list_word_content"><span class="client_def_list_word_bar">你好；喂；您好；哈喽</span>
</div></div></div></div></div>
<div class="client_def_bar"><div class="client_def_title_bar"><span class="client_def_title_web">网络</span></div>
<div class="client_def_list"><div class="client_def_list_item"><div class="client_def_list_word_item">
<div class="client_def_list_word_content"><span class="client_def_list_word_bar">哈罗；哈啰；大家好</span>
</div></div></div></div></div>
</body></html>`

func TestParseBing(t *testing.T) {
	card := parseBing("hello", fixtureBing)
	if card == "" {
		t.Fatal("parseBing returned empty card")
	}
	for _, want := range []string{"hello", "你好", "哈罗", "网络", "int.", "he"} {
		if !strings.Contains(card, want) {
			t.Errorf("card missing %q\n%s", want, card)
		}
	}
	// phonetic nbsp should be normalized (no literal &#160; or U+00A0)
	if strings.Contains(card, "&#160;") || strings.Contains(card, " ") {
		t.Errorf("phonetic not cleaned: %q", card)
	}
}

func TestParseBing_Empty(t *testing.T) {
	// no recognizable structure → empty card
	if got := parseBing("x", "<html><body>nothing here</body></html>"); got != "" {
		t.Errorf("expected empty card, got %q", got)
	}
}

func TestBing_SearchSingle(t *testing.T) {
	b := NewBing()
	res, err := b.Search("hello")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 || res[0].KeyWord != "hello" {
		t.Fatalf("Search = %+v", res)
	}
	if res, _ := b.Search(""); len(res) != 0 {
		t.Fatalf("Search('') should be empty")
	}
}
