<template>
  <div class="template-file-explorer">
    <div class="file-tree">
      <n-tree
        :data="treeData"
        :selected-keys="selectedKeys"
        :expanded-keys="expandedKeys"
        key-field="id"
        label-field="name"
        children-field="children"
        block-line
        @update:selected-keys="onSelectFile"
        @update:expanded-keys="onExpandChange"
      >
        <template #prefix="{ option }">
          <n-icon :class="getFileIcon(option)">
            <component :is="getFileIcon(option)" />
          </n-icon>
        </template>
        
        <template #suffix="{ option }">
          <span v-if="option.type === 'file'" class="file-size">
            {{ formatFileSize(option.size) }}
          </span>
        </template>
      </n-tree>
    </div>

    <!-- 文件内容预览 -->
    <div class="file-content" v-if="selectedFile">
      <div class="content-header">
        <div class="file-info">
          <h4>{{ selectedFile.name }}</h4>
          <span class="file-path">{{ selectedFile.path }}</span>
        </div>
        <div class="content-actions">
          <n-button size="tiny" quaternary @click="refreshContent">
            <template #icon>
              <n-icon><RefreshOutline /></n-icon>
            </template>
          </n-button>
          <n-button size="tiny" quaternary @click="toggleWrap">
            <template #icon>
              <n-icon><ReorderThreeOutline /></n-icon>
            </template>
          </n-button>
        </div>
      </div>

      <div class="content-body">
        <n-scrollbar>
          <div class="code-editor">
            <pre v-if="fileContent" class="code-content" :class="{ 'wrap-text': wrapText }">{{ fileContent }}</pre>
            <div v-else-if="loading" class="loading-content">
              <n-spin size="small" />
              <span>加载文件内容...</span>
            </div>
            <div v-else class="empty-content">
              <n-empty description="无法加载文件内容" size="small" />
            </div>
          </div>
        </n-scrollbar>
      </div>

      <!-- 变量使用分析 -->
      <div class="variable-analysis" v-if="variableUsages.length > 0">
        <div class="analysis-header">
          <h5>变量使用情况</h5>
          <n-tag size="small">{{ variableUsages.length }} 个变量</n-tag>
        </div>
        <div class="usage-list">
          <div
            v-for="usage in variableUsages"
            :key="usage.name"
            class="usage-item"
            @click="highlightVariable(usage)"
          >
            <div class="usage-info">
              <span class="variable-name">{{ usage.name }}</span>
              <span class="usage-count">{{ usage.count }} 次</span>
            </div>
            <div class="usage-lines">
              <n-tag
                v-for="line in usage.lines.slice(0, 3)"
                :key="line"
                size="tiny"
                class="line-tag"
              >
                第{{ line }}行
              </n-tag>
              <span v-if="usage.lines.length > 3" class="more-lines">...</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div class="empty-state" v-else>
      <n-empty description="请选择一个模板文件查看">
        <template #icon>
          <n-icon><DocumentOutline /></n-icon>
        </template>
      </n-empty>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { NTree, NIcon, NButton, NScrollbar, NSpin, NEmpty, NTag, useMessage } from 'naive-ui'
import {
  FolderOutline, FolderOpenOutline, DocumentOutline, CodeSlashOutline,
  RefreshOutline, ReorderThreeOutline, ImageOutline, ArchiveOutline
} from '@vicons/ionicons5'
import request from '@/utils/request'

