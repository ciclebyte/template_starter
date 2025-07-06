<template>
  <div class="templates-page">
    <!-- 分类区域 -->
    <div class="category-tags-section">
      <div class="container">
        <div class="category-tags-title">
          <span class="category-tags-icon">🏷️</span>
          <span class="category-tags-label">分类</span>
          <div class="title-actions">
            <n-button 
              size="small" 
              type="primary" 
              ghost
              class="add-btn"
              @click="showAddCategoryModal = true"
            >
              <template #icon>
                <n-icon><svg viewBox="0 0 24 24" width="16" height="16"><path fill="currentColor" d="M19 13H13V19H11V13H5V11H11V5H13V11H19V13Z"/></svg></n-icon>
              </template>
            </n-button>
          </div>
        </div>
        <div class="category-tags-list">
          <n-tag
            v-for="cat in categoryTags"
            :key="cat.id"
            :type="selectedCategory === cat.id ? 'primary' : 'default'"
            size="large"
            class="category-tag-item"
            @click="selectCategory(cat.id)"
            @contextmenu.prevent.stop="showCategoryDropdown(cat, $event)"
          >
            {{ cat.name }}
          </n-tag>
        </div>
        
        <!-- 分类右键菜单 -->
        <n-dropdown
          v-if="categoryDropdownShow"
          :options="categoryDropdownOptions"
          trigger="manual"
          :show="categoryDropdownShow"
          :x="categoryDropdownX"
          :y="categoryDropdownY"
          @select="key => handleCategoryDropdownSelect(key, categoryDropdownItem)"
          @clickoutside="categoryDropdownShow = false"
        />
      </div>
    </div>

    <!-- Tag列表 -->
    <div class="tags-section">
      <div class="container">
        <div class="tags-title">
          <span class="tags-icon">💻</span>
          <span class="tags-label">语言</span>
          <div class="title-actions">
            <n-button 
              size="small" 
              type="primary" 
              ghost
              class="add-btn"
              @click="showAddLanguageModal = true"
            >
              <template #icon>
                <n-icon><svg viewBox="0 0 24 24" width="16" height="16"><path fill="currentColor" d="M19 13H13V19H11V13H5V11H11V5H13V11H19V13Z"/></svg></n-icon>
              </template>
            </n-button>
          </div>
        </div>
        <div class="tags-list">
          <n-tag 
            v-for="tag in tags" 
            :key="tag.id"
            :type="selectedTag === tag.id ? 'primary' : 'default'"
            size="large"
            class="tag-item"
            @click="selectTag(tag.id)"
            @contextmenu.prevent.stop="showLanguageDropdown(tag, $event)"
          >
            {{ tag.name }}
          </n-tag>
        </div>
        
        <!-- 语言右键菜单 -->
        <n-dropdown
          v-if="languageDropdownShow"
          :options="languageDropdownOptions"
          trigger="manual"
          :show="languageDropdownShow"
          :x="languageDropdownX"
          :y="languageDropdownY"
          @select="key => handleLanguageDropdownSelect(key, languageDropdownItem)"
          @clickoutside="languageDropdownShow = false"
        />
      </div>
    </div>

    <!-- 模板列表 -->
    <div class="templates-section">
      <div class="container">
        <div class="templates-header">
          <div class="templates-title-section">
            <h2>模板列表</h2>
            <div class="templates-count">共 {{ totalTemplates }} 个模板</div>
          </div>
          <div class="templates-header-actions">
            <n-button type="primary" @click="showAddModal = true" class="add-template-btn">
              <template #icon>
                <n-icon><svg viewBox="0 0 24 24" width="18" height="18"><path fill="currentColor" d="M19 13H13V19H11V13H5V11H11V5H13V11H19V13Z"/></svg></n-icon>
              </template>
            </n-button>
          </div>
        </div>
        
        <div class="templates-grid">
          <div
            v-for="template in templates"
            :key="template.id"
            class="template-card"
            :class="{ 'featured': template.isFeatured }"
            @contextmenu.prevent.stop="showDropdown(template, $event)"
          >
            <div class="template-logo">
              <div class="template-icon">
                <n-icon size="48">
                  <component :is="getIconComponent(template.icon) || getDefaultIcon(template)" />
                </n-icon>
              </div>
              <div v-if="template.isFeatured" class="featured-badge">
                <n-icon><svg viewBox="0 0 24 24" width="16" height="16"><path fill="currentColor" d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"/></svg></n-icon>
              </div>
            </div>
            <div class="template-info">
              <h3>{{ template.name }}</h3>
              <p>{{ template.description }}</p>
              <div class="template-tags">
                <!-- 分类标签 -->
                <div class="template-category" v-if="getCategoryName(template.categoryId)">
                  <n-tag 
                    type="success" 
                    size="small"
                    class="category-tag"
                  >
                    {{ getCategoryName(template.categoryId) }}
                  </n-tag>
                </div>
                <!-- 语言标签 -->
                <div class="template-languages">
                  <n-tag 
                    v-for="lang in template.languages" 
                    :key="lang.id"
                    :color="{ color: getLanguageColor(lang.languageId) }"
                    size="small"
                    class="language-tag"
                  >
                    {{ getLanguageName(lang.languageId) }}
                  </n-tag>
                </div>
              </div>
              <div class="template-actions">
                <n-button type="primary" @click="useTemplate(template)">
                  使用模板
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

    <!-- 添加分类弹窗 -->
    <CategoryForm
      :show="showAddCategoryModal"
      title="添加分类"
      :form="categoryForm"
      @update:show="val => showAddCategoryModal = val"
      @submit="handleAddCategory"
      @cancel="handleCancelCategory"
    />

    <!-- 编辑分类弹窗 -->
    <CategoryForm
      :show="showEditCategoryModal"
      title="编辑分类"
      :form="categoryForm"
      :is-edit="true"
      @update:show="val => showEditCategoryModal = val"
      @submit="handleEditCategory"
      @cancel="handleCancelCategory"
    />

    <!-- 添加语言弹窗 -->
    <LanguageForm
      :show="showAddLanguageModal"
      title="添加语言"
      :form="languageForm"
      @update:show="val => showAddLanguageModal = val"
      @submit="handleAddLanguage"
      @cancel="handleCancelLanguage"
    />

    <!-- 编辑语言弹窗 -->
    <LanguageForm
      :show="showEditLanguageModal"
      title="编辑语言"
      :form="languageForm"
      :is-edit="true"
      @update:show="val => showEditLanguageModal = val"
      @submit="handleEditLanguage"
      @cancel="handleCancelLanguage"
    />

    <!-- 删除确认对话框 -->
    <n-modal
      v-model:show="showDeleteModal"
      preset="dialog"
      title="确认删除"
      :show-icon="false"
      :mask-closable="false"
    >
      <template #default>
        <div class="delete-confirm-content">
          <div class="delete-icon">🗑️</div>
          <div class="delete-message">
            确定要删除{{ getDeleteItemType() }} <strong>"{{ deleteTemplateInfo?.name }}"</strong> 吗？
          </div>
          <div class="delete-warning">
            此操作不可恢复，删除后将无法找回。
          </div>
        </div>
      </template>
      <template #action>
        <div class="delete-actions">
          <n-button @click="showDeleteModal = false">取消</n-button>
          <n-button type="error" @click="confirmDelete" :loading="deleteLoading">
            确认删除
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, h, watch } from 'vue'
import { useRouter } from 'vue-router'
import { NIcon, useMessage, NModal, NButton } from 'naive-ui'
import * as IonIcons from '@vicons/ionicons5'
import { useLanguageStore } from '@/stores/languageStore'
import { storeToRefs } from 'pinia'
import { useCategoryStore } from '@/stores/categoryStore'
import { addTemplate, listTemplates, editTemplate, deleteTemplate } from '@/api/templates'
import { deleteCategory } from '@/api/categories'
import { deleteLanguage } from '@/api/languages'
import { addTemplateLanguage } from '@/api/templateLanguages'
import TemplateForm from './components/TemplateForm.vue'
import CategoryForm from './components/CategoryForm.vue'
import LanguageForm from './components/LanguageForm.vue'

