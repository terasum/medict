package static

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/static/tmpl"
	"log"

	"net/http"
	"strings"
)

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

const ContentRootUrl = "/64ca4043daabdc41"

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
	r.Static("/", "./frontend/dist")
	r.NoRoute(func(c *gin.Context) {
		fmt.Printf("NoRoute REQ: %s\n", c.Request.URL.String())
		fmt.Printf("NoRoute REQ: %s\n", c.Request.RequestURI)
		if !strings.HasPrefix(c.Request.URL.String(), ContentRootUrl) {
			c.AbortWithStatus(http.StatusNotFound)
			fmt.Printf("NoRoute REQ ABORT: %s\n", c.Request.RequestURI)
			return
		}
		origin := c.Request.Header.Get("Origin") //请求头部

		rawKeyWord := c.Query("raw_key_word")  // shortcut for c.Request.URL.Query().Get("lastname")
		recordStart := c.Query("record_start") // shortcut for c.Request.URL.Query().Get("lastname")
		if rawKeyWord == "" {
			rawKeyWord = "null"
		}
		if recordStart == "" {
			recordStart = "000"
		}
		result, err := c.Cookie("raw_key_word")
		if err == nil {
			fmt.Printf("cookie:raw_key_word %s\n", result)
		} else {
			fmt.Printf("cookie:raw_key_word null %s\n", err.Error())
		}
		result2, err := c.Cookie("record_start")
		if err == nil {
			fmt.Printf("cookie:record_start %s\n", result2)
		} else {
			fmt.Printf("cookie:record_start null, %s\n", err.Error())
		}

		c.SetCookie("raw_key_word", rawKeyWord, 3600, "/", origin, false, false)
		c.SetCookie("raw_key_word", rawKeyWord, 3600, "/", "localhost", false, false)
		c.SetCookie("record_start", recordStart, 3600, "/", origin, false, false)
		c.SetCookie("record_start", recordStart, 3600, "/", "localhost", false, false)
		c.Data(http.StatusOK, ContentTypeHTML, []byte(fmt.Sprintf(tmpl.WordContainer,
			"raw_key_word", rawKeyWord,
			"record_start", recordStart,
			rawKeyWord, recordStart)))
		return
	})

	// err := r.Run("localhost:" + strconv.Itoa(cfg.StaticServerPort))
	err := r.Run("localhost:19191")
	if err != nil {
		panic(err)
		return
	}
}
