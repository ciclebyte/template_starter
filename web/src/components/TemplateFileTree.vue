<template>
  <div class="template-explorer" :style="{ width: `${panelWidth}px` }" @contextmenu="onTreeAreaContextMenu">
    <div class="explorer-title">模板资源</div>
    <div class="explorer-container">
      <NTree
        :data="treeDataComputed"
        :selected-keys="[currentFile]"
        :render-label="renderLabel"
        :node-props="nodeProps"
        :render-switcher-icon="renderSwitcherIcon"
        @update:selected-keys="onSelectFile"
        @update:expanded-keys="updatePrefixWithExpanded"
      />
      <div v-if="!treeData || treeData.length === 0" style="padding: 32px; color: #888; text-align: center; user-select: none; cursor: context-menu;" @contextmenu="onTreeAreaContextMenu">暂无数据（右键新建）</div>
    </div>
    <n-dropdown
      to="body"
      trigger="manual"
      :show="showDropdown"
      :options="dropdownOptions"
      :x="dropdownX"
      :y="dropdownY"
      @select="handleDropdownSelect"
      @clickoutside="handleDropdownClickoutside"
      class="tree-dropdown-menu"
    />
    <input
      type="file"
      ref="fileInput"
      accept=".zip"
      style="display: none"
      @change="handleFileSelect"
    />
    <input
      type="file"
      ref="codeFileInput"
      accept=".js,.ts,.py,.go,.java,.c,.cpp,.json,.md,.txt,.html,.css,.sh,.php,.rb,.rs,.cs,.xml,.yml,.yaml,.ini,.cfg,.conf,.log,.vue,.jsx,.tsx,.less,.scss,.sass,.bat,.ps1,.swift,.dart,.r,.m,.mm,.pl,.lua,.groovy,.gradle,.pom,.lock,.toml,.env,.gitignore,.dockerfile,.makefile,.cmake,.rst,.tex,.bib,.wiki,.adoc,.asciidoc"
      style="display: none"
      @change="handleCodeFileSelect"
    />
    <!-- 拖拽调整宽度的分隔条 -->
    <div 
      class="resize-handle" 
      @mousedown="startResize"
      :class="{ 'is-resizing': isResizing }"
    ></div>
  </div>
</template>

<script setup>
import { ref, watch, h, computed, onMounted, onUnmounted } from 'vue'
import { NTree, useMessage, NIcon } from 'naive-ui'
import { ChevronForward, FileTrayFullOutline, Folder, FolderOpenOutline, Trash, CreateOutline as Edit } from '@vicons/ionicons5'

const props = defineProps({
  treeData: {
    type: Array,
    default: () => []
  },
  currentFile: {
    type: [String, Number],
    default: ''
  }
})
const emit = defineEmits(['select', 'reload', 'update:treeData', 'rename', 'uploadZip', 'uploadCodeFile'])

const showDropdown = ref(false)
const dropdownOptions = ref([
  { label: '新增文件', key: 'addFile', icon: () => h(NIcon, null, { default: () => h(FileTrayFullOutline) }) },
  { label: '新增文件夹', key: 'addFolder', icon: () => h(NIcon, null, { default: () => h(Folder) }) },
  { type: 'divider', key: 'divider1' },
  { label: '删除节点', key: 'deleteNode', icon: () => h(NIcon, null, { default: () => h(Trash) }) }
])
const dropdownX = ref(0)
const dropdownY = ref(0)
const dropdownNode = ref(null)
const addType = ref('file')
const editingNode = ref(null)
const newName = ref('')
const renamingNode = ref(null)
const message = useMessage()
const fileInput = ref(null)
const codeFileInput = ref(null)
const localTreeData = ref(JSON.parse(JSON.stringify(props.treeData)))
// 维护展开状态的映射
const expandedKeys = ref(new Set())

// 拖拽调整宽度相关
const panelWidth = ref(260)
const isResizing = ref(false)
const resizeStartX = ref(0)
const resizeStartWidth = ref(260)

// 拖拽调整宽度功能
function startResize(event) {
  event.preventDefault()
  isResizing.value = true
  resizeStartX.value = event.clientX
  resizeStartWidth.value = panelWidth.value
  
  document.addEventListener('mousemove', handleResize)
  document.addEventListener('mouseup', stopResize)
  document.body.style.cursor = 'ew-resize'
  document.body.style.userSelect = 'none'
}

