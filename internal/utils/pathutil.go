package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func HomeDir() (string, error) {
	return os.UserHomeDir()
}

func ReplaceHome(origin string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return origin, err
	}
	origin = strings.ReplaceAll(origin, "$HOME", home)
	return origin, nil
}

func FetchBaseDirName(fpath string) string {
	return filepath.Base(fpath)
}
