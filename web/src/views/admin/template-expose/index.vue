<template>
  <div class="template-expose-container">
    <!-- 页头 -->
    <div class="page-header">
      <div class="header-left">
        <n-button quaternary circle @click="$router.back()">
          <template #icon>
            <n-icon><ArrowBackOutline /></n-icon>
          </template>
        </n-button>
        <div class="header-info">
          <h2 class="template-title">{{ templateInfo?.name || '模板变量定义' }}</h2>
          <p class="template-desc">定义模板变量的数据结构和类型信息</p>
        </div>
      </div>
      <div class="header-right">
        <n-button type="primary" @click="saveExpose" :loading="saving">
          <template #icon>
            <n-icon><SaveOutline /></n-icon>
          </template>
          保存变量定义
        </n-button>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="main-content">
      <!-- 左侧：模板文件树 -->
      <div class="left-panel" :style="{ width: leftPanelWidth + 'px' }">
        <div class="panel-header">
          <h3>模板文件</h3>
          <n-tooltip>
            <template #trigger>
              <n-icon class="help-icon"><HelpCircleOutline /></n-icon>
            </template>
            选择模板文件查看其中的变量使用情况
          </n-tooltip>
        </div>
        <TemplateFileExplorer 
          v-model:selectedFile="selectedFile"
          :template-id="templateId"
          @file-content-change="onFileContentChange"
        />
      </div>

      <!-- 左侧调整手柄 -->
      <div class="resize-handle left-resize" @mousedown="startLeftResize"></div>

      <!-- 中间：变量定义编辑 -->
      <div class="center-panel" :style="{ width: centerPanelWidth + 'px' }">
        <div class="panel-header">
          <h3>变量定义</h3>
          <div class="header-actions">
            <n-button size="small" quaternary @click="addVariable">
              <template #icon>
                <n-icon><AddOutline /></n-icon>
              </template>
              添加变量
            </n-button>
          </div>
        </div>
        <VariableDefinePanel
          v-model:variables="variableDefinitions"
          v-model:selectedVariable="selectedVariable"
          :template-id="templateId"
          :current-file="selectedFile"
          @variable-change="onVariableChange"
        />
      </div>

      <!-- 右侧调整手柄 -->
      <div class="resize-handle right-resize" @mousedown="startRightResize"></div>

      <!-- 右侧：预览和代码生成 -->
      <div class="right-panel">
        <div class="panel-header">
          <h3>预览</h3>
          <n-tabs v-model:value="previewTab" size="small">
            <n-tab-pane name="schema" tab="Schema">
              <SchemaPreview :variables="variableDefinitions" />
            </n-tab-pane>
            <n-tab-pane name="usage" tab="使用示例">
              <UsagePreview 
                :variables="variableDefinitions" 
                :selected-file="selectedFile"
                :file-content="fileContent"
              />
            </n-tab-pane>
            <n-tab-pane name="types" tab="类型定义">
              <TypesPreview :variables="variableDefinitions" />
            </n-tab-pane>
          </n-tabs>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { 
  NButton, NIcon, NTooltip, NTabs, NTabPane
} from 'naive-ui'
import { 
  ArrowBackOutline, SaveOutline, HelpCircleOutline, AddOutline 
} from '@vicons/ionicons5'

import TemplateFileExplorer from './components/TemplateFileExplorer.vue'
import VariableDefinePanel from './components/VariableDefinePanel.vue'
import SchemaPreview from './components/SchemaPreview.vue'
import UsagePreview from './components/UsagePreview.vue'
import TypesPreview from './components/TypesPreview.vue'

import request from '@/utils/request'

const route = useRoute()
const router = useRouter()
const message = useMessage()

// 路由参数
const templateId = computed(() => route.params.id)

// 页面状态
const templateInfo = ref(null)
const saving = ref(false)
const previewTab = ref('schema')

// 文件选择
const selectedFile = ref(null)
const fileContent = ref('')

// 变量定义
const variableDefinitions = ref([])
const selectedVariable = ref(null)

// 面板宽度控制
const leftPanelWidth = ref(280)
const centerPanelWidth = ref(400)
const minPanelWidth = 200
const maxPanelWidth = 600

// 拖拽状态
const isLeftResizing = ref(false)
const isRightResizing = ref(false)
const startX = ref(0)
const startLeftWidth = ref(0)
const startCenterWidth = ref(0)

