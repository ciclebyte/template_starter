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
                <span class="variable-name">{{ variable.name }}</span>
                <span v-if="variable.description && variable.description !== '模板变量'" class="variable-desc-inline">- {{ variable.description }}</span>
                <span v-if="variable.isRequired" class="required-mark">*</span>
                <span class="variable-type-badge">{{ getVariableTypeLabel(variable.variableType) }}</span>
                <span v-if="variable.fileCount > 0" class="usage-indicator">
                  在{{ variable.fileCount }}个文件中使用
                </span>
              </label>
              
              <!-- 根据变量类型渲染不同的输入组件 -->
              <!-- 字符串类型 -->
              <n-input 
                v-if="variable.variableType === 'string'"
                v-model:value="formData[variable.name]" 
                :placeholder="variable.description || '请输入字符串'"
                :status="getValidationStatus(variable.name)"
              />
              
              <!-- 数字类型 -->
              <n-input-number 
                v-else-if="variable.variableType === 'number'"
                v-model:value="formData[variable.name]" 
                :placeholder="variable.description || '请输入数字'"
                :status="getValidationStatus(variable.name)"
                style="width: 100%"
              />
              
              <!-- 布尔类型 -->
              <n-switch 
                v-else-if="variable.variableType === 'boolean'"
                v-model:value="formData[variable.name]"
                :checked-value="true"
                :unchecked-value="false"
              >
                <template #checked>是</template>
                <template #unchecked>否</template>
              </n-switch>
              
              <!-- 列表类型 -->
              <div v-else-if="variable.variableType === 'list'" class="list-input">
                <n-dynamic-tags 
                  v-model:value="formData[variable.name]"
                  :placeholder="variable.description || '按回车添加项目'"
                />
                <div class="input-hint">按回车添加项目，点击标签删除</div>
              </div>
              
              <!-- 对象类型 -->
              <div v-else-if="variable.variableType === 'object'" class="object-input">
                <n-input 
                  v-model:value="formData[variable.name]" 
                  type="textarea"
                  :placeholder="variable.description || '请输入JSON格式的对象'"
                  :status="getValidationStatus(variable.name)"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                />
                <div class="input-hint">请输入有效的JSON格式，例如: {"key": "value"}</div>
              </div>
              
              <!-- 其他类型（兼容旧版本） -->
              <n-input 
                v-else
                v-model:value="formData[variable.name]" 
                :placeholder="variable.description || '请输入值'"
                :status="getValidationStatus(variable.name)"
                :type="getInputType(variable.variableType)"
              />
              
              <div class="variable-desc">
                {{ variable.description }}
                <div v-if="variable.files && variable.files.length > 0" class="usage-files">
                  使用文件: {{ variable.files.join(', ') }}
                </div>
              </div>
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
import { getTemplateExpose } from '@/api/templateExpose'
import { analyzeTemplateVariables } from '@/api/templates'

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
  
  // 对象类型的JSON格式验证
  const variable = customVariables.value.find(v => v.name === fieldName)
  if (variable?.variableType === 'object' && value) {
    try {
      JSON.parse(value)
    } catch {
      return 'error'
    }
  }
  
  return undefined
}

// 表单是否有效
const isFormValid = computed(() => {
  const rules = getValidationRules()
  
  for (const [fieldName, rule] of Object.entries(rules)) {
    const value = formData.value[fieldName]
    const variable = customVariables.value.find(v => v.name === fieldName)
    
    // 必填验证
    if (rule.required) {
      if (variable?.variableType === 'list') {
        if (!Array.isArray(value) || value.length === 0) {
          return false
        }
      } else if (variable?.variableType === 'boolean') {
        // 布尔类型不需要必填验证，因为它总是有值（true/false）
      } else if (!value) {
        return false
      }
    }
    
    // 正则验证
    if (rule.pattern && value && !rule.pattern.test(value)) {
      return false
    }
    
    // 对象类型的JSON格式验证
    if (variable?.variableType === 'object' && value) {
      try {
        JSON.parse(value)
      } catch {
        return false
      }
    }
  }
  
  return true
})

// 获取输入类型（兼容旧版本）
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