const router = useRouter()
const message = useMessage()

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
  const arr = [
    { id: 'all', name: '全部' },
    ...languagesList.value.map(lang => ({ id: lang.id, name: lang.name }))
  ]
  return arr
})

// 模拟数据 - 模板
const allTemplates = ref([])

// 计算属性
const filteredTemplates = computed(() => {
  let filtered = allTemplates.value

  // 按分类筛选
  if (selectedCategory.value !== 'all') {
    filtered = filtered.filter(t => t.categoryId === Number(selectedCategory.value))
  }

  // 按标签筛选
  if (selectedTag.value !== 'all') {
    filtered = filtered.filter(t => {
      // 检查模板的语言列表中是否包含选中的语言
      if (t.languages && Array.isArray(t.languages)) {
        return t.languages.some(lang => {
          const langId = lang.languageId || lang.id
          return Number(langId) === Number(selectedTag.value)
        })
      }
      return false
    })
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
const showAddCategoryModal = ref(false)
const showEditCategoryModal = ref(false)
const showAddLanguageModal = ref(false)
const showEditLanguageModal = ref(false)

// 删除确认对话框相关
const showDeleteModal = ref(false)
const deleteTemplateInfo = ref(null)
const deleteLoading = ref(false)
const categoryForm = ref({
  name: '',
  description: '',
  sort: 0
})

const languageForm = ref({
  name: '',
  code: '',
  display_name: '',
  color: '#18a058',
  sort: 0,
  isPopular: 0
})
const addFormRef = ref(null)
const addForm = ref({
  name: '',
  description: '',
  categoryId: null,
  languages: [],
  primary_language: null,
  icon: '',
  introduction: ''
})
const addRules = {
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  description: [{ required: true, message: '请输入模板描述', trigger: 'blur' }],
  categoryId: [{ required: true, type: 'number', message: '请选择分类', trigger: 'change' }],
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





const handleAddTemplate = async () => {
  // 类型转换
  addForm.value.categoryId = Number(addForm.value.categoryId)
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
    introduction: addForm.value.introduction,
    categoryId: addForm.value.categoryId,
    isFeatured: 0,

    icon: addForm.value.icon,
    languages: languagesArr
  })
  showAddModal.value = false
  addForm.value = { name: '', description: '', categoryId: null, languages: [], primary_language: null, icon: '', introduction: '' }
  // 刷新模板列表
  const res = await listTemplates({})
  allTemplates.value = res.data.data.templatesList || []
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
  router.push(`/templates/generate/${template.id}`)
}

// 获取分类名称
const getCategoryName = (categoryId) => {
  if (!categoryId) return null
  const category = categoriesList.value.find(cat => cat.id === Number(categoryId))
  return category ? category.name : null
}

// 获取语言名称
const getLanguageName = (languageId) => {
  if (!languageId) return ''
  const language = languagesList.value.find(lang => lang.id === Number(languageId))
  return language ? language.name : ''
}

// 获取语言颜色
const getLanguageColor = (languageId) => {
  if (!languageId) return '#666'
  const language = languagesList.value.find(lang => lang.id === Number(languageId))
  return language ? language.color : '#666'
}

// 获取图标组件
const getIconComponent = (iconName) => {
  return IonIcons[iconName] || null
}

// 获取默认图标
const getDefaultIcon = (template) => {
  // 根据模板名称或分类选择合适的默认图标
  const name = template.name?.toLowerCase() || ''
  const categoryId = template.categoryId
  
  // 根据分类选择默认图标
  if (categoryId) {
    const category = categoriesList.value.find(cat => cat.id === Number(categoryId))
    if (category) {
      const categoryName = category.name.toLowerCase()
      if (categoryName.includes('web') || categoryName.includes('前端')) return IonIcons.GlobeOutline
      if (categoryName.includes('mobile') || categoryName.includes('移动')) return IonIcons.PhonePortraitOutline
      if (categoryName.includes('backend') || categoryName.includes('后端')) return IonIcons.ServerOutline
      if (categoryName.includes('fullstack') || categoryName.includes('全栈')) return IonIcons.AppsOutline
      if (categoryName.includes('ui') || categoryName.includes('界面')) return IonIcons.GridOutline
      if (categoryName.includes('data') || categoryName.includes('数据')) return IonIcons.BarChartOutline
    }
  }
  
  // 根据模板名称选择默认图标
  if (name.includes('web') || name.includes('前端') || name.includes('vue') || name.includes('react')) {
    return IonIcons.GlobeOutline
  }
  if (name.includes('mobile') || name.includes('移动') || name.includes('app')) {
    return IonIcons.PhonePortraitOutline
  }
  if (name.includes('backend') || name.includes('后端') || name.includes('api') || name.includes('server')) {
    return IonIcons.ServerOutline
  }
  if (name.includes('fullstack') || name.includes('全栈')) {
    return IonIcons.AppsOutline
  }
  if (name.includes('ui') || name.includes('界面') || name.includes('design')) {
    return IonIcons.GridOutline
  }
  if (name.includes('data') || name.includes('数据') || name.includes('chart')) {
    return IonIcons.BarChartOutline
  }
  if (name.includes('admin') || name.includes('管理')) {
    return IonIcons.SettingsOutline
  }
  if (name.includes('blog') || name.includes('博客')) {
    return IonIcons.DocumentTextOutline
  }
  if (name.includes('ecommerce') || name.includes('电商') || name.includes('shop')) {
    return IonIcons.CartOutline
  }
  
  // 默认图标
  return IonIcons.DocumentOutline
}



const dropdownShow = ref(false)
const dropdownTemplate = ref(null)
const dropdownX = ref(0)
const dropdownY = ref(0)

// 分类右键菜单相关
const categoryDropdownShow = ref(false)
const categoryDropdownItem = ref(null)
const categoryDropdownX = ref(0)
const categoryDropdownY = ref(0)

// 语言右键菜单相关
const languageDropdownShow = ref(false)
const languageDropdownItem = ref(null)
const languageDropdownX = ref(0)
const languageDropdownY = ref(0)
const showDropdown = (template, e) => {
  e.preventDefault()
  dropdownShow.value = true
  dropdownTemplate.value = template
  dropdownX.value = e.clientX
  dropdownY.value = e.clientY
}

const showCategoryDropdown = (category, e) => {
  // 排除"全部"分类
  if (category.id === 'all') return
  
  e.preventDefault()
  categoryDropdownShow.value = true
  categoryDropdownItem.value = category
  categoryDropdownX.value = e.clientX
  categoryDropdownY.value = e.clientY
}

const showLanguageDropdown = (language, e) => {
  // 排除"全部"语言
  if (language.id === 'all') return
  
  e.preventDefault()
  languageDropdownShow.value = true
  languageDropdownItem.value = language
  languageDropdownX.value = e.clientX
  languageDropdownY.value = e.clientY
}
const dropdownOptions = computed(() => [
  { label: '编辑模板信息', key: 'editInfo', icon: () => h('span', { style: 'color:#18a058' }, '✏️') },
  { label: '编辑模板内容', key: 'editContent', icon: () => h('span', { style: 'color:#2080f0' }, '📄') },
  { 
    label: dropdownTemplate.value?.isFeatured ? '取消推荐' : '设为推荐', 
    key: 'toggleFeatured', 
    icon: () => h('span', { style: 'color:#f0a020' }, dropdownTemplate.value?.isFeatured ? '⭐' : '☆') 
  },
  { label: '删除模板', key: 'deleteTemplate', icon: () => h('span', { style: 'color:#d03050' }, '🗑️') }
])

const categoryDropdownOptions = [
  { label: '编辑分类', key: 'editCategory', icon: () => h('span', { style: 'color:#18a058' }, '✏️') },
  { label: '删除分类', key: 'deleteCategory', icon: () => h('span', { style: 'color:#d03050' }, '🗑️') }
]

const languageDropdownOptions = [
  { label: '编辑语言', key: 'editLanguage', icon: () => h('span', { style: 'color:#18a058' }, '✏️') },
  { label: '删除语言', key: 'deleteLanguage', icon: () => h('span', { style: 'color:#d03050' }, '🗑️') }
]
const showEditModal = ref(false)
const editFormRef = ref(null)
const editForm = ref({
  id: null,
  name: '',
  description: '',
  categoryId: null,
  languages: [],
  primary_language: null,
  icon: '',
  introduction: '',
  isFeatured: 0
})

const openEditModal = (template) => {
  editForm.value.id = template.id
  editForm.value.name = template.name
  editForm.value.description = template.description
  editForm.value.categoryId = template.categoryId || template.category_id
  editForm.value.introduction = template.introduction || ''
  editForm.value.isFeatured = template.isFeatured || 0
  editForm.value.icon = template.icon || ''
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
  editForm.value.categoryId = Number(editForm.value.categoryId)
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
    introduction: editForm.value.introduction,
    categoryId: editForm.value.categoryId,
    isFeatured: editForm.value.isFeatured,

    icon: editForm.value.icon,
    languages: languagesArr
  })
  showEditModal.value = false
  // 刷新模板列表
  const res = await listTemplates({})
  allTemplates.value = res.data.data.templatesList || []
}

