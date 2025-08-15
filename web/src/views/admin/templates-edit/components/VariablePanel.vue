<template>
  <!-- 变量插入面板 -->
  <div v-show="isOpen" class="variable-panel" :style="{ height: panelHeight + 'px' }">
    <div class="variable-tabs">
      <div class="tab-header">
        <div 
          v-for="tab in variableTabs" 
          :key="tab.key"
          class="tab-item"
          :class="{ active: activeTab === tab.key }"
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
        </div>
      </div>
      
      <!-- 模板语法 Tab -->
      <div v-show="activeTab === 'syntax'" class="tab-content">
        <div class="function-categories">
          <div 
            v-for="category in templateSyntaxCategories" 
            :key="category.name"
            class="category-row"
          >
            <span class="category-label">{{ category.name }}</span>
            <div class="category-tags">
              <div 
                v-for="syntax in category.syntaxes" 
                :key="syntax.name"
                class="variable-tag syntax"
                @click="handleInsertSyntax(syntax)"
                @mouseenter="handleShowSyntaxDetail(syntax, $event)"
                @mouseleave="handleHideFunctionDetail"
              >
                {{ syntax.display_name || syntax.name }}
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 内置函数 Tab -->
      <div v-show="activeTab === 'functions'" class="tab-content">
        <div v-if="loadingFunctions" class="loading-state">
          <n-spin size="small" />
          <span style="margin-left: 8px;">加载函数中...</span>
        </div>
        <div v-else class="function-categories">
          <div 
            v-for="category in builtinFunctionCategories" 
            :key="category.name"
            class="category-row"
          >
            <span class="category-label">{{ category.name }}</span>
            <div class="category-tags">
              <div 
                v-for="func in category.functions" 
                :key="func.name"
                class="variable-tag function"
                @click="handleInsertFunction(formatFunction(func))"
                @mouseenter="handleShowFunctionDetail(func, $event)"
                @mouseleave="handleHideFunctionDetail"
              >
                {{ func.display_name || func.name }}
              </div>
            </div>
          </div>
          
          <!-- 如果没有数据显示提示 -->
          <div v-if="builtinFunctionCategories.length === 0" class="empty-state">
            <span>暂无可用的内置函数</span>
          </div>
        </div>
      </div>
      
      <!-- Sprig函数 Tab -->
      <div v-show="activeTab === 'sprig'" class="tab-content">
        <div v-if="loadingSprigFunctions" class="loading-state">
          <n-spin size="small" />
          <span style="margin-left: 8px;">加载Sprig函数中...</span>
        </div>
        <div v-else class="function-categories">
          <div 
            v-for="category in sprigFunctionCategories" 
            :key="category.name"
            class="category-row"
          >
            <span class="category-label">{{ category.name }}</span>
            <div class="category-tags">
              <div 
                v-for="func in category.functions" 
                :key="func.name"
                class="variable-tag function sprig"
                @click="handleInsertSprigFunction(func)"
                @mouseenter="handleShowSprigFunctionDetail(func, $event)"
                @mouseleave="handleHideFunctionDetail"
              >
                {{ func.display_name || func.name }}
              </div>
            </div>
          </div>
          
          <!-- 如果没有数据显示提示 -->
          <div v-if="sprigFunctionCategories.length === 0" class="empty-state">
            <span>暂无可用的Sprig函数</span>
          </div>
        </div>
      </div>

      <!-- 自定义函数详情面板 -->
      <div 
        v-if="functionDetailVisible && selectedFunction"
        class="function-detail-panel"
        :style="functionDetailStyle"
        @mouseenter="onDetailPanelEnter"
        @mouseleave="handleHideFunctionDetail"
      >
        <div class="detail-header">
          <div class="detail-title">
            <span class="function-icon">⚡</span>
            {{ selectedFunction.display_name || selectedFunction.name }}
          </div>
          <div class="detail-type">{{ selectedFunction.return_type || 'string' }}</div>
        </div>
        
        <div class="detail-body">
          <div class="detail-description">
            {{ selectedFunction.description }}
          </div>
          
          <!-- 如果有usage字段，显示使用说明 -->
          <div v-if="selectedFunction.usage" class="detail-section">
            <div class="section-label">
              <span class="section-icon">📖</span>
              使用说明
            </div>
            <div class="section-content">
              {{ selectedFunction.usage }}
            </div>
          </div>
          
          <!-- 参数信息 -->
          <div v-if="selectedFunction.params && selectedFunction.params.length > 0" class="detail-section">
            <div class="section-label">
              <span class="section-icon">🔧</span>
              参数列表
            </div>
            <div class="section-content">
              <div class="params-list">
                <div 
                  v-for="param in selectedFunction.params" 
                  :key="param.name"
                  class="param-item"
                >
                  <div class="param-header">
                    <span class="param-name">{{ param.name }}</span>
                    <span class="param-type">{{ param.type }}</span>
                    <span v-if="param.required" class="param-required">必需</span>
                    <span v-else class="param-optional">可选</span>
                  </div>
                  <div class="param-description">{{ param.description }}</div>
                  <div v-if="param.default" class="param-default">
                    默认值: <code>{{ param.default }}</code>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <div class="detail-section">
            <div class="section-label">
              <span class="section-icon">💡</span>
              使用示例
            </div>
            <div class="section-content code-content">
              {{ selectedFunction.example }}
            </div>
          </div>
          
          <div class="detail-section">
            <div class="section-label">
              <span class="section-icon">✨</span>
              点击插入
            </div>
            <div class="section-content code-content insert-preview">
              {{ selectedFunction.insert_text }}
            </div>
          </div>
        </div>
      </div>
      
      <!-- 内置变量 Tab -->
      <div v-show="activeTab === 'builtin'" class="tab-content">
        <div class="variable-section">
          <div class="section-title">文本变量</div>
          <div class="variable-tags">
            <div 
              v-for="variable in quickVariables" 
              :key="variable.name"
              class="variable-tag builtin"
              @click="handleInsertVariable(variable.name)"
              :title="`${variable.name} - ${variable.label}`"
            >
              {{ variable.label }}
            </div>
          </div>
        </div>
      </div>
      
      <!-- 用户变量 Tab -->
      <div v-show="activeTab === 'custom'" class="tab-content">
        <div v-if="textVariables.length > 0" class="variable-section">
          <div class="section-title">用户变量</div>
          <div class="variable-tags">
            <n-tag
              v-for="variable in textVariables"
              :key="variable.id"
              :class="['variable-tag', getVariableTagClass(variable.variableType)]"
              @click="handleInsertVariable(variable.name)"
              :title="`${variable.name} (${getVariableTypeLabel(variable.variableType)})${variable.description ? ' - ' + variable.description : ''}`"
            >
              {{ variable.name }}
              <span class="variable-type-badge">{{ getVariableTypeLabel(variable.variableType) }}</span>
            </n-tag>
          </div>
        </div>
        
        <div v-if="conditionalVariables.length > 0" class="variable-section">
          <div class="section-title">条件变量</div>
          <div class="variable-tags">
            <n-tag
              v-for="variable in conditionalVariables"
              :key="variable.id"
              class="variable-tag conditional"
              @click="handleInsertVariable(variable.name)"
              :title="`${variable.name}${variable.description ? ' - ' + variable.description : ''}`"
            >
              {{ variable.name }}
            </n-tag>
          </div>
        </div>
        
        <div v-if="templateVariables.length === 0" class="empty-variables">
          <div class="empty-text">暂无自定义变量</div>
          <n-button text type="primary" size="small" @click="$emit('show-variable-manager')">
            添加变量
          </n-button>
        </div>
      </div>

      <!-- 预设变量 Tab - 显示已订阅的预设变量 -->
      <div v-show="activeTab === 'presets'" class="tab-content preset-tab">
        <div class="preset-variables-summary">
          <div class="summary-header">
            <h4>已订阅的预设变量</h4>
            <n-button 
              type="primary" 
              size="small" 
              @click="showSubscribeModal = true"
            >
              <template #icon>
                <n-icon><AddOutline /></n-icon>
              </template>
              订阅
            </n-button>
          </div>
          
          <div class="subscribed-variables" v-loading="loadingPresets">
            <n-empty v-if="subscribedPresets.length === 0" description="暂无订阅的预设变量" size="small">
              <template #extra>
                <n-button size="small" @click="$emit('show-preset-manager')">立即订阅</n-button>
              </template>
            </n-empty>
            
            <div v-else class="preset-list">
              <div 
                v-for="preset in subscribedPresets" 
                :key="preset.id"
                class="preset-item"
              >
                <div class="preset-info">
                  <div class="preset-name">{{ preset.presetName }}</div>
                  <div class="preset-description" v-if="preset.description">{{ preset.description }}</div>
                </div>
                <div class="preset-actions">
                  <n-button 
                    size="tiny" 
                    quaternary 
                    type="error"
                    @click="unsubscribePreset(preset)"
                    :loading="preset.unsubscribing"
                  >
                    取消订阅
                  </n-button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 订阅预设变量弹窗 -->
    <n-modal v-model:show="showSubscribeModal" :mask-closable="false">
      <n-card style="width: 800px" title="订阅预设变量" :bordered="false" size="huge">
        <template #header-extra>
          <n-button quaternary circle @click="showSubscribeModal = false">
            <template #icon>
              <n-icon><CloseOutline /></n-icon>
            </template>
          </n-button>
        </template>

        <!-- 搜索栏 -->
        <div class="search-bar">
          <n-input 
            v-model:value="searchKeyword" 
            placeholder="搜索预设变量..." 
            clearable
            @update:value="searchPresets"
          >
            <template #prefix>
              <n-icon><SearchOutline /></n-icon>
            </template>
          </n-input>
        </div>

        <!-- 可用预设变量列表 -->
        <div class="available-presets" v-loading="presetsLoading">
          <n-empty v-if="availablePresets.length === 0" description="暂无可用的预设变量" />
          
          <div v-else class="presets-list">
            <div 
              class="preset-item-modal" 
              v-for="preset in availablePresets" 
              :key="preset.id"
            >
              <div class="preset-content">
                <n-checkbox 
                  :checked="isPresetSelected(preset.id)"
                  @update:checked="togglePreset(preset, $event)"
                >
                  <div class="preset-info-modal">
                    <div class="preset-name">{{ preset.name }}</div>
                    <div class="preset-description" v-if="preset.description">
                      {{ preset.description }}
                    </div>
                  </div>
                </n-checkbox>
              </div>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div class="pagination" v-if="totalPresets > pageSize">
          <n-pagination 
            v-model:page="currentPage"
            :page-size="pageSize"
            :item-count="totalPresets"
            @update:page="loadAvailablePresets"
            show-size-picker
            :page-sizes="[10, 20, 30]"
            @update:page-size="handlePageSizeChange"
          />
        </div>

        <template #footer>
          <div class="modal-footer">
            <n-button @click="showSubscribeModal = false">取消</n-button>
            <n-button 
              type="primary" 
              @click="confirmSubscribe" 
              :loading="subscribing"
              :disabled="selectedPresets.length === 0"
            >
              订阅选中的预设 ({{ selectedPresets.length }})
            </n-button>
          </div>
        </template>
      </n-card>
    </n-modal>
    
    <!-- 变量面板拖拽调整手柄 -->
    <div 
      class="variable-panel-resize-handle"
      @mousedown="startResize"
      :class="{ 'resizing': isResizing }"
    >
      <div class="resize-handle-dots">
        <div class="dot"></div>
        <div class="dot"></div>
        <div class="dot"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { NSpin, NButton, NTag, NEmpty, NSwitch, NIcon, NModal, NCard, NInput, NCheckbox, NPagination, useMessage } from 'naive-ui'
