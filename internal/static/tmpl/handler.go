package tmpl

type Handler interface {
	Replace(dictId string, entry string, resource []byte) (string, []byte)
	Match(dictId string, key string) bool
}
