<template>
  <n-message-provider>
    <router-view />
  </n-message-provider>
</template>

<script setup>
// App.vue 现在只作为路由入口
import { NMessageProvider } from 'naive-ui'
</script>

<style>
/* 全局样式 */
</style>
