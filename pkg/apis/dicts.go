package apis

import (
	"encoding/base64"
	"github.com/terasum/medict/internal/config"
	"github.com/terasum/medict/pkg/model"
	"github.com/terasum/medict/pkg/service"
	"strings"
)

type DictsAPI struct {
	dictSvc *service.DictService
}

func NewDictsAPI(config *config.Config) (*DictsAPI, error) {
	svc, err := service.NewDictService(config)
	if err != nil {
		return nil, err
	}
	_, _, err = svc.AllWordIndexing()
	if err != nil {
		return nil, err
	}
	return &DictsAPI{dictSvc: svc}, nil
}

func (dict *DictsAPI) Dicts() []*model.PlainDictItem {
	return dict.dictSvc.Dicts()
}

func (dict *DictsAPI) Suggest(word string) ([]*model.WrappedWordItem, error) {
	simWords, err := dict.dictSvc.SimWords(word)
	if err != nil {
		return nil, err
	}
	return simWords, nil
}
func (dict *DictsAPI) LookupDefinition(dictId string, rawKeyWord string, recordStart uint64) (string, error) {
	def, err := dict.dictSvc.LookupDefinition(dictId, rawKeyWord, recordStart)
	if err != nil {
		return "", err
	}
	def = strings.Replace(def, `\"`, `"`, -1)
	def = `<html><body>` + def + `</body></html>`
	encodedStr := base64.StdEncoding.EncodeToString([]byte(def))
	return encodedStr, nil
}

func (dict *DictsAPI) Destroy() {
	dict.dictSvc.GC()
}
