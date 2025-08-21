<template>
  <n-layout-header bordered class="n-layout-header header-no-padding">
    <div class="nav-bar-center">
      <div class="nav-bar-container">
        <div class="nav-left">
          <div class="brand" @click="goHome" style="cursor:pointer;">
            <div class="logo-icon">
              <svg width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect width="32" height="32" rx="6" fill="url(#brandGradient)" />
                <rect x="8" y="6" width="12" height="16" rx="1" fill="#ffffff" />
                <path d="M18 6 L18 10 L22 10 Z" fill="#e6f7ff" />
                <rect x="10" y="10" width="6" height="1" fill="#52c41a" />
                <rect x="10" y="12" width="4" height="1" fill="#1890ff" />
                <rect x="10" y="14" width="5" height="1" fill="#722ed1" />
                <circle cx="11" cy="17" r="0.5" fill="#ff4d4f" />
                <circle cx="13" cy="17" r="0.5" fill="#ff4d4f" />
                <rect x="14.5" y="16.5" width="2" height="1" fill="#ff4d4f" />
                <path d="M22 20 L26 24 L22 28" stroke="#52c41a" stroke-width="2" stroke-linecap="round"
                  stroke-linejoin="round" fill="none" />
                <defs>
                  <linearGradient id="brandGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" style="stop-color:#1890ff;stop-opacity:1" />
                    <stop offset="100%" style="stop-color:#18a058;stop-opacity:1" />
                  </linearGradient>
                </defs>
              </svg>
            </div>
            <div class="logo-text">
              <span class="logo-main">Template <span class="brand-accent">Starter</span></span>
              <span class="logo-shadow">Template <span class="brand-accent">Starter</span></span>
            </div>
          </div>
          <n-menu class="menu-center" mode="horizontal" :options="menuOptions" :value="activeKey"
            @update:value="handleUpdateValue" />
          <div class="search-inline" ref="searchRef">
            <n-input v-model:value="searchKeyword" round placeholder="搜索模板..." class="search-input" clearable
              @keyup.enter="handleSearch" @update:value="handleSearchChange">
              <template #prefix>
                <n-icon>
                  <SearchOutline />
                </n-icon>
              </template>
            </n-input>
          </div>
        </div>
        <div class="nav-right">
          <!-- 未登录状态 -->
          <div v-if="!authStore.isLoggedIn" class="auth-buttons">
            <n-button text @click="goLogin">登录</n-button>
            <n-button type="primary" round @click="goRegister">注册</n-button>
          </div>
          
          <!-- 已登录状态 -->
          <div v-else class="user-menu">
            <n-dropdown :options="userMenuOptions" @select="handleUserMenuSelect">
              <n-button text>
                <template #icon>
                  <n-icon>
                    <PersonOutline />
                  </n-icon>
                </template>
                {{ authStore.userDisplayName }}
              </n-button>
            </n-dropdown>
            
            <n-button v-if="authStore.isAdmin" type="primary" round @click="goTemplateEditor">
              <template #icon>
                <n-icon>
                  <SettingsOutline />
                </n-icon>
              </template>
              管理后台
            </n-button>
          </div>
        </div>
      </div>
    </div>
  </n-layout-header>
</template>

