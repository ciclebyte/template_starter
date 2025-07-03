// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Languages is the golang structure for table languages.
type Languages struct {
	Id          int         `json:"id"          description:"语言ID，自增主键"`
	Name        string      `json:"name"        description:"语言名称（如JavaScript、Python）"`
	DisplayName string      `json:"displayName" description:"语言显示名称"`
	Code        string      `json:"code"        description:"语言代码（如js、py）"`
	Icon        string      `json:"icon"        description:"语言图标标识或URL"`
	Color       string      `json:"color"       description:"语言代表色（十六进制）"`
	IsPopular   int         `json:"isPopular"   description:"是否热门语言"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:"记录创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:"记录最后更新时间"`
}
