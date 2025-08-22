package middleware

import (
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/ciclebyte/template_starter/library/libJWT"
	"github.com/ciclebyte/template_starter/library/libResponse"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterMiddleware(New())
}

func New() *sMiddleware {
	return &sMiddleware{}
}

type sMiddleware struct{}

func (s *sMiddleware) MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

// RequireAuth 认证中间件 - 验证JWT令牌
func (s *sMiddleware) RequireAuth(r *ghttp.Request) {
	ctx := r.Context()
	
	// 获取Authorization头
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		libResponse.JsonExit(r, 401, "缺少认证令牌")
		return
	}

	// 提取Token
	token := libJWT.ExtractTokenFromHeader(authHeader)
	if token == "" {
		libResponse.JsonExit(r, 401, "认证令牌格式错误")
		return
	}

	// 验证Token
	claims, err := libJWT.GetManager().ValidateToken(token)
	if err != nil {
		g.Log().Warning(ctx, "token validation failed:", err)
		libResponse.JsonExit(r, 401, "认证令牌无效或已过期")
		return
	}

	// 添加调试日志
	g.Log().Debug(ctx, "JWT claims UserID:", claims.UserID, "Username:", claims.Username)

	// 将用户信息存储到请求上下文
	r.SetCtxVar("user_id", claims.UserID)
	r.SetCtxVar("username", claims.Username)
	
	// 检查用户状态
	userInfo, err := service.Auth().GetCurrentUser(r.Context())
	if err != nil {
		g.Log().Warning(ctx, "get user info failed:", err)
		libResponse.JsonExit(r, 401, "用户信息获取失败")
		return
	}

	if userInfo.Status != 1 {
		libResponse.JsonExit(r, 401, "用户账户已被禁用")
		return
	}

	// 存储用户完整信息到上下文
	r.SetCtxVar("user_info", userInfo)
	
	r.Middleware.Next()
}

// OptionalAuth 可选认证中间件 - 支持匿名访问
func (s *sMiddleware) OptionalAuth(r *ghttp.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		token := libJWT.ExtractTokenFromHeader(authHeader)
		if token != "" {
			claims, err := libJWT.GetManager().ValidateToken(token)
			if err == nil {
				// 设置用户ID到上下文
				r.SetCtxVar("user_id", claims.UserID)
				r.SetCtxVar("username", claims.Username)
				
				// 验证成功，设置用户信息
				userInfo, err := service.Auth().GetCurrentUser(r.Context())
				if err == nil && userInfo.Status == 1 {
					r.SetCtxVar("user_info", userInfo)
				}
			}
		}
	}
	
	r.Middleware.Next()
}

// RequirePermission 权限检查中间件
func (s *sMiddleware) RequirePermission(permission string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		ctx := r.Context()
		
		// 获取用户ID
		userIdVar := r.GetCtxVar("user_id")
		if userIdVar == nil {
			libResponse.JsonExit(r, 401, "请先登录")
			return
		}

		userId := gconv.Int64(userIdVar)
		if userId == 0 {
			libResponse.JsonExit(r, 401, "用户信息无效")
			return
		}

		// 检查权限
		hasPermission, err := service.Auth().HasPermission(ctx, userId, permission)
		if err != nil {
			g.Log().Error(ctx, "check permission failed:", err)
			libResponse.JsonExit(r, 500, "权限检查失败")
			return
		}

		if !hasPermission {
			g.Log().Warning(ctx, "permission denied", g.Map{
				"user_id":    userId,
				"permission": permission,
			})
			libResponse.JsonExit(r, 403, "权限不足")
			return
		}

		r.Middleware.Next()
	}
}

// RequireRole 角色检查中间件
func (s *sMiddleware) RequireRole(role string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		ctx := r.Context()
		
		// 获取用户ID
		userIdVar := r.GetCtxVar("user_id")
		if userIdVar == nil {
			libResponse.JsonExit(r, 401, "请先登录")
			return
		}

		userId := gconv.Int64(userIdVar)
		if userId == 0 {
			libResponse.JsonExit(r, 401, "用户信息无效")
			return
		}

		// 检查角色
		hasRole, err := service.Auth().HasRole(ctx, userId, role)
		if err != nil {
			g.Log().Error(ctx, "check role failed:", err)
			libResponse.JsonExit(r, 500, "角色检查失败")
			return
		}

		if !hasRole {
			g.Log().Warning(ctx, "role denied", g.Map{
				"user_id": userId,
				"role":    role,
			})
			libResponse.JsonExit(r, 403, "角色权限不足")
			return
		}

		r.Middleware.Next()
	}
}

// RequireTemplateOwnerOrPermission 模板所有者或权限检查中间件
func (s *sMiddleware) RequireTemplateOwnerOrPermission(permission string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		ctx := r.Context()
		
		// 获取用户ID
		userIdVar := r.GetCtxVar("user_id")
		if userIdVar == nil {
			libResponse.JsonExit(r, 401, "请先登录")
			return
		}

		userId := gconv.Int64(userIdVar)
		if userId == 0 {
			libResponse.JsonExit(r, 401, "用户信息无效")
			return
		}

		// 获取模板ID (从路径参数或请求体)
		templateId := r.Get("id").Int64()
		if templateId == 0 {
			templateId = r.Get("template_id").Int64()
		}

		if templateId > 0 {
			// 检查是否是模板所有者
			template, err := service.Templates().GetById(ctx, templateId)
			if err == nil && template != nil {
				// 如果模板有owner_id字段且当前用户是所有者，则允许访问
				// 注意：这里需要根据实际的模板结构调整
				// if template.OwnerId == userId {
				//     r.Middleware.Next()
				//     return
				// }
			}
		}

		// 否则检查权限
		hasPermission, err := service.Auth().HasPermission(ctx, userId, permission)
		if err != nil {
			g.Log().Error(ctx, "check permission failed:", err)
			libResponse.JsonExit(r, 500, "权限检查失败")
			return
		}

		if !hasPermission {
			g.Log().Warning(ctx, "permission denied", g.Map{
				"user_id":     userId,
				"template_id": templateId,
				"permission":  permission,
			})
			libResponse.JsonExit(r, 403, "权限不足")
			return
		}

		r.Middleware.Next()
	}
}