<script setup>
import { ref, computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NLayoutHeader, NMenu, NInput, NButton, NIcon, NDropdown } from 'naive-ui'
import { SearchOutline, SettingsOutline, PersonOutline, LogOutOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/stores/auth'
import { useMessage } from 'naive-ui'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const message = useMessage()

const menuOptions = [
  { label: '首页', key: 'home' },
  { label: '模板', key: 'templates' },
  { label: '权限演示', key: 'permission-demo' },
]

// 用户菜单选项
const userMenuOptions = [
  { label: '个人资料', key: 'profile', icon: () => h(PersonOutline) },
  { type: 'divider' },
  { label: '退出登录', key: 'logout', icon: () => h(LogOutOutline) },
]

const activeKey = ref(route.name || 'home')

function handleUpdateValue(key) {
  activeKey.value = key
  if (key === 'home') router.push('/')
  else if (key === 'templates') router.push('/templates')
  else if (key === 'permission-demo') router.push('/permission-demo')
}

function goHome() {
  router.push('/')
}

function goLogin() {
  router.push('/login')
}

function goRegister() {
  router.push('/register')
}

function goTemplateEditor() {
  router.push('/admin')
}

// 处理用户菜单选择
async function handleUserMenuSelect(key) {
  if (key === 'logout') {
    try {
      await authStore.doLogout()
      message.success('已退出登录')
      router.push('/')
    } catch (error) {
      message.error('退出登录失败')
    }
  } else if (key === 'profile') {
    router.push('/profile')
  }
}

// 搜索相关
const searchKeyword = ref('')
const searchRef = ref(null)

function handleSearchChange(value) {
  searchKeyword.value = value
}

function handleSearch() {
  if (searchKeyword.value.trim()) {
    router.push({
      path: '/templates',
      query: { search: searchKeyword.value.trim() }
    })
  }
}
</script>

<style scoped>
.header-no-padding {
  padding: 0 !important;
}

.nav-bar-center {
  width: 80vw;
  max-width: 1280px;
  margin: 0 auto;
  display: flex;
  justify-content: center;
}

.nav-bar-container {
  width: 100%;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  box-sizing: border-box;
  padding: 0 32px;
  position: relative;
}

.nav-left {
  display: flex;
  align-items: center;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.auth-buttons {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand {
  display: flex;
  align-items: center;
  margin-right: 36px;
  margin-left: 8px;
  position: relative;
  cursor: pointer;
  transition: all 0.3s ease;
}

.logo-icon {
  margin-right: 12px;
  display: flex;
  align-items: center;
  transition: all 0.3s ease;
}

.brand:hover .logo-icon {
  transform: scale(1.1) rotate(5deg);
}

.brand:hover {
  transform: scale(1.05) rotate(-1deg);
}

.brand:hover .logo-main {
  text-shadow: 0 0 20px rgba(24, 160, 88, 0.8);
  transform: translateY(-2px) scale(1.1);
}

.brand:hover .logo-shadow {
  opacity: 0.8;
  transform: translateY(4px) scale(1.1);
}

.logo-text {
  position: relative;
  display: flex;
  align-items: center;
}

.logo-main {
  font-size: 1.7rem;
  font-weight: 800;
  letter-spacing: 1.5px;
  color: #333;
  font-family: 'Fira Code', 'Lato', 'Segoe UI', 'Arial', sans-serif;
  background: linear-gradient(90deg, #18a058 0%, #2196F3 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  position: relative;
  z-index: 2;
  animation: float 3s ease-in-out infinite;
  transition: all 0.3s ease;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  white-space: nowrap;
}

.logo-shadow {
  font-size: 1.7rem;
  font-weight: 800;
  letter-spacing: 1.5px;
  color: rgba(51, 51, 51, 0.2);
  font-family: 'Fira Code', 'Lato', 'Segoe UI', 'Arial', sans-serif;
  position: absolute;
  top: 2px;
  left: 0;
  right: 0;
  z-index: 1;
  animation: float-shadow 3s ease-in-out infinite;
  transition: all 0.3s ease;
  filter: blur(1px);
  white-space: nowrap;
}

.brand-accent {
  color: #18a058;
  -webkit-text-fill-color: #18a058;
  background: none;
  font-weight: 900;
}

@keyframes float {

  0%,
  100% {
    transform: translateY(0px);
  }

  50% {
    transform: translateY(-3px);
  }
}

@keyframes float-shadow {

  0%,
  100% {
    transform: translateY(2px);
    opacity: 0.2;
  }

  50% {
    transform: translateY(5px);
    opacity: 0.3;
  }
}

.brand::before {
  content: '';
  position: absolute;
  top: -20px;
  left: -20px;
  right: -20px;
  bottom: -20px;
  background: radial-gradient(circle, rgba(24, 160, 88, 0.1) 0%, transparent 70%);
  opacity: 0;
  transition: opacity 0.3s ease;
  pointer-events: none;
}

.brand:hover::before {
  opacity: 1;
}

.menu-center {
  min-width: 200px;
}

.search-inline {
  margin-left: 32px;
  display: flex;
  align-items: center;
  position: relative;
}

.search-input {
  width: 260px;
  background: #f5f6fa;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 18px;
}

:deep(.menu-center .n-menu--horizontal .n-menu__content) {
  justify-content: flex-start !important;
}

:deep(.n-layout-header__content) {
  padding: 0 !important;
}

.n-layout-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 200;
  background: #fff;
  box-shadow: 0 2px 8px 0 rgba(60, 60, 60, 0.04);
}
</style>