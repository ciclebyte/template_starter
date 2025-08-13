<template>
  <div class="templates-edit-fullscreen">
    <EditHeader :is-variable-panel-open="isVariablePanelOpen" @toggle-variable-panel="toggleVariablePanel"
      @show-variable-manager="showVariableManager = true" @close-edit="closeEdit" />

    <!-- 变量插入面板 -->
    <VariablePanel :is-open="isVariablePanelOpen" :template-variables="templateVariables"
      :template-syntax-categories="templateSyntaxCategories" :builtin-function-categories="builtinFunctionCategories"
      :sprig-function-categories="sprigFunctionCategories" :loading-functions="loadingFunctions"
      :loading-sprig-functions="loadingSprigFunctions" :quick-variables="quickVariables" @insert-syntax="insertSyntax"
      @insert-function="insertFunction" @insert-sprig-function="insertSprigFunction" @insert-variable="insertVariable"
      @show-variable-manager="showVariableManager = true" @update:height="variablePanelHeight = $event" />

    <div class="edit-main">
      <!-- 左侧：模板资源管理器 -->
      <TemplateExplorer v-model:treeData="treeData" :currentFile="currentFile" @select="onSelectFile"
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
import TemplatePreview from './components/TemplatePreview.vue'
import AIAssistant from './components/AIAssistant.vue'
import AISDKPanel from './components/AISDKPanel.vue'
import EditHeader from './components/EditHeader.vue'
import ConditionModal from './components/ConditionModal.vue'
import VariablePanel from './components/VariablePanel.vue'
import { useTemplateFileStore } from '@/stores/templateFileStore'
import { useMessage, NIcon, NTag, NButton, NSpin, NForm, NFormItem, NSwitch, NSelect, NRadioGroup, NRadio, NInput } from 'naive-ui'
import { ChevronDown, ChevronUp, Add, Settings, Pricetag } from '@vicons/ionicons5'

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

