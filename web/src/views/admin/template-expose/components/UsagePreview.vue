<template>
  <div class="usage-preview">
    <div class="usage-tabs">
      <n-tabs v-model:value="activeTab" size="small">
        <n-tab-pane name="template" tab="模板示例">
          <div class="template-usage">
            <div class="usage-section">
              <h5>当前文件中的变量使用</h5>
              <div class="file-variables" v-if="fileVariables.length > 0">
                <div
                  v-for="fileVar in fileVariables"
                  :key="fileVar.name"
                  class="file-variable-item"
                  :class="{ 'undefined': !isVariableDefined(fileVar.name) }"
                >
                  <div class="variable-header">
                    <span class="variable-name">{{ fileVar.name }}</span>
                    <span class="usage-count">{{ fileVar.count }} 次使用</span>
                    <n-tag
                      :type="isVariableDefined(fileVar.name) ? 'success' : 'warning'"
                      size="small"
                    >
                      {{ isVariableDefined(fileVar.name) ? '已定义' : '未定义' }}
                    </n-tag>
                  </div>
                  <div class="variable-locations">
                    <span
                      v-for="line in fileVar.lines.slice(0, 5)"
                      :key="line"
                      class="line-number"
                      @click="highlightLine(line)"
                    >
                      第{{ line }}行
                    </span>
                    <span v-if="fileVar.lines.length > 5" class="more-lines">
                      还有{{ fileVar.lines.length - 5 }}处...
                    </span>
                  </div>
                </div>
              </div>
              <n-empty v-else description="当前文件中没有使用任何变量" size="small" />
            </div>

            <div class="usage-section">
              <h5>建议的变量使用方式</h5>
              <div class="suggested-usage">
                <div
                  v-for="variable in variables"
                  :key="variable.id"
                  class="usage-example"
                >
                  <div class="example-header">
                    <strong>{{ variable.displayName || variable.name }}</strong>
                    <span class="variable-type">({{ getTypeLabel(variable.type) }})</span>
                  </div>
                  <div class="example-code">
                    <pre><code>{{ generateUsageExample(variable) }}</code></pre>
                  </div>
                  <div class="example-description" v-if="variable.description">
                    {{ variable.description }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </n-tab-pane>

        <n-tab-pane name="json" tab="JSON数据示例">
          <div class="json-example">
            <div class="example-section">
              <div class="section-header">
                <h5>示例数据结构</h5>
                <n-button size="small" @click="generateRandomData">
                  <template #icon>
                    <n-icon><RefreshOutline /></n-icon>
                  </template>
                  生成随机数据
                </n-button>
              </div>
              <n-scrollbar style="max-height: 400px;">
                <div class="json-editor">
                  <pre class="json-content">{{ formattedExampleData }}</pre>
                </div>
              </n-scrollbar>
            </div>

            <div class="example-actions">
              <n-space>
                <n-button size="small" @click="copyExampleData">
                  <template #icon>
                    <n-icon><CopyOutline /></n-icon>
                  </template>
                  复制数据
                </n-button>
                
                <n-button size="small" @click="downloadExampleData">
                  <template #icon>
                    <n-icon><DownloadOutline /></n-icon>
                  </template>
                  下载示例
                </n-button>
              </n-space>
            </div>
          </div>
        </n-tab-pane>

        <n-tab-pane name="api" tab="API文档">
          <div class="api-docs">
            <div class="api-section">
              <h5>变量接口文档</h5>
              <div class="api-endpoint">
                <div class="endpoint-header">
                  <span class="method">POST</span>
                  <span class="url">/api/v1/templates/:id/generate</span>
                </div>
                <div class="endpoint-description">
                  使用变量数据生成模板内容
                </div>
              </div>
              
              <div class="api-params">
                <h6>请求参数</h6>
                <div class="param-table">
                  <table>
                    <thead>
                      <tr>
                        <th>参数名</th>
                        <th>类型</th>
                        <th>必填</th>
                        <th>描述</th>
                        <th>示例</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="variable in flattenedVariables" :key="variable.path">
                        <td><code>{{ variable.path }}</code></td>
                        <td><span class="type-badge" :class="variable.type">{{ getTypeLabel(variable.type) }}</span></td>
                        <td>
                          <n-tag :type="variable.required ? 'error' : 'default'" size="small">
                            {{ variable.required ? '是' : '否' }}
                          </n-tag>
                        </td>
                        <td>{{ variable.description || '-' }}</td>
                        <td><code>{{ variable.example || '-' }}</code></td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </n-tab-pane>
      </n-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { NTabs, NTabPane, NEmpty, NTag, NButton, NIcon, NSpace, NScrollbar, useMessage } from 'naive-ui'
import { RefreshOutline, CopyOutline, DownloadOutline } from '@vicons/ionicons5'

const props = defineProps({
  variables: {
    type: Array,
    default: () => []
  },
  selectedFile: {
    type: Object,
    default: null
  },
  fileContent: {
    type: String,
    default: ''
  }
})

const message = useMessage()
const activeTab = ref('template')
const exampleData = ref({})

// 从文件内容中解析变量使用
const fileVariables = computed(() => {
  if (!props.fileContent) return []
  
  const variableRegex = /\{\{\s*\.(\w+(?:\.\w+)*)\s*\}\}/g
  const usages = new Map()
  const lines = props.fileContent.split('\n')
  
  lines.forEach((line, index) => {
    let match
    const lineNum = index + 1
    
    while ((match = variableRegex.exec(line)) !== null) {
      const variableName = match[1]
      
      if (!usages.has(variableName)) {
        usages.set(variableName, {
          name: variableName,
          count: 0,
          lines: []
        })
      }
      
      const usage = usages.get(variableName)
      usage.count++
      
      if (!usage.lines.includes(lineNum)) {
        usage.lines.push(lineNum)
      }
    }
  })
  
  return Array.from(usages.values()).sort((a, b) => b.count - a.count)
})

// 扁平化变量列表（用于API文档）
const flattenedVariables = computed(() => {
  const flattened = []
  
  const flattenVariable = (variable, parentPath = '') => {
    const currentPath = parentPath ? `${parentPath}.${variable.name}` : variable.name
    
    flattened.push({
      ...variable,
      path: currentPath
    })
    
    if (variable.children && variable.children.length > 0) {
      variable.children.forEach(child => {
        flattenVariable(child, currentPath)
      })
    }
  }
  
  props.variables.forEach(variable => {
    flattenVariable(variable)
  })
  
  return flattened
})

// 格式化的示例数据
const formattedExampleData = computed(() => {
  try {
    return JSON.stringify(exampleData.value, null, 2)
  } catch (error) {
    return '// 数据生成错误'
  }
})

// 检查变量是否已定义
const isVariableDefined = (variableName) => {
  return flattenedVariables.value.some(v => v.path === variableName)
}

// 获取类型标签
const getTypeLabel = (type) => {
  const typeMap = {
    string: '字符串',
    number: '数字',
    boolean: '布尔',
    date: '日期',
    array: '数组',
    object: '对象',
    file: '文件'
  }
  return typeMap[type] || type
}

// 生成变量使用示例
const generateUsageExample = (variable) => {
  let example = `{{.${variable.name}}}`
  
  switch (variable.type) {
    case 'object':
      if (variable.children && variable.children.length > 0) {
        const childExamples = variable.children.slice(0, 2).map(child => 
          `{{.${variable.name}.${child.name}}}`
        ).join('\n')
        example = childExamples
      }
      break
      
    case 'array':
      example = `{{range .${variable.name}}}
  {{.}}
{{end}}`
      break
      
    case 'boolean':
      example = `{{if .${variable.name}}}
  <!-- 当${variable.displayName}为true时显示 -->
{{end}}`
      break
  }
  
  return example
}

// 生成示例数据
const generateExampleData = () => {
  const data = {}
  
  const generateVariableData = (variable) => {
    switch (variable.type) {
      case 'string':
        if (variable.validation?.enum && variable.validation.enum.length > 0) {
          return variable.validation.enum[0]
        }
        return variable.example || variable.defaultValue || `示例${variable.displayName || variable.name}`
        
      case 'number':
        if (variable.validation?.enum && variable.validation.enum.length > 0) {
          return Number(variable.validation.enum[0])
        }
        return Number(variable.example) || Number(variable.defaultValue) || Math.floor(Math.random() * 100)
        
      case 'boolean':
        return variable.defaultValue !== undefined ? Boolean(variable.defaultValue) : Math.random() > 0.5
        
      case 'date':
        return variable.example || new Date().toISOString().split('T')[0]
        
      case 'array':
        const itemType = variable.itemType || 'string'
        const itemCount = Math.floor(Math.random() * 3) + 1
        return Array.from({ length: itemCount }, (_, i) => {
          switch (itemType) {
            case 'string': return `项目${i + 1}`
            case 'number': return i + 1
            case 'boolean': return i % 2 === 0
            default: return `项目${i + 1}`
          }
        })
        
      case 'object':
        const objData = {}
        if (variable.children) {
          variable.children.forEach(child => {
            objData[child.name] = generateVariableData(child)
          })
        }
        return objData
        
      case 'file':
        return variable.example || 'https://example.com/file.pdf'
        
      default:
        return variable.example || variable.defaultValue || '示例数据'
    }
  }
  
  props.variables.forEach(variable => {
    data[variable.name] = generateVariableData(variable)
  })
  
  return data
}

// 生成随机数据
const generateRandomData = () => {
  exampleData.value = generateExampleData()
}

// 高亮行
const highlightLine = (lineNumber) => {
  console.log('高亮第', lineNumber, '行')
  // 可以实现跳转到对应行的功能
}

// 复制示例数据
const copyExampleData = async () => {
  try {
    await navigator.clipboard.writeText(formattedExampleData.value)
    message.success('示例数据已复制到剪贴板')
  } catch (error) {
    message.error('复制失败')
  }
}

// 下载示例数据
const downloadExampleData = () => {
  try {
    const blob = new Blob([formattedExampleData.value], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'template-example-data.json'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    message.success('示例数据下载成功')
  } catch (error) {
    message.error('下载失败')
  }
}

// 初始化示例数据
watch(() => props.variables, () => {
  exampleData.value = generateExampleData()
}, { immediate: true, deep: true })
</script>

<style scoped>
.usage-preview {
  height: 100%;
  padding: 16px;
}

.usage-tabs {
  height: 100%;
}

.template-usage,
.json-example,
.api-docs {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.usage-section {
  margin-bottom: 24px;
}

.usage-section h5 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.file-variable-item {
  padding: 12px;
  margin-bottom: 8px;
  background: white;
  border-radius: 6px;
  border: 1px solid #e8e8e8;
}

.file-variable-item.undefined {
  border-color: #faad14;
  background: #fffbe6;
}

.variable-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.variable-name {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  font-weight: 600;
  color: #1890ff;
}

.usage-count {
  font-size: 12px;
  color: #666;
}

.variable-locations {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.line-number {
  display: inline-block;
  padding: 2px 6px;
  background: #f0f0f0;
  border-radius: 3px;
  font-size: 11px;
  color: #666;
  cursor: pointer;
}

.line-number:hover {
  background: #e6f7ff;
  color: #1890ff;
}

.more-lines {
  font-size: 11px;
  color: #999;
}

.suggested-usage {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.usage-example {
  padding: 12px;
  background: #f8f8f8;
  border-radius: 6px;
  border-left: 3px solid #1890ff;
}

.example-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.variable-type {
  font-size: 12px;
  color: #666;
}

.example-code {
  margin: 8px 0;
}

.example-code pre {
  margin: 0;
  padding: 8px;
  background: white;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.4;
  color: #333;
}

.example-description {
  font-size: 12px;
  color: #666;
  font-style: italic;
}

.example-section {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.section-header h5 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.json-editor {
  flex: 1;
}

.json-content {
  margin: 0;
  padding: 16px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.5;
  color: #333;
  background: white;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  overflow-x: auto;
}

.example-actions {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.api-section h5 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.api-endpoint {
  margin-bottom: 20px;
  padding: 12px;
  background: #f8f8f8;
  border-radius: 6px;
}

.endpoint-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.method {
  padding: 2px 8px;
  background: #52c41a;
  color: white;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

.url {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  color: #333;
}

.endpoint-description {
  font-size: 13px;
  color: #666;
}

.api-params h6 {
  margin: 0 0 8px 0;
  font-size: 13px;
  font-weight: 600;
  color: #333;
}

.param-table {
  overflow-x: auto;
}

.param-table table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}

.param-table th,
.param-table td {
  padding: 8px 12px;
  border: 1px solid #f0f0f0;
  text-align: left;
}

.param-table th {
  background: #fafafa;
  font-weight: 600;
  color: #333;
}

.param-table code {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  background: #f0f0f0;
  padding: 2px 4px;
  border-radius: 2px;
  font-size: 11px;
}

.type-badge {
  display: inline-block;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 11px;
  font-weight: 500;
  color: white;
}

.type-badge.string {
  background: #52c41a;
}

.type-badge.number {
  background: #1890ff;
}

.type-badge.boolean {
  background: #722ed1;
}

.type-badge.object {
  background: #fa8c16;
}

.type-badge.array {
  background: #eb2f96;
}

.type-badge.date {
  background: #13c2c2;
}

.type-badge.file {
  background: #666;
}
</style>