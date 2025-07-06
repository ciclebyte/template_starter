<template>
  <div class="edit-editor">
    <!-- 文件头部 -->
    <div v-if="currentFileName" class="file-header">
      <div class="file-info">
        <span class="file-name">{{ currentFileName }}</span>
      </div>
      <div class="file-actions">
        <n-button size="small" @click="saveCurrentFile">
          <template #icon>
            <n-icon><Save /></n-icon>
          </template>
          保存
        </n-button>
        <n-button size="small" @click="triggerPreview">
          <template #icon>
            <n-icon><Eye /></n-icon>
          </template>
          预览
        </n-button>
      </div>
    </div>
    
    <!-- 编辑器容器 -->
    <div v-if="!currentFileName" class="no-file-selected">
      <div class="no-file-icon">
        <n-icon size="48" color="#ccc">
          <Document />
        </n-icon>
      </div>
      <div class="no-file-text">请选择左侧文件进行编辑</div>
    </div>
    <div v-else class="codemirror-container" ref="editorContainer"></div>
    
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
import { NButton, NIcon, useNotification } from 'naive-ui'
import { editTemplateFile } from '@/api/templateFiles'
import { Save, Eye, Document } from '@vicons/ionicons5'

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
  currentFileName: {
    type: String,
    default: ''
  },
  currentFileId: {
    type: [String, Number],
    default: ''
  },
  currentFileContent: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['contentChange', 'insertVariable', 'preview'])

const editorContainer = ref(null)
const htmlPreviewFrame = ref(null)
const showHtmlPreviewModal = ref(false)
const isFullscreen = ref(false)
let editorView = null

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