const handleDropdownSelect = async (key, template) => {
  if (key === 'editInfo') {
    openEditModal(template)
  } else if (key === 'editContent') {
    router.push(`/templates/edit/${template.id}`)
  } else if (key === 'toggleFeatured') {
    await handleToggleFeatured(template)
  } else if (key === 'deleteTemplate') {
    await handleDeleteTemplate(template)
  }
  dropdownShow.value = false
}

const handleToggleFeatured = async (template) => {
  try {
    const newFeaturedStatus = template.isFeatured ? 0 : 1
    await editTemplate({
      id: template.id,
      name: template.name,
      description: template.description,
      introduction: template.introduction || '',
      categoryId: template.categoryId,
      isFeatured: newFeaturedStatus,
  
      languages: template.languages
    })
    // 刷新模板列表以确保数据一致性
    const res = await listTemplates({})
    allTemplates.value = res.data.data.templatesList || []
  } catch (error) {
    console.error('切换推荐状态失败:', error)
  }
}

// 删除模板
const handleDeleteTemplate = async (template) => {
  // 设置要删除的模板信息
  deleteTemplateInfo.value = template
  // 显示删除确认对话框
  showDeleteModal.value = true
}

// 确认删除
// 获取删除项目类型
const getDeleteItemType = () => {
  if (!deleteTemplateInfo.value) return '项目'
  if (deleteTemplateInfo.value.type === 'category') return '分类'
  if (deleteTemplateInfo.value.type === 'language') return '语言'
  return '模板'
}

