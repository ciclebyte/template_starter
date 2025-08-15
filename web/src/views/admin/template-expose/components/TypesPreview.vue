<template>
  <div class="types-preview">
    <div class="language-tabs">
      <n-tabs v-model:value="activeLanguage" size="small">
        <n-tab-pane name="go" tab="Go">
          <div class="code-preview">
            <div class="preview-header">
              <h5>Go 结构体定义</h5>
              <n-space>
                <n-button size="small" @click="copyCode('go')">
                  <template #icon>
                    <n-icon><CopyOutline /></n-icon>
                  </template>
                  复制
                </n-button>
                <n-button size="small" @click="downloadCode('go')">
                  <template #icon>
                    <n-icon><DownloadOutline /></n-icon>
                  </template>
                  下载
                </n-button>
              </n-space>
            </div>
            <n-scrollbar>
              <div class="code-editor">
                <pre class="code-content"><code>{{ goCode }}</code></pre>
              </div>
            </n-scrollbar>
          </div>
        </n-tab-pane>

        <n-tab-pane name="typescript" tab="TypeScript">
          <div class="code-preview">
            <div class="preview-header">
              <h5>TypeScript 接口定义</h5>
              <n-space>
                <n-button size="small" @click="copyCode('typescript')">
                  <template #icon>
                    <n-icon><CopyOutline /></n-icon>
                  </template>
                  复制
                </n-button>
                <n-button size="small" @click="downloadCode('typescript')">
                  <template #icon>
                    <n-icon><DownloadOutline /></n-icon>
                  </template>
                  下载
                </n-button>
              </n-space>
            </div>
            <n-scrollbar>
              <div class="code-editor">
                <pre class="code-content"><code>{{ typescriptCode }}</code></pre>
              </div>
            </n-scrollbar>
          </div>
        </n-tab-pane>

        <n-tab-pane name="python" tab="Python">
          <div class="code-preview">
            <div class="preview-header">
              <h5>Python 数据类定义</h5>
              <n-space>
                <n-button size="small" @click="copyCode('python')">
                  <template #icon>
                    <n-icon><CopyOutline /></n-icon>
                  </template>
                  复制
                </n-button>
                <n-button size="small" @click="downloadCode('python')">
                  <template #icon>
                    <n-icon><DownloadOutline /></n-icon>
                  </template>
                  下载
                </n-button>
              </n-space>
            </div>
            <n-scrollbar>
              <div class="code-editor">
                <pre class="code-content"><code>{{ pythonCode }}</code></pre>
              </div>
            </n-scrollbar>
          </div>
        </n-tab-pane>

        <n-tab-pane name="java" tab="Java">
          <div class="code-preview">
            <div class="preview-header">
              <h5>Java 类定义</h5>
              <n-space>
                <n-button size="small" @click="copyCode('java')">
                  <template #icon>
                    <n-icon><CopyOutline /></n-icon>
                  </template>
                  复制
                </n-button>
                <n-button size="small" @click="downloadCode('java')">
                  <template #icon>
                    <n-icon><DownloadOutline /></n-icon>
                  </template>
                  下载
                </n-button>
              </n-space>
            </div>
            <n-scrollbar>
              <div class="code-editor">
                <pre class="code-content"><code>{{ javaCode }}</code></pre>
              </div>
            </n-scrollbar>
          </div>
        </n-tab-pane>
      </n-tabs>
    </div>

    <!-- 类型统计 -->
    <div class="type-stats">
      <div class="stats-header">
        <h5>类型统计</h5>
      </div>
      <div class="stats-content">
        <div class="stats-grid">
          <div class="stat-item" v-for="stat in typeStats" :key="stat.type">
            <div class="stat-icon" :class="stat.type">
              <n-icon><component :is="getTypeIcon(stat.type)" /></n-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">{{ getTypeLabel(stat.type) }}</div>
              <div class="stat-count">{{ stat.count }} 个</div>
            </div>
          </div>
        </div>
        
        <div class="complexity-info">
          <div class="complexity-item">
            <span class="complexity-label">总变量数:</span>
            <span class="complexity-value">{{ totalVariables }}</span>
          </div>
          <div class="complexity-item">
            <span class="complexity-label">嵌套层级:</span>
            <span class="complexity-value">{{ maxDepth }} 层</span>
          </div>
          <div class="complexity-item">
            <span class="complexity-label">必填字段:</span>
            <span class="complexity-value">{{ requiredCount }} 个</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { NTabs, NTabPane, NScrollbar, NSpace, NButton, NIcon, useMessage } from 'naive-ui'
