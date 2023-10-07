package utils

import (
	"testing"
)

func TestStrToUnicode(t *testing.T) {
	uncodeStr := StrToUnicode("十大户￥@！#%……&……*（）——+《》、，。、；‘、配【】")
	t.Logf(uncodeStr)
	uncodeStr = StrToUnicode("國語詞典")
	t.Logf(uncodeStr)
}
