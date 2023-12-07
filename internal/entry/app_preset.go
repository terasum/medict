package entry

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/terasum/medict/internal/utils"
)

//go:embed preset/cc-cedict/cc-cedict.mdx
var PresetMdictMdxFile []byte

//go:embed preset/cc-cedict/cc-cedict.mdd
var PresetMdictMddFile []byte

//go:embed preset/cc-cedict/cc-cedict.css
var PresetMdictCSSFile []byte

//go:embed preset/cc-cedict/cover.png
var PresetMdictCoverImg []byte

func WritePresetDictionary(baseDictDir string) error {
	fullpath := filepath.Join(baseDictDir, "cc-cedict")
	if !utils.FileExists(fullpath) {
		if err := os.MkdirAll(fullpath, 0755); err != nil {
			return err
		}
	}

	mdxfilePath := filepath.Join(fullpath, "cc-cedict.mdx")
	if !utils.FileExists(mdxfilePath) {
		err := os.WriteFile(mdxfilePath, PresetMdictMdxFile, 0644)
		if err != nil {
			return err
		}
	}

	mddfilePath := filepath.Join(fullpath, "cc-cedict.mdd")
	if !utils.FileExists(mddfilePath) {
		err := os.WriteFile(mddfilePath, PresetMdictMddFile, 0644)
		if err != nil {
			return err
		}
	}

	cssfilePath := filepath.Join(fullpath, "cc-cedict.css")
	if !utils.FileExists(cssfilePath) {
		err := os.WriteFile(cssfilePath, PresetMdictCSSFile, 0644)
		if err != nil {
			return err
		}
	}

	coverfilePath := filepath.Join(fullpath, "cover.png")
	if !utils.FileExists(coverfilePath) {
		err := os.WriteFile(coverfilePath, PresetMdictCoverImg, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
