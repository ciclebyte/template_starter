<template>
  <div class="variable-manager">
    <div class="variable-tabs">
      <n-tabs type="line" animated>
        <n-tab-pane name="manage" tab="变量管理">
          <div class="variable-list">
            <div v-for="variable in textVariables" :key="variable.id" class="variable-item">
              <div class="variable-info">
                <div class="variable-name">{{ variable.name }}</div>
                <div class="variable-desc">{{ variable.description }}</div>
              </div>
              <div class="variable-actions">
                <n-button size="small" @click="insertVariable(variable)">插入</n-button>
                <n-button size="small" @click="editVariable(variable)">编辑</n-button>
                <n-button size="small" type="error" @click="deleteVariable(variable.id)">删除</n-button>
              </div>
            </div>
            <div v-for="variable in conditionalVariables" :key="variable.id" class="variable-item">
              <div class="variable-info">
                <div class="variable-name">{{ variable.name }}</div>
                <div class="variable-desc">{{ variable.description }}</div>
              </div>
              <div class="variable-actions">
                <n-button size="small" @click="insertVariable(variable)">插入</n-button>
                <n-button size="small" @click="editVariable(variable)">编辑</n-button>
                <n-button size="small" type="error" @click="deleteVariable(variable.id)">删除</n-button>
              </div>
            </div>
            <div v-if="templateVariables.length === 0" class="empty-state">
              暂无变量
            </div>
          </div>
        </n-tab-pane>
        
        <n-tab-pane name="quick" tab="快速插入">
          <div class="quick-insert-content">
            <div class="quick-grid">
              <n-button 
                v-for="quickVar in quickVariables" 
                :key="quickVar.name"
                size="small" 
                class="quick-button"
                @click="insertQuickVariable(quickVar)"
              >
                {{ quickVar.label }}
              </n-button>
            </div>
          </div>
        </n-tab-pane>
      </n-tabs>
    </div>

    <!-- 变量编辑对话框 -->
    <n-modal v-model:show="showEditModal" preset="card" title="编辑变量" style="width: 600px">
      <n-form ref="formRef" :model="editForm" :rules="rules" label-placement="left" label-width="100px">
        <n-form-item label="变量名称" path="name">
          <n-input v-model:value="editForm.name" placeholder="请输入变量名称" />
        </n-form-item>
        <n-form-item label="变量描述" path="description">
          <n-input v-model:value="editForm.description" placeholder="请输入变量描述" />
        </n-form-item>
        <n-form-item label="变量类型" path="type">
          <n-select v-model:value="editForm.type" :options="typeOptions" />
        </n-form-item>
        <n-form-item label="默认值" path="defaultValue">
          <n-input v-model:value="editForm.defaultValue" placeholder="请输入默认值" />
        </n-form-item>
        <n-form-item label="是否必填" path="isRequired">
          <n-switch v-model:value="editForm.isRequired" />
        </n-form-item>
        <n-form-item label="验证规则" path="validationRegex">
          <n-input v-model:value="editForm.validationRegex" placeholder="请输入验证正则表达式" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showEditModal = false">取消</n-button>
          <n-button type="primary" @click="saveVariable">保存</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { NButton, NIcon, NTabs, NTabPane, NModal, NForm, NFormItem, NInput, NSelect, NSwitch, NSpace, useMessage } from 'naive-ui'

const props = defineProps({
  variables: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['add', 'edit', 'delete', 'insert'])

const message = useMessage()

// 变量类型选项
const typeOptions = [
  { label: '文本变量', value: 'text' },
  { label: '条件变量', value: 'conditional' }
]

// 快速插入变量
const quickVariables = [
  { name: 'project_name', label: '项目名', template: '{{project_name}}' },
  { name: 'author', label: '作者', template: '{{author}}' },
  { name: 'version', label: '版本', template: '{{version}}' },
  { name: 'description', label: '描述', template: '{{description}}' },
  { name: 'date', label: '日期', template: '{{date}}' },
  { name: 'time', label: '时间', template: '{{time}}' }
]

// 计算属性：按类型分组变量
const textVariables = computed(() => {
  return props.variables.filter(v => v.type === 'text' || !v.type)
})

const conditionalVariables = computed(() => {
  return props.variables.filter(v => v.type === 'conditional')
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
  type: 'text',
  defaultValue: '',
  isRequired: true,
  validationRegex: ''
})

const rules = {
  name: { required: true, message: '请输入变量名称' },
  description: { required: true, message: '请输入变量描述' },
  type: { required: true, message: '请选择变量类型' }
}

// 新增变量
function addVariable() {
  editForm.value = {
    id: null,
    name: '',
    description: '',
    type: 'text',
    defaultValue: '',
    isRequired: true,
    validationRegex: ''
  }
  showEditModal.value = true
}

// 编辑变量
function editVariable(variable) {
  editForm.value = { ...variable }
  showEditModal.value = true
}

// 删除变量
function deleteVariable(id) {
  emit('delete', id)
}

// 保存变量
function saveVariable() {
  if (editForm.value.id) {
    emit('edit', editForm.value)
  } else {
    emit('add', editForm.value)
  }
  showEditModal.value = false
}

// 插入变量
function insertVariable(variable) {
  const template = `{{${variable.name}}}`
  emit('insert', template)
}

// 插入快速变量
function insertQuickVariable(quickVar) {
  emit('insert', quickVar.template)
}

// 暴露方法给父组件
defineExpose({
  addVariable
})
</script>

<style scoped>
.variable-manager {
  width: 100%;
  background: #fff;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.variable-tabs {
  flex: 1;
  overflow: hidden;
}

.variable-list {
  padding: 16px;
  overflow-y: auto;
  height: 100%;
}

.variable-item {
  padding: 12px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  margin-bottom: 8px;
  background: #fafafa;
}

.variable-info {
  margin-bottom: 8px;
}

.variable-name {
  font-weight: bold;
  color: #333;
  margin-bottom: 4px;
}

.variable-desc {
  font-size: 12px;
  color: #666;
}

.variable-actions {
  display: flex;
  gap: 8px;
}

.empty-state {
  text-align: center;
  color: #999;
  padding: 32px;
}

.quick-insert-content {
  padding: 16px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  width: 100%;
  max-width: 300px;
}

.quick-button {
  font-size: 12px;
  height: 36px;
}
</style> 