import { AddOutline, CloseOutline, SearchOutline } from '@vicons/ionicons5'
import request from '@/utils/request'

const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true
  },
  templateVariables: {
    type: Array,
    default: () => []
  },
  templateSyntaxCategories: {
    type: Array,
    default: () => []
  },
  builtinFunctionCategories: {
    type: Array,
    default: () => []
  },
  sprigFunctionCategories: {
    type: Array,
    default: () => []
  },
  loadingFunctions: {
    type: Boolean,
    default: false
  },
  loadingSprigFunctions: {
    type: Boolean,
    default: false
  },
  quickVariables: {
    type: Array,
    default: () => []
  },
  templateId: {
    type: [String, Number],
    required: true
  }
})

const emit = defineEmits([
  'insert-syntax',
  'insert-function', 
  'insert-sprig-function',
  'insert-variable',
  'show-variable-manager',
  'show-preset-manager',
  'update:height'
])

// 变量标签页配置
const variableTabs = [
  { key: 'syntax', label: '模板语法' },
  { key: 'functions', label: '内置函数' },
  { key: 'sprig', label: 'Sprig函数' },
  { key: 'builtin', label: '内置变量' },
  { key: 'custom', label: '用户变量' },
  { key: 'presets', label: '+ 预设变量' }
]

