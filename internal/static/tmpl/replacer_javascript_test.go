package tmpl

import (
	"testing"
)

func TestReplacerJs_Replace(t *testing.T) {
	js := &ReplacerJs{}
	_, html := js.Replace("1283", nil, TestHTML)
	t.Logf(html)
}
