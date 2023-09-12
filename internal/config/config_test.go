package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfig(t *testing.T) {
	cfg, err := ReadConfig("./testdata/test.toml")
	assert.Nil(t, err)
	assert.Equal(t, cfg.StaticServerPort, 1234)
	assert.Equal(t, cfg.BaseDictDir, "testdir")
}
