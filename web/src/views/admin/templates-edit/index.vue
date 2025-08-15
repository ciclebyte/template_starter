<template>
  <div class="templates-edit-fullscreen">
    <EditHeader 
      :is-variable-panel-open="isVariablePanelOpen" 
      :is-file-tree-visible="isFileTreeVisible"
      :has-unsaved-changes="hasUnsavedChanges"
      :current-file-name="currentFileName"
      @toggle-variable-panel="toggleVariablePanel" 
      @show-variable-manager="showVariableManager = true" 
      @show-preset-manager="showPresetManager = true"
      @close-edit="closeEdit" 
      @toggle-file-tree="toggleFileTree" 
      @show-settings="showSettings = true"
    />
      
    <!-- 自动保存指示器 -->
    <transition name="auto-save-fade">
      <div v-if="autoSaveIndicator" class="auto-save-indicator">
        <n-icon size="16" style="margin-right: 6px;">
          <svg viewBox="0 0 24 24">
            <path fill="currentColor" d="M17 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14c1.1 0 2-.9 2-2V7l-4-4zm-5 16c-1.66 0-3-1.34-3-3s1.34-3 3-3 3 1.34 3 3-1.34 3-3 3zm3-10H5V5h10v4z"/>
          </svg>
        </n-icon>
        已自动保存
      </div>
    </transition>

    <!-- 变量插入面板 -->
    <VariablePanel :is-open="isVariablePanelOpen" :template-variables="templateVariables"
      :template-syntax-categories="templateSyntaxCategories" :builtin-function-categories="builtinFunctionCategories"
      :sprig-function-categories="sprigFunctionCategories" :loading-functions="loadingFunctions"
      :loading-sprig-functions="loadingSprigFunctions" :quick-variables="quickVariables" @insert-syntax="insertSyntax"
      @insert-function="insertFunction" @insert-sprig-function="insertSprigFunction" @insert-variable="insertVariable"
      @show-variable-manager="showVariableManager = true" @update:height="variablePanelHeight = $event" />

    <div class="edit-main">
      <!-- 左侧：模板资源管理器 -->
      <TemplateExplorer v-show="isFileTreeVisible" v-model:treeData="treeData" :currentFile="currentFile" @select="onSelectFile"
        @reload="onTreeReload" @rename="onRenameFile" @uploadZip="onUploadZip" @uploadCodeFile="onUploadCodeFile"
        @move="onMoveFile" @setCondition="onSetCondition" />

      <!-- 中间：编辑器区域 -->
      <div class="editor-container">
        <TemplateEditor ref="templateEditorRef" :currentFileName="currentFileName" :currentFileId="currentFileId"
          :currentFileContent="currentFileContent" @contentChange="onEditorContentChange"
          @insertVariable="onInsertVariable" @preview="onPreview" @selectionChange="onEditorSelectionChange" />
      </div>

      <!-- 右侧：预览面板 -->
      <TemplatePreview ref="templatePreviewRef" :current-file="currentFileNode" :file-id="previewFileId"
        :variables="variableValues" />
    </div>



    <!-- 变量管理弹框 -->
    <n-modal v-model:show="showVariableManager" preset="card" style="width: 80vw; height: 80vh; max-width: 1200px;"
      :mask-closable="false">
      <template #header>
        <div class="variable-manager-header">
          <span class="modal-title">变量管理</span>
        </div>
      </template>

      <div class="variable-manager-content">
        <VariableManager ref="variableManagerRef" :variables="templateVariables" :template-id="route.params.id"
          @add="onAddVariable" @edit="onEditVariable" @delete="onDeleteVariable" @insert="onInsertVariable"
          @applyTestData="onApplyTestData" />
      </div>
    </n-modal>

    <!-- 预设变量管理弹框 -->
    <n-modal v-model:show="showPresetManager" preset="card" style="width: 90vw; height: 85vh; max-width: 1400px;"
      :mask-closable="false">
      <template #header>
        <div class="preset-manager-header">
          <span class="modal-title">预设变量管理</span>
        </div>
      </template>

      <div class="preset-manager-content">
        <PresetVariableManager :template-id="route.params.id" />
      </div>
    </n-modal>

    <!-- 条件设置弹框 -->
    <ConditionModal ref="conditionModalRef" v-model:show="showConditionModal"
      :selected-file-for-condition="selectedFileForCondition" :template-variables="templateVariables"
      @close="showConditionModal = false" @save="handleConditionSave" />

    <!-- AI助手组件 -->
    <AIAssistant :current-file-name="currentFileName" :current-file-content="currentFileContent"
      :template-variables="templateVariables" :editor-selection="editorSelection" @insert-code="onAIInsertCode"
      @add-variable="onAIAddVariable" @create-file="onAICreateFile" @apply-suggestion="onAIApplySuggestion" />

    <!-- AI SDK面板组件 - 暂时隐藏 -->
    <AISDKPanel :current-file-name="currentFileName" :current-file-content="currentFileContent"
      :template-variables="templateVariables" :editor-selection="editorSelection" @insert-code="onAIInsertCode"
      @add-variable="onAIAddVariable" @create-file="onAICreateFile" @apply-suggestion="onAIApplySuggestion"
      style="display: none;" />
      
    <!-- 编辑器设置面板 -->
    <EditorSettings 
      v-model:show="showSettings" 
      :settings="editorSettings"
      @save-settings="saveSettings" 
    />
    
  </div>
