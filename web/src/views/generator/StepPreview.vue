<template>
  <div class="preview-fullscreen">
    <div class="preview-header">
      <div class="header-left">
        <span class="preview-title">预览确认</span>
        <span class="template-name">{{ templateInfo?.name }}</span>
      </div>
      <div class="header-actions">
        <n-button size="small" @click="$emit('prev')">
          <template #icon>
            <n-icon><ArrowBack /></n-icon>
          </template>
          返回配置
        </n-button>
        <n-button type="primary" size="small" @click="generateProject">
          <template #icon>
            <n-icon><Download /></n-icon>
          </template>
          生成项目
        </n-button>
      </div>
    </div>
    
    <div class="preview-main">
      <!-- 左侧：模板资源管理器 -->
      <div class="file-explorer">
        <div class="explorer-header">
          <span class="explorer-title">模板文件</span>
        </div>
        <div class="explorer-content">
          <div v-if="loading" class="loading-container">
            <n-spin size="small">
              <template #description>
                加载中...
              </template>
            </n-spin>
          </div>
          <div v-else-if="treeData.length === 0" class="empty-container">
            <div class="empty-text">暂无文件</div>
          </div>
          <div v-else class="file-tree">
            <NTree
              :data="treeDataComputed"
              :selected-keys="[currentFile]"
              :render-label="renderLabel"
              :render-switcher-icon="renderSwitcherIcon"
              @update:selected-keys="onSelectFile"
              @update:expanded-keys="updateExpandedKeys"
            />
          </div>
        </div>
      </div>
      
      <!-- 右侧：预览区域 -->
      <div class="preview-container">
        
        <div class="preview-content">
          <div v-if="!currentFile" class="no-file-selected">
            <div class="no-file-icon">
              <n-icon size="48" color="#ccc">
                <Document />
              </n-icon>
            </div>
            <div class="no-file-text">请选择左侧文件进行预览</div>
          </div>
          <div v-else class="file-preview">
            <div class="file-header">
              <div class="file-info">
                <span class="file-name">{{ currentFile }}</span>
              </div>
            </div>
            <div class="file-content">
              <div class="code-preview">
                <div class="codemirror-container" ref="codeContainer"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, h, watch, nextTick, onBeforeUnmount } from 'vue'
import { useMessage, NIcon, NTree } from 'naive-ui'
import { 
  ArrowBack, 
  Download, 
  Close, 
  Document, 
  Folder,
  FolderOpenOutline,
  FileTrayFullOutline,
  ChevronForward
} from '@vicons/ionicons5'
import { getTemplateFileTree, renderFileTree } from '@/api/templateFiles'

// CodeMirror 核心模块
import { EditorView, lineNumbers, highlightActiveLineGutter } from '@codemirror/view'
import { EditorState } from '@codemirror/state'
import { defaultHighlightStyle, syntaxHighlighting } from '@codemirror/language'

// Dracula 主题
import { dracula } from '@uiw/codemirror-theme-dracula'

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
  templateInfo: {
    type: Object,
    default: null
  },
  variables: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['prev', 'next'])

const message = useMessage()

// 文件树数据
const treeData = ref([])
const loading = ref(false)
const currentFile = ref('')
const currentFilePath = ref('')
const currentFileContent = ref('')

// 渲染后的文件数据
const renderedFiles = ref([])
const renderedFilesMap = ref(new Map())

// CodeMirror 相关
const codeContainer = ref(null)
let editorView = null

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

function getLanguageExtension(filename) {
  const ext = filename.split('.').pop()?.toLowerCase()
  return languageMap[ext] || null
}



// 展开状态
const expandedKeys = ref(new Set())

// 转换树数据为NTree格式
const treeDataComputed = computed(() => {
  return treeToNaive(treeData.value)
})

