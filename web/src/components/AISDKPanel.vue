<template>
  <div class="ai-sdk-panel">
    <!-- AI SDK悬浮按钮 -->
    <div class="ai-sdk-float-button" 
         :class="{ 'expanded': aiSDKPanelVisible }"
         @click="toggleAISDKPanel">
      <div class="ai-button-icon">
        <svg v-if="!aiSDKPanelVisible" class="ai-icon" viewBox="0 0 24 24" fill="currentColor">
          <path d="M9,4V6H16L17,7V16A1,1 0 0,1 16,17H3A1,1 0 0,1 2,16V4A1,1 0 0,1 3,3H8A1,1 0 0,1 9,4M4,5V15H15V8H12A1,1 0 0,1 11,7V5H4M13,5V6H15.5L13,8.5V5M7,9H13V11H7V9M7,12H11V14H7V12Z"/>
        </svg>
        <svg v-else class="close-icon" viewBox="0 0 24 24" fill="currentColor">
          <path d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z"/>
        </svg>
      </div>
      <div class="ai-button-pulse" v-show="isLoading"></div>
      <div class="sdk-badge">SDK</div>
    </div>

    <!-- AI SDK对话面板 -->
    <div class="ai-sdk-chat-panel" 
         :class="{ 'visible': aiSDKPanelVisible, 'resizing': isResizing }"
         :style="{ width: panelWidth + 'px', height: panelHeight + 'px' }"
         @click.stop>
      
      <!-- 拖拽调整手柄 -->
      <div class="resize-handle" 
           @mousedown="startResize"
           @touchstart="startResizeTouch"
           :class="{ 'active': isResizing }"
           title="拖拽调整面板大小">
        <div class="resize-dots">
          <span></span>
          <span></span>
          <span></span>
          <span></span>
        </div>
      </div>
      
      <!-- 面板头部 -->
      <div class="ai-panel-header">
        <div class="ai-status-indicator">
          <div class="status-dot" :class="{ 'online': aiConnected }"></div>
          <span class="status-text">AI助手 {{ aiConnected ? '在线' : '离线' }}</span>
          <div class="api-format-switch">
            <button 
              @click="toggleAPIFormat" 
              class="format-switch-btn openai"
              title="OpenAI兼容格式">
              OpenAI
            </button>
          </div>
        </div>
        <div class="ai-context-info">
          <span class="context-file" v-if="contextInfo.fileName">📄 {{ contextInfo.fileName }}</span>
          <span class="context-variables">🔢 {{ contextInfo.variableCount }}个变量</span>
          <span class="context-lines">📝 {{ contextInfo.codeLines }}行代码</span>
          <span class="context-selection" v-if="props.editorSelection.hasSelection">
            ✂️ 已选中 {{ props.editorSelection.selectionLength }} 个字符
          </span>
        </div>
      </div>

      <!-- 对话历史 -->
      <div class="ai-chat-history" ref="chatHistory">
        <div v-for="message in chatMessages" 
             :key="message.id" 
             :class="['chat-message', message.type]">
          
          <!-- 用户消息 -->
          <div v-if="message.type === 'user'" class="user-message">
            <div class="message-content">{{ message.content }}</div>
            <div class="message-time">{{ formatTime(message.timestamp) }}</div>
          </div>

          <!-- AI消息 -->
          <div v-else class="ai-message">
            <div class="ai-avatar">🤖</div>
            <div class="message-body">
              <!-- AI响应内容 -->
              <div class="message-content" v-html="formatAIMessage(message.content)"></div>
              
              <!-- 流式状态指示器 -->
              <div v-if="message.streaming" class="streaming-indicator">
                <div class="typing-dots">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
                <span class="streaming-text">AI正在思考...</span>
              </div>
              
              <div class="message-time" v-if="!message.streaming">{{ formatTime(message.timestamp) }}</div>
            </div>
          </div>
        </div>
        
        <!-- 空状态 -->
        <div v-if="chatMessages.length === 0" class="empty-state">
          <div class="empty-icon">💡</div>
          <div class="empty-text">开始与AI助手对话</div>
          <div class="empty-hint">选择代码片段或文件，然后点击下方快捷操作</div>
        </div>
      </div>

      <!-- 快捷操作区域 -->
      <div class="ai-quick-actions">
        <div class="quick-actions-header">
          <span class="actions-title">快捷操作</span>
          <span class="actions-hint">专注代码解释功能</span>
        </div>
        <div class="quick-actions-grid">
          <button 
            @click="explainCode" 
            :disabled="!currentCode || isLoading"
            class="quick-action-btn primary">
            <span class="action-icon">💡</span>
            <span class="action-text">解释代码</span>
          </button>
          
          <button 
            @click="explainInDetail" 
            :disabled="!currentCode || isLoading"
            class="quick-action-btn">
            <span class="action-icon">🔍</span>
            <span class="action-text">详细分析</span>
          </button>
          
          <button 
            @click="explainLogic" 
            :disabled="!currentCode || isLoading"
            class="quick-action-btn">
            <span class="action-icon">🧠</span>
            <span class="action-text">逻辑解析</span>
          </button>
          
          <button 
            @click="clearChat" 
            :disabled="chatMessages.length === 0"
            class="quick-action-btn danger">
            <span class="action-icon">🗑️</span>
            <span class="action-text">清空</span>
          </button>
        </div>
      </div>
      
      <!-- 输入区域 -->
      <div class="ai-input-area">
        <div class="input-wrapper">
          <textarea 
            v-model="userInput" 
            ref="userInputRef"
            placeholder="输入问题或直接使用上方快捷操作..."
            class="user-input"
            :disabled="isLoading"
            @keydown="handleInputKeydown"
            @input="adjustTextareaHeight"
            rows="1"></textarea>
          <button 
            @click="sendMessage" 
            :disabled="!userInput.trim() || isLoading"
            class="send-btn">
            <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
              <path d="M2,21L23,12L2,3V10L17,12L2,14V21Z"/>
            </svg>
          </button>
        </div>
        <div class="input-hint">
          <span class="hint-text">{{ isLoading ? 'AI正在回复中...' : 'Ctrl+Enter 发送消息' }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick, onMounted, watch } from 'vue'