</template>

<script setup>
import { useRouter, useRoute } from 'vue-router'
import { ref, onMounted, onUnmounted, watch, computed, nextTick } from 'vue'
import { getTemplateFileTree, addTemplateFile, delTemplateFile, getTemplateFileDetail, getTemplateFileContent, renameTemplateFile, uploadZipFile, uploadCodeFile, moveTemplateFile, setFileCondition, getFileCondition } from '@/api/templateFiles'
import { listTemplateVariables, addTemplateVariable, editTemplateVariable, deleteTemplateVariable } from '@/api/templateVariables'
import { getBuiltinFunctions } from '@/api/builtinFunctions'
import { getSprigFunctions } from '@/api/sprigFunctions'
import TemplateExplorer from './components/TemplateFileTree.vue'
import TemplateEditor from './components/TemplateEditor.vue'
import VariableManager from './components/VariableManager.vue'
import PresetVariableManager from './components/PresetVariableManager.vue'
import TemplatePreview from './components/TemplatePreview.vue'
import AIAssistant from './components/AIAssistant.vue'
import AISDKPanel from './components/AISDKPanel.vue'
import EditHeader from './components/EditHeader.vue'
import ConditionModal from './components/ConditionModal.vue'
import VariablePanel from './components/VariablePanel.vue'
import EditorSettings from './components/EditorSettings.vue'
import { templateSyntaxCategories as syntaxData } from './data/templateSyntax.js'
import { useTemplateFileStore } from '@/stores/templateFileStore'
import { useMessage, NIcon, NButton, NSpin, NForm, NFormItem, NSwitch, NSelect, NRadioGroup, NRadio, NInput, NModal, NCard } from 'naive-ui'
import { ChevronDown, ChevronUp, Add, Settings, CloseOutline } from '@vicons/ionicons5'

const router = useRouter()
const route = useRoute()
const message = useMessage()
const closeEdit = () => {
  router.push('/templates')
}

const treeData = ref([])
const loadingTree = ref(true)
const noTreeData = ref(false)
const currentFile = ref('')
const currentFileContent = ref('')
const currentFileName = ref('')
const currentFileId = ref('')
const templateFileStore = useTemplateFileStore()

// 变量相关
const templateVariables = ref([])
const templateEditorRef = ref(null)
const variableManagerRef = ref(null)
const conditionModalRef = ref(null)
const showVariableManager = ref(false)
const showPresetManager = ref(false)

const variableValues = ref({})
const currentFileNode = ref(null)
const templatePreviewRef = ref(null)
const previewFileId = ref(null)

// 编辑器选中状态
const editorSelection = ref({
  hasSelection: false,
  selectedText: '',
  selectionStart: 0,
  selectionEnd: 0,
  selectionLength: 0
})

// 条件设置相关
const showConditionModal = ref(false)
const selectedFileForCondition = ref(null)

// 内置变量列表
const quickVariables = [
  { name: 'ProjectName', label: '项目名称' },
  { name: 'Author', label: '作者' },
  { name: 'PackageName', label: '包名' }
]

// 动态函数分类数据
const builtinFunctionCategories = ref([])
const loadingFunctions = ref(false)

// Sprig函数分类数据
const sprigFunctionCategories = ref([])
const loadingSprigFunctions = ref(false)

// 函数详情面板
const functionDetailVisible = ref(false)
const selectedFunction = ref(null)
const functionDetailStyle = ref({})
let hideTimer = null
let showTimer = null


// 变量面板状态
const isVariablePanelOpen = ref(false)
const variablePanelHeight = ref(300) // 默认高度300px

// 文件树状态（默认显示，从本地存储加载）
const isFileTreeVisible = ref(
  localStorage.getItem('template-file-tree-visible') === null ? true : localStorage.getItem('template-file-tree-visible') !== 'false'
)

// 设置面板状态
const showSettings = ref(false)

const editorSettings = ref({
  autoSave: {
    enabled: true,
    interval: 30
  },
  editor: {
    fontSize: 14,
    lineNumbers: true,
    wordWrap: true
  },
  interface: {
    theme: 'light',
    restoreLayout: true
  },
  preview: {
    realtime: true,
    debounceDelay: 500
  }
})

// Go模板语法数据
const templateSyntaxCategories = ref(syntaxData)


const onVariablePanelResize = (event) => {
  if (!isVariablePanelResizing.value) return

  const deltaY = event.clientY - variablePanelStartY.value
  const newHeight = variablePanelStartHeight.value + deltaY

  // 限制高度范围
  if (newHeight >= minVariablePanelHeight && newHeight <= maxVariablePanelHeight) {
    variablePanelHeight.value = newHeight
  }

  event.preventDefault()
}

