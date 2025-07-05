<template>
  <div class="step-intro">
    <div class="intro-content">
      <!-- 模板基本信息 -->
      <div class="template-info">
        <div class="template-header">
          <h1 class="template-name">{{ templateInfo?.name }}</h1>
          <n-tag :type="getCategoryType(templateInfo?.category)" size="large">
            {{ templateInfo?.category }}
          </n-tag>
        </div>
        
        <div class="template-description">
          {{ templateInfo?.description }}
        </div>
        
        <div class="template-stats">
          <div class="stat-item">
            <n-icon size="20" color="#18a058">
              <DocumentText />
            </n-icon>
            <span>{{ templateInfo?.fileCount }} 个文件</span>
          </div>
        </div>
      </div>

      <!-- 适用场景 -->
      <div class="section">
        <h3 class="section-title">
          <n-icon><Bulb /></n-icon>
          适用场景
        </h3>
        <div class="usage-text">
          {{ templateInfo?.usage }}
        </div>
      </div>

      <!-- 主要特性 -->
      <div class="section">
        <h3 class="section-title">
          <n-icon><Star /></n-icon>
          主要特性
        </h3>
        <div class="features-list">
          <div 
            v-for="feature in templateInfo?.features" 
            :key="feature"
            class="feature-item"
          >
            <n-icon size="16" color="#52c41a">
              <CheckmarkCircle />
            </n-icon>
            <span>{{ feature }}</span>
          </div>
        </div>
      </div>

      <!-- 使用说明 -->
      <div class="section">
        <h3 class="section-title">
          <n-icon><InformationCircle /></n-icon>
          使用说明
        </h3>
        <div class="instructions">
          <p>1. 在下一步中配置项目的基本信息</p>
          <p>2. 填写必要的变量值</p>
          <p>3. 预览生成的文件内容</p>
          <p>4. 确认无误后下载项目文件</p>
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
import { 
  DocumentText, 
  Bulb, 
  Star, 
  CheckmarkCircle, 
  InformationCircle,
  ArrowForward 
} from '@vicons/ionicons5'

const props = defineProps({
  templateInfo: {
    type: Object,
    default: null
  }
})

defineEmits(['next'])

// 获取分类标签类型
const getCategoryType = (category) => {
  const typeMap = {
    '前端框架': 'success',
    '后端框架': 'info',
    '全栈项目': 'warning',
    '工具库': 'error'
  }
  return typeMap[category] || 'default'
}
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

.template-stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #666;
  font-size: 14px;
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

.usage-text {
  font-size: 15px;
  color: #666;
  line-height: 1.6;
  background: #f8f9fa;
  padding: 16px;
  border-radius: 8px;
  border-left: 4px solid #18a058;
}

.features-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 12px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #f8f9fa;
  border-radius: 6px;
  font-size: 14px;
  color: #333;
}

.instructions {
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: 8px;
  padding: 20px;
}

.instructions p {
  margin: 8px 0;
  color: #333;
  font-size: 14px;
}

.step-actions {
  padding: 24px 40px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
  display: flex;
  justify-content: flex-end;
}
</style> 