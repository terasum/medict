package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/terasum/medict/internal/static/tmpl"
	"github.com/terasum/medict/pkg/model"
	"net/http"
	"strconv"
	"strings"
)

func (server *HttpServer) handleResourceQueryReq(c *gin.Context) {
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
	resourceKey = strings.TrimPrefix(resourceKey, ContentRootUrl+"/")

	server.innerResourceQuery(c, resourceKey, dictId)
	return
}

func (server *HttpServer) handleWordQueryReq(c *gin.Context) {
	// 请求地址: http://localhost:8193/__mdict/__tcidem_query?dict_id=f234356c227f82a54afdaa3514de188a&key_word=card&record_start_offset=20477857&record_end_offset=20501885&key_block_idx=26868
	keyWord := c.Query("key_word")
	recordStart := c.Query("record_start_offset")
	recordEnd := c.Query("record_end_offset")
	dictId := c.Query("dict_id")
	entryId := c.Query("entry_id")
	keyBlockIdx := c.Query("key_block_idx")

	entry, err := convertKeyBlockEntry(entryId, recordStart, recordEnd, keyWord, keyBlockIdx)
	if err != nil {

		fmt.Printf("NoRoute REQ ABORT: %s (%s:%s)\n", c.Request.RequestURI, "bad param convert", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	def, err := server.DictService.Locate(dictId, entry)
	if err != nil {

		c.AbortWithStatus(http.StatusInternalServerError)
		return

	}
	htmlContent, err := tmpl.WrapContent(dictId, entry, def)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, ContentTypeHTML, htmlContent)
	return
}

func (server *HttpServer) innerResourceQuery(c *gin.Context, key, dictId string) {
	if key == "" || dictId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if !strings.Contains(key, "/") {
		fmt.Printf("NoRoute[0] handle resource key: %s from dir\n", key)
		resultBytes, err := server.DictService.FindFromDir(dictId, key)
		if err == nil {
			resultBytes, err = tmpl.WrapResource(dictId, key, resultBytes)
			if err != nil {
				wrapContentType(c, key, resultBytes)
				return
			}
			wrapContentType(c, key, resultBytes)
			return
		}

		resultBytes, err = server.DictService.LookupResource(dictId, key)
		if err == nil {
			resultBytes, err = tmpl.WrapResource(dictId, key, resultBytes)
			if err != nil {
				wrapContentType(c, key, resultBytes)
				return
			}

			wrapContentType(c, key, resultBytes)
			return
		}

	}
	// 补全路径重新搜索
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}

	key = strings.ReplaceAll(key, "/", "\\")
	result, err := server.DictService.LookupResource(dictId, key)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	result, err = tmpl.WrapResource(dictId, key, result)
	if err != nil {
		wrapContentType(c, key, result)
		return
	}
	wrapContentType(c, key, result)
	return
}

func convertKeyBlockEntry(entryId, recordStart, recordEnd, keyWord, keyBlockIdx string) (*model.KeyBlockEntry, error) {
	if entryId == "" {
		entryId = "0"
	}
	if recordStart == "" {
		recordStart = "0"
	}
	if recordEnd == "" {
		recordEnd = "0"
	}
	if keyBlockIdx == "" {
		keyBlockIdx = "0"
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
	ikeyBlockIdx, err := strconv.Atoi(keyBlockIdx)
	if err != nil {
		return nil, err
	}

	return &model.KeyBlockEntry{
		ID:                ientryId,
		RecordStartOffset: int64(irecordStart),
		RecordEndOffset:   int64(irecordEnd),
		KeyWord:           keyWord,
		KeyBlockIdx:       int64(ikeyBlockIdx),
	}, nil
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
