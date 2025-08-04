package service

import (
	"context"

	aiApi "github.com/ciclebyte/template_starter/api/v1/ai"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IAIService interface {
	TestConnection(ctx context.Context, req *aiApi.TestConnectionReq) (*aiApi.TestConnectionRes, error)
	GenerateTemplate(ctx context.Context, req *aiApi.GenerateTemplateReq) (*aiApi.GenerateTemplateRes, error)
	SuggestVariables(ctx context.Context, req *aiApi.SuggestVariablesReq) (*aiApi.SuggestVariablesRes, error)
	Chat(ctx context.Context, req *aiApi.ChatReq) (*aiApi.ChatRes, error)
	ChatStream(ctx context.Context, req *aiApi.ChatReq, r *ghttp.Request)
}

var localAIService IAIService

func AI() IAIService {
	if localAIService == nil {
		panic("implement not found for interface IAIService, forgot register?")
	}
	return localAIService
}

func RegisterAIService(i IAIService) {
	localAIService = i
}