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
      <div class="edit-tree">
        <div class="tree-title">文件树</div>
        <pre>treeData: {{ JSON.stringify(treeData.value, null, 2) }}</pre>
        <pre>naiveTreeData: {{ JSON.stringify(naiveTreeData, null, 2) }}</pre>
        <n-tree
          :data="naiveTreeData"
          :default-expand-all="true"
          :selected-keys="[currentFile]"
          @update:selected-keys="onSelectFile"
          @node-contextmenu="onTreeNodeRightClick"
        />
        <div v-if="!naiveTreeData || naiveTreeData.length === 0" style="padding: 32px; color: #888; text-align: center; user-select: none; cursor: context-menu;">暂无数据（右键新建）</div>
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
      <div class="edit-editor">
        <div class="editor-tabs">
          <n-tabs type="card" :value="activeTab" @update:value="onTabChange" closable @close="onTabClose">
            <n-tab v-for="tab in openedTabs" :key="tab.key" :name="tab.key">
              {{ tab.name }}
            </n-tab>
          </n-tabs>
        </div>
        <div class="editor-title">{{ activeTab }}</div>
        <div class="monaco-editor-container" ref="monacoContainer"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter, useRoute } from 'vue-router'
import { NIcon, NDropdown, NModal, NInput, useMessage } from 'naive-ui'
import { ref, onMounted, onBeforeUnmount, watch, h, computed } from 'vue'
import { getTemplateFileTree, addTemplateFile } from '@/api/templateFiles'

const router = useRouter()
const route = useRoute()
const closeEdit = () => {
  router.push('/templates')
}

const treeData = ref([])
const loadingTree = ref(true)
const noTreeData = ref(false)
const showMenu = ref(false)
const menuX = ref(0)
const menuY = ref(0)
const menuNode = ref(null) // 右键的节点，null为根
const message = useMessage()
const body = typeof window !== 'undefined' ? document.body : undefined

const addType = ref('file')

onMounted(async () => {
  await loadTree()
})

