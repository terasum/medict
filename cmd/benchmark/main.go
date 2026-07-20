// cmd/benchmark — 内置词典性能 benchmark 工具。
//
// 为每个词典(ECDICT 离线英汉 + cc-cedict 汉英 mdict)的每个场景采集
// CPU profile + heap profile,解析为结构化 JSON(供 AI 分析)。
//
// 用法:
//
//	go run ./cmd/benchmark                       # 默认:跑所有词典/场景
//	go run ./cmd/benchmark -iterations 10000     # 自定义迭代次数
//	go run ./cmd/benchmark -top 30               # pprof top-N 行数
//	go run ./cmd/benchmark -out /tmp/bench       # 自定义输出目录
//	go run ./cmd/benchmark -dicts ecdict         # 只跑 ECDICT
//	go run ./cmd/benchmark -dicts ccedict        # 只跑 cc-cedict
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service/ecdict"
	"github.com/terasum/medict/pkg/service/mdict"
)

// ─── data structures (JSON output for AI) ───

type FuncStat struct {
	Function string  `json:"function"`
	FlatPct  float64 `json:"flat_pct"`
	CumPct   float64 `json:"cum_pct"`
}

type HeapStat struct {
	Function string  `json:"function"`
	AllocPct float64 `json:"alloc_pct"`
}

type ScenarioResult struct {
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Iterations     int        `json:"iterations"`
	NsPerOp        int64      `json:"ns_per_op"`
	BytesPerOp     int64      `json:"bytes_per_op"`
	AllocsPerOp    int64      `json:"allocs_per_op"`
	GCCount        uint32     `json:"gc_count"`
	GCPauseTotalNS uint64     `json:"gc_pause_total_ns"`
	Goroutines     int        `json:"goroutines"`
	CPUTop         []FuncStat `json:"cpu_top_functions"`
	HeapTop        []HeapStat `json:"heap_top_allocators"`
}

type DictResult struct {
	DictName    string           `json:"dict_name"`
	DictType    string           `json:"dict_type"`
	EntryCount  int              `json:"entry_count"`
	BuildIndexS float64          `json:"build_index_seconds"`
	Scenarios   []ScenarioResult `json:"scenarios"`
}

type Report struct {
	Benchmark   string       `json:"benchmark"`
	Timestamp   string       `json:"timestamp"`
	Iterations  int          `json:"default_iterations"`
	Dictionaries []DictResult `json:"dictionaries"`
}

// ─── generic scenario ───

type Scenario struct {
	Name        string
	Description string
	Warmup      int
	// Iterations override (0 = use global default)
	IterationsOverride int
	Run                func(d model.GeneralDictionary) error
}

// ─── main ───

func main() {
	iterations := flag.Int("iterations", 3000, "default iterations per scenario")
	topN := flag.Int("top", 15, "pprof top-N functions to capture")
	outDir := flag.String("out", "bench_out", "output directory")
	dictsFlag := flag.String("dicts", "all", "which dicts to benchmark: all|ecdict|ccedict")
	flag.Parse()

	repoRoot := findRepoRoot()
	profDir := filepath.Join(*outDir, "prof")
	os.MkdirAll(profDir, 0755)

	var dictResults []DictResult

	// ── ECDICT ──
	if *dictsFlag == "all" || *dictsFlag == "ecdict" {
		ecdictPath := filepath.Join(repoRoot, "internal", "entry", "preset", "ecdict", "ecdict.db")
		r := benchmarkECDICT(ecdictPath, *iterations, *topN, profDir)
		if r != nil {
			dictResults = append(dictResults, *r)
		}
	}

	// ── cc-cedict (mdict) ──
	if *dictsFlag == "all" || *dictsFlag == "ccedict" {
		ccedictDir := filepath.Join(repoRoot, "internal", "entry", "preset", "cc-cedict")
		r := benchmarkCcedict(ccedictDir, *iterations, *topN, profDir)
		if r != nil {
			dictResults = append(dictResults, *r)
		}
	}

	report := Report{
		Benchmark:   "medict-builtin-dicts",
		Timestamp:   time.Now().Format(time.RFC3339),
		Iterations:  *iterations,
		Dictionaries: dictResults,
	}

	// write JSON
	jsonPath := filepath.Join(*outDir, "bench_report.json")
	jsonData, _ := json.MarshalIndent(report, "", "  ")
	os.WriteFile(jsonPath, jsonData, 0644)
	fmt.Printf("\nJSON: %s\n", jsonPath)

	// write markdown
	mdPath := filepath.Join(*outDir, "bench_report.md")
	os.WriteFile(mdPath, []byte(generateMarkdown(&report)), 0644)
	fmt.Printf("Markdown: %s\n", mdPath)
	fmt.Printf("Profiles: %s/*.prof\n", profDir)
}

