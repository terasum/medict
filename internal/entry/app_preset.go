package entry

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/terasum/medict/internal/utils"
)

//go:embed preset/dictd-elements/dictd-elements.dict.dz
var PresetStarDictCEDictDZFile []byte

//go:embed preset/dictd-elements/dictd-elements.ifo
var PresetStarDictCEDictIFOFile []byte

//go:embed preset/dictd-elements/dictd-elements.idx
var PresetStarDictCEDictIDXFile []byte

//go:embed preset/dictd-elements/cover.png
var PresetStarDictCEDictCoverImg []byte

func WritePresetDictionary(baseDictDir string) error {
	fullpath := filepath.Join(baseDictDir, "cedict-gb")
	if !utils.FileExists(fullpath) {
		if err := os.MkdirAll(fullpath, 0755); err != nil {
			return err
		}
	}

	dzfilePath := filepath.Join(fullpath, "cedict-gb.dict.dz")
	if !utils.FileExists(dzfilePath) {
		err := os.WriteFile(dzfilePath, PresetStarDictCEDictDZFile, 0644)
		if err != nil {
			return err
		}
	}

	ifofilePath := filepath.Join(fullpath, "cedict-gb.ifo")
	if !utils.FileExists(ifofilePath) {
		err := os.WriteFile(ifofilePath, PresetStarDictCEDictIFOFile, 0644)
		if err != nil {
			return err
		}
	}

	idxfilePath := filepath.Join(fullpath, "cedict-gb.idx")
	if !utils.FileExists(idxfilePath) {
		err := os.WriteFile(idxfilePath, PresetStarDictCEDictIDXFile, 0644)
		if err != nil {
			return err
		}
	}

	coverfilePath := filepath.Join(fullpath, "cover.png")
	if !utils.FileExists(coverfilePath) {
		err := os.WriteFile(coverfilePath, PresetStarDictCEDictCoverImg, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