// 确认删除
const confirmDelete = async () => {
  if (!deleteTemplateInfo.value) return
  
  try {
    deleteLoading.value = true
    
    if (deleteTemplateInfo.value.type === 'category') {
      // 删除分类
      await deleteCategory({ id: deleteTemplateInfo.value.id })
      message.success('分类删除成功')
      // 清空缓存并重新请求分类列表
      categoryStore.loaded = false
      await categoryStore.fetchCategories(true)
    } else if (deleteTemplateInfo.value.type === 'language') {
      // 删除语言
      await deleteLanguage({ id: deleteTemplateInfo.value.id })
      message.success('语言删除成功')
      // 清空缓存并重新请求语言列表
      languageStore.loaded = false
      await languageStore.fetchLanguages(true)
    } else {
      // 删除模板
      await deleteTemplate({ id: deleteTemplateInfo.value.id })
      message.success('模板删除成功')
      // 刷新模板列表
      const res = await listTemplates({})
      allTemplates.value = res.data.data.templatesList || []
    }
    
    // 关闭对话框并重置状态
    showDeleteModal.value = false
    deleteTemplateInfo.value = null
  } catch (error) {
    console.error('删除失败:', error)
    const itemType = deleteTemplateInfo.value?.type === 'category' ? '分类' : 
                    deleteTemplateInfo.value?.type === 'language' ? '语言' : '模板'
    message.error(`删除${itemType}失败: ` + (error.response?.data?.message || error.message || '未知错误'))
  } finally {
    deleteLoading.value = false
  }
}

