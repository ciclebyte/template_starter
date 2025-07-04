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
        v-model:treeData="treeData"
        :currentFile="currentFile"
        @select="onSelectFile"
        @reload="onTreeReload"
      />
      <TemplateEditor
        :filePath="currentFilePath"
        :fileContent="currentFileContent"
        :openedTabs="openedTabs"
        :activeTab="activeTab"
        :fileMap="fileMap"
        @tabChange="onTabChange"
        @tabClose="onTabClose"
        @contentChange="onEditorContentChange"
      />
    </div>
  </div>
</template>

<script setup>
import { useRouter, useRoute } from 'vue-router'
import { ref, onMounted, watch } from 'vue'
import { getTemplateFileTree, addTemplateFile, delTemplateFile, getTemplateFileDetail,getTemplateFileContent } from '@/api/templateFiles'
import TemplateFileTree from '@/components/TemplateFileTree.vue'
import TemplateEditor from '@/components/TemplateEditor.vue'
import { useTemplateFileStore } from '@/stores/templateFileStore'

const router = useRouter()
const route = useRoute()
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
          name: node.filePath || node.label || String(key),
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
    filePath: payload.name,
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

// 调试：打印treeData变化
watch(treeData, (val) => {
  console.log('父组件 treeData 变化:', val)
}, { deep: true })
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