<template>
  <div class="step-preview">
    <div class="preview-content">
      <!-- 左侧文件列表 -->
      <div class="file-list">
        <h3 class="list-title">生成的文件</h3>
        <div class="file-tree">
          <div 
            v-for="file in fileList" 
            :key="file.id"
            class="file-item"
            :class="{ active: selectedFile?.id === file.id }"
            @click="selectFile(file)"
          >
            <n-icon size="16" :color="file.isDirectory ? '#ffa500' : '#18a058'">
              <Folder v-if="file.isDirectory" />
              <Document v-else />
            </n-icon>
            <span class="file-name">{{ file.fileName }}</span>
          </div>
        </div>
      </div>

      <!-- 右侧代码预览 -->
      <div class="code-preview">
        <div class="preview-header">
          <h3 class="preview-title">
            {{ selectedFile?.fileName || '选择文件预览' }}
          </h3>
          <div class="preview-actions">
            <n-button size="small" @click="copyContent" v-if="selectedFile">
              <template #icon>
                <n-icon><Copy /></n-icon>
              </template>
              复制
            </n-button>
          </div>
        </div>
        
        <div class="preview-editor" ref="previewEditorRef"></div>
      </div>
    </div>

    <!-- 底部操作 -->
    <div class="step-actions">
      <n-button @click="$emit('prev')">
        <template #icon>
          <n-icon><ArrowBack /></n-icon>
        </template>
        上一步
      </n-button>
      
      <n-button 
        type="primary" 
        size="large"
        @click="handleGenerate"
        :loading="generating"
      >
        <template #icon>
          <n-icon><Download /></n-icon>
        </template>
        生成项目
      </n-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { 
  Folder, 
  Document, 
  Copy, 
  ArrowBack, 
  Download 
} from '@vicons/ionicons5'
import { EditorView, highlightActiveLine, highlightActiveLineGutter, lineNumbers } from '@codemirror/view'
import { EditorState } from '@codemirror/state'
import { defaultHighlightStyle, syntaxHighlighting } from '@codemirror/language'
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

// API
import { renderTemplate } from '@/api/templateFiles'

const props = defineProps({
  templateInfo: {
    type: Object,
    default: null
  },
  variables: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['prev', 'generate'])

const message = useMessage()
const previewEditorRef = ref(null)
let previewEditor = null

// 状态
const selectedFile = ref(null)
const generating = ref(false)
const fileList = ref([])

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

// 获取语言扩展
function getLanguageExtension(filename) {
  const ext = filename.split('.').pop()?.toLowerCase()
  return languageMap[ext] || null
}

// 加载模板文件列表
const loadTemplateFiles = async () => {
  if (!props.templateInfo?.id) return
  
  try {
    // TODO: 调用获取模板文件列表的API
    // const res = await getTemplateFiles(props.templateInfo.id)
    // fileList.value = res.data.data.files || []
    
    // 模拟文件列表
    fileList.value = [
      { id: 1, fileName: 'package.json', filePath: 'package.json', isDirectory: false },
      { id: 2, fileName: 'README.md', filePath: 'README.md', isDirectory: false },
      { id: 3, fileName: 'src', filePath: 'src', isDirectory: true },
      { id: 4, fileName: 'main.js', filePath: 'src/main.js', isDirectory: false },
      { id: 5, fileName: 'App.vue', filePath: 'src/App.vue', isDirectory: false },
      { id: 6, fileName: 'index.html', filePath: 'index.html', isDirectory: false },
      { id: 7, fileName: 'vite.config.js', filePath: 'vite.config.js', isDirectory: false }
    ]
  } catch (error) {
    message.error('加载模板文件失败')
    console.error(error)
  }
}

// 选择文件
const selectFile = async (file) => {
  if (file.isDirectory) return
  
  selectedFile.value = file
  await loadFileContent(file)
}

// 加载文件内容 - 直接调用渲染API
const loadFileContent = async (file) => {
  try {
    // 调用渲染API获取文件内容
    const res = await renderTemplate({
      fileId: file.id,
      testVariables: props.variables
    })
    
    const content = res.data.data.fileContent || '文件内容为空'
    updatePreviewContent(content, file.fileName)
  } catch (error) {
    message.error('加载文件内容失败')
    console.error(error)
    
    // 如果API调用失败，显示错误信息
    updatePreviewContent(`// 加载失败: ${error.message}`, file.fileName)
  }
}

// 更新预览内容
const updatePreviewContent = (content, filename) => {
  if (!previewEditor) return
  
  const languageExt = getLanguageExtension(filename)
  
  const newState = EditorState.create({
    doc: content,
    extensions: [
      dracula,
      syntaxHighlighting(defaultHighlightStyle),
      lineNumbers(),
      highlightActiveLine(),
      highlightActiveLineGutter(),
      EditorView.scrollMargins.of(() => ({ top: 10, bottom: 10 })),
      EditorView.theme({
        "&": { height: "100%" },
        ".cm-scroller": { 
          overflow: "auto !important",
          height: "100% !important"
        }
      }),
      ...(languageExt ? [languageExt] : [])
    ]
  })
  
  previewEditor.setState(newState)
}

// 复制内容
const copyContent = () => {
  if (!previewEditor) return
  
  const content = previewEditor.state.doc.toString()
  navigator.clipboard.writeText(content).then(() => {
    message.success('内容已复制到剪贴板')
  })
}

// 生成项目
const handleGenerate = async () => {
  generating.value = true
  try {
    emit('generate')
  } finally {
    generating.value = false
  }
}

// 初始化预览编辑器
onMounted(() => {
  if (previewEditorRef.value) {
    const state = EditorState.create({
      doc: '选择文件查看预览内容',
      extensions: [
        dracula,
        syntaxHighlighting(defaultHighlightStyle),
        lineNumbers(),
        highlightActiveLine(),
        highlightActiveLineGutter(),
        EditorView.scrollMargins.of(() => ({ top: 10, bottom: 10 })),
        EditorView.theme({
          "&": { height: "100%" },
          ".cm-scroller": { 
            overflow: "auto !important",
            height: "100% !important"
          }
        })
      ]
    })
    
    previewEditor = new EditorView({
      state,
      parent: previewEditorRef.value
    })
  }
  
  // 加载模板文件列表
  loadTemplateFiles()
})

// 监听模板信息变化
watch(() => props.templateInfo, () => {
  loadTemplateFiles()
}, { deep: true })

// 监听变量变化，重新加载当前文件
watch(() => props.variables, () => {
  if (selectedFile.value) {
    loadFileContent(selectedFile.value)
  }
}, { deep: true })
</script>

<style scoped>
.step-preview {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.preview-content {
  flex: 1;
  display: flex;
  gap: 0;
  overflow: hidden;
}

.file-list {
  width: 300px;
  border-right: 1px solid #f0f0f0;
  background: #fafafa;
  display: flex;
  flex-direction: column;
}

.list-title {
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin: 0;
  padding: 20px 20px 16px;
  border-bottom: 1px solid #f0f0f0;
}

.file-tree {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 20px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
}

.file-item:hover {
  background: #f0f0f0;
}

.file-item.active {
  background: #e6f7ff;
  color: #1890ff;
}

.file-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.code-preview {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.preview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
  background: #fff;
}

.preview-title {
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin: 0;
}

.preview-editor {
  flex: 1;
  overflow: hidden;
}

.step-actions {
  padding: 24px 40px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
  display: flex;
  justify-content: space-between;
}
</style> 