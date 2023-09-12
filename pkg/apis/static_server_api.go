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
