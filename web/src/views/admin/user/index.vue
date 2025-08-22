<template>
  <div class="user-management">
    <n-card title="用户管理" :bordered="false">
      <!-- 工具栏 -->
      <div class="toolbar">
        <div class="search-bar">
          <n-input
            v-model:value="searchForm.search"
            placeholder="搜索用户名、邮箱或昵称"
            clearable
            style="width: 300px"
            @keyup.enter="loadUsers"
          >
            <template #prefix>
              <n-icon>
                <SearchOutline />
              </n-icon>
            </template>
          </n-input>
          <n-select
            v-model:value="searchForm.status"
            placeholder="选择状态"
            clearable
            style="width: 120px; margin-left: 12px"
            :options="statusOptions"
            @update:value="loadUsers"
          />
          <n-select
            v-model:value="searchForm.role"
            placeholder="选择角色"
            clearable
            style="width: 150px; margin-left: 12px"
            :options="roleOptions"
            @update:value="loadUsers"
          />
          <n-button @click="loadUsers" style="margin-left: 12px">
            <template #icon>
              <n-icon>
                <SearchOutline />
              </n-icon>
            </template>
            搜索
          </n-button>
        </div>
        <div class="actions">
          <n-button type="primary" @click="showUserModal()">
            <template #icon>
              <n-icon>
                <AddOutline />
              </n-icon>
            </template>
            新增用户
          </n-button>
        </div>
      </div>

      <!-- 用户表格 -->
      <n-data-table
        :columns="userColumns"
        :data="userList"
        :loading="userLoading"
        :pagination="userPagination"
        :row-key="row => row.id"
        @update:page="handleUserPageChange"
      />
    </n-card>

    <!-- 用户编辑弹窗 -->
    <n-modal v-model:show="userModalVisible" preset="dialog" title="用户信息" style="width: 600px">
      <template #header>
        <div>{{ userForm.id ? '编辑用户' : '新增用户' }}</div>
      </template>
      
      <n-form
        ref="userFormRef"
        :model="userForm"
        :rules="userRules"
        label-placement="left"
        label-width="80px"
        style="margin-top: 16px"
      >
        <n-form-item label="用户名" path="username">
          <n-input v-model:value="userForm.username" placeholder="请输入用户名" />
        </n-form-item>

        <n-form-item label="邮箱" path="email">
          <n-input v-model:value="userForm.email" placeholder="请输入邮箱地址" />
        </n-form-item>

        <n-form-item v-if="!userForm.id" label="密码" path="password">
          <n-input v-model:value="userForm.password" type="password" placeholder="请输入密码" />
        </n-form-item>

        <n-form-item label="昵称" path="nickname">
          <n-input v-model:value="userForm.nickname" placeholder="请输入昵称" />
        </n-form-item>

        <n-form-item label="手机号" path="phone">
          <n-input v-model:value="userForm.phone" placeholder="请输入手机号" />
        </n-form-item>

        <n-form-item label="状态" path="status">
          <n-select
            v-model:value="userForm.status"
            placeholder="请选择状态"
            :options="statusOptions"
          />
        </n-form-item>

        <n-form-item label="角色" path="roles">
          <n-select
            v-model:value="userForm.roles"
            placeholder="请选择角色"
            multiple
            :options="allRoleOptions"
          />
        </n-form-item>
      </n-form>

      <template #action>
        <n-button @click="userModalVisible = false">取消</n-button>
        <n-button type="primary" @click="saveUser" :loading="userSaving">
          保存
        </n-button>
      </template>
    </n-modal>

    <!-- 重置密码弹窗 -->
    <n-modal v-model:show="passwordModalVisible" preset="dialog" title="重置密码">
      <n-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-placement="left"
        label-width="80px"
        style="margin-top: 16px"
      >
        <n-form-item label="新密码" path="newPassword">
          <n-input v-model:value="passwordForm.newPassword" type="password" placeholder="请输入新密码" />
        </n-form-item>
      </n-form>

      <template #action>
        <n-button @click="passwordModalVisible = false">取消</n-button>
        <n-button type="primary" @click="resetPassword" :loading="passwordResetting">
          确定
        </n-button>
      </template>
    </n-modal>

    <!-- 角色分配弹窗 -->
    <n-modal v-model:show="roleModalVisible" preset="dialog" title="分配角色" style="width: 500px">
      <template #header>
        <div>为用户 "{{ currentUser?.username }}" 分配角色</div>
      </template>

      <n-form
        ref="roleFormRef"
        :model="roleForm"
        label-placement="left"
        label-width="80px"
        style="margin-top: 16px"
      >
        <n-form-item label="角色" path="roles">
          <n-select
            v-model:value="roleForm.roles"
            placeholder="请选择角色"
            multiple
            :options="allRoleOptions"
          />
        </n-form-item>

        <n-form-item label="过期时间" path="expiresAt">
          <n-date-picker
            v-model:value="roleForm.expiresAt"
            type="datetime"
            placeholder="选择过期时间（可选）"
            clearable
            style="width: 100%"
          />
        </n-form-item>
      </n-form>

      <template #action>
        <n-button @click="roleModalVisible = false">取消</n-button>
        <n-button type="primary" @click="assignRoles" :loading="roleAssigning">
          确定
        </n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, h } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import { SearchOutline, AddOutline, CreateOutline, TrashBinOutline, KeyOutline, PersonOutline } from '@vicons/ionicons5'
