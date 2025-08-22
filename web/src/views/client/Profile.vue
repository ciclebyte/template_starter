<template>
  <div class="user-center-container">
    <div class="user-center-wrapper">
      <!-- 页面标题 -->
      <div class="page-header">
        <h1 class="page-title">用户中心</h1>
      </div>

      <div class="user-center-content">
        <!-- 左侧导航 -->
        <div class="sidebar">
          <div class="user-card">
            <n-avatar 
              round 
              size="medium"
              :src="userInfo?.avatar"
            >
              <template #fallback>
                <n-icon size="20">
                  <PersonOutline />
                </n-icon>
              </template>
            </n-avatar>
            <div class="user-brief">
              <div class="user-name">{{ userInfo?.nickname || userInfo?.username }}</div>
              <div class="user-email">{{ userInfo?.email }}</div>
            </div>
          </div>

          <n-menu
            :value="activeTab"
            :options="menuOptions"
            @update:value="handleTabChange"
            accordion
          />
        </div>

        <!-- 右侧内容区 -->
        <div class="main-content">
          <div class="content-card">
            <!-- 基本资料 -->
            <div v-if="activeTab === 'profile'" class="tab-content">
              <div class="content-header">
                <h2>基本资料</h2>
                <p>管理您的个人基本信息</p>
              </div>

              <div class="avatar-section">
                <n-avatar 
                  round 
                  size="large"
                  :src="userInfo?.avatar"
                >
                  <template #fallback>
                    <n-icon size="32">
                      <PersonOutline />
                    </n-icon>
                  </template>
                </n-avatar>
                <div class="avatar-info">
                  <div class="avatar-tips">
                    <p>支持 JPG、PNG 格式，文件大小不超过 2MB</p>
                    <n-button size="small" @click="handleAvatarUpload">
                      更换头像
                    </n-button>
                  </div>
                </div>
              </div>

              <n-form
                ref="profileFormRef"
                :model="profileForm"
                :rules="profileRules"
                label-placement="left"
                label-width="100px"
              >
                <n-form-item label="用户名" path="username">
                  <n-input
                    v-model:value="profileForm.username"
                    disabled
                    placeholder="用户名不可修改"
                  />
                </n-form-item>

                <n-form-item label="昵称" path="nickname">
                  <n-input
                    v-model:value="profileForm.nickname"
                    placeholder="请输入昵称"
                  />
                </n-form-item>

                <n-form-item label="邮箱" path="email">
                  <n-input
                    v-model:value="profileForm.email"
                    disabled
                    placeholder="邮箱不可修改"
                  />
                </n-form-item>

                <n-form-item label="手机号码" path="phone">
                  <n-input
                    v-model:value="profileForm.phone"
                    placeholder="请输入手机号码"
                  />
                </n-form-item>

                <n-form-item>
                  <n-button type="primary" @click="handleSaveProfile" :loading="loading">
                    保存更改
                  </n-button>
                </n-form-item>
              </n-form>
            </div>

            <!-- 安全设置 -->
            <div v-if="activeTab === 'security'" class="tab-content">
              <div class="content-header">
                <h2>安全设置</h2>
                <p>修改密码和其他安全相关设置</p>
              </div>

              <n-form
                ref="securityFormRef"
                :model="securityForm"
                :rules="securityRules"
                label-placement="left"
                label-width="100px"
              >
                <n-form-item label="当前密码" path="currentPassword">
                  <n-input
                    v-model:value="securityForm.currentPassword"
                    type="password"
                    placeholder="请输入当前密码"
                    show-password-on="click"
                  />
                </n-form-item>

                <n-form-item label="新密码" path="newPassword">
                  <n-input
                    v-model:value="securityForm.newPassword"
                    type="password"
                    placeholder="请输入新密码"
                    show-password-on="click"
                  />
                </n-form-item>

                <n-form-item label="确认密码" path="confirmPassword">
                  <n-input
                    v-model:value="securityForm.confirmPassword"
                    type="password"
                    placeholder="请再次输入新密码"
                    show-password-on="click"
                  />
                </n-form-item>

                <n-form-item>
                  <n-button type="primary" @click="handleChangePassword" :loading="loading">
                    修改密码
                  </n-button>
                </n-form-item>
              </n-form>
            </div>

            <!-- 我的角色 -->
            <div v-if="activeTab === 'roles'" class="tab-content">
              <div class="content-header">
                <h2>我的角色</h2>
                <p>查看您拥有的角色和权限信息</p>
              </div>

              <div class="roles-info">
                <div class="roles-overview">
                  <h3>角色概览</h3>
                  <div class="role-cards">
                    <div 
                      v-for="role in authStore.roles" 
                      :key="role"
                      class="role-card"
                    >
                      <div class="role-header">
                        <n-icon size="24" class="role-icon">
                          <ShieldCheckmarkOutline />
                        </n-icon>
                        <div class="role-info">
                          <div class="role-name">{{ getRoleDisplayName(role) }}</div>
                          <div class="role-code">{{ role }}</div>
                        </div>
                      </div>
                      <div class="role-description">
                        {{ getRoleDescription(role) }}
                      </div>
                    </div>
                  </div>
                </div>

                <div class="permissions-info" v-if="authStore.permissions.length > 0">
                  <h3>我的权限</h3>
                  <div class="permission-grid">
                    <n-tag
                      v-for="permission in authStore.permissions"
                      :key="permission"
                      type="info"
                      size="medium"
                      class="permission-tag"
                    >
                      {{ getPermissionDisplayName(permission) }}
                    </n-tag>
                  </div>
                </div>
              </div>
            </div>

            <!-- 账户统计 -->
            <div v-if="activeTab === 'stats'" class="tab-content">
              <div class="content-header">
                <h2>账户统计</h2>
                <p>查看您的账户使用情况和统计信息</p>
              </div>

              <div class="stats-grid">
                <div class="stat-card">
                  <div class="stat-icon">👤</div>
                  <div class="stat-info">
                    <div class="stat-value">{{ stats.loginCount || 0 }}</div>
                    <div class="stat-label">登录次数</div>
                  </div>
                </div>

                <div class="stat-card">
                  <div class="stat-icon">🕒</div>
                  <div class="stat-info">
                    <div class="stat-value">{{ stats.lastLoginAt || '从未' }}</div>
                    <div class="stat-label">最后登录</div>
                  </div>
                </div>

                <div class="stat-card">
                  <div class="stat-icon">📅</div>
                  <div class="stat-info">
                    <div class="stat-value">{{ stats.createdAt || '未知' }}</div>
                    <div class="stat-label">注册时间</div>
                  </div>
                </div>

                <div class="stat-card">
                  <div class="stat-icon">🏷️</div>
                  <div class="stat-info">
                    <div class="stat-value">{{ authStore.roles.length }}</div>
                    <div class="stat-label">角色数量</div>
                  </div>
                </div>
              </div>

              <div class="roles-section">
                <h3>我的角色</h3>
                <div class="roles-list">
                  <n-tag 
                    v-for="role in authStore.roles" 
                    :key="role"
                    type="success"
                    size="medium"
                  >
                    {{ getRoleDisplayName(role) }}
                  </n-tag>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useMessage } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { PersonOutline, LockClosedOutline, BarChartOutline, ShieldCheckmarkOutline } from '@vicons/ionicons5'
