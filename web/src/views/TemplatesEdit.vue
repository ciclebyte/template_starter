<template>
  <div class="templates-edit-fullscreen">
    <div class="edit-header">
      <div class="header-left">
        <span class="edit-title">模板编辑</span>
        <!-- 模板变量展开按钮 -->
        <div class="variable-expand-trigger" @click="toggleVariablePanel">
          <span class="trigger-text">模板变量</span>
          <n-icon class="trigger-icon" :class="{ 'rotated': isVariablePanelOpen }">
            <ChevronDown />
          </n-icon>
        </div>
      </div>
      <div class="header-actions">
        <n-button size="small" @click="showVariableManager = true">
          <template #icon>
            <n-icon><Settings /></n-icon>
          </template>
          变量管理
        </n-button>
        <n-button quaternary circle class="edit-close-btn" @click="closeEdit">
          <template #icon>
            <n-icon><svg viewBox="0 0 24 24" width="20" height="20"><path fill="currentColor" d="M18.3 5.71a1 1 0 0 0-1.41 0L12 10.59 7.11 5.7A1 1 0 0 0 5.7 7.11L10.59 12l-4.89 4.89a1 1 0 1 0 1.41 1.41L12 13.41l4.89 4.89a1 1 0 0 0 1.41-1.41L13.41 12l4.89-4.89a1 1 0 0 0 0-1.4z"/></svg></n-icon>
          </template>
        </n-button>
      </div>
    </div>
    
    <!-- 变量插入面板 -->
    <div v-show="isVariablePanelOpen" class="variable-panel">
      <div class="variable-tabs">
        <div class="tab-header">
          <div 
            v-for="tab in variableTabs" 
            :key="tab.key"
            class="tab-item"
            :class="{ active: activeVariableTab === tab.key }"
            @click="activeVariableTab = tab.key"
          >
            {{ tab.label }}
          </div>
        </div>
        
        <!-- 内置函数 Tab -->
        <div v-show="activeVariableTab === 'functions'" class="tab-content">
          <div class="function-categories">
            <div class="category-row">
              <span class="category-label">时间函数</span>
              <div class="category-tags">
                <div 
                  v-for="func in timeFunctions" 
                  :key="func.name"
                  class="variable-tag function"
                  @click="insertFunction(func)"
                  :title="`${func.label} - ${func.description}`"
                >
                  {{ func.label }}
                </div>
              </div>
            </div>
            
            <div class="category-row">
              <span class="category-label">字符串处理</span>
              <div class="category-tags">
                <div 
                  v-for="func in stringFunctions" 
                  :key="func.name"
                  class="variable-tag function"
                  @click="insertFunction(func)"
                  :title="`${func.label} - ${func.description}`"
                >
                  {{ func.label }}
                </div>
              </div>
            </div>
            
            <div class="category-row">
              <span class="category-label">随机值</span>
              <div class="category-tags">
                <div 
                  v-for="func in randomFunctions" 
                  :key="func.name"
                  class="variable-tag function"
                  @click="insertFunction(func)"
                  :title="`${func.label} - ${func.description}`"
                >
                  {{ func.label }}
                </div>
              </div>
            </div>
            
            <div class="category-row">
              <span class="category-label">条件函数</span>
              <div class="category-tags">
                <div 
                  v-for="func in conditionalFunctions" 
                  :key="func.name"
                  class="variable-tag function"
                  @click="insertFunction(func)"
                  :title="`${func.label} - ${func.description}`"
                >
                  {{ func.label }}
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 内置变量 Tab -->
        <div v-show="activeVariableTab === 'builtin'" class="tab-content">
          <div class="variable-section">
            <div class="section-title">文本变量</div>
            <div class="variable-tags">
              <div 
                v-for="variable in quickVariables" 
                :key="variable.name"
                class="variable-tag builtin"
                @click="insertVariable(variable.name)"
                :title="`${variable.name} - ${variable.label}`"
              >
                {{ variable.label }}
              </div>
            </div>
          </div>
        </div>
        
        <!-- 用户变量 Tab -->
        <div v-show="activeVariableTab === 'custom'" class="tab-content">
          <div v-if="textVariables.length > 0" class="variable-section">
            <div class="section-title">文本变量</div>
            <div class="variable-tags">
              <n-tag
                v-for="variable in textVariables"
                :key="variable.id"
                class="variable-tag text"
                @click="insertVariable(variable.name)"
                :title="`${variable.name}${variable.description ? ' - ' + variable.description : ''}`"
              >
                {{ variable.name }}
              </n-tag>
            </div>
          </div>
          
          <div v-if="conditionalVariables.length > 0" class="variable-section">
            <div class="section-title">条件变量</div>
            <div class="variable-tags">
              <n-tag
                v-for="variable in conditionalVariables"
                :key="variable.id"
                class="variable-tag conditional"
                @click="insertVariable(variable.name)"
                :title="`${variable.name}${variable.description ? ' - ' + variable.description : ''}`"
              >
                {{ variable.name }}
              </n-tag>
            </div>
          </div>
          
          <div v-if="templateVariables.length === 0" class="empty-variables">
            <div class="empty-text">暂无自定义变量</div>
            <n-button text type="primary" size="small" @click="showVariableManager = true">
              添加变量
            </n-button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="edit-main">
      <!-- 左侧：模板资源管理器 -->
      <TemplateExplorer
        v-model:treeData="treeData"
        :currentFile="currentFile"
        @select="onSelectFile"
        @reload="onTreeReload"
        @rename="onRenameFile"
        @uploadZip="onUploadZip"
        @uploadCodeFile="onUploadCodeFile"
        @move="onMoveFile"
      />
      
      <!-- 中间：编辑器区域 -->
      <div class="editor-container">
        <TemplateEditor
          ref="templateEditorRef"
          :currentFileName="currentFileName"
          :currentFileId="currentFileId"
          :currentFileContent="currentFileContent"
          @contentChange="onEditorContentChange"
          @insertVariable="onInsertVariable"
          @preview="onPreview"
        />
      </div>
      
      <!-- 右侧：预览面板 -->
      <TemplatePreview
        ref="templatePreviewRef"
        :current-file="currentFileNode"
        :file-id="previewFileId"
        :variables="variableValues"
      />
    </div>



    <!-- 变量管理弹框 -->
    <n-modal 
      v-model:show="showVariableManager" 
      preset="card" 
      style="width: 80vw; height: 80vh; max-width: 1200px;"
      :mask-closable="false"
    >
      <template #header>
        <div class="variable-manager-header">
          <span class="modal-title">变量管理</span>
        </div>
      </template>
      
      <div class="variable-manager-content">
        <VariableManager
          ref="variableManagerRef"
          :variables="templateVariables"
          :template-id="route.params.id"
          @add="onAddVariable"
          @edit="onEditVariable"
          @delete="onDeleteVariable"
          @insert="onInsertVariable"
          @applyTestData="onApplyTestData"
        />
      </div>
    </n-modal>
  </div>