const stopVariablePanelResize = () => {
  isVariablePanelResizing.value = false

  // 移除全局鼠标事件监听
  document.removeEventListener('mousemove', onVariablePanelResize)
  document.removeEventListener('mouseup', stopVariablePanelResize)

  // 恢复文本选择
  document.body.style.userSelect = ''
}

// 切换变量面板显示状态
const toggleVariablePanel = () => {
  isVariablePanelOpen.value = !isVariablePanelOpen.value
}

// 切换文件树显示状态
const toggleFileTree = () => {
  isFileTreeVisible.value = !isFileTreeVisible.value
  // 保存状态到本地存储
  localStorage.setItem('template-file-tree-visible', isFileTreeVisible.value.toString())
}

// 设置相关方法
const loadSettings = () => {
  const savedSettings = localStorage.getItem('template-editor-settings')
  if (savedSettings) {
    try {
      const parsed = JSON.parse(savedSettings)
      editorSettings.value = { ...editorSettings.value, ...parsed }
    } catch (e) {
      console.warn('Failed to load settings:', e)
    }
  }
}

const saveSettings = (newSettings) => {
  editorSettings.value = { ...newSettings }
  localStorage.setItem('template-editor-settings', JSON.stringify(newSettings))
  
  // 应用设置
  applySettings()
}

const applySettings = () => {
  // 应用自动保存设置
  setupAutoSave()
  
  // 这里可以添加其他设置的应用逻辑
  // 例如：主题切换、字体大小等
}

// 自动保存功能
let autoSaveTimer = null

const setupAutoSave = () => {
  // 清除现有定时器
  if (autoSaveTimer) {
    clearInterval(autoSaveTimer)
    autoSaveTimer = null
  }
  
  // 如果启用自动保存，设置定时器
  if (editorSettings.value.autoSave.enabled) {
    const interval = editorSettings.value.autoSave.interval * 1000 // 转换为毫秒
    autoSaveTimer = setInterval(() => {
      autoSaveCurrentFile()
    }, interval)
  }
}

const saveCurrentFile = async (silent = false) => {
  // 调用编辑器的保存方法
  if (templateEditorRef.value && templateEditorRef.value.saveCurrentFile) {
    const success = await templateEditorRef.value.saveCurrentFile(silent)
    if (success) {
      hasUnsavedChanges.value = false
    }
    return success
  }
  return false
}

const autoSaveCurrentFile = async () => {
  // 检查是否有当前文件和未保存的更改
  if (currentFileId.value && hasUnsavedChanges.value) {
    try {
      const success = await saveCurrentFile(true) // 静默保存
      if (success) {
        console.log('Auto-saved file:', currentFileName.value)
        // 可以在这里添加轻量级的自动保存指示器
        showAutoSaveIndicator()
      }
    } catch (error) {
      console.error('Auto-save failed:', error)
    }
  }
}

// 检查是否有未保存的更改
const hasUnsavedChanges = ref(false)

// 自动保存指示器
const autoSaveIndicator = ref(false)

const showAutoSaveIndicator = () => {
  autoSaveIndicator.value = true
  setTimeout(() => {
    autoSaveIndicator.value = false
  }, 2000) // 2秒后隐藏
}

// 键盘快捷键处理
const handleKeyDown = (event) => {
  // Ctrl+B 或 Cmd+B 切换文件树
  if ((event.ctrlKey || event.metaKey) && event.key === 'b') {
    event.preventDefault()
    toggleFileTree()
  }
  
  // Ctrl+, 打开设置
  if ((event.ctrlKey || event.metaKey) && event.key === ',') {
    event.preventDefault()
    showSettings.value = true
  }
}





onMounted(async () => {
  await loadTree()
  await loadVariables()
  await loadBuiltinFunctions()
  await loadSprigFunctions()
  loadTestData()
  
  // 加载设置并应用
  loadSettings()
  applySettings()
  
  // 添加键盘事件监听
  document.addEventListener('keydown', handleKeyDown)
})

// 组件卸载时清理事件监听器
onUnmounted(() => {
  // 清理键盘事件监听
  document.removeEventListener('keydown', handleKeyDown)
  
  // 清理自动保存定时器
  if (autoSaveTimer) {
    clearInterval(autoSaveTimer)
    autoSaveTimer = null
  }
  
  // 清理变量面板拖拽事件监听器
  document.removeEventListener('mousemove', onVariablePanelResize)
  document.removeEventListener('mouseup', stopVariablePanelResize)

  // 恢复文本选择
  document.body.style.userSelect = ''

  // 清理函数详情面板的定时器
  if (hideTimer) {
    clearTimeout(hideTimer)
  }
  if (showTimer) {
    clearTimeout(showTimer)
  }
})



