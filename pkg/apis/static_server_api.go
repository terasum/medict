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

package apis

import (
	"strconv"

	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/static"
)

type StaticInfos struct {
	cfg *config.Config
}

func NewStaticServersApi(cfg *config.Config) *StaticInfos {
	return &StaticInfos{cfg: cfg}
}

func (si *StaticInfos) StaticServerUrl() string {
	return "http://localhost:" + strconv.Itoa(si.cfg.StaticServerPort) + static.ContentRootUrl
}
