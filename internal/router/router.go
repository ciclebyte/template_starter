package router

import (
	"context"

	"github.com/ciclebyte/template_starter/internal/controller"
	"github.com/ciclebyte/template_starter/internal/service"
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
			controller.TemplateFiles,
			controller.BuiltinFunctions,
			controller.SprigFunctions,
			controller.Index,
			controller.Statistics,
			controller.SystemConfig,
			controller.AI,
		)

		// 手动注册流式AI聊天端点
		group.POST("/ai/chat/stream", controller.AI.ChatStream)
		
		// 手动注册OpenAI兼容端点
		group.POST("/v1/chat/completions", controller.OpenAI.ChatCompletions)

		//自动绑定定义的控制器
		if err := libRouter.RouterAutoBind(ctx, router, group); err != nil {
			panic(err)
		}
	})
}
