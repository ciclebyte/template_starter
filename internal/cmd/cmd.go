package cmd

import (
	"context"
	"strings"

	_ "github.com/ciclebyte/template_starter/internal/logic"
	"github.com/ciclebyte/template_starter/internal/router"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				r := &router.Router{}
				r.BindController(ctx, group)
				
				// 添加SPA路由回退支持，处理Vue Router的HTML5 History模式
				group.Hook("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
					path := r.URL.Path
					
					// 如果是API请求，跳过SPA回退
					if strings.HasPrefix(path, "/api/") {
						return
					}
					
					// 如果是静态资源文件（有文件扩展名），跳过SPA回退
					if strings.Contains(path, ".") && !strings.HasSuffix(path, "/") {
						return
					}
					
					// 对于其他所有路径，都返回index.html，让Vue Router处理
					if path != "/" && !strings.HasPrefix(path, "/api/") {
						r.Response.ServeFile("resource/public/html/index.html")
						r.ExitAll()
					}
				})
			})
			s.Run()
			return nil
		},
	}
)
