//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfig(t *testing.T) {
	cfg, err := ReadConfig("./testdata/test.toml")
	assert.Nil(t, err)
	assert.Equal(t, cfg.BaseDictDir, "testdir")
}

// toStrSlice normalizes a viper slice ([]interface{} or []string) to []string
// so the test doesn't depend on the codec's concrete type.
func toStrSlice(v any) []string {
	switch s := v.(type) {
	case []interface{}:
		out := make([]string, 0, len(s))
		for _, x := range s {
			out = append(out, x.(string))
		}
		return out
	case []string:
		return s
	}
	return nil
}

// TestWritePreferences_RoundTrip writes a base toml, adds preferences via
// WritePreferences, re-reads the same file and asserts the new keys landed AND
// the pre-existing key (BaseDictDir) was preserved.
func TestWritePreferences_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "medict.toml")
	assert.Nil(t, os.WriteFile(path, []byte("BaseDictDir = \"/tmp/dicts\"\n"), 0o644))

	cfg, err := ReadConfig(path)
	assert.Nil(t, err)
	assert.Nil(t, cfg.WritePreferences(map[string]any{
		"MultiDictIds": []string{"a", "b"},
		"FontSize":     18,
	}))

	// re-read the same file from disk
	cfg2, err := ReadConfig(path)
	assert.Nil(t, err)
	prefs := cfg2.Preferences()

	// viper normalizes all keys to lowercase (case-insensitive), so AllSettings
	// returns lowercased keys regardless of the casing used at Set time.
	assert.Equal(t, "/tmp/dicts", prefs["basedictdir"], "BaseDictDir must be preserved")
	assert.Equal(t, []string{"a", "b"}, toStrSlice(prefs["multidictids"]), "MultiDictIds must round-trip")
	assert.Contains(t, prefs, "fontsize")
}