const props = defineProps({
  templateId: {
    type: [String, Number],
    required: true
  },
  selectedFile: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:selectedFile', 'file-content-change'])

const message = useMessage()

// 状态
const treeData = ref([])
const selectedKeys = ref([])
const expandedKeys = ref([])
const fileContent = ref('')
const loading = ref(false)
const wrapText = ref(false)

// 变量使用分析
const variableUsages = ref([])

// 选中的文件
const selectedFile = computed(() => props.selectedFile)

// 加载模板文件树
const loadFileTree = async () => {
  try {
    const response = await request({
      url: `/api/v1/templates/${props.templateId}/files`,
      method: 'GET'
    })
    treeData.value = response.data.data || []
  } catch (error) {
    console.error('加载文件树失败:', error)
    message.error('加载文件树失败')
  }
}

// 获取文件图标
const getFileIcon = (option) => {
  if (option.type === 'directory') {
    return expandedKeys.value.includes(option.id) ? FolderOpenOutline : FolderOutline
  }
  
  const ext = option.name.split('.').pop()?.toLowerCase()
  const iconMap = {
    // 代码文件
    go: CodeSlashOutline,
    js: CodeSlashOutline,
    ts: CodeSlashOutline,
    vue: CodeSlashOutline,
    html: CodeSlashOutline,
    css: CodeSlashOutline,
    scss: CodeSlashOutline,
    json: CodeSlashOutline,
    xml: CodeSlashOutline,
    yaml: CodeSlashOutline,
    yml: CodeSlashOutline,
    // 模板文件
    tpl: DocumentOutline,
    tmpl: DocumentOutline,
    template: DocumentOutline,
    // 图片文件
    png: ImageOutline,
    jpg: ImageOutline,
    jpeg: ImageOutline,
    gif: ImageOutline,
    svg: ImageOutline,
    // 压缩文件
    zip: ArchiveOutline,
    rar: ArchiveOutline,
    tar: ArchiveOutline,
    gz: ArchiveOutline
  }
  
  return iconMap[ext] || DocumentOutline
}

// 格式化文件大小
const formatFileSize = (size) => {
  if (!size) return ''
  
  const units = ['B', 'KB', 'MB', 'GB']
  let unitIndex = 0
  let fileSize = size
  
  while (fileSize >= 1024 && unitIndex < units.length - 1) {
    fileSize /= 1024
    unitIndex++
  }
  
  return `${fileSize.toFixed(unitIndex === 0 ? 0 : 1)}${units[unitIndex]}`
}

// 选择文件
const onSelectFile = (keys) => {
  selectedKeys.value = keys
  if (keys.length > 0) {
    const file = findFileById(keys[0], treeData.value)
    if (file && file.type === 'file') {
      emit('update:selectedFile', file)
      loadFileContent(file)
    }
  }
}

// 递归查找文件
const findFileById = (id, files) => {
  for (const file of files) {
    if (file.id === id) {
      return file
    }
    if (file.children) {
      const found = findFileById(id, file.children)
      if (found) return found
    }
  }
  return null
}

// 展开状态变化
const onExpandChange = (keys) => {
  expandedKeys.value = keys
}

// 加载文件内容
const loadFileContent = async (file) => {
  if (!file || file.type !== 'file') return
  
  loading.value = true
  try {
    const response = await request({
      url: `/api/v1/templates/${props.templateId}/files/${file.id}/content`,
      method: 'GET'
    })
    
    fileContent.value = response.data.data.content || ''
    emit('file-content-change', fileContent.value)
    
    // 分析变量使用情况
    analyzeVariableUsage(fileContent.value)
  } catch (error) {
    console.error('加载文件内容失败:', error)
    message.error('加载文件内容失败')
    fileContent.value = ''
  } finally {
    loading.value = false
  }
}

// 分析变量使用情况
const analyzeVariableUsage = (content) => {
  if (!content) {
    variableUsages.value = []
    return
  }
  
  // 匹配Go模板变量 {{.VariableName}} 或 {{.Object.Field}}
  const variableRegex = /\{\{\s*\.(\w+(?:\.\w+)*)\s*\}\}/g
  const usages = new Map()
  
  const lines = content.split('\n')
  
  lines.forEach((line, index) => {
    let match
    const lineNum = index + 1
    
    while ((match = variableRegex.exec(line)) !== null) {
      const variableName = match[1]
      
      if (!usages.has(variableName)) {
        usages.set(variableName, {
          name: variableName,
          count: 0,
          lines: []
        })
      }
      
      const usage = usages.get(variableName)
      usage.count++
      
      if (!usage.lines.includes(lineNum)) {
        usage.lines.push(lineNum)
      }
    }
  })
  
  variableUsages.value = Array.from(usages.values()).sort((a, b) => b.count - a.count)
}

// 高亮变量
const highlightVariable = (usage) => {
  // 可以实现在代码编辑器中高亮变量的功能
  console.log('高亮变量:', usage.name)
}

// 刷新内容
const refreshContent = () => {
  if (selectedFile.value) {
    loadFileContent(selectedFile.value)
  }
}

// 切换换行
const toggleWrap = () => {
  wrapText.value = !wrapText.value
}

// 组件挂载时加载文件树
nextTick(() => {
  loadFileTree()
})

// 监听模板ID变化
watch(() => props.templateId, () => {
  loadFileTree()
})
</script>

<style scoped>
.template-file-explorer {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.file-tree {
  flex: 1;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  overflow-y: auto;
  min-height: 200px;
}

.file-size {
  font-size: 12px;
  color: #999;
}

.file-content {
  flex: 2;
  display: flex;
  flex-direction: column;
  min-height: 300px;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #fafafa;
  border-bottom: 1px solid #f0f0f0;
}

.file-info h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.file-path {
  font-size: 12px;
  color: #666;
}

.content-actions {
  display: flex;
  gap: 4px;
}

.content-body {
  flex: 1;
  overflow: hidden;
}

.code-editor {
  height: 100%;
}

.code-content {
  margin: 0;
  padding: 16px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.5;
  color: #333;
  background: #f8f8f8;
  overflow-x: auto;
  white-space: pre;
}

.code-content.wrap-text {
  white-space: pre-wrap;
  word-break: break-all;
}

.loading-content,
.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  gap: 12px;
  color: #666;
}

.variable-analysis {
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
  max-height: 200px;
  overflow-y: auto;
}

.analysis-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
}

.analysis-header h5 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.usage-list {
  padding: 8px;
}

.usage-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  margin-bottom: 4px;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.usage-item:hover {
  background: #f0f9ff;
}

.usage-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.variable-name {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  font-weight: 600;
  color: #1890ff;
}

.usage-count {
  font-size: 12px;
  color: #666;
}

.usage-lines {
  display: flex;
  gap: 4px;
  align-items: center;
}

.line-tag {
  font-size: 11px;
}

.more-lines {
  font-size: 12px;
  color: #999;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}
</style>