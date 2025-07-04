<template>
  <div class="home-page">
    <!-- 英雄区域 -->
    <div class="hero-section">
      <div class="hero-content">
        <h1>快速启动您的项目</h1>
        <p>从数百个精心设计的项目模板中选择，一键生成完整项目结构</p>
        <div class="hero-search">
          <n-input-group>
            <n-input placeholder="搜索模板..." size="large" />
            <n-button type="primary" size="large">搜索</n-button>
          </n-input-group>
        </div>
      </div>
    </div>

    <!-- 热门分类 -->
    <div class="categories-section">
      <div class="container">
        <h2 class="section-title">
          <span class="category-icon-mark">📂</span>
          热门分类
        </h2>
        <div class="categories-grid">
          <div 
            v-for="category in categories" 
            :key="category.id" 
            class="category-card"
          >
            <span class="category-badge">分类</span>
            <div class="category-icon">{{ category.icon }}</div>
            <h3>{{ category.name }}</h3>
            <p>{{ category.description }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 推荐模板 -->
    <div class="templates-section">
      <div class="container">
        <h2>推荐模板</h2>
        <div class="templates-grid">
          <div 
            v-for="template in templates" 
            :key="template.id" 
            class="template-card"
          >
            <div class="template-logo">
              <img :src="template.logo" :alt="template.name" />
            </div>
            <div class="template-info">
              <h3>{{ template.name }}</h3>
              <p>{{ template.description }}</p>
              <div class="template-languages">
                <n-tag 
                  v-for="lang in languagesList" 
                  :key="lang.id"
                  :color="{ color: lang.color }"
                  size="small"
                >
                  {{ lang.displayName }}
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
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listLanguages } from '@/api/languages'
import { useLanguageStore } from '@/stores/languageStore'
import { storeToRefs } from 'pinia'

// 模拟数据 - 分类
const categories = ref([
  {
    id: 1,
    name: 'Web应用',
    description: '前端和后端Web应用模板',
    icon: '🌐',
    sort: 1
  },
  {
    id: 2,
    name: '移动应用',
    description: '移动端应用开发模板',
    icon: '📱',
    sort: 2
  },
  {
    id: 3,
    name: '桌面应用',
    description: '桌面应用程序模板',
    icon: '💻',
    sort: 3
  },
  {
    id: 4,
    name: 'API服务',
    description: '后端API服务模板',
    icon: '🔌',
    sort: 4
  }
])

// 模拟数据 - 模板
const templates = ref([
  {
    id: 1,
    name: 'Vue全栈应用',
    description: '基于Vue3 + Node.js的完整全栈应用模板',
    category_id: 1,
    is_featured: true,
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
    is_featured: true,
    logo: '/vite.svg',
    languages: [
      { id: 1, name: 'JavaScript', display_name: 'JS', color: '#f7df1e' },
      { id: 3, name: 'React', display_name: 'React', color: '#61dafb' }
    ]
  }
])

const languageStore = useLanguageStore()
const { languagesList } = storeToRefs(languageStore)

// 使用模板
const useTemplate = (template) => {
  console.log('使用模板:', template.name)
  // TODO: 跳转到模板使用页面
}

// 页面加载时调用语言列表接口
onMounted(async () => {
  try {
    await languageStore.fetchLanguages()
    console.log('语言列表（pinia）:', languagesList.value)
  } catch (e) {
    console.error('语言列表接口调用失败', e)
  }
})
</script>

<style scoped>
.home-page {
  min-height: 100vh;
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
  margin-bottom: 40px;
  opacity: 0.9;
}

.hero-search {
  max-width: 500px;
  margin: 0 auto;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.categories-section {
  padding: 60px 0;
  background: #fff;
}

.categories-section h2 {
  text-align: center;
  margin-bottom: 40px;
  font-size: 2rem;
}

.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 30px;
}

.category-card {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 30px;
  text-align: center;
  transition: all 0.3s ease;
  position: relative;
  cursor: pointer;
}

.category-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0,0,0,0.1);
}

.category-icon {
  font-size: 3rem;
  margin-bottom: 20px;
}

.category-card h3 {
  margin-bottom: 10px;
  color: #333;
}

.category-card p {
  color: #666;
  margin-bottom: 20px;
}

.templates-section {
  padding: 60px 0;
  background: #f5f5f5;
}

.templates-section h2 {
  text-align: center;
  margin-bottom: 40px;
  font-size: 2rem;
}

.templates-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 30px;
}

.template-card {
  background: white;
  border-radius: 12px;
  padding: 30px;
  box-shadow: 0 5px 15px rgba(0,0,0,0.1);
  transition: transform 0.3s ease;
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