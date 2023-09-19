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

func defaultConfigPath() (string, error) {
	home, err := utils.HomeDir()
	if err != nil {
		return "", err
	}
	fmt.Printf("medict app: default config dir %s\n", home)
	configDir := filepath.Join(home, "dicts")
	if _, err = os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	configFile := filepath.Join(home, "medict.toml")
	if _, err = os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(configFile, []byte(configTmpl), 0644)
		if err != nil {
			return "", err
		}
	}
	fmt.Printf("medict app: default config file %s\n", configFile)
	return configFile, nil

}

var cfg *config.Config

func DefaultConfig() (*config.Config, error) {
	var err error
	var configPath = ""
	if cfg != nil {
		return cfg, nil
	}
	configPath, err = defaultConfigPath()
	cfg, err = config.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}
	fmt.Printf("medict app: config dicts dir %s\n", cfg.BaseDictDir)
	return cfg, nil
}
