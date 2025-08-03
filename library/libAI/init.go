package libAI

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	initOnce sync.Once
	isInited bool
)

// InitAI 初始化AI框架
func InitAI(ctx context.Context) error {
	var initErr error
	
	initOnce.Do(func() {
		g.Log().Info(ctx, "开始初始化AI框架")
		
		// 这里可以添加AI相关的初始化逻辑
		// 比如预热模型、加载配置等
		
		g.Log().Info(ctx, "AI框架初始化完成")
		isInited = true
	})
	
	return initErr
}

// IsAIInited 检查AI是否已初始化
func IsAIInited() bool {
	return isInited
}

// GetAIClientWithInit 获取AI客户端并确保已初始化
func GetAIClientWithInit(ctx context.Context) (AIClient, error) {
	// 确保AI已初始化
	if err := InitAI(ctx); err != nil {
		return nil, err
	}
	
	return NewAIClient(ctx)
}

// GetMetrics 获取AI指标
func GetMetrics() map[string]interface{} {
	if !isInited {
		return map[string]interface{}{
			"status": "not_initialized",
			"error":  "AI framework not initialized",
		}
	}
	
	return map[string]interface{}{
		"status":     "initialized",
		"framework":  "simplified_eino",
		"components": []string{"simple_client", "template_generator", "variable_suggester"},
	}
}