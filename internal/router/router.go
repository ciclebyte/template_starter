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
		
		// 认证相关路由 - 使用OptionalAuth中间件，在控制器方法中处理认证检查
		group.Middleware(service.Middleware().OptionalAuth)
		group.Bind(controller.Auth)
		
		// 权限管理路由 (需要认证)
		group.Group("", func(group *ghttp.RouterGroup) {
			group.Middleware(service.Middleware().RequireAuth)
			group.Bind(controller.Permission)
		})
		
		// 其他公开访问路由 (可选认证)
		group.Bind(
			controller.Languages,
			controller.Categories,
			controller.Templates,
			controller.TemplateLanguages,
			controller.TemplateFiles,
			controller.BuiltinFunctions,
			controller.SprigFunctions,
			controller.Index,
			controller.Statistics,
			controller.SystemConfig,
			controller.AI,
			controller.Tags,
			controller.VarPreset,
			controller.TemplateExpose,
			controller.TemplateVariablePresets,
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