// ─── ECDICT benchmark ───

func benchmarkECDICT(dbPath string, defaultIters, topN int, profDir string) *DictResult {
	fmt.Printf("\n=== ECDICT (offline SQLite) ===\n")
	dict, err := ecdict.NewECDict(&model.DirItem{CurrentDir: filepath.Dir(dbPath)})
	if err != nil {
		fmt.Fprintf(os.Stderr, "skip ECDICT: %v\n", err)
		return nil
	}
	defer dict.Close()
	dict.BuildIndex()

	status, _ := dict.Status()
	edition := status.Edition
	if edition == "" {
		edition = "compact"
	}
	fmt.Printf("edition=%s entries=%d\n", edition, status.EntryCount)

	scenarios := []Scenario{
		{"search_1char", "LIKE 'a%' — widest prefix", 50, 0, func(d model.GeneralDictionary) error { _, e := d.Search("a"); return e }},
		{"search_3char", "LIKE 'app%'", 50, 0, func(d model.GeneralDictionary) error { _, e := d.Search("app"); return e }},
		{"search_full", "LIKE 'apple%'", 50, 0, func(d model.GeneralDictionary) error { _, e := d.Search("apple"); return e }},
		{"search_miss", "LIKE 'zzzzzz%'", 50, 0, func(d model.GeneralDictionary) error { _, e := d.Search("zzzzzz"); return e }},
		{"lookup_common", "exact 'apple'", 100, 0, func(d model.GeneralDictionary) error { _, e := d.Lookup("apple"); return e }},
		{"lookup_miss", "exact 'xxnotaword'", 100, 0, func(d model.GeneralDictionary) error { _, e := d.Lookup("xxnotaword"); return e }},
		{"locate_pipeline", "Search('app')→Locate(first)", 50, 0, func(d model.GeneralDictionary) error {
			res, e := d.Search("app")
			if e != nil || len(res) == 0 {
				return e
			}
			_, e = d.Locate(res[0])
			return e
		}},
	}

	results := runScenarios(dict, scenarios, "ecdict", defaultIters, topN, profDir)
	return &DictResult{DictName: "ECDICT", DictType: "SQLite", EntryCount: status.EntryCount, Scenarios: results}
}

// ─── cc-cedict (mdict) benchmark ───

