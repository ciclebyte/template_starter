package main

import (
	"github.com/ciclebyte/template_starter/internal/config"
	_ "github.com/ciclebyte/template_starter/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/ciclebyte/template_starter/internal/cmd"
)

func main() {
	// 初始化配置
	config.InitConfig()
	cmd.Main.Run(gctx.GetInitCtx())
}
