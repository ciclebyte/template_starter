package common

import (
	"github.com/ciclebyte/template_starter/internal/model"
)

// PageReq 公共请求参数
type PageReq struct {
	model.PageReq
}

type Author struct {
	Authorization string `p:"Authorization" in:"header" dc:"Bearer {{token}}"`
}
