package utils

import (
	"testing"
)

func TestStrToUnicode(t *testing.T) {
	uncodeStr := StrToUnicode("十大户￥@！#%……&……*（）——+《》、，。、；‘、配【】")
	t.Logf("%s", uncodeStr)
	uncodeStr = StrToUnicode("國語詞典")
	t.Logf("%s", uncodeStr)
}
