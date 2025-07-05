<template>
  <div 
    class="template-preview" 
    :class="{ collapsed: isCollapsed }"
    :style="{ width: panelWidth + 'px' }"
  >
    <div class="preview-header">
      <div class="file-info">
        <div class="file-path">{{ currentFilePath }}</div>
      </div>
      <div class="preview-actions">
        <n-button size="small" @click="copyContent" v-if="!isCollapsed">
          <template #icon>
            <n-icon><Copy /></n-icon>
          </template>
          复制
        </n-button>
        <n-button 
          size="small" 
          quaternary 
          @click="toggleCollapse"
          class="collapse-btn"
          :class="{ 'collapsed': isCollapsed }"
        >
          <template #icon>
            <n-icon>
              <ChevronForward v-if="!isCollapsed" />
              <ChevronBack v-else />
            </n-icon>
          </template>
        </n-button>
      </div>
    </div>
    
    <div v-show="!isCollapsed" class="preview-content">
      <div class="file-content" ref="previewEditorRef"></div>
    </div>
    
    <!-- 拖动调整器 -->
    <div 
      class="resize-handle"
      @mousedown="startResize"
      :class="{ 'resizing': isResizing }"
    ></div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { NButton, NIcon, useMessage } from 'naive-ui'
import { Copy, ChevronForward, ChevronBack } from '@vicons/ionicons5'
import { EditorView } from '@codemirror/view'
import { EditorState } from '@codemirror/state'
import { defaultHighlightStyle, syntaxHighlighting } from '@codemirror/language'
import { dracula } from '@uiw/codemirror-theme-dracula'
import { renderTemplate } from '@/api/templateFiles'

const props = defineProps({
  currentFile: {
    type: Object,
    default: null
  },
  fileId: {
    type: String,
    default: null
  },
  testVariables: {
    type: Object,
    default: null
  }
})

const message = useMessage()

const previewEditorRef = ref(null)
let previewEditor = null

// 折叠状态
const isCollapsed = ref(true)

// 面板宽度和拖动状态
const panelWidth = ref(400)
const isResizing = ref(false)
const startX = ref(0)
const startWidth = ref(0)

// 从本地存储加载宽度设置
const loadPanelWidth = () => {
  const savedWidth = localStorage.getItem('template-preview-width')
  if (savedWidth) {
    const width = parseInt(savedWidth)
    if (width >= 300 && width <= 1200) {
      panelWidth.value = width
    }
  }
}

// 保存宽度设置到本地存储
const savePanelWidth = (width) => {
  localStorage.setItem('template-preview-width', width.toString())
}

// 计算文件路径
const currentFilePath = computed(() => {
  return props.currentFile?.filePath || props.currentFile?.fileName || '未选择文件'
})

// 更新预览内容
function updatePreviewContent() {
  if (!previewEditor || !props.currentFile) return
  
  const content = props.currentFile.fileContent || props.currentFile.content || ''
  const transaction = previewEditor.state.update({
    changes: {
      from: 0,
      to: previewEditor.state.doc.length,
      insert: content
    }
  })
  previewEditor.dispatch(transaction)
}

// 复制内容
function copyContent() {
  if (!previewEditor) return
  
  const content = previewEditor.state.doc.toString()
  navigator.clipboard.writeText(content).then(() => {
    message.success('内容已复制到剪贴板')
  })
}

// 切换折叠状态
function toggleCollapse() {
  isCollapsed.value = !isCollapsed.value
  if (isCollapsed.value) {
    panelWidth.value = 40
  } else {
    // 展开时使用保存的宽度或默认宽度
    loadPanelWidth()
  }
}

// 开始拖动调整
function startResize(e) {
  if (isCollapsed.value) return
  
  e.preventDefault()
  e.stopPropagation()
  isResizing.value = true
  startX.value = e.clientX
  startWidth.value = panelWidth.value
  
  document.addEventListener('mousemove', handleResize)
  document.addEventListener('mouseup', stopResize)
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
}

// 处理拖动调整
function handleResize(e) {
  if (!isResizing.value) return
  
  e.preventDefault()
  e.stopPropagation()
  
  const deltaX = startX.value - e.clientX
  const newWidth = startWidth.value + deltaX
  
  // 限制最小和最大宽度
  if (newWidth >= 300 && newWidth <= 1200) {
    panelWidth.value = newWidth
  }
}

// 停止拖动调整
function stopResize() {
  isResizing.value = false
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
  
  // 保存当前宽度到本地存储
  savePanelWidth(panelWidth.value)
}

// 渲染模板
const renderTemplateContent = async () => {
  if (!props.fileId || !props.testVariables) {
    return
  }
  
  try {
    loading.value = true
    const response = await renderTemplate({
      fileId: props.fileId,
      testVariables: props.testVariables
    })
    
    if (response.code === 0) {
      renderedContent.value = response.data.fileContent
      fileName.value = response.data.fileName
    } else {
      message.error(response.msg || '渲染失败')
    }
  } catch (error) {
    console.error('渲染模板失败:', error)
    message.error('渲染模板失败')
  } finally {
    loading.value = false
  }
}

// 监听当前文件变化
watch(() => props.currentFile, () => {
  nextTick(() => {
    updatePreviewContent()
  })
})

// 初始化预览编辑器
onMounted(() => {
  // 加载保存的宽度设置
  loadPanelWidth()
  
  if (previewEditorRef.value) {
    const state = EditorState.create({
      doc: '',
      extensions: [
        EditorView.theme({
          '&': { height: '100%' },
          '.cm-editor': { height: '100%' },
          '.cm-scroller': { fontFamily: 'monospace' }
        }),
        syntaxHighlighting(defaultHighlightStyle),
        dracula
      ]
    })
    
    previewEditor = new EditorView({
      state,
      parent: previewEditorRef.value
    })
    
    // 初始化内容
    if (props.currentFile) {
      updatePreviewContent()
    }
  }
})
</script>

<style scoped>
.template-preview {
  background: #fff;
  border-left: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  height: 100%;
  transition: all 0.3s ease;
  position: relative;
  flex-shrink: 0;
}

.template-preview.collapsed {
  width: 40px !important;
}

.template-preview.collapsed .preview-header {
  padding: 8px 6px;
  justify-content: center;
}

.template-preview.collapsed .resize-handle {
  display: none;
}

.preview-header {
  padding: 12px 16px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
  transition: padding 0.3s ease;
  background: #f8f9fa;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-path {
  font-size: 13px;
  color: #666;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.preview-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.collapse-btn {
  transition: all 0.3s ease;
}

.collapse-btn.collapsed {
  width: 24px;
  height: 24px;
  padding: 0;
  border-radius: 4px;
}

.collapse-btn.collapsed:hover {
  background-color: #f0f0f0;
}

.preview-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.file-content {
  flex: 1;
  overflow: hidden;
}

.resize-handle {
  position: absolute;
  left: 0;
  top: 0;
  width: 4px;
  height: 100%;
  cursor: col-resize;
  background: transparent;
  transition: background-color 0.2s;
  z-index: 10;
  pointer-events: auto;
  user-select: none;
  touch-action: none;
}

.resize-handle:hover {
  background-color: #e0e0e0;
}

.resize-handle.resizing {
  background-color: #1890ff;
}

.template-preview.collapsed .resize-handle {
  display: none;
}
</style> 