function handleResize(event) {
  if (!isResizing.value) return
  
  const deltaX = event.clientX - resizeStartX.value
  const newWidth = resizeStartWidth.value + deltaX
  
  // 限制最小和最大宽度
  const minWidth = 200
  const maxWidth = 600
  
  if (newWidth >= minWidth && newWidth <= maxWidth) {
    panelWidth.value = newWidth
  }
}

function stopResize() {
  isResizing.value = false
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
}

// 全局点击事件处理
function handleGlobalClick(event) {
  // 检查是否点击在输入框容器外
  const inputContainer = event.target.closest('.vscode-input-container')
  const isTreeNode = event.target.closest('.n-tree-node')
  
  // 如果有正在编辑的节点，且点击在输入框外，则取消编辑
  if ((editingNode.value || renamingNode.value) && !inputContainer) {
    cancelAddNode()
  }
}

watch(() => props.treeData, (val) => {
  localTreeData.value = JSON.parse(JSON.stringify(val))
}, { deep: true })

// 组件挂载时添加全局点击监听
onMounted(() => {
  document.addEventListener('click', handleGlobalClick, true)
})

// 组件卸载时移除全局点击监听
onUnmounted(() => {
  document.removeEventListener('click', handleGlobalClick, true)
  // 清理拖拽相关事件监听
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
})

// 计算属性，确保展开状态变化时重新渲染
const treeDataComputed = computed(() => {
  return treeToNaive(localTreeData.value)
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
    
    return {
      label: node.isEditing === true ? undefined : (node.fileName || node.label),
      key: nodeKey,
      isLeaf: node.isDirectory === 0,
      isEditing: node.isEditing === true,
      filePath: node.filePath,
      fileName: node.fileName,
      isDirectory: node.isDirectory === 1,
      prefix: node.isDirectory
        ? () => h(NIcon, null, { 
            default: () => h(isExpanded ? FolderOpenOutline : Folder)
          })
        : () => h(NIcon, null, { default: () => h(FileTrayFullOutline) }),
      children: node.children ? treeToNaive(node.children) : []
    }
  })
}

// 更新展开/闭合时的图标
function updatePrefixWithExpanded(
  keys,
  _option,
  meta
) {
  if (!meta.node) return
  
  // 只对目录节点进行图标更新
  if (!meta.node.isDirectory) return
  
  const nodeKey = String(meta.node.key)
  
  switch (meta.action) {
    case 'expand':
      expandedKeys.value.add(nodeKey)
      break
    case 'collapse':
      expandedKeys.value.delete(nodeKey)
      break
  }
}

