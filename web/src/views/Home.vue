<template>
  <div class="home-page">
    <!-- 英雄区域 -->
    <div class="hero-section">
      <div class="hero-content">
        <h1>Template Starter</h1>
        <p>从数百个精心设计的项目模板中选择，一键生成完整项目结构</p>
      </div>
    </div>

    <!-- 推荐模板 -->
    <div class="templates-section" v-if="featuredTemplates && featuredTemplates.length > 0">
      <div class="container">
        <h2 class="section-title">
          <span class="category-icon-mark">⭐</span>
          推荐模板
        </h2>
        <div v-if="loading" class="loading-container">
          <n-spin size="large" />
          <p>加载中...</p>
        </div>
        <div v-else class="templates-grid">
          <div 
            v-for="template in featuredTemplates" 
            :key="template.id" 
            class="template-card"
          >
            <div class="template-logo">
              <div class="template-icon">
                <n-icon size="48">
                  <component :is="getIconComponent(template.icon) || getDefaultIcon(template)" />
                </n-icon>
              </div>
            </div>
            <div class="template-info">
              <h3>{{ template.name }}</h3>
              <p>{{ template.description }}</p>
              <div class="template-languages">
                <n-tag 
                  v-for="lang in template.languages" 
                  :key="lang.id"
                  :color="{ color: getLanguageColor(lang.languageId) }"
                  size="small"
                >
                  {{ getLanguageName(lang.languageId) }}
                </n-tag>
              </div>
              <n-button type="primary" @click="useTemplate(template)">
                使用模板
              </n-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 分类模板 -->
    <div class="categories-templates-section">
      <div class="container">
        <div v-if="loading" class="loading-container">
          <n-spin size="large" />
          <p>加载中...</p>
        </div>
        <div v-else-if="categories.length === 0" class="empty-container">
          <p>暂无分类数据</p>
        </div>
        <div v-else>
          <div 
            v-for="category in categories" 
            :key="category.id" 
            class="category-section"
            v-show="category.templates && category.templates.length > 0"
          >
            <div class="category-header">
              <h2 class="category-title">
                {{ category.name }}
              </h2>
              <p class="category-description">{{ category.description }}</p>
              <span class="template-count-badge">{{ category.templateCount }} 个模板</span>
            </div>
            <div class="templates-grid">
              <div 
                v-for="template in category.templates" 
                :key="template.id" 
                class="template-card"
              >
                <div class="template-logo">
                  <div class="template-icon">
                    <n-icon size="48">
                      <component :is="getIconComponent(template.icon) || getDefaultIcon(template)" />
                    </n-icon>
                  </div>
                </div>
                <div class="template-info">
                  <h3>{{ template.name }}</h3>
                  <p>{{ template.description }}</p>
                  <div class="template-languages">
                    <n-tag 
                      v-for="lang in template.languages" 
                      :key="lang.id"
                      :color="{ color: getLanguageColor(lang.languageId) }"
                      size="small"
                    >
                      {{ getLanguageName(lang.languageId) }}
                    </n-tag>
                  </div>
                  <n-button type="primary" @click="useTemplate(template)">
                    使用模板
                  </n-button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getIndexData } from '@/api/indexData'
import { useLanguageStore } from '@/stores/languageStore'
import { storeToRefs } from 'pinia'
import * as IonIcons from '@vicons/ionicons5'

const router = useRouter()
const languageStore = useLanguageStore()
const { languagesList } = storeToRefs(languageStore)

// 响应式数据
const loading = ref(false)
const statistics = ref({
  totalTemplates: 0,
  totalCategories: 0,
  totalLanguages: 0,
  featuredCount: 0
})
const categories = ref([])
const featuredTemplates = ref([])

// 使用模板
const useTemplate = (template) => {
  router.push(`/templates/generate/${template.id}`)
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
  // 根据模板名称选择合适的默认图标
  const name = template.name?.toLowerCase() || ''
  
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

// 获取首页数据
const fetchIndexData = async () => {
  loading.value = true
  try {
    const response = await getIndexData({
      categoryLimit: 6,
      featuredLimit: 8
    })
    
    if (response.data && response.data.data) {
      statistics.value = response.data.data.statistics || {}
      categories.value = response.data.data.categories || []
      featuredTemplates.value = response.data.data.featuredTemplates || []
    }
  } catch (error) {
  } finally {
    loading.value = false
  }
}

// 页面加载时获取数据
onMounted(async () => {
  try {
    // 先获取语言列表
    await languageStore.fetchLanguages()
    // 再获取首页数据
    await fetchIndexData()
  } catch (error) {
  }
})
</script>

<style scoped>
.home-page {
  width: 100%;
}

.hero-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 80px 0;
  text-align: center;
}

.hero-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 20px;
}

.hero-content h1 {
  font-size: 3rem;
  margin-bottom: 20px;
}

.hero-content p {
  font-size: 1.2rem;
  margin-bottom: 0;
  opacity: 0.9;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.categories-templates-section {
  padding: 60px 0;
  background: #fff;
}

.category-section {
  margin-bottom: 60px;
}

.category-section:last-child {
  margin-bottom: 0;
}

.category-header {
  text-align: left;
  margin-bottom: 40px;
  padding-bottom: 20px;
  border-bottom: 2px solid #f0f0f0;
}

.category-title {
  font-size: 2rem;
  margin-bottom: 10px;
  color: #333;
}

.category-description {
  color: #666;
  font-size: 1.1rem;
  margin-bottom: 15px;
}

.template-count-badge {
  background: #18a058;
  color: white;
  padding: 6px 16px;
  border-radius: 16px;
  font-size: 14px;
  font-weight: 500;
  display: inline-block;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  color: #666;
}

.loading-container p {
  margin-top: 16px;
  font-size: 14px;
}

.empty-container {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  color: #999;
  font-size: 14px;
}

.templates-section {
  padding: 60px 0;
  background: #f8f9fa;
}

.templates-section h2 {
  text-align: center;
  margin-bottom: 40px;
  font-size: 2rem;
}

.templates-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 20px;
}

@media (max-width: 1200px) {
  .templates-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 992px) {
  .templates-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .templates-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .templates-grid {
    grid-template-columns: 1fr;
  }
}

.template-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 5px 15px rgba(0,0,0,0.1);
  transition: transform 0.3s ease;
  min-width: 0;
}

.template-card:hover {
  transform: translateY(-5px);
}

.template-logo {
  text-align: center;
  margin-bottom: 20px;
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

.template-info h3 {
  margin-bottom: 10px;
  color: #333;
}

.template-info p {
  color: #666;
  margin-bottom: 15px;
}

.template-languages {
  margin-bottom: 20px;
}

.template-languages .n-tag {
  margin-right: 8px;
  margin-bottom: 8px;
}

.section-title {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  margin-bottom: 40px;
}

.category-icon-mark {
  font-size: 2.2rem;
  margin-right: 12px;
  color: #18a058;
  display: flex;
  align-items: center;
}

.category-badge {
  position: absolute;
  top: 16px;
  left: 16px;
  background: #18a058;
  color: #fff;
  font-size: 12px;
  padding: 2px 10px;
  border-radius: 10px;
  font-weight: bold;
  letter-spacing: 2px;
  z-index: 1;
  box-shadow: 0 2px 8px rgba(24,160,88,0.08);
}
</style> 