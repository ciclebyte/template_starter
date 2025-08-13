<template>
  <n-layout-header bordered class="admin-header">
    <div class="header-content">
      <!-- 左侧：侧边栏切换按钮 -->
      <div class="header-left">
        <n-button
          text
          @click="$emit('toggle-sidebar')"
          class="sidebar-toggle"
        >
          <template #icon>
            <n-icon size="18">
              <MenuOutline />
            </n-icon>
          </template>
        </n-button>
      </div>

      <!-- 右侧：用户信息和操作 -->
      <div class="header-right">
        <!-- 全屏切换 -->
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-button
              text
              @click="toggleFullscreen"
              class="header-action"
            >
              <template #icon>
                <n-icon size="18">
                  <ScanOutline />
                </n-icon>
              </template>
            </n-button>
          </template>
          全屏
        </n-tooltip>

        <!-- 设置 -->
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-button
              text
              @click="goSettings"
              class="header-action"
            >
              <template #icon>
                <n-icon size="18">
                  <SettingsOutline />
                </n-icon>
              </template>
            </n-button>
          </template>
          设置
        </n-tooltip>

        <!-- 返回前台 -->
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-button
              text
              @click="goFrontend"
              class="header-action"
            >
              <template #icon>
                <n-icon size="18">
                  <ExitOutline />
                </n-icon>
              </template>
            </n-button>
          </template>
          返回前台
        </n-tooltip>

        <!-- 用户头像和菜单 -->
        <n-dropdown
          trigger="hover"
          :options="userMenuOptions"
          @select="handleUserMenuSelect"
        >
          <div class="user-info">
            <n-avatar
              round
              size="small"
              src="https://naive-ui.oss-cn-beijing.aliyuncs.com/07akioni.jpeg"
            />
            <span class="username">管理员</span>
          </div>
        </n-dropdown>
      </div>
    </div>
  </n-layout-header>
</template>

<script setup>
import { h } from 'vue'
import { useRouter } from 'vue-router'
import {
  NLayoutHeader,
  NButton,
  NIcon,
  NTooltip,
  NDropdown,
  NAvatar
} from 'naive-ui'
import {
  MenuOutline,
  ScanOutline,
  SettingsOutline,
  ExitOutline,
  PersonOutline,
  LogOutOutline
} from '@vicons/ionicons5'

const props = defineProps({
  sidebarCollapsed: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['toggle-sidebar'])

const router = useRouter()

const userMenuOptions = [
  {
    label: '个人资料',
    key: 'profile',
    icon: () => h(NIcon, null, { default: () => h(PersonOutline) })
  },
  {
    type: 'divider'
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: () => h(NIcon, null, { default: () => h(LogOutOutline) })
  }
]

function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    if (document.exitFullscreen) {
      document.exitFullscreen()
    }
  }
}

function goSettings() {
  router.push({ name: 'admin-settings' })
}

function goFrontend() {
  router.push('/')
}

function handleUserMenuSelect(key) {
  if (key === 'profile') {
    // 处理个人资料
    console.log('Go to profile')
  } else if (key === 'logout') {
    // 处理退出登录
    router.push('/')
  }
}
</script>

<style scoped>
.admin-header {
  height: 56px;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
  position: sticky;
  top: 0;
  z-index: 99;
}

.header-content {
  height: 100%;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
}

.sidebar-toggle {
  padding: 8px;
  border-radius: 6px;
}

.sidebar-toggle:hover {
  background: #f5f5f5;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-action {
  padding: 8px;
  border-radius: 6px;
  color: #666;
}

.header-action:hover {
  background: #f5f5f5;
  color: #333;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-radius: 20px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.user-info:hover {
  background: #f5f5f5;
}

.username {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}
</style>