function treeToNaive(tree) {
  if (!Array.isArray(tree)) return []
  return tree.map(node => {
    if (node.isEditing) {
      return {
        key: node.key || node.id,
        label: () =>
          h(
            'div',
            { style: 'display:flex;align-items:center;gap:4px;' },
            [
              h('input', {
                style: 'width:120px',
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
                {
                  style: 'padding:0 4px;',
                  onClick: () => confirmAddNode()
                },
                '✔'
              ),
              h(
                'button',
                {
                  style: 'padding:0 4px;',
                  onClick: () => cancelAddNode()
                },
                '✖'
              )
            ]
          ),
        isLeaf: node.isLeaf,
        children: []
      }
    }
    return {
      label: node.filePath ? node.filePath.split('/').pop() : node.label,
      key: node.key || node.id,
      isLeaf: node.isDirectory === 0,
      children: node.children ? treeToNaive(node.children) : []
    }
  })
}

// 文件内容映射
const fileMap = ref({
  'src/App.vue': `<template>\n  <h1>你好，Monaco Editor！</h1>\n</template>\n\n<script setup>\nconst msg = '你好，Monaco！'\n<\/script>\n`,
  'src/main.js': `import { createApp } from 'vue'\nimport App from './App.vue'\ncreateApp(App).mount('#app')`,
  'src/components/HelloWorld.vue': `<template>\n  <div>你好，世界</div>\n</template>`,
  'src/assets/logo.png': '// 二进制文件无法预览',
  'public/index.html': `<!DOCTYPE html>\n<html lang='zh-CN'>\n<head>\n  <meta charset='UTF-8' />\n  <title>Vite + Vue</title>\n</head>\n<body>\n  <div id='app'></div>\n</body>\n</html>`,
  'package.json': `{"name": "template-starter", "version": "1.0.0"}`,
  'vite.config.js': `import { defineConfig } from 'vite'\nexport default defineConfig({})`,
  'README.md': `# 模板项目\n\n这是一个项目模板。`
})
const currentFile = ref('src/App.vue')

// Monaco Editor 集成
import * as monaco from 'monaco-editor'
const monacoContainer = ref(null)
let editorInstance = null

const getLanguage = (file) => {
  if (file.endsWith('.vue')) return 'vue'
  if (file.endsWith('.js')) return 'javascript'
  if (file.endsWith('.json')) return 'json'
  if (file.endsWith('.md')) return 'markdown'
  if (file.endsWith('.html')) return 'html'
  return 'plaintext'
}

const openedTabs = ref([
  { key: 'src/App.vue', name: 'App.vue', content: fileMap.value['src/App.vue'] }
])
const activeTab = ref('src/App.vue')

// 文件树点击，打开新标签或切换
const onSelectFile = (keys) => {
  if (keys && keys.length > 0) {
    const key = keys[0]
    const exist = openedTabs.value.find(tab => tab.key === key)
    if (!exist) {
      openedTabs.value.push({ key, name: key.split('/').pop(), content: fileMap.value[key] || '' })
    }
    activeTab.value = key
    currentFile.value = key
  }
}

// 标签切换
const onTabChange = (key) => {
  activeTab.value = key
  currentFile.value = key
}

// 标签关闭
const onTabClose = (key) => {
  const idx = openedTabs.value.findIndex(tab => tab.key === key)
  if (idx !== -1) {
    openedTabs.value.splice(idx, 1)
    // 切换到相邻标签
    if (openedTabs.value.length > 0) {
      const next = openedTabs.value[Math.max(0, idx - 1)]
      activeTab.value = next.key
      currentFile.value = next.key
    }
  }
}

// Monaco Editor 内容与标签同步
watch([activeTab, openedTabs], () => {
  const tab = openedTabs.value.find(t => t.key === activeTab.value)
  if (tab && editorInstance) {
    const lang = getLanguage(tab.key)
    let model = monaco.editor.getModels().find(m => m.uri.path.endsWith(tab.key))
    if (!model) {
      model = monaco.editor.createModel(tab.content, lang, monaco.Uri.parse(`file:///${tab.key}`))
    }
    editorInstance.setModel(model)
  }
})

onMounted(() => {
  editorInstance = monaco.editor.create(monacoContainer.value, {
    value: openedTabs.value[0].content,
    language: getLanguage(openedTabs.value[0].key),
    theme: 'vs-dark',
    automaticLayout: true,
    fontSize: 15,
    minimap: { enabled: false }
  })
  editorInstance.onDidChangeModelContent(() => {
    const tab = openedTabs.value.find(t => t.key === activeTab.value)
    if (tab) {
      tab.content = editorInstance.getValue()
      fileMap.value[tab.key] = tab.content
    }
  })
})

onBeforeUnmount(() => {
  if (editorInstance) {
    editorInstance.dispose()
  }
})

function onTreeNodeRightClick({ event, node }) {
  event.preventDefault()
  event.stopPropagation()
  showMenu.value = true
  menuX.value = event.clientX
  menuY.value = event.clientY
  menuNode.value = node
}

const editingNode = ref(null) // 当前正在编辑的节点key
const newName = ref('')

function addFile() {
  showMenu.value = false
  addType.value = 'file'
  // 插入编辑节点
  insertEditingNode('file')
}
function addFolder() {
  showMenu.value = false
  addType.value = 'folder'
  insertEditingNode('folder')
}
function insertEditingNode(type) {
  // 只允许一个编辑节点
  removeEditingNode(treeData.value)
  newName.value = ''
  const newKey = '__new__' + Date.now()
  editingNode.value = {
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
  if (editingNode.value.parentId === 0) {
    treeData.value.unshift(editingNode.value)
  } else {
    insertToParent(treeData.value, editingNode.value.parentId, editingNode.value)
  }
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
async function confirmAddNode() {
  if (!newName.value) {
    message.warning('请输入名称')
    return
  }
  const templateId = route.params.id
  const isDirectory = addType.value === 'folder' ? 1 : 0
  try {
    await addTemplateFile({
      templateId,
      filePath: newName.value,
      fileContent: '',
      fileSize: 0,
      isDirectory,
      md5: '',
      sort: 0,
      parentId: editingNode.value.parentId
    })
    message.success('新增成功')
    editingNode.value = null
    await loadTree()
  } catch (e) {
    message.error('新增失败')
  }
}
function cancelAddNode() {
  editingNode.value = null
  removeEditingNode(treeData.value)
}

function onClickOutside() {
  showMenu.value = false
}

onMounted(() => {
  document.addEventListener('click', onClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onClickOutside)
})

async function loadTree() {
  loadingTree.value = true
  try {
    const res = await getTemplateFileTree(route.params.id)
    console.log('接口完整返回', res)
    console.log('res.data', res.data)
    const tree = res.data && res.data.data && res.data.data.tree
    console.log('tree', tree)
    if (tree && tree.length > 0) {
      treeData.value = tree
      noTreeData.value = false
      console.log('接口tree', tree)
      console.log('treeData.value', treeData.value)
    } else {
      treeData.value = []
      noTreeData.value = true
      console.log('接口tree为空', tree)
    }
  } catch (e) {
    treeData.value = []
    noTreeData.value = true
    console.log('接口异常', e)
  }
  loadingTree.value = false
}

function onEmptyTreeRightClick(event) {
  showMenu.value = true
  menuX.value = event.clientX
  menuY.value = event.clientY
  menuNode.value = null
}

const naiveTreeData = computed(() =>
  Array.isArray(treeData.value)
    ? treeData.value.map(item => ({
        key: item.id,
        label: item.filePath,
        isLeaf: item.isDirectory === 0
      }))
    : []
)
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