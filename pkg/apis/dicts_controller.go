package apis

import (
	"encoding/hex"
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
	ds *service.DictService
}

func NewDictsController(ds *service.DictService) *DictsController {
	return &DictsController{ds: ds}
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
	if d := dc.ds.GetDictById(dictId); d != nil && d.DictType == string(model.DictTypeStarDict) {
		dictType = "stardict"
	}
	entry, err := convertKeyIndex(dictType, entryId, recordStart, recordEnd, keyWord, recordBlockDataStartOffset, recordBlockDataCompressSize, recordBlockDataDeCompressSize, keyWordDataStartOffset, keyWordDataEndOffset)
	if err != nil {

		log.Errorf("NoRoute REQ ABORT: %s (bad param convert: %s)", c.Request.RequestURI, err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	def, err := dc.ds.Locate(dictId, entry)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// handle @@@Link=${word}
	def = strings.TrimSpace(def)

	if strings.HasPrefix(def, "@@@LINK=") {
		log.Infof("search @@@LINK=>[%s], hex:[%s]", def, hex.EncodeToString([]byte(def)))
		newWord := strings.TrimPrefix(def, "@@@LINK=")
		newWord = strings.TrimRight(newWord, "\r\n\000")
		result, err1 := dc.ds.Search(dictId, newWord)
		if err1 == nil && len(result) > 0 {
			newEntry := result[0]
			def1, err2 := dc.ds.Locate(dictId, newEntry)
			// handle link jump
			if err2 == nil {
				def = def1
			} else {
				log.Errorf("search @@@link jump failed %s, %v", def, err2)
			}
		}
	}

	dict, ok := dc.ds.GetDictPlain(dictId)
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
	log.Infof("innerResourceQuery(0) search key: [%s]", key)

	// 先从文件夹搜索
	resultBytes, err := dc.ds.FindFromDir(dictId, key)
	if err == nil {
		log.Infof("innerResourceQuery search key hit dir: [%s]", key)
		resultBytes, err = handler.WrapResource(dictId, key, resultBytes)
		if err != nil {
			wrapContentType(c, key, resultBytes)
			return
		}
		wrapContentType(c, key, resultBytes)
		return
	}

	log.Infof("innerResourceQuery(1) search from resource [%s]", key)
	// 再从文件资源中搜索
	resultBytes, err = dc.ds.LookupResource(dictId, key)
	if err == nil {
		log.Infof("innerResourceQuery search key hit resource: [%s]", key)
		resultBytes, err = handler.WrapResource(dictId, key, resultBytes)
		if err != nil {
			wrapContentType(c, key, resultBytes)
			return
		}

		wrapContentType(c, key, resultBytes)
		return
	}

	key = strings.ReplaceAll(key, "/", "\\")
	log.Infof("innerResourceQuery(2) search from resource [%s]", key)
	resultBytes, err = dc.ds.LookupResource(dictId, key)
	if err == nil {
		log.Infof("innerResourceQuery search key hit resource: [%s]", key)
		resultBytes, err = handler.WrapResource(dictId, key, resultBytes)
		if err != nil {
			wrapContentType(c, key, resultBytes)
			return
		}

		wrapContentType(c, key, resultBytes)
		return
	}

	// 补全路径重新搜索
	if !strings.HasPrefix(key, "\\") {
		key = "\\" + key
	}

	log.Infof("innerResourceQuery(3) search from resource [%s]", key)

	resultBytes, err = dc.ds.LookupResource(dictId, key)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	log.Infof("innerResourceQuery(4) search key hit resource with '\\' prefix: [%s]", key)
	result, err := handler.WrapResource(dictId, key, resultBytes)
	if err != nil {
		wrapContentType(c, key, result)
		return
	}
	wrapContentType(c, key, result)
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
