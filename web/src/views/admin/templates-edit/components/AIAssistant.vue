<template>
  <div class="ai-assistant">
    <!-- AI悬浮按钮 -->
    <div class="ai-float-button" 
         :class="{ 'expanded': aiPanelVisible }"
         @click="toggleAIPanel">
      <div class="ai-button-icon">
        <svg v-if="!aiPanelVisible" class="ai-icon" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12,2A2,2 0 0,1 14,4C14,5.1 13.1,6 12,6C10.9,6 10,5.1 10,4A2,2 0 0,1 12,2M21,9V7L15,1H5C3.89,1 3,1.89 3,3V21A2,2 0 0,0 5,23H11V21H5V3H14V8H19V9H21M17.5,12A2.5,2.5 0 0,1 20,14.5A2.5,2.5 0 0,1 17.5,17A2.5,2.5 0 0,1 15,14.5A2.5,2.5 0 0,1 17.5,12M17.5,22C15,22 13,20 13,17.5C13,15 15,13 17.5,13C20,13 22,15 22,17.5C22,20 20,22 17.5,22Z"/>
        </svg>
        <svg v-else class="close-icon" viewBox="0 0 24 24" fill="currentColor">
          <path d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z"/>
        </svg>
      </div>
      <div class="ai-button-pulse" v-show="aiProcessing"></div>
    </div>

    <!-- AI对话面板 -->
    <div class="ai-chat-panel" 
         :class="{ 'visible': aiPanelVisible, 'resizing': isResizing }"
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
              class="format-switch-btn"
              :class="{ 'openai': useOpenAIFormat }"
              :title="useOpenAIFormat ? '当前使用OpenAI兼容格式' : '当前使用原始API格式'">
              {{ useOpenAIFormat ? 'OpenAI' : '原始' }}
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
              <!-- AI响应内容 - 始终显示 -->
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
              
              <!-- JSON响应完成后的操作区域 -->
              <div v-if="message.suggestions && message.suggestions.length > 0 && !message.streaming" class="ai-suggestions">
                <!-- 建议区域折叠控制 -->
                <div class="suggestions-header" @click="toggleSuggestionsCollapse(message)">
                  <div class="suggestions-title">
                    <span class="suggestions-icon">💡</span>
                    <span class="suggestions-text">AI建议</span>
                    <span class="suggestions-count">({{ message.suggestions.length }})</span>
                  </div>
                  <button class="collapse-btn" :class="{ collapsed: message.suggestionsCollapsed }">
                    <span class="collapse-icon">{{ message.suggestionsCollapsed ? '▶' : '▼' }}</span>
                    <span class="collapse-text">{{ message.suggestionsCollapsed ? '展开' : '折叠' }}</span>
                  </button>
                </div>
                
                <!-- 建议内容 -->
                <div class="suggestions-content">
                  <!-- 代码建议 -->
                  <div v-for="suggestion in getCodeSuggestions(message.suggestions)" 
                       :key="suggestion.name" 
                       class="suggestion-item code-suggestion">
                    <div class="suggestion-header">
                      <span class="suggestion-name">{{ suggestion.name }}</span>
                      <div class="suggestion-actions">
                        <button @click="viewSuggestionDetail(suggestion)" class="action-btn secondary">查看详情</button>
                        <button @click="insertCode(suggestion.code)" class="action-btn primary">插入代码</button>
                      </div>
                    </div>
                    <div class="suggestion-description">{{ suggestion.description }}</div>
                  </div>
                  
                  <!-- 变量建议 -->
                  <div v-if="getVariableSuggestions(message.suggestions).length > 0" class="variables-table">
                    <div class="variables-header">
                      <h4 class="variables-title">
                        <span class="variables-icon">🔧</span>
                        模板变量建议
                        <span class="variables-count">({{ getVariableSuggestions(message.suggestions).length }})</span>
                      </h4>
                      <div class="variables-actions">
                        <button @click="addAllVariables(getVariableSuggestions(message.suggestions))" class="action-btn mini primary">
                          全部添加
                        </button>
                      </div>
                    </div>
                    <div class="variables-container">
                      <table class="variables-grid">
                        <thead>
                          <tr>
                            <th class="var-name-col">变量名</th>
                            <th class="var-desc-col">描述</th>
                            <th class="var-code-col">模板代码</th>
                            <th class="var-priority-col">优先级</th>
                            <th class="var-actions-col">操作</th>
                          </tr>
                        </thead>
                        <tbody>
                          <tr v-for="variable in getSortedVariables(getVariableSuggestions(message.suggestions))" 
                              :key="variable.name" 
                              class="variable-row"
                              :class="{ 'high-priority': variable.priority === 'high' }">
                            <td class="variable-name">
                              <span class="var-name-text">{{ variable.name }}</span>
                            </td>
                            <td class="variable-desc">
                              <span class="var-desc-text">{{ variable.description }}</span>
                            </td>
                            <td class="variable-code">
                              <code class="var-code-text">{{ variable.code }}</code>
                            </td>
                            <td class="variable-priority">
                              <span class="priority-badge" :class="variable.priority">
                                {{ getPriorityText(variable.priority) }}
                              </span>
                            </td>
                            <td class="variable-actions">
                              <div class="action-group">
                                <button @click="addVariable(variable)" 
                                        class="action-btn mini primary"
                                        :disabled="isVariableAdded(variable.name)">
                                  {{ isVariableAdded(variable.name) ? '已添加' : '添加' }}
                                </button>
                                <button @click="viewSuggestionDetail(variable)" 
                                        class="action-btn mini secondary">
                                  详情
                                </button>
                              </div>
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </div>
                  </div>
                  
                  <!-- 文件建议 -->
                  <div v-for="suggestion in getFileSuggestions(message.suggestions)" 
                       :key="suggestion.name" 
                       class="suggestion-item file-suggestion">
                    <div class="suggestion-header">
                      <span class="suggestion-name">📄 {{ suggestion.name }}</span>
                      <div class="suggestion-actions">
                        <button @click="viewSuggestionDetail(suggestion)" class="action-btn secondary">预览</button>
                        <button @click="createFile(suggestion)" class="action-btn primary">创建文件</button>
                      </div>
                    </div>
                    <div class="suggestion-description">{{ suggestion.description }}</div>
                  </div>
                  
                  <!-- 通用操作建议 -->
                  <div v-for="suggestion in getActionSuggestions(message.suggestions)" 
                       :key="suggestion.name" 
                       class="suggestion-item action-suggestion">
                    <div class="suggestion-header">
                      <span class="suggestion-name">⚡ {{ suggestion.name }}</span>
                      <div class="suggestion-actions">
                        <button @click="viewSuggestionDetail(suggestion)" class="action-btn secondary">查看详情</button>
                        <button @click="executeSuggestion(suggestion)" class="action-btn primary">执行操作</button>
                      </div>
                    </div>
                    <div class="suggestion-description">{{ suggestion.description }}</div>
                  </div>
                </div>
              </div>
              
              <!-- 传统的消息操作（兼容） -->
              <div class="message-actions" v-if="message.actions && message.actions.length > 0 && !message.streaming && (!message.suggestions || message.suggestions.length === 0)">
                <button v-for="action in message.actions" 
                        :key="action.type"
                        @click="executeAction(action)"
                        class="action-btn">
                  {{ action.label }}
                </button>
              </div>
              
              <div class="message-time" v-if="!message.streaming">{{ formatTime(message.timestamp) }}</div>
            </div>
          </div>
        </div>

        <!-- 打字指示器 -->
        <div v-if="aiTyping" class="typing-indicator">
          <div class="ai-avatar">🤖</div>
          <div class="typing-dots">
            <span></span><span></span><span></span>
          </div>
        </div>
      </div>

      <!-- 快捷操作按钮 -->
      <div class="ai-quick-actions">
        <!-- AI建议的快捷操作 -->
        <div v-if="latestAISuggestions.length > 0" class="action-row ai-suggestions-row">
          <button v-for="suggestion in latestAISuggestions" 
                  :key="suggestion.name"
                  @click="executeSuggestionAction(suggestion)" 
                  class="quick-btn"
                  :class="getPriorityClass(suggestion.priority)">
            {{ getSuggestionIcon(suggestion.type) }} {{ suggestion.name }}
          </button>
        </div>
        
        <!-- 默认快捷操作 -->
        <div class="action-row default-actions-row">
          <template v-if="props.editorSelection.hasSelection">
            <!-- 有选中文本时的操作 -->
            <button @click="quickAction('optimize-selection')" class="quick-btn selection-action">
              ⚡ 优化选中代码
            </button>
            <button @click="quickAction('explain-selection')" class="quick-btn selection-action">
              💡 解释选中代码
            </button>
            <button @click="quickAction('refactor-selection')" class="quick-btn selection-action">
              🔧 重构选中代码
            </button>
            <button @click="quickAction('comment-selection')" class="quick-btn selection-action">
              📝 添加注释
            </button>
          </template>
          <template v-else>
            <!-- 没有选中文本时的默认操作 -->
            <button @click="quickAction('optimize')" class="quick-btn">
              ⚡ 优化代码
            </button>
            <button @click="quickAction('suggest-variables')" class="quick-btn">
              📝 建议变量
            </button>
            <button @click="quickAction('generate-template')" class="quick-btn">
              🎯 生成模板
            </button>
            <button @click="quickAction('explain')" class="quick-btn">
              💡 解释代码
            </button>
          </template>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="ai-input-area">
        <div class="input-wrapper">
          <textarea v-model="userInput" 
                    @keydown="handleKeydown"
                    placeholder="向AI描述你的需求..."
                    class="chat-input"
                    rows="1"
                    ref="chatInput"></textarea>
          <button @click="sendMessage" 
                  :disabled="!userInput.trim() || aiProcessing"
                  class="send-button">
            <svg v-if="!aiProcessing" viewBox="0 0 24 24" fill="currentColor">
              <path d="M2,21L23,12L2,3V10L17,12L2,14V21Z"/>
            </svg>
            <div v-else class="loading-spinner"></div>
          </button>
        </div>
        <div class="input-hint">
          支持 Enter 换行，⌘+Enter 或 Ctrl+Enter 发送
        </div>
      </div>
    </div>
    
    <!-- 详情弹窗 -->
    <div v-if="detailModalVisible" class="detail-modal-overlay" @click.self="detailModalVisible = false">
      <div class="detail-modal" @click.stop>
        <div class="detail-modal-header">
          <h3 class="detail-modal-title">
            <span class="detail-icon">{{ getSuggestionIcon(currentSuggestion?.type) }}</span>
            {{ currentSuggestion?.name }}
          </h3>
          <button @click="detailModalVisible = false" class="detail-modal-close">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z"/>
            </svg>
          </button>
        </div>
        
        <div class="detail-modal-content">
          <!-- 基本信息 -->
          <div class="detail-section">
            <h4 class="detail-section-title">描述信息</h4>
            <p class="detail-description">{{ currentSuggestion?.description }}</p>
          </div>
          
          <!-- 优先级和类型 -->
          <div class="detail-section" v-if="currentSuggestion?.priority || currentSuggestion?.type">
            <h4 class="detail-section-title">属性信息</h4>
            <div class="detail-meta">
              <span v-if="currentSuggestion?.type" class="detail-tag">
                类型: {{ getSuggestionTypeText(currentSuggestion.type) }}
              </span>
              <span v-if="currentSuggestion?.priority" class="detail-tag priority" :class="currentSuggestion.priority">
                优先级: {{ getPriorityText(currentSuggestion.priority) }}
              </span>
            </div>
          </div>
          
          <!-- 代码内容 -->
          <div class="detail-section" v-if="currentSuggestion?.code">
            <h4 class="detail-section-title">代码内容</h4>
            <div class="detail-code-editor">
              <div class="code-editor-header">
                <span class="code-language">{{ getCodeLanguage(currentSuggestion.code) }}</span>
                <button @click="copyCode(currentSuggestion.code)" class="copy-btn">
                  <svg viewBox="0 0 24 24" fill="currentColor">
                    <path d="M19,21H8V7H19M19,5H8A2,2 0 0,0 6,7V21A2,2 0 0,0 8,23H19A2,2 0 0,0 21,21V7A2,2 0 0,0 19,5M16,1H4A2,2 0 0,0 2,3V17H4V3H16V1Z"/>
                  </svg>
                  复制
                </button>
              </div>
              <pre class="code-content"><code>{{ currentSuggestion.code }}</code></pre>
            </div>
          </div>
        </div>
        
        <div class="detail-modal-actions">
          <button @click="detailModalVisible = false" class="detail-btn secondary">关闭</button>
          <button v-if="currentSuggestion?.code" @click="explainCodeFromModal" class="detail-btn tertiary">
            💡 解释代码
          </button>
          <button v-if="currentSuggestion?.code" @click="insertCodeFromModal" class="detail-btn primary">
            插入代码
          </button>
          <button v-if="currentSuggestion?.type === 'variable'" @click="addVariableFromModal" class="detail-btn primary">
            添加变量
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick, onMounted, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { testAIConnection, chatWithAI, chatWithAIStream, generateTemplate, suggestVariables, optimizeCode, explainCode } from '@/api/ai'
import { explainCodeStreamWithOpenAI, optimizeCodeWithOpenAI, suggestVariablesWithOpenAI } from '@/api/ai/openai'

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
const aiPanelVisible = ref(false)
const aiConnected = ref(true) // 默认连接状态
const aiProcessing = ref(false)
const aiTyping = ref(false)
const chatMessages = ref([])
const userInput = ref('')
const chatHistory = ref(null)
const chatInput = ref(null)

