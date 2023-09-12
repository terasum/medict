package tmpl

import "testing"

func TestReplacerCss_Replace(t *testing.T) {
	css := &ReplacerCss{}
	css.Replace("1236", nil, TestHTML)
}
