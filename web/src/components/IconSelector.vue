<template>
  <div class="icon-selector">
    <div class="icon-selector-header">
      <n-input 
        v-model:value="searchText" 
        placeholder="搜索图标..." 
        clearable
        size="small"
      >
        <template #prefix>
          <n-icon><SearchOutline /></n-icon>
        </template>
      </n-input>
    </div>
    
    <div class="icon-list">
      <div 
        v-for="icon in filteredIcons" 
        :key="icon.name"
        class="icon-item"
        :class="{ 'selected': selectedIcon === icon.name }"
        @click="selectIcon(icon)"
        :title="icon.name"
      >
        <div class="icon-wrapper">
          <component :is="icon.component" />
        </div>
        <div class="icon-name">{{ formatIconName(icon.name) }}</div>
      </div>
    </div>
    
    <div v-if="filteredIcons.length === 0" class="no-results">
      <n-empty description="未找到匹配的图标" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { SearchOutline } from '@vicons/ionicons5'
import * as IonIcons from '@vicons/ionicons5'

// 响应式数据
const searchText = ref('')
const selectedIcon = ref('')
const allIcons = ref([])

// 计算过滤后的图标
const filteredIcons = computed(() => {
  if (!searchText.value) {
    return allIcons.value
  }
  
  const searchLower = searchText.value.toLowerCase()
  return allIcons.value.filter(icon => 
    icon.name.toLowerCase().includes(searchLower) ||
    formatIconName(icon.name).toLowerCase().includes(searchLower)
  )
})

// 格式化图标名称显示
const formatIconName = (name) => {
  return name
    .replace(/Outline$/, '')
    .replace(/([A-Z])/g, ' $1')
    .trim()
}

// 选择图标
const selectIcon = (icon) => {
  selectedIcon.value = icon.name
  emit('select', icon.name)
}

// 初始化图标列表 - 动态遍历所有图标
const initIcons = () => {
  const icons = []
  
  // 遍历IonIcons对象的所有属性
  Object.keys(IonIcons).forEach(iconName => {
    // 只处理以Outline结尾的图标（轮廓风格）
    if (iconName.endsWith('Outline') && typeof IonIcons[iconName] === 'function') {
      icons.push({
        name: iconName,
        component: IonIcons[iconName]
      })
    }
  })
  
  // 按名称排序
  icons.sort((a, b) => a.name.localeCompare(b.name))
  
  allIcons.value = icons
  console.log(`加载了 ${icons.length} 个图标`)
}

// 暴露方法
const reset = () => {
  selectedIcon.value = ''
  searchText.value = ''
}

const setSelectedIcon = (iconName) => {
  selectedIcon.value = iconName
}

// 定义事件
const emit = defineEmits(['select'])

// 暴露方法给父组件
defineExpose({
  reset,
  setSelectedIcon
})

// 组件挂载时初始化
onMounted(() => {
  initIcons()
})
</script>

<style scoped>
.icon-selector {
  width: 100%;
  max-height: 400px;
  display: flex;
  flex-direction: column;
}

.icon-selector-header {
  padding: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.icon-list {
  flex: 1;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 8px;
  padding: 12px;
}

.icon-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 2px solid transparent;
}

.icon-item:hover {
  background-color: #f5f5f5;
  border-color: #e0e0e0;
}

.icon-item.selected {
  background-color: #e8f5e8;
  border-color: #18a058;
}

.icon-wrapper {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
}

.icon-wrapper :deep(svg) {
  width: 24px;
  height: 24px;
  color: #666;
}

.icon-item.selected .icon-wrapper :deep(svg) {
  color: #18a058;
}

.icon-name {
  font-size: 12px;
  color: #666;
  text-align: center;
  line-height: 1.2;
  word-break: break-word;
}

.icon-item.selected .icon-name {
  color: #18a058;
  font-weight: 500;
}

.no-results {
  padding: 40px 20px;
  text-align: center;
}

/* 滚动条样式 */
.icon-list::-webkit-scrollbar {
  width: 6px;
}

.icon-list::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.icon-list::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.icon-list::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style> 