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
    
    <!-- HTML 预览弹框 -->
    <n-modal v-model:show="showHtmlPreviewModal" preset="card" :style="modalStyle">
      <template #header>
        <div class="modal-header">
          <span>HTML 预览</span>
          <div class="modal-actions">
            <n-button size="small" quaternary circle @click="toggleFullscreen">
              <template #icon>
                <svg v-if="!isFullscreen" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M8 3H5a2 2 0 0 0-2 2v3m18 0V5a2 2 0 0 0-2-2h-3m0 18h3a2 2 0 0 0 2-2v-3M3 16v3a2 2 0 0 0 2 2h3"/>
                </svg>
                <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M8 3v3a2 2 0 0 1-2 2H3m18 0h-3a2 2 0 0 1-2-2V3m0 18v-3a2 2 0 0 1 2-2h3M3 16h3a2 2 0 0 1 2 2v3"/>
                </svg>
              </template>
            </n-button>
          </div>
        </div>
      </template>
      <div class="html-preview-container">
        <iframe 
          ref="htmlPreviewFrame" 
          class="html-preview-frame"
          sandbox="allow-scripts allow-same-origin"
        ></iframe>
      </div>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount, nextTick, computed } from 'vue'
import { NTabs, NTab, useNotification } from 'naive-ui'
import { useTemplateFileStore } from '@/stores/templateFileStore'
import { storeToRefs } from 'pinia'
import { editTemplateFile } from '@/api/templateFiles'
import { useRoute } from 'vue-router'

// CodeMirror 核心模块 - 按照官方示例导入
import { EditorView, keymap, highlightSpecialChars, drawSelection, 
         highlightActiveLine, dropCursor, rectangularSelection, 
         crosshairCursor, lineNumbers, highlightActiveLineGutter } from '@codemirror/view'
import { EditorState } from '@codemirror/state'
import { defaultKeymap, indentWithTab, history, historyKeymap } from '@codemirror/commands'
import { defaultHighlightStyle, syntaxHighlighting, indentOnInput, 
         bracketMatching, foldGutter, foldKeymap } from '@codemirror/language'
import { searchKeymap, highlightSelectionMatches } from '@codemirror/search'
import { autocompletion, completionKeymap, closeBrackets, closeBracketsKeymap } from '@codemirror/autocomplete'

// Dracula 主题
import { dracula } from '@uiw/codemirror-theme-dracula'

// 语言支持
import { javascript } from '@codemirror/lang-javascript'
import { html } from '@codemirror/lang-html'
import { css } from '@codemirror/lang-css'
import { json } from '@codemirror/lang-json'
import { markdown } from '@codemirror/lang-markdown'
import { python } from '@codemirror/lang-python'
import { java } from '@codemirror/lang-java'
import { cpp } from '@codemirror/lang-cpp'
import { rust } from '@codemirror/lang-rust'
import { go } from '@codemirror/lang-go'
import { sql } from '@codemirror/lang-sql'
import { xml } from '@codemirror/lang-xml'
import { yaml } from '@codemirror/lang-yaml'
import { vue } from '@codemirror/lang-vue'

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
const htmlPreviewFrame = ref(null)
const showHtmlPreviewModal = ref(false)
const isFullscreen = ref(false)
let editorView = null

const templateFileStore = useTemplateFileStore()
const { currentFileContent } = storeToRefs(templateFileStore)

const route = useRoute()
const notification = useNotification()

// 计算属性
const modalStyle = computed(() => {
  if (isFullscreen.value) {
    return {
      width: '100vw',
      height: '100vh',
      margin: '0',
      borderRadius: '0'
    }
  } else {
    return {
      width: '90vw',
      height: '80vh'
    }
  }
})

// 语言映射
const languageMap = {
  'js': javascript(),
  'javascript': javascript(),
  'ts': javascript({typescript: true}),
  'typescript': javascript({typescript: true}),
  'jsx': javascript({jsx: true}),
  'tsx': javascript({typescript: true, jsx: true}),
  'vue': vue(),
  'html': html(),
  'htm': html(),
  'css': css(),
  'scss': css(),
  'sass': css(),
  'less': css(),
  'json': json(),
  'md': markdown(),
  'markdown': markdown(),
  'py': python(),
  'python': python(),
  'java': java(),
  'cpp': cpp(),
  'cc': cpp(),
  'cxx': cpp(),
  'c': cpp(),
  'rs': rust(),
  'rust': rust(),
  'go': go(),
  'sql': sql(),
  'xml': xml(),
  'yaml': yaml(),
  'yml': yaml()
}

