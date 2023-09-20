//go:build darwin

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kirsle/configdir"
)

func ReplaceHome(origin string) (string, error) {
	configPath, err := HomeDir()
	if err != nil {
		return "", err
	}

	origin = strings.ReplaceAll(origin, "$HOME", configPath)
	fmt.Printf("medict app: ReplaceHome dictionary dir %s\n", origin)
	return origin, nil
}

func HomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("userhomedir err: %s", err.Error())
	}
	configPath := filepath.Join(home, ".medict")

	err = configdir.MakePath(configPath)
	if err != nil {
		return "", fmt.Errorf("base home dir mkall failed, %s", err.Error())
	}

	return configPath, nil
}

func AppConfigDir() (string, error) {
	configPath := configdir.LocalConfig("medict")
	fmt.Printf("medict app: AppConfigDir %s\n", configPath)
	err := configdir.MakePath(configPath) // Ensure it exists.
	if err != nil {
		return "", fmt.Errorf("app config dir mkall failed, %s", err.Error())
	}
	return configPath, nil
}
