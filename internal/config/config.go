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

import "github.com/spf13/viper"

type Config struct {
	BaseDictDir      string `toml:"baseDictDir"`
	StaticServerPort int
	CacheFileDir     string
}

func ReadConfig(configFilePath string) (*Config, error) {
	vip := viper.New()
	vip.SetConfigFile(configFilePath)
	vip.SetConfigType("toml")
	err := vip.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = vip.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