</template>

<script setup>
import { useRouter, useRoute } from 'vue-router'
import { ref, onMounted, onUnmounted, watch, computed, nextTick } from 'vue'
import { getTemplateFileTree, addTemplateFile, delTemplateFile, getTemplateFileDetail, getTemplateFileContent, renameTemplateFile, uploadZipFile, uploadCodeFile, moveTemplateFile } from '@/api/templateFiles'
import { listTemplateVariables, addTemplateVariable, editTemplateVariable, deleteTemplateVariable } from '@/api/templateVariables'
import TemplateExplorer from '@/components/TemplateFileTree.vue'
import TemplateEditor from '@/components/TemplateEditor.vue'
import VariableManager from '@/components/VariableManager.vue'
import TemplatePreview from '@/components/TemplatePreview.vue'
import { useTemplateFileStore } from '@/stores/templateFileStore'
import { useMessage, NIcon, NTag, NButton } from 'naive-ui'
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
const currentFilePath = ref('')
const currentFileContent = ref('')
const currentFileName = ref('')
const currentFileId = ref('')
const templateFileStore = useTemplateFileStore()

// 变量相关
const templateVariables = ref([])
const templateEditorRef = ref(null)
const variableManagerRef = ref(null)
const showVariableManager = ref(false)

const variableValues = ref({})
const currentFileNode = ref(null)
const templatePreviewRef = ref(null)
const previewFileId = ref(null)

