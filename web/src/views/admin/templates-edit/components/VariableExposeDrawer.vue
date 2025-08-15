<template>
  <n-drawer v-model:show="visible" :width="'90%'" placement="bottom" :height="drawerHeight + 'px'">
    <n-drawer-content title="模板变量定义" :native-scrollbar="false">
      <template #header-extra>
        <n-space>
          <n-button size="small" @click="saveVariables" type="primary" :loading="saving">
            <template #icon>
              <n-icon><SaveOutline /></n-icon>
            </template>
            保存定义
          </n-button>
          <n-button size="small" quaternary @click="visible = false">
            <template #icon>
              <n-icon><CloseOutline /></n-icon>
            </template>
          </n-button>
        </n-space>
      </template>

      <!-- 拖拽手柄 -->
      <div class="resize-handle" @mousedown="startResize">
        <div class="handle-bar"></div>
      </div>

      <div class="expose-layout">
        <!-- 左侧：文件列表 -->
        <div class="left-panel">
          <div class="panel-title">模板文件</div>
          <div class="file-list">
            <div 
              v-for="file in templateFiles" 
              :key="file.id"
              class="file-item"
              :class="{ active: selectedFile?.id === file.id }"
              @click="selectFile(file)"
            >
              <n-icon class="file-icon">
                <DocumentOutline />
              </n-icon>
              {{ file.name }}
            </div>
          </div>
          
          <!-- 变量使用情况 -->
          <div v-if="selectedFile && variableUsages.length > 0" class="usage-section">
            <div class="section-title">变量使用 ({{ variableUsages.length }})</div>
            <div class="usage-list">
              <div v-for="usage in variableUsages" :key="usage.name" class="usage-item">
                <span class="var-name">{{ usage.name }}</span>
                <span class="var-count">{{ usage.count }}次</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 中间：变量定义 -->
        <div class="center-panel">
          <div class="panel-title">
            变量定义
            <n-button size="small" @click="addVariable">
              <template #icon>
                <n-icon><AddOutline /></n-icon>
              </template>
              添加
            </n-button>
          </div>
          
          <div class="variables-list">
            <div 
              v-for="(variable, index) in variableDefinitions" 
              :key="variable.id"
              class="variable-card"
              :class="{ active: selectedVariable?.id === variable.id }"
              @click="selectVariable(variable)"
            >
              <div class="variable-header">
                <div class="variable-info">
                  <span class="variable-name">{{ variable.name || '未命名' }}</span>
                  <span class="variable-type">{{ getTypeLabel(variable.type) }}</span>
                </div>
                <n-button size="tiny" quaternary type="error" @click.stop="removeVariable(index)">
                  <template #icon>
                    <n-icon><TrashOutline /></n-icon>
                  </template>
                </n-button>
              </div>
              <div class="variable-desc">{{ variable.description || '暂无描述' }}</div>
            </div>
          </div>

          <!-- 变量编辑表单 -->
          <div v-if="selectedVariable" class="variable-form">
            <n-form ref="formRef" :model="selectedVariable" size="small">
              <n-grid :cols="2" :x-gap="12" :y-gap="8">
                <n-grid-item>
                  <n-form-item label="变量名" path="name">
                    <n-input v-model:value="selectedVariable.name" placeholder="例如：userName" />
                  </n-form-item>
                </n-grid-item>
                <n-grid-item>
                  <n-form-item label="显示名">
                    <n-input v-model:value="selectedVariable.displayName" placeholder="例如：用户名称" />
                  </n-form-item>
                </n-grid-item>
                <n-grid-item>
                  <n-form-item label="数据类型">
                    <n-select v-model:value="selectedVariable.type" :options="typeOptions" />
                  </n-form-item>
                </n-grid-item>
                <n-grid-item>
                  <n-form-item label="必填">
                    <n-switch v-model:value="selectedVariable.required" />
                  </n-form-item>
                </n-grid-item>
                <n-grid-item :span="2">
                  <n-form-item label="描述">
                    <n-input v-model:value="selectedVariable.description" type="textarea" :rows="2" />
                  </n-form-item>
                </n-grid-item>
                <n-grid-item>
                  <n-form-item label="示例">
                    <n-input v-model:value="selectedVariable.example" />
                  </n-form-item>
                </n-grid-item>
                <n-grid-item>
                  <n-form-item label="默认值">
                    <n-input v-model:value="selectedVariable.defaultValue" />
                  </n-form-item>
                </n-grid-item>
              </n-grid>
            </n-form>
          </div>
        </div>

        <!-- 右侧：预览 -->
        <div class="right-panel">
          <div class="panel-title">预览</div>
          <n-tabs size="small">
            <n-tab-pane name="schema" tab="JSON Schema">
              <div class="preview-content">
                <pre class="json-preview">{{ jsonSchema }}</pre>
              </div>
            </n-tab-pane>
            <n-tab-pane name="usage" tab="使用示例">
              <div class="preview-content">
                <div class="usage-examples">
                  <div v-for="variable in variableDefinitions" :key="variable.id" class="usage-example">
                    <div class="example-title">{{ variable.displayName || variable.name }}</div>
                    <pre class="example-code">{{ getVariableTemplate(variable.name) }}</pre>
                  </div>
                </div>
              </div>
            </n-tab-pane>
          </n-tabs>
        </div>
      </div>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { 
  NDrawer, NDrawerContent, NButton, NIcon, NSpace, NForm, NFormItem, 
  NGrid, NGridItem, NInput, NSelect, NSwitch, NTabs, NTabPane, useMessage 
} from 'naive-ui'
import { 
  SaveOutline, CloseOutline, DocumentOutline, AddOutline, TrashOutline 
} from '@vicons/ionicons5'
import request from '@/utils/request'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  templateId: {
    type: [String, Number],
    required: true
  },
  templateVariables: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:show'])