async function loadTree() {
  loadingTree.value = true
  try {
    const res = await getTemplateFileTree(route.params.id)
    const tree = res.data && res.data.data && res.data.data.tree
    if (tree && tree.length > 0) {
      treeData.value = tree
      noTreeData.value = false
    } else {
      treeData.value = []
      noTreeData.value = true
    }
  } catch (e) {
    treeData.value = []
    noTreeData.value = true
  }
  loadingTree.value = false
}

async function loadVariables() {
  try {
    const res = await listTemplateVariables({ templateId: route.params.id })
    templateVariables.value = res.data.data.templateVariablesList || []
  } catch (e) {
    templateVariables.value = []
    console.error('加载变量失败:', e)
  }
}

// 加载内置函数
async function loadBuiltinFunctions() {
  loadingFunctions.value = true
  try {
    const res = await getBuiltinFunctions()
    if (res.data && res.data.data) {
      builtinFunctionCategories.value = res.data.data.categories || []
    }
  } catch (error) {
    console.error('加载内置函数失败:', error)
    message.error('加载内置函数失败')
    builtinFunctionCategories.value = []
  } finally {
    loadingFunctions.value = false
  }
}

// 加载Sprig函数
async function loadSprigFunctions() {
  loadingSprigFunctions.value = true
  try {
    const res = await getSprigFunctions()
    if (res.data && res.data.data) {
      sprigFunctionCategories.value = res.data.data.categories || []
    }
  } catch (error) {
    console.error('加载Sprig函数失败:', error)
    message.error('加载Sprig函数失败')
    sprigFunctionCategories.value = []
  } finally {
    loadingSprigFunctions.value = false
  }
}

// 加载测试数据
function loadTestData() {
  const templateId = route.params.id
  const testDataKey = `template_test_data_${templateId}`
  const savedTestData = localStorage.getItem(testDataKey)
  if (savedTestData) {
    try {
      variableValues.value = JSON.parse(savedTestData)
    } catch (e) {
      variableValues.value = {}
    }
  } else {
    variableValues.value = {}
  }
}

function findNodeByKey(list, key) {
  console.log('查找节点:', { key, listLength: list?.length })

  for (const item of list) {
    const itemKey = String(item.key || item.id)
    const targetKey = String(key)

    console.log('检查节点:', {
      itemKey,
      targetKey,
      matches: itemKey === targetKey,
      fileName: item.fileName,
      isDirectory: item.isDirectory
    })

    if (itemKey === targetKey) {
      console.log('找到匹配节点:', item)
      return item
    }

    if (item.children) {
      const found = findNodeByKey(item.children, key)
      if (found) {
        console.log('在子节点中找到:', found)
        return found
      }
    }
  }

  console.log('未找到匹配的节点')
  return null
}

async function onSelectFile(key) {
  console.log('选择文件:', { key, treeDataLength: treeData.value.length })

  // 重置编辑器选中状态
  resetEditorSelection()
  
  // 重置未保存状态
  hasUnsavedChanges.value = false

  currentFile.value = key
  const node = findNodeByKey(treeData.value, key)
  currentFileNode.value = node

  console.log('找到节点:', { node, isDirectory: node?.isDirectory })

  if (node && node.isDirectory === 0) {
    try {
      // 先设置文件名，避免编辑器被销毁
      const fileName = node.fileName || node.label || String(key)
      console.log('设置文件名:', fileName)

      currentFileName.value = fileName
      currentFileId.value = String(key)

      // 然后异步加载文件内容
      console.log('开始加载文件内容...')
      const res = await getTemplateFileContent(key)
      const content = res.data.data.fileContent

      console.log('文件内容加载完成:', { contentLength: content?.length })

      // 设置文件内容
      currentFileContent.value = content
      templateFileStore.setCurrentFileContent(content)

    } catch (e) {
      console.error('加载文件内容失败:', e)
      currentFileContent.value = ''
      currentFileName.value = ''
      currentFileId.value = ''
      templateFileStore.setCurrentFileContent('')
    }
  } else {
    // 对于目录，清空所有文件信息
    console.log('选择的是目录，清空文件信息')
    currentFileContent.value = ''
    currentFileName.value = ''
    currentFileId.value = ''
    templateFileStore.setCurrentFileContent('')
  }
}

function onEditorContentChange({ content }) {
  currentFileContent.value = content
  templateFileStore.setCurrentFileContent(content)
  hasUnsavedChanges.value = true
}

// 处理编辑器选中变化
function onEditorSelectionChange(selectionInfo) {
  editorSelection.value = selectionInfo
}

// 重置编辑器选中状态
function resetEditorSelection() {
  editorSelection.value = {
    hasSelection: false,
    selectedText: '',
    selectionStart: 0,
    selectionEnd: 0,
    selectionLength: 0
  }
}

function onTreeReload(payload) {
  if (payload && payload.type === 'delete') {
    // 调用真实删除接口
    delTemplateFile(payload.id).then(() => {
      loadTree()
    })
    return
  }
  // 这里处理新增文件/文件夹逻辑，后续可接接口
  // 调用真实接口
  const templateId = route.params.id
  const isDirectory = payload.type === 'folder' ? 1 : 0
  addTemplateFile({
    templateId,
    fileName: payload.name,
    fileContent: '',
    fileSize: 0,
    isDirectory,
    md5: '',
    sort: 0,
    parentId: payload.parentId
  }).then(() => {
    loadTree()
  })
}

