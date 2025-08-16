<template>
  <!-- 条件设置弹框 -->
  <n-modal 
    v-model:show="modalVisible" 
    preset="card" 
    style="width: 600px;"
    :mask-closable="false"
  >
    <template #header>
      <div class="condition-modal-header">
        <span class="modal-title">设置生成条件</span>
        <span v-if="selectedFileForCondition" class="file-path">{{ selectedFileForCondition.fileName }}</span>
      </div>
    </template>
    
    <div class="condition-form">
      <n-form ref="conditionFormRef" :model="conditionForm" label-placement="left" label-width="120px">
        <n-form-item label="启用条件" path="enabled">
          <n-switch v-model:value="conditionForm.enabled" />
          <span class="form-help">开启后，文件将根据条件决定是否生成</span>
        </n-form-item>
        
        <template v-if="conditionForm.enabled">
          <n-form-item label="变量名称" path="variableName" :rule="{ required: true, message: '请选择变量' }">
            <n-select 
              v-model:value="conditionForm.variableName" 
              :options="booleanVariableOptions"
              placeholder="请选择一个布尔类型变量"
              filterable
            />
          </n-form-item>
          
          <n-form-item label="期望值" path="expectedValue">
            <n-radio-group v-model:value="conditionForm.expectedValue">
              <n-radio :value="true">为真时生成</n-radio>
              <n-radio :value="false">为假时生成</n-radio>
            </n-radio-group>
          </n-form-item>
          
          <n-form-item label="条件描述" path="description">
            <n-input 
              v-model:value="conditionForm.description" 
              placeholder="可选：描述此条件的用途"
              type="textarea"
              :rows="3"
            />
          </n-form-item>
        </template>
      </n-form>
      
      <div class="modal-actions">
        <n-button @click="$emit('close')">取消</n-button>
        <n-button type="primary" @click="handleSave">保存</n-button>
      </div>
    </div>
  </n-modal>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { NModal, NForm, NFormItem, NSwitch, NSelect, NRadioGroup, NRadio, NInput, NButton } from 'naive-ui'

const props = defineProps({
  show: {
    type: Boolean,
    required: true
  },
  selectedFileForCondition: {
    type: Object,
    default: null
  },
  templateVariables: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:show', 'close', 'save'])

const conditionFormRef = ref(null)
const conditionForm = ref({
  enabled: false,
  variableName: '',
  expectedValue: true,
  description: ''
})

const modalVisible = computed({
  get: () => props.show,
  set: (value) => emit('update:show', value)
})

// 布尔类型变量选项 - 模板变量功能已移除，改为使用变量定义
const booleanVariableOptions = computed(() => {
  // TODO: 集成变量定义的布尔变量
  return []
})

const handleSave = async () => {
  try {
    await conditionFormRef.value?.validate()
    emit('save', conditionForm.value)
  } catch (error) {
    // 验证失败
    console.error('表单验证失败:', error)
  }
}

// 重置表单
const resetForm = () => {
  conditionForm.value = {
    enabled: false,
    variableName: '',
    expectedValue: true,
    description: ''
  }
}

// 设置表单数据
const setFormData = (data) => {
  conditionForm.value = {
    enabled: data.enabled || false,
    variableName: data.variableName || '',
    expectedValue: data.expectedValue !== undefined ? data.expectedValue : true,
    description: data.description || ''
  }
}

// 监听文件选择变化，重置表单
watch(() => props.selectedFileForCondition, () => {
  resetForm()
}, { immediate: true })

// 暴露方法给父组件
defineExpose({
  resetForm,
  setFormData
})
</script>

<style scoped>
.condition-modal-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.condition-modal-header .file-path {
  font-size: 12px;
  color: #666;
  font-weight: normal;
}

.condition-form {
  padding: 16px 0;
}

.form-help {
  margin-left: 8px;
  font-size: 12px;
  color: #666;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e9ecef;
}
</style>