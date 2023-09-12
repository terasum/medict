package tmpl

import (
	"testing"
)

func TestReplacerFont_Replace(t *testing.T) {
	fonts := &HandlerFont{}
	_, html := fonts.Replace("1723", nil, TestHTML)
	t.Log(html)
}
