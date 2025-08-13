import request from '@/utils/request'

// AI连接测试
export function testAIConnection() {
  return request({
    url: '/api/v1/ai/testConnection',
    method: 'post'
  })
}

// 统一AI聊天接口（推荐使用）
export function chatWithAI(data) {
  return request({
    url: '/api/v1/ai/chat',
    method: 'post',
    data: {
      action: data.action || 'general_chat',
      context: data.context || {},
      userInput: data.userInput || data.message || '',
      preferences: data.preferences || {},
      chatHistory: data.chatHistory || []
    },
    timeout: 300000 // 增加到5分钟，支持长文本处理
  })
}

// 流式AI聊天接口（推荐用于长对话）
export async function chatWithAIStream(data, onChunk, onComplete, onError) {
  const requestData = {
    action: data.action || 'general_chat',
    context: data.context || {},
    userInput: data.userInput || data.message || '',
    preferences: data.preferences || {},
    chatHistory: data.chatHistory || [],
    stream: true // 设置为流式响应
  }

  let fullContent = ''
  let metadata = null
  let suggestions = []
  let controller = new AbortController()

  try {
    const response = await fetch('/api/v1/ai/chat/stream', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(requestData),
      signal: controller.signal
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const reader = response.body.getReader()
    const decoder = new TextDecoder()

    while (true) {
      const { done, value } = await reader.read()
      
      if (done) break

      const chunk = decoder.decode(value)
      const lines = chunk.split('\n')

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          try {
            const jsonStr = line.slice(6) // 移除 'data: ' 前缀
            if (jsonStr.trim() === '') continue
            
            const data = JSON.parse(jsonStr)
            
            switch(data.type) {
              case 'start':
                if (onChunk) onChunk({ type: 'start', content: data.content })
                break
              case 'chunk':
                fullContent += data.content
                if (onChunk) onChunk({ type: 'chunk', content: data.content, fullContent })
                break
              case 'suggestions':
                suggestions = data.suggestions
                break
              case 'metadata':
                metadata = data.metadata
                break
              case 'done':
                if (onComplete) {
                  // 尝试解析AI返回的JSON内容
                  let parsedContent = fullContent
                  let parsedSuggestions = suggestions
                  
                  try {
                    // 检查内容是否是JSON格式
                    if (fullContent.trim().startsWith('{') && fullContent.trim().endsWith('}')) {
                      const jsonResponse = JSON.parse(fullContent)
                      if (jsonResponse.content) {
                        parsedContent = jsonResponse.content
                      }
                      if (jsonResponse.suggestions && Array.isArray(jsonResponse.suggestions)) {
                        parsedSuggestions = jsonResponse.suggestions
                      }
                    }
                  } catch (jsonParseError) {
                    console.debug('内容不是JSON格式，使用原始内容:', jsonParseError)
                  }
                  
                  onComplete({
                    content: parsedContent,
                    suggestions: parsedSuggestions,
                    metadata: metadata
                  })
                }
                return
              case 'error':
                if (onError) onError(new Error(data.content))
                return
            }
          } catch (parseError) {
            console.error('解析流式响应失败:', parseError, 'Line:', line)
          }
        }
      }
    }
  } catch (error) {
    console.error('流式连接错误:', error)
    if (onError) onError(error)
  }

  // 返回控制对象
  return {
    close: () => controller.abort()
  }
}

// === 便捷方法（基于统一接口） ===

// AI代码优化
export function optimizeCode(data) {
  return chatWithAI({
    action: 'optimize_code',
    context: {
      fileName: data.fileName,
      fileContent: data.code,
      selectedText: data.selectionInfo?.selectedText,
      hasSelection: data.hasSelection || false,
      techStack: inferTechStack(data.fileName, data.code)
    },
    userInput: data.requirements || '请优化这段代码的性能和可读性',
    preferences: {
      codeStyle: 'modern',
      verbosity: 'detailed'
    }
  })
}

// AI代码解释
export function explainCode(data) {
  return chatWithAI({
    action: 'explain_code',
    context: {
      fileName: data.fileName,
      fileContent: data.fullFileContent || data.code,
      selectedText: data.hasSelection ? data.code : null,
      hasSelection: data.hasSelection || false,
      techStack: inferTechStack(data.fileName, data.code)
    },
    userInput: data.question || '请解释这段代码的功能和逻辑',
    preferences: {
      verbosity: 'detailed',
      includeExamples: true
    }
  })
}

// AI模板生成
export function generateTemplate(data) {
  return chatWithAI({
    action: 'generate_template',
    context: {
      projectType: data.projectType,
      techStack: data.techStack,
      variables: data.variables,
      features: data.features
    },
    userInput: data.description,
    preferences: {
      templateStyle: 'comprehensive',
      includeComments: true
    }
  })
}

// AI变量建议
export function suggestVariables(data) {
  return chatWithAI({
    action: 'suggest_variables',
    context: {
      projectType: data.projectType,
      techStack: data.techStack,
      existingVariables: data.existingVariables,
      templateContent: data.description
    },
    userInput: data.description,
    preferences: {
      variableNaming: 'camelCase',
      includeTypes: true
    }
  })
}

// 推断技术栈（辅助函数）
function inferTechStack(fileName, content) {
  const techStack = []
  
  if (!fileName || !content) return techStack
  
  const ext = fileName.split('.').pop()?.toLowerCase()
  const contentLower = content.toLowerCase()
  
  // 根据文件扩展名
  const extMap = {
    'js': 'javascript',
    'ts': 'typescript', 
    'vue': 'vue',
    'jsx': 'react',
    'tsx': 'react',
    'py': 'python',
    'go': 'golang',
    'java': 'java',
    'php': 'php'
  }
  
  if (extMap[ext]) {
    techStack.push(extMap[ext])
  }
  
  // 根据内容检测框架
  if (contentLower.includes('vue')) techStack.push('vue')
  if (contentLower.includes('react')) techStack.push('react')
  if (contentLower.includes('angular')) techStack.push('angular')
  if (contentLower.includes('express')) techStack.push('express')
  if (contentLower.includes('fastify')) techStack.push('fastify')
  
  return [...new Set(techStack)] // 去重
}