// Go模板语法数据
const templateSyntaxCategories = ref([
  {
    name: '条件语句',
    syntaxes: [
      {
        name: 'if',
        display_name: 'if 条件',
        description: '条件判断语句',
        syntax: '{{if .condition}}...{{end}}',
        usage: '当条件为真时执行内容。条件可以是变量、比较表达式或函数调用。',
        example: '{{if .isEnabled}}启用状态{{end}}',
        insertText: '{{if .condition}}\n  content\n{{end}}',
        params: [
          { name: 'condition', type: 'bool', required: true, description: '条件表达式' }
        ]
      },
      {
        name: 'if-else',
        display_name: 'if-else 条件',
        description: '条件判断语句（带else分支）',
        syntax: '{{if .condition}}...{{else}}...{{end}}',
        usage: '当条件为真时执行第一个分支，否则执行else分支。',
        example: '{{if .isEnabled}}启用{{else}}禁用{{end}}',
        insertText: '{{if .condition}}\n  true branch\n{{else}}\n  false branch\n{{end}}',
        params: [
          { name: 'condition', type: 'bool', required: true, description: '条件表达式' }
        ]
      },
      {
        name: 'if-else-if',
        display_name: 'if-else if 条件',
        description: '多重条件判断语句',
        syntax: '{{if .condition1}}...{{else if .condition2}}...{{else}}...{{end}}',
        usage: '按顺序检查多个条件，执行第一个为真的分支。',
        example: '{{if eq .status "active"}}活跃{{else if eq .status "inactive"}}非活跃{{else}}未知{{end}}',
        insertText: '{{if .condition1}}\n  branch1\n{{else if .condition2}}\n  branch2\n{{else}}\n  default\n{{end}}',
        params: [
          { name: 'condition1', type: 'bool', required: true, description: '第一个条件' },
          { name: 'condition2', type: 'bool', required: true, description: '第二个条件' }
        ]
      }
    ]
  },
  {
    name: '循环语句',
    syntaxes: [
      {
        name: 'range',
        display_name: 'range 循环',
        description: '遍历数组、切片或映射',
        syntax: '{{range .items}}...{{end}}',
        usage: '遍历集合中的每个元素。在循环体内，. 代表当前元素。',
        example: '{{range .users}}<p>{{.name}}</p>{{end}}',
        insertText: '{{range .items}}\n  {{.}}\n{{end}}',
        params: [
          { name: 'items', type: 'array|slice|map', required: true, description: '要遍历的集合' }
        ]
      },
      {
        name: 'range-index',
        display_name: 'range 带索引',
        description: '遍历时获取索引和值',
        syntax: '{{range $index, $element := .items}}...{{end}}',
        usage: '遍历集合并获取索引（或键）和对应的值。',
        example: '{{range $i, $user := .users}}<p>{{$i}}: {{$user.name}}</p>{{end}}',
        insertText: '{{range $index, $element := .items}}\n  {{$index}}: {{$element}}\n{{end}}',
        params: [
          { name: 'items', type: 'array|slice|map', required: true, description: '要遍历的集合' }
        ]
      },
      {
        name: 'range-empty',
        display_name: 'range 带空检查',
        description: '遍历集合，支持空值处理',
        syntax: '{{range .items}}...{{else}}...{{end}}',
        usage: '当集合为空时，执行else分支。',
        example: '{{range .users}}<p>{{.name}}</p>{{else}}<p>没有用户</p>{{end}}',
        insertText: '{{range .items}}\n  {{.}}\n{{else}}\n  empty content\n{{end}}',
        params: [
          { name: 'items', type: 'array|slice|map', required: true, description: '要遍历的集合' }
        ]
      }
    ]
  },
  {
    name: '变量操作',
    syntaxes: [
      {
        name: 'with',
        display_name: 'with 作用域',
        description: '设置新的上下文作用域',
        syntax: '{{with .value}}...{{end}}',
        usage: '在with块内，. 代表指定的值。如果值为空，则不执行块内容。',
        example: '{{with .user}}<p>Hello {{.name}}</p>{{end}}',
        insertText: '{{with .value}}\n  {{.}}\n{{end}}',
        params: [
          { name: 'value', type: 'any', required: true, description: '新的上下文值' }
        ]
      },
      {
        name: 'with-else',
        display_name: 'with 带else',
        description: '设置作用域，支持空值处理',
        syntax: '{{with .value}}...{{else}}...{{end}}',
        usage: '当value不为空时执行第一个分支，否则执行else分支。',
        example: '{{with .user}}Hello {{.name}}{{else}}No user{{end}}',
        insertText: '{{with .value}}\n  {{.}}\n{{else}}\n  default content\n{{end}}',
        params: [
          { name: 'value', type: 'any', required: true, description: '上下文值' }
        ]
      },
      {
        name: 'variable',
        display_name: '变量赋值',
        description: '定义和使用变量',
        syntax: '{{$var := .value}}',
        usage: '将值赋给变量，变量名以$开头。变量在整个模板中有效。',
        example: '{{$name := .user.name}}Hello {{$name}}',
        insertText: '{{$var := .value}}',
        params: [
          { name: 'var', type: 'string', required: true, description: '变量名（以$开头）' },
          { name: 'value', type: 'any', required: true, description: '变量值' }
        ]
      }
    ]
  },
  {
    name: '输出控制',
    syntaxes: [
      {
        name: 'printf',
        display_name: 'printf 格式化',
        description: '格式化输出',
        syntax: '{{printf "%s: %d" .name .count}}',
        usage: '使用格式化字符串输出内容，类似C语言的printf函数。',
        example: '{{printf "用户: %s, 年龄: %d" .name .age}}',
        insertText: '{{printf "%s" .value}}',
        params: [
          { name: 'format', type: 'string', required: true, description: '格式化字符串' },
          { name: 'args', type: 'any...', required: false, description: '格式化参数' }
        ]
      },
      {
        name: 'print',
        display_name: 'print 输出',
        description: '简单输出（空格分隔）',
        syntax: '{{print .value1 .value2}}',
        usage: '输出多个值，用空格分隔。',
        example: '{{print "Hello" .name "!"}}',
        insertText: '{{print .value}}',
        params: [
          { name: 'values', type: 'any...', required: true, description: '要输出的值' }
        ]
      },
      {
        name: 'println',
        display_name: 'println 输出',
        description: '输出并换行',
        syntax: '{{println .value}}',
        usage: '输出值并在末尾添加换行符。',
        example: '{{println "Line 1"}}{{println "Line 2"}}',
        insertText: '{{println .value}}',
        params: [
          { name: 'values', type: 'any...', required: true, description: '要输出的值' }
        ]
      }
    ]
  },
  {
    name: '模板控制',
    syntaxes: [
      {
        name: 'template',
        display_name: 'template 引用',
        description: '引用其他模板',
        syntax: '{{template "name" .}}',
        usage: '调用已定义的模板，传递当前上下文或指定数据。',
        example: '{{template "header" .}}{{template "content" .user}}',
        insertText: '{{template "templateName" .}}',
        params: [
          { name: 'name', type: 'string', required: true, description: '模板名称' },
          { name: 'data', type: 'any', required: false, description: '传递给模板的数据' }
        ]
      },
      {
        name: 'define',
        display_name: 'define 定义',
        description: '定义可重用的模板块',
        syntax: '{{define "name"}}...{{end}}',
        usage: '定义一个命名的模板块，可以被template调用。',
        example: '{{define "header"}}<h1>{{.title}}</h1>{{end}}',
        insertText: '{{define "templateName"}}\n  content\n{{end}}',
        params: [
          { name: 'name', type: 'string', required: true, description: '模板名称' }
        ]
      },
      {
        name: 'block',
        display_name: 'block 块定义',
        description: '定义可覆盖的默认块',
        syntax: '{{block "name" .}}...{{end}}',
        usage: '定义一个可以被子模板覆盖的默认内容块。',
        example: '{{block "content" .}}<p>默认内容</p>{{end}}',
        insertText: '{{block "blockName" .}}\n  default content\n{{end}}',
        params: [
          { name: 'name', type: 'string', required: true, description: '块名称' },
          { name: 'data', type: 'any', required: false, description: '传递给块的数据' }
        ]
      }
    ]
  }
])


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


