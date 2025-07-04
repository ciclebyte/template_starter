<template>
  <div class="templates-page">
    <!-- 分类区域 -->
    <div class="category-tags-section">
      <div class="container">
        <div class="category-tags-title">
          <span class="category-tags-icon">🏷️</span>
          <span class="category-tags-label">分类</span>
        </div>
        <div class="category-tags-list">
          <n-tag
            v-for="cat in categoryTags"
            :key="cat.id"
            :type="selectedCategory === cat.id ? 'primary' : 'default'"
            size="large"
            class="category-tag-item"
            @click="selectCategory(cat.id)"
          >
            {{ cat.name }}
          </n-tag>
        </div>
      </div>
    </div>

    <!-- Tag列表 -->
    <div class="tags-section">
      <div class="container">
        <div class="tags-list">
          <n-tag 
            v-for="tag in tags" 
            :key="tag.id"
            :type="selectedTag === tag.id ? 'primary' : 'default'"
            size="large"
            class="tag-item"
            @click="selectTag(tag.id)"
          >
            {{ tag.name }}
          </n-tag>
        </div>
      </div>
    </div>

    <!-- 模板列表 -->
    <div class="templates-section">
      <div class="container">
        <div class="templates-header">
          <h2>模板列表</h2>
          <div class="templates-header-actions">
            <div class="templates-count">共 {{ totalTemplates }} 个模板</div>
            <n-button type="primary" @click="showAddModal = true">
              <template #icon>
                <n-icon><svg viewBox="0 0 24 24" width="18" height="18"><path fill="currentColor" d="M19 13H13V19H11V13H5V11H11V5H13V11H19V13Z"/></svg></n-icon>
              </template>
              添加模板
            </n-button>
          </div>
        </div>
        
        <div class="templates-grid">
          <div
            v-for="template in templates"
            :key="template.id"
            class="template-card"
            @contextmenu.prevent.stop="showDropdown(template, $event)"
          >
            <div class="template-logo">
              <img :src="template.logo || DEFAULT_LOGO" :alt="template.name" />
            </div>
            <div class="template-info">
              <h3>{{ template.name }}</h3>
              <p>{{ template.description }}</p>
              <div class="template-languages">
                <n-tag 
                  v-for="lang in template.languages" 
                  :key="lang.id"
                  :color="{ color: lang.color }"
                  size="small"
                >
                  {{ lang.display_name }}
                </n-tag>
              </div>
              <div class="template-actions">
                <n-button type="primary" @click="useTemplate(template)">
                  使用模板
                </n-button>
                <n-button @click="previewTemplate(template)">
                  预览
                </n-button>
              </div>
            </div>
          </div>
          <n-dropdown
            v-if="dropdownShow"
            :options="dropdownOptions"
            trigger="manual"
            :show="dropdownShow"
            :x="dropdownX"
            :y="dropdownY"
            @select="key => handleDropdownSelect(key, dropdownTemplate)"
            @clickoutside="dropdownShow = false"
          />
        </div>

        <!-- 分页 -->
        <div class="pagination-section">
          <n-pagination
            v-model:page="currentPage"
            :page-count="totalPages"
            :page-sizes="[20, 40, 60]"
            :page-size="pageSize"
            show-size-picker
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </div>
    </div>

    <!-- 新增弹窗 -->
    <TemplateForm
      :show="showAddModal"
      title="添加模板"
      :form="addForm"
      :rules="addRules"
      :categorySelectOptions="categorySelectOptions"
      @update:show="val => showAddModal = val"
      @submit="handleAddTemplate"
      @cancel="showAddModal = false"
    />

    <!-- 编辑弹窗 -->
    <TemplateForm
      :show="showEditModal"
      title="编辑模板"
      :form="editForm"
      :rules="addRules"
      :categorySelectOptions="categorySelectOptions"
      @update:show="val => showEditModal = val"
      @submit="handleEditTemplate"
      @cancel="showEditModal = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, h, watch } from 'vue'
import { useRouter } from 'vue-router'
import { NIcon } from 'naive-ui'
import { useLanguageStore } from '@/stores/languageStore'
import { storeToRefs } from 'pinia'
import { useCategoryStore } from '@/stores/categoryStore'
import { addTemplate, listTemplates, editTemplate } from '@/api/templates'
import { addTemplateLanguage } from '@/api/templateLanguages'
import TemplateForm from './components/TemplateForm.vue'

const router = useRouter()

const languageStore = useLanguageStore()
const { languagesList } = storeToRefs(languageStore)

const categoryStore = useCategoryStore()
const { categoriesList } = storeToRefs(categoryStore)

// 分类tag数据
const categoryTags = computed(() => [
  { id: 'all', name: '全部' },
  ...categoriesList.value.map(cat => ({ id: cat.id, name: cat.name }))
])
const selectedCategory = ref('all')

