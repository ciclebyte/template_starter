<template>
  <n-modal
    :show="show"
    @update:show="val => emit('update:show', val)"
    :title="title"
    preset="dialog"
    :mask-closable="false"
    style="width: 90vw; max-width: 1200px;"
  >
    <div class="tabs-container">
      <n-tabs v-model:value="activeTab" type="line" animated>
        <!-- 基础信息 Tab -->
        <n-tab-pane name="basic" tab="基础信息">
          <div class="tab-content">
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
          </div>
        </n-tab-pane>

        <!-- 详细介绍 Tab -->
        <n-tab-pane name="introduction" tab="详细介绍">
          <div class="tab-content">
            <div class="introduction-tab">
              <div class="introduction-header">
                <h3>模板详细介绍</h3>
                <p class="introduction-tip">支持 Markdown 格式，可以添加代码示例、使用说明、特性介绍等</p>
              </div>
                          <div class="markdown-editor-container">
              <MdEditor
                v-model="form.introduction"
                :toolbars="toolbars"
                :preview="preview"
                :codeTheme="codeTheme"
                :theme="theme"
                placeholder="请输入模板详细介绍，支持 Markdown 格式..."
                style="width: 100%;"
              />
            </div>
            </div>
          </div>
        </n-tab-pane>
      </n-tabs>
    </div>

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
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

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
const activeTab = ref('basic')

// MdEditor 配置
const toolbars = [
  'bold', 'underline', 'italic', 'strikethrough', 'title', 'sub', 'sup', 'quote', 'unordered-list', 'ordered-list',
  'task-list', '-', 'code', 'code-block', 'link', 'image', 'table', 'mermaid', 'katex', '-',
  'preview', 'fullscreen', 'preview-only'
]
const preview = 'live'
const codeTheme = 'github'
const theme = 'light'

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

<style scoped>
.tabs-container {
  min-height: 700px;
}

.tab-content {
  min-height: 650px;
  padding: 20px 0;
}

.introduction-tab {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.introduction-header {
  margin-bottom: 20px;
  flex-shrink: 0;
}

.introduction-header h3 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 18px;
  font-weight: bold;
}

.introduction-tip {
  margin: 0;
  color: #666;
  font-size: 14px;
  line-height: 1.5;
}

.markdown-editor-container {
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  overflow: hidden;
  width: 100%;
  flex: 1;
  min-height: 550px;
  height: 550px;
}

.markdown-editor-container :deep(.md-editor) {
  border: none;
  width: 100% !important;
  height: 100% !important;
}

.markdown-editor-container :deep(.md-editor-toolbar) {
  border-bottom: 1px solid #f0f0f0;
  background: #fafafa;
  width: 100%;
}

.markdown-editor-container :deep(.md-editor-content) {
  background: #fff;
  width: 100%;
  height: calc(100% - 40px) !important;
}

.markdown-editor-container :deep(.md-editor-input) {
  width: 100% !important;
  min-height: 500px !important;
}

.markdown-editor-container :deep(.md-editor-preview) {
  width: 100% !important;
  min-height: 500px !important;
}

.markdown-editor-container :deep(.md-editor-textarea) {
  min-height: 500px !important;
}

/* Tab 样式优化 */
:deep(.n-tabs-tab) {
  font-size: 16px;
  font-weight: 500;
}

:deep(.n-tabs-tab-pad) {
  padding: 16px 24px;
}

:deep(.n-tab-pane) {
  padding: 0;
}

:deep(.n-tabs-content) {
  height: 100%;
}

:deep(.n-tab-pane) {
  height: 100%;
}
</style> 