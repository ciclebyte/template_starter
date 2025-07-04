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
    key: 'home',
    icon: () => h('span', '🏠')
  },
  {
    label: '模板',
    key: 'templates',
    icon: () => h('span', '📁')
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
}

.nav-menu {
  flex: 1;
  margin: 0 40px;
}

.search-box {
  width: 300px;
}

.main-content {
  flex: 1;
  background: #f5f5f5;
  width: 100%;
  min-height: 0;
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