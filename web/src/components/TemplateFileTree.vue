<template>
  <div class="edit-tree" @contextmenu="onTreeAreaContextMenu">
    <div class="tree-title">文件树</div>
    <NTree
      :data="treeToNaive(localTreeData)"
      :default-expand-all="true"
      :selected-keys="[currentFile]"
      :render-label="renderLabel"
      :node-props="nodeProps"
      :render-switcher-icon="renderSwitcherIcon"
      @update:selected-keys="onSelectFile"
    />
    <div v-if="!treeData || treeData.length === 0" style="padding: 32px; color: #888; text-align: center; user-select: none; cursor: context-menu;" @contextmenu="onTreeAreaContextMenu">暂无数据（右键新建）</div>
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
  </div>
</template>

<script setup>
import { ref, watch, h } from 'vue'
import { NTree, useMessage, NIcon } from 'naive-ui'
import { ChevronForward, FileTrayFullOutline, Folder, FolderOpenOutline, Trash } from '@vicons/ionicons5'

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
const message = useMessage()
const localTreeData = ref(JSON.parse(JSON.stringify(props.treeData)))

watch(() => props.treeData, (val) => {
  localTreeData.value = JSON.parse(JSON.stringify(val))
}, { deep: true })

function treeToNaive(tree) {
  if (!Array.isArray(tree)) return []
  const customSort = (a, b) => {
    if ((b.isDirectory || 0) - (a.isDirectory || 0) !== 0) {
      return (b.isDirectory || 0) - (a.isDirectory || 0)
    }
    const nameA = a.filePath ? a.filePath.split('/').pop().toLowerCase() : (a.label || '').toLowerCase()
    const nameB = b.filePath ? b.filePath.split('/').pop().toLowerCase() : (b.label || '').toLowerCase()
    const aIsA = nameA.startsWith('a') ? 1 : 0
    const bIsA = nameB.startsWith('a') ? 1 : 0
    if (aIsA !== bIsA) return bIsA - aIsA
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
function nodeProps({ option }) {
  return {
    onContextmenu(e) {
      e.preventDefault();
      dropdownNode.value = option
      if (option.isLeaf) {
        dropdownOptions.value = [
          { label: '删除节点', key: 'deleteNode', icon: () => h(NIcon, null, { default: () => h(Trash) }) }
        ]
      } else {
        dropdownOptions.value = [
          { label: '新增文件', key: 'addFile', icon: () => h(NIcon, null, { default: () => h(FileTrayFullOutline) }) },
          { label: '新增文件夹', key: 'addFolder', icon: () => h(NIcon, null, { default: () => h(Folder) }) },
          { type: 'divider', key: 'divider1' },
          { label: '删除节点', key: 'deleteNode', icon: () => h(NIcon, null, { default: () => h(Trash) }) },
          { type: 'divider', key: 'divider2' },
          { label: '上传代码文件', key: 'uploadFile', icon: () => h(NIcon, null, { default: () => h(FileTrayFullOutline) }) }
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
    { label: '上传代码文件', key: 'uploadFile', icon: () => h(NIcon, null, { default: () => h(FileTrayFullOutline) }) },
    { label: '上传zip包', key: 'uploadZip', icon: () => h(NIcon, null, { default: () => h(Folder) }) }
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
  emit('reload', { name: newName.value, type: addType.value, parentId: editingNode.value.parentId, sort: editingNode.value.sort })
  editingNode.value = null
  removeEditingNode(localTreeData.value)
}
function cancelAddNode() {
  editingNode.value = null
  removeEditingNode(localTreeData.value)
}
function handleDropdownSelect(key) {
  showDropdown.value = false
  if (key === 'addFile') addFile()
  else if (key === 'addFolder') addFolder()
  else if (key === 'deleteNode') deleteNode()
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
const renderSwitcherIcon = () =>
  h(NIcon, null, { default: () => h(ChevronForward) })
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
.tree-dropdown-menu {
  z-index: 2147483647 !important;
}
</style> 