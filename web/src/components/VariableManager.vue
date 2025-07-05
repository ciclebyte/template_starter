<template>
  <div class="variable-manager">
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
import { ref, computed } from 'vue'
import { NButton, NIcon, NModal, NForm, NFormItem, NInput, NSelect, NSwitch, NSpace, NTag, useMessage } from 'naive-ui'
import { Add, DocumentText, CodeSlash } from '@vicons/ionicons5'

const props = defineProps({
  variables: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['add', 'edit', 'delete', 'insert'])

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
</style> 