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
        <div class="tags-title">
          <span class="tags-icon">💻</span>
          <span class="tags-label">语言</span>
        </div>
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

    <!-- 模板类型列表 -->
    <div class="template-types-section">
      <div class="container">
        <div class="template-types-title">
          <span class="template-types-icon">🏗️</span>
          <span class="template-types-label">类型</span>
        </div>
        <div class="template-types-list">
          <n-tag
            v-for="type in templateTypes"
            :key="type.value"
            :type="selectedTemplateType === type.value ? 'primary' : 'default'"
            size="large"
            class="template-type-item"
            @click="selectTemplateType(type.value)"
          >
            {{ type.label }}
          </n-tag>
        </div>
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
          <div class="view-controls">
            <n-button-group>
              <n-button 
                :type="viewMode === 'grid' ? 'primary' : 'default'"
                @click="viewMode = 'grid'"
                size="small"
              >
                <template #icon>
                  <n-icon><GridOutline /></n-icon>
                </template>
                卡片
              </n-button>
              <n-button 
                :type="viewMode === 'list' ? 'primary' : 'default'"
                @click="viewMode = 'list'"
                size="small"
              >
                <template #icon>
                  <n-icon><ListOutline /></n-icon>
                </template>
                列表
              </n-button>
            </n-button-group>
          </div>
        </div>
        
        <n-spin :show="loading">
          <!-- 卡片视图 -->
          <div v-if="viewMode === 'grid'" class="templates-grid">
            <div
              v-for="template in templates"
              :key="template.id"
              class="template-card"
              :class="{ 'featured': template.isFeatured, 'selected': selectedTemplateId === template.id }"
              @click="selectTemplate(template)"
            >
              <div class="template-logo">
                <div class="template-icon">
                  <n-icon size="48">
                    <component :is="getIconComponent(template.icon) || getDefaultIcon(template)" />
                  </n-icon>
                </div>
                <div v-if="template.isFeatured" class="featured-badge">
                  <n-icon><StarOutline /></n-icon>
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
                  <n-button type="primary" @click.stop="useTemplate(template)">
                    使用模板
                  </n-button>
                </div>
              </div>
            </div>
          </div>

          <!-- 列表视图 -->
          <div v-else class="templates-list">
            <div
              v-for="template in templates"
              :key="template.id"
              class="template-list-card"
              :class="{ 'featured': template.isFeatured, 'selected': selectedTemplateId === template.id }"
              @click="selectTemplate(template)"
            >
              <div class="card-left">
                <div class="template-icon-medium">
                  <n-icon size="40">
                    <component :is="getIconComponent(template.icon) || getDefaultIcon(template)" />
                  </n-icon>
                </div>
              </div>
              <div class="card-content">
                <div class="title-row">
                  <h4>{{ template.name }}</h4>
                  <div class="badges">
                    <n-tag v-if="template.isFeatured" type="warning" size="small">
                      <template #icon>
                        <n-icon><StarOutline /></n-icon>
                      </template>
                      推荐
                    </n-tag>
                    <n-tag type="success" size="small" v-if="getCategoryName(template.categoryId)">
                      {{ getCategoryName(template.categoryId) }}
                    </n-tag>
                  </div>
                </div>
                <p class="description">{{ template.description }}</p>
                <div class="footer-row">
                  <div class="languages">
                    <n-tag 
                      v-for="lang in template.languages" 
                      :key="lang.id"
                      :color="{ color: getLanguageColor(lang.languageId) }"
                      size="small"
                    >
                      {{ getLanguageName(lang.languageId) }}
                    </n-tag>
                  </div>
                  <div class="list-stats">
                    <span class="stat-item">
                      <n-icon size="14"><EyeOutline /></n-icon>
                      {{ formatNumber(template.views || 0) }}
                    </span>
                    <span class="stat-item">
                      <n-icon size="14"><HeartOutline /></n-icon>
                      {{ formatNumber(template.likes || 0) }}
                    </span>
                    <span class="stat-item">
                      <n-icon size="14"><DownloadOutline /></n-icon>
                      {{ formatNumber(template.downloads || 0) }}
                    </span>
                  </div>
                  <div class="actions">
                    <n-button size="small" @click.stop="previewTemplate(template)">
                      预览
                    </n-button>
                    <n-button size="small" type="primary" @click.stop="useTemplate(template)">
                      使用模板
                    </n-button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </n-spin>

        <!-- 分页 -->
        <div class="pagination-section" v-if="totalPages > 1">
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

    <!-- 模板预览抽屉 -->
    <n-drawer
      v-model:show="showPreview"
      :width="800"
      placement="right"
    >
      <n-drawer-content title="模板预览" closable>
        <TemplatePreview v-if="selectedTemplate" :template="selectedTemplate" />
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NIcon, NTag, NSpin, NPagination, NButton, NDrawer, NDrawerContent, NButtonGroup, useMessage } from 'naive-ui'
import * as IonIcons from '@vicons/ionicons5'
import { GridOutline, ListOutline, StarOutline, EyeOutline, HeartOutline, DownloadOutline } from '@vicons/ionicons5'
import { useLanguageStore } from '@/stores/languageStore'
import { storeToRefs } from 'pinia'
import { useCategoryStore } from '@/stores/categoryStore'
import { listTemplates, getTemplateTypes } from '@/api/templates'
import TemplatePreview from '@/components/TemplatePreview.vue'

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
const selectedTemplateType = ref('all')
const currentPage = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const showPreview = ref(false)
const selectedTemplate = ref(null)
const selectedTemplateId = ref(null)
const viewMode = ref('grid') // 'grid' | 'list'

