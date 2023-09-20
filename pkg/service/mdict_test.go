package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMdict_Name(t *testing.T) {
	mdict := &Mdict{
		mdxFilePath:       "/User/yourname/test/dict/test.mdx",
		mddFilePaths:      nil,
		mdxins:            nil,
		mddinss:           nil,
		hasBuildIndex:     false,
		buildingIndexLock: nil,
	}
	t.Logf("name is %s", mdict.Name())
	assert.Equal(t, "test", mdict.Name())
}
