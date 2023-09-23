package model

type DictType string
type ImgType string

const DictTypeMdict DictType = "Mdict"
const DictTypeStarDict DictType = "StarDict"
const ImgTypeJPG ImgType = "jpg"
const ImgTypePNG ImgType = "png"

type GeneralDictionary interface {
	BuildIndex() error
	// DictType 返回词典类型, 目前仅支持 mdict 和 stardict
	DictType() DictType
	// Description 返回词典描述
	Description() *PlainDictionaryInfo
	// Name 返回词典名称(主要是为了解决mdict的外部资源问题)
	Name() string
	// Lookup 直接输入 keyword 遍历词典搜索
	Lookup(keyword string) ([]byte, error)
	// LookupResource 搜索词典资源，包括 css/js/图片/字体等, 仅限制支持在词典所在文件夹
	LookupResource(keyword string) ([]byte, error)
	// Locate 使用索引定位词条，并返回 html 释义
	Locate(entry *KeyIndex) ([]byte, error)
	// Search 返回近似词条索引列表，用于后续 Locate
	Search(keyword string) ([]*KeyIndex, error)
}
