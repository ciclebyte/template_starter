package router

import (
	"context"

	controller "github.com/ciclebyte/template_starter/internal/controller"
	service "github.com/ciclebyte/template_starter/internal/service"
	"github.com/ciclebyte/template_starter/library/libRouter"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Router struct{}

func (router *Router) BindController(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().MiddlewareCORS)
		group.Bind(
			controller.Languages,
			controller.Categories,
			controller.Templates,
			controller.TemplateVariables,
			controller.TemplateLanguages,
		)

		//自动绑定定义的控制器
		if err := libRouter.RouterAutoBind(ctx, router, group); err != nil {
			panic(err)
		}
	})
}
