package config

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	MINIO_ENDPOINT   = ""
	MINIO_ACCESS_KEY = ""
	MINIO_SECRET_KEY = ""
	MINIO_BUCKET     = ""
	MINIO_USE_SSL    = false
	MINIO_BASE_PATH  = ""
)

// InitConfig 初始化配置
func InitConfig() {
	ctx := context.Background()
	// 从配置文件读取MinIO配置
	MINIO_ENDPOINT = g.Cfg().MustGet(ctx, "minio.endpoint").String()
	MINIO_ACCESS_KEY = g.Cfg().MustGet(ctx, "minio.accessKey").String()
	MINIO_SECRET_KEY = g.Cfg().MustGet(ctx, "minio.secretKey").String()
	MINIO_BUCKET = g.Cfg().MustGet(ctx, "minio.bucket").String()
	MINIO_USE_SSL = g.Cfg().MustGet(ctx, "minio.useSSL").Bool()
	MINIO_BASE_PATH = g.Cfg().MustGet(ctx, "minio.basePath").String()
}