// 加载模板信息
const loadTemplateInfo = async () => {
  try {
    const response = await request({
      url: `/api/v1/templates/${templateId.value}`,
      method: 'GET'
    })
    templateInfo.value = response.data.data
  } catch (error) {
    console.error('加载模板信息失败:', error)
    message.error('加载模板信息失败')
  }
}

// 加载变量定义
const loadVariableDefinitions = async () => {
  try {
    const response = await request({
      url: `/api/v1/templates/${templateId.value}/expose`,
      method: 'GET'
    })
    variableDefinitions.value = response.data.data?.variables || []
  } catch (error) {
    console.error('加载变量定义失败:', error)
    // 初始化为空数组，允许用户从零开始创建
    variableDefinitions.value = []
  }
}

// 保存变量定义
const saveExpose = async () => {
  saving.value = true
  try {
    await request({
      url: `/api/v1/templates/${templateId.value}/expose`,
      method: 'POST',
      data: {
        variables: variableDefinitions.value
      }
    })
    message.success('变量定义保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    message.error('保存变量定义失败')
  } finally {
    saving.value = false
  }
}

// 添加新变量
const addVariable = () => {
  const newVariable = {
    id: Date.now().toString(),
    name: '',
    displayName: '',
    description: '',
    type: 'string',
    required: false,
    defaultValue: '',
    example: '',
    validation: {
      minLength: null,
      maxLength: null,
      pattern: '',
      enum: []
    },
    children: []
  }
  variableDefinitions.value.push(newVariable)
  selectedVariable.value = newVariable
}

// 文件内容变化处理
const onFileContentChange = (content) => {
  fileContent.value = content
}

// 变量变化处理
const onVariableChange = (variable) => {
  // 实时预览更新等逻辑
}

// 左侧面板拖拽
const startLeftResize = (e) => {
  isLeftResizing.value = true
  startX.value = e.clientX
  startLeftWidth.value = leftPanelWidth.value
  document.addEventListener('mousemove', onLeftResize)
  document.addEventListener('mouseup', stopLeftResize)
}

const onLeftResize = (e) => {
  if (!isLeftResizing.value) return
  const deltaX = e.clientX - startX.value
  const newWidth = Math.min(Math.max(startLeftWidth.value + deltaX, minPanelWidth), maxPanelWidth)
  leftPanelWidth.value = newWidth
}

const stopLeftResize = () => {
  isLeftResizing.value = false
  document.removeEventListener('mousemove', onLeftResize)
  document.removeEventListener('mouseup', stopLeftResize)
}

// 右侧面板拖拽
const startRightResize = (e) => {
  isRightResizing.value = true
  startX.value = e.clientX
  startCenterWidth.value = centerPanelWidth.value
  document.addEventListener('mousemove', onRightResize)
  document.addEventListener('mouseup', stopRightResize)
}

const onRightResize = (e) => {
  if (!isRightResizing.value) return
  const deltaX = e.clientX - startX.value
  const newWidth = Math.min(Math.max(startCenterWidth.value + deltaX, minPanelWidth), maxPanelWidth)
  centerPanelWidth.value = newWidth
}

const stopRightResize = () => {
  isRightResizing.value = false
  document.removeEventListener('mousemove', onRightResize)
  document.removeEventListener('mouseup', stopRightResize)
}

// 组件挂载
onMounted(async () => {
  await loadTemplateInfo()
  await loadVariableDefinitions()
})

// 组件卸载
onUnmounted(() => {
  document.removeEventListener('mousemove', onLeftResize)
  document.removeEventListener('mouseup', stopLeftResize)
  document.removeEventListener('mousemove', onRightResize)
  document.removeEventListener('mouseup', stopRightResize)
})
</script>

<style scoped>
.template-expose-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f5f5;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: white;
  border-bottom: 1px solid #e8e8e8;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.template-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

.template-desc {
  margin: 0;
  font-size: 14px;
  color: #666;
}

.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.left-panel,
.center-panel,
.right-panel {
  background: white;
  border-right: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
}

.right-panel {
  flex: 1;
  border-right: none;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  background: #fafafa;
}

.panel-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.help-icon {
  color: #999;
  font-size: 16px;
  cursor: help;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.resize-handle {
  width: 4px;
  background: transparent;
  cursor: col-resize;
  position: relative;
  z-index: 10;
}

.resize-handle:hover {
  background: #1890ff;
}

.resize-handle.left-resize,
.resize-handle.right-resize {
  background: transparent;
}
</style>