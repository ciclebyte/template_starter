<template>
  <div class="variable-manager">
    <div class="variable-tabs">
      <n-tabs type="line" animated>
        <n-tab-pane name="variables" tab="变量管理">
          <div class="variable-content">
            <div class="variable-list">
              <div class="variable-header">
                <div class="header-info">
                  <span class="header-title">变量列表</span>
                  <span class="header-count">共 {{ templateVariables.length }} 个变量</span>
                </div>
                <div class="header-actions">
                  <n-button type="primary" size="small" @click="addVariable">
                    <template #icon>
                      <n-icon><Add /></n-icon>
                    </template>
                    新增变量
                  </n-button>
                </div>
              </div>
              
              <!-- 文本变量 -->
              <div class="variable-section">
                <div class="section-title">
                  <n-icon><DocumentText /></n-icon>
                  文本变量 ({{ textVariables.length }})
                </div>
                <div class="variable-grid">
                  <div v-for="variable in textVariables" :key="variable.id" class="variable-card">
                    <div class="variable-header-card">
                      <div class="variable-name">{{ variable.name }}</div>
                      <div class="variable-actions">
                        <n-button size="tiny" @click="insertVariable(variable)">插入</n-button>
                        <n-button size="tiny" @click="editVariable(variable)">编辑</n-button>
                        <n-button size="tiny" type="error" @click="deleteVariable(variable.id)">删除</n-button>
                      </div>
                    </div>
                    <div class="variable-desc">{{ variable.description }}</div>
                    <div class="variable-meta">
                      <n-tag size="small" :type="variable.isRequired === 1 ? 'error' : 'default'">
                        {{ variable.isRequired === 1 ? '必填' : '可选' }}
                      </n-tag>
                      <span class="meta-text" v-if="variable.defaultValue">默认值: {{ variable.defaultValue }}</span>
                      <span class="meta-text" v-if="variable.validationRegex">验证: {{ variable.validationRegex }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 条件变量 -->
              <div class="variable-section" v-if="conditionalVariables.length > 0">
                <div class="section-title">
                  <n-icon><CodeSlash /></n-icon>
                  条件变量 ({{ conditionalVariables.length }})
                </div>
                <div class="variable-grid">
                  <div v-for="variable in conditionalVariables" :key="variable.id" class="variable-card">
                    <div class="variable-header-card">
                      <div class="variable-name">{{ variable.name }}</div>
                      <div class="variable-actions">
                        <n-button size="tiny" @click="insertVariable(variable)">插入</n-button>
                        <n-button size="tiny" @click="editVariable(variable)">编辑</n-button>
                        <n-button size="tiny" type="error" @click="deleteVariable(variable.id)">删除</n-button>
                      </div>
                    </div>
                    <div class="variable-desc">{{ variable.description }}</div>
                    <div class="variable-meta">
                      <n-tag size="small" :type="variable.isRequired === 1 ? 'error' : 'default'">
                        {{ variable.isRequired === 1 ? '必填' : '可选' }}
                      </n-tag>
                      <span class="meta-text" v-if="variable.defaultValue">默认值: {{ variable.defaultValue }}</span>
                      <span class="meta-text" v-if="variable.validationRegex">验证: {{ variable.validationRegex }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 空状态 -->
              <div v-if="templateVariables.length === 0" class="empty-state">
                <div class="empty-icon">📝</div>
                <div class="empty-title">暂无变量</div>
                <div class="empty-desc">点击"新增变量"按钮开始创建变量</div>
              </div>
            </div>
          </div>
        </n-tab-pane>
        
        <n-tab-pane name="testData" tab="测试数据">
          <div class="test-data-content">
            <div class="test-data-header">
              <div class="header-info">
                <span class="header-title">测试数据设置</span>
                <span class="header-desc">设置变量测试值，用于预览模板效果</span>
              </div>
              <div class="header-actions">
                <n-button size="small" @click="resetTestData">重置</n-button>
                <n-button type="primary" size="small" @click="applyTestData">应用</n-button>
              </div>
            </div>
            
            <!-- 内置变量测试值 -->
            <div class="test-data-section">
              <div class="section-title">
                <n-icon><Settings /></n-icon>
                内置变量测试值
              </div>
              <div class="test-data-grid">
                <div 
                  v-for="variable in builtinVariables" 
                  :key="variable.name"
                  class="test-data-item"
                >
                  <label>{{ variable.label }}:</label>
                  <n-input 
                    v-model:value="testData[variable.name]" 
                    :placeholder="variable.description"
                    size="small"
                  />
                </div>
              </div>
            </div>
            
            <!-- 自定义变量测试值 -->
            <div class="test-data-section" v-if="templateVariables.length > 0">
              <div class="section-title">
                <n-icon><DocumentText /></n-icon>
                自定义变量测试值
              </div>
              <div class="test-data-grid">
                <div 
                  v-for="variable in templateVariables" 
                  :key="variable.id"
                  class="test-data-item"
                >
                  <label>{{ variable.name }}:</label>
                  <n-input 
                    v-model:value="testData[variable.name]" 
                    :placeholder="variable.description || '请输入测试值'"
                    size="small"
                  />
                </div>
              </div>
            </div>
          </div>
        </n-tab-pane>
      </n-tabs>
    </div>

    <!-- 变量编辑对话框 -->
    <n-modal v-model:show="showEditModal" preset="card" :title="editForm.id ? '编辑变量' : '新增变量'" style="width: 600px" @update:show="onModalShowChange">
      <n-form ref="formRef" :model="editForm" :rules="rules" label-placement="left" label-width="100px">
        <n-form-item label="变量名称 *" path="name">
          <n-input v-model:value="editForm.name" placeholder="请输入变量名称，如：project_name" />
          <template #feedback>
            <span style="font-size: 12px; color: #999;">变量名称用于在模板中引用，建议使用下划线命名</span>
          </template>
        </n-form-item>
        <n-form-item label="变量描述 *" path="description">
          <n-input v-model:value="editForm.description" placeholder="请输入变量描述，如：项目名称" />
          <template #feedback>
            <span style="font-size: 12px; color: #999;">变量描述用于说明变量的用途</span>
          </template>
        </n-form-item>
        <n-form-item label="变量类型 *" path="variableType">
          <n-select v-model:value="editForm.variableType" :options="typeOptions" />
          <template #feedback>
            <span style="font-size: 12px; color: #999;">文本变量用于普通文本，条件变量用于控制流程</span>
          </template>
        </n-form-item>
        <n-form-item label="默认值" path="defaultValue">
          <n-input v-model:value="editForm.defaultValue" placeholder="请输入默认值（可选）" />
          <template #feedback>
            <span style="font-size: 12px; color: #999;">当用户未输入值时使用的默认值</span>
          </template>
        </n-form-item>
        <n-form-item label="是否必填" path="isRequired">
          <n-switch v-model:value="editForm.isRequired" />
          <template #feedback>
            <span style="font-size: 12px; color: #999;">必填变量在生成模板时必须提供值</span>
          </template>
        </n-form-item>
        <n-form-item label="验证规则" path="validationRegex">
          <n-input v-model:value="editForm.validationRegex" placeholder="请输入验证正则表达式（可选）" />
          <template #feedback>
            <span style="font-size: 12px; color: #999;">用于验证用户输入值的格式，如：^[a-zA-Z0-9_]+$</span>
          </template>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="cancelEdit">取消</n-button>
          <n-button type="primary" @click="saveVariable">{{ editForm.id ? '更新' : '保存' }}</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { NButton, NIcon, NModal, NForm, NFormItem, NInput, NSelect, NSwitch, NSpace, NTag, useMessage } from 'naive-ui'
import { Add, DocumentText, CodeSlash, Settings } from '@vicons/ionicons5'

const props = defineProps({
  variables: {
    type: Array,
    default: () => []
  },
  templateId: {
    type: [String, Number],
    required: true
  }
})

const emit = defineEmits(['add', 'edit', 'delete', 'insert', 'applyTestData'])

const message = useMessage()
const formRef = ref(null)

// 变量类型选项
const typeOptions = [
  { label: '文本变量', value: 'text', description: '用于普通文本内容' },
  { label: '条件变量', value: 'conditional', description: '用于控制模板生成流程' }
]

// 计算属性：按类型分组变量
const textVariables = computed(() => {
  return props.variables.filter(v => v.variableType === 'text' || !v.variableType)
})

const conditionalVariables = computed(() => {
  return props.variables.filter(v => v.variableType === 'conditional')
})

const templateVariables = computed(() => {
  return props.variables
})

// 编辑表单
const showEditModal = ref(false)
const editForm = ref({
  id: null,
  name: '',
  description: '',
  variableType: 'text',
  defaultValue: '',
  isRequired: true,
  validationRegex: ''
})

const rules = {
  name: { required: true, message: '请输入变量名称' },
  description: { required: true, message: '请输入变量描述' },
  variableType: { required: true, message: '请选择变量类型' }
}

// 测试数据相关
const testData = ref({})

// 内置变量定义
const builtinVariables = [
  { name: 'project_name', label: '项目名', description: '项目名称' },
  { name: 'project_description', label: '项目描述', description: '项目的详细描述信息' },
  { name: 'author', label: '作者', description: '项目作者姓名' },
  { name: 'author_email', label: '作者邮箱', description: '作者联系邮箱' },
  { name: 'author_github', label: '作者GitHub', description: '作者GitHub用户名' },
  { name: 'current_time', label: '当前时间', description: '当前时间戳' },
  { name: 'package_name', label: '包名', description: '项目包名' },
  { name: 'module_name', label: '模块名', description: '模块名称' },
  { name: 'namespace', label: '命名空间', description: '代码命名空间' },
  { name: 'port', label: '端口号', description: '服务端口号' }
]

// 获取测试数据存储键名
const getTestDataKey = () => {
  return `template_test_data_${props.templateId}`
}

// 保存测试数据到本地存储
const saveTestDataToStorage = (data) => {
  try {
    localStorage.setItem(getTestDataKey(), JSON.stringify(data))
  } catch (error) {
    console.error('保存测试数据失败:', error)
  }
}

// 从本地存储加载测试数据
const loadTestDataFromStorage = () => {
  try {
    const savedData = localStorage.getItem(getTestDataKey())
    if (savedData) {
      return JSON.parse(savedData)
    }
  } catch (error) {
    console.error('加载测试数据失败:', error)
  }
  return null
}

// 初始化测试数据
const initTestData = () => {
  // 首先尝试从本地存储加载
  const savedData = loadTestDataFromStorage()
  
  if (savedData) {
    // 如果有保存的数据，使用保存的数据，但补充缺失的变量
    const data = { ...savedData }
    
    // 补充内置变量（如果保存的数据中没有）
    builtinVariables.forEach(variable => {
      if (!(variable.name in data)) {
        data[variable.name] = getDefaultValue(variable.name)
      }
    })
    
    // 补充自定义变量（如果保存的数据中没有）
    templateVariables.value.forEach(variable => {
      if (!(variable.name in data)) {
        data[variable.name] = variable.defaultValue || ''
      }
    })
    
    testData.value = data
  } else {
    // 如果没有保存的数据，使用默认值
    const data = {}
    
    // 初始化内置变量
    builtinVariables.forEach(variable => {
      data[variable.name] = getDefaultValue(variable.name)
    })
    
    // 初始化自定义变量
    templateVariables.value.forEach(variable => {
      data[variable.name] = variable.defaultValue || ''
    })
    
    testData.value = data
  }
}

// 获取默认值
const getDefaultValue = (variableName) => {
  const defaults = {
    project_name: 'my-project',
    project_description: '这是一个示例项目',
    author: '开发者',
    author_email: 'developer@example.com',
    author_github: 'developer',
    current_time: new Date().toLocaleString(),
    package_name: 'com.example.project',
    module_name: 'main',
    namespace: 'Example',
    port: '8080'
  }
  return defaults[variableName] || ''
}

// 重置测试数据
const resetTestData = () => {
  // 清除本地存储的测试数据
  try {
    localStorage.removeItem(getTestDataKey())
  } catch (error) {
    console.error('清除测试数据失败:', error)
  }
  
  // 重新初始化默认数据
  initTestData()
  message.success('测试数据已重置')
}

// 应用测试数据
const applyTestData = () => {
  // 保存到本地存储
  saveTestDataToStorage(testData.value)
  
  // 发送给父组件
  emit('applyTestData', testData.value)
  message.success('测试数据已应用并保存')
}

// 监听变量变化，重新初始化测试数据
watch(() => props.variables, () => {
  initTestData()
}, { deep: true })

// 监听测试数据变化，自动保存
watch(() => testData.value, (newData) => {
  // 延迟保存，避免频繁写入
  setTimeout(() => {
    saveTestDataToStorage(newData)
  }, 1000)
}, { deep: true })

// 组件挂载时初始化测试数据
onMounted(() => {
  initTestData()
})

// 新增变量
function addVariable() {
  editForm.value = {
    id: null,
    name: '',
    description: '',
    variableType: 'text',
    defaultValue: '',
    isRequired: true,
    validationRegex: ''
  }
  showEditModal.value = true
}

// 编辑变量
function editVariable(variable) {
  editForm.value = { 
    id: variable.id,
    name: variable.name || '',
    description: variable.description || '',
    variableType: variable.variableType || 'text',
    defaultValue: variable.defaultValue || '',
    isRequired: variable.isRequired === 1,
    validationRegex: variable.validationRegex || ''
  }
  showEditModal.value = true
}

// 删除变量
function deleteVariable(id) {
  // 使用简单的确认对话框
  if (window.confirm('确定要删除这个变量吗？')) {
    emit('delete', id)
  }
}

// 保存变量
async function saveVariable() {
  try {
    // 使用Naive UI的表单验证
    await formRef.value?.validate()
    
    // 准备保存的数据，只包含必要字段
    const saveData = {
      id: editForm.value.id,
      name: editForm.value.name.trim(),
      description: editForm.value.description.trim(),
      variableType: editForm.value.variableType,
      defaultValue: editForm.value.defaultValue.trim() || '',
      isRequired: editForm.value.isRequired,
      validationRegex: editForm.value.validationRegex.trim() || ''
    }
    
    if (editForm.value.id) {
      emit('edit', saveData)
    } else {
      emit('add', saveData)
    }
    showEditModal.value = false
    // 重置表单数据
    editForm.value = {
      id: null,
      name: '',
      description: '',
      variableType: 'text',
      defaultValue: '',
      isRequired: true,
      validationRegex: ''
    }
    // 重置表单
    formRef.value?.restoreValidation()
  } catch (error) {
    // 表单验证失败，错误信息会自动显示
    console.error('表单验证失败:', error)
  }
}

// 插入变量
function insertVariable(variable) {
  const template = `{{${variable.name}}}`
  emit('insert', template)
}

// 暴露方法给父组件
defineExpose({
  addVariable
})

// 模态框显示状态变化处理
function onModalShowChange(show) {
  if (!show) {
    // 模态框关闭时重置表单数据
    editForm.value = {
      id: null,
      name: '',
      description: '',
      variableType: 'text',
      defaultValue: '',
      isRequired: true,
      validationRegex: ''
    }
    // 重置表单
    formRef.value?.restoreValidation()
  }
}

function cancelEdit() {
  showEditModal.value = false
  // 重置表单数据
  editForm.value = {
    id: null,
    name: '',
    description: '',
    variableType: 'text',
    defaultValue: '',
    isRequired: true,
    validationRegex: ''
  }
  // 清除验证状态
  formRef.value?.restoreValidation()
}
</script>

<style scoped>
.variable-manager {
  width: 100%;
  height: 100%;
  background: #fff;
  display: flex;
  flex-direction: column;
  border-radius: 6px;
  overflow: hidden;
}

.variable-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.variable-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.variable-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-title {
  font-size: 1.1rem;
  font-weight: bold;
  color: #333;
}

.header-count {
  font-size: 14px;
  color: #666;
  background: #f5f5f5;
  padding: 4px 8px;
  border-radius: 4px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.variable-section {
  margin-bottom: 32px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}

.variable-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.variable-card {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 16px;
  background: #fafafa;
  transition: all 0.2s;
}

.variable-card:hover {
  border-color: #18a058;
  background: #f0f9ff;
  box-shadow: 0 2px 8px rgba(24, 160, 88, 0.1);
}

.variable-header-card {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 8px;
}

.variable-name {
  font-weight: bold;
  color: #333;
  font-size: 14px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.variable-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

.variable-desc {
  color: #666;
  font-size: 13px;
  margin-bottom: 12px;
  line-height: 1.4;
}

.variable-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.meta-text {
  font-size: 12px;
  color: #999;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 8px;
  color: #666;
}

.empty-desc {
  font-size: 14px;
  color: #999;
}

/* 滚动条样式 */
.variable-list::-webkit-scrollbar {
  width: 6px;
}

.variable-list::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.variable-list::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.variable-list::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* 测试数据样式 */
.variable-tabs {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.variable-tabs .n-tabs {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.variable-tabs .n-tabs-content {
  flex: 1;
  overflow: hidden;
}

.variable-tabs .n-tab-pane {
  height: 100%;
  overflow: hidden;
}

.test-data-content {
  height: 100%;
  overflow-y: auto;
  padding: 20px;
}

.test-data-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.test-data-header .header-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.test-data-header .header-title {
  font-size: 1.1rem;
  font-weight: bold;
  color: #333;
}

.test-data-header .header-desc {
  font-size: 13px;
  color: #666;
}

.test-data-section {
  margin-bottom: 32px;
}

.test-data-section .section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}

.test-data-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.test-data-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.test-data-item label {
  font-size: 13px;
  font-weight: 500;
  color: #333;
}

.test-data-item .n-input {
  width: 100%;
}
</style> 