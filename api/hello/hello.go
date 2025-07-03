// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package hello

import (
	"context"

	v1 "github.com/ciclebyte/template_starter/api/hello/v1"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

type IHelloV1 interface {
	Hello(ctx context.Context, req *v1.HelloReq) (res *v1.HelloRes, err error)
}
