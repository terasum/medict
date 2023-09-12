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