async function onRenameFile(payload) {
  const { id, oldName, newName, node } = payload

  // 验证新名称
  if (!newName || newName.trim() === '' || newName === oldName) {
    return
  }

  try {
    await renameTemplateFile({
      id: parseInt(id),
      fileName: newName.trim()
    })

    message.success('重命名成功')

    // 重新加载文件树
    await loadTree()

    // 如果重命名的是当前文件，更新当前文件名
    if (currentFileId.value === String(id)) {
      currentFileName.value = newName.trim()
    }
  } catch (error) {
    console.error('重命名失败:', error)
    message.error('重命名失败: ' + (error.response?.data?.message || error.message || '未知错误'))
  }
}

async function onUploadZip(payload) {
  const { file } = payload

  try {
    const res = await uploadZipFile(route.params.id, file)
    const { successCount, failedFiles, message: resultMessage } = res.data.data

    if (failedFiles && failedFiles.length > 0) {
      message.warning(`${resultMessage}，成功 ${successCount} 个文件，失败 ${failedFiles.length} 个文件`)
      console.log('失败的文件:', failedFiles)
    } else {
      message.success(`${resultMessage}，成功解压 ${successCount} 个文件`)
    }

    // 重新加载文件树
    await loadTree()
  } catch (error) {
    console.error('上传错误:', error)
    message.error('ZIP包上传失败：' + (error.message || '未知错误'))
  }
}

async function onUploadCodeFile(payload) {
  const { file, parentId } = payload
  try {
    const res = await uploadCodeFile(route.params.id, file, parentId)
    const { fileName, isTextFile, message: resultMessage } = res.data.data
    if (isTextFile) {
      message.success(`${resultMessage}：${fileName}（文本文件）`)
    } else {
      message.success(`${resultMessage}：${fileName}（非文本文件）`)
    }
    await loadTree()
  } catch (error) {
    message.error('代码文件上传失败：' + (error.response?.data?.message || error.message || '未知错误'))
  }
}

async function onMoveFile(payload) {
  const { sourceId, targetId, sourceNode, targetNode, isRootDrop } = payload

  // 处理特殊情况：如果是根目录拖拽但没有具体的 sourceNode 数据
  if (isRootDrop && (!sourceId || sourceId === 'unknown')) {
    console.log('根目录拖拽，但缺少源节点信息，忽略此次移动')
    message.warning('拖拽移动需要明确的源文件信息')
    return
  }

  // 验证必要的参数
  if (!sourceId || sourceId === 'unknown') {
    console.error('移动失败：缺少源文件ID')
    message.error('移动失败：缺少源文件信息')
    return
  }

  try {
    // 处理根目录的情况：targetId 为 '0' 时传递 null 表示移动到根目录
    const newParentId = targetId === '0' ? null : parseInt(targetId)

    console.log('移动文件参数:', {
      sourceId: parseInt(sourceId),
      targetId,
      newParentId,
      isRootMove: targetId === '0'
    })

    await moveTemplateFile({
      id: parseInt(sourceId),
      newParentId: newParentId
    })

    const targetName = targetId === '0' ? '根目录' : (targetNode?.fileName || '未知目录')
    const sourceName = sourceNode?.fileName || sourceNode?.label || '未知文件'
    message.success(`已将 "${sourceName}" 移动到 "${targetName}"`)

    // 保存当前正在编辑的文件信息
    const wasCurrentFile = currentFileId.value === String(sourceId)
    const currentContent = currentFileContent.value
    const currentName = currentFileName.value

    // 重新加载文件树
    await loadTree()

    // 如果移动的是当前文件，重新选择它以保持编辑状态
    if (wasCurrentFile) {
      // 重新设置当前文件信息（文件ID没有变化，只是位置变了）
      currentFileId.value = String(sourceId)
      currentFileName.value = currentName
      currentFileContent.value = currentContent
      templateFileStore.setCurrentFileContent(currentContent)

      // 更新文件树选中状态
      currentFile.value = String(sourceId)
    }
  } catch (error) {
    console.error('移动失败:', error)

    // 检查是否是404错误（接口未实现）
    if (error.response?.status === 404) {
      message.error('移动功能暂未实现，请联系管理员添加后端接口')
    } else {
      message.error('移动失败: ' + (error.response?.data?.message || error.message || '未知错误'))
    }
  }
}

// 调试：打印treeData变化
watch(treeData, (val) => {
}, { deep: true })

// 变量相关事件处理
async function onAddVariable(variable) {
  try {
    await addTemplateVariable({
      templateId: parseInt(route.params.id),
      name: variable.name,
      variableType: variable.variableType || 'string',
      description: variable.description,
      defaultValue: variable.defaultValue || '',
      isRequired: variable.isRequired ? 1 : 0,
      sort: templateVariables.value.length + 1
    })
    await loadVariables()
  } catch (error) {
    message.error('变量添加失败: ' + (error.response?.data?.message || error.message || '未知错误'))
  }
}

