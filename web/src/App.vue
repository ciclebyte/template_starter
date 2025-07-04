<template>
  <n-message-provider>
    <MainLayout v-if="!isEditPage" />
    <router-view />
  </n-message-provider>
</template>

<script setup>
// App.vue 现在只作为路由入口
import { NMessageProvider } from 'naive-ui'
import MainLayout from './layouts/MainLayout.vue'
import { useRoute } from 'vue-router'
import { computed, watchEffect } from 'vue'
const route = useRoute()
const isEditPage = computed(() => route.path.startsWith('/templates/edit'))
</script>

<style>
/* 全局样式 */
</style>
