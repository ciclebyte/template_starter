<template>
  <div class="main-layout">
    <!-- 顶部导航栏 -->
    <header class="header">
      <div class="header-content">
        <div class="logo">
          <h2>Template Starter</h2>
        </div>
        <div class="nav-menu">
          <n-menu 
            mode="horizontal" 
            :options="menuOptions" 
            :value="activeKey"
            @update:value="handleMenuClick"
          />
        </div>
        <div class="search-box">
          <n-input placeholder="搜索模板..." />
        </div>
      </div>
    </header>

    <!-- 主要内容区域 -->
    <main class="main-content">
      <router-view />
    </main>

    <!-- 底部 -->
    <footer class="footer">
      <div class="footer-content">
        <p>&copy; 2024 Template Starter. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { h, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

// 菜单配置
const menuOptions = [
  {
    label: '首页',
    key: 'home'
  },
  {
    label: '模板',
    key: 'templates'
  }
]

const activeKey = ref(route.name || 'home')

const handleMenuClick = (key) => {
  if (key === 'home') {
    router.push('/')
  } else if (key === 'templates') {
    router.push('/templates')
  }
}
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: #f0f9f4;
  border-bottom: 1px solid #e6f7e6;
  padding: 0;
  width: 100%;
  flex-shrink: 0;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 100;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.logo h2 {
  margin: 0;
  color: #18a058;
  font-weight: 600;
  font-size: 20px;
}

.nav-menu {
  flex: 1;
  margin: 0 40px;
}

/* 自定义菜单样式 */
.nav-menu :deep(.n-menu-item) {
  font-weight: 500;
  font-size: 15px;
  color: #666;
  transition: color 0.2s ease;
  padding: 8px 20px;
  border-radius: 6px;
  margin: 0 4px;
  cursor: pointer;
}

.nav-menu :deep(.n-menu-item:hover) {
  color: #18a058;
  background: rgba(24, 160, 88, 0.06);
}

.nav-menu :deep(.n-menu-item--selected) {
  color: #18a058;
  background: rgba(24, 160, 88, 0.1);
  font-weight: 600;
}

.nav-menu :deep(.n-menu-item-content) {
  padding: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-menu :deep(.n-menu-item-content__text) {
  width: 100%;
  text-align: center;
}

.search-box {
  width: 300px;
}

.search-box :deep(.n-input) {
  border-radius: 6px;
}

.main-content {
  flex: 1;
  background: #f5f5f5;
  width: 100%;
  padding-top: 64px;
}

.footer {
  background: #fff;
  border-top: 1px solid #f0f0f0;
  width: 100%;
  flex-shrink: 0;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  text-align: center;
  color: #666;
}
</style> 