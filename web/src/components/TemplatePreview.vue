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
import { EditorView, highlightActiveLine, highlightActiveLineGutter, lineNumbers } from '@codemirror/view'
import { EditorState } from '@codemirror/state'
import { defaultHighlightStyle, syntaxHighlighting } from '@codemirror/language'
import { dracula } from '@uiw/codemirror-theme-dracula'
import { renderTemplate } from '@/api/templateFiles'

// 语言支持
import { javascript } from '@codemirror/lang-javascript'
import { html } from '@codemirror/lang-html'
import { css } from '@codemirror/lang-css'
import { json } from '@codemirror/lang-json'
import { markdown } from '@codemirror/lang-markdown'
import { python } from '@codemirror/lang-python'
import { java } from '@codemirror/lang-java'
import { cpp } from '@codemirror/lang-cpp'
import { rust } from '@codemirror/lang-rust'
import { go } from '@codemirror/lang-go'
import { sql } from '@codemirror/lang-sql'
import { xml } from '@codemirror/lang-xml'
import { yaml } from '@codemirror/lang-yaml'
import { vue } from '@codemirror/lang-vue'

const props = defineProps({
  currentFile: {
    type: Object,
    default: null
  },
  fileId: {
    type: String,
    default: null
  },
  variables: {
    type: Object,
    default: null
  }
})

const message = useMessage()

const previewEditorRef = ref(null)
let previewEditor = null

// 渲染相关状态
const loading = ref(false)
const renderedContent = ref('')
const fileName = ref('')

// 语言映射
const languageMap = {
  'js': javascript(),
  'javascript': javascript(),
  'ts': javascript({typescript: true}),
  'typescript': javascript({typescript: true}),
  'jsx': javascript({jsx: true}),
  'tsx': javascript({typescript: true, jsx: true}),
  'vue': vue(),
  'html': html(),
  'htm': html(),
  'css': css(),
  'scss': css(),
  'sass': css(),
  'less': css(),
  'json': json(),
  'md': markdown(),
  'markdown': markdown(),
  'py': python(),
  'python': python(),
  'java': java(),
  'cpp': cpp(),
  'cc': cpp(),
  'cxx': cpp(),
  'c': cpp(),
  'rs': rust(),
  'rust': rust(),
  'go': go(),
  'sql': sql(),
  'xml': xml(),
  'yaml': yaml(),
  'yml': yaml()
}

// 获取语言扩展
function getLanguageExtension(filename) {
  const ext = filename.split('.').pop()?.toLowerCase()
  return languageMap[ext] || null
}

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
  if (fileName.value) {
    return `预览: ${fileName.value}`
  }
  return props.currentFile?.filePath || props.currentFile?.fileName || '未选择文件'
})

// 更新预览内容
function updatePreviewContent() {
  if (!previewEditor) {
    return
  }
  
  let content = ''
  let currentFileName = ''
  
  // 如果有渲染后的内容，优先显示
  if (renderedContent.value) {
    content = renderedContent.value
    currentFileName = fileName.value
  } else if (props.currentFile) {
    content = props.currentFile.fileContent || props.currentFile.content || ''
    currentFileName = props.currentFile.fileName || props.currentFile.name || ''
  }
  // 获取语言扩展
  const languageExt = getLanguageExtension(currentFileName)
  
  // 更新编辑器状态，包括语言支持
  const newState = EditorState.create({
    doc: content,
    extensions: [
      // Dracula 主题（与主编辑器一致）
      dracula,
      // 语法高亮
      syntaxHighlighting(defaultHighlightStyle),
      // 行号
      lineNumbers(),
      // 当前行高亮
      highlightActiveLine(),
      // 当前行行号高亮
      highlightActiveLineGutter(),
      // 确保滚动正常工作
      EditorView.scrollMargins.of(() => ({ top: 10, bottom: 10 })),
      // 强制启用滚动
      EditorView.theme({
        "&": { height: "100%" },
        ".cm-scroller": { 
          overflow: "auto !important",
          height: "100% !important"
        }
      }),
      // 添加语言支持
      ...(languageExt ? [languageExt] : [])
    ]
  })
  
  previewEditor.setState(newState)
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

// 展开面板（如果已折叠）
function expandPanel() {
  if (isCollapsed.value) {
    isCollapsed.value = false
    loadPanelWidth()
  }
}

// 暴露方法给父组件
defineExpose({
  toggleCollapse,
  expandPanel
})

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
  if (!props.fileId) {
    return
  }
  
  // 确保有测试变量，如果没有则使用空对象
  const variables = props.variables || {}
  
  try {
    loading.value = true
    
    const response = await renderTemplate({
      fileId: props.fileId,
      variables: variables
    })
    
    if (response.data.code === 0) {
      renderedContent.value = response.data.data.fileContent
      fileName.value = response.data.data.fileName
      
      // 立即更新预览内容
      nextTick(() => {
        updatePreviewContent()
      })
    } else {
      message.error(response.data.message || response.data.msg || '渲染失败')
    }
  } catch (error) {
    console.error('渲染模板失败:', error)
    message.error('渲染模板失败: ' + (error.message || '未知错误'))
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

// 监听文件ID和测试变量变化，自动渲染
watch([() => props.fileId, () => props.variables], async () => {
  
  if (props.fileId) {
    await renderTemplateContent()
    // 渲染成功后会自动调用 updatePreviewContent
  }
}, { deep: true, immediate: false })

// 初始化预览编辑器
onMounted(() => {
  // 加载保存的宽度设置
  loadPanelWidth()
  
  if (previewEditorRef.value) {
    const state = EditorState.create({
      doc: '',
      extensions: [
        // Dracula 主题（与主编辑器一致）
        dracula,
        // 语法高亮
        syntaxHighlighting(defaultHighlightStyle),
        // 行号
        lineNumbers(),
        // 当前行高亮
        highlightActiveLine(),
        // 当前行行号高亮
        highlightActiveLineGutter(),
        // 确保滚动正常工作
        EditorView.scrollMargins.of(() => ({ top: 10, bottom: 10 })),
        // 强制启用滚动
        EditorView.theme({
          "&": { height: "100%" },
          ".cm-scroller": { 
            overflow: "auto !important",
            height: "100% !important"
          }
        })
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