import { aiSDKConfig } from '@/config/ai-sdk.js'

// Props
const props = defineProps({
  currentFileName: {
    type: String,
    default: ''
  },
  currentFileContent: {
    type: String,
    default: ''
  },
  templateVariables: {
    type: Array,
    default: () => []
  },
  editorSelection: {
    type: Object,
    default: () => ({
      hasSelection: false,
      selectedText: '',
      selectionStart: 0,
      selectionEnd: 0,
      selectionLength: 0
    })
  }
})

// Emits
const emit = defineEmits(['insertCode', 'addVariable', 'createFile', 'applySuggestion'])

// 响应式数据
const aiSDKPanelVisible = ref(false)
const chatHistory = ref(null)
const userInputRef = ref(null)
const isLoading = ref(false)
const userInput = ref('')
const chatMessages = ref([])
const messageIdCounter = ref(0)

// 面板调整相关
const isResizing = ref(false)
const panelWidth = ref(400)
const panelHeight = ref(600)
const resizeStartX = ref(0)
const resizeStartY = ref(0)
const resizeStartWidth = ref(0)
const resizeStartHeight = ref(0)

// AI连接状态
const aiConnected = ref(true)

// 计算属性
const contextInfo = computed(() => ({
  fileName: props.currentFileName || '未选择文件',
  variableCount: props.templateVariables.length,
  codeLines: props.currentFileContent ? props.currentFileContent.split('\n').length : 0,
  hasSelection: props.editorSelection.hasSelection
}))

const currentCode = computed(() => {
  if (props.editorSelection.hasSelection) {
    return props.editorSelection.selectedText
  }
  return props.currentFileContent
})

const currentLanguage = computed(() => {
  return getFileLanguage(props.currentFileName)
})

// 方法
const toggleAISDKPanel = () => {
  aiSDKPanelVisible.value = !aiSDKPanelVisible.value
  if (aiSDKPanelVisible.value) {
    nextTick(() => {
      scrollToBottom()
      if (userInputRef.value) {
        userInputRef.value.focus()
      }
    })
  }
}

const explainCode = async () => {
  if (!currentCode.value || isLoading.value) return
  
  const prompt = `请详细解释这段${currentLanguage.value}代码的功能和逻辑：\n\n\`\`\`${currentLanguage.value}\n${currentCode.value}\n\`\`\``
  
  await sendAIMessage(prompt)
}

