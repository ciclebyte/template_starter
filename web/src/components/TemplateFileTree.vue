<template>
  <div class="template-explorer" :style="{ width: `${panelWidth}px` }" @contextmenu="onTreeAreaContextMenu">
    <div class="explorer-title">模板资源</div>
    <div class="explorer-container"
         :class="{ 'drag-over-root': isDragOverRoot }"
         ref="explorerContainer"
         @dragover.capture="onContainerDragOver"
         @drop.capture="onContainerDrop"
         @dragenter.capture="onContainerDragEnter"
         @click="onContainerClick">
      <NTree
        :data="treeDataComputed"
        :selected-keys="[currentFile]"
        :render-label="renderLabel"
        :node-props="nodeProps"
        :render-switcher-icon="renderSwitcherIcon"
        @update:selected-keys="onSelectFile"
        @update:expanded-keys="updatePrefixWithExpanded"
        draggable
        @dragstart="onDragStart"
        @drag-enter="onDragEnter"
        @drag-leave="onDragLeave"
        @drag-over="onDragOver"
        @drop="onDrop"
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
const emit = defineEmits(['select', 'reload', 'update:treeData', 'rename', 'uploadZip', 'uploadCodeFile', 'move'])

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

// 拖拽移动相关
const draggedNode = ref(null)
const dragOverNode = ref(null)
const isDragging = ref(false)
const isDragOverRoot = ref(false)
const lastMousePosition = ref({ x: 0, y: 0 })
const explorerContainer = ref(null)

