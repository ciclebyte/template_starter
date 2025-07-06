package controller

import (
	"context"

	indexApi "github.com/ciclebyte/template_starter/api/v1/index"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

type IndexController struct{}

var Index = IndexController{}

// Index 首页数据
func (c *IndexController) Index(ctx context.Context, req *indexApi.IndexReq) (res *indexApi.IndexRes, err error) {
	g.Log().Debug(ctx, "Index called with req:", req)

	// 设置默认值
	if req.CategoryLimit <= 0 {
		req.CategoryLimit = 6
	}
	if req.FeaturedLimit <= 0 {
		req.FeaturedLimit = 8
	}

	// 调用service层获取数据
	data, err := service.Index().GetIndexData(ctx, req)
	if err != nil {
		g.Log().Error(ctx, "GetIndexData failed:", err)
		return nil, err
	}

	return data, nil
}
