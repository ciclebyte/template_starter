<template>
  <div class="templates-edit-fullscreen">
    <div class="edit-header">
      <div class="header-left">
        <span class="edit-title">模板编辑</span>
      </div>
      <div class="header-center">
        <!-- 变量插入快捷按钮组 -->
        <div class="variable-insert-group">
          <n-button 
            v-for="quickVar in quickVariables" 
            :key="quickVar.name"
            size="small" 
            class="variable-insert-btn"
            @click="insertQuickVariable(quickVar)"
          >
            {{ quickVar.template }}
          </n-button>
          <n-dropdown 
            :options="variableDropdownOptions" 
            @select="insertVariableFromDropdown"
            trigger="click"
          >
            <n-button size="small" class="variable-more-btn">
              更多变量
              <template #icon>
                <n-icon><ChevronDown /></n-icon>
              </template>
            </n-button>
          </n-dropdown>
        </div>
      </div>
      <div class="edit-actions">
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
      />
      
      <!-- 中间：编辑器区域 -->
      <div class="editor-container">
        <TemplateEditor
          ref="templateEditorRef"
          :filePath="currentFilePath"
          :fileContent="currentFileContent"
          :openedTabs="openedTabs"
          :activeTab="activeTab"
          :fileMap="fileMap"
          @tabChange="onTabChange"
          @tabClose="onTabClose"
          @contentChange="onEditorContentChange"
          @insertVariable="onInsertVariable"
        />
      </div>
      
      <!-- 右侧：预览面板 -->
      <TemplatePreview
        :current-file="currentFileNode"
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
          <n-button type="primary" size="small" @click="addVariable">
            <template #icon>
              <n-icon><Add /></n-icon>
            </template>
            新增变量
          </n-button>
        </div>
      </template>
      
      <div class="variable-manager-content">
        <VariableManager
          ref="variableManagerRef"
          :variables="templateVariables"
          @add="onAddVariable"
          @edit="onEditVariable"
          @delete="onDeleteVariable"
          @insert="onInsertVariable"
        />
      </div>
    </n-modal>
  </div>
</template>

<script setup>
import { useRouter, useRoute } from 'vue-router'
import { ref, onMounted, watch, computed } from 'vue'
import { getTemplateFileTree, addTemplateFile, delTemplateFile, getTemplateFileDetail, getTemplateFileContent, renameTemplateFile, uploadZipFile, uploadCodeFile } from '@/api/templateFiles'
import { listTemplateVariables, addTemplateVariable, editTemplateVariable, deleteTemplateVariable } from '@/api/templateVariables'
import TemplateExplorer from '@/components/TemplateFileTree.vue'
import TemplateEditor from '@/components/TemplateEditor.vue'
import VariableManager from '@/components/VariableManager.vue'
import TemplatePreview from '@/components/TemplatePreview.vue'
import { useTemplateFileStore } from '@/stores/templateFileStore'
import { useMessage, NIcon, NDropdown } from 'naive-ui'
import { ChevronDown, Add, Settings } from '@vicons/ionicons5'

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
const openedTabs = ref([])
const activeTab = ref('')
const fileMap = ref({})
const templateFileStore = useTemplateFileStore()

// 变量相关
const templateVariables = ref([])
const templateEditorRef = ref(null)
const variableManagerRef = ref(null)
const showVariableManager = ref(false)
const variableValues = ref({})
const currentFileNode = ref(null)

// 快速插入变量
const quickVariables = [
  { name: 'project_name', label: '项目名', template: '{{project_name}}' },
  { name: 'author', label: '作者', template: '{{author}}' },
  { name: 'date', label: '日期', template: '{{date}}' },
  { name: 'version', label: '版本', template: '{{version}}' }
]

// 变量下拉菜单选项
const variableDropdownOptions = computed(() => {
  return templateVariables.value.map(v => ({
    label: `${v.name} - ${v.description}`,
    key: v.name,
    template: `{{${v.name}}}`
  }))
})

