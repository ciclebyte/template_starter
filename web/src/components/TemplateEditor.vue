<template>
  <div class="edit-editor">
    <div class="editor-tabs">
      <n-tabs type="card" :value="activeTab" @update:value="onTabChange" closable @close="onTabClose">
        <n-tab v-for="tab in openedTabs" :key="tab.key" :name="tab.key">
          {{ tab.name }}
        </n-tab>
      </n-tabs>
    </div>
    <div class="monaco-editor-container" ref="monacoContainer"></div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { NTabs, NTab, useMessage, useNotification } from 'naive-ui'
import * as monaco from 'monaco-editor'
import { useTemplateFileStore } from '@/stores/templateFileStore'
import { storeToRefs } from 'pinia'
import { editTemplateFile } from '@/api/templateFiles'
import { useRoute } from 'vue-router'
// shiki + shikiToMonaco
import { createHighlighter } from 'shiki'
import { shikiToMonaco } from '@shikijs/monaco'

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
let highlighter = null

const templateFileStore = useTemplateFileStore()
const { currentFileContent } = storeToRefs(templateFileStore)

const route = useRoute()
const notification = useNotification()

function getLanguage(file) {
  return 'vue'
}

function onTabChange(key) {
  emit('tabChange', key)
}
function onTabClose(key) {
  emit('tabClose', key)
}

watch(() => [props.activeTab, props.openedTabs, currentFileContent.value], () => {
  const tab = props.openedTabs.find(t => t.key === props.activeTab)
  if (tab && editorInstance) {
    const lang = tab ? getLanguage(tab.name) : 'plaintext'
    let model = monaco.editor.getModels().find(m => m.uri.path.endsWith(tab.key))
    if (!model) {
      model = monaco.editor.createModel(tab.content, lang, monaco.Uri.parse(`file:///${tab.key}`))
    }
    if (model.getValue() !== currentFileContent.value) {
      model.setValue(currentFileContent.value || '')
    }
    editorInstance.setModel(model)
  }
})

watch(currentFileContent, (val) => {
  const tab = props.openedTabs.find(t => t.key === props.activeTab)
  if (tab && tab.content !== val) {
    tab.content = val
    if (editorInstance) {
      editorInstance.setValue(val || '')
    }
  }
})

async function saveCurrentFile() {
  const tab = props.openedTabs.find(t => t.key === props.activeTab)
  if (!tab) return
  const templateId = route.params.id
  try {
    await editTemplateFile({
      id: tab.key,
      templateId,
      filePath: tab.name,
      fileContent: tab.content,
      isDirectory: 0
    })
    notification.success({
      title: '保存成功',
      content: '文件已成功保存',
      duration: 2500
    })
  } catch (e) {
    notification.error({
      title: '保存失败',
      content: '请检查网络或稍后重试',
      duration: 2500
    })
  }
}

function handleKeydown(e) {
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    saveCurrentFile()
  }
}

onMounted(async () => {
  highlighter = await createHighlighter({
    themes: ['github-dark', 'nord', 'github-light'],
    langs: [
      'vue', 'javascript', 'typescript', 'json', 'markdown', 'html',
      'go', 'java', 'python', 'c', 'cpp', 'csharp', 'php', 'ruby', 'shell', 'xml', 'yaml'
    ]
  })
  console.log('shiki highlighter:', highlighter)
  await shikiToMonaco(highlighter, monaco)
  console.log('shikiToMonaco 注入完成')
  monaco.editor.setTheme('github-dark')
  console.log('monaco theme set:', 'github-dark')

  const tab = props.openedTabs.find(t => t.key === props.activeTab)
  const lang = tab ? getLanguage(tab.name) : 'plaintext'
  console.log('monaco create language:', lang)
  editorInstance = monaco.editor.create(monacoContainer.value, {
    value: tab ? tab.content : '',
    language: lang,
    theme: 'github-dark',
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
  window.addEventListener('keydown', handleKeydown)

  // 注册右键菜单保存命令
  if (editorInstance) {
    editorInstance.addAction({
      id: 'save-file',
      label: '保存 (Ctrl+S)',
      keybindings: [
        monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS
      ],
      precondition: null,
      keybindingContext: null,
      contextMenuGroupId: 'navigation',
      contextMenuOrder: 1.5,
      run: function(ed) {
        saveCurrentFile()
      }
    })
  }
})

onBeforeUnmount(() => {
  if (editorInstance) {
    editorInstance.dispose()
  }
  window.removeEventListener('keydown', handleKeydown)
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