package entry

import (
	"errors"
	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

func defaultConfigPath() string {
	home, err := utils.HomeDir()
	if err != nil {
		panic(err)
	}
	configDir := filepath.Join(home, ".medictapp", "dicts")
	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	configFile := filepath.Join(home, ".medictapp", "medict.toml")
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		err = ioutil.WriteFile(configFile, []byte(configTmpl), 0644)
		if err != nil {
			panic(err)
		}
	}
	return configFile

}

var cfg *config.Config

func DefaultConfig() (*config.Config, error) {
	var err error
	if cfg == nil {
		cfg, err = config.ReadConfig(defaultConfigPath())
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}
