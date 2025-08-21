<template>
  <div class="permission-demo">
    <n-card title="权限控制演示">
      <n-space vertical size="large">
        <!-- 基本权限控制 -->
        <div class="demo-section">
          <h3>基本权限控制</h3>
          <PermissionWrapper :permissions="['template:read']">
            <n-alert type="success">
              您有 template:read 权限，可以看到此内容
            </n-alert>
          </PermissionWrapper>
          
          <PermissionWrapper :permissions="['template:create']" show-fallback>
            <n-alert type="info">
              您有 template:create 权限，可以看到此内容
            </n-alert>
            <template #fallback>
              <n-alert type="warning">
                您没有 template:create 权限
              </n-alert>
            </template>
          </PermissionWrapper>
        </div>
        
        <!-- 角色控制 -->
        <div class="demo-section">
          <h3>角色控制</h3>
          <PermissionWrapper :roles="['super_admin', 'system_admin']" show-fallback>
            <n-alert type="success">
              您是管理员，可以看到此管理员内容
            </n-alert>
            <template #fallback>
              <n-alert type="error">
                只有管理员可以看到此内容
              </n-alert>
            </template>
          </PermissionWrapper>
        </div>
        
        <!-- 复合权限控制 -->
        <div class="demo-section">
          <h3>复合权限控制</h3>
          <PermissionWrapper 
            :permissions="['template:edit', 'template:delete']" 
            mode="any"
            show-fallback
          >
            <n-alert type="info">
              您可以编辑或删除模板
            </n-alert>
            <template #fallback>
              <n-alert type="warning">
                您没有编辑或删除模板的权限
              </n-alert>
            </template>
          </PermissionWrapper>
        </div>
        
        <!-- 按钮权限控制 -->
        <div class="demo-section">
          <h3>按钮权限控制</h3>
          <n-space>
            <PermissionWrapper :permissions="['template:create']">
              <n-button type="primary">创建模板</n-button>
            </PermissionWrapper>
            
            <PermissionWrapper :permissions="['template:edit']">
              <n-button type="info">编辑模板</n-button>
            </PermissionWrapper>
            
            <PermissionWrapper :permissions="['template:delete']">
              <n-button type="error">删除模板</n-button>
            </PermissionWrapper>
            
            <PermissionWrapper :roles="['super_admin']">
              <n-button type="warning">超级管理功能</n-button>
            </PermissionWrapper>
          </n-space>
        </div>
      </n-space>
    </n-card>
  </div>
</template>

<script setup>
import PermissionWrapper from '@/components/PermissionWrapper.vue'
</script>

<style scoped>
.permission-demo {
  padding: 20px;
}

.demo-section {
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.demo-section:last-child {
  border-bottom: none;
}

.demo-section h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
}
</style>