function onSelectFile(keys) {
  if (keys && keys.length > 0) {
    emit('select', keys[0])
  }
}
function nodeProps({ option }) {
  return {
    onContextmenu(e) {
      e.preventDefault();
      dropdownNode.value = option
      if (option.isLeaf) {
        dropdownOptions.value = [
          { label: '重命名', key: 'renameNode', icon: () => h(NIcon, null, { default: () => h(Edit) }) },
          { type: 'divider', key: 'divider1' },
          { label: '删除节点', key: 'deleteNode', icon: () => h(NIcon, null, { default: () => h(Trash) }) }
        ]
      } else {
        dropdownOptions.value = [
          { label: '新增文件', key: 'addFile', icon: () => h(NIcon, null, { default: () => h(FileTrayFullOutline) }) },
          { label: '新增文件夹', key: 'addFolder', icon: () => h(NIcon, null, { default: () => h(Folder) }) },
          { type: 'divider', key: 'divider1' },
          { label: '上传代码文件', key: 'uploadCodeFile', icon: () => h(NIcon, null, { default: () => h(FileTrayFullOutline) }) },
          { label: '重命名', key: 'renameNode', icon: () => h(NIcon, null, { default: () => h(Edit) }) },
          { type: 'divider', key: 'divider2' },
          { label: '删除节点', key: 'deleteNode', icon: () => h(NIcon, null, { default: () => h(Trash) }) }
        ]
      }
      showDropdown.value = true
      dropdownX.value = e.clientX
      dropdownY.value = e.clientY
    }
  }
}
function onTreeAreaContextMenu(event) {
  if (event.target.closest('.n-tree')) return
  event.preventDefault()
  event.stopPropagation()
  dropdownNode.value = null
  dropdownOptions.value = [
    { label: '新增文件', key: 'addFile', icon: () => h(NIcon, null, { default: () => h(FileTrayFullOutline) }) },
    { label: '新增文件夹', key: 'addFolder', icon: () => h(NIcon, null, { default: () => h(Folder) }) },
    { type: 'divider', key: 'divider1' },
    { label: '上传ZIP包', key: 'uploadZip', icon: () => h(NIcon, null, { default: () => h(Folder) }) },
    { label: '上传代码文件', key: 'uploadCodeFile', icon: () => h(NIcon, null, { default: () => h(FileTrayFullOutline) }) }
  ]
  showDropdown.value = true
  dropdownX.value = event.clientX
  dropdownY.value = event.clientY
}
function addFile() {
  showDropdown.value = false
  addType.value = 'file'
  insertEditingNode('file')
}
function addFolder() {
  showDropdown.value = false
  addType.value = 'folder'
  insertEditingNode('folder')
}
function insertEditingNode(type) {
  const newTreeData = JSON.parse(JSON.stringify(localTreeData.value))
  removeEditingNode(newTreeData)
  newName.value = ''
  const newKey = '__new__' + Date.now() + Math.random().toString(36).slice(2)
  let siblings = []
  const parentId = dropdownNode.value ? String(dropdownNode.value.id || dropdownNode.value.key) : '0'
  if (parentId === '0') {
    siblings = newTreeData
  } else {
    function findParent(list, pid) {
      for (const item of list) {
        if (String(item.id) === pid || String(item.key) === pid) return item
        if (item.children) {
          const found = findParent(item.children, pid)
          if (found) return found
        }
      }
      return null
    }
    const parent = findParent(newTreeData, parentId)
    siblings = parent && parent.children ? parent.children : []
  }
  let maxSort = 0
  siblings.forEach(n => {
    const s = Number(n.sort)
    if (!isNaN(s) && s > maxSort) maxSort = s
  })
  const newNode = {
    key: newKey,
    id: newKey,
    label: '',
    filePath: '',
    isEditing: true,
    isLeaf: type === 'file',
    isDirectory: type === 'folder' ? 1 : 0,
    parentId,
    sort: maxSort + 1,
    children: []
  }
  editingNode.value = newNode
  if (parentId === '0') {
    newTreeData.unshift(newNode)
  } else {
    insertToParent(newTreeData, parentId, newNode)
  }
  localTreeData.value = newTreeData
}
function insertToParent(list, parentId, node) {
  for (const item of list) {
    if (String(item.id) === String(parentId)) {
      if (!item.children) item.children = []
      item.children.unshift(node)
      return true
    }
    if (item.children && insertToParent(item.children, parentId, node)) return true
  }
  return false
}
function removeEditingNode(list) {
  for (let i = list.length - 1; i >= 0; i--) {
    if (list[i].isEditing) {
      // 只处理新增状态，删除临时节点
      list.splice(i, 1)
    } else if (list[i].children) {
      removeEditingNode(list[i].children)
    }
  }
}
function confirmAddNode() {
  if (!newName.value) {
    message.warning('请输入名称')
    return
  }
  
  if (renamingNode.value) {
    // 处理重命名
    emit('rename', { 
      id: renamingNode.value.id || renamingNode.value.key, 
      oldName: renamingNode.value.fileName || renamingNode.value.label,
      newName: newName.value,
      node: renamingNode.value 
    })
    // 清除重命名状态
    const newTreeData = JSON.parse(JSON.stringify(localTreeData.value))
    function clearEditingState(list) {
      for (const item of list) {
        if (String(item.id || item.key) === String(renamingNode.value.id || renamingNode.value.key)) {
          item.isEditing = false
          return true
        }
        if (item.children && clearEditingState(item.children)) {
          return true
        }
      }
      return false
    }
    clearEditingState(newTreeData)
    localTreeData.value = newTreeData
    renamingNode.value = null
  } else {
    // 处理新增
    emit('reload', { name: newName.value, type: addType.value, parentId: editingNode.value.parentId, sort: editingNode.value.sort })
    editingNode.value = null
    removeEditingNode(localTreeData.value)
  }
}