const explainInDetail = async () => {
  if (!currentCode.value || isLoading.value) return
  
  const prompt = `请深入分析这段${currentLanguage.value}代码，包括：\n1. 代码功能和目的\n2. 实现原理和算法\n3. 关键变量和数据结构\n4. 可能的优化点\n\n代码：\n\`\`\`${currentLanguage.value}\n${currentCode.value}\n\`\`\``
  
  await sendAIMessage(prompt)
}

const explainLogic = async () => {
  if (!currentCode.value || isLoading.value) return
  
  const prompt = `请分析这段${currentLanguage.value}代码的逻辑流程：\n1. 执行步骤和流程图\n2. 条件判断和分支逻辑\n3. 循环和递归逻辑\n4. 异常处理机制\n\n代码：\n\`\`\`${currentLanguage.value}\n${currentCode.value}\n\`\`\``
  
  await sendAIMessage(prompt)
}

const sendMessage = async () => {
  if (!userInput.value.trim() || isLoading.value) return
  
  await sendAIMessage(userInput.value.trim())
  userInput.value = ''
  adjustTextareaHeight()
}

const sendAIMessage = async (message) => {
  // 添加用户消息
  const userMessage = {
    id: `user-${messageIdCounter.value++}`,
    type: 'user',
    content: message,
    timestamp: new Date()
  }
  chatMessages.value.push(userMessage)
  
  // 添加AI消息占位符
  const aiMessage = {
    id: `ai-${messageIdCounter.value++}`,
    type: 'ai',
    content: '',
    streaming: true,
    timestamp: new Date()
  }
  chatMessages.value.push(aiMessage)
  
  scrollToBottom()
  isLoading.value = true
  
  try {
    const response = await fetch(aiSDKConfig.apiEndpoint, {
      method: 'POST',
      headers: aiSDKConfig.requestConfig.headers,
      body: JSON.stringify({
        model: aiSDKConfig.defaultModel,
        messages: [
          {
            role: 'system',
            content: '你是一个专业的代码分析师，擅长解释代码功能、分析逻辑结构。请用中文提供清晰、结构化的解释。'
          },
          {
            role: 'user',
            content: message
          }
        ],
        stream: true,
        temperature: 0.3,
        max_tokens: 2000
      })
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let fullContent = ''

    while (true) {
      const { done, value } = await reader.read()
      
      if (done) break

      const chunk = decoder.decode(value)
      const lines = chunk.split('\n')

      for (const line of lines) {
        if (line.trim().startsWith('data: ')) {
          const jsonStr = line.slice(6).trim()
          if (jsonStr === '' || jsonStr === '[DONE]') {
            console.log('AI SDK Panel - 收到结束信号')
            continue
          }
          
          try {
            const data = JSON.parse(jsonStr)
            console.log('AI SDK Panel - 解析的数据:', data)
            
            // 标准OpenAI流式响应格式
            if (data.choices && data.choices[0] && data.choices[0].delta) {
              const deltaContent = data.choices[0].delta.content
              if (deltaContent) {
                fullContent += deltaContent
                aiMessage.content = fullContent
                
                // 强制触发Vue的响应式更新 - 直接修改原数组
                const lastIndex = chatMessages.value.length - 1
                if (lastIndex >= 0) {
                  chatMessages.value[lastIndex] = { ...aiMessage }
                }
                
                // 立即滚动到底部
                scrollToBottom()
                
                console.log('AI SDK Panel - 收到chunk:', deltaContent, '总长度:', fullContent.length)
              }
              
              // 检查是否完成
              if (data.choices[0].finish_reason === 'stop') {
                aiMessage.streaming = false
                aiMessage.timestamp = new Date()
                chatMessages.value[chatMessages.value.length - 1] = { ...aiMessage }
                console.log('AI SDK Panel - 流式响应完成')
                break
              }
            }
          } catch (parseError) {
            console.error('AI SDK Panel - JSON解析失败:', parseError, '原始数据:', jsonStr)
          }
        }
      }
    }
  } catch (error) {
    console.error('AI SDK Panel - 请求失败:', error)
    aiMessage.content = '抱歉，AI服务暂时不可用。请稍后再试。'
    aiMessage.streaming = false
  } finally {
    isLoading.value = false
    scrollToBottom()
  }
}

const clearChat = () => {
  chatMessages.value = []
}

const toggleAPIFormat = () => {
  // 兼容原有面板的API格式切换功能，但SDK面板固定使用OpenAI格式
  console.log('AI SDK Panel - 固定使用OpenAI标准格式')
}

const getFileLanguage = (fileName) => {
  if (!fileName) return 'text'
  
  const ext = fileName.split('.').pop()?.toLowerCase()
  const langMap = {
    'js': 'javascript',
    'ts': 'typescript',
    'vue': 'vue',
    'jsx': 'javascript',
    'tsx': 'typescript',
    'py': 'python',
    'go': 'go',
    'java': 'java',
    'php': 'php',
    'cpp': 'cpp',
    'c': 'c',
    'css': 'css',
    'html': 'html',
    'json': 'json',
    'yml': 'yaml',
    'yaml': 'yaml'
  }
  
  return langMap[ext] || 'text'
}

const scrollToBottom = () => {
  nextTick(() => {
    if (chatHistory.value) {
      chatHistory.value.scrollTop = chatHistory.value.scrollHeight
    }
  })
}

const formatTime = (timestamp) => {
  if (!timestamp) return ''
  return new Date(timestamp).toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const formatAIMessage = (content) => {
  return content
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>') // 粗体
    .replace(/\*(.*?)\*/g, '<em>$1</em>') // 斜体
    .replace(/`(.*?)`/g, '<code>$1</code>') // 行内代码
    .replace(/```(.*?)```/gs, '<pre><code>$1</code></pre>') // 代码块
    .replace(/\n\n/g, '<br><br>') // 段落
    .replace(/\n/g, '<br>') // 换行
}

const handleInputKeydown = (event) => {
  if ((event.ctrlKey || event.metaKey) && event.key === 'Enter') {
    event.preventDefault()
    sendMessage()
  }
}

const adjustTextareaHeight = () => {
  nextTick(() => {
    if (userInputRef.value) {
      userInputRef.value.style.height = 'auto'
      userInputRef.value.style.height = Math.min(userInputRef.value.scrollHeight, 120) + 'px'
    }
  })
}

// 面板调整功能
const startResize = (event) => {
  isResizing.value = true
  resizeStartX.value = event.clientX
  resizeStartY.value = event.clientY
  resizeStartWidth.value = panelWidth.value
  resizeStartHeight.value = panelHeight.value
  
  document.addEventListener('mousemove', handleResize)
  document.addEventListener('mouseup', stopResize)
  event.preventDefault()
}

const startResizeTouch = (event) => {
  const touch = event.touches[0]
  isResizing.value = true
  resizeStartX.value = touch.clientX
  resizeStartY.value = touch.clientY
  resizeStartWidth.value = panelWidth.value
  resizeStartHeight.value = panelHeight.value
  
  document.addEventListener('touchmove', handleResizeTouch)
  document.addEventListener('touchend', stopResizeTouch)
  event.preventDefault()
}

const handleResize = (event) => {
  if (!isResizing.value) return
  
  const deltaX = resizeStartX.value - event.clientX
  const deltaY = event.clientY - resizeStartY.value
  
  panelWidth.value = Math.max(320, Math.min(800, resizeStartWidth.value + deltaX))
  panelHeight.value = Math.max(400, Math.min(800, resizeStartHeight.value + deltaY))
}

const handleResizeTouch = (event) => {
  if (!isResizing.value) return
  
  const touch = event.touches[0]
  const deltaX = resizeStartX.value - touch.clientX
  const deltaY = touch.clientY - resizeStartY.value
  
  panelWidth.value = Math.max(320, Math.min(800, resizeStartWidth.value + deltaX))
  panelHeight.value = Math.max(400, Math.min(800, resizeStartHeight.value + deltaY))
}

const stopResize = () => {
  isResizing.value = false
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
  
  // 保存面板尺寸到localStorage
  localStorage.setItem('aiSDKPanelWidth', panelWidth.value.toString())
  localStorage.setItem('aiSDKPanelHeight', panelHeight.value.toString())
}

const stopResizeTouch = () => {
  isResizing.value = false
  document.removeEventListener('touchmove', handleResizeTouch)
  document.removeEventListener('touchend', stopResizeTouch)
  
  // 保存面板尺寸到localStorage
  localStorage.setItem('aiSDKPanelWidth', panelWidth.value.toString())
  localStorage.setItem('aiSDKPanelHeight', panelHeight.value.toString())
}

// 生命周期
onMounted(() => {
  console.log('AI SDK Panel mounted - OpenAI Standard Format')
  console.log('Using API endpoint:', aiSDKConfig.apiEndpoint)
  
  // 恢复面板尺寸
  const savedWidth = localStorage.getItem('aiSDKPanelWidth')
  const savedHeight = localStorage.getItem('aiSDKPanelHeight')
  
  if (savedWidth) panelWidth.value = parseInt(savedWidth)
  if (savedHeight) panelHeight.value = parseInt(savedHeight)
})
</script>

<style scoped>
.ai-sdk-panel {
  position: fixed;
  top: 25%;
  right: 16px;
  transform: translateY(-50%);
  z-index: 1001;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
}

/* 参考原有AI面板的悬浮按钮样式 */
.ai-sdk-float-button {
  position: relative;
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.3);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: visible;
}

.ai-sdk-float-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(102, 126, 234, 0.4);
}

.ai-sdk-float-button.expanded {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
}

.ai-sdk-float-button.expanded:hover {
  background: linear-gradient(135deg, #ff5252 0%, #e84118 100%);
}

.ai-button-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: white;
}

.ai-button-icon svg {
  width: 24px;
  height: 24px;
  transition: all 0.2s;
}

.ai-button-pulse {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.6);
  animation: pulse 2s infinite;
}