const handleCategoryDropdownSelect = async (key, category) => {
  if (key === 'editCategory') {
    openEditCategoryModal(category)
  } else if (key === 'deleteCategory') {
    await handleDeleteCategory(category)
  }
  categoryDropdownShow.value = false
}

const handleLanguageDropdownSelect = async (key, language) => {
  if (key === 'editLanguage') {
    openEditLanguageModal(language)
  } else if (key === 'deleteLanguage') {
    await handleDeleteLanguage(language)
  }
  languageDropdownShow.value = false
}

// 分类相关处理函数
const handleAddCategory = async () => {
  try {
    // 清空缓存并重新请求分类列表
    categoryStore.loaded = false
    await categoryStore.fetchCategories(true)
    // 重置表单
    categoryForm.value = {
      name: '',
      description: '',
      sort: 0
    }
    showAddCategoryModal.value = false
  } catch (error) {
    console.error('添加分类失败:', error)
  }
}

const handleCancelCategory = () => {
  // 重置表单
  categoryForm.value = {
    name: '',
    description: '',
    sort: 0
  }
  showAddCategoryModal.value = false
  showEditCategoryModal.value = false
}

const openEditCategoryModal = (category) => {
  categoryForm.value = {
    id: category.id,
    name: category.name,
    description: category.description || category.Description || '',
    sort: category.sort || category.Sort || 0
  }
  showEditCategoryModal.value = true
}

