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
	"github.com/terasum/medict/pkg/apis"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/static"
	"github.com/terasum/medict/pkg/service"
)

type BackServer struct {
	Config *config.Config
	ctx    context.Context

	StopChan chan int
	ErrChan  chan error

	Srv        *http.Server
	ListenAddr net.Addr

	Ready   bool
	DevMode bool

	GinEngine *gin.Engine
	DictCon   *apis.DictsController
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
// handlers/routes. The DictService is owned by App and passed in (issue #727).
func (bs *BackServer) SetUp(dictsSvc *service.DictService) error {
	bs.DictCon = apis.NewDictsController(dictsSvc)

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