function treeToNaive(tree) {
  if (!Array.isArray(tree)) return []
  
  // 使用与模板编辑界面一致的排序逻辑
  const customSort = (a, b) => {
    if ((b.isDirectory || 0) - (a.isDirectory || 0) !== 0) {
      return (b.isDirectory || 0) - (a.isDirectory || 0)
    }
    const nameA = (a.fileName || a.label || '').toLowerCase()
    const nameB = (b.fileName || b.label || '').toLowerCase()
    const aIsA = nameA.startsWith('a') ? 1 : 0
    const bIsA = nameB.startsWith('a') ? 1 : 0
    if (aIsA !== bIsA) return bIsA - aIsA
    return nameA.localeCompare(nameB)
  }
  
  const sorted = [...tree].sort(customSort)
  
  return sorted.map(node => {
    const nodeKey = node.key || node.id
    const isExpanded = expandedKeys.value.has(String(nodeKey))
    
    // 获取显示名称：直接使用fileName，后端已提供正确的文件名
    const getDisplayName = (node) => {
      return node.fileName || node.name || ''
    }
    
    return {
      label: getDisplayName(node),
      key: nodeKey,
      isLeaf: !node.children || node.children.length === 0,
      filePath: node.filePath,
      fileName: node.fileName,
      prefix: () => h(NIcon, null, { 
        default: () => h(node.children && node.children.length > 0 ? (isExpanded ? FolderOpenOutline : Folder) : FileTrayFullOutline)
      }),
      children: node.children ? treeToNaive(node.children) : []
    }
  })
}

// 加载文件树和渲染数据
const loadTree = async () => {
  console.log('StepPreview 开始加载文件树，templateInfo:', props.templateInfo)
  if (!props.templateInfo?.id) {
    console.log('templateInfo.id 不存在，跳过加载')
    return
  }
  
  console.log('调用API获取文件树，templateId:', props.templateInfo.id)
  loading.value = true
  try {
    // 1. 获取文件树结构
    const treeRes = await getTemplateFileTree(props.templateInfo.id)
    console.log('文件树API返回结果:', treeRes)
    const tree = treeRes.data?.data?.tree
    if (tree && tree.length > 0) {
      treeData.value = tree
      console.log('设置文件树数据:', treeData.value)
    } else {
      treeData.value = []
      console.log('文件树为空')
    }
    
    // 2. 获取渲染后的文件内容
    const renderRes = await renderFileTree({
      templateId: props.templateInfo.id,
      testVariables: props.variables || {}
    })
    console.log('渲染文件树API返回结果:', renderRes)
    const files = renderRes.data?.data?.files || []
    renderedFiles.value = files
    
    // 构建文件映射，方便快速查找
    renderedFilesMap.value.clear()
    files.forEach(file => {
      renderedFilesMap.value.set(file.id, file)
    })
    
    console.log('设置渲染文件数据:', renderedFiles.value)
  } catch (error) {
    console.error('加载文件树失败:', error)
    treeData.value = []
    renderedFiles.value = []
  } finally {
    loading.value = false
  }
}

// 选择文件
const onSelectFile = async (keys) => {
  if (!keys || keys.length === 0) return
  
  const selectedKey = keys[0]
  
  // 从树数据中查找对应的节点
  const findNodeByKey = (nodes, key) => {
    for (const node of nodes) {
      if (node.key === key || node.id === key) {
        return node
      }
      if (node.children) {
        const found = findNodeByKey(node.children, key)
        if (found) return found
      }
    }
    return null
  }
  
  const selectedNode = findNodeByKey(treeData.value, selectedKey)
  if (!selectedNode || selectedNode.isDirectory === 1) {
    // 如果是文件夹，不处理
    return
  }
  
  const fileId = selectedNode.id || selectedKey
  
  // 从渲染后的文件映射中获取内容
  const renderedFile = renderedFilesMap.value.get(fileId)
  if (renderedFile) {
    currentFile.value = renderedFile.fileName
    currentFilePath.value = renderedFile.filePath
    currentFileContent.value = renderedFile.fileContent
    
    // 等待DOM更新后创建或更新CodeMirror编辑器
    await nextTick()
    createOrUpdateEditor()
  } else {
    console.error('未找到渲染后的文件:', fileId)
    currentFileContent.value = '文件未找到'
  }
}