// 状态
const activeTab = ref('syntax')
const panelHeight = ref(300)
const isResizing = ref(false)

// 预设变量相关状态
const message = useMessage()
const subscribedPresets = ref([])
const loadingPresets = ref(false)

// 订阅弹窗相关状态
const showSubscribeModal = ref(false)
const availablePresets = ref([])
const presetsLoading = ref(false)
const selectedPresets = ref([])
const subscribing = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const totalPresets = ref(0)
const searchKeyword = ref('')

// 函数详情面板
const functionDetailVisible = ref(false)
const selectedFunction = ref(null)
const functionDetailStyle = ref({})
let hideTimer = null
let showTimer = null

// 面板大小调整
const minPanelHeight = 150
const maxPanelHeight = 600
let startY = 0
let startHeight = 0

// 计算属性：按类型分组变量
const textVariables = computed(() => {
  return props.templateVariables.filter(v => 
    v.variableType === 'text' || 
    v.variableType === 'string' || 
    v.variableType === '字符串' ||
    v.variableType === 'number' || 
    v.variableType === '数字' ||
    v.variableType === 'boolean' || 
    v.variableType === '布尔值' ||
    v.variableType === 'list' || 
    v.variableType === '列表' ||
    v.variableType === 'object' || 
    v.variableType === '对象' ||
    !v.variableType
  )
})

