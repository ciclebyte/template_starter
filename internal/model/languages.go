package model

type LanguagesInfo struct {
	Id          int    `orm:"id"  json:"id"`                    // 语言ID，自增主键
	Name        string `orm:"name"  json:"name"`                // 语言名称（如JavaScript、Python）
	DisplayName string `orm:"display_name"  json:"displayName"` // 语言显示名称
	Code        string `orm:"code"  json:"code"`                // 语言代码（如js、py）
	Icon        string `orm:"icon"  json:"icon"`                // 语言图标标识或URL
	Color       string `orm:"color"  json:"color"`              // 语言代表色（十六进制）
	IsPopular   int    `orm:"is_popular"  json:"isPopular"`     // 是否热门语言
}
