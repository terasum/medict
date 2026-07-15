package apis

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/terasum/medict/internal/static"
	"github.com/terasum/medict/internal/static/handler"
	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service"
)

var log = logging.MustGetLogger("apis")

type DictsController struct {
	svc *service.DictService
}

func NewDictsController(svc *service.DictService) *DictsController {
	return &DictsController{svc: svc}
}

// linkTargets returns the @@@LINK target words in an MDict definition, handling
// multiple targets and "</>"/newline sub-record separators. Empty if def is not
// a redirect. (Extracted for unit testing — #260.)
func linkTargets(def string) []string {
	normalized := strings.ReplaceAll(strings.TrimSpace(def), "</>", "\n")
	var targets []string
	for _, line := range strings.Split(normalized, "\n") {
		line = strings.TrimSpace(line)
		if t := strings.TrimPrefix(line, "@@@LINK="); t != line && t != "" {
			targets = append(targets, strings.TrimSpace(t))
		}
	}
	return targets
}

// resolveLinkRedirects resolves MDict @@@LINK cross-references in a definition.
// A def may contain one or more "@@@LINK=<word>" lines (sub-records separated
// by "</>" or newlines); each target is located recursively (depth- and
// cycle-guarded), and multi-target defs concatenate. Fixes #260: the old code
// followed only the first target of a single hop, so multi-target redirects
// showed blank and chained redirects (A→B→C) stopped at the first hop or leaked
// the raw "@@@LINK=" text.
func (dc *DictsController) resolveLinkRedirects(dictId, def string, depth int, visited map[string]bool) string {
	const maxDepth = 8
	if depth > maxDepth {
		return def
	}
	targets := linkTargets(def)
	if len(targets) == 0 {
		return def // not a redirect; show original
	}
	var resolved []string
	for _, target := range targets {
		if visited[target] {
			continue // cycle guard
		}
		visited[target] = true
		result, err := dc.svc.Search(dictId, target)
		if err != nil || len(result) == 0 {
			continue
		}
		targetDef, err := dc.svc.Locate(dictId, result[0])
		if err != nil {
			continue
		}
		resolved = append(resolved, dc.resolveLinkRedirects(dictId, targetDef, depth+1, visited))
	}
	if len(resolved) == 0 {
		return def // nothing resolved; show original rather than blank
	}
	return strings.Join(resolved, "\n")
}