const conditionalVariables = computed(() => {
  return props.templateVariables.filter(v => v.variableType === 'conditional')
})

// 获取变量类型标签
const getVariableTypeLabel = (type) => {
  const typeLabels = {
    'string': '字符串',
    '字符串': '字符串',
    'text': '文本',
    '文本': '文本',
    'number': '数字',
    '数字': '数字',
    'boolean': '布尔值',
    '布尔值': '布尔值',
    'list': '列表',
    '列表': '列表',
    'object': '对象',
    '对象': '对象',
    'conditional': '条件'
  }
  return typeLabels[type] || type || '文本'
}

// 获取变量标签样式类
const getVariableTagClass = (type) => {
  const classMap = {
    'string': 'string',
    '字符串': 'string',
    'text': 'string',
    '文本': 'string',
    'number': 'number',
    '数字': 'number',
    'boolean': 'boolean',
    '布尔值': 'boolean',
    'list': 'list',
    '列表': 'list',
    'object': 'object',
    '对象': 'object',
    'conditional': 'conditional'
  }
  return classMap[type] || 'string'
}

// 格式化函数
const formatFunction = (func) => {
  if (func.insert_text) {
    return func.insert_text
  }
  
  if (func.params && func.params.length > 0) {
    const params = func.params.map(p => `${p.name}`).join(' ')
    return `{{${func.name} ${params}}}`
  }
  
  return `{{${func.name}}}`
}