function getLanguageExtension(filename) {
  const ext = filename.split('.').pop()?.toLowerCase()
  return languageMap[ext] || null
}

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



function runHtmlFile() {
  const tab = props.openedTabs.find(t => t.key === props.activeTab)
  if (!tab) return
  
  // 直接预览，不做校验
  showHtmlPreview(tab.content)
}

function showHtmlPreview(content) {
  // 显示预览弹框
  showHtmlPreviewModal.value = true
  isFullscreen.value = false // 重置全屏状态
  
  // 等待 DOM 更新后设置 iframe 内容
  nextTick(() => {
    if (htmlPreviewFrame.value) {
      const iframe = htmlPreviewFrame.value
      
      // 设置超时保护
      const timeout = setTimeout(() => {
        notification.error({
          title: '预览超时',
          content: 'HTML 渲染超时，可能存在无限循环或错误',
          duration: 3000
        })
        showHtmlPreviewModal.value = false
      }, 3000) // 3秒超时
      
      // 尝试渲染内容
      try {
        const doc = iframe.contentDocument || iframe.contentWindow.document
        
        // 写入 HTML 内容
        doc.open()
        doc.write(content)
        doc.close()
        
        // 清除超时
        clearTimeout(timeout)
        
        notification.success({
          title: '预览成功',
          content: 'HTML 内容已在弹框中显示',
          duration: 2500
        })
      } catch (e) {
        // 清除超时
        clearTimeout(timeout)
        
        // 如果渲染失败，显示纯文本
        try {
          const doc = iframe.contentDocument || iframe.contentWindow.document
          doc.open()
          doc.write(`<html><body><pre style="white-space: pre-wrap; font-family: monospace; padding: 20px;">${content}</pre></body></html>`)
          doc.close()
          
          notification.warning({
            title: '渲染失败',
            content: 'HTML 渲染失败，已显示为纯文本',
            duration: 3000
          })
        } catch (fallbackError) {
          // 最后的保护措施
          showHtmlPreviewModal.value = false
          notification.error({
            title: '预览失败',
            content: '无法显示内容',
            duration: 3000
          })
        }
      }
    }
  })
}

function toggleFullscreen() {
  isFullscreen.value = !isFullscreen.value
}

function handleKeydown(e) {
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    saveCurrentFile()
  }
  
  // HTML 文件运行快捷键
  if ((e.ctrlKey || e.metaKey) && e.key === 'r') {
    e.preventDefault()
    const tab = props.openedTabs.find(t => t.key === props.activeTab)
    if (tab) {
      const fileExt = tab.name.split('.').pop()?.toLowerCase()
      if (fileExt === 'html' || fileExt === 'htm') {
        runHtmlFile()
      }
    }
  }
}

function updateEditorContent(content, filename = '') {
  if (!editorView) return

  const languageExt = getLanguageExtension(filename)
  
  const newState = EditorState.create({
    doc: content || '',
    extensions: createEditorExtensions(languageExt)
  })

  editorView.setState(newState)
}

