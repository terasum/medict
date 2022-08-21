package apis

import (
	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/internal/static"
	"strconv"
)

type StaticInfos struct {
	cfg *config.Config
}

func NewStaticInfosApi(cfg *config.Config) *StaticInfos {
	return &StaticInfos{cfg: cfg}
}

func (si *StaticInfos) StaticSrvUrl() string {
	return "http://localhost:" + strconv.Itoa(si.cfg.StaticServerPort) + static.ContentRootUrl
}
