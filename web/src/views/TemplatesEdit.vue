<template>
  <div class="templates-edit-fullscreen">
    <div class="edit-header">
      <span class="edit-title">模板编辑</span>
      <n-button quaternary circle class="edit-close-btn" @click="closeEdit">
        <template #icon>
          <n-icon><svg viewBox="0 0 24 24" width="20" height="20"><path fill="currentColor" d="M18.3 5.71a1 1 0 0 0-1.41 0L12 10.59 7.11 5.7A1 1 0 0 0 5.7 7.11L10.59 12l-4.89 4.89a1 1 0 1 0 1.41 1.41L12 13.41l4.89 4.89a1 1 0 0 0 1.41-1.41L13.41 12l4.89-4.89a1 1 0 0 0 0-1.4z"/></svg></n-icon>
        </template>
      </n-button>
    </div>
    <div class="edit-main">
      <TemplateFileTree
        :treeData="treeData.value || []"
        :currentFile="currentFile.value || ''"
        @select="onSelectFile"
        @reload="onTreeReload"
      />
      <TemplateEditor
        :openedTabs="openedTabs.value || []"
        :activeTab="activeTab.value || ''"
        :fileMap="fileMap.value || {}"
        @tabChange="onTabChange"
        @tabClose="onTabClose"
        @contentChange="onEditorContentChange"
      />
    </div>
  </div>
</template>

<script setup>
import { useRouter, useRoute } from 'vue-router'
import { ref, onMounted } from 'vue'
import { getTemplateFileTree, addTemplateFile } from '@/api/templateFiles'
import TemplateFileTree from '@/components/TemplateFileTree.vue'
import TemplateEditor from '@/components/TemplateEditor.vue'

const router = useRouter()
const route = useRoute()
const closeEdit = () => {
  router.push('/templates')
}

const treeData = ref([])
const loadingTree = ref(true)
const noTreeData = ref(false)
const currentFile = ref('')
const openedTabs = ref([])
const activeTab = ref('')
const fileMap = ref({})

onMounted(async () => {
  await loadTree()
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

function onSelectFile(key) {
  if (!key) return
  // 打开新标签或切换
  const exist = openedTabs.value.find(tab => tab.key === key)
  if (!exist) {
    openedTabs.value.push({ key, name: key.split('/').pop(), content: fileMap.value[key] || '' })
  }
  activeTab.value = key
  currentFile.value = key
}

function onTabChange(key) {
  activeTab.value = key
  currentFile.value = key
}

function onTabClose(key) {
  const idx = openedTabs.value.findIndex(tab => tab.key === key)
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
  // 这里处理新增文件/文件夹逻辑，后续可接接口
  // 你可以在这里调用 addTemplateFile 或其它接口，然后刷新 treeData
  // 目前仅console.log
  console.log('onTreeReload', payload)
}
</script>

<style scoped>
.templates-edit-fullscreen {
  position: fixed;
  inset: 0;
  background: #f5f5f5;
  z-index: 9999;
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
.edit-title {
  font-size: 1.2rem;
  font-weight: bold;
  color: #18a058;
}
.edit-close-btn {
  margin-left: 16px;
}
.edit-main {
  flex: 1;
  display: flex;
  min-height: 0;
}
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
.edit-editor {
  flex: 1;
  background: #f8fafc;
  padding: 24px 24px 0 24px;
  display: flex;
  flex-direction: column;
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
</style> 