// 内置变量列表
const quickVariables = [
  { name: 'ProjectName', label: '项目名称' },
  { name: 'Author', label: '作者' },
  { name: 'PackageName', label: '包名' }
]

// 快捷函数列表
const quickFunctions = [
  { name: 'now', label: '当前时间', code: '{{now}}', description: '返回当前时间' },
  { name: 'date', label: '格式化日期', code: '{{date "2006-01-02"}}', description: '按指定格式返回当前日期' },
  { name: 'lower', label: '转小写', code: '{{lower .变量名}}', description: '将变量转换为小写' },
  { name: 'upper', label: '转大写', code: '{{upper .变量名}}', description: '将变量转换为大写' },
  { name: 'camelcase', label: '驼峰命名', code: '{{camelcase .变量名}}', description: '转换为驼峰命名格式' },
  { name: 'snakecase', label: '下划线命名', code: '{{snakecase .变量名}}', description: '转换为下划线命名格式' },
  { name: 'randInt', label: '随机整数', code: '{{randInt 1 100}}', description: '生成1-100之间的随机整数' },
  { name: 'uuid', label: 'UUID', code: '{{uuid}}', description: '生成UUID' },
  { name: 'default', label: '默认值', code: '{{default "默认值" .变量名}}', description: '如果变量为空则使用默认值' }
]

// 完整的函数分类
const timeFunctions = [
  { name: 'now', label: '当前时间', code: '{{now}}', description: '返回当前时间' },
  { name: 'date', label: '格式化日期', code: '{{date "2006-01-02"}}', description: '按指定格式返回当前日期' },
  { name: 'datetime', label: '格式化时间', code: '{{date "2006-01-02 15:04:05"}}', description: '按指定格式返回当前时间' },
  { name: 'timestamp', label: '时间戳', code: '{{now | unixEpoch}}', description: '返回Unix时间戳' },
  { name: 'year', label: '当前年份', code: '{{date "2006"}}', description: '返回当前年份' },
  { name: 'month', label: '当前月份', code: '{{date "01"}}', description: '返回当前月份' },
  { name: 'day', label: '当前日期', code: '{{date "02"}}', description: '返回当前日期' }
]

const stringFunctions = [
  { name: 'lower', label: '转小写', code: '{{lower .变量名}}', description: '将变量转换为小写' },
  { name: 'upper', label: '转大写', code: '{{upper .变量名}}', description: '将变量转换为大写' },
  { name: 'title', label: '首字母大写', code: '{{title .变量名}}', description: '将变量首字母大写' },
  { name: 'camelcase', label: '驼峰命名', code: '{{camelcase .变量名}}', description: '转换为驼峰命名格式' },
  { name: 'snakecase', label: '下划线命名', code: '{{snakecase .变量名}}', description: '转换为下划线命名格式' },
  { name: 'kebabcase', label: '短横线命名', code: '{{kebabcase .变量名}}', description: '转换为短横线命名格式' },
  { name: 'trim', label: '去除空格', code: '{{trim .变量名}}', description: '去除变量首尾空格' },
  { name: 'trunc', label: '截断字符串', code: '{{trunc 10 .变量名}}', description: '截断字符串到指定长度' }
]

const randomFunctions = [
  { name: 'randInt', label: '随机整数', code: '{{randInt 1 100}}', description: '生成1-100之间的随机整数' },
  { name: 'randAlpha', label: '随机字母', code: '{{randAlpha 10}}', description: '生成10位随机字母' },
  { name: 'randAlphaNum', label: '随机字母数字', code: '{{randAlphaNum 8}}', description: '生成8位随机字母数字' },
  { name: 'randNumeric', label: '随机数字', code: '{{randNumeric 6}}', description: '生成6位随机数字' },
  { name: 'uuid', label: 'UUID', code: '{{uuid}}', description: '生成UUID' }
]