import { 
  CopyOutline, DownloadOutline, DocumentTextOutline, CodeOutline, 
  ToggleOutline, CalendarOutline, ListOutline, ArchiveOutline 
} from '@vicons/ionicons5'

const props = defineProps({
  variables: {
    type: Array,
    default: () => []
  }
})

const message = useMessage()
const activeLanguage = ref('go')

// 生成Go代码
const goCode = computed(() => {
  if (props.variables.length === 0) {
    return '// 暂无变量定义'
  }
  
  let code = `package main

import (
	"time"
)

// TemplateVariables 模板变量定义
type TemplateVariables struct {
${generateGoFields(props.variables, 1)}
}

`
  
  // 为对象类型生成子结构体
  const subStructs = generateGoSubStructs(props.variables)
  if (subStructs) {
    code += subStructs
  }
  
  return code
})

// 生成Go字段
const generateGoFields = (variables, indent = 1) => {
  return variables.map(variable => {
    const indentStr = '\t'.repeat(indent)
    const fieldName = capitalize(variable.name)
    const jsonTag = `json:"${variable.name}${variable.required ? '' : ',omitempty'}"`
    
    let goType = getGoType(variable)
    let comment = ''
    
    if (variable.description) {
      comment = ` // ${variable.description}`
    }
    
    return `${indentStr}${fieldName} ${goType} \`${jsonTag}\`${comment}`
  }).join('\n')
}

// 生成Go子结构体
const generateGoSubStructs = (variables) => {
  const subStructs = []
  
  const processVariable = (variable) => {
    if (variable.type === 'object' && variable.children && variable.children.length > 0) {
      const structName = capitalize(variable.name)
      let structCode = `// ${structName} ${variable.description || variable.displayName || variable.name}
type ${structName} struct {
${generateGoFields(variable.children, 1)}
}

`
      subStructs.push(structCode)
      
      // 递归处理子对象
      variable.children.forEach(processVariable)
    }
  }
  
  variables.forEach(processVariable)
  return subStructs.join('')
}

// 获取Go类型
const getGoType = (variable) => {
  const typeMap = {
    string: 'string',
    number: 'float64',
    boolean: 'bool',
    date: 'time.Time',
    file: 'string'
  }
  
  if (variable.type === 'array') {
    const itemType = variable.itemType || 'string'
    return `[]${typeMap[itemType] || itemType}`
  }
  
  if (variable.type === 'object') {
    return capitalize(variable.name)
  }
  
  return typeMap[variable.type] || 'interface{}'
}

// 生成TypeScript代码
const typescriptCode = computed(() => {
  if (props.variables.length === 0) {
    return '// 暂无变量定义'
  }
  
  let code = `// 模板变量接口定义

export interface TemplateVariables {
${generateTypeScriptFields(props.variables, 1)}
}

`
  
  // 为对象类型生成子接口
  const subInterfaces = generateTypeScriptSubInterfaces(props.variables)
  if (subInterfaces) {
    code += subInterfaces
  }
  
  return code
})

// 生成TypeScript字段
const generateTypeScriptFields = (variables, indent = 1) => {
  return variables.map(variable => {
    const indentStr = '  '.repeat(indent)
    const fieldName = variable.name
    const optional = variable.required ? '' : '?'
    const tsType = getTypeScriptType(variable)
    
    let comment = ''
    if (variable.description) {
      comment = `\n${indentStr}/** ${variable.description} */`
    }
    
    return `${comment}${comment ? '\n' + indentStr : indentStr}${fieldName}${optional}: ${tsType};`
  }).join('\n')
}

// 生成TypeScript子接口
const generateTypeScriptSubInterfaces = (variables) => {
  const subInterfaces = []
  
  const processVariable = (variable) => {
    if (variable.type === 'object' && variable.children && variable.children.length > 0) {
      const interfaceName = capitalize(variable.name)
      let interfaceCode = `/**
 * ${variable.description || variable.displayName || variable.name}
 */
export interface ${interfaceName} {
${generateTypeScriptFields(variable.children, 1)}
}

`
      subInterfaces.push(interfaceCode)
      
      // 递归处理子对象
      variable.children.forEach(processVariable)
    }
  }
  
  variables.forEach(processVariable)
  return subInterfaces.join('')
}

// 获取TypeScript类型
const getTypeScriptType = (variable) => {
  const typeMap = {
    string: 'string',
    number: 'number',
    boolean: 'boolean',
    date: 'string | Date',
    file: 'string'
  }
  
  if (variable.type === 'array') {
    const itemType = variable.itemType || 'string'
    return `${typeMap[itemType] || itemType}[]`
  }
  
  if (variable.type === 'object') {
    return capitalize(variable.name)
  }
  
  if (variable.validation?.enum && variable.validation.enum.length > 0) {
    return variable.validation.enum.map(val => `'${val}'`).join(' | ')
  }
  
  return typeMap[variable.type] || 'any'
}

