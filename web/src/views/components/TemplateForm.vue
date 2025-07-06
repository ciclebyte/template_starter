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
              <n-form-item label="分类" path="categoryId">
                <n-select v-model:value="form.categoryId" :options="categorySelectOptions" placeholder="请选择分类" />
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
              <n-form-item label="图标">
                <div class="icon-selector-wrapper">
                  <div class="current-icon" @click="showIconSelector = true">
                    <n-icon v-if="form.icon" size="32">
                      <component :is="getIconComponent(form.icon)" />
                    </n-icon>
                    <div v-else class="no-icon">
                      <n-icon size="24" color="#ccc">
                        <svg viewBox="0 0 24 24" fill="currentColor">
                          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                        </svg>
                      </n-icon>
                      <span>选择图标</span>
                    </div>
                  </div>
                  <div class="icon-actions">
                    <n-button 
                      size="small" 
                      type="primary" 
                      ghost
                      @click="showIconSelector = true"
                    >
                      选择图标
                    </n-button>
                    <n-button 
                      v-if="form.icon" 
                      size="small" 
                      type="error" 
                      ghost
                      @click="clearIcon"
                    >
                      清除
                    </n-button>
                  </div>
                </div>
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

  <!-- 图标选择器弹框 -->
  <n-modal
    v-model:show="showIconSelector"
    title="选择图标"
    preset="dialog"
    style="width: 90vw; max-width: 900px;"
  >
    <IconSelector @select="onIconSelect" />
    <template #action>
      <n-button @click="showIconSelector = false">取消</n-button>
      <n-button type="primary" @click="showIconSelector = false">确定</n-button>
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
import IconSelector from '@/components/IconSelector.vue'
import * as IonIcons from '@vicons/ionicons5'

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
const showIconSelector = ref(false)

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

// 图标相关方法
const getIconComponent = (iconName) => {
  return IonIcons[iconName] || null
}

const onIconSelect = (iconName) => {
  props.form.icon = iconName
  showIconSelector.value = false
}

const clearIcon = () => {
  props.form.icon = ''
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
}

.icon-selector-wrapper {
  display: flex;
  align-items: center;
  gap: 16px;
}

.current-icon {
  width: 64px;
  height: 64px;
  border: 2px dashed #d9d9d9;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #fafafa;
  position: relative;
  overflow: hidden;
}

.current-icon:hover {
  border-color: #18a058;
  background-color: #f6ffed;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(24, 160, 88, 0.15);
}

.current-icon:active {
  transform: translateY(0);
}

.no-icon {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  color: #999;
  font-size: 12px;
}

.no-icon span {
  font-weight: 500;
}

.icon-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.icon-actions .n-button {
  border-radius: 6px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.icon-actions .n-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
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