onMounted(async () => {
  await loadTree()
  await loadVariables()
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

function findNodeByKey(list, key) {
  for (const item of list) {
    if (String(item.key || item.id) === String(key)) return item
    if (item.children) {
      const found = findNodeByKey(item.children, key)
      if (found) return found
    }
  }
  return null
}

async function onSelectFile(key) {
  console.log('点击的 key:', key)
  currentFile.value = key
  const node = findNodeByKey(treeData.value, key)
  console.log('找到的 node:', node)
  currentFileNode.value = node
  if (node && node.isDirectory === 0) {
    try {
      const res = await getTemplateFileContent(key)
      console.log('接口返回:', res)
      console.log('接口返回的 fileContent:', res.data?.data?.fileContent)
      const content = res.data.data.fileContent
      currentFileContent.value = content
      templateFileStore.setCurrentFileContent(content)

      // 自动打开tab
      let tab = openedTabs.value.find(t => t.key === String(key))
      if (!tab) {
        openedTabs.value.push({
          key: String(key),
          name: node.fileName || node.label || String(key),
          content
        })
        tab = openedTabs.value[openedTabs.value.length - 1]
      } else {
        tab.content = content
      }
      activeTab.value = String(key)
      // 调试：打印tab信息
      console.log('新建/切换 tab:', tab)
      console.log('tab.name:', tab?.name)
      // 这里 getLanguage 需要和 TemplateEditor.vue 保持一致
      // 你可以在这里 import getLanguage 或直接写一份
      // 这里只做演示
      if (typeof tab?.name === 'string') {
        let lang = 'plaintext'
        if (tab.name.endsWith('.vue')) lang = 'vue'
        else if (tab.name.endsWith('.js')) lang = 'javascript'
        else if (tab.name.endsWith('.ts')) lang = 'typescript'
        else if (tab.name.endsWith('.json')) lang = 'json'
        else if (tab.name.endsWith('.md')) lang = 'markdown'
        else if (tab.name.endsWith('.html')) lang = 'html'
        else if (tab.name.endsWith('.go')) lang = 'go'
        else if (tab.name.endsWith('.java')) lang = 'java'
        else if (tab.name.endsWith('.py')) lang = 'python'
        else if (tab.name.endsWith('.c')) lang = 'c'
        else if (tab.name.endsWith('.cpp')) lang = 'cpp'
        else if (tab.name.endsWith('.cs')) lang = 'csharp'
        else if (tab.name.endsWith('.php')) lang = 'php'
        else if (tab.name.endsWith('.rb')) lang = 'ruby'
        else if (tab.name.endsWith('.sh')) lang = 'shell'
        else if (tab.name.endsWith('.xml')) lang = 'xml'
        else if (tab.name.endsWith('.yml') || tab.name.endsWith('.yaml')) lang = 'yaml'
        console.log('getLanguage(tab.name):', lang)
      }
    } catch (e) {
      currentFileContent.value = ''
      templateFileStore.setCurrentFileContent('')
    }
  } else {
    currentFileContent.value = ''
    templateFileStore.setCurrentFileContent('')
  }
}

function onTabChange(key) {
  activeTab.value = String(key)
  currentFile.value = String(key)
}

function onTabClose(key) {
  const idx = openedTabs.value.findIndex(tab => tab.key === String(key))
  if (idx !== -1) {
    openedTabs.value.splice(idx, 1)
    if (openedTabs.value.length > 0) {
      const next = openedTabs.value[Math.max(0, idx - 1)]
      activeTab.value = next.key
      currentFile.value = next.key
    }
  }
}

function onEditorContentChange({ key, content }) {
  fileMap.value[key] = content
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
    
    // 如果重命名的文件在打开的标签页中，更新标签页名称
    const tabIndex = openedTabs.value.findIndex(tab => tab.key === String(id))
    if (tabIndex !== -1) {
      openedTabs.value[tabIndex].name = newName.trim()
    }
  } catch (error) {
    console.error('重命名失败:', error)
    message.error('重命名失败: ' + (error.response?.data?.message || error.message || '未知错误'))
  }
}

async function onUploadZip(payload) {
  const { file } = payload
  
  try {
    console.log('开始上传文件:', file)
    const res = await uploadZipFile(route.params.id, file)
    console.log('上传结果:', res)
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

// 调试：打印treeData变化
watch(treeData, (val) => {
  console.log('父组件 treeData 变化:', val)
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
    message.success('变量添加成功')
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

// 快速插入变量
function insertQuickVariable(quickVar) {
  onInsertVariable(quickVar.template)
}

// 从下拉菜单插入变量
function insertVariableFromDropdown(key) {
  const option = variableDropdownOptions.value.find(opt => opt.key === key)
  if (option) {
    onInsertVariable(option.template)
  }
}

// 新增变量（直接打开编辑对话框）
function addVariable() {
  if (variableManagerRef.value) {
    variableManagerRef.value.addVariable()
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
  gap: 24px;
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.edit-title {
  font-size: 1.2rem;
  font-weight: bold;
  color: #18a058;
}

.variable-insert-group {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.variable-insert-btn {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  padding: 4px 8px;
  height: 28px;
  border: 1px solid #e0e0e0;
  background: #fafafa;
  color: #333;
  transition: all 0.2s;
}

.variable-insert-btn:hover {
  border-color: #18a058;
  background: #f0f9ff;
  color: #18a058;
}

.variable-more-btn {
  height: 28px;
  padding: 4px 8px;
  font-size: 12px;
}

.edit-actions {
  display: flex;
  align-items: center;
  gap: 12px;
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

/* 变量管理弹框样式 */
.variable-manager-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
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