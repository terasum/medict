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

package entry

import (
	"github.com/terasum/medict/pkg/service"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config, err := LoadApp()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(config.ConfigPath)
	t.Logf(config.BaseDictDir)
	t.Logf(config.ConfigStruct.BaseDictDir)
}

func TestEnsureConfigDir(t *testing.T) {
	config, err := LoadApp()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(config.ConfigPath)
	t.Logf(config.BaseDictDir)
	t.Logf(config.ConfigStruct.BaseDictDir)
	t.Logf(config.EnsureDictsDir())
}

func TestLoadApp(t *testing.T) {
	conf, err := LoadApp()
	if err != nil {
		t.Fatal(err)
	}

	dir := conf.EnsureDictsDir()
	t.Logf(dir)

	svc, err := service.NewDictService(conf)
	if err != nil {
		t.Fatal(err)
	}

	dicts := svc.Dicts()

	for _, dict := range dicts {
		t.Logf("dict name: %s", dict.Name)

	}

}
