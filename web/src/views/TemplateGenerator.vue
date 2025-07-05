<template>
  <div class="template-generator">
    <!-- 顶部导航 -->
    <div class="generator-header">
      <div class="header-left">
        <n-button text @click="goBack">
          <template #icon>
            <n-icon><ArrowBack /></n-icon>
          </template>
          返回
        </n-button>
        <div class="template-title">
          <span class="title-text">使用模板</span>
          <span class="template-name" v-if="templateInfo">{{ templateInfo.name }}</span>
        </div>
      </div>
      
      <!-- 步骤指示器 -->
      <div class="step-indicator">
        <div 
          v-for="(step, index) in steps" 
          :key="index"
          class="step-item"
          :class="{ 
            'active': currentStep === index + 1,
            'completed': currentStep > index + 1
          }"
        >
          <div class="step-number">{{ index + 1 }}</div>
          <div class="step-label">{{ step.label }}</div>
        </div>
      </div>
    </div>

    <!-- 步骤内容 -->
    <div class="generator-content">
      <!-- 步骤1: 模板介绍 -->
      <div v-show="currentStep === 1" class="step-content">
        <StepIntro 
          :template-info="templateInfo"
          @next="nextStep"
        />
      </div>

      <!-- 步骤2: 变量配置 -->
      <div v-show="currentStep === 2" class="step-content">
        <StepVariables 
          :template-info="templateInfo"
          :variables="variables"
          @prev="prevStep"
          @next="nextStep"
          @update-variables="updateVariables"
        />
      </div>

      <!-- 步骤3: 预览确认 -->
      <div v-show="currentStep === 3" class="step-content">
        <StepPreview 
          :template-info="templateInfo"
          :variables="variables"
          @prev="prevStep"
          @generate="generateProject"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { ArrowBack } from '@vicons/ionicons5'
import StepIntro from './generator/StepIntro.vue'
import StepVariables from './generator/StepVariables.vue'
import StepPreview from './generator/StepPreview.vue'

const route = useRoute()
const router = useRouter()
const message = useMessage()

// 步骤配置
const steps = [
  { label: '模板介绍', key: 'intro' },
  { label: '配置变量', key: 'variables' },
  { label: '预览确认', key: 'preview' }
]

// 状态
const currentStep = ref(1)
const templateInfo = ref(null)
const variables = ref({})

// 获取模板信息
const loadTemplateInfo = async () => {
  try {
    // TODO: 调用API获取模板信息
    // const res = await getTemplateInfo(route.params.id)
    // templateInfo.value = res.data.data
    
    // 临时使用模拟数据
    templateInfo.value = {
      id: route.params.id,
      name: 'Vue3项目模板',
      description: '现代化的Vue3项目模板，包含Vue3、Vite、TypeScript等最新技术栈',
      category: '前端框架',
      fileCount: 15,
      usage: '适用于快速搭建现代化的前端项目',
      features: [
        'Vue3 Composition API',
        'TypeScript支持',
        'Vite构建工具',
        'ESLint + Prettier',
        '单元测试配置'
      ]
    }
  } catch (error) {
    message.error('加载模板信息失败')
    console.error(error)
  }
}

// 步骤导航
const nextStep = () => {
  if (currentStep.value < 3) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

// 更新变量
const updateVariables = (newVariables) => {
  variables.value = { ...newVariables }
}

// 生成项目
const generateProject = async () => {
  try {
    message.loading('正在生成项目...')
    // TODO: 调用生成API
    // const res = await generateProject(route.params.id, variables.value)
    
    // 模拟生成过程
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    message.success('项目生成成功！')
    // TODO: 处理下载逻辑
  } catch (error) {
    message.error('生成失败：' + error.message)
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}

onMounted(() => {
  loadTemplateInfo()
})
</script>

<style scoped>
.template-generator {
  min-height: 100vh;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
}

.generator-header {
  background: #fff;
  border-bottom: 1px solid #e0e0e0;
  padding: 16px 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 8px rgba(0,0,0,0.03);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.template-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-text {
  font-size: 16px;
  color: #666;
}

.template-name {
  font-size: 18px;
  font-weight: bold;
  color: #18a058;
}

.step-indicator {
  display: flex;
  align-items: center;
  gap: 32px;
}

.step-item {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.step-number {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #f0f0f0;
  color: #999;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  transition: all 0.3s;
}

.step-label {
  font-size: 14px;
  color: #666;
  transition: all 0.3s;
}

.step-item.active .step-number {
  background: #18a058;
  color: #fff;
}

.step-item.active .step-label {
  color: #18a058;
  font-weight: bold;
}

.step-item.completed .step-number {
  background: #52c41a;
  color: #fff;
}

.step-item.completed .step-label {
  color: #52c41a;
}

.generator-content {
  flex: 1;
  padding: 32px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.step-content {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.08);
  overflow: hidden;
}
</style> 