.sdk-badge {
  position: absolute;
  top: -8px;
  right: -8px;
  background: #f59e0b;
  color: white;
  font-size: 10px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 10px;
  border: 2px solid white;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* 参考原有AI面板的对话面板样式 */
.ai-sdk-chat-panel {
  position: absolute;
  right: 72px;
  top: 50%;
  transform: translateY(-50%);
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 20px 64px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(0, 0, 0, 0.06);
  opacity: 0;
  visibility: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.ai-sdk-chat-panel.visible {
  opacity: 1;
  visibility: visible;
}

.ai-sdk-chat-panel.resizing {
  transition: none;
}

/* 调整手柄样式 */
.resize-handle {
  position: absolute;
  top: 8px;
  left: 8px;
  width: 24px;
  height: 24px;
  cursor: nw-resize;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  opacity: 0.5;
  transition: all 0.2s;
}

.resize-handle:hover,
.resize-handle.active {
  opacity: 1;
  background: rgba(102, 126, 234, 0.1);
}

.resize-dots {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2px;
  width: 8px;
  height: 8px;
}

.resize-dots span {
  width: 3px;
  height: 3px;
  background: #94a3b8;
  border-radius: 50%;
  transition: background 0.2s;
}

.resize-handle:hover .resize-dots span {
  background: #667eea;
}

/* 面板头部样式 */
.ai-panel-header {
  padding: 16px;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  border-radius: 16px 16px 0 0;
}

.ai-status-indicator {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #ef4444;
  margin-right: 8px;
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.2);
  transition: all 0.3s;
}

.status-dot.online {
  background: #10b981;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
}

.status-text {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

.api-format-switch {
  display: flex;
  align-items: center;
}

.format-switch-btn {
  padding: 4px 8px;
  font-size: 11px;
  font-weight: 600;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: #ffffff;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
}

.format-switch-btn.openai {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-color: transparent;
}

.ai-context-info {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  font-size: 12px;
  color: #6b7280;
}

.context-file, .context-variables, .context-lines, .context-selection {
  background: rgba(0, 0, 0, 0.05);
  padding: 2px 6px;
  border-radius: 4px;
}

/* 对话历史样式 */
.ai-chat-history {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.chat-message {
  display: flex;
  flex-direction: column;
}

.user-message {
  align-self: flex-end;
  max-width: 80%;
}

.user-message .message-content {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 12px 16px;
  border-radius: 16px 16px 4px 16px;
  font-size: 14px;
  line-height: 1.5;
  word-wrap: break-word;
}

.user-message .message-time {
  font-size: 11px;
  color: #9ca3af;
  text-align: right;
  margin-top: 4px;
}

.ai-message {
  align-self: flex-start;
  max-width: 90%;
  display: flex;
  gap: 12px;
}

.ai-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
}

.message-body {
  flex: 1;
  min-width: 0;
}

.ai-message .message-content {
  background: #f8fafc;
  color: #374151;
  padding: 12px 16px;
  border-radius: 4px 16px 16px 16px;
  font-size: 14px;
  line-height: 1.6;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.streaming-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  padding: 8px 12px;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(102, 126, 234, 0.1);
}

.typing-dots {
  display: flex;
  gap: 3px;
}

.typing-dots span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #667eea;
  animation: typing 1.5s infinite;
}

.typing-dots span:nth-child(2) {
  animation-delay: 0.3s;
}

.typing-dots span:nth-child(3) {
  animation-delay: 0.6s;
}

.streaming-text {
  font-size: 13px;
  color: #667eea;
}

.ai-message .message-time {
  font-size: 11px;
  color: #9ca3af;
  margin-top: 4px;
}

/* 空状态样式 */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  color: #9ca3af;
  padding: 40px 20px;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #6b7280;
}