async function onEditVariable(variable) {
  try {
    await editTemplateVariable({
      id: variable.id,
      templateId: parseInt(route.params.id),
      name: variable.name,
      variableType: variable.variableType || 'string',
      description: variable.description,
      defaultValue: variable.defaultValue || '',
      isRequired: variable.isRequired ? 1 : 0,
      sort: variable.sort || 0
    })
    message.success('变量更新成功')
    await loadVariables()
  } catch (error) {
    message.error('变量更新失败: ' + (error.response?.data?.message || error.message || '未知错误'))
  }
}

async function onDeleteVariable(id) {
  try {
    await deleteTemplateVariable({ id })
    message.success('变量删除成功')
    await loadVariables()
  } catch (error) {
    message.error('变量删除失败: ' + (error.response?.data?.message || error.message || '未知错误'))
  }
}

function onInsertVariable(template) {
  // 通过 ref 调用编辑器的插入变量方法
  if (templateEditorRef.value) {
    templateEditorRef.value.insertVariable(template)
  }
}

// 预览事件处理
function onPreview({ fileId, fileName }) {

  // 设置预览文件ID
  previewFileId.value = fileId

  // 从 localStorage 获取测试数据
  const templateId = route.params.id
  const testDataKey = `template_test_data_${templateId}`
  const savedTestData = localStorage.getItem(testDataKey)
  if (savedTestData) {
    try {
      variableValues.value = JSON.parse(savedTestData)
    } catch (e) {
      console.error('解析测试数据失败:', e)
      variableValues.value = {}
    }
  } else {
    variableValues.value = {}
  }

  // 触发预览面板展开（如果已折叠）
  if (templatePreviewRef.value) {
    templatePreviewRef.value.expandPanel()
  }

  message.success(`正在预览: ${fileName}`)
}

// 快速插入变量
function insertQuickVariable(quickVar) {
  const goTemplateVar = `{{.${quickVar.name}}}`
  onInsertVariable(goTemplateVar)
}

// 插入自定义变量
function insertVariable(variableName) {
  const goTemplateVar = `{{.${variableName}}}`
  if (templateEditorRef.value) {
    templateEditorRef.value.insertVariable(goTemplateVar)
  }
}

// 显示函数详情 - 增加延迟显示，避免误触发
function showFunctionDetail(func, event) {
  clearHideTimer()
  clearShowTimer()

  // 延迟显示，只有悬停800ms后才显示详情
  showTimer = setTimeout(() => {
    selectedFunction.value = func

    // 计算面板位置 - 显示在鼠标右下角
    const panelWidth = 300  // 减小面板宽度
    const panelHeight = 200 // 减小面板高度
    const offset = 8

    let left = event.clientX + offset
    let top = event.clientY + offset

    // 边界检查 - 确保面板不超出屏幕
    if (left + panelWidth > window.innerWidth - 10) {
      left = event.clientX - panelWidth - offset // 显示在鼠标左下角
    }
    if (top + panelHeight > window.innerHeight - 10) {
      top = event.clientY - panelHeight - offset // 显示在鼠标右上角
    }

    functionDetailStyle.value = {
      left: `${left}px`,
      top: `${top}px`
    }

    functionDetailVisible.value = true
  }, 800) // 增加到800ms延迟，减少误触发
}

// 清除隐藏计时器
function clearHideTimer() {
  if (hideTimer) {
    clearTimeout(hideTimer)
    hideTimer = null
  }
}

// 清除显示计时器
function clearShowTimer() {
  if (showTimer) {
    clearTimeout(showTimer)
    showTimer = null
  }
}

// 鼠标进入详情面板时的处理
function onDetailPanelEnter() {
  clearHideTimer()
  clearShowTimer()
}

// 隐藏函数详情 - 增加延迟隐藏时间
function hideFunctionDetail() {
  clearShowTimer() // 取消可能的显示操作

  hideTimer = setTimeout(() => {
    functionDetailVisible.value = false
    selectedFunction.value = null
  }, 300) // 增加到300ms延迟，给用户更多时间移动到详情面板
}

// 格式化函数为适合插入的格式
function formatFunction(func) {
  return {
    name: func.name,
    label: func.display_name || func.name,
    code: func.insert_text || `{{ ${func.name} }}`,
    description: func.description
  }
}

// 插入函数
function insertFunction(func) {
  // 点击插入时，立即取消详情面板的显示和计时器
  clearShowTimer()
  clearHideTimer()
  functionDetailVisible.value = false

  let code = func.code

  // 如果代码中包含"变量名"占位符，需要特殊处理
  if (code.includes('变量名')) {
    // 这里可以获取编辑器当前选中的文本，如果没有选中则插入占位符
    // 暂时使用占位符，后续可以通过props传递选中文本
    code = code.replace('变量名', '')
  }

  if (templateEditorRef.value) {
    templateEditorRef.value.insertVariable(code)
  }
}