// API格式选择
const useOpenAIFormat = ref(JSON.parse(localStorage.getItem('ai-use-openai-format') || 'false'))

// 面板尺寸调整相关
const panelWidth = ref(parseInt(localStorage.getItem('ai-panel-width') || '400'))
const panelHeight = ref(parseInt(localStorage.getItem('ai-panel-height') || '600'))
const isResizing = ref(false)
const resizeStartX = ref(0)
const resizeStartY = ref(0)
const startWidth = ref(0)
const startHeight = ref(0)

// 详情弹窗相关
const detailModalVisible = ref(false)
const currentSuggestion = ref(null)

const message = useMessage()

// 计算属性：最新AI建议
const latestAISuggestions = computed(() => {
  // 获取最后一条AI消息的建议
  const lastAIMessage = chatMessages.value
    .filter(msg => msg.type === 'ai' && msg.suggestions && msg.suggestions.length > 0)
    .slice(-1)[0]
  
  return lastAIMessage?.suggestions || []
})

// 计算属性：上下文信息
const contextInfo = computed(() => ({
  fileName: props.currentFileName,
  variableCount: props.templateVariables.length,
  codeLines: props.currentFileContent ? props.currentFileContent.split('\n').length : 0
}))

// 生成消息ID
const generateMessageId = () => {
  return Date.now() + Math.random().toString(36).substr(2, 9)
}

// 切换AI面板
const toggleAIPanel = () => {
  aiPanelVisible.value = !aiPanelVisible.value
  if (aiPanelVisible.value) {
    initAIChat()
    nextTick(() => {
      if (chatInput.value) {
        chatInput.value.focus()
      }
    })
  }
}

// 切换API格式
const toggleAPIFormat = () => {
  useOpenAIFormat.value = !useOpenAIFormat.value
  localStorage.setItem('ai-use-openai-format', JSON.stringify(useOpenAIFormat.value))
  
  // 添加系统消息提示用户
  const formatName = useOpenAIFormat.value ? 'OpenAI兼容格式' : '原始API格式'
  addAIMessage(`已切换到 ${formatName} 🔄\n\n${useOpenAIFormat.value ? 
    '现在使用标准OpenAI接口，更好地兼容第三方AI SDK。' : 
    '现在使用原始内部接口，支持所有自定义功能。'}`)
}

// 初始化AI对话
const initAIChat = () => {
  if (chatMessages.value.length === 0) {
    addAIMessage(
      `你好！我是AI助手 🤖\n\n` +
      `**我可以帮你：**\n` +
      `• 💡 优化和改进当前代码\n` +
      `• 📝 建议合适的模板变量\n` +
      `• 🎯 生成新的模板文件\n` +
      `• 🔍 解释复杂的代码逻辑\n` +
      `• 🛠️ 修复代码问题\n\n` +
      `试试下面的快捷操作，或者直接告诉我你的需求！`,
      []
    )
  }
}

// 添加用户消息
const addUserMessage = (content) => {
  const message = {
    id: generateMessageId(),
    type: 'user',
    content,
    timestamp: new Date()
  }
  chatMessages.value.push(message)
  scrollToBottom()
}

// 添加AI消息
const addAIMessage = (content, actions = []) => {
  const message = {
    id: generateMessageId(),
    type: 'ai',
    content,
    actions,
    timestamp: new Date()
  }
  chatMessages.value.push(message)
  scrollToBottom()
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (chatHistory.value) {
      chatHistory.value.scrollTop = chatHistory.value.scrollHeight
    }
  })
}

// 格式化时间
const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 格式化AI消息（支持Markdown基本语法和代码解释结构化）
const formatAIMessage = (content) => {
  // 检查是否是代码解释类型的内容
  if (isCodeExplanationContent(content)) {
    return formatCodeExplanation(content)
  }
  
  // 常规消息格式化
  return content
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>') // 粗体
    .replace(/\*(.*?)\*/g, '<em>$1</em>') // 斜体
    .replace(/`(.*?)`/g, '<code>$1</code>') // 行内代码
    .replace(/\n/g, '<br>') // 换行
}

// 检查是否是代码解释内容
const isCodeExplanationContent = (content) => {
  // 判断是否包含代码解释的关键词和结构
  const indicators = [
    '主要功能',
    '具体逻辑',
    '实现原理',
    '代码功能',
    '这段代码',
    '函数作用'
  ]
  return indicators.some(indicator => content.includes(indicator))
}

