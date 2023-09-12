package tmpl

import (
	"fmt"
	"github.com/terasum/medict/pkg/model"
)

type ContentPreHandlePipeline struct {
	replacers []Replacer
	handlers  []Handler
}

var handler = &ContentPreHandlePipeline{
	replacers: []Replacer{
		&ReplacerCss{},
		&ReplacerJs{},
		&ReplacerImage{},
		&ReplacerSound{},
	},
	handlers: []Handler{
		&HandlerFont{},
	},
}

func WrapContent(dictId string, keyEntry *model.KeyBlockEntry, definition string) ([]byte, error) {
	content := handleContent(dictId, keyEntry, definition)
	return []byte(fmt.Sprintf(WordDefinitionTempl, content)), nil
}

func WrapResource(dictId string, keyWord string, resource []byte) ([]byte, error) {
	content := handleResource(dictId, keyWord, resource)
	return content, nil
}

func handleContent(dictId string, keyEntry *model.KeyBlockEntry, definition string) string {
	for _, replacer := range handler.replacers {
		keyEntry, definition = replacer.Replace(dictId, keyEntry, definition)
	}
	return definition
}

func handleResource(dictId string, keyWord string, resource []byte) []byte {
	for _, han := range handler.handlers {
		if han.Match(dictId, keyWord) {
			fmt.Printf("resource handle matched [%s](%s)\n", keyWord, dictId)
			keyWord, resource = han.Replace(dictId, keyWord, resource)
		}
	}
	return resource
}
