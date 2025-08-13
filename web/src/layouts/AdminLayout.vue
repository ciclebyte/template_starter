<template>
  <n-layout has-sider class="admin-layout">
    <!-- 侧边栏 -->
    <AdminSidebar :collapsed="sidebarCollapsed" @update:collapsed="sidebarCollapsed = $event" />

    <!-- 主内容区域 -->
    <n-layout class="main-layout" :class="{ collapsed: sidebarCollapsed }">
      <!-- 顶部导航栏 -->
      <AdminHeader :sidebar-collapsed="sidebarCollapsed" @toggle-sidebar="sidebarCollapsed = !sidebarCollapsed" />

      <!-- 面包屑 -->
      <AdminBreadcrumb />

      <!-- 页面内容 -->
      <n-layout-content class="admin-content">
        <div class="content-wrapper">
          <router-view v-slot="{ Component }">
            <transition mode="out-in" @before-enter="beforeEnter" @enter="enter" @leave="leave">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup>
import { ref } from 'vue'
import { NLayout, NLayoutContent } from 'naive-ui'
import AdminSidebar from '@/components/admin/AdminSidebar.vue'
import AdminHeader from '@/components/admin/AdminHeader.vue'
import AdminBreadcrumb from '@/components/admin/AdminBreadcrumb.vue'

const sidebarCollapsed = ref(false)

function beforeEnter(el) {
  el.style.opacity = 0
  el.style.transform = 'translateY(20px)'
}

function enter(el, done) {
  el.offsetHeight // trigger reflow
  el.style.transition = 'opacity 0.3s ease, transform 0.3s ease'
  el.style.opacity = 1
  el.style.transform = 'translateY(0)'
  setTimeout(done, 300)
}

function leave(el, done) {
  el.style.transition = 'opacity 0.2s ease, transform 0.2s ease'
  el.style.opacity = 0
  el.style.transform = 'translateY(-20px)'
  setTimeout(done, 200)
}
</script>

<style scoped>
.admin-layout {
  height: 100vh;
}

.main-layout {
  background: #f5f5f5;
  margin-left: 240px;
  transition: margin-left 0.3s ease;
}

.main-layout.collapsed {
  margin-left: 64px;
}

.admin-content {
  padding: 20px;
  background: #f5f5f5;
}

.content-wrapper {
  /* 移除所有padding，让内容自然流动 */
}
</style>