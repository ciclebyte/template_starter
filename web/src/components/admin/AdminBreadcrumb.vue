<template>
  <div class="breadcrumb-container">
    <n-breadcrumb class="admin-breadcrumb">
      <n-breadcrumb-item
        v-for="(item, index) in breadcrumbItems"
        :key="`breadcrumb-${index}-${item.title}`"
        :clickable="item.clickable"
        @click="item.clickable ? handleBreadcrumbClick(item) : null"
      >
        <template #icon v-if="item.icon">
          <n-icon>
            <component :is="item.icon" />
          </n-icon>
        </template>
        {{ item.title }}
      </n-breadcrumb-item>
    </n-breadcrumb>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NBreadcrumb, NBreadcrumbItem, NIcon } from 'naive-ui'
import {
  GridOutline,
  DocumentTextOutline,
  LayersOutline,
  LanguageOutline,
  SettingsOutline,
  StatsChartOutline,
  AddOutline,
  CreateOutline
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()

// 路由对应的面包屑配置
const breadcrumbConfig = {
  'admin-dashboard': {
    title: '仪表盘',
    icon: GridOutline,
    path: '/admin'
  },
  'admin-templates-list': {
    title: '模板管理',
    icon: DocumentTextOutline,
    path: '/admin/templates'
  },
  'admin-templates-create': {
    title: '创建模板',
    icon: AddOutline,
    parent: 'admin-templates-list',
    path: '/admin/templates/create'
  },
  'admin-templates-edit': {
    title: '编辑模板',
    icon: CreateOutline,
    parent: 'admin-templates-list',
    path: '/admin/templates/edit'
  },
  'admin-categories': {
    title: '分类管理',
    icon: LayersOutline,
    path: '/admin/categories'
  },
  'admin-categories-create': {
    title: '创建分类',
    icon: AddOutline,
    parent: 'admin-categories',
    path: '/admin/categories/create'
  },
  'admin-categories-edit': {
    title: '编辑分类',
    icon: CreateOutline,
    parent: 'admin-categories',
    path: '/admin/categories/edit'
  },
  'admin-languages': {
    title: '语言管理',
    icon: LanguageOutline,
    path: '/admin/languages'
  },
  'admin-languages-create': {
    title: '创建语言',
    icon: AddOutline,
    parent: 'admin-languages',
    path: '/admin/languages/create'
  },
  'admin-languages-edit': {
    title: '编辑语言',
    icon: CreateOutline,
    parent: 'admin-languages',
    path: '/admin/languages/edit'
  },
  'admin-analytics': {
    title: '统计分析',
    icon: StatsChartOutline,
    path: '/admin/analytics'
  },
  'admin-settings': {
    title: '系统设置',
    icon: SettingsOutline,
    path: '/admin/settings'
  }
}

const breadcrumbItems = computed(() => {
  const currentRouteName = route.name
  const items = []
  
  // 始终添加仪表盘作为首页
  items.push({
    title: '仪表盘',
    icon: GridOutline,
    path: '/admin',
    clickable: currentRouteName !== 'admin-dashboard'
  })

  if (currentRouteName && currentRouteName !== 'admin-dashboard') {
    const currentConfig = breadcrumbConfig[currentRouteName]
    
    if (currentConfig) {
      // 如果有父级，先添加父级
      if (currentConfig.parent) {
        const parentConfig = breadcrumbConfig[currentConfig.parent]
        if (parentConfig) {
          items.push({
            title: parentConfig.title,
            icon: parentConfig.icon,
            path: parentConfig.path,
            clickable: true
          })
        }
      }
      
      // 添加当前页面
      items.push({
        title: currentConfig.title,
        icon: currentConfig.icon,
        path: currentConfig.path,
        clickable: false
      })
    }
  }

  return items
})

function handleBreadcrumbClick(item) {
  if (item.clickable && item.path) {
    router.push(item.path)
  }
}
</script>

<style scoped>
.breadcrumb-container {
  padding: 16px 24px;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
}

.admin-breadcrumb {
  font-size: 14px;
}

:deep(.n-breadcrumb-item:not(:last-child) .n-breadcrumb-item__link) {
  color: #666;
  transition: color 0.2s ease;
}

:deep(.n-breadcrumb-item:not(:last-child) .n-breadcrumb-item__link:hover) {
  color: #18a058;
}

:deep(.n-breadcrumb-item:last-child .n-breadcrumb-item__link) {
  color: #333;
  font-weight: 500;
}

:deep(.n-breadcrumb-item__icon) {
  margin-right: 4px;
}
</style>