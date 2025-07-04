<template>
  <div class="step-intro">
    <div class="intro-content">
      <!-- 模板基本信息 -->
      <div class="template-info">
        <div class="template-header">
          <h1 class="template-name">{{ templateInfo?.name }}</h1>
          <n-tag v-if="templateInfo?.categoryId" :type="getCategoryType(templateInfo?.categoryId)" size="large">
            {{ getCategoryName(templateInfo?.categoryId) }}
          </n-tag>
        </div>
        
        <div class="template-description">
          {{ templateInfo?.description }}
        </div>
        
        <!-- 模板统计信息 -->
        <div v-if="templateInfo?.isFeatured" class="template-stats">
          <div class="stat-item">
            <n-icon size="20" color="#f0a020">
              <Star />
            </n-icon>
            <span>推荐模板</span>
          </div>
        </div>
        
        <!-- 支持的语言 -->
        <div v-if="getTemplateLanguages.length > 0" class="template-languages">
          <h3 class="languages-title">支持的语言</h3>
          <div class="languages-list">
            <n-tag 
              v-for="lang in getTemplateLanguages" 
              :key="lang.id"
              :color="{ color: lang.isPrimary === 1 ? '#e8f5e8' : lang.color }"
              size="large"
              style="margin-right: 12px; margin-bottom: 8px;"
            >
              {{ lang.displayName }}
              <template v-if="lang.isPrimary === 1">
                (主语言)
              </template>
            </n-tag>
          </div>
        </div>
        
        <!-- 详细介绍 -->
        <div v-if="templateInfo?.introduction" class="template-introduction">
          <h3 class="intro-title">详细介绍</h3>
          <div class="markdown-content">
            <MdPreview :modelValue="templateInfo.introduction" />
          </div>
        </div>
      </div>


    </div>

    <!-- 底部操作 -->
    <div class="step-actions">
      <n-button 
        type="primary" 
        size="large" 
        @click="$emit('next')"
        :disabled="!templateInfo"
      >
        开始配置
        <template #icon>
          <n-icon><ArrowForward /></n-icon>
        </template>
      </n-button>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, watch } from 'vue'
import { 
  DocumentText, 
  Star, 
  CodeSlash,
  ArrowForward 
} from '@vicons/ionicons5'
import { MdPreview } from 'md-editor-v3'
import { useLanguageStore } from '@/stores/languageStore'
import { useCategoryStore } from '@/stores/categoryStore'

const props = defineProps({
  templateInfo: {
    type: Object,
    default: null
  }
})

defineEmits(['next'])

// 使用 Pinia stores
const languageStore = useLanguageStore()
const categoryStore = useCategoryStore()



// 获取分类信息
const getCategoryInfo = computed(() => {
  if (!props.templateInfo?.categoryId) return null
  return categoryStore.categoriesList.find(cat => cat.id === props.templateInfo.categoryId)
})

// 获取分类标签类型
const getCategoryType = (categoryId) => {
  console.log('查找分类ID:', categoryId)
  console.log('Pinia分类列表:', categoryStore.categoriesList)
  const category = categoryStore.categoriesList.find(cat => cat.id === categoryId)
  console.log('找到的分类:', category)
  return 'success' // 统一使用success样式
}

// 获取分类名称
const getCategoryName = (categoryId) => {
  const category = categoryStore.categoriesList.find(cat => cat.id === categoryId)
  console.log('分类名称查找结果:', category)
  return category ? category.name : '未知分类'
}

// 获取模板的语言详细信息
const getTemplateLanguages = computed(() => {
  if (!props.templateInfo?.languages?.length) return []
  
  console.log('模板语言数据:', props.templateInfo.languages)
  console.log('Pinia语言列表:', languageStore.languagesList)
  
  const result = props.templateInfo.languages.map(templateLang => {
    const language = languageStore.languagesList.find(lang => lang.id === templateLang.languageId)
    console.log(`查找语言ID ${templateLang.languageId}:`, language)
    return {
      ...templateLang,
      name: language?.name || '未知语言',
      displayName: language?.displayName || language?.name || '未知语言',
      color: language?.color || '#18a058'
    }
  })
  
  console.log('最终语言数组:', result)
  return result
})

// 初始化数据
onMounted(async () => {
  try {
    await Promise.all([
      languageStore.getLanguages(),
      categoryStore.getCategories()
    ])
  } catch (error) {
    console.error('加载语言或分类数据失败:', error)
  }
})

watch(
  () => props.templateInfo,
  (val) => {
    console.log('StepIntro.vue 接收到的 templateInfo:', val)
  },
  { immediate: true, deep: true }
)
</script>

<style scoped>
.step-intro {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.intro-content {
  flex: 1;
  padding: 40px;
  overflow-y: auto;
}

.template-info {
  margin-bottom: 40px;
}

.template-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.template-name {
  font-size: 28px;
  font-weight: bold;
  color: #333;
  margin: 0;
}

.template-description {
  font-size: 16px;
  color: #666;
  line-height: 1.6;
  margin-bottom: 24px;
}

.template-languages {
  margin-bottom: 24px;
}

.languages-title {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin-bottom: 12px;
}

.template-introduction {
  margin-bottom: 24px;
}

.intro-title {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin-bottom: 12px;
}

.markdown-content {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.template-stats {
  display: flex;
  gap: 24px;
  margin-bottom: 24px;
  padding: 12px 16px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #666;
  font-size: 14px;
  font-weight: 500;
}

.section {
  margin-bottom: 32px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin-bottom: 16px;
}

.languages-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.step-actions {
  padding: 24px 40px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
  display: flex;
  justify-content: flex-end;
}
</style> 