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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/terasum/medict/internal/utils"
)

type ConfigStruct struct {
	BaseDictDir string `toml:"BaseDictDir"`
}

type Config struct {
	*ConfigStruct
	ConfigPath    string
	ViperInstance *viper.Viper
}

func ReadConfig(configFilePath string) (*Config, error) {
	vip := viper.New()
	vip.SetConfigFile(configFilePath)
	vip.SetConfigType("toml")
	err := vip.ReadInConfig()
	if err != nil {
		return nil, err
	}

	configs := &ConfigStruct{}
	err = vip.Unmarshal(configs)
	if err != nil {
		return nil, err
	}

	return &Config{
		ConfigStruct:  configs,
		ViperInstance: vip,
		ConfigPath:    configFilePath,
	}, nil
}

func (c *Config) Write() error {
	return c.ViperInstance.WriteConfig()
}

func (c *Config) EnsureDictsDir() string {
	dirPath := c.BaseDictDir
	defer func() {
		fmt.Printf("[medict-init]: ensure app config directory: %s\n", dirPath)
	}()

	if dirPath == "" {
		dirPath, _ = utils.AppConfigDir()
		dirPath = filepath.Join(dirPath, "medict", "dicts")
	}

	home, _ := utils.HomeDir()
	dirPath = strings.ReplaceAll(dirPath, "$HOME", home)
	appdir, _ := utils.AppConfigDir()
	dirPath = strings.ReplaceAll(dirPath, "$APPCONFDIR", appdir)
	if utils.FileExists(dirPath) {
		return dirPath
	}
	// 不存在就创建
	err := os.MkdirAll(dirPath, 0755)
	if err == nil {
		return dirPath
	}
	// 依旧创建失败, 选择默认配置
	if !strings.HasPrefix(dirPath, home) || strings.HasPrefix(dirPath, appdir) {
		dirPath, _ = utils.AppConfigDir()
		dirPath = filepath.Join(dirPath, "medict", "dicts")
	}

	if !utils.FileExists(dirPath) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	return dirPath
}
