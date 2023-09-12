package tmpl

import "github.com/terasum/medict/pkg/model"

type Replacer interface {
	Replace(dictId string, entry *model.KeyBlockEntry, htmlContent string) (*model.KeyBlockEntry, string)
}