// 事件处理
const handleInsertSyntax = (syntax) => {
  emit('insert-syntax', syntax)
}

const handleInsertFunction = (func) => {
  emit('insert-function', func)
}

const handleInsertSprigFunction = (func) => {
  emit('insert-sprig-function', func)
}

const handleInsertVariable = (variableName) => {
  emit('insert-variable', variableName)
}

// 函数详情显示
const handleShowSyntaxDetail = (syntax, event) => {
  showFunctionDetail(syntax, event)
}

const handleShowFunctionDetail = (func, event) => {
  showFunctionDetail(func, event)
}

const handleShowSprigFunctionDetail = (func, event) => {
  showFunctionDetail(func, event)
}

const showFunctionDetail = (func, event) => {
  clearTimeout(hideTimer)
  
  showTimer = setTimeout(() => {
    selectedFunction.value = func
    
    const rect = event.target.getBoundingClientRect()
    functionDetailStyle.value = {
      position: 'fixed',
      left: rect.right + 10 + 'px',
      top: rect.top + 'px',
      zIndex: 1000
    }
    
    functionDetailVisible.value = true
  }, 800)
}

const handleHideFunctionDetail = () => {
  clearTimeout(showTimer)
  
  hideTimer = setTimeout(() => {
    functionDetailVisible.value = false
    selectedFunction.value = null
  }, 200)
}

const onDetailPanelEnter = () => {
  clearTimeout(hideTimer)
}

// 面板大小调整
const startResize = (e) => {
  isResizing.value = true
  startY = e.clientY
  startHeight = panelHeight.value
  
  document.addEventListener('mousemove', onResize)
  document.addEventListener('mouseup', stopResize)
  
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'ns-resize'
}

const onResize = (e) => {
  if (!isResizing.value) return
  
  const deltaY = e.clientY - startY
  const newHeight = Math.max(minPanelHeight, Math.min(maxPanelHeight, startHeight + deltaY))
  
  panelHeight.value = newHeight
  emit('update:height', newHeight)
}

const stopResize = () => {
  isResizing.value = false
  
  document.removeEventListener('mousemove', onResize)
  document.removeEventListener('mouseup', stopResize)
  
  document.body.style.userSelect = ''
  document.body.style.cursor = ''
}

// 预设变量API调用
const getSubscribedPresets = (templateId) => {
  return request({
    url: `/api/v1/templates/${templateId}/preset-variables`,
    method: 'GET'
  })
}

const updatePresetMapping = (templateId, id, data) => {
  return request({
    url: `/api/v1/templates/${templateId}/preset-variables/${id}`,
    method: 'PUT',
    data: { template_id: templateId, id: id, ...data }
  })
}

// 加载已订阅的预设变量
const loadSubscribedPresets = async () => {
  if (!props.templateId) return
  
  loadingPresets.value = true
  try {
    const response = await getSubscribedPresets(props.templateId)
    subscribedPresets.value = response.data.data?.list || []
  } catch (error) {
    console.error('加载预设变量失败:', error)
    subscribedPresets.value = []
  } finally {
    loadingPresets.value = false
  }
}

// 取消订阅预设变量
const unsubscribePreset = async (preset) => {
  preset.unsubscribing = true
  try {
    const response = await request({
      url: `/api/v1/templates/${props.templateId}/preset-variables/${preset.id}`,
      method: 'DELETE'
    })
    message.success('取消订阅成功')
    await loadSubscribedPresets()
  } catch (error) {
    console.error('取消订阅失败:', error)
    message.error('取消订阅失败')
  } finally {
    preset.unsubscribing = false
  }
}

