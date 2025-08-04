package controller

import (
	"context"

	aiApi "github.com/ciclebyte/template_starter/api/v1/ai"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var AI = cAI{}

type cAI struct {
	BaseController
}

// TestConnection 测试AI连接
func (c *cAI) TestConnection(ctx context.Context, req *aiApi.TestConnectionReq) (res *aiApi.TestConnectionRes, err error) {
	g.Log().Debug(ctx, "AI.TestConnection called")
	return service.AI().TestConnection(ctx, req)
}

// GenerateTemplate 生成模板
func (c *cAI) GenerateTemplate(ctx context.Context, req *aiApi.GenerateTemplateReq) (res *aiApi.GenerateTemplateRes, err error) {
	g.Log().Debug(ctx, "AI.GenerateTemplate called")
	return service.AI().GenerateTemplate(ctx, req)
}

// SuggestVariables 建议变量
func (c *cAI) SuggestVariables(ctx context.Context, req *aiApi.SuggestVariablesReq) (res *aiApi.SuggestVariablesRes, err error) {
	g.Log().Debug(ctx, "AI.SuggestVariables called")
	return service.AI().SuggestVariables(ctx, req)
}

// Chat 统一AI聊天接口
func (c *cAI) Chat(ctx context.Context, req *aiApi.ChatReq) (res *aiApi.ChatRes, err error) {
	g.Log().Debug(ctx, "AI.Chat called with action:", req.Action, "stream:", req.Stream)
	
	// 如果请求流式响应，返回错误提示使用专门的流式端点
	if req.Stream {
		return nil, gerror.New("流式响应请使用 /api/v1/ai/chat/stream 端点")
	}
	
	return service.AI().Chat(ctx, req)
}

// ChatStream 流式AI聊天接口 - 专门处理流式响应
func (c *cAI) ChatStream(r *ghttp.Request) {
	ctx := r.Context()
	g.Log().Debug(ctx, "AI.ChatStream called")
	
	// 解析请求体
	var req *aiApi.ChatReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteStatus(400, "请求参数解析失败: "+err.Error())
		return
	}
	
	// 强制设置为流式响应
	req.Stream = true
	
	g.Log().Debug(ctx, "AI.ChatStream parsed request with action:", req.Action)
	service.AI().ChatStream(ctx, req, r)
}