// 格式化代码解释内容
const formatCodeExplanation = (content) => {
  let formatted = content
  
  // 1. 突出显示主要功能描述
  formatted = formatted.replace(
    /(主要功能[是：]?)(.*?)(?=[，。；]|$)/g,
    '<div class="explanation-section main-function"><div class="section-header">🎯 主要功能</div><div class="section-content">$2</div></div>'
  )
  
  // 2. 格式化逻辑分解（数字列表）
  formatted = formatted.replace(
    /(\d+\.?\s*)(.*?)(?=\d+\.|\n|$)/g,
    '<div class="explanation-step"><span class="step-number">$1</span><span class="step-content">$2</span></div>'
  )
  
  // 3. 突出显示技术关键词
  const techKeywords = [
    'GoFrame框架', 'HTTP服务', '配置文件', '上下文', '定时任务',
    'InitConfig', 'cmd.Main.Run', 'gctx.GetInitCtx',
    'main函数', '应用配置', '核心功能', '启动流程'
  ]
  
  techKeywords.forEach(keyword => {
    const regex = new RegExp(`(${keyword})`, 'g')
    formatted = formatted.replace(regex, '<span class="tech-keyword">$1</span>')
  })
  
  // 4. 格式化代码片段（行内代码）
  formatted = formatted.replace(/`([^`]+)`/g, '<code class="inline-code">$1</code>')
  
  // 5. 格式化分号分隔的逻辑点
  formatted = formatted.replace(/；/g, '；<br>')
  
  // 6. 包装在解释容器中
  formatted = `<div class="code-explanation">${formatted}</div>`
  
  return formatted
}

// 处理键盘事件
const handleKeydown = (event) => {
  if (event.key === 'Enter') {
    if (event.metaKey || event.ctrlKey) {
      // Cmd+Enter 或 Ctrl+Enter 发送消息
      event.preventDefault()
      sendMessage()
    }
    // 普通 Enter 换行（默认行为）
  }
}

// 检测消息的操作类型
const detectActionType = (message) => {
  const msg = message.toLowerCase()
  
  if (msg.includes('优化') || msg.includes('改进') || msg.includes('性能')) {
    return 'optimize_code'
  } else if (msg.includes('解释') || msg.includes('说明') || msg.includes('这是什么')) {
    return 'explain_code'
  } else if (msg.includes('变量') || msg.includes('建议') || msg.includes('参数')) {
    return 'suggest_variables'
  } else if (msg.includes('生成') || msg.includes('创建') || msg.includes('模板')) {
    return 'generate_template'
  } else if (msg.includes('重构') || msg.includes('重写')) {
    return 'refactor_code'
  } else if (msg.includes('注释') || msg.includes('文档')) {
    return 'add_comments'
  } else {
    return 'general_chat'
  }
}

// 发送消息
const sendMessage = async () => {
  if (!userInput.value.trim() || aiProcessing.value) return
  
  const message = userInput.value.trim()
  userInput.value = ''
  
  // 添加用户消息
  addUserMessage(message)
  
  // 显示AI正在输入
  aiTyping.value = true
  aiProcessing.value = true
  
  // 创建一个临时的AI消息用于流式更新
  const aiMessageId = generateMessageId()
  const aiMessage = {
    id: aiMessageId,
    type: 'ai',
    content: '',
    actions: [],
    timestamp: new Date(),
    streaming: true
  }
  chatMessages.value.push(aiMessage)
  scrollToBottom()
  
  try {
    // 获取当前上下文
    const context = getCurrentContext()
    
    // 构建请求数据
    const requestData = {
      message,
      action: detectActionType(message),
      context: {
        fileName: context.fileName,
        fileContent: context.fileContent,
        variables: context.variables,
        selectedText: context.selectedText || '',
        hasSelection: context.hasSelection
      },
      chatHistory: chatMessages.value
        .filter(msg => msg.type !== 'ai' || !msg.streaming) // 排除正在流式更新的消息
        .slice(-10) // 只发送最近10条消息作为上下文
        .map(msg => ({
          role: msg.type === 'user' ? 'user' : 'assistant',
          content: msg.content,
          timestamp: msg.timestamp
        }))
    }
    
    // 使用流式API
    const streamControl = await chatWithAIStream(
      requestData,
      // onChunk - 处理流式数据片段
      (chunk) => {
        const targetMessage = chatMessages.value.find(msg => msg.id === aiMessageId)
        if (targetMessage) {
          if (chunk.type === 'chunk') {
            targetMessage.content = chunk.fullContent
            scrollToBottom()
          }
        }
      },
      // onComplete - 流式响应完成
      (response) => {
        const targetMessage = chatMessages.value.find(msg => msg.id === aiMessageId)
        if (targetMessage) {
          targetMessage.content = response.content
          targetMessage.streaming = false
          targetMessage.metadata = response.metadata
          
          // 添加建议操作（新的结构化格式）
          if (response.suggestions && response.suggestions.length > 0) {
            targetMessage.suggestions = response.suggestions
            targetMessage.suggestionsCollapsed = false // 建议默认展开
          }
          
          scrollToBottom()
        }
        
        aiTyping.value = false
        aiProcessing.value = false
      },
      // onError - 处理错误
      (error) => {
        console.error('流式AI API调用失败:', error)
        
        const targetMessage = chatMessages.value.find(msg => msg.id === aiMessageId)
        if (targetMessage) {
          targetMessage.content = '抱歉，我遇到了一些问题 😅\n\n请稍后再试，或者检查网络连接。'
          targetMessage.streaming = false
          targetMessage.actions = [{
            type: 'retry',
            label: '🔄 重试',
            action: () => sendMessage()
          }]
        }
        
        aiTyping.value = false
        aiProcessing.value = false
      }
    )
    
    // 存储流控制对象，用于可能的取消操作
    // 可以在组件中添加取消按钮调用 streamControl.close()
    
  } catch (error) {
    console.error('启动流式AI API失败:', error)
    
    const targetMessage = chatMessages.value.find(msg => msg.id === aiMessageId)
    if (targetMessage) {
      targetMessage.content = '抱歉，启动AI对话失败 😅\n\n请检查网络连接后重试。'
      targetMessage.streaming = false
      targetMessage.actions = [{
        type: 'retry',
        label: '🔄 重试',
        action: () => sendMessage()
      }]
    }
    
    aiTyping.value = false
    aiProcessing.value = false
  }
}

// 获取当前编辑器上下文
const getCurrentContext = () => {
  const hasSelection = props.editorSelection.hasSelection
  const targetContent = hasSelection ? props.editorSelection.selectedText : props.currentFileContent
  
  return {
    fileName: props.currentFileName,
    fileContent: props.currentFileContent,
    targetContent, // 优先使用选中的内容，如果没有选中则使用整个文件内容
    variables: props.templateVariables,
    variableCount: props.templateVariables.length,
    codeLines: props.currentFileContent ? props.currentFileContent.split('\n').length : 0,
    hasSelection,
    selectedText: props.editorSelection.selectedText,
    selectionLength: props.editorSelection.selectionLength,
    selectionStart: props.editorSelection.selectionStart,
    selectionEnd: props.editorSelection.selectionEnd
  }
}

// 调用AI API
const callAIAPI = async ({ message, context, chatHistory }) => {
  try {
    // 构建请求数据
    const requestData = {
      message,
      context: {
        fileName: context.fileName,
        fileContent: context.fileContent,
        variables: context.variables,
        selectedText: context.selectedText || '',
        cursorPosition: context.cursorPosition || 0
      },
      chatHistory: chatHistory.map(msg => ({
        role: msg.type === 'user' ? 'user' : 'assistant',
        content: msg.content,
        timestamp: msg.timestamp
      }))
    }
    
    // 根据消息类型和API格式选择不同的API
    let apiResponse
    
    if (message.includes('优化') || message.includes('改进')) {
      // 调用代码优化API
      if (useOpenAIFormat.value) {
        apiResponse = await optimizeCodeWithOpenAI({
          code: context.targetContent,
          language: getFileLanguage(context.fileName),
          requirements: message
        })
      } else {
        apiResponse = await optimizeCode({
          code: context.targetContent, // 使用目标内容（选中文本或整个文件）
          fileName: context.fileName,
          requirements: message,
          hasSelection: context.hasSelection,
          selectionInfo: context.hasSelection ? {
            selectedText: context.selectedText,
            selectionStart: context.selectionStart,
            selectionEnd: context.selectionEnd
          } : null
        })
      }
    } else if (message.includes('变量') || message.includes('建议')) {
      // 调用变量建议API
      if (useOpenAIFormat.value) {
        apiResponse = await suggestVariablesWithOpenAI({
          projectType: inferProjectType(context.fileName),
          techStack: inferTechStack(context.fileName, context.fileContent),
          content: context.targetContent
        })
      } else {
        apiResponse = await suggestVariables({
          projectType: inferProjectType(context.fileName),
          techStack: inferTechStack(context.fileName, context.fileContent),
          description: message,
          existingVariables: context.variables
        })
      }
    } else if (message.includes('生成') || message.includes('模板')) {
      // 调用模板生成API
      apiResponse = await generateTemplate({
        description: message,
        projectType: inferProjectType(context.fileName),
        techStack: inferTechStack(context.fileName, context.fileContent),
        variables: context.variables.reduce((acc, v) => {
          acc[v.name] = v.defaultValue || ''
          return acc
        }, {}),
        features: []
      })
    } else if (message.includes('解释') || message.includes('说明')) {
      // 调用代码解释API
      if (useOpenAIFormat.value) {
        // 使用流式OpenAI API进行代码解释
        return handleOpenAIStreamExplanation(message, context, aiMessageId)
      } else {
        apiResponse = await explainCode({
          code: context.targetContent, // 使用目标内容
          fileName: context.fileName,
          question: message,
          hasSelection: context.hasSelection,
          fullFileContent: context.fileContent // 保留完整文件内容作为上下文
        })
      }
    } else {
      // 通用对话API
      apiResponse = await chatWithAI({
        action: 'general_chat',
        context: {
          fileName: context.fileName,
          fileContent: context.fileContent,
          selectedText: context.selectedText,
          hasSelection: context.hasSelection,
          variables: context.variables
        },
        userInput: message,
        chatHistory: chatHistory.map(msg => ({
          role: msg.type === 'user' ? 'user' : 'assistant',
          content: msg.content,
          timestamp: msg.timestamp
        }))
      })
    }
    
    // 统一处理响应格式
    return processAIResponse(apiResponse, message, context)
    
  } catch (error) {
    console.error('AI API调用失败:', error)
    
    // 如果是网络错误或API不可用，回退到模拟响应
    if (error.code === 'NETWORK_ERROR' || error.response?.status >= 500) {
      return getFallbackResponse(message, context)
    }
    
    throw error
  }
}

// 处理AI响应
const processAIResponse = (apiResponse, message, context) => {
  const response = {
    content: '',
    actions: []
  }
  
  if (apiResponse.data && apiResponse.data.data) {
    const data = apiResponse.data.data
    
    // 根据API返回的数据结构处理
    if (data.content || data.response) {
      response.content = data.content || data.response
    }
    
    if (data.suggestions) {
      // 处理建议操作
      data.suggestions.forEach(suggestion => {
        switch (suggestion.type) {
          case 'code':
            response.actions.push({
              type: 'insert-code',
              label: `📝 插入${suggestion.name || '代码'}`,
              code: suggestion.code
            })
            break
          case 'variable':
            response.actions.push({
              type: 'add-variable',
              label: `➕ 添加变量: ${suggestion.name}`,
              variable: suggestion
            })
            break
          case 'file':
            response.actions.push({
              type: 'create-file',
              label: `📁 创建文件: ${suggestion.fileName}`,
              fileName: suggestion.fileName,
              content: suggestion.content,
              fileType: suggestion.isDirectory ? 'folder' : 'file'
            })
            break
        }
      })
    }
    
    if (data.projectStructure) {
      // 处理项目结构
      response.actions.push({
        type: 'apply-suggestion',
        label: '🎯 应用项目结构',
        suggestion: {
          type: 'template-structure',
          structure: data.projectStructure
        }
      })
    }
    
    if (data.variables) {
      // 处理变量建议
      response.actions.push({
        type: 'apply-suggestion',
        label: '📝 应用变量建议',
        suggestion: {
          type: 'variable-suggestion',
          variables: data.variables
        }
      })
    }
  }
  
  // 如果没有内容，提供默认响应
  if (!response.content) {
    response.content = '我已经处理了你的请求，请查看建议的操作按钮。'
  }
  
  return response
}

// 回退响应（当API不可用时）
const getFallbackResponse = (message, context) => {
  let response = {
    content: '',
    actions: []
  }
  
  // 根据消息内容生成不同的响应
  if (message.includes('优化') || message.includes('改进')) {
    response.content = `我来帮你优化代码！ 💡\n\n基于当前的 **${context.fileName}** 文件，我建议：\n\n• 🔧 **代码结构优化**：提高可读性\n• ⚡ **性能优化**：减少不必要的计算\n• 📝 **注释完善**：添加必要的说明\n\n*注：当前使用离线模式，功能有限*`
    response.actions = [
      { type: 'optimize-structure', label: '🔧 优化结构' },
      { type: 'optimize-performance', label: '⚡ 优化性能' },
      { type: 'add-comments', label: '📝 添加注释' }
    ]
  } else if (message.includes('变量') || message.includes('建议')) {
    response.content = `根据你的模板内容，我建议添加以下变量：\n\n• **projectName** (字符串) - 项目名称\n• **author** (字符串) - 作者信息\n• **version** (字符串) - 版本号\n• **description** (字符串) - 项目描述\n\n*注：当前使用离线模式，建议有限*`
    response.actions = [
      { type: 'add-variable', label: '➕ 添加建议的变量', variable: { name: 'projectName', type: 'string', description: '项目名称' } }
    ]
  } else {
    response.content = `抱歉，AI服务当前不可用 😔\n\n我正在使用离线模式为你提供基础帮助：\n\n• 📝 **基础代码建议**\n• 🔧 **简单优化提示**\n• 💡 **通用最佳实践**\n\n请稍后重试连接AI服务获得更好的帮助。`
    response.actions = [
      { type: 'retry', label: '🔄 重试连接' }
    ]
  }
  
  return response
}

// 推断项目类型
const inferProjectType = (fileName) => {
  if (!fileName) return 'web'
  
  const ext = fileName.split('.').pop()?.toLowerCase()
  const typeMap = {
    'vue': 'web',
    'js': 'web',
    'ts': 'web',
    'jsx': 'web',
    'tsx': 'web',
    'html': 'web',
    'css': 'web',
    'go': 'backend',
    'py': 'backend',
    'java': 'backend',
    'php': 'backend',
    'swift': 'mobile',
    'kt': 'mobile',
    'dart': 'mobile'
  }
  
  return typeMap[ext] || 'web'
}

// 推断技术栈
const inferTechStack = (fileName, fileContent) => {
  const techStack = []
  
  if (!fileName || !fileContent) return techStack
  
  const ext = fileName.split('.').pop()?.toLowerCase()
  const content = fileContent.toLowerCase()
  
  // 根据文件扩展名
  if (['vue', 'js', 'ts'].includes(ext)) {
    techStack.push('javascript')
    if (ext === 'ts' || content.includes('typescript')) {
      techStack.push('typescript')
    }
    if (ext === 'vue' || content.includes('vue')) {
      techStack.push('vue')
    }
  }
  
  // 根据内容检测
  if (content.includes('react')) techStack.push('react')
  if (content.includes('angular')) techStack.push('angular')
  if (content.includes('express')) techStack.push('express')
  if (content.includes('fastify')) techStack.push('fastify')
  if (content.includes('mysql')) techStack.push('mysql')
  if (content.includes('postgresql')) techStack.push('postgresql')
  if (content.includes('mongodb')) techStack.push('mongodb')
  
  return [...new Set(techStack)] // 去重
}

// 获取文件类型
const getFileType = (fileName) => {
  const ext = fileName.split('.').pop()?.toLowerCase()
  const typeMap = {
    'js': 'JavaScript',
    'ts': 'TypeScript',
    'vue': 'Vue组件',
    'html': 'HTML',
    'css': 'CSS',
    'json': 'JSON配置',
    'md': 'Markdown文档',
    'go': 'Go语言',
    'py': 'Python'
  }
  return typeMap[ext] || '文本文件'
}

// 快捷操作
const quickAction = async (action) => {
  const hasSelection = props.editorSelection.hasSelection
  const selectionLength = props.editorSelection.selectionLength
  
  const actionMessages = {
    // 针对整个文件的操作
    'optimize': '请帮我优化当前的代码，提高可读性和性能',
    'suggest-variables': '基于当前模板内容，建议一些有用的变量',
    'generate-template': '我想生成一个新的模板文件',
    'explain': '请解释一下当前代码的功能和逻辑',
    
    // 针对选中文本的操作
    'optimize-selection': `请帮我优化选中的这${selectionLength}个字符的代码，提高可读性和性能`,
    'explain-selection': `请解释一下我选中的这${selectionLength}个字符的代码的功能和逻辑`,
    'refactor-selection': `请帮我重构选中的这${selectionLength}个字符的代码，提供更好的实现方案`,
    'comment-selection': `请为我选中的这${selectionLength}个字符的代码添加详细的注释说明`
  }
  
  const message = actionMessages[action]
  if (message) {
    userInput.value = message
    await sendMessage()
  }
}

// 执行AI建议的操作
const executeAction = (action) => {
  switch (action.type) {
    case 'insert-code':
      // 插入代码到编辑器
      emit('insertCode', action.code)
      message.success('代码已插入到编辑器')
      break
    case 'add-variable':
      // 添加变量
      emit('addVariable', action.variable)
      message.success('变量建议已发送')
      break
    case 'create-file':
      // 创建新文件
      emit('createFile', {
        type: action.fileType,
        name: action.fileName,
        content: action.content,
        parentId: action.parentId
      })
      message.success('文件创建请求已发送')
      break
    case 'apply-suggestion':
      // 应用建议
      emit('applySuggestion', action.suggestion)
      message.success('建议已应用')
      break
    case 'retry':
      // 重试
      if (action.action) {
        action.action()
      }
      break
    default:
      // 其他操作转换为消息
      userInput.value = getActionMessage(action.type)
      sendMessage()
      break
  }
}

// 获取操作对应的消息
const getActionMessage = (actionType) => {
  const actionMessages = {
    'optimize-structure': '请帮我优化代码结构，提高可读性',
    'optimize-performance': '请帮我优化代码性能',
    'add-comments': '请帮我添加必要的代码注释',
    'suggest-more': '请给我更多的变量建议',
    'generate-web': '请帮我生成一个Web前端项目模板',
    'generate-api': '请帮我生成一个后端API项目模板',
    'generate-mobile': '请帮我生成一个移动应用项目模板',
    'explain-structure': '请解释一下代码的整体结构',
    'explain-variables': '请解释一下变量的作用和用法',
    'get-help': '我需要更多帮助，请告诉我你能做什么',
    'show-examples': '请给我一些模板使用的示例'
  }
  return actionMessages[actionType] || '请帮助我'
}

// 切换消息折叠状态
const toggleMessageCollapse = (message) => {
  message.collapsed = !message.collapsed
}

// 处理建议区域折叠
const toggleSuggestionsCollapse = (message) => {
  message.suggestionsCollapsed = !message.suggestionsCollapsed
}

// 获取消息摘要
const getMessageSummary = (message) => {
  if (message.suggestions && message.suggestions.length > 0) {
    const types = [...new Set(message.suggestions.map(s => s.type))]
    const typeNames = {
      code: '代码建议',
      variable: '变量建议', 
      file: '文件建议',
      action: '操作建议'
    }
    const typeLabels = types.map(type => typeNames[type] || type).join('、')
    return `AI提供了${message.suggestions.length}个${typeLabels}`
  }
  
  // 如果没有建议，显示内容的前50个字符
  const content = message.content || ''
  return content.length > 50 ? content.substring(0, 50) + '...' : content
}

// 根据类型过滤建议
const getCodeSuggestions = (suggestions) => {
  return suggestions.filter(s => s.type === 'code')
}

const getVariableSuggestions = (suggestions) => {
  return suggestions.filter(s => s.type === 'variable')
}

const getFileSuggestions = (suggestions) => {
  return suggestions.filter(s => s.type === 'file')
}

const getActionSuggestions = (suggestions) => {
  return suggestions.filter(s => s.type === 'action')
}

// 变量排序和处理
const getSortedVariables = (variables) => {
  return variables.sort((a, b) => {
    // 先按优先级排序
    const priorityOrder = { high: 3, medium: 2, low: 1 }
    const aPriority = priorityOrder[a.priority] || 0
    const bPriority = priorityOrder[b.priority] || 0
    
    if (aPriority !== bPriority) {
      return bPriority - aPriority
    }
    
    // 默认排序
    return 0
  })
}

// 获取优先级文本
const getPriorityText = (priority) => {
  const priorityMap = {
    high: '高',
    medium: '中',
    low: '低'
  }
  return priorityMap[priority] || '未知'
}


// 检查变量是否已添加
const isVariableAdded = (variableName) => {
  return props.templateVariables.some(v => v.name === variableName)
}

// 添加所有变量
const addAllVariables = (variables) => {
  const newVariables = variables.filter(v => !isVariableAdded(v.name))
  
  if (newVariables.length === 0) {
    message.info('所有变量都已添加')
    return
  }
  
  newVariables.forEach(variable => {
    emit('addVariable', {
      name: variable.name,
      description: variable.description,
      defaultValue: variable.code?.replace(/[{}\.]/g, '') || variable.name,
      type: 'string',
      required: variable.priority === 'high'
    })
  })
  
  message.success(`已添加 ${newVariables.length} 个新变量`)
}

// 处理建议操作
const viewSuggestionDetail = (suggestion) => {
  currentSuggestion.value = suggestion
  detailModalVisible.value = true
}

// 获取建议类型图标
const getSuggestionIcon = (type) => {
  const iconMap = {
    'code': '🔧',
    'variable': '📝',
    'file': '📄',
    'action': '⚡'
  }
  return iconMap[type] || '💡'
}

// 获取建议类型文本
const getSuggestionTypeText = (type) => {
  const typeMap = {
    'code': '代码建议',
    'variable': '变量建议',
    'file': '文件建议',
    'action': '操作建议'
  }
  return typeMap[type] || '未知类型'
}

// 获取代码语言
const getCodeLanguage = (code) => {
  if (!code) return 'text'
  
  // 根据代码内容推断语言
  if (code.includes('function') || code.includes('const') || code.includes('let')) return 'JavaScript'
  if (code.includes('def ') || code.includes('import ')) return 'Python'
  if (code.includes('<?php') || code.includes('$')) return 'PHP'
  if (code.includes('<template>') || code.includes('<div>')) return 'Vue/HTML'
  if (code.includes('func ') || code.includes('package ')) return 'Go'
  if (code.includes('class ') && code.includes('public ')) return 'Java'
  
  return 'Code'
}

// 复制代码
const copyCode = async (code) => {
  try {
    await navigator.clipboard.writeText(code)
    message.success('代码已复制到剪贴板')
  } catch (err) {
    // 回退方案
    const textArea = document.createElement('textarea')
    textArea.value = code
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    message.success('代码已复制到剪贴板')
  }
}

// 从弹窗插入代码
const insertCodeFromModal = () => {
  if (currentSuggestion.value?.code) {
    insertCode(currentSuggestion.value.code)
    detailModalVisible.value = false
  }
}

// 从弹窗添加变量
const addVariableFromModal = () => {
  if (currentSuggestion.value?.type === 'variable') {
    addVariable(currentSuggestion.value)
    detailModalVisible.value = false
  }
}

// 执行建议操作
const executeSuggestionAction = (suggestion) => {
  // 根据建议类型执行不同的操作
  switch(suggestion.type) {
    case 'action':
      // 将建议描述转换为用户消息发送给AI
      userInput.value = suggestion.description
      sendMessage()
      break
    case 'code':
      // 如果有代码，可以插入或解释
      if (suggestion.code) {
        insertCode(suggestion.code)
      }
      break
    case 'variable':
      // 添加变量
      addVariable(suggestion)
      break
    default:
      // 其他类型转换为消息
      userInput.value = `请帮我：${suggestion.description}`
      sendMessage()
      break
  }
}

// 获取优先级样式类
const getPriorityClass = (priority) => {
  return `priority-${priority || 'medium'}`
}

// 从弹窗解释代码
const explainCodeFromModal = () => {
  if (currentSuggestion.value?.code) {
    // 关闭弹窗
    detailModalVisible.value = false
    
    // 构建解释代码的消息
    const explainMessage = `请详细解释这段代码的功能和实现原理：\n\n${currentSuggestion.value.code}`
    
    // 发送消息
    userInput.value = explainMessage
    sendMessage()
  }
}

// 处理OpenAI流式代码解释
const handleOpenAIStreamExplanation = async (message, context, aiMessageId) => {
  try {
    const streamControl = await explainCodeStreamWithOpenAI({
      code: context.targetContent,
      language: getFileLanguage(context.fileName),
      hasSelection: context.hasSelection,
      selectedCode: context.hasSelection ? context.selectedText : null,
      fullCode: context.fileContent
    }, 
    // onChunk
    (chunk) => {
      const targetMessage = chatMessages.value.find(msg => msg.id === aiMessageId)
      if (targetMessage) {
        if (chunk.type === 'chunk') {
          targetMessage.content = chunk.fullContent
          scrollToBottom()
        }
      }
    },
    // onComplete
    (response) => {
      const targetMessage = chatMessages.value.find(msg => msg.id === aiMessageId)
      if (targetMessage) {
        targetMessage.content = response.content
        targetMessage.streaming = false
        targetMessage.metadata = response.metadata
        
        // 处理OpenAI格式的建议
        if (response.suggestions && response.suggestions.length > 0) {
          targetMessage.suggestions = response.suggestions
          targetMessage.suggestionsCollapsed = false
        }
        
        scrollToBottom()
      }
      
      aiTyping.value = false
      aiProcessing.value = false
    },
    // onError
    (error) => {
      console.error('OpenAI流式API调用失败:', error)
      
      const targetMessage = chatMessages.value.find(msg => msg.id === aiMessageId)
      if (targetMessage) {
        targetMessage.content = '抱歉，OpenAI API调用遇到了问题 😅\n\n请稍后再试，或者切换到原始API格式。'
        targetMessage.streaming = false
      }
      
      aiTyping.value = false
      aiProcessing.value = false
    })
  } catch (error) {
    console.error('启动OpenAI流式调用失败:', error)
    
    const targetMessage = chatMessages.value.find(msg => msg.id === aiMessageId)
    if (targetMessage) {
      targetMessage.content = '启动OpenAI流式调用失败，请检查网络连接。'
      targetMessage.streaming = false
    }
    
    aiTyping.value = false
    aiProcessing.value = false
  }
}

// 获取文件语言类型
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

const insertCode = (code) => {
  // 向父组件发射插入代码的事件
  emit('insertCode', code)
  message.success('代码已插入到编辑器')
}

const addVariable = (variable) => {
  // 向父组件发射添加变量的事件
  emit('addVariable', {
    name: variable.name,
    description: variable.description,
    defaultValue: variable.code?.replace(/[{}\.]/g, '') || variable.name,
    type: 'string',
    required: false
  })
  message.success(`变量 ${variable.name} 已添加`)
}

const createFile = (fileSuggestion) => {
  // 向父组件发射创建文件的事件
  emit('createFile', {
    name: fileSuggestion.name,
    content: fileSuggestion.code || '',
    description: fileSuggestion.description,
    type: 'file'
  })
  message.success(`文件 ${fileSuggestion.name} 创建请求已发送`)
}

const executeSuggestion = (suggestion) => {
  // 根据建议类型执行相应操作
  switch (suggestion.type) {
    case 'action':
      // 向父组件发射执行操作的事件
      emit('applySuggestion', {
        action: suggestion.name,
        description: suggestion.description,
        data: suggestion
      })
      message.success(`操作 ${suggestion.name} 已执行`)
      break
    case 'code':
      insertCode(suggestion.code)
      break
    case 'variable':
      addVariable(suggestion)
      break
    case 'file':
      createFile(suggestion)
      break
    default:
      message.info(`执行操作: ${suggestion.name}`)
  }
}

// 监听面板显示状态，自动滚动到底部
watch(aiPanelVisible, (visible) => {
  if (visible) {
    nextTick(() => {
      scrollToBottom()
    })
  }
})

// 监听消息变化，自动滚动到底部
watch(chatMessages, () => {
  scrollToBottom()
}, { deep: true })

// 监听文件名变化，当切换文件时添加系统提示
watch(() => props.currentFileName, (newFileName, oldFileName) => {
  if (aiPanelVisible.value && newFileName && oldFileName && newFileName !== oldFileName) {
    // 当切换文件时，添加一条简洁的系统提示消息
    addAIMessage(
      `📄 已切换到 **${newFileName}**`,
      []
    )
  }
})

// 组件挂载时检查AI连接状态
onMounted(() => {
  // TODO: 实际检查AI服务连接状态
  checkAIConnection()
})

// 检查AI连接状态
const checkAIConnection = async () => {
  try {
    const response = await testAIConnection()
    if (response.data && response.data.data) {
      aiConnected.value = response.data.data.success || false
      
      // 如果连接成功，更新状态信息
      if (aiConnected.value) {
        console.log('AI服务连接成功:', {
          provider: response.data.data.provider,
          model: response.data.data.model,
          latency: response.data.data.latency
        })
      }
    } else {
      aiConnected.value = false
    }
  } catch (error) {
    aiConnected.value = false
    console.error('AI连接检查失败:', error)
    
    // 如果是401或403错误，说明可能是配置问题
    if (error.response?.status === 401 || error.response?.status === 403) {
      console.warn('AI服务认证失败，请检查配置')
    }
  }
}

// 开始调整面板大小
const startResize = (event) => {
  isResizing.value = true
  resizeStartX.value = event.clientX
  resizeStartY.value = event.clientY
  startWidth.value = panelWidth.value
  startHeight.value = panelHeight.value
  
  // 添加全局事件监听
  document.addEventListener('mousemove', handleResize)
  document.addEventListener('mouseup', stopResize)
  
  // 防止文本选择
  event.preventDefault()
}

// 处理面板大小调整
const handleResize = (event) => {
  if (!isResizing.value) return
  
  const deltaX = resizeStartX.value - event.clientX
  const deltaY = resizeStartY.value - event.clientY
  
  // 根据屏幕尺寸动态调整限制
  const maxWidth = Math.min(1200, window.innerWidth - 48) // 增加最大宽度到1200px
  const maxHeight = Math.min(800, window.innerHeight - 120) // 留出按钮和边距
  
  // 计算新的宽度和高度
  const newWidth = Math.max(300, Math.min(maxWidth, startWidth.value + deltaX))
  const newHeight = Math.max(400, Math.min(maxHeight, startHeight.value + deltaY))
  
  panelWidth.value = newWidth
  panelHeight.value = newHeight
}

// 停止调整面板大小
const stopResize = () => {
  isResizing.value = false
  
  // 保存尺寸到本地存储
  localStorage.setItem('ai-panel-width', panelWidth.value.toString())
  localStorage.setItem('ai-panel-height', panelHeight.value.toString())
  
  // 移除全局事件监听
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
  document.removeEventListener('touchmove', handleResizeTouch)
  document.removeEventListener('touchend', stopResize)
}

// 触摸开始调整
const startResizeTouch = (event) => {
  const touch = event.touches[0]
  if (!touch) return
  
  isResizing.value = true
  resizeStartX.value = touch.clientX
  resizeStartY.value = touch.clientY
  startWidth.value = panelWidth.value
  startHeight.value = panelHeight.value
  
  // 添加触摸事件监听
  document.addEventListener('touchmove', handleResizeTouch, { passive: false })
  document.addEventListener('touchend', stopResize)
  
  // 防止默认行为
  event.preventDefault()
}

// 处理触摸调整
const handleResizeTouch = (event) => {
  if (!isResizing.value) return
  
  const touch = event.touches[0]
  if (!touch) return
  
  const deltaX = resizeStartX.value - touch.clientX
  const deltaY = resizeStartY.value - touch.clientY
  
  // 根据屏幕尺寸动态调整限制
  const maxWidth = Math.min(1200, window.innerWidth - 48) // 增加最大宽度到1200px
  const maxHeight = Math.min(800, window.innerHeight - 120)
  
  // 计算新的宽度和高度
  const newWidth = Math.max(300, Math.min(maxWidth, startWidth.value + deltaX))
  const newHeight = Math.max(400, Math.min(maxHeight, startHeight.value + deltaY))
  
  panelWidth.value = newWidth
  panelHeight.value = newHeight
  
  // 防止默认行为
  event.preventDefault()
}
</script>

<style scoped>
.ai-assistant {
  position: relative;
  z-index: 1000;
}

/* 悬浮按钮样式 */
.ai-float-button {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.3);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 1001;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ai-float-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(102, 126, 234, 0.4);
}

.ai-float-button.expanded {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
}

.ai-float-button.expanded:hover {
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
  top: -4px;
  left: -4px;
  right: -4px;
  bottom: -4px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  animation: pulse 2s infinite;
  z-index: -1;
}

/* 对话面板样式 */
.ai-chat-panel {
  position: fixed;
  bottom: 96px;
  right: 24px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 64px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  transform: translateY(100%) scale(0.8);
  opacity: 0;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 1000;
  overflow: hidden;
  border: 1px solid rgba(0, 0, 0, 0.06);
  min-width: 300px;
  max-width: 1200px;
  min-height: 400px;
  max-height: 800px;
  resize: both;
}

.ai-chat-panel.visible {
  transform: translateY(0) scale(1);
  opacity: 1;
}

.ai-chat-panel.resizing {
  transition: none;
  user-select: none;
}

/* 拖拽调整手柄 */
.resize-handle {
  position: absolute;
  top: 8px;
  left: 8px;
  width: 20px;
  height: 20px;
  cursor: nw-resize;
  z-index: 1001;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.resize-handle:hover {
  background: rgba(102, 126, 234, 0.1);
}

.resize-handle.active {
  background: rgba(102, 126, 234, 0.2);
}

.resize-dots {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
  gap: 2px;
  width: 8px;
  height: 8px;
}

.resize-dots span {
  width: 2px;
  height: 2px;
  background: #64748b;
  border-radius: 50%;
  transition: all 0.2s ease;
}

.resize-handle:hover .resize-dots span {
  background: #667eea;
  transform: scale(1.2);
}

.resize-handle.active .resize-dots span {
  background: #667eea;
  transform: scale(1.5);
}

/* 面板头部 */
.ai-panel-header {
  padding: 16px 20px 16px 32px; /* 为拖拽手柄留出空间 */
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  border-bottom: 1px solid #e2e8f0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .ai-chat-panel {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    width: 100% !important;
    height: 100% !important;
    border-radius: 0;
    max-width: none;
    max-height: none;
  }
  
  .resize-handle {
    display: none;
  }
  
  .ai-panel-header {
    padding: 16px 20px;
  }
}

@media (max-width: 480px) {
  .ai-context-info {
    flex-direction: column;
    gap: 4px;
  }
  
  .variables-grid {
    font-size: 12px;
  }
  
  .variables-grid th,
  .variables-grid td {
    padding: 6px 8px;
  }
}

.ai-status-indicator {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 12px;
}

.api-format-switch {
  display: flex;
  align-items: center;
}

.format-switch-btn {
  padding: 4px 8px;
  font-size: 11px;
  font-weight: 600;
  border: 1.5px solid #e2e8f0;
  border-radius: 6px;
  background: #ffffff;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.format-switch-btn:hover {
  border-color: #cbd5e1;
  background: #f8fafc;
  color: #475569;
}

.format-switch-btn.openai {
  border-color: #10b981;
  background: #ecfdf5;
  color: #059669;
}

.format-switch-btn.openai:hover {
  border-color: #059669;
  background: #d1fae5;
  color: #047857;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #94a3b8;
  transition: all 0.3s;
}

.status-dot.online {
  background: #10b981;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.4);
}

.status-text {
  font-size: 14px;
  font-weight: 600;
  color: #334155;
}

.ai-context-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.context-file {
  font-size: 12px;
  color: #667eea;
  font-weight: 500;
}

.context-variables,
.context-lines,
.context-selection {
  font-size: 11px;
  color: #64748b;
}

.context-selection {
  color: #667eea;
  font-weight: 500;
}

/* 对话历史 */
.ai-chat-history {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  scroll-behavior: smooth;
}

.chat-message {
  margin-bottom: 16px;
}

/* 用户消息 */
.chat-message.user .user-message {
  text-align: right;
}

.user-message .message-content {
  display: inline-block;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 12px 16px;
  border-radius: 18px 18px 6px 18px;
  max-width: 80%;
  word-wrap: break-word;
  line-height: 1.4;
}

.user-message .message-time {
  font-size: 11px;
  color: #94a3b8;
  margin-top: 4px;
}

/* AI消息 */
.chat-message.ai .ai-message {
  display: flex;
  gap: 12px;
  align-items: flex-start;
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
  margin-top: 2px;
}

.message-body {
  flex: 1;
  min-width: 0;
}

.ai-message .message-content {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  padding: 12px 16px;
  border-radius: 6px 18px 18px 18px;
  color: #334155;
  line-height: 1.5;
}

.ai-message .message-content :deep(strong) {
  color: #1e293b;
  font-weight: 600;
}

.ai-message .message-content :deep(em) {
  color: #475569;
  font-style: italic;
}

.ai-message .message-content :deep(code) {
  background: #e2e8f0;
  color: #1e293b;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
}

.message-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
  flex-wrap: wrap;
}

.action-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 8px 12px;
  border-radius: 16px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.action-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.ai-message .message-time {
  font-size: 11px;
  color: #94a3b8;
  margin-top: 8px;
}

/* 打字指示器 */
.typing-indicator {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 16px;
}

.typing-dots {
  display: flex;
  gap: 4px;
  align-items: center;
}

.typing-dots span {
  width: 6px;
  height: 6px;
  background: #94a3b8;
  border-radius: 50%;
  animation: typing 1.4s infinite ease-in-out both;
}

.typing-dots span:nth-child(1) { animation-delay: -0.32s; }
.typing-dots span:nth-child(2) { animation-delay: -0.16s; }

/* 流式状态指示器 */
.streaming-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  padding: 8px 12px;
  background: #f1f5f9;
  border-radius: 6px;
  font-size: 12px;
  color: #64748b;
}

.streaming-indicator .typing-dots {
  display: flex;
  gap: 3px;
}

.streaming-indicator .typing-dots span {
  width: 4px;
  height: 4px;
  background: #64748b;
  border-radius: 50%;
  animation: streaming-pulse 1.4s infinite ease-in-out both;
}

.streaming-indicator .typing-dots span:nth-child(1) { animation-delay: -0.32s; }
.streaming-indicator .typing-dots span:nth-child(2) { animation-delay: -0.16s; }
.streaming-indicator .typing-dots span:nth-child(3) { animation-delay: 0s; }

.streaming-text {
  font-size: 12px;
  color: #64748b;
  font-style: italic;
}

@keyframes streaming-pulse {
  0%, 80%, 100% {
    transform: scale(0.8);
    opacity: 0.5;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

/* AI建议区域 */
.ai-suggestions {
  margin-top: 12px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #f8fafc;
}

.suggestions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: linear-gradient(135deg, #f1f5f9 0%, #e2e8f0 100%);
  border-radius: 8px 8px 0 0;
  cursor: pointer;
  user-select: none;
  transition: all 0.2s;
  border-bottom: 1px solid #e2e8f0;
}

.suggestions-header:hover {
  background: linear-gradient(135deg, #e2e8f0 0%, #cbd5e1 100%);
  transform: translateY(-1px);
}

.suggestions-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.suggestions-icon {
  font-size: 16px;
}

.suggestions-text {
  font-weight: 600;
  font-size: 14px;
  color: #334155;
}

.suggestions-count {
  font-size: 12px;
  color: #64748b;
  background: #cbd5e1;
  padding: 2px 6px;
  border-radius: 8px;
}

.collapse-btn {
  background: none;
  border: none;
  font-size: 12px;
  color: #64748b;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 6px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 4px;
}

.collapse-btn:hover {
  background: rgba(99, 102, 241, 0.1);
  color: #334155;
  transform: translateY(-1px);
}

.collapse-btn.collapsed {
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.1);
}

.collapse-icon {
  font-weight: bold;
  transition: transform 0.2s;
}

.collapse-text {
  font-weight: 500;
}

.suggestions-content {
  padding: 16px;
}

/* 建议项 */
.suggestion-item {
  margin-bottom: 16px;
  padding: 16px;
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.suggestion-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, #667eea, transparent);
  transform: translateX(-100%);
  transition: transform 0.6s ease;
}

.suggestion-item:hover::before {
  transform: translateX(100%);
}

.suggestion-item:hover {
  border-color: #c7d2fe;
  box-shadow: 0 4px 20px rgba(99, 102, 241, 0.15);
  transform: translateY(-2px);
}

.suggestion-item:last-child {
  margin-bottom: 0;
}

.suggestion-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
  gap: 12px;
}

.suggestion-name {
  font-weight: 600;
  font-size: 15px;
  color: #1e293b;
  line-height: 1.4;
  flex: 1;
}

.suggestion-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.suggestion-description {
  font-size: 13px;
  color: #64748b;
  line-height: 1.6;
  margin-top: 8px;
}

/* 代码建议 */
.code-suggestion {
  border-left: 4px solid #10b981;
}

.code-suggestion .suggestion-name::before {
  content: "💻 ";
}

/* 文件建议 */
.file-suggestion {
  border-left: 4px solid #f59e0b;
}

/* 操作建议 */
.action-suggestion {
  border-left: 4px solid #6366f1;
}

/* 变量表格 */
.variables-table {
  margin: 20px 0;
  background: white;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  overflow: hidden;
}

.variables-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-bottom: 1px solid #e2e8f0;
}

.variables-title {
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: #1e293b;
}

.variables-icon {
  font-size: 18px;
}

.variables-count {
  font-size: 12px;
  color: #64748b;
  background: #e2e8f0;
  padding: 2px 8px;
  border-radius: 12px;
  font-weight: 500;
}

.variables-actions {
  display: flex;
  gap: 8px;
}

.variables-container {
  overflow-x: auto;
}

.variables-grid {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  min-width: 800px;
}

.variables-grid th {
  background: #f8fafc;
  padding: 12px 16px;
  text-align: left;
  font-weight: 600;
  color: #374151;
  border-bottom: 2px solid #e2e8f0;
  position: sticky;
  top: 0;
  z-index: 10;
}

.variables-grid td {
  padding: 16px;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}

.variable-row {
  transition: all 0.2s ease;
}

.variable-row:hover {
  background: #f8fafc;
}

.variable-row.high-priority {
  border-left: 4px solid #f59e0b;
}

/* 列宽设置 */
.var-name-col { width: 18%; }
.var-desc-col { width: 35%; }
.var-code-col { width: 25%; }
.var-priority-col { width: 12%; }
.var-actions-col { width: 10%; }

/* 变量名样式 */
.variable-name {
  display: flex;
  align-items: center;
}

.var-name-text {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-weight: 700;
  color: #7c3aed;
  font-size: 14px;
}

/* 描述样式 */
.variable-desc {
  color: #64748b;
  line-height: 1.5;
}

.var-desc-text {
  display: block;
  word-break: break-word;
}

/* 代码样式 */
.variable-code {
  background: #f1f5f9;
  border-radius: 4px;
  padding: 2px 0;
}

.var-code-text {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  background: #f1f5f9;
  color: #1e293b;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

/* 优先级徽章 */
.variable-priority {
  text-align: center;
}

.priority-badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
  text-align: center;
  min-width: 40px;
}

.priority-badge.high {
  background: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%);
  color: white;
}

.priority-badge.medium {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
}

.priority-badge.low {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}


/* 操作按钮组 */
.variable-actions {
  text-align: center;
}

.action-group {
  display: flex;
  gap: 6px;
  justify-content: center;
  align-items: center;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.action-btn:disabled:hover {
  transform: none;
  box-shadow: none;
}

/* 按钮样式 */
.action-btn {
  padding: 8px 16px;
  font-size: 12px;
  font-weight: 500;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  border: none;
  position: relative;
  overflow: hidden;
}

.action-btn:active {
  transform: scale(0.98);
}

.action-btn.mini {
  padding: 6px 12px;
  font-size: 11px;
  border-radius: 4px;
}

.action-btn.secondary {
  background: #f8fafc;
  color: #64748b;
  border: 1px solid #e2e8f0;
}

.action-btn.secondary:hover {
  background: #f1f5f9;
  color: #475569;
  border-color: #cbd5e1;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(100, 116, 139, 0.15);
}

.action-btn.primary {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: white;
  border: none;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
}

.action-btn.primary:hover {
  background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(59, 130, 246, 0.4);
}

/* 消息摘要 */
.message-summary {
  padding: 12px 16px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 8px;
}

.message-summary:hover {
  background: linear-gradient(135deg, #f1f5f9 0%, #e2e8f0 100%);
  border-color: #c7d2fe;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.1);
}

.summary-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-main {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.summary-icon {
  font-size: 16px;
}

.summary-text {
  font-size: 14px;
  color: #475569;
  font-weight: 500;
  line-height: 1.4;
}

.summary-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.suggestions-badge {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
}

.expand-icon {
  font-size: 12px;
  color: #64748b;
  transition: all 0.2s;
  font-weight: bold;
}

/* 快捷操作 */
.ai-quick-actions {
  padding: 24px 16px 12px 16px; /* 增加顶部边距为AI建议标签留空间 */
  border-top: 1px solid #e2e8f0;
}

.quick-btn {
  background: transparent;
  border: 1px solid #e2e8f0;
  color: #64748b;
  padding: 6px 12px;
  border-radius: 16px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.quick-btn:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
  color: #334155;
  transform: translateY(-1px);
}

/* 选中文本时的快捷按钮样式 */
.quick-btn.selection-action {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-color: transparent;
}

.quick-btn.selection-action:hover {
  background: linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

/* 输入区域 */
.ai-input-area {
  padding: 16px;
  border-top: 1px solid #e2e8f0;
  background: #fafafa;
}

.input-wrapper {
  display: flex;
  gap: 8px;
  align-items: flex-end;
}

.chat-input {
  flex: 1;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 12px 16px;
  resize: none;
  outline: none;
  transition: all 0.2s;
  max-height: 100px;
  font-family: inherit;
  font-size: 14px;
  line-height: 1.4;
  background: white;
}

.chat-input:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.send-button {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 50%;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.send-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.send-button:disabled {
  background: #94a3b8;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.send-button svg {
  width: 20px;
  height: 20px;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.input-hint {
  font-size: 11px;
  color: #94a3b8;
  margin-top: 8px;
  text-align: center;
}

/* 动画 */
@keyframes pulse {
  0% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.7;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

@keyframes typing {
  0%, 80%, 100% {
    transform: scale(0.8);
    opacity: 0.5;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 480px) {
  .ai-chat-panel {
    width: calc(100vw - 48px);
    right: 24px;
    left: 24px;
  }
  
  .ai-quick-actions {
    flex-direction: column;
  }
  
  .quick-btn {
    width: 100%;
    justify-content: center;
  }
}

/* 详情弹窗样式 */
.detail-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  animation: fadeIn 0.2s ease;
}

.detail-modal {
  background: white;
  border-radius: 12px;
  box-shadow: 0 24px 48px rgba(0, 0, 0, 0.2);
  max-width: 800px;
  max-height: 80vh;
  width: 90vw;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: slideIn 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideIn {
  from { 
    opacity: 0;
    transform: translateY(-20px) scale(0.95);
  }
  to { 
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.detail-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
}

.detail-modal-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  display: flex;
  align-items: center;
  gap: 8px;
}

.detail-icon {
  font-size: 20px;
}

.detail-modal-close {
  background: none;
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
}

.detail-modal-close:hover {
  background: rgba(0, 0, 0, 0.1);
  color: #374151;
}

.detail-modal-close svg {
  width: 18px;
  height: 18px;
}

.detail-modal-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.detail-section {
  margin-bottom: 24px;
}

.detail-section:last-child {
  margin-bottom: 0;
}

.detail-section-title {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.detail-description {
  margin: 0;
  color: #4b5563;
  line-height: 1.6;
  font-size: 15px;
}

.detail-meta {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.detail-tag {
  display: inline-flex;
  align-items: center;
  padding: 6px 12px;
  background: #f3f4f6;
  color: #374151;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
}

.detail-tag.priority.high {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  color: #92400e;
}

.detail-tag.priority.medium {
  background: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%);
  color: #1e40af;
}

.detail-tag.priority.low {
  background: linear-gradient(135deg, #d1fae5 0%, #a7f3d0 100%);
  color: #065f46;
}

.detail-code-editor {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
  background: #1f2937;
}

.code-editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #374151;
  border-bottom: 1px solid #4b5563;
}

.code-language {
  font-size: 12px;
  font-weight: 600;
  color: #9ca3af;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.copy-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  background: none;
  border: none;
  color: #9ca3af;
  font-size: 12px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.copy-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #d1d5db;
}

.copy-btn svg {
  width: 14px;
  height: 14px;
}

.code-content {
  margin: 0;
  padding: 16px;
  background: #1f2937;
  color: #e5e7eb;
  font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Roboto Mono', monospace;
  font-size: 13px;
  line-height: 1.5;
  overflow-x: auto;
  white-space: pre;
}

.code-content code {
  color: inherit;
  background: none;
  padding: 0;
  font-family: inherit;
}

.detail-modal-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid #e5e7eb;
  background: #f9fafb;
}

.detail-btn {
  padding: 10px 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
  min-width: 80px;
}

.detail-btn.secondary {
  background: #f3f4f6;
  color: #374151;
}

.detail-btn.secondary:hover {
  background: #e5e7eb;
}

.detail-btn.primary {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
}

.detail-btn.primary:hover {
  background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
}

.detail-btn.tertiary {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.detail-btn.tertiary:hover {
  background: linear-gradient(135deg, #059669 0%, #047857 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .detail-modal {
    width: 95vw;
    max-height: 85vh;
  }
  
  .detail-modal-header {
    padding: 16px 20px;
  }
  
  .detail-modal-content {
    padding: 20px;
  }
  
  .detail-modal-actions {
    padding: 16px 20px;
    flex-direction: column-reverse;
  }
  
  .detail-btn {
    width: 100%;
    justify-content: center;
  }
}

/* 代码解释样式 */
.code-explanation {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 20px;
  margin: 8px 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.explanation-section {
  margin: 16px 0;
  padding: 16px;
  border-radius: 8px;
  background: white;
  border-left: 4px solid #3b82f6;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}

.explanation-section.main-function {
  border-left-color: #10b981;
  background: linear-gradient(135deg, #ecfdf5 0%, #f0fdf4 100%);
}

.section-header {
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.section-content {
  color: #4b5563;
  line-height: 1.6;
  font-size: 15px;
}

.explanation-step {
  display: flex;
  align-items: flex-start;
  margin: 12px 0;
  padding: 12px;
  background: white;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
  transition: all 0.2s ease;
}

.explanation-step:hover {
  border-color: #3b82f6;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.1);
}

.step-number {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 12px;
  margin-right: 12px;
  margin-top: 2px;
}

.step-content {
  flex: 1;
  color: #374151;
  line-height: 1.6;
  font-size: 14px;
}

.tech-keyword {
  background: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%);
  color: #1e40af;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 500;
  font-size: 13px;
  border: 1px solid #93c5fd;
  white-space: nowrap;
}

.inline-code {
  background: #f1f5f9;
  color: #475569;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Roboto Mono', monospace;
  font-size: 13px;
  border: 1px solid #cbd5e1;
  font-weight: 500;
}

/* 响应式优化 */
@media (max-width: 768px) {
  .code-explanation {
    padding: 16px;
    margin: 6px 0;
  }
  
  .explanation-section {
    padding: 12px;
    margin: 12px 0;
  }
  
  .explanation-step {
    padding: 10px;
    flex-direction: column;
    align-items: flex-start;
  }
  
  .step-number {
    margin-bottom: 8px;
    margin-right: 0;
  }
  
  .tech-keyword {
    display: inline-block;
    margin: 2px;
  }
}

/* 快捷操作行布局 */
.action-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 8px;
}

.action-row:last-child {
  margin-bottom: 0;
}

.ai-suggestions-row {
  position: relative;
}

.ai-suggestions-row::before {
  content: "🤖 AI建议";
  position: absolute;
  top: -20px;
  left: 0;
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

/* 优先级样式修饰 - 微妙的边框和背景变化 */
.quick-btn.priority-high {
  border-color: #f59e0b;
  background: linear-gradient(135deg, #fefbf3 0%, #fef3c7 100%);
}

.quick-btn.priority-high:hover {
  border-color: #d97706;
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.2);
}

.quick-btn.priority-medium {
  border-color: #3b82f6;
  background: linear-gradient(135deg, #f8fafc 0%, #dbeafe 100%);
}

.quick-btn.priority-medium:hover {
  border-color: #2563eb;
  background: linear-gradient(135deg, #dbeafe 0%, #bfdbfe 100%);
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.2);
}

.quick-btn.priority-low {
  border-color: #10b981;
  background: linear-gradient(135deg, #f0fdf9 0%, #d1fae5 100%);
}

.quick-btn.priority-low:hover {
  border-color: #059669;
  background: linear-gradient(135deg, #d1fae5 0%, #a7f3d0 100%);
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.2);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .action-row {
    flex-direction: column;
  }
  
  .ai-suggestions-row::before {
    position: static;
    margin-bottom: 8px;
    display: block;
  }
}
</style>