const handleEditCategory = async () => {
  try {
    // 清空缓存并重新请求分类列表
    categoryStore.loaded = false
    await categoryStore.fetchCategories(true)
    // 重置表单
    categoryForm.value = {
      name: '',
      description: '',
      sort: 0
    }
    showEditCategoryModal.value = false
  } catch (error) {
    console.error('编辑分类失败:', error)
  }
}

const handleDeleteCategory = async (category) => {
  // 设置要删除的分类信息
  deleteTemplateInfo.value = { ...category, type: 'category' }
  // 显示删除确认对话框
  showDeleteModal.value = true
}

// 语言相关处理函数
const handleAddLanguage = async () => {
  try {
    // 清空缓存并重新请求语言列表
    languageStore.loaded = false
    await languageStore.fetchLanguages(true)
    // 重置表单
    languageForm.value = {
      name: '',
      code: '',
      display_name: '',
      color: '#18a058',
      sort: 0,
      isPopular: 0
    }
    showAddLanguageModal.value = false
  } catch (error) {
    console.error('添加语言失败:', error)
  }
}

const handleCancelLanguage = () => {
  // 重置表单
  languageForm.value = {
    name: '',
    code: '',
    display_name: '',
    color: '#18a058',
    sort: 0,
    isPopular: 0
  }
  showAddLanguageModal.value = false
  showEditLanguageModal.value = false
}

