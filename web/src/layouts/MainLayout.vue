<template>
  <div class="main-layout">
    <!-- 顶部导航栏 -->
    <header class="header">
      <div class="logo">
        <h2>Template Starter</h2>
      </div>
      <div class="header-content">
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

<style>
/* 全局样式 - 确保始终显示滚动条 */
html {
  overflow-y: scroll;
}
</style>

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
  display: flex;
  align-items: center;
  height: 64px;
}

.logo {
  position: absolute;
  left: 20px;
  z-index: 101;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  width: 100%;
}

.logo h2 {
  margin: 0;
  color: #18a058;
  font-weight: 600;
  font-size: 20px;
}

.nav-menu {
  flex: 1;
  display: flex;
  justify-content: flex-start;
}

/* 自定义菜单样式 */
.nav-menu :deep(.n-menu-item) {
  font-weight: 500;
  font-size: 15px;
  color: #666;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  padding: 16px 32px;
  border-radius: 0;
  margin: 0 8px;
  cursor: pointer;
  min-height: 56px;
  display: flex;
  align-items: center;
  position: relative;
  background: transparent;
}

.nav-menu :deep(.n-menu-item::after) {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  width: 80%;
  height: 3px;
  background: rgba(24, 160, 88, 0.1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform: translateX(-50%);
  border-radius: 2px;
}

.nav-menu :deep(.n-menu-item:hover) {
  color: #18a058;
  background: transparent;
  transform: translateY(-1px);
}

.nav-menu :deep(.n-menu-item:hover::after) {
  background: rgba(24, 160, 88, 0.4);
}

.nav-menu :deep(.n-menu-item--selected) {
  color: #18a058;
  background: transparent;
  font-weight: 600;
}

.nav-menu :deep(.n-menu-item--selected::after) {
  background: #18a058 !important;
  opacity: 1 !important;
}

.nav-menu :deep(.n-menu-item-content) {
  padding: 0;
  width: 100%;
  height: 100%;
  min-height: 56px;
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
  display: flex;
  justify-content: center;
}

.main-content > * {
  max-width: 1200px;
  width: 100%;
  padding: 0 20px;
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