// 加载可用预设变量
const loadAvailablePresets = async (page = 1) => {
  presetsLoading.value = true
  try {
    const response = await request({
      url: '/api/v1/templates/preset-variables/available',
      method: 'GET',
      params: {
        pageNum: page,
        pageSize: pageSize.value,
        keyword: searchKeyword.value
      }
    })
    
    const data = response.data.data || {}
    availablePresets.value = data.list || []
    totalPresets.value = data.total || 0
    currentPage.value = data.pageNum || 1
  } catch (error) {
    console.error('加载可用预设变量失败:', error)
    message.error('加载可用预设变量失败')
  } finally {
    presetsLoading.value = false
  }
}

// 搜索预设变量
let searchTimeout = null
const searchPresets = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
    loadAvailablePresets()
  }, 300)
}

// 分页处理
const handlePageSizeChange = (newSize) => {
  pageSize.value = newSize
  currentPage.value = 1
  loadAvailablePresets()
}

// 变量选择相关
const isPresetSelected = (presetId) => {
  return selectedPresets.value.includes(presetId)
}

const togglePreset = (preset, checked) => {
  if (checked) {
    if (!selectedPresets.value.includes(preset.id)) {
      selectedPresets.value.push(preset.id)
    }
  } else {
    const index = selectedPresets.value.indexOf(preset.id)
    if (index > -1) {
      selectedPresets.value.splice(index, 1)
    }
  }
}

// 确认订阅
const confirmSubscribe = async () => {
  if (selectedPresets.value.length === 0) {
    message.warning('请选择要订阅的预设变量')
    return
  }

  subscribing.value = true
  try {
    await request({
      url: `/api/v1/templates/${props.templateId}/preset-variables/subscribe`,
      method: 'POST',
      data: {
        template_id: props.templateId,
        preset_ids: selectedPresets.value
      }
    })
    message.success('订阅成功')
    showSubscribeModal.value = false
    selectedPresets.value = []
    await loadSubscribedPresets()
  } catch (error) {
    console.error('订阅失败:', error)
    message.error('订阅失败')
  } finally {
    subscribing.value = false
  }
}

// 更新预设变量状态
const updatePresetStatus = async (preset) => {
  try {
    await updatePresetMapping(props.templateId, preset.id, {
      mapped_name: preset.mapped_name,
      is_active: preset.is_active ? 1 : 0,
      sort: preset.sort || 0
    })
    message.success('状态更新成功')
  } catch (error) {
    console.error('更新状态失败:', error)
    message.error('更新状态失败')
    // 恢复原状态
    preset.is_active = !preset.is_active
  }
}

// 监听模板ID变化
watch(() => props.templateId, () => {
  if (props.templateId) {
    loadSubscribedPresets()
  }
}, { immediate: true })

// 监听活动tab变化，如果切换到预设变量tab则刷新数据
watch(() => activeTab.value, (newTab) => {
  if (newTab === 'presets' && props.templateId) {
    loadSubscribedPresets()
  }
})

// 监听弹窗打开，加载可用预设变量
watch(() => showSubscribeModal.value, (show) => {
  if (show) {
    selectedPresets.value = []
    currentPage.value = 1
    searchKeyword.value = ''
    loadAvailablePresets()
  }
})

// 清理
onUnmounted(() => {
  clearTimeout(hideTimer)
  clearTimeout(showTimer)
  document.removeEventListener('mousemove', onResize)
  document.removeEventListener('mouseup', stopResize)
})
</script>

<style scoped>
/* 变量面板样式 */
.variable-panel {
  background: #fff;
  border-bottom: 1px solid #e0e0e0;
  box-shadow: 0 2px 8px rgba(0,0,0,0.03);
  overflow: hidden;
  position: relative;
}

/* 变量面板拖拽调整样式 */
.variable-panel-resize-handle {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 12px;
  background: transparent;
  cursor: ns-resize;
  display: flex;
  align-items: center;
  justify-content: center;
}