const message = useMessage()

// 抽屉显示状态
const visible = computed({
  get: () => props.show,
  set: (value) => emit('update:show', value)
})

// 数据状态
const saving = ref(false)
const templateFiles = ref([])
const selectedFile = ref(null)
const variableUsages = ref([])
const variableDefinitions = ref([])
const selectedVariable = ref(null)

// 拖拽相关状态
const drawerHeight = ref(Math.floor(window.innerHeight * 0.67)) // 默认三分之二高度
const isResizing = ref(false)
const minHeight = 300
const maxHeight = window.innerHeight - 100

// 数据类型选项
const typeOptions = [
  { label: '字符串', value: 'string' },
  { label: '数字', value: 'number' },
  { label: '布尔值', value: 'boolean' },
  { label: '对象', value: 'object' },
  { label: '数组', value: 'array' }
]

// 获取类型标签
const getTypeLabel = (type) => {
  const typeMap = {
    string: '字符串',
    number: '数字',
    boolean: '布尔',
    object: '对象',
    array: '数组'
  }
  return typeMap[type] || type
}

// 获取变量模板语法
const getVariableTemplate = (varName) => {
  return `{{.${varName || 'varName'}}}`
}

// JSON Schema 预览
const jsonSchema = computed(() => {
  const schema = {
    type: 'object',
    properties: {},
    required: []
  }
  
  variableDefinitions.value.forEach(variable => {
    if (variable.name) {
      schema.properties[variable.name] = {
        type: variable.type,
        title: variable.displayName || variable.name,
        description: variable.description || ''
      }
      if (variable.required) {
        schema.required.push(variable.name)
      }
    }
  })
  
  return JSON.stringify(schema, null, 2)
})

// 加载模板文件
const loadTemplateFiles = async () => {
  try {
    const response = await request({
      url: `/api/v1/templates/${props.templateId}/files`,
      method: 'GET'
    })
    templateFiles.value = flattenFiles(response.data.data || [])
  } catch (error) {
    console.error('加载文件失败:', error)
  }
}

// 扁平化文件树
const flattenFiles = (files) => {
  const flattened = []
  const flatten = (items, path = '') => {
    items.forEach(item => {
      if (item.type === 'file') {
        flattened.push({
          ...item,
          path: path ? `${path}/${item.name}` : item.name
        })
      } else if (item.children) {
        flatten(item.children, path ? `${path}/${item.name}` : item.name)
      }
    })
  }
  flatten(files)
  return flattened
}

// 选择文件
const selectFile = async (file) => {
  selectedFile.value = file
  try {
    const response = await request({
      url: `/api/v1/templates/${props.templateId}/files/${file.id}/content`,
      method: 'GET'
    })
    analyzeVariableUsage(response.data.data.content || '')
  } catch (error) {
    console.error('加载文件内容失败:', error)
  }
}

// 分析变量使用情况
const analyzeVariableUsage = (content) => {
  const regex = /\{\{\s*\.(\w+(?:\.\w+)*)\s*\}\}/g
  const usages = new Map()
  
  let match
  while ((match = regex.exec(content)) !== null) {
    const varName = match[1]
    if (!usages.has(varName)) {
      usages.set(varName, { name: varName, count: 0 })
    }
    usages.get(varName).count++
  }
  
  variableUsages.value = Array.from(usages.values())
}

// 选择变量
const selectVariable = (variable) => {
  selectedVariable.value = variable
}

// 添加变量
const addVariable = () => {
  const newVariable = {
    id: Date.now().toString(),
    name: '',
    displayName: '',
    description: '',
    type: 'string',
    required: false,
    example: '',
    defaultValue: ''
  }
  variableDefinitions.value.push(newVariable)
  selectedVariable.value = newVariable
}