async function saveCurrentFile() {
  if (!props.currentFileId) {
    notification.warning({
      title: '无法保存',
      content: '请先选择一个文件',
      duration: 2500
    })
    return
  }
  
  try {
    const content = editorView ? editorView.state.doc.toString() : props.currentFileContent
    await editTemplateFile({
      id: props.currentFileId,
      fileContent: content
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

// 插入变量到编辑器
function insertVariable(template) {
  if (!editorView) return
  
  const { state, dispatch } = editorView
  const { selection } = state
  
  // 在光标位置插入变量模板
  const transaction = state.update({
    changes: {
      from: selection.main.head,
      insert: template
    }
  })
  
  dispatch(transaction)
  
  // 触发内容变化事件
  const content = editorView.state.doc.toString()
  emit('contentChange', { content })
}

// 触发预览
function triggerPreview() {
  if (!props.currentFileId) {
    notification.warning({
      title: '无法预览',
      content: '请先选择一个文件',
      duration: 2500
    })
    return
  }
  
  emit('preview', { fileId: props.currentFileId, fileName: props.currentFileName })
}

// 暴露方法给父组件
defineExpose({
  insertVariable
})

function runHtmlFile() {
  const content = editorView ? editorView.state.doc.toString() : props.currentFileContent
  if (!content) return
  
  // 直接预览，不做校验
  showHtmlPreview(content)
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

function updateEditorContent(content, filename = '') {
  if (!editorView) return

  const languageExt = getLanguageExtension(filename)
  
  const newState = EditorState.create({
    doc: content || '',
    extensions: createEditorExtensionsWithListener(languageExt)
  })

  editorView.setState(newState)
  
  // 确保滚动正常工作
  setTimeout(() => {
    if (editorView && editorView.scrollDOM) {
      editorView.requestMeasure()
    }
  }, 100)
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
    // 确保滚动正常工作
    EditorView.scrollMargins.of(() => ({ top: 10, bottom: 10 })),
    // 强制启用滚动
    EditorView.theme({
      "&": { height: "100%" },
      ".cm-scroller": { 
        overflow: "auto !important",
        fontFamily: "monospace"
      }
    }),
    // 语法高亮
    syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
    // 快捷键
    keymap.of([
      // 保存快捷键
      {
        key: "Ctrl-s",
        run: () => {
          saveCurrentFile()
          return true
        }
      },
      // 预览快捷键
      {
        key: "Ctrl-r",
        run: () => {
          triggerPreview()
          return true
        }
      },
      // 默认快捷键
      ...defaultKeymap,
      ...historyKeymap,
      ...foldKeymap,
      ...completionKeymap,
      ...closeBracketsKeymap,
      ...searchKeymap,
      { key: "Tab", run: indentWithTab }
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
    })
  ]

  // 添加语言支持
  if (languageExtension) {
    extensions.push(languageExtension)
  }

  return extensions
}

function createEditorExtensionsWithListener(languageExtension = null) {
  return [
    ...createEditorExtensions(languageExtension),
    EditorView.updateListener.of((update) => {
      if (update.docChanged) {
        const content = update.state.doc.toString()
        emit('contentChange', { content })
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
    { label: '撤销 (Ctrl+Z)', action: () => document.execCommand('undo') },
    { label: '重做 (Ctrl+Y)', action: () => document.execCommand('redo') },
    { type: 'separator' },
    { label: '剪切 (Ctrl+X)', action: () => document.execCommand('cut') },
    { label: '复制 (Ctrl+C)', action: () => document.execCommand('copy') },
    { label: '粘贴 (Ctrl+V)', action: () => document.execCommand('paste') },
    { type: 'separator' },
    { label: '全选 (Ctrl+A)', action: () => view.dispatch({ selection: { anchor: 0, head: view.state.doc.length } }) }
  ]
  
  // 添加预览选项
  menuItems.push(
    { type: 'separator' },
    { label: '预览 (Ctrl+R)', action: () => triggerPreview() }
  )
  
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

// 监听文件内容变化
watch(() => props.currentFileContent, (newContent) => {
  if (editorView && newContent !== editorView.state.doc.toString()) {
    updateEditorContent(newContent, props.currentFileName)
  }
}, { immediate: false })

// 监听文件名变化，创建或更新编辑器
watch(() => props.currentFileName, (newFileName) => {
  if (newFileName) {
    // 有文件名时，创建或更新编辑器
    if (editorView) {
      // 更新现有编辑器
      const languageExt = getLanguageExtension(newFileName)
      const newState = EditorState.create({
        doc: editorView.state.doc.toString(),
        extensions: createEditorExtensionsWithListener(languageExt)
      })
      editorView.setState(newState)
    } else if (editorContainer.value) {
      // 创建新编辑器
      const languageExt = getLanguageExtension(newFileName)
      const state = EditorState.create({
        doc: props.currentFileContent || '',
        extensions: createEditorExtensionsWithListener(languageExt)
      })
      editorView = new EditorView({
        state,
        parent: editorContainer.value
      })
    }
  } else {
    // 没有文件名时，销毁编辑器
    if (editorView) {
      editorView.destroy()
      editorView = null
    }
  }
})

onMounted(() => {
  // 只有在有文件名时才创建编辑器
  if (editorContainer.value && props.currentFileName) {
    const languageExt = getLanguageExtension(props.currentFileName)
    
    const state = EditorState.create({
      doc: props.currentFileContent || '',
      extensions: createEditorExtensionsWithListener(languageExt)
    })

    editorView = new EditorView({
      state,
      parent: editorContainer.value
    })
  }
})

onBeforeUnmount(() => {
  if (editorView) {
    editorView.destroy()
  }
})
</script>

<style scoped>
.edit-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
}

/* 文件头部样式 */
.file-header {
  height: 48px;
  background: #f8f9fa;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
}

.file-info {
  display: flex;
  align-items: center;
}

.file-name {
  font-size: 14px;
  font-weight: bold;
  color: #333;
}

.file-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 编辑器容器 */
.codemirror-container {
  flex: 1;
  overflow: auto;
  background: #1e1e1e;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.03);
  min-height: 400px;
}

/* 空状态样式 */
.no-file-selected {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #999;
  background: #fff;
}

.no-file-icon {
  margin-bottom: 16px;
}

.no-file-text {
  font-size: 16px;
  color: #999;
}

/* 模态框样式 */
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.modal-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.html-preview-container {
  height: 100%;
  overflow: hidden;
}

.html-preview-frame {
  width: 100%;
  height: 100%;
  border: none;
  border-radius: 4px;
}

/* 右键菜单样式 */
.context-menu {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.context-menu-item:hover {
  background: #3c3c3c !important;
}
</style> 