.variable-panel-resize-handle:hover {
  background: rgba(24, 160, 88, 0.05);
}

.variable-panel-resize-handle.resizing {
  background: rgba(24, 160, 88, 0.1);
}

.resize-handle-dots {
  display: flex;
  gap: 3px;
}

.resize-handle-dots .dot {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: #999;
}

.variable-panel-resize-handle:hover .dot {
  background: #18a058;
}

/* 标签页样式 */
.variable-tabs {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.tab-header {
  display: flex;
  background: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  padding: 0 16px;
}

.tab-item {
  padding: 12px 16px;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
  user-select: none;
}

.tab-item:hover {
  color: #18a058;
  background: rgba(24, 160, 88, 0.05);
}

.tab-item.active {
  color: #18a058;
  border-bottom-color: #18a058;
  background: rgba(24, 160, 88, 0.05);
  font-weight: 500;
}

.tab-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

/* 分类和标签样式 */
.function-categories {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.category-row {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.category-label {
  min-width: 80px;
  font-size: 13px;
  color: #666;
  font-weight: 500;
  padding: 4px 0;
  flex-shrink: 0;
}

.category-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  flex: 1;
}

.variable-tag {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;
  border: 1px solid transparent;
  white-space: nowrap;
}

.variable-tag.syntax {
  background: #e8f4fd;
  color: #1890ff;
  border-color: #d4eafc;
}

.variable-tag.syntax:hover {
  background: #bae7ff;
  border-color: #91d5ff;
}

.variable-tag.function {
  background: #f6ffed;
  color: #52c41a;
  border-color: #d9f7be;
}

.variable-tag.function:hover {
  background: #d9f7be;
  border-color: #95de64;
}

.variable-tag.function.sprig {
  background: #fff7e6;
  color: #fa8c16;
  border-color: #ffd591;
}

.variable-tag.function.sprig:hover {
  background: #ffd591;
  border-color: #ffb347;
}

.variable-tag.builtin {
  background: #f0f5ff;
  color: #2f54eb;
  border-color: #d6e4ff;
}

.variable-tag.builtin:hover {
  background: #d6e4ff;
  border-color: #adc6ff;
}

.variable-tag.string {
  background: #e8f4fd;
  color: #1890ff;
  border-color: #d4eafc;
}

.variable-tag.number {
  background: #f6ffed;
  color: #52c41a;
  border-color: #d9f7be;
}

.variable-tag.boolean {
  background: #fff7e6;
  color: #fa8c16;
  border-color: #ffd591;
}

.variable-tag.list {
  background: #f9f0ff;
  color: #722ed1;
  border-color: #efdbff;
}

.variable-tag.object {
  background: #fff0f6;
  color: #eb2f96;
  border-color: #ffd6e7;
}

.variable-tag.conditional {
  background: #fff2e8;
  color: #fa541c;
  border-color: #ffd8bf;
}

.variable-type-badge {
  margin-left: 4px;
  font-size: 10px;
  opacity: 0.7;
}

/* 变量区域样式 */
.variable-section {
  margin-bottom: 20px;
}

.section-title {
  font-size: 13px;
  color: #333;
  font-weight: 600;
  margin-bottom: 8px;
  padding-bottom: 4px;
  border-bottom: 1px solid #f0f0f0;
}

.variable-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

/* 空状态样式 */
.empty-state,
.empty-variables {
  text-align: center;
  padding: 40px 20px;
  color: #999;
}

.empty-text {
  margin-bottom: 12px;
  font-size: 14px;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: #666;
}

/* 函数详情面板样式 */
.function-detail-panel {
  background: #fff;
  border: 1px solid #d9d9d9;
  border-radius: 8px;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
  max-width: 400px;
  min-width: 300px;
  max-height: 500px;
  overflow-y: auto;
  z-index: 1000;
}

.detail-header {
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  background: #fafafa;
}

.detail-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.function-icon {
  font-size: 18px;
}

.detail-type {
  font-size: 12px;
  color: #666;
  background: #f0f0f0;
  padding: 2px 6px;
  border-radius: 4px;
  display: inline-block;
}

.detail-body {
  padding: 16px;
}

.detail-description {
  font-size: 14px;
  color: #333;
  line-height: 1.5;
  margin-bottom: 16px;
}

.detail-section {
  margin-bottom: 16px;
}

.section-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.section-icon {
  font-size: 14px;
}

.section-content {
  font-size: 13px;
  color: #666;
  line-height: 1.5;
}

.section-content.code-content {
  background: #f8f8f8;
  padding: 8px 12px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  border: 1px solid #e8e8e8;
}

.section-content.insert-preview {
  background: #f6ffed;
  border-color: #d9f7be;
  color: #52c41a;
}

.params-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.param-item {
  padding: 8px;
  background: #fafafa;
  border-radius: 4px;
  border: 1px solid #f0f0f0;
}

.param-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.param-name {
  font-weight: 600;
  color: #333;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.param-type {
  font-size: 11px;
  color: #666;
  background: #f0f0f0;
  padding: 1px 4px;
  border-radius: 2px;
}

.param-required {
  font-size: 10px;
  color: #ff4d4f;
  background: #fff2f0;
  padding: 1px 4px;
  border-radius: 2px;
}

.param-optional {
  font-size: 10px;
  color: #52c41a;
  background: #f6ffed;
  padding: 1px 4px;
  border-radius: 2px;
}

.param-description {
  font-size: 12px;
  color: #666;
  line-height: 1.4;
}

.param-default {
  font-size: 12px;
  color: #666;
  margin-top: 4px;
}

.param-default code {
  background: #f0f0f0;
  padding: 1px 4px;
  border-radius: 2px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

/* 预设变量Tab特殊样式 */
.tab-item:has-text("+ 预设变量"), 
.tab-item:last-child {
  position: relative;
}

.tab-item:has-text("+ 预设变量"):before,
.tab-item:last-child:before {
  content: "+";
  margin-right: 4px;
  font-weight: bold;
  color: #1890ff;
}

.preset-tab {
  padding: 0;
  height: 100%;
}

/* 预设变量概览样式 */
.preset-variables-summary {
  padding: 16px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.preset-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.preset-item {
  display: flex;
  flex-direction: column;
  width: calc(20% - 6.4px);
  min-width: 120px;
  padding: 8px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  background: #fff;
  transition: all 0.2s;
  position: relative;
}

.preset-item:hover {
  border-color: #d9d9d9;
  background: #fafafa;
}

.preset-info {
  flex: 1;
  min-width: 0;
  margin-bottom: 8px;
}

.preset-name {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.preset-description {
  color: #666;
  font-size: 11px;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.preset-actions {
  display: flex;
  justify-content: flex-end;
}

.summary-header {
  display: flex;
  justify-content: space-between;
}

/* 订阅弹窗样式 */
.search-bar {
  margin-bottom: 16px;
}

.available-presets {
  max-height: 400px;
  overflow-y: auto;
}

.preset-item-modal {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
}

.preset-info-modal .preset-name {
  margin: 0 0 8px 0;
  color: #333;
  font-weight: 500;
}

.preset-info-modal .preset-description {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.pagination {
  margin-top: 16px;
  text-align: center;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}

.summary-header h4 {
  margin: 0;
  font-size: 14px;
  color: #333;
  font-weight: 600;
}

.subscribed-variables {
  flex: 1;
  overflow-y: auto;
}

.preset-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.preset-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.preset-item:hover {
  background: #e9ecef;
  border-color: #18a058;
}

.preset-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.preset-name {
  font-size: 13px;
  font-weight: 500;
  color: #333;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.preset-path {
  font-size: 11px;
  color: #666;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

</style>