<template>
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
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { NTabs, NTab } from 'naive-ui'
import * as monaco from 'monaco-editor'

const props = defineProps({
  openedTabs: {
    type: Array,
    default: () => []
  },
  activeTab: {
    type: String,
    default: ''
  },
  fileMap: {
    type: Object,
    default: () => ({})
  }
})
const emit = defineEmits(['tabChange', 'tabClose', 'contentChange'])

const monacoContainer = ref(null)
let editorInstance = null

function getLanguage(file) {
  if (file.endsWith('.vue')) return 'vue'
  if (file.endsWith('.js')) return 'javascript'
  if (file.endsWith('.json')) return 'json'
  if (file.endsWith('.md')) return 'markdown'
  if (file.endsWith('.html')) return 'html'
  return 'plaintext'
}

function onTabChange(key) {
  emit('tabChange', key)
}
function onTabClose(key) {
  emit('tabClose', key)
}

watch(() => [props.activeTab, props.openedTabs], () => {
  const tab = props.openedTabs.find(t => t.key === props.activeTab)
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
  const tab = props.openedTabs.find(t => t.key === props.activeTab)
  editorInstance = monaco.editor.create(monacoContainer.value, {
    value: tab ? tab.content : '',
    language: tab ? getLanguage(tab.key) : 'plaintext',
    theme: 'vs-dark',
    automaticLayout: true,
    fontSize: 15,
    minimap: { enabled: false }
  })
  editorInstance.onDidChangeModelContent(() => {
    const tab = props.openedTabs.find(t => t.key === props.activeTab)
    if (tab) {
      tab.content = editorInstance.getValue()
      emit('contentChange', { key: tab.key, content: tab.content })
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