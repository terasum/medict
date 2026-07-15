//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package ankiexport

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// validatorScript imports Anki's own AnkiPackageImporter (the exact code path
// Anki uses when a user imports an .apkg), imports the package into a fresh
// collection, and prints a JSON summary on stdout. It mirrors /tmp/validate_apkg.py.
const validatorScript = `import json, os, sys, tempfile, warnings
import anki.lang
from anki.collection import Collection
from anki.importing.apkg import AnkiPackageImporter
warnings.filterwarnings("ignore")  # silence Anki's deprecation notice
anki.lang.set_lang("en")
apkg_path, out_path = sys.argv[1], sys.argv[2]
work = tempfile.mkdtemp(prefix="anki-compat-")
dst = Collection(os.path.join(work, "v.anki2"))
imp = AnkiPackageImporter(dst, apkg_path)
imp.run()
notes = dst.find_notes("")
sample = None
if notes:
    n = dst.get_note(notes[0])
    sample = {fname: n[fname] for fname in n.keys()}
res = {
    "notes": len(notes),
    "cards": len(dst.find_cards("")),
    "decks": [d.name for d in dst.decks.all_names_and_ids()],
    "definition": (sample or {}).get("Definition", ""),
}
dst.close()
with open(out_path, "w") as f:
    json.dump(res, f)
`

// sampleExportRows builds two cards across two notebooks, one carrying an
// inlined image, exercising the full builder.
func sampleExportRows() []ExportRow {
	return []ExportRow{
		{
			Word: "apple", DictName: "OALD", NotebookName: "Fruit",
			HTML: `<html><head><style>p{color:#000}</style></head><body><p>a round fruit</p><img src="data:image/png;base64,` + png1x1B64 + `"></body></html>`,
		},
		{Word: "dog", DictName: "OALD", NotebookName: "Animals",
			HTML: `<html><body><p>a common pet</p></body></html>`},
	}
}

// TestAnkiPackageCompat is the gold-standard compatibility check: it feeds a
// generated .apkg to Anki's own importer (AnkiPackageImporter) and asserts the
// notes/cards/decks/media land. It is gated behind the ANKI_VALIDATE_PYTHON env
// var (a path to a Python interpreter with the `anki` package installed) so the
// default test run needs no Python; set the env var to opt in.
//
//	pip install anki genanki
//	ANKI_VALIDATE_PYTHON=/path/to/venv/bin/python go test ./pkg/service/ankiexport/ -run TestAnkiPackageCompat -count=1
func TestAnkiPackageCompat(t *testing.T) {
	py := os.Getenv("ANKI_VALIDATE_PYTHON")
	if py == "" {
		t.Skip("set ANKI_VALIDATE_PYTHON=<path to python with anki installed> to validate against Anki's importer")
	}

	apkg, err := Build(sampleExportRows())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	dir := t.TempDir()
	apkgPath := filepath.Join(dir, "sample.apkg")
	if err := os.WriteFile(apkgPath, apkg, 0o644); err != nil {
		t.Fatal(err)
	}
	scriptPath := filepath.Join(dir, "validate_apkg.py")
	if err := os.WriteFile(scriptPath, []byte(validatorScript), 0o644); err != nil {
		t.Fatal(err)
	}
	resultPath := filepath.Join(dir, "result.json")

	cmd := exec.Command(py, scriptPath, apkgPath, resultPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Anki importer rejected the .apkg: %v\n%s", err, out)
	}

	resBytes, err := os.ReadFile(resultPath)
	if err != nil {
		t.Fatalf("read importer result (%s): %v", resultPath, err)
	}
	var res struct {
		Notes      int      `json:"notes"`
		Cards      int      `json:"cards"`
		Decks      []string `json:"decks"`
		Definition string   `json:"definition"`
	}
	if err := json.Unmarshal(resBytes, &res); err != nil {
		t.Fatalf("parse importer result: %v\nraw: %s", err, resBytes)
	}
	if res.Notes != 2 || res.Cards != 2 {
		t.Fatalf("imported notes=%d cards=%d, want 2/2", res.Notes, res.Cards)
	}
	decks := strings.Join(res.Decks, ",")
	if !strings.Contains(decks, "Medict::Fruit") || !strings.Contains(decks, "Medict::Animals") {
		t.Fatalf("notebook decks not created: %v", res.Decks)
	}
	// Image must have been extracted to an Anki media file and referenced by name.
	if !strings.Contains(res.Definition, "medict-") || !strings.Contains(res.Definition, ".png") {
		t.Fatalf("definition should reference an extracted media file: %q", res.Definition)
	}
	t.Logf("Anki importer OK: %d notes, %d cards, decks=%v", res.Notes, res.Cards, res.Decks)
}
