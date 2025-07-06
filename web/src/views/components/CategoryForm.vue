<template>
  <n-modal
    v-model:show="showModal"
    preset="card"
    :title="title"
    :style="{ width: '500px' }"
    :mask-closable="false"
    :closable="true"
    @close="handleCancel"
  >
    <n-form
      ref="formRef"
      :model="localForm"
      :rules="rules"
      label-placement="left"
      label-width="80px"
      require-mark-placement="right-hanging"
    >
      <n-form-item label="分类名称" path="name">
        <n-input
          v-model:value="localForm.name"
          placeholder="请输入分类名称"
          maxlength="50"
          show-count
        />
      </n-form-item>
      
      <n-form-item label="排序" path="sort">
        <n-input-number
          v-model:value="localForm.sort"
          placeholder="请输入排序值"
          :min="0"
          :max="9999"
          style="width: 100%"
        />
      </n-form-item>
      
      <n-form-item label="分类描述" path="description">
        <n-input
          v-model:value="localForm.description"
          type="textarea"
          placeholder="请输入分类描述"
          :rows="3"
          maxlength="200"
          show-count
        />
      </n-form-item>
    </n-form>

    <template #footer>
      <div class="modal-footer">
        <n-button @click="handleCancel">取消</n-button>
        <n-button type="primary" :loading="loading" @click="handleSubmit">
          确定
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup>
import { ref, watch } from 'vue'
import { addCategory, editCategory } from '@/api/categories'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: '添加分类'
  },
  form: {
    type: Object,
    default: () => ({
      name: '',
      description: '',
      sort: 0
    })
  },
  isEdit: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:show', 'submit', 'cancel'])

const formRef = ref(null)
const loading = ref(false)

const showModal = ref(false)

// 创建本地表单数据的响应式副本
const localForm = ref({
  name: '',
  description: '',
  sort: 0
})

// 监听props.form的变化，更新本地表单数据
watch(() => props.form, (newForm) => {
  localForm.value = {
    name: newForm.name || '',
    description: newForm.description || '',
    sort: newForm.sort || 0
  }
}, { immediate: true, deep: true })

const rules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 1, max: 50, message: '分类名称长度在 1 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 200, message: '分类描述不能超过 200 个字符', trigger: 'blur' }
  ],
  sort: [
    { type: 'number', min: 0, max: 9999, message: '排序值必须在 0-9999 之间', trigger: 'blur' }
  ]
}

watch(() => props.show, (val) => {
  showModal.value = val
})

watch(showModal, (val) => {
  emit('update:show', val)
})

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    loading.value = true
    
    const submitData = {
      ...localForm.value,
      id: props.form.id
    }
    
    if (props.isEdit) {
      await editCategory(submitData)
    } else {
      await addCategory(submitData)
    }
    
    emit('submit')
    showModal.value = false
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  emit('cancel')
  showModal.value = false
}
</script>

<style scoped>
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style> 