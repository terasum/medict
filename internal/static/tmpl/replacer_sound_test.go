package tmpl

import "testing"

func TestReplacerSound_Replace(t *testing.T) {
	sound := &ReplacerSound{}
	_, html := sound.Replace("1236", nil, TestHTML)
	t.Logf(html)

}
