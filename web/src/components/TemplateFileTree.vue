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
  showDropdown.value = true
  dropdownX.value = event.clientX
  dropdownY.value = event.clientY
}
function handleDropdownSelect(key) {
  showDropdown.value = false
  if (key === 'addFile') message.info('新增文件')
  else if (key === 'addFolder') message.info('新增文件夹')
  else if (key === 'deleteNode') message.info('删除节点')
}
function handleDropdownClickoutside() {
  showDropdown.value = false
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