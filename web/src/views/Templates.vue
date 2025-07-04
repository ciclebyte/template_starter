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
              <img :src="template.logo" :alt="template.name" />
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

    <!-- 添加模板弹窗 -->
    <n-modal v-model:show="showAddModal" preset="dialog" title="添加模板" :mask-closable="false">
      <n-form :model="addForm" :rules="addRules" ref="addFormRef" label-width="80">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="addForm.name" placeholder="请输入模板名称" />
        </n-form-item>
        <n-form-item label="描述" path="description">
          <n-input v-model:value="addForm.description" placeholder="请输入模板描述" />
        </n-form-item>
        <n-form-item label="分类" path="category_id">
          <n-select v-model:value="addForm.category_id" :options="categorySelectOptions" placeholder="请选择分类" />
        </n-form-item>
        <n-form-item label="标签" path="tags">
          <n-select v-model:value="addForm.tags" :options="tagSelectOptions" multiple placeholder="请选择标签" />
        </n-form-item>
        <n-form-item label="Logo">
          <n-input v-model:value="addForm.logo" placeholder="Logo图片URL，可选" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showAddModal = false">取消</n-button>
        <n-button type="primary" @click="handleAddTemplate">确定</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, h, watch } from 'vue'
import { useRouter } from 'vue-router'
import { NIcon } from 'naive-ui'
import { useLanguageStore } from '@/stores/languageStore'
import { storeToRefs } from 'pinia'
import { useCategoryStore } from '@/stores/categoryStore'

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
const allTemplates = ref([
  {
    id: 1,
    name: 'Vue全栈应用',
    description: '基于Vue3 + Node.js的完整全栈应用模板',
    category_id: 1,
    tags: ['vue', 'nodejs', 'typescript'],
    logo: '/vite.svg',
    languages: [
      { id: 1, name: 'JavaScript', display_name: 'JS', color: '#f7df1e' },
      { id: 2, name: 'Vue', display_name: 'Vue', color: '#42b883' }
    ]
  },
  {
    id: 2,
    name: 'React管理后台',
    description: '企业级React管理后台模板',
    category_id: 1,
    tags: ['react', 'typescript'],
    logo: '/vite.svg',
    languages: [
      { id: 1, name: 'JavaScript', display_name: 'JS', color: '#f7df1e' },
      { id: 3, name: 'React', display_name: 'React', color: '#61dafb' }
    ]
  },
  {
    id: 3,
    name: 'SpringBoot微服务',
    description: '基于SpringBoot的微服务架构模板',
    category_id: 5,
    tags: ['java', 'springboot'],
    logo: '/vite.svg',
    languages: [
      { id: 4, name: 'Java', display_name: 'Java', color: '#007396' }
    ]
  },
  {
    id: 4,
    name: 'Gin Web框架',
    description: '基于Gin的Web应用模板',
    category_id: 6,
    tags: ['go', 'gin'],
    logo: '/vite.svg',
    languages: [
      { id: 5, name: 'Go', display_name: 'Go', color: '#00ADD8' }
    ]
  }
])

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
  tags: [],
  logo: ''
})
const addRules = {
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  description: [{ required: true, message: '请输入模板描述', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }],
}
const categorySelectOptions = categoryTags.value.filter(c => c.id !== 'all').map(c => ({ label: c.name, value: c.id }))
const tagSelectOptions = tags.value.filter(t => t.id !== 'all').map(t => ({ label: t.name, value: t.id }))

const handleAddTemplate = async () => {
  await addFormRef.value?.validate()
  // 简单模拟添加
  allTemplates.value.unshift({
    id: Date.now(),
    name: addForm.value.name,
    description: addForm.value.description,
    category_id: addForm.value.category_id,
    tags: addForm.value.tags,
    logo: addForm.value.logo || '/vite.svg',
    languages: []
  })
  showAddModal.value = false
  addForm.value = { name: '', description: '', category_id: null, tags: [], logo: '' }
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
  { label: '编辑', key: 'edit', icon: () => h('span', { style: 'color:#18a058' }, '✏️') }
]
const handleDropdownSelect = (key, template) => {
  if (key === 'edit') {
    router.push(`/templates/edit/${template.id}`)
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