// 标签数据
const tags = computed(() => {
  const arr = [
    { id: 'all', name: '全部' },
    ...languagesList.value.map(lang => ({ id: lang.id, name: lang.name }))
  ]
  return arr
})

// 模板数据
const allTemplates = ref([])

// 模板类型数据
const allTemplateTypes = ref([])
const templateTypes = computed(() => {
  return [
    { value: 'all', label: '全部' },
    ...allTemplateTypes.value
  ]
})

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

  // 按模板类型筛选
  if (selectedTemplateType.value !== 'all') {
    filtered = filtered.filter(t => t.templateType === selectedTemplateType.value)
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

// 方法
const selectCategory = (catId) => {
  selectedCategory.value = catId === 'all' ? 'all' : Number(catId)
  currentPage.value = 1
}

const selectTag = (tagId) => {
  selectedTag.value = tagId
  currentPage.value = 1
}

const selectTemplateType = (typeValue) => {
  selectedTemplateType.value = typeValue
  currentPage.value = 1
}

const handlePageChange = (page) => {
  currentPage.value = page
}

const handlePageSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
}

const selectTemplate = (template) => {
  selectedTemplateId.value = template.id
}

const useTemplate = (template) => {
  router.push(`/template-generator/${template.id}`)
}

const previewTemplate = (template) => {
  selectedTemplate.value = template
  showPreview.value = true
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

// 格式化数字显示
const formatNumber = (num) => {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

// 加载模板数据
const loadTemplates = async () => {
  try {
    loading.value = true
    const res = await listTemplates({})
    allTemplates.value = res.data.data.templatesList || []
  } catch (error) {
    console.error('获取模板列表失败:', error)
    message.error('获取模板列表失败')
  } finally {
    loading.value = false
  }
}

// 加载模板类型数据
const loadTemplateTypes = async () => {
  try {
    const res = await getTemplateTypes()
    allTemplateTypes.value = res.data.data.templateTypes || []
  } catch (error) {
    console.error('获取模板类型失败:', error)
    // 静默失败，使用默认的空数组
  }
}

// 初始化
onMounted(async () => {
  await languageStore.getLanguages()
  await categoryStore.getCategories()
  await loadTemplateTypes()
  await loadTemplates()
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

/* 模板类型区域样式 */
.template-types-section {
  background: #f0f2f5;
  padding: 20px 0;
}

.template-types-title {
  display: flex;
  align-items: center;
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  position: relative;
}

.template-types-icon {
  font-size: 1.4rem;
  margin-right: 8px;
  color: #18a058;
}

.template-types-label {
  font-size: 1.1rem;
  color: #333;
}

.template-types-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.template-type-item {
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 18px;
  font-weight: 500;
  font-size: 14px;
  padding: 6px 16px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  background: #fff;
}

.template-type-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(24, 160, 88, 0.1), transparent);
  transition: left 0.4s;
}

.template-type-item:hover::before {
  left: 100%;
}

.template-type-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0,0,0,0.12);
}