const openEditLanguageModal = (language) => {
  languageForm.value = {
    id: language.id,
    name: language.name,
    code: language.code || language.Code || '',
    display_name: language.displayName || language.display_name || language.DisplayName || language.name,
    color: language.color || language.Color || '#18a058',
    sort: language.sort || language.Sort || 0,
    isPopular: language.isPopular || language.IsPopular || 0
  }
  showEditLanguageModal.value = true
}

const handleEditLanguage = async () => {
  try {
    // 清空缓存并重新请求语言列表
    languageStore.loaded = false
    await languageStore.fetchLanguages(true)
    // 重置表单
    languageForm.value = {
      name: '',
      code: '',
      display_name: '',
      color: '#18a058',
      sort: 0,
      isPopular: 0
    }
    showEditLanguageModal.value = false
  } catch (error) {
    console.error('编辑语言失败:', error)
  }
}

const handleDeleteLanguage = async (language) => {
  // 设置要删除的语言信息
  deleteTemplateInfo.value = { ...language, type: 'language' }
  // 显示删除确认对话框
  showDeleteModal.value = true
}

watch(languagesList, (val) => {
})

// 初始化
onMounted(async () => {
  await languageStore.getLanguages()
  await categoryStore.getCategories()
  // 获取真实模板数据
  const res = await listTemplates({})
  allTemplates.value = res.data.data.templatesList || []
  // 可以在这里加载数据
})
</script>

<style scoped>
.templates-page {
  width: 100%;
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
  padding: 24px 0 16px 0;
}

.category-tags-title {
  display: flex;
  align-items: center;
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  position: relative;
}

.category-tags-title .title-actions {
  margin-left: auto;
  opacity: 0;
  transform: translateX(10px);
  transition: all 0.3s ease;
}

.category-tags-section:hover .title-actions,
.tags-section:hover .title-actions {
  opacity: 1;
  transform: translateX(0);
}

.add-btn {
  font-size: 12px;
  padding: 4px;
  height: 28px;
  width: 28px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.category-tags-icon {
  font-size: 1.4rem;
  margin-right: 8px;
  color: #18a058;
}

.category-tags-label {
  font-size: 1.1rem;
  color: #333;
}

.category-tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.category-tag-item {
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 20px;
  font-weight: 500;
  font-size: 14px;
  padding: 6px 16px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.category-tag-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(24, 160, 88, 0.1), transparent);
  transition: left 0.4s;
}

.category-tag-item:hover::before {
  left: 100%;
}

.category-tag-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0,0,0,0.12);
}

.category-tag-item:active {
  transform: translateY(0);
  transition: all 0.1s;
}

/* 标签列表样式 */
.tags-section {
  background: #f8f9fa;
  padding: 20px 0;
}

.tags-title {
  display: flex;
  align-items: center;
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  position: relative;
}

.tags-title .title-actions {
  margin-left: auto;
  opacity: 0;
  transform: translateX(10px);
  transition: all 0.3s ease;
}

.tags-icon {
  font-size: 1.4rem;
  margin-right: 8px;
  color: #18a058;
}

.tags-label {
  font-size: 1.1rem;
  color: #333;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.tag-item {
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 16px;
  font-weight: 500;
  font-size: 13px;
  padding: 5px 14px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 6px rgba(0,0,0,0.08);
}

.tag-item::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  background: rgba(24, 160, 88, 0.1);
  border-radius: 50%;
  transform: translate(-50%, -50%);
  transition: all 0.3s ease;
}

.tag-item:hover::before {
  width: 120%;
  height: 120%;
}

.tag-item:hover {
  transform: translateY(-1px);
  box-shadow: 0 3px 12px rgba(0,0,0,0.12);
}

.tag-item:active {
  transform: translateY(0);
  transition: all 0.1s;
}

/* 模板列表样式 */
.templates-section {
  padding: 40px 0;
  background: #fff;
}

.templates-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 30px;
}

