<template>
  <div 
    class="template-preview" 
    :class="{ collapsed: isCollapsed }"
  >
    <div class="preview-header">
      <div class="preview-title" v-if="!isCollapsed">预览</div>
      <div class="preview-actions">
        <n-button size="small" @click="copyContent" v-if="!isCollapsed">
          <template #icon>
            <n-icon><Copy /></n-icon>
          </template>
          复制
        </n-button>
        <n-button 
          size="small" 
          quaternary 
          @click="toggleCollapse"
          class="collapse-btn"
          :class="{ 'collapsed': isCollapsed }"
        >
          <template #icon>
            <n-icon>
              <ChevronForward v-if="!isCollapsed" />
              <ChevronBack v-else />
            </n-icon>
          </template>
        </n-button>
      </div>
    </div>
    
    <div v-show="!isCollapsed" class="preview-content">
      <div class="file-info">
        <div class="file-path">{{ currentFilePath }}</div>
      </div>
      
      <div class="file-content" ref="previewEditorRef"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { NButton, NIcon, useMessage } from 'naive-ui'
import { Copy, ChevronForward, ChevronBack } from '@vicons/ionicons5'
import { EditorView } from '@codemirror/view'
import { EditorState } from '@codemirror/state'
import { defaultHighlightStyle, syntaxHighlighting } from '@codemirror/language'
import { dracula } from '@uiw/codemirror-theme-dracula'

const props = defineProps({
  currentFile: {
    type: Object,
    default: null
  }
})

const message = useMessage()

const previewEditorRef = ref(null)
let previewEditor = null

// 折叠状态
const isCollapsed = ref(true)

// 计算文件路径
const currentFilePath = computed(() => {
  return props.currentFile?.filePath || props.currentFile?.fileName || '未选择文件'
})

// 更新预览内容
function updatePreviewContent() {
  if (!previewEditor || !props.currentFile) return
  
  const content = props.currentFile.fileContent || props.currentFile.content || ''
  const transaction = previewEditor.state.update({
    changes: {
      from: 0,
      to: previewEditor.state.doc.length,
      insert: content
    }
  })
  previewEditor.dispatch(transaction)
}

// 复制内容
function copyContent() {
  if (!previewEditor) return
  
  const content = previewEditor.state.doc.toString()
  navigator.clipboard.writeText(content).then(() => {
    message.success('内容已复制到剪贴板')
  })
}

// 切换折叠状态
function toggleCollapse() {
  isCollapsed.value = !isCollapsed.value
}

// 监听当前文件变化
watch(() => props.currentFile, () => {
  nextTick(() => {
    updatePreviewContent()
  })
})

// 初始化预览编辑器
onMounted(() => {
  if (previewEditorRef.value) {
    const state = EditorState.create({
      doc: '',
      extensions: [
        EditorView.theme({
          '&': { height: '100%' },
          '.cm-editor': { height: '100%' },
          '.cm-scroller': { fontFamily: 'monospace' }
        }),
        syntaxHighlighting(defaultHighlightStyle),
        dracula
      ]
    })
    
    previewEditor = new EditorView({
      state,
      parent: previewEditorRef.value
    })
    
    // 初始化内容
    if (props.currentFile) {
      updatePreviewContent()
    }
  }
})
</script>

<style scoped>
.template-preview {
  width: 400px;
  background: #fff;
  border-left: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  height: 100%;
  transition: width 0.3s ease;
}

.template-preview.collapsed {
  width: 32px;
}

.template-preview.collapsed .preview-header {
  padding: 8px 4px;
  justify-content: center;
}

.preview-header {
  padding: 16px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
  transition: padding 0.3s ease;
}

.preview-title {
  font-weight: bold;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.preview-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.collapse-btn {
  transition: all 0.3s ease;
}

.collapse-btn.collapsed {
  width: 24px;
  height: 24px;
  padding: 0;
  border-radius: 4px;
}

.collapse-btn.collapsed:hover {
  background-color: #f0f0f0;
}

.preview-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.file-info {
  padding: 12px 16px;
  border-bottom: 1px solid #e0e0e0;
  background: #f8f9fa;
}

.file-path {
  font-size: 13px;
  color: #666;
  word-break: break-all;
}

.file-content {
  flex: 1;
  overflow: hidden;
}
</style> 