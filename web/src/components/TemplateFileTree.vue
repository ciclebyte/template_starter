<template>
  <div class="edit-tree" @contextmenu="onTreeAreaContextMenu">
    <div class="tree-title">文件树</div>
    <NTree
      :data="treeToNaive(treeData)"
      :default-expand-all="true"
      :selected-keys="[currentFile]"
      @update:selected-keys="onSelectFile"
      @node-contextmenu="onTreeNodeRightClick"
      @contextmenu.stop
      :render-label="renderLabel"
    >
    </NTree>
    <div v-if="!treeData || treeData.length === 0" style="padding: 32px; color: #888; text-align: center; user-select: none; cursor: context-menu;" @contextmenu="onTreeAreaContextMenu">暂无数据（右键新建）</div>
    <!-- 原生右键菜单 -->
    <div
      v-if="showMenu"
      class="custom-menu"
      :style="{ left: menuX + 'px', top: menuY + 'px' }"
    >
      <div class="menu-item" @click="addFile">新增文件</div>
      <div class="menu-item" @click="addFolder">新增文件夹</div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, h } from 'vue'
import { NTree, useMessage } from 'naive-ui'

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

function treeToNaive(tree) {
  if (!Array.isArray(tree)) return []
  const result = tree.map(node => {
    return {
      label: node.isEditing === true ? undefined : (node.filePath ? node.filePath.split('/').pop() : node.label),
      key: node.key || node.id,
      isLeaf: node.isDirectory === 0,
      isEditing: node.isEditing === true,
      filePath: node.filePath,
      children: node.children ? treeToNaive(node.children) : []
    }
  })
  console.log('treeToNaive result:', JSON.stringify(result, null, 2))
  return result
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
  console.log("@@@@insertEditingNode@@@",type)
  const newTreeData = JSON.parse(JSON.stringify(props.treeData))
  removeEditingNode(newTreeData)
  newName.value = ''
  const newKey = '__new__' + Date.now()
  const newNode = {
    key: newKey,
    id: newKey,
    label: '',
    filePath: '',
    isEditing: true,
    isLeaf: type === 'file',
    isDirectory: type === 'folder' ? 1 : 0,
    parentId: menuNode.value ? (menuNode.value.key || menuNode.value.id) : 0,
    children: []
  }
  editingNode.value = newNode
  console.log("newNode@@@",newNode)
  if (newNode.parentId === 0) {
    newTreeData.unshift(newNode)
  } else {
    insertToParent(newTreeData, newNode.parentId, newNode)
  }
  console.log("newTreeData@@@",newTreeData)
  emit('update:treeData', newTreeData)
}
function insertToParent(list, parentId, node) {
  for (const item of list) {
    if (item.key === parentId) {
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
  emit('reload', { name: newName.value, type: addType.value, parentId: editingNode.value.parentId })
  editingNode.value = null
  removeEditingNode(props.treeData)
}
function cancelAddNode() {
  editingNode.value = null
  removeEditingNode(props.treeData)
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