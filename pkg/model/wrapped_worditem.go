package model

type WrappedWordItem struct {
	DictId string `json:"dict_id"`
	*WordItem
}