onMounted(async () => {
  await loadTree()
  await loadVariables()
  await loadBuiltinFunctions()
  await loadSprigFunctions()
  loadTestData()
})

// 组件卸载时清理事件监听器
onUnmounted(() => {
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


/* 变量面板样式 */
.variable-panel {
  background: #fff;
  border-bottom: 1px solid #e0e0e0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);
  overflow: hidden;
  position: relative;
}

/* 变量面板拖拽调整样式 */
.variable-panel-resize-handle {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 8px;
  background: transparent;
  cursor: ns-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

.variable-panel-resize-handle:hover {
  background: rgba(24, 160, 88, 0.1);
}

.variable-panel-resize-handle.resizing {
  background: rgba(24, 160, 88, 0.2);
}

.resize-handle-dots {
  display: flex;
  gap: 3px;
  align-items: center;
}

.resize-handle-dots .dot {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: #ccc;
  transition: background-color 0.2s;
}

.variable-panel-resize-handle:hover .dot {
  background: #18a058;
}

.variable-panel-resize-handle.resizing .dot {
  background: #18a058;
}

.variable-tabs {
  height: calc(100% - 8px);
  /* 为拖拽手柄留出空间 */
  display: flex;
  flex-direction: column;
}

.tab-header {
  display: flex;
  background: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.tab-item {
  padding: 12px 20px;
  text-align: center;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  transition: all 0.2s;
  border-bottom: 2px solid transparent;
  white-space: nowrap;
}

.tab-item:hover {
  background: #e9ecef;
  color: #333;
}

.tab-item.active {
  color: #18a058;
  border-bottom-color: #18a058;
  background: #fff;
}

.tab-content {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  /* 移除固定最大高度，让内容根据面板高度动态调整 */
}

.function-categories {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.category-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.category-label {
  width: 140px;
  min-width: 140px;
  font-size: 12px;
  color: #666;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  flex-shrink: 0;
}

.category-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  flex: 1;
}

.edit-close-btn {
  margin-left: 16px;
}

.edit-main {
  flex: 1;
  display: flex;
  min-height: 0;
}

.editor-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
}

.editor-title {
  font-weight: bold;
  margin-bottom: 16px;
  color: #333;
}

.editor-tabs {
  margin-bottom: 4px;
}

.monaco-editor-container {
  flex: 1;
  min-height: 400px;
  background: #1e1e1e;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);
  overflow: hidden;
}

.custom-menu {
  position: fixed;
  z-index: 9999;
  background: #fff;
  border: 1px solid #eee;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border-radius: 6px;
  min-width: 120px;
  padding: 4px 0;
}

.menu-item {
  padding: 8px 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.menu-item:hover {
  background: #f5f5f5;
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

.variable-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  max-width: 100%;
}

.variable-tag {
  cursor: pointer;
  font-size: 12px;
  padding: 6px 12px;
  border-radius: 4px;
  transition: all 0.2s;
  user-select: none;
  border: none;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 0;
}

.variable-tag.builtin {
  background: #18a058;
  color: #fff;
}

.variable-tag.builtin:hover {
  background: #16a34a;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(24, 160, 88, 0.3);
}

.variable-tag.text {
  background: #1890ff;
  color: #fff;
}

.variable-tag.text:hover {
  background: #1677ff;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
}

.variable-tag.conditional {
  background: #f0a020;
  color: #fff;
}

.variable-tag.conditional:hover {
  background: #e69500;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(240, 160, 32, 0.3);
}

.variable-tag.function {
  background: #722ed1;
  color: #fff;
}

.variable-tag.function:hover {
  background: #531dab;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(114, 46, 209, 0.3);
}

.variable-tag.function.sprig {
  background: #38b2ac;
  color: #fff;
}

.variable-tag.function.sprig:hover {
  background: #319795;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(56, 178, 172, 0.3);
}

.variable-tag.syntax {
  background: #ff7875;
  color: #fff;
}

.variable-tag.syntax:hover {
  background: #ff4d4f;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(255, 120, 117, 0.3);
}

/* 新变量类型样式 */
.variable-tag.string {
  background: #1890ff;
  color: #fff;
}

.variable-tag.string:hover {
  background: #1677ff;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
}

.variable-tag.number {
  background: #52c41a;
  color: #fff;
}

.variable-tag.number:hover {
  background: #389e0d;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(82, 196, 26, 0.3);
}

.variable-tag.boolean {
  background: #fa541c;
  color: #fff;
}

.variable-tag.boolean:hover {
  background: #d4380d;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(250, 84, 28, 0.3);
}

.variable-tag.list {
  background: #722ed1;
  color: #fff;
}

.variable-tag.list:hover {
  background: #531dab;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(114, 46, 209, 0.3);
}

.variable-tag.object {
  background: #13c2c2;
  color: #fff;
}

.variable-tag.object:hover {
  background: #08979c;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(19, 194, 194, 0.3);
}

/* 变量类型标识 */
.variable-type-badge {
  margin-left: 8px;
  font-size: 10px;
  opacity: 0.8;
  padding: 2px 4px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
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

/* 条件设置弹框样式 */
</style>