// 生成Python代码
const pythonCode = computed(() => {
  if (props.variables.length === 0) {
    return '# 暂无变量定义'
  }
  
  let code = `from dataclasses import dataclass
from typing import Optional, List, Dict, Any, Union
from datetime import datetime

@dataclass
class TemplateVariables:
    """模板变量数据类"""
${generatePythonFields(props.variables, 1)}

`
  
  // 为对象类型生成子类
  const subClasses = generatePythonSubClasses(props.variables)
  if (subClasses) {
    code += subClasses
  }
  
  return code
})

// 生成Python字段
const generatePythonFields = (variables, indent = 1) => {
  return variables.map(variable => {
    const indentStr = '    '.repeat(indent)
    const fieldName = variable.name
    const pythonType = getPythonType(variable)
    const defaultValue = variable.required ? '' : ' = None'
    
    let comment = ''
    if (variable.description) {
      comment = `  # ${variable.description}`
    }
    
    return `${indentStr}${fieldName}: ${pythonType}${defaultValue}${comment}`
  }).join('\n')
}

// 生成Python子类
const generatePythonSubClasses = (variables) => {
  const subClasses = []
  
  const processVariable = (variable) => {
    if (variable.type === 'object' && variable.children && variable.children.length > 0) {
      const className = capitalize(variable.name)
      let classCode = `@dataclass
class ${className}:
    """${variable.description || variable.displayName || variable.name}"""
${generatePythonFields(variable.children, 1)}

`
      subClasses.push(classCode)
      
      // 递归处理子对象
      variable.children.forEach(processVariable)
    }
  }
  
  variables.forEach(processVariable)
  return subClasses.join('')
}

// 获取Python类型
const getPythonType = (variable) => {
  const typeMap = {
    string: 'str',
    number: 'float',
    boolean: 'bool',
    date: 'datetime',
    file: 'str'
  }
  
  if (variable.type === 'array') {
    const itemType = variable.itemType || 'string'
    return `List[${typeMap[itemType] || 'Any'}]`
  }
  
  if (variable.type === 'object') {
    return capitalize(variable.name)
  }
  
  const baseType = typeMap[variable.type] || 'Any'
  return variable.required ? baseType : `Optional[${baseType}]`
}

// 生成Java代码
const javaCode = computed(() => {
  if (props.variables.length === 0) {
    return '// 暂无变量定义'
  }
  
  let code = `package com.example.template;

import java.util.List;
import java.util.Date;
import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * 模板变量类
 */
public class TemplateVariables {
${generateJavaFields(props.variables, 1)}
}

`
  
  // 为对象类型生成子类
  const subClasses = generateJavaSubClasses(props.variables)
  if (subClasses) {
    code += subClasses
  }
  
  return code
})

// 生成Java字段
const generateJavaFields = (variables, indent = 1) => {
  return variables.map(variable => {
    const indentStr = '    '.repeat(indent)
    const fieldName = variable.name
    const javaType = getJavaType(variable)
    const jsonProperty = `@JsonProperty("${variable.name}")`
    
    let comment = ''
    if (variable.description) {
      comment = `${indentStr}/**\n${indentStr} * ${variable.description}\n${indentStr} */\n`
    }
    
    return `${comment}${indentStr}${jsonProperty}\n${indentStr}private ${javaType} ${fieldName};`
  }).join('\n\n')
}

// 生成Java子类
const generateJavaSubClasses = (variables) => {
  const subClasses = []
  
  const processVariable = (variable) => {
    if (variable.type === 'object' && variable.children && variable.children.length > 0) {
      const className = capitalize(variable.name)
      let classCode = `/**
 * ${variable.description || variable.displayName || variable.name}
 */
public class ${className} {
${generateJavaFields(variable.children, 1)}
}

`
      subClasses.push(classCode)
      
      // 递归处理子对象
      variable.children.forEach(processVariable)
    }
  }
  
  variables.forEach(processVariable)
  return subClasses.join('')
}

// 获取Java类型
const getJavaType = (variable) => {
  const typeMap = {
    string: 'String',
    number: 'Double',
    boolean: 'Boolean',
    date: 'Date',
    file: 'String'
  }
  
  if (variable.type === 'array') {
    const itemType = variable.itemType || 'string'
    return `List<${typeMap[itemType] || 'Object'}>`
  }
  
  if (variable.type === 'object') {
    return capitalize(variable.name)
  }
  
  return typeMap[variable.type] || 'Object'
}