// 状态管理
const selectedTag = ref('all')
const currentPage = ref(1)
const pageSize = ref(20)

// 模拟数据 - 标签
const tags = computed(() => {
  console.log('languagesList for tags:', languagesList.value)
  const arr = [
    { id: 'all', name: '全部' },
    ...languagesList.value.map(lang => ({ id: lang.id, name: lang.name }))
  ]
  console.log('computed tags:', arr)
  return arr
})

// 模拟数据 - 模板
const allTemplates = ref([])

// 计算属性
const filteredTemplates = computed(() => {
  let filtered = allTemplates.value

  // 按分类筛选
  if (selectedCategory.value !== 'all') {
    filtered = filtered.filter(t => t.category_id === Number(selectedCategory.value))
  }

  // 按标签筛选
  if (selectedTag.value !== 'all') {
    filtered = filtered.filter(t => t.tags.includes(selectedTag.value))
  }

  return filtered
})

const totalTemplates = computed(() => filteredTemplates.value.length)

const totalPages = computed(() => Math.ceil(totalTemplates.value / pageSize.value))

const templates = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredTemplates.value.slice(start, end)
})

// 添加模板弹窗相关
const showAddModal = ref(false)
const addFormRef = ref(null)
const addForm = ref({
  name: '',
  description: '',
  category_id: null,
  languages: [],
  primary_language: null,
  logo: ''
})
const addRules = {
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  description: [{ required: true, message: '请输入模板描述', trigger: 'blur' }],
  category_id: [{ required: true, type: 'number', message: '请选择分类', trigger: 'change' }],
  languages: [{ required: true, type: 'array', min: 1, message: '请选择支持的语言', trigger: 'change' }],
  primary_language: [{ required: true, type: 'number', message: '请选择主语言', trigger: 'change' }]
}
const categorySelectOptions = computed(() =>
  categoriesList.value.map(c => ({ label: c.name, value: Number(c.id) }))
)
const languageSelectOptions = computed(() =>
  languagesList.value.map(lang => ({ label: lang.name, value: Number(lang.id) }))
)
const primaryLanguageOptions = computed(() =>
  languagesList.value
    .filter(lang => editForm.value.languages.includes(Number(lang.id)))
    .map(lang => ({
      label: lang.name,
      value: Number(lang.id)
    }))
)
const onLanguagesChange = (val) => {
  if (!val.includes(addForm.value.primary_language)) {
    addForm.value.primary_language = null
  }
}

// 工具方法：获取语言名称
const getLanguageName = (id) => {
  const lang = languagesList.value.find(l => Number(l.id) === Number(id))
  return lang ? lang.name : id
}

const DEFAULT_LOGO = '/vite.svg'

const handleAddTemplate = async () => {
  // 类型转换
  addForm.value.category_id = Number(addForm.value.category_id)
  addForm.value.primary_language = Number(addForm.value.primary_language)
  addForm.value.languages = addForm.value.languages.map(Number)
  await addFormRef.value?.validate()
  // 组装languages结构体
  const languagesArr = addForm.value.languages.map(langId => ({
    languageId: langId,
    isPrimary: langId === addForm.value.primary_language ? 1 : 0
  }))
  // 1. 添加模板（一次性提交所有信息）
  await addTemplate({
    name: addForm.value.name,
    description: addForm.value.description,
    categoryId: addForm.value.category_id,
    isFeatured: 0,
    logo: addForm.value.logo || DEFAULT_LOGO,
    languages: languagesArr
  })
  showAddModal.value = false
  addForm.value = { name: '', description: '', category_id: null, languages: [], primary_language: null, logo: '' }
  // TODO: 刷新模板列表
}

// 方法
const selectCategory = (catId) => {
  selectedCategory.value = catId === 'all' ? 'all' : Number(catId)
  currentPage.value = 1
}

const selectTag = (tagId) => {
  selectedTag.value = tagId
  currentPage.value = 1
}

const handlePageChange = (page) => {
  currentPage.value = page
}

const handlePageSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
}

const useTemplate = (template) => {
  console.log('使用模板:', template.name)
  // TODO: 跳转到模板使用页面
}

const previewTemplate = (template) => {
  console.log('预览模板:', template.name)
  // TODO: 跳转到模板预览页面
}

const dropdownShow = ref(false)
const dropdownTemplate = ref(null)
const dropdownX = ref(0)
const dropdownY = ref(0)
const showDropdown = (template, e) => {
  e.preventDefault()
  dropdownShow.value = true
  dropdownTemplate.value = template
  dropdownX.value = e.clientX
  dropdownY.value = e.clientY
}
const dropdownOptions = [
  { label: '编辑模板信息', key: 'editInfo', icon: () => h('span', { style: 'color:#18a058' }, '✏️') },
  { label: '编辑模板内容', key: 'editContent', icon: () => h('span', { style: 'color:#2080f0' }, '📄') }
]
const showEditModal = ref(false)
const editFormRef = ref(null)
const editForm = ref({
  id: null,
  name: '',
  description: '',
  category_id: null,
  languages: [],
  primary_language: null,
  logo: ''
})