.empty-hint {
  font-size: 14px;
  line-height: 1.5;
}

/* 快捷操作区域样式 */
.ai-quick-actions {
  padding: 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
  background: #fafbfc;
}

.quick-actions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.actions-title {
  font-size: 13px;
  font-weight: 600;
  color: #374151;
}

.actions-hint {
  font-size: 11px;
  color: #9ca3af;
}

.quick-actions-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.quick-action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 12px;
  font-size: 13px;
  font-weight: 500;
  border: 1.5px solid #e5e7eb;
  border-radius: 8px;
  background: white;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.quick-action-btn:hover:not(:disabled) {
  border-color: #667eea;
  background: #f8fafc;
  transform: translateY(-1px);
}

.quick-action-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-color: transparent;
  grid-column: 1 / -1;
}

.quick-action-btn.primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.quick-action-btn.danger {
  border-color: #fecaca;
  color: #ef4444;
}

.quick-action-btn.danger:hover:not(:disabled) {
  border-color: #ef4444;
  background: #fef2f2;
}

.quick-action-btn:disabled {
  background: #f9fafb;
  color: #d1d5db;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
  border-color: #f3f4f6;
}

.action-icon {
  font-size: 14px;
}

.action-text {
  font-size: 12px;
}

/* 输入区域样式 */
.ai-input-area {
  padding: 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
  background: white;
  border-radius: 0 0 16px 16px;
}