// 全局拖拽事件监听
let globalDragOverHandler = null
let globalDragEndHandler = null

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
  // 检查是否点击在输入框外
  const inputElement = event.target.closest('.vscode-tree-input')
  const isTreeNode = event.target.closest('.n-tree-node')
  
  // 如果有正在编辑的节点，且点击在输入框外，则确认编辑
  if ((editingNode.value || renamingNode.value) && !inputElement) {
    confirmAddNode()
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
  // 清理拖拽状态和监听器
  clearDragState()
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
      children: node.children ? treeToNaive(node.children) : [],
      // 添加拖拽状态的类名
      class: getDragClass(nodeKey)
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

// 检查是否为子节点
function isDescendant(parentNode, childNode) {
  if (!parentNode.children || parentNode.children.length === 0) return false
  
  for (const child of parentNode.children) {
    if (child.key === childNode.key) return true
    if (isDescendant(child, childNode)) return true
  }
  return false
}

// 根据 key 查找节点
function findNodeByKey(nodes, key) {
  if (!nodes || !Array.isArray(nodes)) return null
  
  for (const node of nodes) {
    if (String(node.key) === String(key) || String(node.id) === String(key)) {
      return node
    }
    
    if (node.children && node.children.length > 0) {
      const found = findNodeByKey(node.children, key)
      if (found) return found
    }
  }
  
  return null
}

// 获取拖拽状态的类名
function getDragClass(nodeKey) {
  const classes = []
  
  if (draggedNode.value && String(draggedNode.value.key) === String(nodeKey)) {
    classes.push('dragging')
  }
  
  if (dragOverNode.value && String(dragOverNode.value.key) === String(nodeKey)) {
    classes.push('drag-over')
  }
  
  return classes.join(' ')
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

// 设置全局拖拽监听
function setupGlobalDragListeners() {
  globalDragEndHandler = () => {
    console.log('🔚 globalDragEnd - clearing drag state')
    clearDragState()
  }
  
  document.addEventListener('dragend', globalDragEndHandler)
}

// 清理全局拖拽监听
function clearGlobalDragListeners() {
  if (globalDragEndHandler) {
    document.removeEventListener('dragend', globalDragEndHandler)
    globalDragEndHandler = null
  }
}

// 清理拖拽状态
function clearDragState() {
  draggedNode.value = null
  dragOverNode.value = null
  isDragging.value = false
  isDragOverRoot.value = false
  clearGlobalDragListeners()
}

// NTree 拖拽事件处理
function onDragStart(info) {
  console.log('🚀 drag start:', info)
  console.log('🚀 dragNode:', info.dragNode)
  console.log('🚀 event:', info.event)
  if (info.dragNode) {
    draggedNode.value = info.dragNode
    isDragging.value = true
    console.log('✅ Starting drag operation for:', info.dragNode.fileName || info.dragNode.label)
    console.log('✅ isDragging set to:', isDragging.value)
    
    // 将拖拽数据存储到 DataTransfer 中，以便在容器事件中恢复
    if (info.event && info.event.dataTransfer) {
      info.event.dataTransfer.setData('application/json', JSON.stringify(info.dragNode))
    }
    
    setupGlobalDragListeners()
  } else {
    console.log('❌ No dragNode in info')
  }
}

function onDragEnter(info) {
  console.log('🎯 drag enter:', info)
  console.log('🎯 dragNode:', info.dragNode)
  console.log('🎯 event:', info.event)
  if (info.dragNode && !isDragging.value) {
    draggedNode.value = info.dragNode
    isDragging.value = true
    console.log('✅ Starting drag operation for:', info.dragNode.fileName || info.dragNode.label)
    console.log('✅ isDragging set to:', isDragging.value)
    
    // 同样在 dragEnter 中设置 DataTransfer
    if (info.event && info.event.dataTransfer) {
      info.event.dataTransfer.setData('application/json', JSON.stringify(info.dragNode))
      console.log('✅ DataTransfer set in dragEnter')
    }
    
    setupGlobalDragListeners()
  } else {
    console.log('❌ No dragNode in info or already dragging')
  }
}

function onDragLeave(info) {
  console.log('drag leave:', info)
  // 清除节点拖拽状态
  dragOverNode.value = null
}

function onDragOver(info) {
  console.log('🔄 drag over:', info)
  const { event, node } = info
  
  // 更新鼠标位置
  lastMousePosition.value = { x: event.clientX, y: event.clientY }
  
  // 只有文件夹可以作为目标
  if (node && node.isDirectory) {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'move'
    dragOverNode.value = node
    isDragOverRoot.value = false
    console.log('📁 dragOver - folder target:', node.fileName)
  } else {
    // 如果不是文件夹，清除dragOverNode状态
    dragOverNode.value = null
    console.log('📄 dragOver - file target, checking root area')
    
    // 检查是否在根目录区域
    checkRootAreaFromEvent(event)
  }
}

// 从拖拽事件中检查根目录区域
function checkRootAreaFromEvent(event) {
  if (!isDragging.value) return
  
  const { clientX, clientY } = event
  console.log('🔍 checkRootAreaFromEvent - mouse:', clientX, clientY)
  
  // 检查是否在容器内
  if (explorerContainer.value) {
    const containerRect = explorerContainer.value.getBoundingClientRect()
    
    if (clientX >= containerRect.left && clientX <= containerRect.right &&
        clientY >= containerRect.top && clientY <= containerRect.bottom) {
      
      // 检查鼠标位置的元素
      const elementAtPoint = document.elementFromPoint(clientX, clientY)
      const treeNode = elementAtPoint?.closest('.n-tree-node')
      
      console.log('🎯 checkRootAreaFromEvent - elementAtPoint:', elementAtPoint?.tagName, 'treeNode:', !!treeNode)
      
      if (!treeNode) {
        // 在容器内但不在树节点上，说明在根目录区域
        console.log('✅ checkRootAreaFromEvent - ROOT AREA DETECTED!')
        isDragOverRoot.value = true
        event.preventDefault()
        event.dataTransfer.dropEffect = 'move'
      } else {
        console.log('❌ checkRootAreaFromEvent - on tree node, not root')
        isDragOverRoot.value = false
      }
    } else {
      console.log('❌ checkRootAreaFromEvent - outside container')
      isDragOverRoot.value = false
    }
  }
}

function onDrop(info) {
  console.log('drop:', info)
  const { event, node, dragNode } = info
  
  // 检查是否是根目录拖拽
  if (isDragOverRoot.value || !node) {
    console.log('Dropping to root directory')
    handleRootDrop(dragNode)
    return
  }
  
  if (!dragNode || !node.isDirectory || node.isEditing) {
    console.log('Invalid drop target:', { dragNode, node })
    clearDragState()
    return
  }
  
  const sourceId = dragNode.key
  const targetId = node.key
  
  // 不能移动到自己
  if (sourceId === targetId) {
    console.log('Cannot move to self')
    clearDragState()
    return
  }
  
  // 不能移动到自己的子节点
  if (isDescendant(dragNode, node)) {
    console.log('Cannot move to descendant')
    clearDragState()
    return
  }
  
  console.log('Moving from', sourceId, 'to', targetId)
  
  // 触发移动事件
  emit('move', {
    sourceId,
    targetId,
    sourceNode: dragNode,
    targetNode: node
  })
  
  // 清理状态
  clearDragState()
}

// 处理根目录拖拽
function handleRootDrop(dragNode) {
  if (!dragNode) return
  
  // 使用最后的鼠标位置进行额外检查
  const containerElement = document.querySelector('.explorer-container')
  if (!containerElement) return
  
  const containerRect = containerElement.getBoundingClientRect()
  const { x, y } = lastMousePosition.value
  
  // 确认鼠标在容器内
  if (x >= containerRect.left && x <= containerRect.right &&
      y >= containerRect.top && y <= containerRect.bottom) {
    
    // 再次检查是否在树节点上
    const elementAtPoint = document.elementFromPoint(x, y)
    const treeNode = elementAtPoint?.closest('.n-tree-node')
    
    if (!treeNode) {
      console.log('Confirmed root drop at position:', x, y)
      
      // 移动到根目录
      emit('move', {
        sourceId: dragNode.key,
        targetId: '0',
        sourceNode: dragNode,
        targetNode: { key: '0', isDirectory: true }
      })
    } else {
      console.log('Drop cancelled: found tree node at position')
    }
  } else {
    console.log('Drop cancelled: mouse outside container')
  }
  
  // 清理状态
  clearDragState()
}

// 容器级别的拖拽事件处理（作为NTree事件的补充）
function onContainerDragOver(event) {
  console.log('🌊 containerDragOver - ALWAYS TRIGGERED')
  console.log('🌊 isDragging.value:', isDragging.value)
  console.log('🌊 event.target:', event.target)
  console.log('🌊 event type:', event.type)
  
  // 检查是否有拖拽数据
  if (event.dataTransfer && event.dataTransfer.types.length > 0) {
    console.log('🌊 dataTransfer types:', event.dataTransfer.types)
    
    // 如果我们检测到拖拽事件但没有拖拽状态，尝试设置它
    if (!isDragging.value) {
      console.log('🌊 containerDragOver - setting dragging state from container event')
      isDragging.value = true
      
      // 尝试从数据传输中获取拖拽的节点信息
      // 这是一个备用方案，当NTree事件不工作时
      try {
        const dragData = event.dataTransfer.getData('application/json')
        if (dragData) {
          const nodeData = JSON.parse(dragData)
          draggedNode.value = nodeData
          console.log('🌊 recovered drag node from dataTransfer:', nodeData)
        }
      } catch (e) {
        console.log('🌊 could not parse drag data, continuing anyway')
      }
      
      // 如果还是没有拖拽节点，但有拖拽数据，创建一个备用状态
      if (!draggedNode.value && event.dataTransfer.types.includes('application/json')) {
        console.log('🌊 containerDragOver - creating backup drag state')
        setupGlobalDragListeners()
      }
    }
  }
  
  // 如果没有拖拽状态但有 DataTransfer，可能是内部拖拽，继续处理
  if (!isDragging.value && (!event.dataTransfer || event.dataTransfer.types.length === 0)) {
    console.log('🌊 containerDragOver - no dragging state and no dataTransfer, exiting')
    return
  }
  
  if (!isDragging.value) {
    console.log('🌊 containerDragOver - no internal dragging state but has dataTransfer, continuing')
  }
  
  console.log('🌊 containerDragOver - capture event with dragging state')
  
  // 检查是否在树节点上
  const treeNode = event.target.closest('.n-tree-node')
  console.log('🌊 treeNode found:', !!treeNode)
  
  if (treeNode) {
    console.log('🌊 containerDragOver - on tree node, let NTree handle')
    return // 让NTree处理
  }
  
  // 在空白区域
  console.log('🌊 containerDragOver - in empty area, checking root')
  lastMousePosition.value = { x: event.clientX, y: event.clientY }
  
  if (explorerContainer.value) {
    const containerRect = explorerContainer.value.getBoundingClientRect()
    const { clientX, clientY } = event
    
    console.log('🌊 mouse position:', clientX, clientY)
    console.log('🌊 container rect:', containerRect)
    
    if (clientX >= containerRect.left && clientX <= containerRect.right &&
        clientY >= containerRect.top && clientY <= containerRect.bottom) {
      
      console.log('✅ containerDragOver - ROOT AREA DETECTED!')
      isDragOverRoot.value = true
      event.preventDefault()
      event.dataTransfer.dropEffect = 'move'
    }
  }
}

function onContainerDrop(event) {
  console.log('🌊 containerDrop - capture event')
  console.log('🌊 isDragging.value:', isDragging.value)
  console.log('🌊 isDragOverRoot.value:', isDragOverRoot.value)
  console.log('🌊 draggedNode.value:', draggedNode.value)
  
  // 添加详细的调试信息
  console.log('🌊 containerDrop - dataTransfer check:', {
    hasDataTransfer: !!event.dataTransfer,
    typesLength: event.dataTransfer?.types?.length || 0,
    types: event.dataTransfer?.types || []
  })
  
  // 尝试恢复拖拽状态，如果没有的话
  if (!isDragging.value && event.dataTransfer) {
    console.log('🌊 containerDrop - trying to recover drag state')
    
    try {
      const dragData = event.dataTransfer.getData('application/json')
      console.log('🌊 containerDrop - dragData from dataTransfer:', dragData)
      
      if (dragData) {
        const nodeData = JSON.parse(dragData)
        draggedNode.value = nodeData
        isDragging.value = true
        console.log('🌊 containerDrop - recovered drag state:', nodeData)
      } else {
        // 尝试其他可能的数据类型
        console.log('🌊 containerDrop - trying other data types')
        for (const type of event.dataTransfer.types) {
          console.log('🌊 containerDrop - trying type:', type)
          const data = event.dataTransfer.getData(type)
          console.log('🌊 containerDrop - data for type', type, ':', data)
          
          // 如果找到了包含节点信息的数据
          if (data && data.includes('"key"') && data.includes('"fileName"')) {
            try {
              const nodeData = JSON.parse(data)
              draggedNode.value = nodeData
              isDragging.value = true
              console.log('🌊 containerDrop - recovered drag state from type', type, ':', nodeData)
              break
            } catch (e) {
              console.log('🌊 containerDrop - failed to parse data from type', type, ':', e)
            }
          }
        }
      }
    } catch (e) {
      console.log('🌊 containerDrop - could not recover drag state:', e)
    }
  }
  
  // 更宽松的退出条件 - 如果这是一个拖拽事件（有 DataTransfer），就继续处理
  if (!isDragging.value && !draggedNode.value && !event.dataTransfer) {
    console.log('🌊 containerDrop - no dragging context at all, exiting')
    return
  }
  
  // 如果有 DataTransfer 但没有我们的拖拽状态，可能是从 NTree 内部拖拽过来的
  if (!isDragging.value && !draggedNode.value && event.dataTransfer) {
    console.log('🌊 containerDrop - have dataTransfer but no internal state, treating as potential drag')
    // 在这种情况下，我们继续处理，尝试恢复状态
  }
  
  // 检查是否在树节点上
  const treeNode = event.target.closest('.n-tree-node')
  if (treeNode) {
    console.log('🌊 containerDrop - on tree node, let NTree handle')
    return // 让NTree处理
  }
  
  // 在空白区域 - 直接处理根目录拖拽
  console.log('🌊 containerDrop - handling root drop')
  event.preventDefault()
  event.stopPropagation()
  
  // 最后检查拖拽位置是否在容器内
  lastMousePosition.value = { x: event.clientX, y: event.clientY }
  
  // 优先使用已有的draggedNode，否则从dataTransfer恢复
  let nodeToMove = draggedNode.value
  
  if (!nodeToMove) {
    console.log('🌊 containerDrop - no draggedNode, trying to get from dataTransfer')
    try {
      const dragData = event.dataTransfer.getData('application/json')
      if (dragData) {
        nodeToMove = JSON.parse(dragData)
        console.log('🌊 containerDrop - recovered node data:', nodeToMove)
      }
    } catch (e) {
      console.log('🌊 containerDrop - error parsing drag data:', e)
    }
  }
  
  if (nodeToMove) {
    console.log('🌊 containerDrop - executing root drop for node:', nodeToMove.fileName || nodeToMove.label)
    handleRootDrop(nodeToMove)
  } else {
    console.log('🌊 containerDrop - no node data available')
    // 在没有节点数据的情况下，尝试使用当前选中的文件作为源
    console.log('🌊 containerDrop - current file info:', {
      currentFile: props.currentFile,
      hasCurrentFile: !!props.currentFile
    })
    
    if (props.currentFile) {
      // 使用当前选中的文件作为拖拽源
      console.log('🌊 containerDrop - using current file as drag source')
      
      // 从树数据中查找当前文件的详细信息
      const currentFileNode = findNodeByKey(localTreeData.value, props.currentFile)
      
      emit('move', {
        sourceId: props.currentFile,
        targetId: '0',
        sourceNode: currentFileNode,
        targetNode: { key: '0', isDirectory: true },
        isRootDrop: true,
        mousePosition: lastMousePosition.value
      })
    } else {
      console.log('🌊 containerDrop - no current file available, cannot determine drag source')
      message.warning('无法确定拖拽的源文件，请重新尝试')
    }
    
    // 清理拖拽状态
    clearDragState()
  }
}

// 测试函数 - 确认事件绑定工作
function onContainerClick(event) {
  console.log('🖱️ Container clicked!', event.target)
}

function onContainerDragEnter(event) {
  console.log('🌊 Container dragenter triggered!', event.target)
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
    const placeholder = isRenaming ? '' : ''
    
    return h('input', {
      class: 'vscode-tree-input',
      value: newName.value,
      autofocus: true,
      placeholder: placeholder,
      onInput: e => (newName.value = e.target.value),
      onKeydown: e => {
        if (e.key === 'Enter') {
          e.preventDefault()
          confirmAddNode()
        }
        if (e.key === 'Escape') {
          e.preventDefault()
          cancelAddNode()
        }
      },
      onBlur: () => {
        // 延迟执行，避免与点击事件冲突
        setTimeout(() => {
          if (editingNode.value || renamingNode.value) {
            confirmAddNode()
          }
        }, 100)
      }
    })
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

/* VSCode 风格的文件树输入框样式 */
:deep(.vscode-tree-input) {
  width: 100%;
  height: 22px;
  padding: 1px 4px;
  font-size: 13px;
  font-family: 'Segoe UI', 'Consolas', 'Monaco', monospace;
  background: #ffffff;
  border: 1px solid #007acc;
  border-radius: 0;
  outline: none;
  color: #333333;
  line-height: 18px;
  box-shadow: 0 0 0 1px #007acc;
  margin: 0;
  display: block;
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

/* 拖拽移动样式 */
:deep(.n-tree-node.dragging) {
  opacity: 0.5;
  background: rgba(0, 123, 204, 0.1);
}

:deep(.n-tree-node.drag-over) {
  background: rgba(0, 123, 204, 0.2);
  border: 2px dashed #007acc;
  border-radius: 4px;
}

:deep(.n-tree-node.drag-over .n-tree-node-content) {
  background: rgba(0, 123, 204, 0.1);
}

/* 拖拽时的全局样式 */
.template-explorer.dragging {
  user-select: none;
}

:deep(.n-tree-node[draggable="true"]) {
  cursor: grab;
}

:deep(.n-tree-node[draggable="true"]:active) {
  cursor: grabbing;
}

/* 根目录拖拽样式 */
.explorer-container.drag-over-root {
  background: rgba(0, 123, 204, 0.1);
  border: 2px dashed #007acc;
  border-radius: 4px;
  position: relative;
  animation: drag-over-pulse 1s infinite;
}

.explorer-container.drag-over-root::before {
  content: '拖拽到此处移动到根目录';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #007acc;
  font-size: 14px;
  font-weight: 500;
  background: rgba(255, 255, 255, 0.95);
  padding: 8px 12px;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 123, 204, 0.3);
  pointer-events: none;
  z-index: 1000;
  border: 1px solid #007acc;
}

@keyframes drag-over-pulse {
  0% { border-color: #007acc; }
  50% { border-color: #4db3ff; }
  100% { border-color: #007acc; }
}
</style> 