const openEditModal = (template) => {
  editForm.value.id = template.id
  editForm.value.name = template.name
  editForm.value.description = template.description
  editForm.value.category_id = template.categoryId || template.category_id
  // 语言回显
  if (template.languages && template.languages.length > 0) {
    editForm.value.languages = template.languages.map(l => Number(l.languageId || l.id))
    // 赋值主语言
    const primary = template.languages.find(l => l.isPrimary === 1 || l.is_primary === 1)
    editForm.value.primary_language = primary ? Number(primary.languageId || primary.id) : null
  } else {
    editForm.value.languages = []
    editForm.value.primary_language = null
  }
  editForm.value.logo = template.logo
  showEditModal.value = true
}

watch(
  () => editForm.value.languages,
  (langs) => {
    if (!langs.includes(editForm.value.primary_language)) {
      editForm.value.primary_language = null
    }
  }
)

const handleEditTemplate = async () => {
  // 类型转换
  editForm.value.category_id = Number(editForm.value.category_id)
  editForm.value.primary_language = Number(editForm.value.primary_language)
  editForm.value.languages = editForm.value.languages.map(Number)
  await editFormRef.value?.validate()
  const languagesArr = editForm.value.languages.map(langId => ({
    languageId: langId,
    isPrimary: langId === editForm.value.primary_language ? 1 : 0
  }))
  await editTemplate({
    id: editForm.value.id,
    name: editForm.value.name,
    description: editForm.value.description,
    categoryId: editForm.value.category_id,
    isFeatured: 0,
    logo: editForm.value.logo || DEFAULT_LOGO,
    languages: languagesArr
  })
  showEditModal.value = false
  // TODO: 刷新模板列表
}

const handleDropdownSelect = (key, template) => {
  if (key === 'editInfo') {
    openEditModal(template)
  } else if (key === 'editContent') {
    router.push(`/templates/content/${template.id}`)
  }
  dropdownShow.value = false
}

watch(languagesList, (val) => {
  console.log('languagesList changed:', val)
})

// 初始化
onMounted(async () => {
  await languageStore.getLanguages()
  await categoryStore.getCategories()
  // 获取真实模板数据
  const res = await listTemplates({})
  allTemplates.value = res.data.data.templatesList || []
  console.log('onMounted languagesList:', languagesList.value)
  // 可以在这里加载数据
})
</script>

<style scoped>
.templates-page {
  min-height: 100vh;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

/* 分类tag区域样式 */
.category-tags-section {
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
  padding: 24px 0 8px 0;
}
.category-tags-title {
  display: flex;
  align-items: center;
  font-size: 1.2rem;
  font-weight: bold;
  color: #18a058;
  margin-bottom: 16px;
}
.category-tags-icon {
  font-size: 1.5rem;
  margin-right: 8px;
}
.category-tags-label {
  font-size: 1.1rem;
}
.category-tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}
.category-tag-item {
  cursor: pointer;
  transition: all 0.2s;
}
.category-tag-item:hover {
  transform: translateY(-1px);
}

/* 标签列表样式 */
.tags-section {
  background: #f8f9fa;
  padding: 20px 0;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.tag-item {
  cursor: pointer;
  transition: all 0.3s ease;
}

.tag-item:hover {
  transform: translateY(-1px);
}

/* 模板列表样式 */
.templates-section {
  padding: 40px 0;
  background: #fff;
}

.templates-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}
.templates-header-actions {
  display: flex;
  align-items: center;
  gap: 18px;
}
.templates-header h2 {
  margin: 0;
  color: #333;
}

.templates-count {
  color: #666;
  font-size: 14px;
}

.templates-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 30px;
  margin-bottom: 40px;
}

.template-card {
  background: white;
  border-radius: 12px;
  padding: 25px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  transition: all 0.3s ease;
  border: 1px solid #f0f0f0;
}

.template-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0,0,0,0.15);
}

.template-logo {
  text-align: center;
  margin-bottom: 20px;
}

.template-logo img {
  width: 60px;
  height: 60px;
}

.template-info h3 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 18px;
}

.template-info p {
  color: #666;
  margin-bottom: 15px;
  line-height: 1.5;
}

.template-languages {
  margin-bottom: 20px;
}

.template-languages .n-tag {
  margin-right: 8px;
  margin-bottom: 8px;
}

.template-actions {
  display: flex;
  gap: 10px;
}

/* 分页样式 */
.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .templates-grid {
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 20px;
  }
  
  .template-actions {
    flex-direction: column;
  }
}
</style> 