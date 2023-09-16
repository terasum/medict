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

package service

import (
	"context"
	"github.com/terasum/medict/internal/config"
)

type StaticService struct {
	cfg          *config.Config
	ctx          context.Context
	errorChannel chan error
	stopChannel  chan int
	ready        bool

	httpServer *HttpServer
}

func NewStaticServer(ctx context.Context, conf *config.Config, dictService *DictService, errorChannel chan error, stopChannel chan int) *StaticService {
	return &StaticService{
		cfg:          conf,
		ctx:          ctx,
		errorChannel: errorChannel,
		stopChannel:  stopChannel,
		httpServer: &HttpServer{
			DictService: dictService,
			Config:      conf,
			StopChan:    stopChannel,
			ErrChan:     errorChannel,
		},
	}
}

func (si *StaticService) Start() {
	si.ready = true
	si.httpServer.StartStaticServer()
}

func (si *StaticService) StaticServerBaseUrl() string {
	return si.httpServer.BaseUrl()
}
