<template>
  <n-config-provider>
    <n-notification-provider>
      <n-message-provider>
        <n-dialog-provider>
          <router-view v-slot="{ Component }">
            <component :is="Component" />
          </router-view>
        </n-dialog-provider>
      </n-message-provider>
    </n-notification-provider>
  </n-config-provider>
</template>

<script setup>
// App.vue 现在只作为路由入口
import { NMessageProvider, NNotificationProvider, NConfigProvider, NDialogProvider } from 'naive-ui'
</script>

<style>
/* 全局样式 */
</style>