// 插入Sprig函数
function insertSprigFunction(func) {
  // 点击插入时，立即取消详情面板的显示和计时器
  clearShowTimer()
  clearHideTimer()
  functionDetailVisible.value = false

  const code = func.insert_text || `{{ ${func.name} }}`

  if (templateEditorRef.value) {
    templateEditorRef.value.insertVariable(code)
  }
}

// 插入模板语法
function insertSyntax(syntax) {
  // 点击插入时，立即取消详情面板的显示和计时器
  clearShowTimer()
  clearHideTimer()
  functionDetailVisible.value = false

  const code = syntax.insertText || syntax.syntax

  if (templateEditorRef.value) {
    templateEditorRef.value.insertVariable(code)
  }
}


// 应用测试数据
function onApplyTestData(testData) {
  // 更新本地变量值
  variableValues.value = testData
  message.success('测试数据已应用并保存')
}

// 设置文件生成条件
async function onSetCondition(fileNode) {
  selectedFileForCondition.value = fileNode

  // 重置组件表单为默认值
  if (conditionModalRef.value) {
    conditionModalRef.value.resetForm()
  }

  // 从后端获取当前文件的条件配置
  try {
    const fileId = fileNode.key || fileNode.id
    const res = await getFileCondition(fileId)
    if (res.data && res.data.data && res.data.data.condition) {
      const condition = res.data.data.condition
      // 设置组件表单数据
      if (conditionModalRef.value) {
        conditionModalRef.value.setFormData({
          enabled: condition.enabled || false,
          variableName: condition.variableName || '',
          expectedValue: condition.expectedValue !== undefined ? condition.expectedValue : true,
          description: condition.description || ''
        })
      }
    }
  } catch (error) {
    console.error('获取条件配置失败:', error)
    // 如果获取失败，使用默认值
  }

  showConditionModal.value = true
}

// 处理条件保存（来自组件的事件）
async function handleConditionSave(formData) {
  if (!selectedFileForCondition.value) return

  try {
    const fileId = selectedFileForCondition.value.key || selectedFileForCondition.value.id

    // 构建条件数据
    const conditionData = {
      id: fileId,
      enabled: formData.enabled,
      variableName: formData.variableName,
      expectedValue: formData.expectedValue,
      description: formData.description
    }

    // 调用后端API保存条件设置
    await setFileCondition(conditionData)

    message.success('条件设置已保存')
    showConditionModal.value = false

    // 重新加载文件树以显示条件标识
    await loadTree()
  } catch (error) {
    console.error('保存条件失败:', error)
    message.error('保存条件失败: ' + (error.response?.data?.message || error.message || '未知错误'))
  }
}

// AI助手事件处理
// AI插入代码
function onAIInsertCode(code) {
  if (templateEditorRef.value) {
    templateEditorRef.value.insertVariable(code)
    message.success('AI代码已插入到编辑器')
  } else {
    message.warning('请先选择一个文件进行编辑')
  }
}

// AI添加变量
async function onAIAddVariable(variable) {
  try {
    await addTemplateVariable({
      templateId: parseInt(route.params.id),
      name: variable.name,
      variableType: variable.type || 'string',
      description: variable.description || '',
      defaultValue: variable.defaultValue || '',
      isRequired: variable.required ? 1 : 0,
      sort: templateVariables.value.length + 1
    })
    await loadVariables()
    message.success(`AI建议的变量 "${variable.name}" 已添加`)
  } catch (error) {
    message.error('AI变量添加失败: ' + (error.response?.data?.message || error.message || '未知错误'))
  }
}

// AI创建文件
async function onAICreateFile(fileInfo) {
  try {
    const templateId = route.params.id
    const isDirectory = fileInfo.type === 'folder' ? 1 : 0

    await addTemplateFile({
      templateId,
      fileName: fileInfo.name,
      fileContent: fileInfo.content || '',
      fileSize: (fileInfo.content || '').length,
      isDirectory,
      md5: '',
      sort: 0,
      parentId: fileInfo.parentId || null
    })

    await loadTree()
    message.success(`AI生成的${isDirectory ? '文件夹' : '文件'} "${fileInfo.name}" 已创建`)
  } catch (error) {
    message.error('AI文件创建失败: ' + (error.response?.data?.message || error.message || '未知错误'))
  }
}

// AI应用建议
function onAIApplySuggestion(suggestion) {
  switch (suggestion.type) {
    case 'code-optimization':
      // 代码优化建议
      if (suggestion.code && templateEditorRef.value) {
        templateEditorRef.value.insertVariable(suggestion.code)
        message.success('AI优化建议已应用')
      }
      break
    case 'variable-suggestion':
      // 变量建议
      if (suggestion.variables) {
        suggestion.variables.forEach(variable => {
          onAIAddVariable(variable)
        })
      }
      break
    case 'template-structure':
      // 模板结构建议
      if (suggestion.structure) {
        suggestion.structure.forEach(file => {
          onAICreateFile(file)
        })
      }
      break
    default:
      message.info('AI建议已接收')
  }
}
</script>

