package tmpl

import (
	"testing"
)

func TestReplacerImage_Replace(t *testing.T) {
	image := &ReplacerImage{}
	_, html := image.Replace("1723", nil, TestHTML)
	t.Logf(html)
}
