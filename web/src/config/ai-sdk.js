// AI SDK 配置文件
// 配置Vercel AI SDK使用我们的OpenAI兼容端点

import { createOpenAI } from '@ai-sdk/openai'

// 创建自定义OpenAI客户端，指向我们的后端API
export const customOpenAI = createOpenAI({
  baseURL: '/api/v1/v1', // 使用我们的OpenAI兼容API端点
  apiKey: 'template-starter-dummy-key', // 占位符，我们的后端不验证这个
  // 可选：添加其他配置
  compatibility: {
    // 确保与我们的API格式兼容
    strict: false
  }
})

// AI SDK 配置选项
export const aiSDKConfig = {
  // 默认使用的端点
  apiEndpoint: '/api/v1/v1/chat/completions',
  
  // 默认模型配置
  defaultModel: 'gpt-3.5-turbo',
  
  // 请求配置
  requestConfig: {
    headers: {
      'Content-Type': 'application/json',
    },
    timeout: 300000, // 5分钟超时
  },
  
  // 聊天配置
  chatConfig: {
    temperature: 0.7,
    max_tokens: 2000,
    stream: true, // 默认启用流式响应
  },
  
  // 错误处理配置
  errorConfig: {
    maxRetries: 3,
    retryDelay: 1000, // 1秒
  },
  
  // UI配置
  uiConfig: {
    showTypingIndicator: true,
    showTokenCount: false,
    showModelInfo: true,
    enableMarkdown: true,
  },
  
  // 工具调用配置
  toolsConfig: {
    enabled: true,
    availableTools: [
      'code_explanation',
      'code_optimization', 
      'variable_suggestion',
      'template_generation'
    ]
  }
}

// 预定义的提示模板
export const promptTemplates = {
  explainCode: (code, language = 'text') => ({
    role: 'user',
    content: `请详细解释这段${language}代码的功能和逻辑：\n\n\`\`\`${language}\n${code}\n\`\`\``
  }),
  
  optimizeCode: (code, language = 'text') => ({
    role: 'user', 
    content: `请优化这段${language}代码的性能和可读性：\n\n\`\`\`${language}\n${code}\n\`\`\``
  }),
  
  suggestVariables: (content, projectType = 'general') => ({
    role: 'user',
    content: `基于以下${projectType}项目内容，请建议合适的模板变量：\n\n${content}`
  }),
  
  generateTemplate: (description, techStack = []) => ({
    role: 'user',
    content: `请生成${techStack.join(', ')}技术栈的代码模板：\n\n需求：${description}`
  })
}

// 系统提示
export const systemPrompts = {
  codeAssistant: `你是一个专业的代码助手，擅长解释代码、优化性能、建议最佳实践。请用中文回答，提供清晰、结构化的解释。`,
  
  templateHelper: `你是一个模板生成助手，能够分析项目需求并生成相应的代码模板和变量建议。请提供实用、可维护的模板方案。`,
  
  generalAssistant: `你是一个通用AI助手，能够帮助用户解决各种技术问题。请提供准确、有用的回答。`
}

// 导出默认配置
export default {
  customOpenAI,
  aiSDKConfig,
  promptTemplates,
  systemPrompts
}