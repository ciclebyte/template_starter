<template>
  <div class="step-variables">
    <div class="variables-content">
      <!-- 变量配置表单 -->
      <div class="variables-form">
        <h2 class="form-title">配置项目变量</h2>
        <p class="form-desc">请填写以下信息来配置您的项目</p>
        
        <!-- 基本信息 -->
        <div class="variable-section">
          <h3 class="section-title">
            <n-icon><Person /></n-icon>
            基本信息
          </h3>
          <div class="variable-grid">
            <div class="variable-item">
              <label>项目名称 *</label>
              <n-input 
                v-model:value="formData.project_name" 
                placeholder="请输入项目名称"
                :status="getValidationStatus('project_name')"
              />
              <div class="variable-desc">项目名称，用于package.json等配置文件</div>
            </div>
            
            <div class="variable-item">
              <label>项目描述</label>
              <n-input 
                v-model:value="formData.project_description" 
                placeholder="请输入项目描述"
                type="textarea"
                :rows="3"
              />
              <div class="variable-desc">项目的详细描述信息</div>
            </div>
            
            <div class="variable-item">
              <label>作者 *</label>
              <n-input 
                v-model:value="formData.author" 
                placeholder="请输入作者姓名"
                :status="getValidationStatus('author')"
              />
              <div class="variable-desc">项目作者姓名</div>
            </div>
            
            <div class="variable-item">
              <label>作者邮箱</label>
              <n-input 
                v-model:value="formData.author_email" 
                placeholder="请输入作者邮箱"
                type="email"
              />
              <div class="variable-desc">作者联系邮箱</div>
            </div>
            
            <div class="variable-item">
              <label>GitHub用户名</label>
              <n-input 
                v-model:value="formData.author_github" 
                placeholder="请输入GitHub用户名"
              />
              <div class="variable-desc">GitHub用户名，用于生成链接</div>
            </div>
          </div>
        </div>

        <!-- 项目配置 -->
        <div class="variable-section">
          <h3 class="section-title">
            <n-icon><Settings /></n-icon>
            项目配置
          </h3>
          <div class="variable-grid">
            <div class="variable-item">
              <label>包名</label>
              <n-input 
                v-model:value="formData.package_name" 
                placeholder="com.example.project"
              />
              <div class="variable-desc">项目包名，如：com.example.project</div>
            </div>
            
            <div class="variable-item">
              <label>模块名</label>
              <n-input 
                v-model:value="formData.module_name" 
                placeholder="main"
              />
              <div class="variable-desc">主要模块名称</div>
            </div>
            
            <div class="variable-item">
              <label>命名空间</label>
              <n-input 
                v-model:value="formData.namespace" 
                placeholder="Example"
              />
              <div class="variable-desc">代码命名空间</div>
            </div>
            
            <div class="variable-item">
              <label>端口号</label>
              <n-input 
                v-model:value="formData.port" 
                placeholder="8080"
                type="number"
              />
              <div class="variable-desc">服务运行端口</div>
            </div>
          </div>
        </div>

        <!-- 自定义变量 -->
        <div class="variable-section" v-if="customVariables.length > 0">
          <h3 class="section-title">
            <n-icon><CodeSlash /></n-icon>
            自定义变量
          </h3>
          <div class="variable-grid">
            <div 
              v-for="variable in customVariables" 
              :key="variable.name"
              class="variable-item"
            >
              <label>{{ variable.label || variable.name }}</label>
              <n-input 
                v-model:value="formData[variable.name]" 
                :placeholder="variable.description || '请输入值'"
                :status="getValidationStatus(variable.name)"
              />
              <div class="variable-desc">{{ variable.description }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧预览 -->
      <div class="variables-preview">
        <h3 class="preview-title">实时预览</h3>
        <div class="preview-content">
          <div class="preview-file">
            <div class="file-header">
              <span class="file-name">package.json</span>
            </div>
            <div class="file-content">
              <pre>{{ previewContent }}</pre>
            </div>
          </div>
        </div>
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
        @click="handleNext"
        :disabled="!isFormValid"
      >
        下一步
        <template #icon>
          <n-icon><ArrowForward /></n-icon>
        </template>
      </n-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { 
  Person, 
  Settings, 
  CodeSlash, 
  ArrowBack, 
  ArrowForward 
} from '@vicons/ionicons5'

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

const emit = defineEmits(['prev', 'next', 'update-variables'])

// 表单数据
const formData = ref({
  project_name: '',
  project_description: '这是一个示例项目',
  author: '',
  author_email: 'developer@example.com',
  author_github: 'developer',
  package_name: 'com.example.project',
  module_name: 'main',
  namespace: 'Example',
  port: '8080'
})

// 自定义变量（模拟数据）
const customVariables = ref([
  {
    name: 'custom_feature',
    label: '自定义功能',
    description: '是否启用自定义功能'
  }
])

// 表单验证
const validationRules = {
  project_name: {
    required: true,
    pattern: /^[a-zA-Z0-9_-]+$/,
    message: '项目名称只能包含字母、数字、下划线和连字符'
  },
  author: {
    required: true,
    message: '作者姓名不能为空'
  }
}

// 验证状态
const getValidationStatus = (fieldName) => {
  const rule = validationRules[fieldName]
  if (!rule) return undefined
  
  const value = formData.value[fieldName]
  
  if (rule.required && !value) {
    return 'error'
  }
  
  if (rule.pattern && value && !rule.pattern.test(value)) {
    return 'error'
  }
  
  return undefined
}

// 表单是否有效
const isFormValid = computed(() => {
  for (const [fieldName, rule] of Object.entries(validationRules)) {
    const value = formData.value[fieldName]
    
    if (rule.required && !value) {
      return false
    }
    
    if (rule.pattern && value && !rule.pattern.test(value)) {
      return false
    }
  }
  
  return true
})

// 预览内容
const previewContent = computed(() => {
  const data = formData.value
  return `{
  "name": "${data.project_name || 'my-project'}",
  "version": "1.0.0",
  "description": "${data.project_description || '这是一个示例项目'}",
  "author": "${data.author || '开发者'}",
  "email": "${data.author_email || 'developer@example.com'}",
  "github": "${data.author_github || 'developer'}",
  "main": "index.js",
  "scripts": {
    "start": "node index.js",
    "dev": "nodemon index.js"
  },
  "dependencies": {
    "express": "^4.18.2"
  },
  "devDependencies": {
    "nodemon": "^3.0.1"
  }
}`
})

// 监听表单变化
watch(formData, (newData) => {
  emit('update-variables', newData)
}, { deep: true })

// 初始化表单数据
watch(() => props.variables, (newVariables) => {
  if (newVariables && Object.keys(newVariables).length > 0) {
    formData.value = { ...formData.value, ...newVariables }
  }
}, { immediate: true })

// 下一步
const handleNext = () => {
  if (isFormValid.value) {
    emit('update-variables', formData.value)
    emit('next')
  }
}
</script>

<style scoped>
.step-variables {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.variables-content {
  flex: 1;
  display: flex;
  gap: 32px;
  padding: 40px;
  overflow: hidden;
}

.variables-form {
  flex: 1;
  overflow-y: auto;
}

.form-title {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  margin-bottom: 8px;
}

.form-desc {
  color: #666;
  margin-bottom: 32px;
}

.variable-section {
  margin-bottom: 40px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.variable-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.variable-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.variable-item label {
  font-weight: bold;
  color: #333;
  font-size: 14px;
}

.variable-desc {
  font-size: 12px;
  color: #999;
  line-height: 1.4;
}

.variables-preview {
  width: 400px;
  border-left: 1px solid #f0f0f0;
  padding-left: 32px;
}

.preview-title {
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin-bottom: 16px;
}

.preview-content {
  background: #f8f9fa;
  border-radius: 8px;
  overflow: hidden;
}

.preview-file {
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  overflow: hidden;
}

.file-header {
  background: #f0f0f0;
  padding: 8px 12px;
  border-bottom: 1px solid #e0e0e0;
}

.file-name {
  font-size: 12px;
  color: #666;
  font-family: 'Monaco', 'Menlo', monospace;
}

.file-content {
  padding: 16px;
  max-height: 400px;
  overflow-y: auto;
}

.file-content pre {
  margin: 0;
  font-size: 12px;
  line-height: 1.4;
  color: #333;
  font-family: 'Monaco', 'Menlo', monospace;
  white-space: pre-wrap;
}

.step-actions {
  padding: 24px 40px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
  display: flex;
  justify-content: space-between;
}
</style> 