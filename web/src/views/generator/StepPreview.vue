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
                <span class="file-name">{{ currentFile.split('/').pop() }}</span>
                <span class="file-path">{{ currentFile }}</span>
              </div>
            </div>
            <div class="file-content">
              <div class="code-preview">
                <div class="code-container">
                  <div class="line-numbers">
                    <div 
                      v-for="(line, index) in codeLines" 
                      :key="index" 
                      class="line-number"
                    >
                      {{ index + 1 }}
                    </div>
                  </div>
                  <div class="code-content">
                    <pre><code :class="codeLanguageClass">{{ currentFileContent }}</code></pre>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, h, watch, nextTick } from 'vue'
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
import { getTemplateFileTree, getTemplateFileContent } from '@/api/templateFiles'

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

// 加载文件树
const loadTree = async () => {
  console.log('StepPreview 开始加载文件树，templateInfo:', props.templateInfo)
  if (!props.templateInfo?.id) {
    console.log('templateInfo.id 不存在，跳过加载')
    return
  }
  
  console.log('调用API获取文件树，templateId:', props.templateInfo.id)
  loading.value = true
  try {
    const res = await getTemplateFileTree(props.templateInfo.id)
    console.log('文件树API返回结果:', res)
    const tree = res.data?.data?.tree
    if (tree && tree.length > 0) {
      treeData.value = tree
      console.log('设置文件树数据:', treeData.value)
    } else {
      treeData.value = []
      console.log('文件树为空')
    }
  } catch (error) {
    console.error('加载文件树失败:', error)
    treeData.value = []
  } finally {
    loading.value = false
  }
}

// 选择文件
const onSelectFile = async (keys) => {
  if (!keys || keys.length === 0 || !props.templateInfo?.id) return
  
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
  
  const filePath = selectedNode.filePath || selectedNode.fileName || String(selectedKey)
  const fileId = selectedNode.id || selectedKey
  currentFile.value = filePath
  currentFilePath.value = filePath
  
  // 加载文件内容
  try {
    // 使用文件ID获取内容
    const res = await getTemplateFileContent(fileId)
    currentFileContent.value = res.data?.data?.fileContent || ''
    
    // 等待DOM更新后应用代码高亮
    await nextTick()
    applyCodeHighlighting()
  } catch (error) {
    console.error('加载文件内容失败:', error)
    currentFileContent.value = '加载失败'
  }
}

// 应用代码高亮
const applyCodeHighlighting = () => {
  // 这里可以集成Prism.js或其他代码高亮库
  // 暂时使用简单的语法高亮
  const codeElement = document.querySelector('.code-content code')
  if (codeElement) {
    // 移除之前的语言类
    codeElement.className = codeElement.className.replace(/language-\w+/g, '')
    // 添加新的语言类
    codeElement.classList.add(codeLanguageClass.value)
  }
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

// 代码行数计算
const codeLines = computed(() => {
  if (!currentFileContent.value) return []
  return currentFileContent.value.split('\n')
})

// 根据文件扩展名获取语言类名
const codeLanguageClass = computed(() => {
  if (!currentFile.value) return ''
  const extension = currentFile.value.split('.').pop()?.toLowerCase()
  
  const languageMap = {
    'js': 'language-javascript',
    'jsx': 'language-javascript',
    'ts': 'language-typescript',
    'tsx': 'language-typescript',
    'vue': 'language-vue',
    'html': 'language-html',
    'css': 'language-css',
    'scss': 'language-scss',
    'sass': 'language-sass',
    'less': 'language-less',
    'json': 'language-json',
    'xml': 'language-xml',
    'yaml': 'language-yaml',
    'yml': 'language-yaml',
    'go': 'language-go',
    'py': 'language-python',
    'java': 'language-java',
    'c': 'language-c',
    'cpp': 'language-cpp',
    'cs': 'language-csharp',
    'php': 'language-php',
    'rb': 'language-ruby',
    'rs': 'language-rust',
    'sql': 'language-sql',
    'sh': 'language-bash',
    'bash': 'language-bash',
    'md': 'language-markdown',
    'txt': 'language-plaintext'
  }
  
  return languageMap[extension] || 'language-plaintext'
})



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
  flex-direction: column;
  gap: 2px;
}

.file-name {
  font-size: 14px;
  font-weight: bold;
  color: #333;
}

.file-path {
  font-size: 12px;
  color: #666;
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
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

.code-container {
  display: flex;
  min-height: 100%;
}

.line-numbers {
  display: flex;
  flex-direction: column;
  background: #282a36;
  border-right: 1px solid #44475a;
  padding: 16px 8px 16px 16px;
  user-select: none;
  min-width: 50px;
}

.line-number {
  font-family: 'Fira Code', 'Consolas', 'Monaco', monospace;
  font-size: 15px;
  line-height: 1.5;
  color: #6272a4;
  text-align: right;
  padding-right: 8px;
}

.code-content {
  flex: 1;
  overflow: auto;
  padding: 16px;
}

.code-content pre {
  margin: 0;
  font-family: 'Fira Code', 'Consolas', 'Monaco', monospace;
  font-size: 15px;
  line-height: 1.5;
  color: #f8f8f2;
  background: transparent;
}

.code-content code {
  font-family: inherit;
  background: transparent;
}

/* 代码高亮主题 - Dracula主题 */
.code-content .language-javascript,
.code-content .language-typescript,
.code-content .language-vue,
.code-content .language-html,
.code-content .language-css,
.code-content .language-scss,
.code-content .language-sass,
.code-content .language-less,
.code-content .language-json,
.code-content .language-xml,
.code-content .language-yaml,
.code-content .language-go,
.code-content .language-python,
.code-content .language-java,
.code-content .language-c,
.code-content .language-cpp,
.code-content .language-csharp,
.code-content .language-php,
.code-content .language-ruby,
.code-content .language-rust,
.code-content .language-sql,
.code-content .language-bash,
.code-content .language-markdown {
  color: #f8f8f2;
}

/* 关键字高亮 - Dracula主题 */
.code-content .keyword {
  color: #ff79c6;
}

/* 字符串高亮 - Dracula主题 */
.code-content .string {
  color: #f1fa8c;
}

/* 注释高亮 - Dracula主题 */
.code-content .comment {
  color: #6272a4;
}

/* 数字高亮 - Dracula主题 */
.code-content .number {
  color: #bd93f9;
}

/* 函数名高亮 - Dracula主题 */
.code-content .function {
  color: #50fa7b;
}

/* 类名高亮 - Dracula主题 */
.code-content .class-name {
  color: #8be9fd;
}

/* 变量名高亮 - Dracula主题 */
.code-content .variable {
  color: #f8f8f2;
}

/* 属性名高亮 - Dracula主题 */
.code-content .property {
  color: #66d9ef;
}

/* 标签名高亮 - Dracula主题 */
.code-content .tag {
  color: #ff79c6;
}

/* 操作符高亮 - Dracula主题 */
.code-content .operator {
  color: #ff79c6;
}
</style> 