import { NIcon } from 'naive-ui'
import { 
  getProfile, 
  updateProfile, 
  changePassword, 
  updateEmail, 
  uploadAvatar, 
  getSecurityInfo, 
  getLoginHistory 
} from '@/api/profile'

const message = useMessage()
const authStore = useAuthStore()

const loading = ref(false)
const activeTab = ref('profile')
const profileFormRef = ref(null)
const securityFormRef = ref(null)

// 用户信息
const userInfo = computed(() => authStore.user)

// 菜单选项
const menuOptions = [
  {
    label: '基本资料',
    key: 'profile',
    icon: () => h(NIcon, { size: 16 }, { default: () => h(PersonOutline) })
  },
  {
    label: '安全设置',
    key: 'security',
    icon: () => h(NIcon, { size: 16 }, { default: () => h(LockClosedOutline) })
  },
  {
    label: '我的角色',
    key: 'roles',
    icon: () => h(NIcon, { size: 16 }, { default: () => h(ShieldCheckmarkOutline) })
  },
  {
    label: '账户统计',
    key: 'stats',
    icon: () => h(NIcon, { size: 16 }, { default: () => h(BarChartOutline) })
  }
]

// 基本资料表单
const profileForm = reactive({
  username: '',
  nickname: '',
  email: '',
  phone: ''
})

