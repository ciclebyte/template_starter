<template>
  <div class="schema-preview">
    <div class="preview-content">
      <n-scrollbar>
        <div class="json-editor">
          <pre class="json-content">{{ formattedSchema }}</pre>
        </div>
      </n-scrollbar>
    </div>

    <div class="preview-actions">
      <n-space>
        <n-button size="small" @click="copyToClipboard">
          <template #icon>
            <n-icon><CopyOutline /></n-icon>
          </template>
          复制Schema
        </n-button>
        
        <n-button size="small" @click="downloadSchema">
          <template #icon>
            <n-icon><DownloadOutline /></n-icon>
          </template>
          下载Schema
        </n-button>
        
        <n-button size="small" @click="validateSchema">
          <template #icon>
            <n-icon><CheckmarkCircleOutline /></n-icon>
          </template>
          验证Schema
        </n-button>
      </n-space>
    </div>

    <!-- 验证结果 -->
    <div class="validation-result" v-if="validationResult">
      <n-alert
        :type="validationResult.valid ? 'success' : 'error'"
        :title="validationResult.valid ? '验证通过' : '验证失败'"
        closable
        @close="validationResult = null"
      >
        <template v-if="!validationResult.valid">
          <div class="error-details">
            <p v-for="error in validationResult.errors" :key="error.path">
              <strong>{{ error.path }}:</strong> {{ error.message }}
            </p>
          </div>
        </template>
        <template v-else>
          Schema格式正确，包含 {{ variables.length }} 个变量定义
        </template>
      </n-alert>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { NScrollbar, NSpace, NButton, NIcon, NAlert, useMessage } from 'naive-ui'
import { CopyOutline, DownloadOutline, CheckmarkCircleOutline } from '@vicons/ionicons5'

const props = defineProps({
  variables: {
    type: Array,
    default: () => []
  }
})

const message = useMessage()
const validationResult = ref(null)

// 生成JSON Schema
const formattedSchema = computed(() => {
  try {
    const schema = generateJsonSchema(props.variables)
    return JSON.stringify(schema, null, 2)
  } catch (error) {
    return `// Schema生成错误: ${error.message}`
  }
})

// 生成JSON Schema结构
const generateJsonSchema = (variables) => {
  const schema = {
    $schema: "http://json-schema.org/draft-07/schema#",
    type: "object",
    title: "模板变量Schema",
    description: "模板变量的数据结构定义",
    properties: {},
    required: []
  }
  
  variables.forEach(variable => {
    schema.properties[variable.name] = generateVariableSchema(variable)
    if (variable.required) {
      schema.required.push(variable.name)
    }
  })
  
  return schema
}