.input-wrapper {
  display: flex;
  gap: 8px;
  align-items: flex-end;
}

.user-input {
  flex: 1;
  min-height: 36px;
  max-height: 120px;
  padding: 8px 12px;
  border: 1.5px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.5;
  resize: none;
  background: #fafbfc;
  transition: all 0.2s;
}

.user-input:focus {
  outline: none;
  border-color: #667eea;
  background: white;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.user-input:disabled {
  background: #f9fafb;
  color: #d1d5db;
  cursor: not-allowed;
}

.send-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  flex-shrink: 0;
}

.send-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.send-btn:disabled {
  background: #e5e7eb;
  color: #9ca3af;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.input-hint {
  margin-top: 8px;
  display: flex;
  justify-content: center;
}

.hint-text {
  font-size: 11px;
  color: #9ca3af;
}

/* 动画 */
@keyframes pulse {
  0% {
    transform: translate(-50%, -50%) scale(1);
    opacity: 1;
  }
  100% {
    transform: translate(-50%, -50%) scale(1.3);
    opacity: 0;
  }
}

@keyframes typing {
  0%, 20% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-8px);
  }
  80%, 100% {
    transform: translateY(0);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .ai-sdk-chat-panel {
    right: 8px;
    width: calc(100vw - 88px) !important;
    max-width: 450px;
    height: 500px !important;
  }
  
  .quick-actions-grid {
    grid-template-columns: 1fr;
  }
  
  .quick-action-btn.primary {
    grid-column: 1;
  }
}

/* 滚动条样式 */
.ai-chat-history::-webkit-scrollbar {
  width: 4px;
}

.ai-chat-history::-webkit-scrollbar-track {
  background: transparent;
}

.ai-chat-history::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 2px;
}

.ai-chat-history::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
</style>