// 安全设置表单
const securityForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 统计信息
const stats = reactive({
  loginCount: 0,
  lastLoginAt: '',
  createdAt: ''
})

// 基本资料验证规则
const profileRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 1, max: 50, message: '昵称长度在 1 到 50 个字符', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ]
}

// 安全设置验证规则
const securityRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 64, message: '密码长度在 6 到 64 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value) => {
        if (value !== securityForm.newPassword) {
          return new Error('两次输入的密码不一致')
        }
        return true
      },
      trigger: 'blur'
    }
  ]
}

// 角色显示名称映射
const getRoleDisplayName = (role) => {
  const roleMap = {
    'super_admin': '超级管理员',
    'system_admin': '系统管理员',
    'org_admin': '组织管理员',
    'user': '普通用户',
    'guest': '访客'
  }
  return roleMap[role] || role
}

// 角色描述映射
const getRoleDescription = (role) => {
  const roleDescMap = {
    'super_admin': '拥有系统的最高权限，可以管理所有功能和用户',
    'system_admin': '系统管理员，可以管理系统配置和用户',
    'org_admin': '组织管理员，可以管理本组织的用户和资源',
    'user': '普通用户，可以使用基本功能',
    'guest': '访客用户，仅可查看公开内容'
  }
  return roleDescMap[role] || '暂无描述'
}

// 权限显示名称映射
const getPermissionDisplayName = (permission) => {
  const permMap = {
    'template:read': '查看模板',
    'template:use': '使用模板',
    'template:create': '创建模板',
    'template:edit': '编辑模板',
    'template:delete': '删除模板',
    'template:share': '分享模板',
    'template:manage': '管理模板',
    'category:read': '查看分类',
    'category:manage': '管理分类',
    'language:read': '查看语言',
    'language:manage': '管理语言',
    'tag:read': '查看标签',
    'tag:manage': '管理标签',
    'user:read': '查看用户',
    'user:manage': '管理用户',
    'user:assign_role': '分配角色',
    'role:read': '查看角色',
    'role:manage': '管理角色',
    'role:assign_permission': '分配权限',
    'system:config': '系统配置',
    'system:audit': '审计日志',
    'system:analytics': '统计分析'
  }
  return permMap[permission] || permission
}

// Tab切换
const handleTabChange = (key) => {
  activeTab.value = key
}

// 初始化表单数据
const initFormData = async () => {
  try {
    const response = await getProfile()
    if (response.data.code === 0) {
      const profile = response.data.data.profile
      profileForm.username = profile.username || ''
      profileForm.nickname = profile.nickname || ''
      profileForm.email = profile.email || ''
      profileForm.phone = profile.phone || ''
      
      // 统计信息
      stats.loginCount = 0 // 这个需要从安全信息中获取
      stats.lastLoginAt = profile.lastLoginAt ? 
        new Date(profile.lastLoginAt).toLocaleString() : '从未'
      stats.createdAt = profile.createdAt ? 
        new Date(profile.createdAt).toLocaleDateString() : '未知'
        
      // 加载安全信息
      loadSecurityInfo()
    }
  } catch (error) {
    console.error('获取个人资料失败:', error)
    // 如果获取失败，使用authStore中的信息作为后备
    if (userInfo.value) {
      profileForm.username = userInfo.value.username || ''
      profileForm.nickname = userInfo.value.nickname || ''
      profileForm.email = userInfo.value.email || ''
      profileForm.phone = userInfo.value.phone || ''
    }
  }
}

// 加载安全信息
const loadSecurityInfo = async () => {
  try {
    const response = await getSecurityInfo()
    if (response.data.code === 0) {
      const securityInfo = response.data.data.securityInfo
      stats.loginCount = securityInfo.loginCount || 0
    }
  } catch (error) {
    console.error('获取安全信息失败:', error)
  }
}

// 处理头像上传
const handleAvatarUpload = () => {
  message.info('头像上传功能开发中...')
}

