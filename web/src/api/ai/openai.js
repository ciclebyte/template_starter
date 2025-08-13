import request from '@/utils/request'

// OpenAI兼容的聊天完成接口
export async function chatCompletions(data) {
  return request({
    url: '/api/v1/v1/chat/completions',
    method: 'post',
    data: {
      model: data.model || 'gpt-3.5-turbo',
      messages: data.messages || [],
      stream: data.stream || false,
      temperature: data.temperature || 0.7,
      max_tokens: data.max_tokens || 2000,
      functions: data.functions || [],
      function_call: data.function_call || 'auto',
      user: data.user,
      stop: data.stop,
      frequency_penalty: data.frequency_penalty || 0,
      presence_penalty: data.presence_penalty || 0
    },
    timeout: 300000 // 5分钟超时
  })
}

// OpenAI兼容的流式聊天完成接口
export async function chatCompletionsStream(data, onChunk, onComplete, onError) {
  const requestData = {
    model: data.model || 'gpt-3.5-turbo',
    messages: data.messages || [],
    stream: true,
    temperature: data.temperature || 0.7,
    max_tokens: data.max_tokens || 2000,
    functions: data.functions || [],
    function_call: data.function_call || 'auto',
    user: data.user,
    stop: data.stop,
    frequency_penalty: data.frequency_penalty || 0,
    presence_penalty: data.presence_penalty || 0
  }

  let fullContent = ''
  let functionCallData = null
  let metadata = null
  let controller = new AbortController()

  try {
    const response = await fetch('/api/v1/v1/chat/completions', {
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
            if (jsonStr.trim() === '' || jsonStr.trim() === '[DONE]') continue
            
            const data = JSON.parse(jsonStr)
            
            // 处理OpenAI格式的流式响应
            if (data.choices && data.choices[0]) {
              const choice = data.choices[0]
              
              if (choice.delta) {
                // 处理内容增量
                if (choice.delta.content) {
                  fullContent += choice.delta.content
                  if (onChunk) onChunk({ 
                    type: 'chunk', 
                    content: choice.delta.content, 
                    fullContent 
                  })
                }
                
                // 处理函数调用
                if (choice.delta.function_call) {
                  if (!functionCallData) {
                    functionCallData = {
                      name: choice.delta.function_call.name || '',
                      arguments: choice.delta.function_call.arguments || ''
                    }
                  } else {
                    if (choice.delta.function_call.arguments) {
                      functionCallData.arguments += choice.delta.function_call.arguments
                    }
                  }
                }
              }
              
              // 检查是否完成
              if (choice.finish_reason) {
                if (onComplete) {
                  // 解析函数调用参数
                  let suggestions = []
                  if (functionCallData && functionCallData.arguments) {
                    try {
                      const args = JSON.parse(functionCallData.arguments)
                      suggestions = args.suggestions || []
                    } catch (e) {
                      console.error('解析函数调用参数失败:', e)
                    }
                  }
                  
                  onComplete({
                    content: fullContent,
                    suggestions: suggestions,
                    metadata: metadata || {},
                    function_call: functionCallData
                  })
                }
                return
              }
            }
          } catch (parseError) {
            console.error('解析OpenAI流式响应失败:', parseError, 'Line:', line)
          }
        }
      }
    }
  } catch (error) {
    console.error('OpenAI流式连接错误:', error)
    if (onError) onError(error)
  }

  // 返回控制对象
  return {
    close: () => controller.abort()
  }
}

// === 便捷方法（基于OpenAI格式） ===

// 代码解释（使用OpenAI格式）
export function explainCodeWithOpenAI(data) {
  const messages = [
    {
      role: 'system',
      content: '你是一个专业的代码助手，擅长解释代码功能和逻辑。请用中文回答，并提供结构化的解释。'
    },
    {
      role: 'user', 
      content: `请解释这段代码的功能和逻辑：\n\n\`\`\`${data.language || 'text'}\n${data.code}\n\`\`\``
    }
  ]

  // 如果有选中的代码，优先解释选中部分
  if (data.hasSelection && data.selectedCode) {
    messages[1].content = `请解释这段选中的代码片段：\n\n\`\`\`${data.language || 'text'}\n${data.selectedCode}\n\`\`\`\n\n完整文件上下文：\n\`\`\`${data.language || 'text'}\n${data.fullCode}\n\`\`\``
  }

  return chatCompletions({
    model: 'gpt-3.5-turbo',
    messages: messages,
    temperature: 0.3,
    max_tokens: 2000,
    functions: [
      {
        name: 'suggest_actions',
        description: '为代码提供改进建议',
        parameters: {
          type: 'object',
          properties: {
            suggestions: {
              type: 'array',
              items: {
                type: 'object',
                properties: {
                  type: {type: 'string', enum: ['code', 'action', 'variable']},
                  name: {type: 'string'},
                  description: {type: 'string'},
                  code: {type: 'string'},
                  priority: {type: 'string', enum: ['high', 'medium', 'low']}
                }
              }
            }
          }
        }
      }
    ],
    function_call: 'auto'
  })
}