import {
  getUsers,
  createUser,
  updateUser,
  deleteUser,
  resetUserPassword,
  updateUserStatus,
  assignUserRoles
} from '@/api/user'
import { getRoles } from '@/api/permission'

const message = useMessage()
const dialog = useDialog()

// ============================================================================
// 数据定义
// ============================================================================

const userList = ref([])
const userLoading = ref(false)
const userModalVisible = ref(false)
const userSaving = ref(false)
const userFormRef = ref(null)

const passwordModalVisible = ref(false)
const passwordResetting = ref(false)
const passwordFormRef = ref(null)

const roleModalVisible = ref(false)
const roleAssigning = ref(false)
const roleFormRef = ref(null)
const currentUser = ref(null)

const allRoleOptions = ref([])
const roleOptions = ref([])

// 搜索表单
const searchForm = reactive({
  search: '',
  status: null,
  role: '',
  page: 1,
  size: 20
})

// 用户表单
const userForm = reactive({
  id: null,
  username: '',
  email: '',
  password: '',
  nickname: '',
  phone: '',
  status: 1,
  roles: []
})

// 密码表单
const passwordForm = reactive({
  userId: null,
  newPassword: ''
})

// 角色表单
const roleForm = reactive({
  userId: null,
  roles: [],
  expiresAt: null
})

// 分页
const userPagination = reactive({
  page: 1,
  pageSize: 20,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page) => {
    userPagination.page = page
    searchForm.page = page
    loadUsers()
  },
  onUpdatePageSize: (pageSize) => {
    userPagination.pageSize = pageSize
    userPagination.page = 1
    searchForm.size = pageSize
    searchForm.page = 1
    loadUsers()
  }
})

// 状态选项
const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

// 表格列定义
const userColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '用户名', key: 'username' },
  { title: '邮箱', key: 'email' },
  { title: '昵称', key: 'nickname' },
  { title: '手机号', key: 'phone' },
  { 
    title: '状态', 
    key: 'status',
    width: 100,
    render(row) {
      return h('n-tag', {
        type: row.status === 1 ? 'success' : 'default'
      }, {
        default: () => row.status === 1 ? '启用' : '禁用'
      })
    }
  },
  {
    title: '角色',
    key: 'roles',
    render(row) {
      if (!row.roles || row.roles.length === 0) {
        return h('span', { style: { color: '#999' } }, '无角色')
      }
      return h('div', {}, row.roles.map(role => 
        h('n-tag', { 
          size: 'small',
          style: { marginRight: '4px' }
        }, { default: () => role.name })
      ))
    }
  },
  { title: '最后登录', key: 'lastLoginAt', width: 180 },
  { title: '创建时间', key: 'createdAt', width: 180 },
  {
    title: '操作',
    key: 'actions',
    width: 240,
    render(row) {
      return [
        h('n-button', {
          size: 'small',
          type: 'primary',
          style: { marginRight: '8px' },
          onClick: () => showUserModal(row)
        }, { default: () => '编辑' }),
        h('n-button', {
          size: 'small',
          style: { marginRight: '8px' },
          onClick: () => showPasswordModal(row)
        }, { 
          icon: () => h('n-icon', {}, { default: () => h(KeyOutline) }),
          default: () => '重置密码'
        }),
        h('n-button', {
          size: 'small',
          style: { marginRight: '8px' },
          onClick: () => showRoleModal(row)
        }, { 
          icon: () => h('n-icon', {}, { default: () => h(PersonOutline) }),
          default: () => '分配角色'
        }),
        h('n-button', {
          size: 'small',
          type: 'error',
          onClick: () => handleDeleteUser(row)
        }, { default: () => '删除' })
      ]
    }
  }
]