function cancelAddNode() {
  if (renamingNode.value) {
    // 取消重命名：恢复编辑状态，不清除节点
    const newTreeData = JSON.parse(JSON.stringify(localTreeData.value))
    function clearEditingState(list) {
      for (const item of list) {
        if (String(item.id || item.key) === String(renamingNode.value.id || renamingNode.value.key)) {
          item.isEditing = false
          return true
        }
        if (item.children && clearEditingState(item.children)) {
          return true
        }
      }
      return false
    }
    clearEditingState(newTreeData)
    localTreeData.value = newTreeData
    renamingNode.value = null
  } else {
    // 取消新增：删除临时节点
    editingNode.value = null
    removeEditingNode(localTreeData.value)
  }
}
function handleDropdownSelect(key) {
  showDropdown.value = false
  if (key === 'addFile') addFile()
  else if (key === 'addFolder') addFolder()
  else if (key === 'deleteNode') deleteNode()
  else if (key === 'renameNode') renameNode()
  else if (key === 'uploadZip') uploadZip()
  else if (key === 'uploadCodeFile') uploadCodeFile()
}
function handleDropdownClickoutside() {
  showDropdown.value = false
}
function deleteNode() {
  if (!dropdownNode.value) return
  const id = String(dropdownNode.value.id || dropdownNode.value.key)
  function recursiveDelete(list) {
    for (let i = list.length - 1; i >= 0; i--) {
      if (String(list[i].id || list[i].key) === id) {
        list.splice(i, 1)
      } else if (list[i].children) {
        recursiveDelete(list[i].children)
      }
    }
  }
  const newTree = JSON.parse(JSON.stringify(localTreeData.value))
  recursiveDelete(newTree)
  localTreeData.value = newTree
  emit('reload', { type: 'delete', id })
}

function renameNode() {
  if (!dropdownNode.value) return
  const node = dropdownNode.value
  const oldName = node.fileName || node.label
  newName.value = oldName
  renamingNode.value = node
  
  // 设置节点为编辑状态
  const newTreeData = JSON.parse(JSON.stringify(localTreeData.value))
  function setEditingState(list) {
    for (const item of list) {
      if (String(item.id || item.key) === String(node.id || node.key)) {
        item.isEditing = true
        return true
      }
      if (item.children && setEditingState(item.children)) {
        return true
      }
    }
    return false
  }
  setEditingState(newTreeData)
  localTreeData.value = newTreeData
}

function uploadZip() {
  fileInput.value.click()
}

function handleFileSelect({ target }) {
  const file = target.files[0]
  if (!file) return
  
  // 文件验证
  if (!file.name.endsWith('.zip')) {
    message.error('请选择ZIP格式的文件')
    target.value = ''
    return
  }
  
  if (file.size > 1024 * 1024) {
    message.error('文件大小不能超过1MB')
    target.value = ''
    return
  }
  
  // 发送上传事件给父组件
  emit('uploadZip', { file, parentId: dropdownNode.value ? String(dropdownNode.value.id || dropdownNode.value.key) : '0' })
  
  // 清空文件输入
  target.value = ''
}

function uploadCodeFile() {
  codeFileInput.value.value = '' // reset
  codeFileInput.value.click()
}

function handleCodeFileSelect({ target }) {
  const file = target.files[0]
  if (!file) return
  const parentId = dropdownNode.value ? String(dropdownNode.value.id || dropdownNode.value.key) : undefined
  emit('uploadCodeFile', { file, parentId })
  target.value = ''
}

function renderLabel({ option }) {
  if (option.isEditing === true) {
    const isRenaming = renamingNode.value && String(option.id || option.key) === String(renamingNode.value.id || renamingNode.value.key)
    const placeholder = isRenaming ? '请输入新名称' : '请输入名称'
    
    return h(
      'div',
      { 
        class: 'vscode-input-container',
        style: 'display:flex;align-items:center;gap:6px;padding:2px 0;' 
      },
      [
        h('input', {
          class: isRenaming ? 'vscode-input rename' : 'vscode-input create',
          value: newName.value,
          autofocus: true,
          placeholder: placeholder,
          onInput: e => (newName.value = e.target.value),
          onKeydown: e => {
            if (e.key === 'Enter') confirmAddNode()
            if (e.key === 'Escape') cancelAddNode()
          }
        }),
        h(
          'button',
          { 
            class: 'vscode-action-btn confirm',
            onClick: () => confirmAddNode(),
            title: '确认'
          },
          '✓'
        ),
        h(
          'button',
          { 
            class: 'vscode-action-btn cancel',
            onClick: () => cancelAddNode(),
            title: '取消'
          },
          '✕'
        )
      ]
    )
  }
  return option.fileName || option.label
}
const renderSwitcherIcon = () =>
  h(NIcon, null, { default: () => h(ChevronForward) })
