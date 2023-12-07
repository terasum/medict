package mdict_idxer

import (
	"github.com/op/go-logging"
	"time"
)

var log = logging.MustGetLogger("default")

func logstart(method string, args interface{}) time.Time {
	log.Infof("%s|STR|: %v", method, args)
	return time.Now()
}

func logend(method string, startTime time.Time, err error) {
	log.Infof("%s|END|%s|E:%v", method, time.Now().Sub(startTime).String(), err)
}
