package controller

import (
	"context"

	aiApi "github.com/ciclebyte/template_starter/api/v1/ai"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/frame/g"
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
	g.Log().Debug(ctx, "AI.Chat called with action:", req.Action)
	return service.AI().Chat(ctx, req)
}