// 获取变量类型标签
const getVariableTypeLabel = (variableType) => {
  switch (variableType) {
    case 'string':
      return '字符串'
    case 'number':
      return '数字'
    case 'boolean':
      return '布尔值'
    case 'list':
      return '列表'
    case 'object':
      return '对象'
    default:
      return '文本'
  }
}

// 根据变量类型获取默认值
const getDefaultValue = (variable) => {
  if (variable.defaultValue !== undefined && variable.defaultValue !== null) {
    // 如果有设置默认值，根据类型转换
    switch (variable.variableType) {
      case 'boolean':
        return variable.defaultValue === 'true' || variable.defaultValue === true
      case 'number':
        return Number(variable.defaultValue) || 0
      case 'list':
        if (typeof variable.defaultValue === 'string') {
          try {
            const parsed = JSON.parse(variable.defaultValue)
            return Array.isArray(parsed) ? parsed : [variable.defaultValue]
          } catch {
            return variable.defaultValue.split(',').map(s => s.trim())
          }
        }
        return Array.isArray(variable.defaultValue) ? variable.defaultValue : []
      case 'object':
        if (typeof variable.defaultValue === 'string') {
          try {
            return JSON.parse(variable.defaultValue)
          } catch {
            return variable.defaultValue
          }
        }
        return variable.defaultValue
      default:
        return String(variable.defaultValue)
    }
  }
  
  // 没有默认值时，根据类型返回默认值
  switch (variable.variableType) {
    case 'boolean':
      return false
    case 'number':
      return 0
    case 'list':
      return []
    case 'object':
      return '{}'
    default:
      return ''
  }
}



// 加载模板变量数据
const loadTemplateVariables = async () => {
  if (!props.templateInfo?.id) {
    return
  }
  loading.value = true
  try {
    // 并行调用两个接口获取更全面的数据
    const [exposeRes, analyzeRes] = await Promise.all([
      getTemplateExpose({ templateId: props.templateInfo.id }),
      analyzeTemplateVariables(props.templateInfo.id)
    ])
    
    // 解析用户定义的变量结构
    let userDefinedVariables = []
    if (exposeRes.data?.code === 0 && exposeRes.data?.data?.templateExpose) {
      const fieldSchemaJson = exposeRes.data.data.templateExpose.fieldSchemaJson
      if (fieldSchemaJson) {
        try {
          const parsedSchema = JSON.parse(fieldSchemaJson)
          userDefinedVariables = convertSchemaToCustomVariables(parsedSchema)
        } catch (parseError) {
          console.error('解析变量定义失败:', parseError)
        }
      }
    }
    
    // 解析模板分析结果
    let detectedVariables = []
    let builtinVars = []
    let templateFuncs = []
    let stats = null
    
    if (analyzeRes.data?.code === 0 && analyzeRes.data?.data) {
      detectedVariables = analyzeRes.data.data.detectedVariables || []
      builtinVars = analyzeRes.data.data.builtinVariables || []
      templateFuncs = analyzeRes.data.data.templateFunctions || []
      stats = analyzeRes.data.data.statistics || null
    }
    
    // 合并用户定义的变量和检测到的变量
    const mergedVariables = mergeVariables(userDefinedVariables, detectedVariables)
    
    customVariables.value = mergedVariables
    builtinVariables.value = builtinVars
    templateFunctions.value = templateFuncs
    variableStatistics.value = stats
    
    // 初始化表单数据
    const newFormData = { ...formData.value }
    
    // 初始化自定义变量
    customVariables.value.forEach(variable => {
      if (!(variable.name in newFormData)) {
        newFormData[variable.name] = getDefaultValue(variable)
      }
    })
    
    // 初始化内置变量
    builtinVariables.value.forEach(variable => {
      const fieldName = variable.name.toLowerCase()
      if (!newFormData[fieldName]) {
        newFormData[fieldName] = ''
      }
    })
    
    // 一次性更新表单数据
    formData.value = newFormData
    
  } catch (error) {
    console.error('加载模板变量失败:', error)
    customVariables.value = []
    builtinVariables.value = []
    templateFunctions.value = []
    variableStatistics.value = null
  } finally {
    loading.value = false
  }
}