// 代码优化（使用OpenAI格式）
export function optimizeCodeWithOpenAI(data) {
  const messages = [
    {
      role: 'system',
      content: '你是一个专业的代码优化助手，能够分析代码并提供性能和可读性改进建议。'
    },
    {
      role: 'user',
      content: `请优化这段代码的性能和可读性：\n\n\`\`\`${data.language || 'text'}\n${data.code}\n\`\`\`\n\n${data.requirements || '请重点关注代码的可读性、性能和最佳实践。'}`
    }
  ]

  return chatCompletions({
    model: 'gpt-3.5-turbo',
    messages: messages,
    temperature: 0.3,
    max_tokens: 3000,
    functions: [
      {
        name: 'suggest_optimizations', 
        description: '提供代码优化建议',
        parameters: {
          type: 'object',
          properties: {
            suggestions: {
              type: 'array',
              items: {
                type: 'object',
                properties: {
                  type: {type: 'string', enum: ['code', 'refactor', 'performance']},
                  name: {type: 'string'},
                  description: {type: 'string'},
                  code: {type: 'string'},
                  priority: {type: 'string', enum: ['high', 'medium', 'low']}
                }
              }
            }
          }
        }
      }
    ],
    function_call: 'auto'
  })
}

// 变量建议（使用OpenAI格式）
export function suggestVariablesWithOpenAI(data) {
  const messages = [
    {
      role: 'system',
      content: '你是一个模板变量建议助手，能够分析代码内容并推荐合适的模板变量。'
    },
    {
      role: 'user',
      content: `基于以下内容，请建议合适的模板变量：\n\n项目类型：${data.projectType || '未指定'}\n技术栈：${data.techStack?.join(', ') || '未指定'}\n\n代码内容：\n\`\`\`\n${data.content}\n\`\`\``
    }
  ]

  return chatCompletions({
    model: 'gpt-3.5-turbo',
    messages: messages,
    temperature: 0.5,
    max_tokens: 2000,
    functions: [
      {
        name: 'suggest_variables',
        description: '推荐模板变量',
        parameters: {
          type: 'object',
          properties: {
            suggestions: {
              type: 'array',
              items: {
                type: 'object',
                properties: {
                  type: {type: 'string', enum: ['variable']},
                  name: {type: 'string'},
                  description: {type: 'string'},
                  variable_name: {type: 'string'},
                  variable_type: {type: 'string'},
                  default_value: {type: 'string'},
                  priority: {type: 'string', enum: ['high', 'medium', 'low']}
                }
              }
            }
          }
        }
      }
    ],
    function_call: 'auto'
  })
}

// 通用聊天（使用OpenAI格式）
export function generalChatWithOpenAI(data) {
  return chatCompletions({
    model: data.model || 'gpt-3.5-turbo',
    messages: data.messages,
    temperature: data.temperature || 0.7,
    max_tokens: data.max_tokens || 2000
  })
}

// 流式代码解释
export function explainCodeStreamWithOpenAI(data, onChunk, onComplete, onError) {
  const messages = [
    {
      role: 'system',
      content: '你是一个专业的代码助手，擅长解释代码功能和逻辑。请用中文回答，并提供结构化的解释。'
    },
    {
      role: 'user',
      content: `请解释这段代码的功能和逻辑：\n\n\`\`\`${data.language || 'text'}\n${data.code}\n\`\`\``
    }
  ]

  return chatCompletionsStream({
    model: 'gpt-3.5-turbo',
    messages: messages,
    temperature: 0.3,
    max_tokens: 2000,
    functions: [
      {
        name: 'suggest_actions',
        description: '为代码提供改进建议',
        parameters: {
          type: 'object',
          properties: {
            suggestions: {
              type: 'array',
              items: {
                type: 'object',
                properties: {
                  type: {type: 'string'},
                  name: {type: 'string'},
                  description: {type: 'string'},
                  code: {type: 'string'},
                  priority: {type: 'string'}
                }
              }
            }
          }
        }
      }
    ],
    function_call: 'auto'
  }, onChunk, onComplete, onError)
}