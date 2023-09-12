package entry

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/utils"
)

func defaultConfigPath() string {
	home, err := utils.HomeDir()
	if err != nil {
		panic(err)
	}
	configDir := filepath.Join(home, ".medict", "dicts")
	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	configFile := filepath.Join(home, ".medict", "medict.toml")
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(configFile, []byte(configTmpl), 0644)
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
