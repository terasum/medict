package support

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestWalkDir(t *testing.T) {
	result, err := WalkDir("./testdata/dicts")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	expectPath, _ := filepath.Abs("testdata/dicts/testdict/testdict.mdx")
	assert.Equal(t, result[0].MdxAbsPath, expectPath)
}