// 保存基本资料
const handleSaveProfile = async () => {
  try {
    await profileFormRef.value?.validate()
    loading.value = true
    
    await updateProfile({
      nickname: profileForm.nickname,
      phone: profileForm.phone,
      avatar: userInfo.value?.avatar || ''
    })
    
    message.success('个人资料更新成功')
    
    // 刷新用户信息
    await authStore.fetchUserInfo()
    
  } catch (error) {
    console.error('保存个人资料失败:', error)
    message.error('保存失败，请重试')
  } finally {
    loading.value = false
  }
}

// 修改密码
const handleChangePassword = async () => {
  try {
    await securityFormRef.value?.validate()
    loading.value = true
    
    await changePassword({
      oldPassword: securityForm.currentPassword,
      newPassword: securityForm.newPassword
    })
    
    message.success('密码修改成功')
    
    // 清空表单
    securityForm.currentPassword = ''
    securityForm.newPassword = ''
    securityForm.confirmPassword = ''
    
  } catch (error) {
    console.error('修改密码失败:', error)
    message.error('修改密码失败，请检查原密码是否正确')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  initFormData()
})
</script>

<style scoped>
.user-center-container {
  min-height: 100vh;
  background: #f5f6fa;
  padding: 24px;
}

.user-center-wrapper {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.user-center-content {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 24px;
}

/* 左侧导航 */
.sidebar {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  height: fit-content;
}

.user-card {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.user-brief {
  flex: 1;
}

.user-name {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 4px 0;
}

.user-email {
  font-size: 12px;
  color: #6c757d;
  margin: 0;
}

/* 右侧内容区 */
.main-content {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.content-card {
  padding: 32px;
}

.tab-content {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.content-header {
  margin-bottom: 32px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e9ecef;
}

.content-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
}

.content-header p {
  font-size: 14px;
  color: #6c757d;
  margin: 0;
}

/* 头像区域 */
.avatar-section {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 32px;
  padding: 24px;
  background: #f8f9fa;
  border-radius: 8px;
}

.avatar-info {
  flex: 1;
}

.avatar-tips p {
  font-size: 14px;
  color: #6c757d;
  margin: 0 0 12px 0;
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.stat-icon {
  font-size: 24px;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 18px;
  font-weight: 700;
  color: #2c3e50;
  margin: 0 0 4px 0;
}

.stat-label {
  font-size: 14px;
  color: #6c757d;
  margin: 0;
}

/* 角色区域 */
.roles-section {
  margin-top: 32px;
}

.roles-section h3 {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 16px 0;
}

.roles-list {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

/* 我的角色页面样式 */
.roles-info h3 {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 16px 0;
}

.role-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
  margin-bottom: 32px;
}

.role-card {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 20px;
  transition: all 0.2s ease;
}

.role-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.role-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.role-icon {
  color: #18a058;
}

.role-info {
  flex: 1;
}

.role-name {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 4px 0;
}

.role-code {
  font-size: 12px;
  color: #6c757d;
  font-family: 'Monaco', 'Consolas', monospace;
  background: #e9ecef;
  padding: 2px 6px;
  border-radius: 4px;
  display: inline-block;
}

.role-description {
  font-size: 14px;
  color: #6c757d;
  line-height: 1.5;
}

.permissions-info {
  margin-top: 32px;
}

.permission-grid {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.permission-tag {
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-center-container {
    padding: 16px;
  }
  
  .user-center-content {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .sidebar {
    padding: 20px;
  }
  
  .content-card {
    padding: 24px;
  }
  
  .avatar-section {
    flex-direction: column;
    text-align: center;
    gap: 16px;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .page-title {
    font-size: 24px;
  }
  
  .user-card {
    flex-direction: column;
    text-align: center;
    gap: 8px;
  }
  
  .content-header h2 {
    font-size: 20px;
  }
}

/* 菜单样式优化 */
:deep(.n-menu .n-menu-item) {
  border-radius: 8px;
  margin-bottom: 4px;
}

:deep(.n-menu .n-menu-item:hover) {
  background-color: #f8f9fa;
}

:deep(.n-menu .n-menu-item.n-menu-item--selected) {
  background-color: #18a058;
  color: white;
}

:deep(.n-menu .n-menu-item.n-menu-item--selected .n-icon) {
  color: white;
}
</style>