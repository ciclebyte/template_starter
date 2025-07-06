<template>
  <div class="step-variables">
    <div class="variables-content">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-container">
        <n-spin size="large">
          <template #description>
            正在加载模板变量信息...
          </template>
        </n-spin>
      </div>
      
      <!-- 变量配置表单 -->
      <div v-else class="variables-form">
        <h2 class="form-title">配置项目变量</h2>
        <p class="form-desc">请填写以下信息来配置您的项目</p>
        
        <!-- 自定义变量 -->
        <div class="variable-section" v-if="customVariables.length > 0">
          <h3 class="section-title">
            <n-icon><CodeSlash /></n-icon>
            自定义变量
            <span class="section-subtitle">({{ variableStatistics?.totalCustomVariables || 0 }} 个变量)</span>
          </h3>
          <div class="variable-grid">
            <div 
              v-for="variable in customVariables" 
              :key="variable.name"
              class="variable-item"
            >
              <label>
                {{ variable.description || variable.name }}
                <span v-if="variable.isRequired" class="required-mark">*</span>
              </label>
              <n-input 
                v-model:value="formData[variable.name]" 
                :placeholder="variable.description || '请输入值'"
                :status="getValidationStatus(variable.name)"
                :type="getInputType(variable.variableType)"
              />
              <div class="variable-desc">{{ variable.description }}</div>
            </div>
          </div>
        </div>

        <!-- 内置变量 -->
        <div class="variable-section" v-if="builtinVariables.length > 0">
          <h3 class="section-title">
            <n-icon><Settings /></n-icon>
            内置变量
            <span class="section-subtitle">({{ variableStatistics?.totalBuiltinVariables || 0 }} 个变量)</span>
          </h3>
          <div class="variable-grid">
            <div 
              v-for="variable in builtinVariables" 
              :key="variable.name"
              class="variable-item builtin-variable"
            >
              <label>
                {{ variable.label }}
                <span class="usage-count">(使用 {{ variable.usageCount }} 次)</span>
              </label>
              <n-input 
                v-model:value="formData[variable.name.toLowerCase()]" 
                :placeholder="variable.description"
              />
              <div class="variable-desc">
                {{ variable.description }}
                <div class="usage-files" v-if="variable.files.length > 0">
                  使用文件: {{ variable.files.join(', ') }}
                </div>
              </div>
            </div>
          </div>
        </div>



        <!-- 模板函数信息 -->
        <div class="variable-section" v-if="templateFunctions.length > 0">
          <n-collapse>
            <n-collapse-item title="模板函数信息" name="functions">
              <template #header>
                <div class="collapse-header">
                  <n-icon><CodeSlash /></n-icon>
                  <span>模板函数</span>
                  <span class="section-subtitle">({{ variableStatistics?.totalFunctions || 0 }} 个函数)</span>
                </div>
              </template>
              <div class="function-list">
                <div 
                  v-for="func in templateFunctions" 
                  :key="func.name"
                  class="function-item"
                >
                  <div class="function-info">
                    <span class="function-name">{{ func.label }}</span>
                    <span class="usage-count">(使用 {{ func.usageCount }} 次)</span>
                  </div>
                  <div class="function-desc">{{ func.description }}</div>
                  <div class="usage-files" v-if="func.files.length > 0">
                    使用文件: {{ func.files.join(', ') }}
                  </div>
                </div>
              </div>
            </n-collapse-item>
          </n-collapse>
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
import { ref, computed, watch, onMounted } from 'vue'
import { 
  Person, 
  Settings, 
  CodeSlash, 
  ArrowBack, 
  ArrowForward 
} from '@vicons/ionicons5'
import { getTemplateVariables } from '@/api/templates'

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
const formData = ref({})

// 变量数据
const customVariables = ref([])
const builtinVariables = ref([])
const templateFunctions = ref([])
const variableStatistics = ref(null)
const loading = ref(false)

// 动态验证规则
const getValidationRules = () => {
  const rules = {}
  
  // 为必填的自定义变量添加验证规则
  customVariables.value.forEach(variable => {
    if (variable.isRequired) {
      rules[variable.name] = {
        required: true,
        message: `${variable.description || variable.name}不能为空`
      }
      
      // 如果有验证正则表达式，添加模式验证
      if (variable.validationRegex) {
        try {
          rules[variable.name].pattern = new RegExp(variable.validationRegex)
          rules[variable.name].message = `${variable.description || variable.name}格式不正确`
        } catch (e) {
          console.warn('无效的验证正则表达式:', variable.validationRegex)
        }
      }
    }
  })
  
  return rules
}