</script>

<style scoped>
.template-explorer {
  min-width: 200px;
  max-width: 600px;
  background: #fff;
  border-right: 1px solid #e0e0e0;
  padding: 24px 12px 0 12px;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  position: relative;
  flex-shrink: 0;
}
.explorer-title {
  font-weight: bold;
  margin-bottom: 16px;
  color: #333;
  flex-shrink: 0;
}
.explorer-container {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  min-height: 0;
  padding-right: 4px;
}
.explorer-container::-webkit-scrollbar {
  width: 6px;
}
.explorer-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}
.explorer-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}
.explorer-container::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
/* Firefox 滚动条样式 */
.explorer-container {
  scrollbar-width: thin;
  scrollbar-color: #c1c1c1 #f1f1f1;
}
.tree-dropdown-menu {
  z-index: 2147483647 !important;
}

/* VSCode 风格的输入框和按钮样式 */
:deep(.vscode-input-container) {
  background: transparent;
  border-radius: 3px;
  padding: 1px 2px;
}

:deep(.vscode-input) {
  width: 140px;
  height: 22px;
  padding: 2px 6px;
  font-size: 13px;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background: #ffffff;
  border: 1px solid #cccccc;
  border-radius: 2px;
  outline: none;
  color: #333333;
  transition: all 0.15s ease;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.05);
}

:deep(.vscode-input:focus) {
  border-color: #007acc;
  box-shadow: 0 0 0 1px #007acc;
  background: #ffffff;
}

:deep(.vscode-input.create) {
  border-color: #28a745;
}

:deep(.vscode-input.create:focus) {
  border-color: #28a745;
  box-shadow: 0 0 0 1px #28a745;
}

:deep(.vscode-input.rename) {
  border-color: #ffc107;
  background: #fffbf0;
}

:deep(.vscode-input.rename:focus) {
  border-color: #ffc107;
  box-shadow: 0 0 0 1px #ffc107;
  background: #ffffff;
}

:deep(.vscode-input::placeholder) {
  color: #999999;
  font-style: italic;
}

:deep(.vscode-action-btn) {
  width: 20px;
  height: 20px;
  padding: 0;
  margin: 0;
  border: 1px solid #cccccc;
  border-radius: 2px;
  background: #ffffff;
  color: #333333;
  font-size: 12px;
  font-weight: bold;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
  line-height: 1;
}

:deep(.vscode-action-btn:hover) {
  background: #f3f3f3;
  border-color: #adadad;
}

:deep(.vscode-action-btn:active) {
  background: #e5e5e5;
  transform: translateY(1px);
}

:deep(.vscode-action-btn.confirm) {
  color: #28a745;
  border-color: #28a745;
}

:deep(.vscode-action-btn.confirm:hover) {
  background: #f8fff9;
  border-color: #1e7e34;
}

:deep(.vscode-action-btn.cancel) {
  color: #dc3545;
  border-color: #dc3545;
}

:deep(.vscode-action-btn.cancel:hover) {
  background: #fff5f5;
  border-color: #c82333;
}

/* 拖拽调整分隔条样式 */
.resize-handle {
  position: absolute;
  top: 0;
  right: 0;
  width: 4px;
  height: 100%;
  background: transparent;
  cursor: ew-resize;
  z-index: 10;
  transition: background-color 0.2s ease;
}

.resize-handle:hover {
  background: rgba(0, 123, 204, 0.3);
}

.resize-handle.is-resizing {
  background: rgba(0, 123, 204, 0.6);
}

/* 拖拽时的全局样式 */
.resize-handle:active,
.resize-handle.is-resizing {
  background: rgba(0, 123, 204, 0.6);
}
</style> 