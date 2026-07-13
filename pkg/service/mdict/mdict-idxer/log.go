package mdict_idxer

import (
	"github.com/op/go-logging"
	"time"
)

var log = logging.MustGetLogger("mdict.idxer")

func logstart(method string, args interface{}) time.Time {
	log.Debugf("%s|STR|: %v", method, args)
	return time.Now()
}

func logend(method string, startTime time.Time, err error) {
	log.Debugf("%s|END|%s|E:%v", method, time.Now().Sub(startTime).String(), err)
}
