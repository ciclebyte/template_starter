<template>
  <div class="edit-editor">
    <div class="editor-tabs">
      <n-tabs type="card" :value="activeTab" @update:value="onTabChange" closable @close="onTabClose">
        <n-tab v-for="tab in openedTabs" :key="tab.key" :name="tab.key">
          {{ tab.name }}
        </n-tab>
      </n-tabs>
    </div>
    <div class="codemirror-container" ref="editorContainer"></div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { NTabs, NTab, useNotification } from 'naive-ui'
import { useTemplateFileStore } from '@/stores/templateFileStore'
import { storeToRefs } from 'pinia'
import { editTemplateFile } from '@/api/templateFiles'
import { useRoute } from 'vue-router'

// CodeMirror 核心模块
import { EditorView } from '@codemirror/view'
import { EditorState } from '@codemirror/state'
import { keymap } from '@codemirror/view'
import { defaultKeymap, indentWithTab } from '@codemirror/commands'

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

const editorContainer = ref(null)
let editorView = null

const templateFileStore = useTemplateFileStore()
const { currentFileContent } = storeToRefs(templateFileStore)

const route = useRoute()
const notification = useNotification()

function onTabChange(key) {
  emit('tabChange', key)
}

function onTabClose(key) {
  emit('tabClose', key)
}

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

function updateEditorContent(content) {
  if (!editorView) return

  const newState = EditorState.create({
    doc: content || '',
    extensions: createEditorExtensions()
  })

  editorView.setState(newState)
}

function createEditorExtensions() {
  return [
    // 基础快捷键
    keymap.of([
      ...defaultKeymap,
      indentWithTab,
      {
        key: 'Ctrl-s',
        run: () => {
          saveCurrentFile()
          return true
        }
      }
    ]),
    
    // 右键菜单
    EditorView.domEventHandlers({
      contextmenu: (event, view) => {
        event.preventDefault()
        showContextMenu(event, view)
        return true
      }
    }),
    
    // 内容变化监听
    EditorView.updateListener.of((update) => {
      if (update.docChanged) {
        const content = update.state.doc.toString()
        const tab = props.openedTabs.find(t => t.key === props.activeTab)
        if (tab && tab.content !== content) {
          tab.content = content
          emit('contentChange', { key: tab.key, content: tab.content })
        }
      }
    })
  ]
}

function showContextMenu(event, view) {
  const { clientX, clientY } = event
  
  // 创建右键菜单
  const menu = document.createElement('div')
  menu.className = 'context-menu'
  menu.style.cssText = `
    position: fixed;
    top: ${clientY}px;
    left: ${clientX}px;
    background: #2d2d30;
    border: 1px solid #3c3c3c;
    border-radius: 4px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.3);
    z-index: 1000;
    min-width: 160px;
    padding: 4px 0;
  `
  
  const menuItems = [
    { label: '保存 (Ctrl+S)', action: () => saveCurrentFile() },
    { type: 'separator' },
    { label: '剪切 (Ctrl+X)', action: () => document.execCommand('cut') },
    { label: '复制 (Ctrl+C)', action: () => document.execCommand('copy') },
    { label: '粘贴 (Ctrl+V)', action: () => document.execCommand('paste') },
    { type: 'separator' },
    { label: '全选 (Ctrl+A)', action: () => view.dispatch({ selection: { anchor: 0, head: view.state.doc.length } }) }
  ]
  
  menuItems.forEach(item => {
    if (item.type === 'separator') {
      const separator = document.createElement('div')
      separator.style.cssText = `
        height: 1px;
        background: #3c3c3c;
        margin: 4px 0;
      `
      menu.appendChild(separator)
    } else {
      const menuItem = document.createElement('div')
      menuItem.className = 'context-menu-item'
      menuItem.textContent = item.label
      menuItem.style.cssText = `
        padding: 8px 16px;
        cursor: pointer;
        color: #cccccc;
        font-size: 14px;
        user-select: none;
      `
      
      menuItem.addEventListener('mouseenter', () => {
        menuItem.style.background = '#3c3c3c'
      })
      
      menuItem.addEventListener('mouseleave', () => {
        menuItem.style.background = 'transparent'
      })
      
      menuItem.addEventListener('click', () => {
        try {
          item.action()
        } catch (e) {
          console.warn('菜单操作失败:', e)
        }
        document.body.removeChild(menu)
      })
      
      menu.appendChild(menuItem)
    }
  })
  
  // 添加菜单到页面
  document.body.appendChild(menu)
  
  // 点击其他地方关闭菜单
  const closeMenu = (e) => {
    if (!menu.contains(e.target)) {
      if (document.body.contains(menu)) {
        document.body.removeChild(menu)
      }
      document.removeEventListener('click', closeMenu)
      document.removeEventListener('contextmenu', closeMenu)
    }
  }
  
  // 延迟添加事件监听，避免立即触发
  setTimeout(() => {
    document.addEventListener('click', closeMenu)
    document.addEventListener('contextmenu', closeMenu)
  }, 0)
}

watch(() => [props.activeTab, props.openedTabs, currentFileContent.value], async () => {
  const tab = props.openedTabs.find(t => t.key === props.activeTab)
  if (tab && editorView) {
    const currentContent = editorView.state.doc.toString()
    if (currentContent !== tab.content) {
      updateEditorContent(tab.content)
    }
  }
})

watch(currentFileContent, (val) => {
  const tab = props.openedTabs.find(t => t.key === props.activeTab)
  if (tab && tab.content !== val) {
    tab.content = val
    if (editorView) {
      const currentContent = editorView.state.doc.toString()
      if (currentContent !== val) {
        updateEditorContent(val)
      }
    }
  }
})

onMounted(async () => {
  await nextTick()

  // 创建 CodeMirror 编辑器
  const tab = props.openedTabs.find(t => t.key === props.activeTab)
  
  const state = EditorState.create({
    doc: tab ? tab.content : '',
    extensions: createEditorExtensions()
  })

  editorView = new EditorView({
    state,
    parent: editorContainer.value
  })

  window.addEventListener('keydown', handleKeydown)
})

onBeforeUnmount(() => {
  if (editorView) {
    editorView.destroy()
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

.editor-tabs {
  margin-bottom: 4px;
}

.codemirror-container {
  flex: 1;
  min-height: 400px;
  background: #1e1e1e;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.03);
  overflow: hidden;
}

/* CodeMirror 基础样式 */
:deep(.cm-editor) {
  height: 100%;
  font-size: 15px;
  color: #cccccc;
}

:deep(.cm-editor .cm-scroller) {
  font-family: 'Fira Code', 'Consolas', 'Monaco', monospace;
}

:deep(.cm-editor .cm-content) {
  padding: 16px;
  background: #1e1e1e;
}

:deep(.cm-editor .cm-line) {
  padding: 0;
}

:deep(.cm-editor .cm-cursor) {
  border-left: 2px solid #007acc;
}

:deep(.cm-editor .cm-selectionBackground) {
  background: rgba(0, 122, 204, 0.3);
}

/* 右键菜单样式 */
.context-menu {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.context-menu-item:hover {
  background: #3c3c3c !important;
}
</style> 