func benchmarkCcedict(dir string, defaultIters, topN int, profDir string) *DictResult {
	fmt.Printf("\n=== cc-cedict (mdict) ===\n")
	dirItem := &model.DirItem{
		CurrentDir:      dir,
		MdictMdxAbsPath: filepath.Join(dir, "cc-cedict.mdx"),
		MdictMddAbsPath: []string{filepath.Join(dir, "cc-cedict.mdd")},
		DictType:        model.DictTypeMdict,
	}
	dict, err := mdict.NewMdictSvc(dirItem)
	if err != nil {
		fmt.Fprintf(os.Stderr, "skip cc-cedict: %v\n", err)
		return nil
	}
	defer dict.Close()

	// BuildIndex — time it (first build, then cached)
	fmt.Printf("BuildIndex ... ")
	biStart := time.Now()
	dict.BuildIndex()
	biElapsed := time.Since(biStart)
	fmt.Printf("%.2fs\n", biElapsed.Seconds())

	// find test words: search common prefixes to get real entries
	testWords := findTestWords(dict, []string{"hello", "test", "你", "apple"})
	fmt.Printf("test words: %v\n", testWords)

	if len(testWords) == 0 {
		testWords = []string{"a"} // fallback
	}
	word := testWords[0]
	prefix := word
	if len([]rune(prefix)) > 3 {
		prefix = string([]rune(prefix)[:3])
	}

	scenarios := []Scenario{
		{"search_prefix", fmt.Sprintf("Search('%s')", prefix), 30, 0, func(d model.GeneralDictionary) error {
			_, e := d.Search(prefix)
			return e
		}},
		{"search_miss", "Search('zzzzzzz')", 30, 0, func(d model.GeneralDictionary) error {
			_, e := d.Search("zzzzzzz")
			return e
		}},
		{"lookup_word", fmt.Sprintf("Lookup('%s')", word), 30, 0, func(d model.GeneralDictionary) error {
			_, e := d.Lookup(word)
			return e
		}},
		{"locate_pipeline", fmt.Sprintf("Search('%s')→Locate(first)", prefix), 20, defaultIters / 3, func(d model.GeneralDictionary) error {
			res, e := d.Search(prefix)
			if e != nil || len(res) == 0 {
				return e
			}
			_, e = d.Locate(res[0])
			return e
		}},
	}

	results := runScenarios(dict, scenarios, "ccedict", defaultIters, topN, profDir)

	// count entries (approximate: search for very broad prefix)
	entryCount := 0
	if res, err := dict.Search("a"); err == nil {
		entryCount = len(res) // just the first 50 results, not total
	}

	return &DictResult{
		DictName:    "cc-cedict",
		DictType:    "mdict",
		EntryCount:  entryCount,
		BuildIndexS: biElapsed.Seconds(),
		Scenarios:   results,
	}
}

// ─── scenario runner ───

func runScenarios(dict model.GeneralDictionary, scenarios []Scenario, prefix string, defaultIters, topN int, profDir string) []ScenarioResult {
	var results []ScenarioResult
	for _, sc := range scenarios {
		iters := defaultIters
		if sc.IterationsOverride > 0 {
			iters = sc.IterationsOverride
		}
		if strings.Contains(sc.Name, "1char") {
			iters /= 5
		}
		if iters < 50 {
			iters = 50
		}

		scenarioName := prefix + "_" + sc.Name
		fmt.Printf("  %-30s ... ", scenarioName)
		res := runScenario(dict, sc, scenarioName, iters, profDir, topN)
		results = append(results, res)
		fmt.Printf("%10d ns/op  %8d B/op  %4d allocs\n", res.NsPerOp, res.BytesPerOp, res.AllocsPerOp)
	}
	return results
}

func runScenario(dict model.GeneralDictionary, sc Scenario, scenarioName string, iterations int, profDir string, topN int) ScenarioResult {
	// warmup
	for i := 0; i < sc.Warmup; i++ {
		sc.Run(dict)
	}
	runtime.GC()

	// CPU profile
	cpuFile := filepath.Join(profDir, scenarioName+"_cpu.prof")
	cf, err := os.Create(cpuFile)
	if err == nil {
		pprof.StartCPUProfile(cf)
	}

	var msBefore, msAfter runtime.MemStats
	runtime.ReadMemStats(&msBefore)
	gcBefore := msBefore.NumGC

	start := time.Now()
	for i := 0; i < iterations; i++ {
		sc.Run(dict)
	}
	elapsed := time.Since(start)

	runtime.ReadMemStats(&msAfter)
	if cf != nil {
		pprof.StopCPUProfile()
		cf.Close()
	}

	bytesPerOp := int64(0)
	allocsPerOp := int64(0)
	if iterations > 0 {
		bytesPerOp = int64(msAfter.TotalAlloc-msBefore.TotalAlloc) / int64(iterations)
		allocsPerOp = int64(msAfter.Mallocs-msBefore.Mallocs) / int64(iterations)
	}

	// heap profile
	memFile := filepath.Join(profDir, scenarioName+"_mem.prof")
	mf, _ := os.Create(memFile)
	pprof.WriteHeapProfile(mf)
	mf.Close()

	cpuTop := parsePprofText(cpuFile, topN)
	heapTop := parseHeapText(memFile, topN)

	return ScenarioResult{
		Name:           scenarioName,
		Description:    sc.Description,
		Iterations:     iterations,
		NsPerOp:        elapsed.Nanoseconds() / int64(iterations),
		BytesPerOp:     bytesPerOp,
		AllocsPerOp:    allocsPerOp,
		GCCount:        msAfter.NumGC - gcBefore,
		GCPauseTotalNS: msAfter.PauseTotalNs - msBefore.PauseTotalNs,
		Goroutines:     runtime.NumGoroutine(),
		CPUTop:         cpuTop,
		HeapTop:        heapTop,
	}
}