// 合并用户定义的变量和检测到的变量
const mergeVariables = (userDefined, detected) => {
  const mergedMap = new Map()
  
  // 首先添加用户定义的变量（有完整的类型和验证信息）
  userDefined.forEach(variable => {
    mergedMap.set(variable.name, {
      ...variable,
      isUserDefined: true,
      files: [],
      contexts: [],
      fileCount: 0
    })
  })
  
  // 然后添加或更新检测到的变量信息
  detected.forEach(detectedVar => {
    if (mergedMap.has(detectedVar.name)) {
      // 如果用户已定义该变量，更新使用信息
      const existing = mergedMap.get(detectedVar.name)
      existing.files = detectedVar.files || []
      existing.contexts = detectedVar.contexts || []
      existing.fileCount = (detectedVar.files || []).length
      existing.isDetected = true
    } else {
      // 如果用户未定义该变量，添加为新变量
      mergedMap.set(detectedVar.name, {
        name: detectedVar.name,
        variableType: convertTypeToOldFormat(detectedVar.type),
        description: detectedVar.suggestions || `模板中使用的 ${detectedVar.name} 变量`,
        defaultValue: getDefaultValueByType(detectedVar.type),
        isRequired: false,
        isUserDefined: false,
        isDetected: true,
        files: detectedVar.files || [],
        contexts: detectedVar.contexts || [],
        fileCount: (detectedVar.files || []).length
      })
    }
  })
  
  return Array.from(mergedMap.values())
}

// 根据类型获取默认值
const getDefaultValueByType = (type) => {
  switch (type) {
    case 'boolean':
      return false
    case 'number':
    case 'integer':
      return 0
    case 'array':
      return []
    case 'object':
      return '{}'
    default:
      return ''
  }
}

// 将schema格式转换为自定义变量列表格式
const convertSchemaToCustomVariables = (schema) => {
  const variablesList = []
  
  if (!schema || typeof schema !== 'object') {
    return variablesList
  }
  
  Object.entries(schema).forEach(([key, variable]) => {
    if (variable && typeof variable === 'object') {
      variablesList.push({
        name: key,
        variableType: convertTypeToOldFormat(variable.type),
        description: variable.description || variable.title || key,
        defaultValue: variable.default,
        isRequired: variable.required || false,
        validationRegex: variable.pattern
      })
    }
  })
  
  return variablesList
}


// 将新的类型格式转换为旧的格式
const convertTypeToOldFormat = (newType) => {
  switch (newType) {
    case 'string':
      return 'string'
    case 'integer':
    case 'number':
      return 'number'
    case 'boolean':
      return 'boolean'
    case 'array':
      return 'list'
    case 'object':
    case 'object_arr':
      return 'object'
    default:
      return 'string'
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
  min-height: 0;
}

.variables-content {
  flex: 1;
  padding: 40px;
  overflow-y: auto;
  overflow-x: hidden;
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
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
}

.variable-name {
  font-weight: bold;
  color: #333;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.variable-desc-inline {
  font-weight: normal;
  color: #666;
  font-size: 13px;
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

.usage-info {
  font-size: 11px;
  color: #18a058;
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





.variable-type-badge {
  display: inline-block;
  background: #f0f0f0;
  color: #666;
  font-size: 11px;
  font-weight: normal;
  padding: 2px 6px;
  border-radius: 3px;
  margin-left: 8px;
}

.list-input {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.object-input {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.input-hint {
  font-size: 11px;
  color: #999;
  font-style: italic;
}

/* 为不同变量类型的徽章添加颜色 */
.variable-type-badge {
  background: #e6f7ff;
  color: #1890ff;
}

.variable-item:has([data-variable-type="boolean"]) .variable-type-badge {
  background: #f6ffed;
  color: #52c41a;
}

.variable-item:has([data-variable-type="number"]) .variable-type-badge {
  background: #fff7e6;
  color: #fa8c16;
}

.variable-item:has([data-variable-type="list"]) .variable-type-badge {
  background: #f9f0ff;
  color: #722ed1;
}

.variable-item:has([data-variable-type="object"]) .variable-type-badge {
  background: #e6fffb;
  color: #13c2c2;
}

.usage-indicator {
  font-size: 11px;
  color: #18a058;
  font-weight: normal;
  margin-left: 8px;
  background: #f0f9f0;
  padding: 2px 6px;
  border-radius: 3px;
}

.step-actions {
  padding: 24px 40px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
  display: flex;
  justify-content: space-between;
}
</style> 