<template>
  <div class="profile-container">
    <n-card title="个人资料" class="profile-card">
      <div v-if="authStore.user" class="user-info">
        <n-descriptions :column="1" label-placement="left">
          <n-descriptions-item label="用户名">
            {{ authStore.user.username }}
          </n-descriptions-item>
          <n-descriptions-item label="邮箱">
            {{ authStore.user.email }}
          </n-descriptions-item>
          <n-descriptions-item label="昵称">
            {{ authStore.user.nickname || '未设置' }}
          </n-descriptions-item>
          <n-descriptions-item label="状态">
            <n-tag :type="authStore.user.status === 1 ? 'success' : 'error'">
              {{ authStore.user.status === 1 ? '正常' : '禁用' }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="角色">
            <n-space>
              <n-tag 
                v-for="role in authStore.roles" 
                :key="role" 
                type="info"
              >
                {{ getRoleName(role) }}
              </n-tag>
            </n-space>
          </n-descriptions-item>
          <n-descriptions-item label="最后登录">
            {{ authStore.user.last_login_at || '未知' }}
          </n-descriptions-item>
        </n-descriptions>
        
        <div class="actions">
          <n-button type="primary" @click="editProfile">
            编辑资料
          </n-button>
          <n-button @click="changePassword">
            修改密码
          </n-button>
        </div>
      </div>
      
      <div v-else class="loading">
        <n-spin size="large" />
      </div>
    </n-card>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { useMessage } from 'naive-ui'

const authStore = useAuthStore()
const message = useMessage()

const roleNames = {
  'super_admin': '超级管理员',
  'system_admin': '系统管理员',
  'org_admin': '组织管理员',
  'template_admin': '模板管理员',
  'developer': '开发者',
  'user': '普通用户',
  'guest': '访客'
}

function getRoleName(role) {
  return roleNames[role] || role
}

function editProfile() {
  message.info('编辑资料功能开发中')
}

function changePassword() {
  message.info('修改密码功能开发中')
}
</script>

<style scoped>
.profile-container {
  padding: 40px;
  max-width: 800px;
  margin: 0 auto;
}

.profile-card {
  margin: 0 auto;
}

.user-info {
  padding: 20px 0;
}

.actions {
  margin-top: 30px;
  display: flex;
  gap: 16px;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
}
</style>