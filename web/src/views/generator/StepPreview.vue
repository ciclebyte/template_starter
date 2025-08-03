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
import { renderFileTree, downloadZip } from '@/api/templateFiles'
import { listTemplateVariables } from '@/api/templateVariables'

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
const renderedFilesMap = ref(new Map())

// 模板变量信息
const templateVariables = ref([])

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
  
  // 排序逻辑：目录在前，文件在后，同类型按名称排序
  const customSort = (a, b) => {
    // 首先按目录/文件排序：目录在前，文件在后
    if ((b.isDirectory || 0) - (a.isDirectory || 0) !== 0) {
      return (b.isDirectory || 0) - (a.isDirectory || 0)
    }
    // 同类型按名称排序
    const nameA = (a.fileName || a.label || '').toLowerCase()
    const nameB = (b.fileName || b.label || '').toLowerCase()
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
      isLeaf: node.isDirectory !== 1, // 目录永远不是叶子节点，文件永远是叶子节点
      filePath: node.filePath,
      fileName: node.fileName,
      prefix: () => h(NIcon, null, { 
        default: () => h(node.isDirectory === 1 ? (isExpanded ? FolderOpenOutline : Folder) : FileTrayFullOutline)
      }),
      children: node.children ? treeToNaive(node.children) : []
    }
  })
}

// 获取模板变量信息
const loadTemplateVariables = async () => {
  if (!props.templateInfo?.id) {
    return
  }
  
  try {
    const res = await listTemplateVariables({
      templateId: props.templateInfo.id,
      pageNum: 1,
      pageSize: 1000 // 获取所有变量
    })
    templateVariables.value = res.data?.data?.templateVariablesList || []
  } catch (error) {
    console.error('获取模板变量失败:', error)
    templateVariables.value = []
  }
}

// 根据变量类型转换变量值
const convertVariableTypes = (variables) => {
  const converted = { ...variables }
  const variableTypeMap = {}
  
  // 创建变量类型映射
  templateVariables.value.forEach(v => {
    variableTypeMap[v.name] = v.variableType
  })
  
  // 转换变量类型
  Object.keys(converted).forEach(key => {
    const variableType = variableTypeMap[key]
    if (!variableType) return
    
    const value = converted[key]
    const valueStr = String(value)
    
    switch (variableType) {
      case 'string':
        converted[key] = valueStr
        break
      case 'boolean':
        // 字符串 "false" 转换为 false，其他非空值转换为 true
        converted[key] = valueStr === 'false' ? false : Boolean(valueStr && valueStr !== '0')
        break
      case 'number':
        // 尝试转换为数字
        const numValue = Number(valueStr)
        converted[key] = isNaN(numValue) ? 0 : numValue
        break
      case 'list':
        // 如果是JSON格式的数组字符串，尝试解析
        if (valueStr.startsWith('[') && valueStr.endsWith(']')) {
          try {
            converted[key] = JSON.parse(valueStr)
          } catch {
            // 如果解析失败，按逗号分割
            converted[key] = valueStr.slice(1, -1).split(',').map(s => s.trim())
          }
        } else {
          // 按逗号分割字符串
          converted[key] = valueStr.split(',').map(s => s.trim())
        }
        break
      case 'object':
        // 如果是JSON格式的对象字符串，尝试解析
        if ((valueStr.startsWith('{') && valueStr.endsWith('}')) ||
            (valueStr.startsWith('[') && valueStr.endsWith(']'))) {
          try {
            converted[key] = JSON.parse(valueStr)
          } catch {
            converted[key] = valueStr
          }
        } else {
          converted[key] = valueStr
        }
        break
      default:
        // 未知类型，保持原值
        break
    }
  })
  
  return converted
}

// 加载渲染后的文件树
const loadTree = async () => {
  if (!props.templateInfo?.id) {
    console.log('templateInfo.id 不存在，跳过加载')
    return
  }
  loading.value = true
  try {
    // 先获取模板变量信息
    await loadTemplateVariables()
    
    // 转换变量类型
    const convertedVariables = convertVariableTypes(props.variables || {})
    
    // 获取渲染后的文件树
    const renderRes = await renderFileTree({
      templateId: props.templateInfo.id,
      variables: convertedVariables
    })
    const tree = renderRes.data?.data?.tree || []
    treeData.value = tree
    
    // 构建文件映射，方便快速查找
    renderedFilesMap.value.clear()
    const flattenFiles = (nodes) => {
      nodes.forEach(node => {
        renderedFilesMap.value.set(node.id, node)
        if (node.children && node.children.length > 0) {
          flattenFiles(node.children)
        }
      })
    }
    flattenFiles(tree)
    
  } catch (error) {
    console.error('加载文件树失败:', error)
    treeData.value = []
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
const generateProject = async () => {
  try {
    message.loading('正在生成项目...', { duration: 0 })
    
    // 转换变量类型
    const convertedVariables = convertVariableTypes(props.variables || {})
    
    const response = await downloadZip({
      templateId: props.templateInfo.id,
      variables: convertedVariables,
      fileName: props.templateInfo.name
    })
    
    // 创建下载链接
    const blob = new Blob([response.data], { type: 'application/zip' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${props.templateInfo.name}.zip`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    message.destroyAll()
    message.success('项目生成成功！')
  } catch (error) {
    message.destroyAll()
    console.error('生成项目失败:', error)
    message.error('生成项目失败，请重试')
  }
}



// 监听templateInfo变化，当有数据时加载文件树
watch(() => props.templateInfo, (newTemplateInfo) => {
  if (newTemplateInfo?.id) {
    // 延迟加载，避免在第一步就加载
    setTimeout(() => {
      loadTree()
    }, 100)
  }
}, { immediate: false })

onMounted(() => {
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