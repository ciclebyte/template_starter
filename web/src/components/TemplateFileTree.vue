<template>
  <div class="edit-tree" @contextmenu="onTreeAreaContextMenu">
    <div class="tree-title">文件树</div>
    <NTree
      :data="treeToNaive(localTreeData)"
      :default-expand-all="true"
      :selected-keys="[currentFile]"
      :render-label="renderLabel"
      :node-props="nodeProps"
      @update:selected-keys="onSelectFile"
    >
    </NTree>
    <div v-if="!treeData || treeData.length === 0" style="padding: 32px; color: #888; text-align: center; user-select: none; cursor: context-menu;" @contextmenu="onTreeAreaContextMenu">暂无数据（右键新建）</div>
    <!-- 原生右键菜单 -->
    <div
      v-if="showMenu"
      class="custom-menu"
      :style="{ left: menuX + 'px', top: menuY + 'px' }"
    >
      <div
        v-for="item in menuOptions"
        :key="item.key"
        class="menu-item"
        @click="handleMenuSelect(item.key)"
      >
        {{ item.label }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, h } from 'vue'
import { NTree, useMessage, NIcon } from 'naive-ui'
import { FileTrayFullOutline, Folder, FolderOpenOutline } from '@vicons/ionicons5'

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
const emit = defineEmits(['select', 'reload', 'update:treeData'])

const showMenu = ref(false)
const menuX = ref(0)
const menuY = ref(0)
const menuNode = ref(null)
const addType = ref('file')
const editingNode = ref(null)
const newName = ref('')
const message = useMessage()
const localTreeData = ref(JSON.parse(JSON.stringify(props.treeData)))

// 右键菜单选项
const menuOptions = [
  { label: '新增文件', key: 'addFile' },
  { label: '新增文件夹', key: 'addFolder' },
  { label: '删除节点', key: 'deleteNode' }
]

watch(() => props.treeData, (val) => {
  localTreeData.value = JSON.parse(JSON.stringify(val))
}, { deep: true })

function treeToNaive(tree) {
  if (!Array.isArray(tree)) return []
  const customSort = (a, b) => {
    // 1. 文件夹在前
    if ((b.isDirectory || 0) - (a.isDirectory || 0) !== 0) {
      return (b.isDirectory || 0) - (a.isDirectory || 0)
    }
    // 2. a开头在前
    const nameA = a.filePath ? a.filePath.split('/').pop().toLowerCase() : (a.label || '').toLowerCase()
    const nameB = b.filePath ? b.filePath.split('/').pop().toLowerCase() : (b.label || '').toLowerCase()
    const aIsA = nameA.startsWith('a') ? 1 : 0
    const bIsA = nameB.startsWith('a') ? 1 : 0
    if (aIsA !== bIsA) return bIsA - aIsA
    // 3. 名称升序
    return nameA.localeCompare(nameB)
  }
  const sorted = [...tree].sort(customSort)
  return sorted.map(node => ({
    label: node.isEditing === true ? undefined : (node.filePath ? node.filePath.split('/').pop() : node.label),
    key: node.key || node.id,
    isLeaf: node.isDirectory === 0,
    isEditing: node.isEditing === true,
    filePath: node.filePath,
    prefix: node.isDirectory
      ? () => h(NIcon, null, { default: () => h(Folder) })
      : () => h(NIcon, null, { default: () => h(FileTrayFullOutline) }),
    children: node.children ? treeToNaive(node.children) : []
  }))
}

