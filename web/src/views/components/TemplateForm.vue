<template>
  <n-modal
    :show="show"
    @update:show="val => emit('update:show', val)"
    :title="title"
    preset="dialog"
    :mask-closable="false"
  >
    <n-form :model="form" :rules="rules" ref="formRef" label-width="80">
      <n-form-item label="名称" path="name">
        <n-input v-model:value="form.name" placeholder="请输入模板名称" />
      </n-form-item>
      <n-form-item label="描述" path="description">
        <n-input v-model:value="form.description" placeholder="请输入模板描述" />
      </n-form-item>
      <n-form-item label="分类" path="category_id">
        <n-select v-model:value="form.category_id" :options="categorySelectOptions" placeholder="请选择分类" />
      </n-form-item>
      <n-form-item label="支持语言" path="languages">
        <n-select
          v-model:value="form.languages"
          :options="languageSelectOptions"
          multiple
          placeholder="请选择支持的语言"
          @update:value="onLanguagesChange"
        />
      </n-form-item>
      <n-form-item label="主语言" path="primary_language">
        <n-select
          v-model:value="form.primary_language"
          :options="primaryLanguageOptions"
          placeholder="请选择主语言"
        />
      </n-form-item>
      <n-form-item label="Logo">
        <n-input v-model:value="form.logo" placeholder="Logo图片URL，可选" />
      </n-form-item>
      <div v-if="form.languages.length > 0" style="margin-bottom: 12px;">
        <n-tag
          v-for="langId in form.languages"
          :key="langId"
          :type="langId === form.primary_language ? 'primary' : 'default'"
          style="margin-right: 8px;"
        >
          {{ getLanguageName(langId) }}<template v-if="langId === form.primary_language">（主语言）</template>
        </n-tag>
      </div>
    </n-form>
    <template #action>
      <n-button @click="$emit('cancel')">取消</n-button>
      <n-button type="primary" @click="handleSubmit">确定</n-button>
    </template>
  </n-modal>
</template>

<script setup>
import { ref, computed, watch, toRefs } from 'vue'
import { storeToRefs } from 'pinia'
import { useLanguageStore } from '@/stores/languageStore'
import { useCategoryStore } from '@/stores/categoryStore'

const props = defineProps({
  show: Boolean,
  title: { type: String, default: '模板表单' },
  form: Object,
  rules: Object,
  categorySelectOptions: Array,
  onSubmit: Function
})
const emit = defineEmits(['update:show', 'submit', 'cancel'])
const formRef = ref(null)

const languageStore = useLanguageStore()
const { languagesList } = storeToRefs(languageStore)

const languageSelectOptions = computed(() =>
  languagesList.value.map(lang => ({ label: lang.name, value: Number(lang.id) }))
)
const primaryLanguageOptions = computed(() =>
  languagesList.value
    .filter(lang => props.form.languages.includes(Number(lang.id)))
    .map(lang => ({ label: lang.name, value: Number(lang.id) }))
)
const onLanguagesChange = (val) => {
  if (!val.includes(props.form.primary_language)) {
    props.form.primary_language = null
  }
}
const getLanguageName = (id) => {
  const lang = languagesList.value.find(l => Number(l.id) === Number(id))
  return lang ? lang.name : id
}

const handleSubmit = async () => {
  await formRef.value?.validate()
  emit('submit')
}
</script> 