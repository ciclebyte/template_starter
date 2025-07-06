package service

import (
	"context"

	indexApi "github.com/ciclebyte/template_starter/api/v1/index"
)

type IIndexService interface {
	GetIndexData(ctx context.Context, req *indexApi.IndexReq) (*indexApi.IndexRes, error)
}

var localIndexService IIndexService

func Index() IIndexService {
	if localIndexService == nil {
		panic("implement not found for interface IIndexService, forgot register?")
	}
	return localIndexService
}

func RegisterIndexService(i IIndexService) {
	localIndexService = i
}