function createEditorExtensions(languageExtension = null) {
  const extensions = [
    // Dracula 主题
    dracula,
    // 启用编辑和选择功能
    EditorView.editable.of(true),
    // 行号
    lineNumbers(),
    // 代码折叠标记
    foldGutter(),
    // 高亮特殊字符
    highlightSpecialChars(),
    // 撤销历史
    history(),
    // 拖拽时的光标
    dropCursor(),
    // 输入时自动缩进
    indentOnInput(),
    // 括号匹配高亮
    bracketMatching(),
    // 自动关闭括号
    closeBrackets(),
    // 自动完成
    autocompletion(),
    // 矩形选择
    rectangularSelection(),
    // Alt+拖拽时显示十字光标
    crosshairCursor(),
    // 当前行高亮
    highlightActiveLine(),
    // 当前行行号高亮
    highlightActiveLineGutter(),
    // 高亮匹配的选中文本 - 重新添加，但使用更安全的配置
    highlightSelectionMatches(),
    
    // 快捷键配置
    keymap.of([
      // 保存快捷键
      {
        key: 'Ctrl-s',
        run: () => {
          saveCurrentFile()
          return true
        }
      },
      // 基础快捷键
      ...closeBracketsKeymap,
      ...defaultKeymap,
      ...searchKeymap,
      ...historyKeymap,
      ...foldKeymap,
      ...completionKeymap,
      // Tab 缩进
      indentWithTab
    ]),
    
    // 右键菜单
    EditorView.domEventHandlers({
      contextmenu: (event, view) => {
        event.preventDefault()
        showContextMenu(event, view)
        return true
      },
      mousedown: (event, view) => {
        // 确保点击时编辑器获得焦点
        view.focus()
        return false
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
  
  // 添加语言支持
  if (languageExtension) {
    extensions.push(languageExtension)
  }
  
  return extensions
}

function showContextMenu(event, view) {
  const { clientX, clientY } = event
  
  // 获取当前文件扩展名
  const tab = props.openedTabs.find(t => t.key === props.activeTab)
  const fileExt = tab ? tab.name.split('.').pop()?.toLowerCase() : ''
  
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
    { label: '撤销 (Ctrl+Z)', action: () => document.execCommand('undo') },
    { label: '重做 (Ctrl+Y)', action: () => document.execCommand('redo') },
    { type: 'separator' },
    { label: '剪切 (Ctrl+X)', action: () => document.execCommand('cut') },
    { label: '复制 (Ctrl+C)', action: () => document.execCommand('copy') },
    { label: '粘贴 (Ctrl+V)', action: () => document.execCommand('paste') },
    { type: 'separator' },
    { label: '全选 (Ctrl+A)', action: () => view.dispatch({ selection: { anchor: 0, head: view.state.doc.length } }) }
  ]
  
  // 如果是 HTML 文件，添加运行选项
  if (fileExt === 'html' || fileExt === 'htm') {
    menuItems.push(
      { type: 'separator' },
      { label: '运行 HTML (Ctrl+R)', action: () => runHtmlFile() }
    )
  }
  
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
        // 确保菜单被移除
        if (document.body.contains(menu)) {
          document.body.removeChild(menu)
        }
        // 移除事件监听器
        document.removeEventListener('click', closeMenu)
        document.removeEventListener('contextmenu', closeMenu)
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
      updateEditorContent(tab.content, tab.name)
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
        updateEditorContent(val, tab.name)
      }
    }
  }
})

onMounted(async () => {
  await nextTick()

  // 创建 CodeMirror 编辑器
  const tab = props.openedTabs.find(t => t.key === props.activeTab)
  const languageExt = tab ? getLanguageExtension(tab.name) : null
  
  const state = EditorState.create({
    doc: tab ? tab.content : '',
    extensions: createEditorExtensions(languageExt)
  })

  editorView = new EditorView({
    state,
    parent: editorContainer.value
  })

  // 自动聚焦编辑器
  setTimeout(() => {
    if (editorView) {
      editorView.focus()
    }
  }, 100)

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
  position: relative;
}

/* CodeMirror 基础样式 - 适配 Dracula 主题 */
:deep(.cm-editor) {
  height: 100% !important;
  font-size: 15px;
  outline: none !important;
}

:deep(.cm-editor .cm-scroller) {
  font-family: 'Fira Code', 'Consolas', 'Monaco', monospace;
}

:deep(.cm-editor .cm-content) {
  padding: 16px;
}

:deep(.cm-editor .cm-line) {
  padding: 0;
}

/* 选中样式 - 适配 Dracula 主题 */
:deep(.cm-selectionBackground) {
  background: rgba(68, 71, 90, 0.6) !important;
}

/* 匹配文本高亮 - 与选中文本匹配的其他位置 */
:deep(.cm-selectionMatch) {
  background: rgba(255, 184, 108, 0.3) !important;
  border-radius: 2px;
}

/* 确保匹配高亮不会覆盖主要选择 */
:deep(.cm-selectionBackground .cm-selectionMatch) {
  background: rgba(68, 71, 90, 0.6) !important;
}

:deep(.cm-editor .cm-activeLine) {
  background: rgba(255, 255, 255, 0.05) !important;
}

:deep(.cm-editor .cm-activeLineGutter) {
  background: rgba(255, 255, 255, 0.05) !important;
}

/* 确保光标可见 */
:deep(.cm-editor .cm-cursor) {
  border-left: 2px solid #007acc !important;
  border-right: none !important;
  width: 2px !important;
  background: transparent !important;
}

:deep(.cm-editor .cm-cursor-primary) {
  border-left: 2px solid #007acc !important;
  border-right: none !important;
  width: 2px !important;
  background: transparent !important;
}

/* 右键菜单样式 */
.context-menu {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.context-menu-item:hover {
  background: #3c3c3c !important;
}

/* HTML 预览弹框样式 */
.html-preview-container {
  width: 100%;
  height: 100%;
  min-height: 60vh;
}

.html-preview-frame {
  width: 100%;
  height: 100%;
  border: none;
  border-radius: 4px;
  background: #ffffff;
}

/* 弹框头部样式 */
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 8px 0;
}

.modal-actions {
  display: flex;
  gap: 4px;
}
</style> 