.templates-title-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.templates-header-actions {
  display: flex;
  align-items: center;
  gap: 18px;
}

.add-template-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.add-template-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 3px 10px rgba(24, 160, 88, 0.2);
}
.templates-header h2 {
  margin: 0;
  color: #333;
}

.templates-count {
  color: #666;
  font-size: 14px;
  margin-top: 4px;
}

.templates-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 20px;
  margin-bottom: 40px;
}

.template-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.08);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid #f0f0f0;
  position: relative;
  overflow: hidden;
}

.template-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: #18a058;
  transform: scaleX(0);
  transition: transform 0.3s ease;
}

.template-card:hover::before {
  transform: scaleX(1);
}

.template-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 6px 20px rgba(0,0,0,0.12);
  border-color: rgba(24, 160, 88, 0.2);
}

.template-card:active {
  transform: translateY(-2px);
  transition: all 0.1s;
}

.template-logo {
  text-align: center;
  margin-bottom: 20px;
  position: relative;
}

.template-logo img {
  width: 60px;
  height: 60px;
}

.template-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 60px;
  height: 60px;
  margin: 0 auto;
}

.template-icon :deep(svg) {
  width: 48px;
  height: 48px;
  color: #18a058;
}

.featured-badge {
  position: absolute;
  top: -5px;
  right: 50%;
  transform: translateX(50%);
  background: #f0a020;
  color: white;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(240, 160, 32, 0.3);
  z-index: 1;
}

.template-card.featured {
  border-color: #f0a020;
  box-shadow: 0 2px 12px rgba(240, 160, 32, 0.15);
}

.template-card.featured::before {
  background: #f0a020;
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

.template-tags {
  margin-bottom: 20px;
}

.template-category {
  margin-bottom: 12px;
}

.category-tag {
  border-radius: 8px;
  font-weight: 500;
  font-size: 12px;
  padding: 3px 10px;
  transition: all 0.2s ease;
  background: #f0f9f4;
  color: #18a058;
  border: 1px solid #d9f2e6;
}

.category-tag:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(24, 160, 88, 0.15);
  background: #e6f7ed;
}

.template-languages {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.language-tag {
  border-radius: 8px;
  font-weight: 500;
  font-size: 12px;
  padding: 3px 10px;
  transition: all 0.2s ease;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
}

.language-tag:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  background: #f1f3f4;
}

.template-actions {
  display: flex;
  gap: 12px;
}

.template-actions .n-button {
  border-radius: 6px;
  font-weight: 500;
  font-size: 13px;
  padding: 6px 16px;
  transition: all 0.2s ease;
}

.template-actions .n-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 3px 10px rgba(0,0,0,0.1);
}

.template-actions .n-button:active {
  transform: translateY(0);
  transition: all 0.1s;
}

/* 分页样式 */
.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .templates-grid {
    grid-template-columns: repeat(4, 1fr);
    gap: 18px;
  }
}

@media (max-width: 992px) {
  .templates-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
  }
}

@media (max-width: 768px) {
  .templates-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 15px;
  }
  
  .template-actions {
    flex-direction: column;
  }
}

@media (max-width: 480px) {
  .templates-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }
}

/* 删除确认对话框样式 */
.delete-confirm-content {
  text-align: center;
  padding: 20px 0;
}

.delete-icon {
  font-size: 48px;
  margin-bottom: 16px;
  animation: shake 0.5s ease-in-out;
}

.delete-message {
  font-size: 16px;
  color: #333;
  margin-bottom: 12px;
  line-height: 1.5;
}

.delete-message strong {
  color: #d03050;
  font-weight: 600;
}

.delete-warning {
  font-size: 14px;
  color: #d03050;
  background: #fff2f0;
  padding: 12px;
  border-radius: 6px;
  border: 1px solid #ffccc7;
  margin-top: 12px;
}

.delete-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-5px); }
  75% { transform: translateX(5px); }
}
</style> 