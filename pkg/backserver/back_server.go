//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package backserver

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/static"
	"github.com/terasum/medict/pkg/apis"
	"github.com/terasum/medict/pkg/service"
)

var log = logging.MustGetLogger("backserver")

type BackServer struct {
	Config *config.Config
	ctx    context.Context

	StopChan chan int
	ErrChan  chan error

	Srv        *http.Server
	ListenAddr net.Addr

	Ready   bool
	DevMode bool

	GinEngine  *gin.Engine
	Controller *apis.DictsController
}

func NewStaticServer(conf *config.Config) (*BackServer, error) {
	r := gin.Default()
	bs := &BackServer{
		Config:    conf,
		DevMode:   false,
		GinEngine: r,
	}
	return bs, nil
}

// SetUp wires the injected DictService into the controller and registers
// routes. The DictService is owned by App and passed in (issue #727).
func (bs *BackServer) SetUp(dictsSvc *service.DictService) error {
	bs.Controller = apis.NewDictsController(dictsSvc)

	if err := bs.setUpRouters(); err != nil {
		return err
	}
	return nil
}

func (bs *BackServer) SetDebug() {
	bs.DevMode = true
}

func (bs *BackServer) Start() {
	if bs.DevMode {
		bs.startStaticServer("localhost:9081")
	} else {
		bs.startStaticServer("localhost:0")
	}
}

func (bs *BackServer) GracefulStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	if bs == nil || bs.Srv == nil {
		return
	}
	if err := bs.Srv.Shutdown(ctx); err != nil {
		log.Errorf("Server Shutdown: %s", err)
	}
	// catching ctx.Done(). timeout of 0.4 seconds.
	select {
	case <-ctx.Done():
		log.Info("timeout of 0.4 seconds.")
	}
	log.Info("Server exiting")
}

func (bs *BackServer) StaticServerBaseUrl() string {
	//return "http://localhost:" + strconv.Itoa(bs.Config.StaticServerPort) + static.ContentRootUrl
	listenAddr := ""
	if bs.Ready {
		listenAddr = bs.ListenAddr.String() + static.ContentRootUrl
		if !strings.HasPrefix(listenAddr, "http://") {
			return "http://" + listenAddr
		}
		return listenAddr
	}
	return ""
}

func (bs *BackServer) startStaticServer(listenAddr string) {
	if listenAddr == "" {
		listenAddr = "localhost:0"
	}
	if listenAddr == ":0" {
		listenAddr = "localhost:0"
	}

	srv := &http.Server{
		Addr:    listenAddr, // use next port available
		Handler: bs.GinEngine,
	}

	log.Infof("start listening... %s\n", srv.Addr)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil && err != http.ErrServerClosed {
		log.Infof("backserver listen err: %s\n", err)
	} else {
		log.Infof("backserver listen at: %s\n", ln.Addr().String())
	}

	bs.ListenAddr = ln.Addr()

	go func() {
		bs.Ready = true
		err = srv.Serve(ln)
		if err != nil && err != http.ErrServerClosed {
			log.Infof("backserver serve err: %s\n", err)
		}
	}()

	bs.Srv = srv
}

// allowedOrigin reports whether the given Origin is permitted to make
// cross-origin requests to the embedded resource server.
//
// 正式功能不依赖跨域：前端控制面走 Wails IPC，释义 HTML 及其子资源经同源
// iframe（http://localhost:<port>）加载，本身不需要 CORS。此处仅放行本地
// wails webview 与 loopback 调试来源，避免原实现反射任意 Origin 的安全隐患。
func allowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	switch origin {
	case "http://wails.localhost", "https://wails.localhost",
		"wails://wails.localhost", "wails://localhost":
		return true
	}
	for _, scheme := range []string{"http", "https"} {
		for _, host := range []string{"localhost", "127.0.0.1"} {
			if origin == scheme+"://"+host || strings.HasPrefix(origin, scheme+"://"+host+":") {
				return true
			}
		}
	}
	return false
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if allowedOrigin(origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Infof("Panic info is: %v\n", err)
			}
		}()

		c.Next()
	}
}

func (bs *BackServer) setUpRouters() error {
	bs.GinEngine.Use(cors())

	// Word lookup: an explicit route (previously a string-HasPrefix branch
	// inside NoRoute — issue #728). Discoverable, middleware-able, structured
	// errors via the handler instead of catch-all dispatch.
	bs.GinEngine.GET(static.ContentRootUrl+static.WordQueryMagicPath, bs.Controller.HandleWordQueryReq)

	// Everything else is a resource lookup. Resource paths are arbitrary
	// (css / images / fonts / ... under the content root), so they stay a
	// catch-all rather than enumerated routes.
	bs.GinEngine.NoRoute(func(c *gin.Context) {
		log.Debugf("resource request: %s", c.Request.RequestURI)
		bs.Controller.HandleResourceQueryReq(c)
	})
	return nil
}
