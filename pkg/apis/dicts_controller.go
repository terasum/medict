package apis

import (
	"encoding/hex"
	"fmt"
	"github.com/terasum/medict/internal/static/handler"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/terasum/medict/internal/static"
	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service"
)

var log = logging.MustGetLogger("default")

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

	entry, err := convertKeyIndex("medict", entryId, recordStart, recordEnd, keyWord, recordBlockDataStartOffset, recordBlockDataCompressSize, recordBlockDataDeCompressSize, keyWordDataStartOffset, keyWordDataEndOffset)
	if err != nil {

		fmt.Printf("NoRoute REQ ABORT: %s (%s:%s)\n", c.Request.RequestURI, "bad param convert", err.Error())
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

	ientryId, err := strconv.Atoi(entryId)
	if err != nil {
		return nil, err
	}
	irecordStart, err := strconv.Atoi(recordStart)
	if err != nil {
		return nil, err
	}
	irecordEnd, err := strconv.Atoi(recordEnd)
	if err != nil {
		return nil, err
	}

	iRecordBlockDataStartOffset, err := strconv.Atoi(recordBlockDataStartOffset)
	if err != nil {
		return nil, err
	}
	iRecordBlockDataCompressSize, err := strconv.Atoi(recordBlockDataCompressSize)
	if err != nil {
		return nil, err
	}
	iRecordBlockDataDeCompressSize, err := strconv.Atoi(recordBlockDataDeCompressSize)
	if err != nil {
		return nil, err
	}
	iKeyWordDataStartOffset, err := strconv.Atoi(keyWordDataStartOffset)
	if err != nil {
		return nil, err
	}
	iKeyWordDataEndOffset, err := strconv.Atoi(keyWordDataEndOffset)
	if err != nil {
		return nil, err
	}

	queryIndex := &model.KeyQueryIndex{
		IndexType: idxtype,
		MdictKeyWordIndex: &model.MdictKeyWordIndex{
			ID:                            ientryId,
			KeyWord:                       keyWord,
			RecordLocateStartOffset:       int64(irecordStart),
			RecordLocateEndOffset:         int64(irecordEnd),
			RecordBlockDataStartOffset:    int64(iRecordBlockDataStartOffset),
			RecordBlockDataCompressSize:   int64(iRecordBlockDataCompressSize),
			RecordBlockDataDeCompressSize: int64(iRecordBlockDataDeCompressSize),
			KeyWordDataStartOffset:        int64(iKeyWordDataStartOffset),
			KeyWordDataEndOffset:          int64(iKeyWordDataEndOffset),
		},
	}
	log.Infof("KeyWord: %s", keyWord)
	log.Infof("RecordLocateStartOffset: %d", int64(irecordStart))
	log.Infof("RecordLocateEndOffset: %d", int64(irecordEnd))
	log.Infof("RecordBlockDataStartOffset: %d", int64(iRecordBlockDataStartOffset))
	log.Infof("RecordBlockDataCompressSize: %d", int64(iRecordBlockDataCompressSize))
	log.Infof("RecordBlockDataDeCompressSize: %d", int64(iRecordBlockDataDeCompressSize))
	log.Infof("KeyWordDataStartOffset: %d", int64(iKeyWordDataStartOffset))
	log.Infof("KeyWordDataEndOffset: %d", int64(iKeyWordDataEndOffset))
	return queryIndex, nil
}

func wrapContentType(c *gin.Context, key string, data []byte) {
	if strings.HasSuffix(key, ".css") {
		c.Data(http.StatusOK, "text/css", data)
	} else if strings.HasSuffix(key, ".js") {
		c.Data(http.StatusOK, "text/javascript", data)
	} else if strings.HasSuffix(key, ".jpeg") {
		c.Data(http.StatusOK, "image/jpeg", data)
	} else if strings.HasSuffix(key, ".png") {
		c.Data(http.StatusOK, "image/png", data)
	} else if strings.HasSuffix(key, ".gif") {
		c.Data(http.StatusOK, "image/gif", data)
	} else if strings.HasSuffix(key, ".jpg") {
		c.Data(http.StatusOK, "image/jpeg", data)
	} else if strings.HasSuffix(key, ".svg") {
		c.Data(http.StatusOK, "image/svg+xml", data)
	} else if strings.HasSuffix(key, ".webp") {
		c.Data(http.StatusOK, "image/webp", data)
	} else if strings.HasSuffix(key, ".mp4") {
		c.Data(http.StatusOK, "video/mp4", data)
	} else if strings.HasSuffix(key, ".wav") {
		c.Data(http.StatusOK, "audio/wav", data)
	} else if strings.HasSuffix(key, ".mp3") {
		c.Data(http.StatusOK, "audio/mp3", data)
	} else if strings.HasSuffix(key, ".ogg") {
		c.Data(http.StatusOK, "audio/ogg", data)
	} else if strings.HasSuffix(key, ".ttf") {
		c.Data(http.StatusOK, "font/ttf", data)
	} else if strings.HasSuffix(key, ".otf") {
		c.Data(http.StatusOK, "font/otf", data)
	} else if strings.HasSuffix(key, ".woff") {
		c.Data(http.StatusOK, "font/woff", data)
	} else if strings.HasSuffix(key, ".woff2") {
		c.Data(http.StatusOK, "font/woff", data)
	} else {
		c.AbortWithStatus(http.StatusUnsupportedMediaType)
	}
}
