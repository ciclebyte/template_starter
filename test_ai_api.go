package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	// 测试API调用
	testAPICall()
}

func testAPICall() {
	ctx := context.Background()
	
	// 构建请求体
	requestBody := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": "Hello, this is a test",
			},
		},
		"max_tokens":  100,
		"temperature": 0.7,
		"stream":      false,
	}
	
	// 创建HTTP客户端
	client := g.Client()
	client.SetTimeout(60 * 1000) // 60秒
	client.SetHeader("Authorization", "Bearer sk-YOUR-API-KEY-HERE") // 替换为真实API密钥
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("User-Agent", "GoFrame-Test/1.0")
	
	apiURL := "https://api.deepseek.com/v1/chat/completions"
	
	fmt.Printf("发送请求到: %s\n", apiURL)
	fmt.Printf("请求体: %+v\n", requestBody)
	
	start := time.Now()
	
	// 使用PostContent方法
	responseBody := client.PostContent(ctx, apiURL, requestBody)
	
	duration := time.Since(start)
	fmt.Printf("请求耗时: %v\n", duration)
	fmt.Printf("响应长度: %d\n", len(responseBody))
	fmt.Printf("响应内容: %s\n", responseBody)
	
	// 解析响应
	var apiResponse map[string]interface{}
	if err := json.Unmarshal([]byte(responseBody), &apiResponse); err != nil {
		fmt.Printf("解析响应失败: %v\n", err)
		return
	}
	
	fmt.Printf("解析成功: %+v\n", apiResponse)
}