const conditionalFunctions = [
  { name: 'default', label: '默认值', code: '{{default "默认值" .变量名}}', description: '如果变量为空则使用默认值' },
  { name: 'if', label: '条件判断', code: '{{if .条件}}值1{{else}}值2{{end}}', description: '条件判断语句' },
  { name: 'eq', label: '相等判断', code: '{{eq .变量1 .变量2}}', description: '判断两个变量是否相等' },
  { name: 'ne', label: '不等判断', code: '{{ne .变量1 .变量2}}', description: '判断两个变量是否不相等' }
]

// 变量面板状态
const isVariablePanelOpen = ref(false)
const activeVariableTab = ref('functions')

// 变量标签页配置
const variableTabs = [
  { key: 'functions', label: '内置函数' },
  { key: 'builtin', label: '内置变量' },
  { key: 'custom', label: '用户变量' }
]

// 计算属性：按类型分组变量
const textVariables = computed(() => {
  return templateVariables.value.filter(v => v.variableType === 'text' || !v.variableType)
})

const conditionalVariables = computed(() => {
  return templateVariables.value.filter(v => v.variableType === 'conditional')
})



onMounted(async () => {
  await loadTree()
  await loadVariables()
  loadTestData()
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
  const { sourceId, targetId, sourceNode, targetNode } = payload
  
  try {
    await moveTemplateFile({
      id: parseInt(sourceId),
      newParentId: parseInt(targetId)
    })
    
    message.success(`已将 "${sourceNode.fileName}" 移动到 "${targetNode.fileName}" 文件夹`)
    
    // 重新加载文件树
    await loadTree()
    
    // 如果移动的是当前文件，清除当前文件状态
    if (currentFileId.value === String(sourceId)) {
      currentFileContent.value = ''
      currentFileName.value = ''
      currentFileId.value = ''
      templateFileStore.setCurrentFileContent('')
    }
  } catch (error) {
    console.error('移动失败:', error)
    message.error('移动失败: ' + (error.response?.data?.message || error.message || '未知错误'))
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
      variableType: variable.variableType || 'text',
      description: variable.description,
      defaultValue: variable.defaultValue || '',
      isRequired: variable.isRequired ? 1 : 0,
      validationRegex: variable.validationRegex || '',
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
      variableType: variable.variableType || 'text',
      description: variable.description,
      defaultValue: variable.defaultValue || '',
      isRequired: variable.isRequired ? 1 : 0,
      validationRegex: variable.validationRegex || '',
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

// 插入函数
function insertFunction(func) {
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

// 切换变量面板
function toggleVariablePanel() {
  isVariablePanelOpen.value = !isVariablePanelOpen.value
}



// 新增变量（直接打开编辑对话框）
function addVariable() {
  if (variableManagerRef.value) {
    variableManagerRef.value.addVariable()
  }
}

// 应用测试数据
function onApplyTestData(testData) {
  // 更新本地变量值
  variableValues.value = testData
  message.success('测试数据已应用并保存')
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

.edit-header {
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

.edit-title {
  font-size: 1.2rem;
  font-weight: bold;
  color: #18a058;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* 变量展开按钮样式 */
.variable-expand-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;
}

.variable-expand-trigger:hover {
  background: #e9ecef;
  border-color: #18a058;
}

.trigger-text {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.trigger-icon {
  font-size: 16px;
  color: #666;
  transition: transform 0.2s;
}

.trigger-icon.rotated {
  transform: rotate(180deg);
}

/* 变量面板样式 */
.variable-panel {
  background: #fff;
  border-bottom: 1px solid #e0e0e0;
  box-shadow: 0 2px 8px rgba(0,0,0,0.03);
  overflow: hidden;
}

.variable-tabs {
  height: 100%;
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
  max-height: 400px;
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
  min-width: 80px;
  font-size: 12px;
  color: #666;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
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
  box-shadow: 0 2px 8px rgba(0,0,0,0.03);
  overflow: hidden;
}

.custom-menu {
  position: fixed;
  z-index: 9999;
  background: #fff;
  border: 1px solid #eee;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
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

.empty-variables {
  text-align: center;
  padding: 20px;
  color: #999;
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