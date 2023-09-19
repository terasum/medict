//go:build windows
// +build windows

package utils

import (
	"os"
	"strings"

	"github.com/kirsle/configdir"
)

func ReplaceHome(origin string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return origin, err
	}

	origin = strings.ReplaceAll(origin, "$HOME", home)
	return origin, nil
}

func HomeDir() (string, error) {
	configPath := configdir.LocalConfig("medict")
	err := configdir.MakePath(configPath) // Ensure it exists.
	if err != nil {
		return "", err
	}
	return configPath, nil
}