// 生成单个变量的Schema
const generateVariableSchema = (variable) => {
  const baseSchema = {
    title: variable.displayName || variable.name,
    description: variable.description || ''
  }
  
  // 根据类型设置不同的Schema属性
  switch (variable.type) {
    case 'string':
      baseSchema.type = 'string'
      if (variable.defaultValue) {
        baseSchema.default = variable.defaultValue
      }
      if (variable.example) {
        baseSchema.examples = [variable.example]
      }
      
      // 验证规则
      if (variable.validation) {
        if (variable.validation.minLength !== null && variable.validation.minLength !== undefined) {
          baseSchema.minLength = variable.validation.minLength
        }
        if (variable.validation.maxLength !== null && variable.validation.maxLength !== undefined) {
          baseSchema.maxLength = variable.validation.maxLength
        }
        if (variable.validation.pattern) {
          baseSchema.pattern = variable.validation.pattern
        }
        if (variable.validation.enum && variable.validation.enum.length > 0) {
          baseSchema.enum = variable.validation.enum
        }
      }
      break
      
    case 'number':
      baseSchema.type = 'number'
      if (variable.defaultValue !== null && variable.defaultValue !== undefined) {
        baseSchema.default = Number(variable.defaultValue)
      }
      if (variable.example) {
        baseSchema.examples = [Number(variable.example)]
      }
      if (variable.validation?.enum && variable.validation.enum.length > 0) {
        baseSchema.enum = variable.validation.enum.map(Number)
      }
      break
      
    case 'boolean':
      baseSchema.type = 'boolean'
      if (variable.defaultValue !== null && variable.defaultValue !== undefined) {
        baseSchema.default = Boolean(variable.defaultValue)
      }
      if (variable.example) {
        baseSchema.examples = [Boolean(variable.example)]
      }
      break
      
    case 'array':
      baseSchema.type = 'array'
      if (variable.itemType) {
        baseSchema.items = {
          type: variable.itemType
        }
      }
      if (variable.validation) {
        if (variable.validation.minItems !== null && variable.validation.minItems !== undefined) {
          baseSchema.minItems = variable.validation.minItems
        }
        if (variable.validation.maxItems !== null && variable.validation.maxItems !== undefined) {
          baseSchema.maxItems = variable.validation.maxItems
        }
      }
      if (variable.example) {
        try {
          baseSchema.examples = [JSON.parse(variable.example)]
        } catch {
          baseSchema.examples = [variable.example]
        }
      }
      break
      
    case 'object':
      baseSchema.type = 'object'
      baseSchema.properties = {}
      baseSchema.required = []
      
      if (variable.children && variable.children.length > 0) {
        variable.children.forEach(child => {
          baseSchema.properties[child.name] = generateVariableSchema(child)
          if (child.required) {
            baseSchema.required.push(child.name)
          }
        })
      }
      break
      
    case 'date':
      baseSchema.type = 'string'
      baseSchema.format = 'date-time'
      if (variable.example) {
        baseSchema.examples = [variable.example]
      }
      break
      
    case 'file':
      baseSchema.type = 'string'
      baseSchema.format = 'uri'
      if (variable.example) {
        baseSchema.examples = [variable.example]
      }
      break
      
    default:
      baseSchema.type = 'string'
  }
  
  return baseSchema
}

// 复制到剪贴板
const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(formattedSchema.value)
    message.success('Schema已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    message.error('复制失败')
  }
}

// 下载Schema
const downloadSchema = () => {
  try {
    const blob = new Blob([formattedSchema.value], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'template-variables-schema.json'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    message.success('Schema文件下载成功')
  } catch (error) {
    console.error('下载失败:', error)
    message.error('下载失败')
  }
}

// 验证Schema
const validateSchema = () => {
  try {
    const schema = JSON.parse(formattedSchema.value)
    const errors = []
    
    // 基本格式检查
    if (!schema.properties) {
      errors.push({ path: 'root', message: '缺少properties字段' })
    }
    
    // 检查每个变量的定义
    if (schema.properties) {
      Object.keys(schema.properties).forEach(key => {
        const prop = schema.properties[key]
        if (!prop.type) {
          errors.push({ path: key, message: '缺少type字段' })
        }
        if (!prop.title && !prop.description) {
          errors.push({ path: key, message: '缺少标题或描述' })
        }
      })
    }
    
    validationResult.value = {
      valid: errors.length === 0,
      errors: errors
    }
    
    if (errors.length === 0) {
      message.success('Schema验证通过')
    } else {
      message.error(`Schema验证失败：${errors.length}个错误`)
    }
  } catch (error) {
    validationResult.value = {
      valid: false,
      errors: [{ path: 'JSON', message: 'JSON格式错误：' + error.message }]
    }
    message.error('Schema验证失败：JSON格式错误')
  }
}
</script>

<style scoped>
.schema-preview {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.preview-content {
  flex: 1;
  overflow: hidden;
}

.json-editor {
  height: 100%;
  min-height: 400px;
}

.json-content {
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

.preview-actions {
  padding: 16px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
}

.validation-result {
  margin-top: 16px;
  padding: 0 16px;
}

.error-details p {
  margin: 4px 0;
  font-size: 14px;
}
</style>