<style scoped>
.templates-edit-fullscreen {
  position: fixed;
  inset: 0;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
}






.edit-close-btn {
  margin-left: 16px;
}

.edit-main {
  flex: 1;
  display: flex;
  min-height: 0;
  transition: all 0.3s ease;
}

.editor-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
  transition: all 0.3s ease;
}

/* 文件树隐藏时的样式优化 */
.template-explorer {
  transition: all 0.3s ease;
}

/* 自动保存指示器样式 */
.auto-save-indicator {
  position: fixed;
  top: 80px;
  right: 24px;
  background: #52c41a;
  color: #fff;
  padding: 8px 16px;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(82, 196, 26, 0.3);
  font-size: 14px;
  font-weight: 500;
  z-index: 2000;
  display: flex;
  align-items: center;
}

/* 自动保存指示器动画 */
.auto-save-fade-enter-active,
.auto-save-fade-leave-active {
  transition: all 0.3s ease;
}

.auto-save-fade-enter-from {
  opacity: 0;
  transform: translateY(-10px) scale(0.95);
}

.auto-save-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.95);
}




.variable-section {
  margin-bottom: 16px;
}

.variable-section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 12px;
  color: #666;
  margin-bottom: 8px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}


.empty-variables {
  text-align: center;
  padding: 20px;
  color: #999;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  color: #666;
}

.empty-state {
  text-align: center;
  padding: 20px;
  color: #999;
  font-style: italic;
}

/* 自定义函数详情面板 */
.function-detail-panel {
  position: fixed;
  z-index: 10000;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12), 0 2px 8px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  animation: panelFadeIn 0.2s ease-out;
  border: 1px solid rgba(0, 0, 0, 0.06);
  pointer-events: auto;
  max-width: 350px;
  width: 350px;
  max-height: 500px;
}

.function-detail-panel::before {
  content: '';
  position: absolute;
  top: -6px;
  left: 12px;
  width: 12px;
  height: 12px;
  background: #667eea;
  transform: rotate(45deg);
  border-radius: 2px 0 0 0;
}

@keyframes panelFadeIn {
  from {
    opacity: 0;
    transform: translateY(-8px) scale(0.95);
  }

  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.detail-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 12px 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.detail-title {
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
}

.function-icon {
  font-size: 16px;
  filter: drop-shadow(0 0 4px rgba(255, 255, 255, 0.3));
}

.detail-type {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  padding: 3px 6px;
  border-radius: 10px;
  font-size: 10px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.detail-body {
  padding: 16px;
  overflow-y: auto;
  max-height: 420px;
  /* 给header留出空间 */
}

.detail-description {
  font-size: 13px;
  color: #4a5568;
  line-height: 1.5;
  margin-bottom: 12px;
  padding: 10px 12px;
  background: linear-gradient(135deg, #f7fafc 0%, #edf2f7 100%);
  border-radius: 6px;
  border-left: 3px solid #667eea;
}

.detail-section {
  margin-bottom: 12px;
}

.detail-section:last-child {
  margin-bottom: 0;
}

.section-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.section-icon {
  font-size: 14px;
}

.section-content {
  padding: 10px 14px;
  border-radius: 6px;
  font-size: 13px;
  line-height: 1.5;
}

.code-content {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #4a5568;
}

.insert-preview {
  background: linear-gradient(135deg, #f0fff4 0%, #e6fffa 100%);
  border: 1px solid #9ae6b4;
  color: #2f855a;
  font-weight: 500;
}

/* 参数列表样式 */
.params-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.param-item {
  padding: 12px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  transition: all 0.2s;
}

.param-item:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
}

.param-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.param-name {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-weight: 600;
  color: #1e40af;
  background: #dbeafe;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 13px;
}

.param-type {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  color: #7c3aed;
  background: #ede9fe;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 12px;
  font-weight: 500;
}

.param-required {
  background: #fee2e2;
  color: #dc2626;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
}

.param-optional {
  background: #f3f4f6;
  color: #6b7280;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
}

.param-description {
  color: #4b5563;
  font-size: 13px;
  line-height: 1.4;
  margin-bottom: 4px;
}

.param-default {
  color: #6b7280;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.param-default code {
  background: #e5e7eb;
  color: #374151;
  padding: 1px 4px;
  border-radius: 2px;
  font-size: 11px;
}

.empty-text {
  font-size: 12px;
  margin-bottom: 8px;
}

/* 变量管理弹框样式 */
.variable-manager-header {
  display: flex;
  align-items: center;
  width: 100%;
}

.modal-title {
  font-size: 1.1rem;
  font-weight: bold;
  color: #333;
}

.variable-manager-content {
  height: calc(80vh - 120px);
  overflow: hidden;
}

.variable-manager-content .variable-manager {
  height: 100%;
  overflow: hidden;
}

</style>
