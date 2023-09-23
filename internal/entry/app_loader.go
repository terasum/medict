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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/utils"
)

var cfg *config.Config

func defaultConfigPath() (string, error) {
	// should raise error
	appConfigDir, err := utils.AppConfigDir()
	if err != nil {
		return "", err
	}
	fmt.Printf("[medict-init]: default config dir: %s\n", appConfigDir)

	configFile := filepath.Join(appConfigDir, "medict.toml")
	if _, err = os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(configFile, []byte(config.ConfigTmpl), 0644)
		if err != nil {
			return "", err
		}
	}

	fmt.Printf("[medict-init]: default config file: %s\n", configFile)
	return configFile, nil
}

func loadConfig() (*config.Config, error) {
	var err error
	var configPath = ""
	if cfg != nil {
		return cfg, nil
	}

	configPath, err = defaultConfigPath()
	if err != nil {
		return nil, err
	}

	cfg, err = config.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}
	fmt.Printf("[medict-init]: config dicts dir: %s\n", cfg.BaseDictDir)
	return cfg, nil
}

func LoadApp() (*config.Config, error) {
	conf, err := loadConfig()
	if err != nil {
		return nil, err
	}
	conf.BaseDictDir = conf.EnsureDictsDir()

	err = WritePresetDictionary(conf.BaseDictDir)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