// 删除变量
const removeVariable = (index) => {
  variableDefinitions.value.splice(index, 1)
  if (selectedVariable.value === variableDefinitions.value[index]) {
    selectedVariable.value = null
  }
}

// 保存变量定义
const saveVariables = async () => {
  saving.value = true
  try {
    await request({
      url: `/api/v1/templates/${props.templateId}/expose`,
      method: 'POST',
      data: {
        variables: variableDefinitions.value
      }
    })
    message.success('变量定义保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    message.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 拖拽调整高度
const startResize = (e) => {
  isResizing.value = true
  const startY = e.clientY
  const startHeight = drawerHeight.value

  const onMouseMove = (e) => {
    if (!isResizing.value) return
    
    const deltaY = startY - e.clientY // 向上为正，向下为负
    const newHeight = Math.min(Math.max(startHeight + deltaY, minHeight), maxHeight)
    drawerHeight.value = newHeight
  }

  const onMouseUp = () => {
    isResizing.value = false
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    document.body.style.userSelect = ''
    document.body.style.cursor = ''
  }

  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'ns-resize'
}

// 监听抽屉显示状态
watch(visible, (show) => {
  if (show) {
    loadTemplateFiles()
  }
})
</script>

<style scoped>
.resize-handle {
  position: relative;
  height: 12px;
  cursor: ns-resize;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 8px;
}

.resize-handle:hover {
  background-color: rgba(24, 144, 255, 0.1);
}

.handle-bar {
  width: 60px;
  height: 4px;
  background: #d9d9d9;
  border-radius: 2px;
  transition: background-color 0.2s;
}

.resize-handle:hover .handle-bar {
  background: #1890ff;
}

.expose-layout {
  display: flex;
  height: calc(100% - 20px);
  gap: 16px;
}

.left-panel,
.center-panel,
.right-panel {
  background: #f8f8f8;
  border-radius: 6px;
  padding: 16px;
  overflow-y: auto;
}

.left-panel {
  width: 200px;
  flex-shrink: 0;
}

.center-panel {
  flex: 1;
  min-width: 300px;
}

.right-panel {
  width: 300px;
  flex-shrink: 0;
}

.panel-title {
  font-weight: 600;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.file-list {
  margin-bottom: 16px;
}

.file-item {
  display: flex;
  align-items: center;
  padding: 8px;
  margin-bottom: 4px;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.file-item:hover {
  background: #e6f7ff;
}

.file-item.active {
  background: #1890ff;
  color: white;
}

.file-icon {
  margin-right: 8px;
}

.usage-section {
  border-top: 1px solid #e0e0e0;
  padding-top: 12px;
}

.section-title {
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #666;
}

.usage-list {
  max-height: 150px;
  overflow-y: auto;
}

.usage-item {
  display: flex;
  justify-content: space-between;
  padding: 4px 8px;
  margin-bottom: 2px;
  background: white;
  border-radius: 3px;
  font-size: 12px;
}

.var-name {
  color: #1890ff;
  font-family: monospace;
}

.var-count {
  color: #666;
}

.variables-list {
  margin-bottom: 16px;
  max-height: 200px;
  overflow-y: auto;
}

.variable-card {
  background: white;
  border-radius: 4px;
  padding: 12px;
  margin-bottom: 8px;
  cursor: pointer;
  border: 1px solid #e0e0e0;
  transition: all 0.2s;
}

.variable-card:hover {
  border-color: #1890ff;
}

.variable-card.active {
  border-color: #1890ff;
  background: #f0f9ff;
}

.variable-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.variable-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.variable-name {
  font-weight: 600;
}

.variable-type {
  font-size: 12px;
  padding: 2px 6px;
  background: #f0f0f0;
  border-radius: 3px;
  color: #666;
}

.variable-desc {
  font-size: 12px;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.variable-form {
  border-top: 1px solid #e0e0e0;
  padding-top: 16px;
}

.preview-content {
  height: 400px;
  overflow-y: auto;
}

.json-preview {
  font-family: monospace;
  font-size: 12px;
  background: white;
  padding: 12px;
  border-radius: 4px;
  border: 1px solid #e0e0e0;
  margin: 0;
  white-space: pre-wrap;
}

.usage-examples {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.usage-example {
  background: white;
  padding: 12px;
  border-radius: 4px;
  border: 1px solid #e0e0e0;
}

.example-title {
  font-weight: 600;
  margin-bottom: 8px;
  color: #333;
}

.example-code {
  font-family: monospace;
  font-size: 12px;
  background: #f8f8f8;
  padding: 8px;
  border-radius: 3px;
  margin: 0;
}
</style>