// 创建或更新CodeMirror编辑器
const createOrUpdateEditor = async () => {
  if (!codeContainer.value) return
  
  // 销毁现有编辑器
  if (editorView) {
    editorView.destroy()
    editorView = null
  }
  
  // 获取语言扩展
  const languageExt = currentFile.value ? getLanguageExtension(currentFile.value) : null
  
  // 创建编辑器扩展
  const extensions = [
    // Dracula 主题
    dracula,
    // 只读模式
    EditorView.editable.of(false),
    // 行号
    lineNumbers(),
    // 当前行行号高亮
    highlightActiveLineGutter(),
    // 语法高亮
    syntaxHighlighting(defaultHighlightStyle),
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
  
  // 添加语言支持
  if (languageExt) {
    extensions.push(languageExt)
  }
  
  // 创建编辑器状态
  const state = EditorState.create({
    doc: currentFileContent.value,
    extensions
  })
  
  // 创建编辑器视图
  editorView = new EditorView({
    state,
    parent: codeContainer.value
  })
  
  // 确保编辑器可以滚动
  setTimeout(() => {
    if (editorView && editorView.scrollDOM) {
      editorView.requestMeasure()
    }
  }, 100)
}

// 更新展开状态
const updateExpandedKeys = (keys) => {
  expandedKeys.value = new Set(keys)
}

// 渲染标签
const renderLabel = ({ option }) => {
  return option.label
}

// 渲染切换图标 - 使用默认的展开/折叠箭头
const renderSwitcherIcon = () => h(NIcon, null, { default: () => h(ChevronForward) })





// 生成项目
const generateProject = () => {
  message.success('项目生成功能开发中...')
  // TODO: 实现项目生成逻辑
}



// 监听templateInfo变化，当有数据时加载文件树
watch(() => props.templateInfo, (newTemplateInfo) => {
  console.log('StepPreview 监听到 templateInfo 变化:', newTemplateInfo)
  if (newTemplateInfo?.id) {
    // 延迟加载，避免在第一步就加载
    setTimeout(() => {
      loadTree()
    }, 100)
  }
}, { immediate: false })

onMounted(() => {
  console.log('StepPreview 组件挂载，templateInfo:', props.templateInfo)
  console.log('StepPreview 组件挂载，templateInfo.id:', props.templateInfo?.id)
  // 如果已经有 templateInfo，则加载数据
  if (props.templateInfo?.id) {
    console.log('StepPreview 组件挂载时直接加载文件树')
    loadTree()
  }
})

onBeforeUnmount(() => {
  if (editorView) {
    editorView.destroy()
  }
})
</script>

<style scoped>
.preview-fullscreen {
  position: fixed;
  inset: 0;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
}

.preview-header {
  height: 56px;
  background: #fff;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.03);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.preview-title {
  font-size: 1.2rem;
  font-weight: bold;
  color: #18a058;
}

.template-name {
  font-size: 1rem;
  color: #666;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.preview-main {
  flex: 1;
  display: flex;
  min-height: 0;
}

/* 文件资源管理器 */
.file-explorer {
  width: 280px;
  background: #fff;
  border-right: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
}

.explorer-header {
  height: 48px;
  background: #f8f9fa;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  padding: 0 16px;
}

.explorer-title {
  font-size: 14px;
  font-weight: bold;
  color: #333;
}

.explorer-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.empty-container {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.empty-text {
  color: #999;
  font-size: 14px;
}

/* 文件树样式 */
.file-tree {
  height: 100%;
}

/* 预览容器 */
.preview-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  background: #fff;
}



.preview-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.no-file-selected {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #999;
}

.no-file-icon {
  margin-bottom: 16px;
}

.no-file-text {
  font-size: 16px;
}

.file-preview {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.file-header {
  height: 48px;
  background: #f8f9fa;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  padding: 0 16px;
}

.file-info {
  display: flex;
  align-items: center;
}

.file-name {
  font-size: 14px;
  font-weight: bold;
  color: #333;
}

.file-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.code-preview {
  flex: 1;
  overflow: auto;
  background: #1e1e1e;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.03);
}

.codemirror-container {
  height: 100%;
  min-height: 400px;
}

/* CodeMirror 样式覆盖 - 确保与编辑页面一致 */
:deep(.cm-editor) {
  height: 100% !important;
  font-size: 15px;
  outline: none !important;
}

:deep(.cm-editor .cm-scroller) {
  font-family: 'Fira Code', 'Consolas', 'Monaco', monospace;
  overflow: auto !important;
  height: 100% !important;
  max-height: none !important;
}

:deep(.cm-editor .cm-line) {
  padding: 0;
}

/* 只读模式下的光标隐藏 */
:deep(.cm-editor .cm-cursor) {
  display: none !important;
}

:deep(.cm-editor .cm-cursor-primary) {
  display: none !important;
}
</style> 