// 验证状态
const getValidationStatus = (fieldName) => {
  const rules = getValidationRules()
  const rule = rules[fieldName]
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
  const rules = getValidationRules()
  
  for (const [fieldName, rule] of Object.entries(rules)) {
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

// 获取输入类型
const getInputType = (variableType) => {
  switch (variableType) {
    case 'email':
      return 'email'
    case 'number':
      return 'number'
    case 'textarea':
      return 'textarea'
    case 'password':
      return 'password'
    default:
      return 'text'
  }
}



// 加载模板变量数据
const loadTemplateVariables = async () => {
  console.log('开始加载模板变量，templateInfo:', props.templateInfo)
  if (!props.templateInfo?.id) {
    console.log('templateInfo.id 不存在，跳过加载')
    return
  }
  
  console.log('调用API获取模板变量，templateId:', props.templateInfo.id)
  loading.value = true
  try {
    const res = await getTemplateVariables(props.templateInfo.id)
    console.log('API返回结果:', res)
    if (res.data && res.data.data) {
      customVariables.value = res.data.data.customVariables || []
      builtinVariables.value = res.data.data.builtinVariables || []
      templateFunctions.value = res.data.data.templateFunctions || []
      variableStatistics.value = res.data.data.statistics
      
      console.log('解析后的变量数据:', {
        customVariables: customVariables.value,
        builtinVariables: builtinVariables.value,
        templateFunctions: templateFunctions.value,
        statistics: variableStatistics.value
      })
      
      // 初始化内置变量到表单数据，避免触发不必要的更新
      const newFormData = { ...formData.value }
      builtinVariables.value.forEach(variable => {
        const fieldName = variable.name.toLowerCase()
        if (!newFormData[fieldName]) {
          newFormData[fieldName] = ''
        }
      })
      
      // 初始化自定义变量到表单数据
      customVariables.value.forEach(variable => {
        if (!newFormData[variable.name]) {
          newFormData[variable.name] = variable.defaultValue || ''
        }
      })
      
      // 一次性更新表单数据
      formData.value = newFormData
    }
  } catch (error) {
    console.error('加载模板变量失败:', error)
  } finally {
    loading.value = false
  }
}

// 监听表单变化，使用防抖避免频繁触发
let updateTimeout = null
watch(formData, (newData) => {
  if (updateTimeout) {
    clearTimeout(updateTimeout)
  }
  updateTimeout = setTimeout(() => {
    emit('update-variables', newData)
  }, 100)
}, { deep: true })

// 初始化表单数据，只在组件初始化时执行一次
let isInitialized = false
watch(() => props.variables, (newVariables) => {
  if (!isInitialized && newVariables && Object.keys(newVariables).length > 0) {
    formData.value = { ...formData.value, ...newVariables }
    isInitialized = true
  }
}, { immediate: true })

// 监听templateInfo变化，当有数据时加载变量
let lastTemplateId = null
watch(() => props.templateInfo, (newTemplateInfo) => {
  console.log('StepVariables 监听到 templateInfo 变化:', newTemplateInfo)
  if (newTemplateInfo?.id && newTemplateInfo.id !== lastTemplateId) {
    lastTemplateId = newTemplateInfo.id
    // 延迟一点加载，确保组件完全渲染
    setTimeout(() => {
      loadTemplateVariables()
    }, 100)
  }
}, { immediate: false })

// 组件挂载时不立即加载，等待 templateInfo 传入
onMounted(() => {
  console.log('StepVariables 组件挂载，templateInfo:', props.templateInfo)
  // 如果已经有 templateInfo，则加载数据
  if (props.templateInfo?.id) {
    lastTemplateId = props.templateInfo.id
    setTimeout(() => {
      loadTemplateVariables()
    }, 100)
  }
})

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
  padding: 40px;
  overflow: hidden;
}

.loading-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.variables-form {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
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

.section-subtitle {
  font-size: 14px;
  font-weight: normal;
  color: #666;
  margin-left: 8px;
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

.builtin-variable {
  background: #f8f9fa;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #e9ecef;
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

.usage-count {
  font-size: 12px;
  color: #18a058;
  font-weight: normal;
}

.usage-files {
  font-size: 11px;
  color: #666;
  margin-top: 4px;
  font-style: italic;
}

.required-mark {
  color: #d03050;
  margin-left: 4px;
}

.collapse-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.function-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.function-item {
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.function-info {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.function-name {
  font-weight: bold;
  color: #333;
  font-size: 14px;
}

.function-desc {
  font-size: 12px;
  color: #666;
  line-height: 1.4;
  margin-bottom: 8px;
}





.step-actions {
  padding: 24px 40px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
  display: flex;
  justify-content: space-between;
}
</style> 