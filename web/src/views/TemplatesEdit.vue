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
        <n-tree :data="treeData" :default-expand-all="true" block-line :selected-keys="[currentFile]" @update:selected-keys="onSelectFile" />
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
import { useRouter } from 'vue-router'
import { NIcon } from 'naive-ui'
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
const router = useRouter()
const closeEdit = () => {
  router.push('/templates')
}

const treeData = ref([
  {
    label: 'src',
    key: 'src',
    children: [
      { label: 'App.vue', key: 'src/App.vue' },
      { label: 'main.js', key: 'src/main.js' },
      {
        label: 'components',
        key: 'src/components',
        children: [
          { label: 'HelloWorld.vue', key: 'src/components/HelloWorld.vue' }
        ]
      },
      {
        label: 'assets',
        key: 'src/assets',
        children: [
          { label: 'logo.png', key: 'src/assets/logo.png' }
        ]
      }
    ]
  },
  { label: 'public', key: 'public', children: [ { label: 'index.html', key: 'public/index.html' } ] },
  { label: 'package.json', key: 'package.json' },
  { label: 'vite.config.js', key: 'vite.config.js' },
  { label: 'README.md', key: 'README.md' }
])

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
</style> 