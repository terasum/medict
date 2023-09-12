package static

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/static/tmpl"
	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

const ContentRootUrl = "/__mdict"

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

func StartStaticServer(cfg *config.Config) {
	r := gin.Default()
	r.Use(Cors()) //开启中间件 允许使用跨域请求

	var err error
	var dictService *service.DictService

	dictService, err = service.NewDictService(cfg)
	if err != nil {
		panic(err)
	}

	r.NoRoute(func(c *gin.Context) {
		fmt.Printf("NoRoute REQ: %s\n", c.Request.URL.String())
		fmt.Printf("NoRoute REQ: %s\n", c.Request.RequestURI)
		if !strings.HasPrefix(c.Request.URL.String(), ContentRootUrl+"/__tcidem_query") {
			handleAssetsQueryReq(c, dictService)
			return
		}

		handleWordQueryReq(c, dictService)
		return
	})

	err = r.Run("localhost:" + strconv.Itoa(cfg.StaticServerPort))
	if err != nil {
		panic(err)
		return
	}
}

func handleAssetsQueryReq(c *gin.Context, dictService *service.DictService) {
	dictId := c.Query("dict_id")
	rawkeys := strings.SplitN(c.Request.RequestURI, "?", 2)
	if len(rawkeys) < 1 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	resourceKey := rawkeys[0]
	resourceKey = strings.TrimPrefix(resourceKey, ContentRootUrl+"/")

	handleWordResourceReq(c, dictService, resourceKey, dictId)
	return
}

func handleWordQueryReq(c *gin.Context, dictService *service.DictService) {
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

	def, err := dictService.Locate(dictId, entry)
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

func handleWordResourceReq(c *gin.Context, dictService *service.DictService, key, dictId string) {
	if key == "" || dictId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if !strings.Contains(key, "/") {
		fmt.Printf("NoRoute[0] handle resource key: %s from dir\n", key)
		resultBytes, err := dictService.FindFromDir(dictId, key)
		if err == nil {
			resultBytes, err = tmpl.WrapResource(dictId, key, resultBytes)
			if err != nil {
				WrapContentType(c, key, resultBytes)
				return
			}
			WrapContentType(c, key, resultBytes)
			return
		}

		resultBytes, err = dictService.LookupResource(dictId, key)
		if err == nil {
			resultBytes, err = tmpl.WrapResource(dictId, key, resultBytes)
			if err != nil {
				WrapContentType(c, key, resultBytes)
				return
			}

			WrapContentType(c, key, resultBytes)
			return
		}

	}
	// 补全路径重新搜索
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}

	key = strings.ReplaceAll(key, "/", "\\")
	result, err := dictService.LookupResource(dictId, key)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	result, err = tmpl.WrapResource(dictId, key, result)
	if err != nil {
		WrapContentType(c, key, result)
		return
	}
	WrapContentType(c, key, result)
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

func WrapContentType(c *gin.Context, key string, data []byte) {
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