func (dc *DictsController) HandleWordQueryReq(c *gin.Context) {
	// 请求地址: http://localhost:8193/__mdict/__tcidem_query?dict_id=f234356c227f82a54afdaa3514de188a&keyword=card&record_start_offset=20477857&record_end_offset=20501885&key_block_idx=26868
	keyWord := c.Query("keyword")
	recordStart := c.Query("record_start_offset")
	recordEnd := c.Query("record_end_offset")
	dictId := c.Query("dict_id")
	entryId := c.Query("entry_id")
	recordBlockDataStartOffset := c.Query("record_block_data_start_offset")
	recordBlockDataCompressSize := c.Query("record_block_data_compress_size")
	recordBlockDataDeCompressSize := c.Query("record_block_data_decompress_size")
	keyWordDataStartOffset := c.Query("keyword_data_start_offset")
	keyWordDataEndOffset := c.Query("keyword_data_end_offset")

	// 根据词典实际类型决定索引类型，避免 stardict 被当作 mdict 处理（见 #678）
	dictType := "medict"
	if d := dc.svc.GetDictById(dictId); d != nil && d.DictType == string(model.DictTypeStarDict) {
		dictType = "stardict"
	}
	entry, err := convertKeyIndex(dictType, entryId, recordStart, recordEnd, keyWord, recordBlockDataStartOffset, recordBlockDataCompressSize, recordBlockDataDeCompressSize, keyWordDataStartOffset, keyWordDataEndOffset)
	if err != nil {

		log.Errorf("NoRoute REQ ABORT: %s (bad param convert: %s)", c.Request.RequestURI, err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	def, err := dc.svc.Locate(dictId, entry)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// handle @@@Link=${word} redirects (single, multi-target, and chained)
	def = dc.resolveLinkRedirects(dictId, strings.TrimSpace(def), 0, map[string]bool{})

	dict, ok := dc.svc.GetDictPlain(dictId)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	htmlContent, err := handler.WrapContent(dict, entry.MdictKeyWordIndex, def)
	if err != nil {

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, static.ContentTypeHTML, htmlContent)
	return
}

// RenderSnapshot renders a self-contained HTML snapshot of a word's definition
// for offline bookmarking: it re-runs the word-query pipeline (Search → Locate →
// @@@LINK resolve → WrapContent) and inlines every resource (CSS/images/fonts)
// as data: URLs so the snapshot still renders after the source dictionary is
// unloaded. The first Search match is snapshotted; "" if no match. On error the
// caller should fall back to bookmarking without a snapshot.
func (dc *DictsController) RenderSnapshot(dictId, word string) (string, error) {
	if dictId == "" || word == "" {
		return "", fmt.Errorf("snapshot: dictId and word required")
	}
	result, err := dc.svc.Search(dictId, word)
	if err != nil {
		return "", fmt.Errorf("snapshot search %q: %w", word, err)
	}
	if len(result) == 0 {
		return "", nil
	}
	entry := result[0]
	def, err := dc.svc.Locate(dictId, entry)
	if err != nil {
		return "", fmt.Errorf("snapshot locate %q: %w", word, err)
	}
	def = dc.resolveLinkRedirects(dictId, strings.TrimSpace(def), 0, map[string]bool{})

	dict, ok := dc.svc.GetDictPlain(dictId)
	if !ok {
		return "", fmt.Errorf("snapshot: dict %s not found", dictId)
	}
	htmlContent, err := handler.WrapContent(dict, entry.MdictKeyWordIndex, def)
	if err != nil {
		return "", fmt.Errorf("snapshot wrap %q: %w", word, err)
	}

	inlined := handler.InlineResources(htmlContent, func(key string) ([]byte, bool) {
		// mirror innerResourceQuery: dict folder first, then key variants
		if raw, ferr := dc.svc.FindFromDir(dictId, key); ferr == nil {
			return raw, true
		}
		for _, cand := range resourceKeyCandidates(key) {
			if raw, ferr := dc.svc.LookupResource(dictId, cand); ferr == nil {
				return raw, true
			}
		}
		return nil, false
	})
	return string(inlined), nil
}

func (dc *DictsController) HandleResourceQueryReq(c *gin.Context) {
	dictId := c.Query("dict_id")
	rawKeys := strings.SplitN(c.Request.RequestURI, "?", 2)
	if len(rawKeys) < 2 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if dictId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	resourceKey := rawKeys[0]
	resourceKey = strings.TrimPrefix(resourceKey, static.ContentRootUrl+"/")

	dc.innerResourceQuery(c, resourceKey, dictId)
	return
}

func (dc *DictsController) innerResourceQuery(c *gin.Context, key, dictId string) {
	if key == "" || dictId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Debugf("innerResourceQuery key: [%s]", key)

	// 1) dict folder first (css / cover image / etc. shipped alongside the .mdx).
	if raw, err := dc.svc.FindFromDir(dictId, key); err == nil {
		log.Debugf("resource hit dir: [%s]", key)
		dc.serveResource(c, dictId, key, raw)
		return
	}

	// 2) resource lookup across key variants (dicts may store paths with
	//    backslashes or a leading separator).
	for _, candidate := range resourceKeyCandidates(key) {
		if raw, err := dc.svc.LookupResource(dictId, candidate); err == nil {
			log.Debugf("resource hit: [%s]", candidate)
			dc.serveResource(c, dictId, candidate, raw)
			return
		}
	}

	c.AbortWithStatus(http.StatusNotFound)
}

// resourceKeyCandidates returns the key variants to try for a resource lookup:
// the key as-is, then "/" -> "\", then the same with a leading "\". Duplicates
// are removed so a key without slashes isn't looked up three times (#736).
func resourceKeyCandidates(key string) []string {
	seen := make(map[string]struct{})
	cands := make([]string, 0, 3)
	add := func(s string) {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			cands = append(cands, s)
		}
	}
	add(key)
	backslash := strings.ReplaceAll(key, "/", "\\")
	add(backslash)
	if !strings.HasPrefix(backslash, "\\") {
		add("\\" + backslash)
	}
	return cands
}

// serveResource applies the resource pipeline and writes the response. WrapResource
// currently never errors (it transforms in place); if it ever can, fall back to
// the raw bytes rather than dropping the response (#736).
func (dc *DictsController) serveResource(c *gin.Context, dictId, key string, raw []byte) {
	out := raw
	if wrapped, err := handler.WrapResource(dictId, key, raw); err == nil {
		out = wrapped
	}
	wrapContentType(c, key, out)
}

func convertKeyIndex(dictType, entryId, recordStart, recordEnd, keyWord, recordBlockDataStartOffset, recordBlockDataCompressSize, recordBlockDataDeCompressSize, keyWordDataStartOffset, keyWordDataEndOffset string) (*model.KeyQueryIndex, error) {
	if entryId == "" {
		entryId = "0"
	}
	if recordStart == "" {
		recordStart = "0"
	}
	if recordEnd == "" {
		recordEnd = "0"
	}

	idxtype := model.IndexTypeMdict
	if dictType == "stardict" {
		idxtype = model.IndexTypeStardict
	}

	// Parse the 8 integer params via one loop instead of 8× repeat Atoi+err.
	raws := []string{entryId, recordStart, recordEnd, recordBlockDataStartOffset, recordBlockDataCompressSize, recordBlockDataDeCompressSize, keyWordDataStartOffset, keyWordDataEndOffset}
	vals := make([]int64, len(raws))
	for i, s := range raws {
		v, err1 := strconv.Atoi(s)
		if err1 != nil {
			return nil, fmt.Errorf("convertKeyIndex param %d (%q): %w", i, s, err1)
		}
		vals[i] = int64(v)
	}
	queryIndex := &model.KeyQueryIndex{
		IndexType: idxtype,
		MdictKeyWordIndex: &model.MdictKeyWordIndex{
			ID:                            int(vals[0]),
			KeyWord:                       keyWord,
			RecordLocateStartOffset:       vals[1],
			RecordLocateEndOffset:         vals[2],
			RecordBlockDataStartOffset:    vals[3],
			RecordBlockDataCompressSize:   vals[4],
			RecordBlockDataDeCompressSize: vals[5],
			KeyWordDataStartOffset:        vals[6],
			KeyWordDataEndOffset:          vals[7],
		},
	}
	log.Debugf("query index: kw=%s offsets=%v", keyWord, vals)
	return queryIndex, nil
}

// contentTypes maps resource file extensions to their MIME types, looked up by
// the extension of the resource key (#738 — replaced a 16-branch if chain).
var contentTypes = map[string]string{
	".css":   "text/css",
	".js":    "text/javascript",
	".jpeg":  "image/jpeg",
	".jpg":   "image/jpeg",
	".png":   "image/png",
	".gif":   "image/gif",
	".svg":   "image/svg+xml",
	".webp":  "image/webp",
	".mp4":   "video/mp4",
	".wav":   "audio/wav",
	".mp3":   "audio/mpeg", // correct MIME for MP3 (was "audio/mp3")
	".ogg":   "audio/ogg",
	".flac":  "audio/flac",
	".spx":   "audio/speex",
	".ttf":   "font/ttf",
	".otf":   "font/otf",
	".woff":  "font/woff",
	".woff2": "font/woff2",
}

func wrapContentType(c *gin.Context, key string, data []byte) {
	if ct, ok := contentTypes[strings.ToLower(filepath.Ext(key))]; ok {
		c.Data(http.StatusOK, ct, data)
		return
	}
	// Unknown extension: sniff from the content instead of hard-failing with 415.
	c.Data(http.StatusOK, http.DetectContentType(data), data)
}