function onSelectFile(keys) {
  if (keys && keys.length > 0) {
    emit('select', keys[0])
  }
}
function onTreeNodeRightClick({ event, node }) {
  event.preventDefault()
  event.stopPropagation()
  showMenu.value = true
  menuX.value = event.clientX
  menuY.value = event.clientY
  menuNode.value = node
}
function onTreeAreaContextMenu(event) {
  if (event.target.closest('.n-tree')) return
  event.preventDefault()
  event.stopPropagation()
  showMenu.value = true
  menuX.value = event.clientX
  menuY.value = event.clientY
  menuNode.value = null
}
function addFile() {
  showMenu.value = false
  addType.value = 'file'
  insertEditingNode('file')
}
function addFolder() {
  showMenu.value = false
  addType.value = 'folder'
  insertEditingNode('folder')
}
function insertEditingNode(type) {
  const newTreeData = JSON.parse(JSON.stringify(localTreeData.value))
  removeEditingNode(newTreeData)
  newName.value = ''
  const newKey = '__new__' + Date.now() + Math.random().toString(36).slice(2)
  // 计算同级最大sort
  let siblings = []
  const parentId = menuNode.value ? String(menuNode.value.id || menuNode.value.key) : '0'
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
  console.log('insertEditingNode siblings:', siblings)
  console.log('insertEditingNode maxSort:', maxSort)
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
  console.log('localTreeData after insert:', JSON.stringify(localTreeData.value, null, 2))
  console.log('treeToNaive after insert:', JSON.stringify(treeToNaive(localTreeData.value), null, 2))
}
function insertToParent(list, parentId, node) {
  for (const item of list) {
    console.log('insertToParent check', item.id, parentId, typeof item.id, typeof parentId)
    if (String(item.id) === String(parentId)) {
      if (!item.children) item.children = []
      item.children.unshift(node)
      console.log('inserted to parent', parentId)
      return true
    }
    if (item.children && insertToParent(item.children, parentId, node)) return true
  }
  console.log('parent not found for', parentId)
  return false
}
function removeEditingNode(list) {
  for (let i = list.length - 1; i >= 0; i--) {
    if (list[i].isEditing) {
      list.splice(i, 1)
    } else if (list[i].children) {
      removeEditingNode(list[i].children)
    }
  }
}
async function confirmAddNode() {
  if (!newName.value) {
    message.warning('请输入名称')
    return
  }
  emit('reload', { name: newName.value, type: addType.value, parentId: editingNode.value.parentId, sort: editingNode.value.sort })
  editingNode.value = null
  removeEditingNode(localTreeData.value)
}
function cancelAddNode() {
  editingNode.value = null
  removeEditingNode(localTreeData.value)
}
function nodeProps({ option }) {
  return {
    onContextmenu(e) {
      e.preventDefault();
      showMenu.value = true;
      menuX.value = e.clientX;
      menuY.value = e.clientY;
      menuNode.value = option;
    }
  }
}
function renderLabel({ option }) {
  if (option.isEditing === true) {
    return h(
      'div',
      { style: 'display:flex;align-items:center;gap:4px;' },
      [
        h('input', {
          style: 'width:120px;background: #ffe0e0; border: 1px solid #f00;',
          value: newName.value,
          autofocus: true,
          placeholder: '请输入名称',
          onInput: e => (newName.value = e.target.value),
          onKeydown: e => {
            if (e.key === 'Enter') confirmAddNode()
            if (e.key === 'Escape') cancelAddNode()
          }
        }),
        h(
          'button',
          { style: 'padding:0 4px;', onClick: () => confirmAddNode() },
          '✔'
        ),
        h(
          'button',
          { style: 'padding:0 4px;', onClick: () => cancelAddNode() },
          '✖'
        )
      ]
    )
  }
  return option.filePath ? option.filePath.split('/').pop() : option.label
}

function handleMenuSelect(key) {
  if (key === 'addFile') addFile()
  else if (key === 'addFolder') addFolder()
  else if (key === 'deleteNode') deleteNode()
}

function deleteNode() {
  if (!menuNode.value) return
  const id = String(menuNode.value.id || menuNode.value.key)
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
  showMenu.value = false
  emit('reload', { type: 'delete', id })
}

function updatePrefixWithExpaned(_keys, _option, meta) {
  if (!meta.node) return;
  switch (meta.action) {
    case 'expand':
      meta.node.prefix = () => h(NIcon, null, { default: () => h(FolderOpenOutline) });
      break;
    case 'collapse':
      meta.node.prefix = () => h(NIcon, null, { default: () => h(Folder) });
      break;
  }
}
</script>

<style scoped>
.edit-tree {
  width: 260px;
  background: #fff;
  border-right: 1px solid #e0e0e0;
  padding: 24px 12px 0 12px;
  display: flex;
  flex-direction: column;
}
.tree-title {
  font-weight: bold;
  margin-bottom: 16px;
  color: #333;
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
</style> 