// findTestWords searches common words and returns ones that have results.
func findTestWords(dict model.GeneralDictionary, candidates []string) []string {
	var found []string
	for _, w := range candidates {
		res, err := dict.Search(w)
		if err == nil && len(res) > 0 {
			found = append(found, res[0].KeyWord)
		}
	}
	return found
}

// ─── pprof parsing ───

func parsePprofText(profFile string, topN int) []FuncStat {
	out, err := exec.Command("go", "tool", "pprof", "-text",
		fmt.Sprintf("-nodecount=%d", topN), profFile).CombinedOutput()
	if err != nil {
		return nil
	}
	return parsePprofLines(string(out))
}

func parseHeapText(profFile string, topN int) []HeapStat {
	out, err := exec.Command("go", "tool", "pprof", "-text", "-alloc_space",
		fmt.Sprintf("-nodecount=%d", topN), profFile).CombinedOutput()
	if err != nil {
		return nil
	}
	return parseHeapLines(string(out))
}

func parsePprofLines(text string) []FuncStat {
	var stats []FuncStat
	inData := false
	for _, raw := range strings.Split(text, "\n") {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "flat") && strings.Contains(line, "flat%") {
			inData = true
			continue
		}
		if !inData || line == "" || !strings.Contains(line, "%") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		flatPct := parseFloat(strings.TrimSuffix(fields[1], "%"))
		cumPct := parseFloat(strings.TrimSuffix(fields[4], "%"))
		fn := strings.Join(fields[5:], " ")
		stats = append(stats, FuncStat{Function: fn, FlatPct: flatPct, CumPct: cumPct})
	}
	return stats
}

func parseHeapLines(text string) []HeapStat {
	var stats []HeapStat
	inData := false
	for _, raw := range strings.Split(text, "\n") {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "flat") && strings.Contains(line, "flat%") {
			inData = true
			continue
		}
		if !inData || line == "" || !strings.Contains(line, "%") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		allocPct := parseFloat(strings.TrimSuffix(fields[1], "%"))
		fn := strings.Join(fields[5:], " ")
		stats = append(stats, HeapStat{Function: fn, AllocPct: allocPct})
	}
	return stats
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

// ─── markdown report ───

func generateMarkdown(r *Report) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# Benchmark Report: %s\n\n", r.Benchmark))
	b.WriteString(fmt.Sprintf("- **Timestamp**: %s\n", r.Timestamp))
	b.WriteString(fmt.Sprintf("- **Iterations**: %d\n\n", r.Iterations))

	for _, dr := range r.Dictionaries {
		b.WriteString(fmt.Sprintf("## %s (%s)\n\n", dr.DictName, dr.DictType))
		if dr.BuildIndexS > 0 {
			b.WriteString(fmt.Sprintf("BuildIndex: **%.2fs**\n\n", dr.BuildIndexS))
		}
		b.WriteString("| Scenario | ns/op | µs/op | B/op | allocs | GC |\n")
		b.WriteString("|---|---:|---:|---:|---:|---:|\n")
		for _, s := range dr.Scenarios {
			b.WriteString(fmt.Sprintf("| %s | %d | %.1f | %d | %d | %d |\n",
				s.Name, s.NsPerOp, float64(s.NsPerOp)/1000, s.BytesPerOp, s.AllocsPerOp, s.GCCount))
		}
		b.WriteString("\n")
		for _, s := range dr.Scenarios {
			if len(s.CPUTop) > 0 {
				b.WriteString(fmt.Sprintf("### %s — CPU Top\n\n", s.Name))
				b.WriteString("| Function | flat%% | cum%% |\n|---|---:|---:|\n")
				for _, f := range s.CPUTop {
					b.WriteString(fmt.Sprintf("| `%s` | %.1f | %.1f |\n", f.Function, f.FlatPct, f.CumPct))
				}
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}

// ─── helpers ───

func findRepoRoot() string {
	// walk up from cwd to find go.mod
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}
