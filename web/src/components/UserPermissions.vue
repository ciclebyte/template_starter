<template>
  <div v-if="authStore.isLoggedIn && showPermissions" class="permissions-info">
    <n-card size="small" title="当前权限">
      <template #header-extra>
        <n-button text @click="showPermissions = false">
          <n-icon>
            <CloseOutline />
          </n-icon>
        </n-button>
      </template>
      
      <div class="permission-section">
        <h4>角色</h4>
        <n-space>
          <n-tag v-for="role in authStore.roles" :key="role" type="info">
            {{ getRoleName(role) }}
          </n-tag>
        </n-space>
      </div>
      
      <div class="permission-section">
        <h4>权限</h4>
        <n-space>
          <n-tag 
            v-for="permission in authStore.permissions" 
            :key="permission" 
            type="success"
            size="small"
          >
            {{ permission }}
          </n-tag>
        </n-space>
      </div>
    </n-card>
  </div>
  
  <n-button 
    v-if="authStore.isLoggedIn && !showPermissions" 
    class="show-permissions-btn"
    size="small"
    @click="showPermissions = true"
  >
    查看权限
  </n-button>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { CloseOutline } from '@vicons/ionicons5'

const authStore = useAuthStore()
const showPermissions = ref(false)

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
</script>

<style scoped>
.permissions-info {
  position: fixed;
  top: 80px;
  right: 20px;
  width: 300px;
  z-index: 1000;
}

.show-permissions-btn {
  position: fixed;
  top: 80px;
  right: 20px;
  z-index: 1000;
}

.permission-section {
  margin-bottom: 16px;
}

.permission-section h4 {
  margin: 0 0 8px 0;
  font-size: 14px;
  font-weight: 600;
}
</style>