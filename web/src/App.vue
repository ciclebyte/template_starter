<template>
  <n-notification-provider>
    <n-message-provider>
      <router-view />
    </n-message-provider>
  </n-notification-provider>
</template>

<script setup>
// App.vue 现在只作为路由入口
import { NMessageProvider, NNotificationProvider } from 'naive-ui'
</script>

<style>
/* 全局样式 */
</style>