// 表单验证规则
const userRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }
  ],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

const passwordRules = {
  newPassword: [{ required: true, message: '请输入新密码', trigger: 'blur' }]
}

// ============================================================================
// 方法定义
// ============================================================================

// 加载用户列表
async function loadUsers() {
  try {
    userLoading.value = true
    const response = await getUsers(searchForm)
    
    if (response.data.code === 0) {
      userList.value = response.data.data.list || []
      userPagination.itemCount = response.data.data.total || 0
    }
  } catch (error) {
    message.error('加载用户列表失败')
  } finally {
    userLoading.value = false
  }
}

// 加载角色选项
async function loadRoleOptions() {
  try {
    const response = await getRoles({ page: 1, size: 1000 })
    
    if (response.data.code === 0) {
      const roles = response.data.data.list || []
      allRoleOptions.value = roles.map(role => ({
        label: role.name,
        value: role.id
      }))
      roleOptions.value = roles.map(role => ({
        label: role.name,
        value: role.code
      }))
    }
  } catch (error) {
    console.error('加载角色选项失败:', error)
  }
}

// 显示用户编辑弹窗
function showUserModal(user = null) {
  if (user) {
    Object.assign(userForm, {
      ...user,
      roles: user.roles ? user.roles.map(r => r.id) : []
    })
  } else {
    Object.assign(userForm, {
      id: null,
      username: '',
      email: '',
      password: '',
      nickname: '',
      phone: '',
      status: 1,
      roles: []
    })
  }
  userModalVisible.value = true
}

// 保存用户
async function saveUser() {
  try {
    await userFormRef.value?.validate()
    userSaving.value = true

    const data = { ...userForm }
    delete data.id

    if (userForm.id) {
      await updateUser(userForm.id, data)
      message.success('用户更新成功')
    } else {
      await createUser(data)
      message.success('用户创建成功')
    }

    userModalVisible.value = false
    loadUsers()
  } catch (error) {
    if (error.message) {
      message.error(error.message)
    }
  } finally {
    userSaving.value = false
  }
}

// 显示密码重置弹窗
function showPasswordModal(user) {
  passwordForm.userId = user.id
  passwordForm.newPassword = ''
  passwordModalVisible.value = true
}

// 重置密码
async function resetPassword() {
  try {
    await passwordFormRef.value?.validate()
    passwordResetting.value = true

    await resetUserPassword(passwordForm.userId, {
      newPassword: passwordForm.newPassword
    })
    
    message.success('密码重置成功')
    passwordModalVisible.value = false
  } catch (error) {
    message.error('密码重置失败')
  } finally {
    passwordResetting.value = false
  }
}

// 显示角色分配弹窗
function showRoleModal(user) {
  currentUser.value = user
  roleForm.userId = user.id
  roleForm.roles = user.roles ? user.roles.map(r => r.id) : []
  roleForm.expiresAt = null
  roleModalVisible.value = true
}

// 分配角色
async function assignRoles() {
  try {
    roleAssigning.value = true

    const data = {
      roles: roleForm.roles,
      expiresAt: roleForm.expiresAt ? new Date(roleForm.expiresAt).toISOString() : ''
    }

    await assignUserRoles(roleForm.userId, data)
    
    message.success('角色分配成功')
    roleModalVisible.value = false
    loadUsers()
  } catch (error) {
    message.error('角色分配失败')
  } finally {
    roleAssigning.value = false
  }
}

// 删除用户
function handleDeleteUser(user) {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除用户"${user.username}"吗？此操作不可恢复。`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteUser(user.id)
        message.success('用户删除成功')
        loadUsers()
      } catch (error) {
        message.error('用户删除失败')
      }
    }
  })
}

// 用户分页变化
function handleUserPageChange(page) {
  userPagination.page = page
  searchForm.page = page
  loadUsers()
}

// 初始化
onMounted(() => {
  loadUsers()
  loadRoleOptions()
})
</script>

<style scoped>
.user-management {
  padding: 16px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.search-bar {
  display: flex;
  align-items: center;
}
</style>