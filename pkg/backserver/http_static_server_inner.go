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
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"github.com/terasum/medict/internal/static"
	"net"
	"net/http"
	"strings"
)

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

	log.Infof("start listening... %s", srv.Addr)

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

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
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
				log.Infof("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

func (bs *BackServer) setUpRouters() error {
	bs.GinEngine.Use(cors())
	bs.GinEngine.NoRoute(func(c *gin.Context) {
		log.Infof("NoRoute REQ REQUEST URI: %s\n", c.Request.RequestURI)

		if strings.HasPrefix(c.Request.URL.String(), static.ContentRootUrl+static.WordQueryMagicPath) {
			bs.DictCon.HandleWordQueryReq(c)
			return
		} else {
			bs.DictCon.HandleResourceQueryReq(c)
			return
		}
	})
	return nil
}

func (bs *BackServer) setupHandlers() {
	bs.handlerMap.Store("GetAllDicts", bs.DictCon.GetAllDicts)
	bs.handlerMap.Store("SearchWord", bs.DictCon.SearchWord)
}