.template-type-item:active {
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

.view-controls {
  display: flex;
  align-items: center;
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

/* 卡片视图样式 */
.templates-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 20px;
  margin-bottom: 40px;
}

.template-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.08);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid #f0f0f0;
  position: relative;
  overflow: hidden;
  cursor: pointer;
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

.template-card.featured {
  border-color: #f0a020;
  box-shadow: 0 2px 12px rgba(240, 160, 32, 0.15);
}

.template-card.featured::before {
  background: #f0a020;
}

.template-card.selected {
  border-color: #18a058;
  box-shadow: 0 4px 16px rgba(24, 160, 88, 0.2);
  transform: translateY(-2px);
}

.template-card.selected::before {
  transform: scaleX(1);
}

.template-logo {
  text-align: center;
  margin-bottom: 16px;
  position: relative;
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

.template-info {
}

.template-info h3 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 16px;
  font-weight: 600;
}

.template-info p {
  color: #666;
  margin-bottom: 12px;
  line-height: 1.4;
  font-size: 14px;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.template-tags {
  margin-bottom: 16px;
}

.template-category {
  margin-bottom: 8px;
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


/* 列表视图样式 */
.templates-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 40px;
}

.template-list-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid #f0f0f0;
  cursor: pointer;
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.template-list-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
  border-color: rgba(24, 160, 88, 0.2);
}

.template-list-card.featured {
  border-color: #f0a020;
  background: rgba(240, 160, 32, 0.02);
}

.template-list-card.selected {
  border-color: #18a058;
  box-shadow: 0 4px 12px rgba(24, 160, 88, 0.15);
  transform: translateY(-1px);
  background: rgba(24, 160, 88, 0.02);
}

.card-left {
  flex-shrink: 0;
}

.template-icon-medium {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  background: rgba(24, 160, 88, 0.1);
  border-radius: 8px;
}

.template-icon-medium :deep(svg) {
  color: #18a058;
}

.card-content {
  flex: 1;
  min-width: 0;
}

.title-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 8px;
}

.title-row h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #333;
  flex: 1;
  min-width: 0;
}

.badges {
  display: flex;
  gap: 6px;
  margin-left: 12px;
}

.description {
  color: #666;
  font-size: 14px;
  line-height: 1.5;
  margin: 0 0 16px 0;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.footer-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.languages {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  flex: 1;
}

.list-stats {
  display: flex;
  gap: 16px;
  margin-right: 16px;
}

.actions {
  display: flex;
  gap: 8px;
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
  
  .footer-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .list-stats {
    margin-right: 0;
  }
}

@media (max-width: 768px) {
  .templates-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 15px;
  }
  
  .templates-header {
    flex-direction: column;
    gap: 16px;
  }
  
  .template-list-card {
    flex-direction: column;
    text-align: center;
  }
  
  .title-row {
    flex-direction: column;
    align-items: center;
    gap: 8px;
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
  
  .category-tags-list,
  .tags-list {
    justify-content: center;
  }
}
</style>