// 类型统计
const typeStats = computed(() => {
  const stats = new Map()
  
  const countTypes = (variables) => {
    variables.forEach(variable => {
      const type = variable.type
      stats.set(type, (stats.get(type) || 0) + 1)
      
      if (variable.children) {
        countTypes(variable.children)
      }
    })
  }
  
  countTypes(props.variables)
  
  return Array.from(stats.entries()).map(([type, count]) => ({
    type,
    count
  })).sort((a, b) => b.count - a.count)
})

// 总变量数
const totalVariables = computed(() => {
  let count = 0
  
  const countVariables = (variables) => {
    variables.forEach(variable => {
      count++
      if (variable.children) {
        countVariables(variable.children)
      }
    })
  }
  
  countVariables(props.variables)
  return count
})

// 最大嵌套深度
const maxDepth = computed(() => {
  let depth = 0
  
  const calculateDepth = (variables, currentDepth = 1) => {
    depth = Math.max(depth, currentDepth)
    variables.forEach(variable => {
      if (variable.children && variable.children.length > 0) {
        calculateDepth(variable.children, currentDepth + 1)
      }
    })
  }
  
  if (props.variables.length > 0) {
    calculateDepth(props.variables)
  }
  
  return depth
})

// 必填字段数量
const requiredCount = computed(() => {
  let count = 0
  
  const countRequired = (variables) => {
    variables.forEach(variable => {
      if (variable.required) {
        count++
      }
      if (variable.children) {
        countRequired(variable.children)
      }
    })
  }
  
  countRequired(props.variables)
  return count
})

// 工具函数
const capitalize = (str) => {
  return str.charAt(0).toUpperCase() + str.slice(1)
}

const getTypeIcon = (type) => {
  const iconMap = {
    string: DocumentTextOutline,
    number: CodeOutline,
    boolean: ToggleOutline,
    date: CalendarOutline,
    array: ListOutline,
    object: ArchiveOutline,
    file: DocumentTextOutline
  }
  return iconMap[type] || DocumentTextOutline
}

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

// 复制代码
const copyCode = async (language) => {
  const codeMap = {
    go: goCode.value,
    typescript: typescriptCode.value,
    python: pythonCode.value,
    java: javaCode.value
  }
  
  try {
    await navigator.clipboard.writeText(codeMap[language])
    message.success(`${language.toUpperCase()} 代码已复制到剪贴板`)
  } catch (error) {
    message.error('复制失败')
  }
}

// 下载代码
const downloadCode = (language) => {
  const codeMap = {
    go: { code: goCode.value, ext: 'go' },
    typescript: { code: typescriptCode.value, ext: 'ts' },
    python: { code: pythonCode.value, ext: 'py' },
    java: { code: javaCode.value, ext: 'java' }
  }
  
  try {
    const { code, ext } = codeMap[language]
    const blob = new Blob([code], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `template-variables.${ext}`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    message.success(`${language.toUpperCase()} 代码下载成功`)
  } catch (error) {
    message.error('下载失败')
  }
}
</script>

<style scoped>
.types-preview {
  height: 100%;
  padding: 16px;
  display: flex;
  flex-direction: column;
}

.language-tabs {
  flex: 1;
  margin-bottom: 20px;
}

.code-preview {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.preview-header h5 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.code-editor {
  flex: 1;
  min-height: 400px;
}

.code-content {
  margin: 0;
  padding: 16px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.5;
  color: #333;
  background: #f8f8f8;
  border-radius: 4px;
  overflow-x: auto;
  white-space: pre;
}

.type-stats {
  border-top: 1px solid #f0f0f0;
  padding-top: 16px;
}

.stats-header h5 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.stats-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background: white;
  border-radius: 6px;
  border: 1px solid #f0f0f0;
}

.stat-icon {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-icon.string {
  background: #52c41a;
}

.stat-icon.number {
  background: #1890ff;
}

.stat-icon.boolean {
  background: #722ed1;
}

.stat-icon.object {
  background: #fa8c16;
}

.stat-icon.array {
  background: #eb2f96;
}

.stat-icon.date {
  background: #13c2c2;
}

.stat-icon.file {
  background: #666;
}

.stat-info {
  flex: 1;
}

.stat-label {
  display: block;
  font-size: 12px;
  color: #666;
}

.stat-count {
  display: block;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.complexity-info {
  display: flex;
  gap: 20px;
  padding: 12px;
  background: #f8f8f8;
  border-radius: 6px;
}

.complexity-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.complexity-label {
  font-size: 12px;
  color: #666;
}

.complexity-value {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}
</style>