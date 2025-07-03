// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Languages is the golang structure of table languages for DAO operations like Where/Data.
type Languages struct {
	g.Meta      `orm:"table:languages, do:true"`
	Id          interface{} // 语言ID，自增主键
	Name        interface{} // 语言名称（如JavaScript、Python）
	DisplayName interface{} // 语言显示名称
	Code        interface{} // 语言代码（如js、py）
	Icon        interface{} // 语言图标标识或URL
	Color       interface{} // 语言代表色（十六进制）
	IsPopular   interface{} // 是否热门语言
	CreatedAt   *gtime.Time // 记录创建时间
	UpdatedAt   *gtime.Time // 记录最后更新时间
}
